package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/reports"
)

type (
	generator struct {
		storage reports.Storage
		writer  ReportWriter
		logger  zerolog.Logger
	}
)

const dataKeyCategory = "/category"

const sourceData = "Data"
const sourceFileID = "FileID"
const sourceVersion = "Version"
const sourcePatientID = "PatientID"
const sourceCreatedAt = "CreatedAt"
const sourceUpdatedAt = "UpdatedAt"

// Generate generates report
func (g *generator) Generate(ctx context.Context, writer ReportWriter, reportSpec ReportSpec, createdAtStart *strfmt.DateTime, createdAtEnd *strfmt.DateTime) error {
	files, err := g.storage.Find("", map[string]string{dataKeyCategory: reportSpec.FileCategory}, createdAtStart, createdAtEnd)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch files")
	}

	err = writer.Write(reportSpec.Columns)
	if err != nil {
		return errors.Wrapf(err, "failed to write report header")
	}

	for _, file := range *files {
		var dataMap map[string]interface{}
		if err := json.Unmarshal([]byte(file.Data), &dataMap); err != nil {
			return errors.Wrapf(err, "failed to unmarshal file data")
		}

		row := []string{}
		for _, column := range reportSpec.Columns {
			spec, ok := reportSpec.ColumnsSpecs[column]
			if !ok {
				return errors.Errorf("could not find a spec for column '%s'", column)
			}
			switch spec.Source {
			case sourceData:
				_, value := g.generateValueFromData(spec, &file, &dataMap, "")
				row = append(row, strings.TrimSpace(value))
			case sourceFileID:
				row = append(row, file.FileID)
			case sourceVersion:
				row = append(row, file.Version)
			case sourcePatientID:
				row = append(row, file.PatientID)
			case sourceCreatedAt:
				row = append(row, file.CreatedAt.String())
			case sourceUpdatedAt:
				row = append(row, file.UpdatedAt.String())
			}
		}
		err = writer.Write(row)
		if err != nil {
			return errors.Wrapf(err, "failed to write report row")
		}
	}

	return nil
}

func (g *generator) generateValueFromData(spec ValueSpec, file *reports.File, data *map[string]interface{}, prefix string) (found bool, value string) {
	found = false
	switch spec.Type {
	case "multipleValues":
		values := []interface{}{}
		for _, fieldSpec := range spec.Properties {
			found, value = g.generateValueFromData(fieldSpec, file, data, prefix)
			values = append(values, value)
		}
		return found, fmt.Sprintf(spec.Format, values...)
	case "array":
		values := []interface{}{}
		for i := spec.IncludeItems.Start; true; i++ {
			elementFound := false
			elementValues := []interface{}{}

			for _, fieldSpec := range spec.Properties {
				valueFound, value := g.generateValueFromData(fieldSpec, file, data, fmt.Sprintf("%s%s:%d", prefix, spec.EhrPath, i))
				if valueFound {
					elementFound = true
				}
				elementValues = append(elementValues, value)
			}
			if elementFound {
				values = append(values, fmt.Sprintf(spec.Format, elementValues...))
				found = true
			}

			if elementFound == false || (spec.IncludeItems.End != -1 && i == spec.IncludeItems.End) {
				break
			}
		}

		value := ""
		for _, v := range values {
			if value == "" {
				value = v.(string)
			} else {
				value = fmt.Sprintf("%s, %s", value, v)
			}
		}

		return found, value

	default:
		return g.getData(data, fmt.Sprintf("%s%s", prefix, spec.EhrPath))
	}

	return found, ""
}

func (g *generator) getData(data *map[string]interface{}, fullEhrPath string) (found bool, value string) {
	if val, ok := (*data)[fullEhrPath]; ok {
		switch val.(type) {
		case string:
			return true, val.(string)
		case int:
			return true, strconv.Itoa(val.(int))
		case float32:
			return true, strconv.FormatFloat(float64(val.(float32)), 'E', -1, 64)
		case float64:
			return true, strconv.FormatFloat(val.(float64), 'E', -1, 64)
		case bool:
			return true, strconv.FormatBool(val.(bool))
		default:
			return false, ""
		}
	}

	return false, ""
}

// New initializes a new instance of generator
func New(storage reports.Storage, logger zerolog.Logger) (*generator, error) {
	s := &generator{
		storage: storage,
		logger:  logger.With().Str("component", "reports/generator").Logger(),
	}

	return s, nil
}
