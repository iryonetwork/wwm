package generator

//go:generate ../../bin/mockgen.sh reports/generator ReportWriter $GOFILE
import (
	"context"

	"github.com/go-openapi/strfmt"
)

type (
	Generator interface {
		Generate(ctx context.Context, writer ReportWriter, reportSpec ReportSpec, createdAtStart *strfmt.DateTime, createdAtEnd *strfmt.DateTime) error
	}

	ReportWriter interface {
		// Write writes single report row
		Write(row []string) error
		// Flush flushes to report any data that has been buffered
		Flush()
		// Error returns any that has occured during previoud WriteHeader, Write or Flush
		Error() error
	}

	// ReportSpec represents specficiation for generating report
	ReportSpec struct {
		Type         string               `json:"type"`
		IncludeAll   bool                 `json:"includeAll"`
		FileCategory string               `json:"fileCategory"`
		Columns      []string             `json:"columns"`
		ColumnsSpecs map[string]ValueSpec `json:"columnsSpecs"`
	}

	// Column represents specification for report column
	ValueSpec struct {
		Type         string             `json:"type"`
		EhrPath      string             `json:"ehrPath"`
		Source       string             `json:"source"`
		Format       string             `json:"format"`
		Properties   []ValueSpec        `json:"properties"`
		IncludeItems IncludeItemsStruct `json:"includeItems"`
	}

	IncludeItemsStruct struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
)
