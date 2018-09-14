package xlsxWriter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type xlsxWriter struct {
	filename    string
	sheet       string
	currentLine int
	xlsx        *excelize.File
	err         error
}

func (w *xlsxWriter) Write(row []string) error {
	w.currentLine++
	w.xlsx.SetSheetRow(w.sheet, fmt.Sprintf("A%d", w.currentLine), &row)

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
	w.xlsx.SetCellStyle(w.sheet, "A1", "ZZ1", boldStyle)

	// save
	w.err = w.xlsx.SaveAs(fmt.Sprintf("./%s", w.filename))
}

func (w *xlsxWriter) Error() error {
	return w.err
}

func (w *xlsxWriter) Filename() string {
	return w.filename
}

func New(reportType string) (*xlsxWriter, error) {
	filename := fmt.Sprintf("%s-%s.xlsx", reportType, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))

	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(reportType)
	xlsx.SetActiveSheet(index)
	xlsx.DeleteSheet("Sheet1")

	err := xlsx.SaveAs(fmt.Sprintf("./%s", filename))

	if err != nil {
		return nil, err
	}

	return &xlsxWriter{
		filename:    fmt.Sprintf("%s-%s.xlsx", reportType, strconv.FormatInt(time.Now().UTC().UnixNano(), 10)),
		sheet:       reportType,
		currentLine: 0,
		xlsx:        xlsx,
	}, nil
}
