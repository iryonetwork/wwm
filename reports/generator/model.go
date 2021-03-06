package generator

//go:generate ../../bin/mockgen.sh reports/generator ReportWriter $GOFILE
import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
)

type (
	Generator interface {
		Generate(ctx context.Context, writer ReportWriter, reportSpec ReportSpec, createdAtStart *strfmt.DateTime, createdAtEnd *strfmt.DateTime) (bool, error)
	}

	ReportWriter interface {
		// Write writes single report row
		Write(row []string) error
		// Flush flushes to report any data that has been buffered
		Flush()
		// Error returns any that has occured during previous Write or Flush
		Error() error
	}

	// ReportSpec represents specficiation for generating report
	ReportSpec struct {
		Type             string               `json:"type"`
		GroupByPatientID bool                 `json:"groupByPatientID"`
		FileCategory     string               `json:"fileCategory"`
		Columns          []string             `json:"columns"`
		ColumnsSpecs     map[string]ValueSpec `json:"columnsSpecs"`
	}

	// ValueSpec represents specification for report column
	ValueSpec struct {
		Type            string             `json:"type"`
		MetaField       string             `json:"metaField"`
		EhrPath         string             `json:"ehrPath"`
		TimestampFormat string             `json:"timestampFormat"`
		Format          string             `json:"format"`
		Unit            string             `json:"unit"`
		Properties      []ValueSpec        `json:"properties"`
		IncludeItems    IncludeItemsStruct `json:"includeItems"`
	}

	IncludeItemsStruct struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
)

const TYPE_FILE_META = "fileMeta"
const TYPE_VALUE = "value"
const FIXED_VALUE = "fixedValue"
const TYPE_BOOLEAN = "boolean"
const TYPE_QUANTITY = "quantity"
const TYPE_INTEGER = "integer"
const TYPE_ARRAY = "array"
const TYPE_DATETIME = "datetime"
const TYPE_CODE = "code"

const META_FIELD_FILE_ID = "fileID"
const META_FIELD_VERSION = "version"
const META_FIELD_PATIENT_ID = "patientID"
const META_FIELD_CREATED_AT = "createdAt"
const META_FIELD_UPDATED_AT = "updatedAt"

const TIMESTAMP_FORMAT_DATETIME = "datetime"
const TIMESTAMP_FORMAT_DATE = "date"
const TIMESTAMP_FORMAT_TIME = "time"

var TIMESTAMP_FORMAT_LAYOUTS = map[string]string{
	TIMESTAMP_FORMAT_DATETIME: time.RFC3339,
	TIMESTAMP_FORMAT_DATE:     "2006-01-02",
	TIMESTAMP_FORMAT_TIME:     "15:04:05",
}
