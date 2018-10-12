package xlsxWriter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/reports/generator"
)

type xlsxWriter struct {
	spec          generator.ReportSpec
	filename      string
	sheet         string
	xlsx          *excelize.File
	currentRow    int
	currentColumn string
	err           error
	logger        zerolog.Logger
}

const EXCEL_DATETIME_LAYOUT = "2006-01-02 15:04:05"

var TIMESTAMP_FORMAT_NUMBER_FORMATS = map[string]int{
	generator.TIMESTAMP_FORMAT_DATETIME: 22,
	generator.TIMESTAMP_FORMAT_DATE:     14,
	generator.TIMESTAMP_FORMAT_TIME:     21,
}

func (w *xlsxWriter) Write(row []string) error {
	w.nextRow()

	for i, value := range row {
		cell := w.nextCell()
		switch w.spec.ColumnsSpecs[w.spec.Columns[i]].Type {
		case generator.TYPE_DATETIME:
			timestampFormat := w.spec.ColumnsSpecs[w.spec.Columns[i]].TimestampFormat
			if timestampFormat == "" {
				timestampFormat = generator.TIMESTAMP_FORMAT_DATETIME
			}

			switch timestampFormat {
			case generator.TIMESTAMP_FORMAT_TIME:
				w.xlsx.SetCellFormula(w.sheet, cell, fmt.Sprintf("=TIMEVALUE(\"%s\")", value))
			case generator.TIMESTAMP_FORMAT_DATE:
				w.xlsx.SetCellFormula(w.sheet, cell, fmt.Sprintf("=DATEVALUE(\"%s\")", value))
			default:
				t, err := time.Parse(generator.TIMESTAMP_FORMAT_LAYOUTS[timestampFormat], value)
				if err != nil {
					w.logger.Error().Err(err).Msg("failed to parse datetime")
					w.err = err
					return err
				}
				value = t.Format(EXCEL_DATETIME_LAYOUT)
				w.xlsx.SetCellFormula(w.sheet, cell, fmt.Sprintf("=DATEVALUE(\"%s\") + TIMEVALUE(\"%s\")", value, value))
			}

			style, err := w.xlsx.NewStyle(fmt.Sprintf("{\"number_format\": %d}", TIMESTAMP_FORMAT_NUMBER_FORMATS[timestampFormat]))
			if err != nil {
				w.err = err
				w.logger.Error().Err(err).Msg("failed to create number formatting style")
				return err
			}
			w.xlsx.SetCellStyle(w.sheet, cell, cell, style)
		default:
			w.xlsx.SetCellValue(w.sheet, cell, value)

		}
	}

	return nil
}

func (w *xlsxWriter) Flush() {
	// set wide col width to fit at least UUID (there's no autofit support in lib yet)
	w.xlsx.SetColWidth(w.sheet, "A", "ZZ", 34)

	// freeze first row (header)
	w.xlsx.SetPanes(w.sheet, `{"freeze":true,"split":false,"x_split":0,"y_split":1,"top_left_cell":"A2","active_pane":"bottomLeft"}`)

	// set bold to first row (header)
	boldStyle, err := w.xlsx.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		w.err = err
		return
	}
	w.xlsx.SetCellStyle(w.sheet, "A1", fmt.Sprintf("%s1", w.currentColumn), boldStyle)

	// save
	w.err = w.xlsx.SaveAs(fmt.Sprintf("./%s", w.filename))
}

func (w *xlsxWriter) Error() error {
	return w.err
}

func (w *xlsxWriter) Filename() string {
	return w.filename
}

func New(spec generator.ReportSpec, logger zerolog.Logger) (*xlsxWriter, error) {
	filename := fmt.Sprintf("%s-%s.xlsx", spec.Type, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(spec.Type)
	xlsx.SetActiveSheet(index)
	xlsx.DeleteSheet("Sheet1")

	err := xlsx.SaveAs(fmt.Sprintf("./%s", filename))

	if err != nil {
		return nil, err
	}

	// write header with column names
	w := &xlsxWriter{
		spec:          spec,
		filename:      filename,
		sheet:         spec.Type,
		xlsx:          xlsx,
		currentRow:    0,
		currentColumn: "",
		logger:        logger.With().Str("component", "reports/generator/xlsxWriter").Logger(),
	}

	err = w.writeHeader()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *xlsxWriter) writeHeader() error {
	w.nextRow()
	w.xlsx.SetSheetRow(w.sheet, fmt.Sprintf("A%d", w.currentRow), &w.spec.Columns)

	return nil
}

func (w *xlsxWriter) nextCell() string {
	if w.currentColumn == "" {
		w.currentColumn = "A"
	} else if w.currentColumn[len(w.currentColumn)-1] == 'Z' {
		w.currentColumn = w.currentColumn[:len(w.currentColumn)-1] + "AA"
	} else {
		w.currentColumn = w.currentColumn[:len(w.currentColumn)-1] + string(w.currentColumn[len(w.currentColumn)-1]+1)
	}

	return fmt.Sprintf("%s%d", w.currentColumn, w.currentRow)
}

func (w *xlsxWriter) nextRow() {
	w.currentRow++
	w.currentColumn = ""
}
