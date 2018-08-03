package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/reports"
	generatorMock "github.com/iryonetwork/wwm/reports/generator/mock"
	storageMock "github.com/iryonetwork/wwm/reports/mock"
)

var (
	noErrors   = false
	withErrors = true

	testSpec1 = ReportSpec{
		Type:         "testReport",
		FileCategory: "openehr::111|test|",
		Columns:      []string{"file ID", "version", "patient ID", "createdAt", "updatedAt", "multiple values from data", "twice nested array first item", "twice nested array"},
		ColumnsSpecs: map[string]ValueSpec{
			"file ID": ValueSpec{
				Type:   "value",
				Source: "FileID",
			},
			"version": ValueSpec{
				Type:   "value",
				Source: "Version",
			},
			"patient ID": ValueSpec{
				Type:   "value",
				Source: "PatientID",
			},
			"createdAt": ValueSpec{
				Type:   "value",
				Source: "CreatedAt",
			},
			"updatedAt": ValueSpec{
				Type:   "value",
				Source: "UpdatedAt",
			},
			"multiple values from data": ValueSpec{
				Type:   "multipleValues",
				Source: "Data",
				Format: "%s %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "value",
						EhrPath: "/userID",
					},
					ValueSpec{
						Type:    "value",
						EhrPath: "/userName",
					},
				},
			},
			"twice nested array first item": ValueSpec{
				Type:    "array",
				Source:  "Data",
				EhrPath: "/arrayLevel1",
				IncludeItems: IncludeItemsStruct{
					Start: 0,
					End:   0,
				},
				Format: "%s - %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "array",
						Source:  "Data",
						EhrPath: "/arrayLevel2",
						IncludeItems: IncludeItemsStruct{
							Start: 0,
							End:   -1,
						},
						Format: "%s/%s",
						Properties: []ValueSpec{
							ValueSpec{
								Type:    "value",
								EhrPath: "/nestedArrayItemID",
							},
							ValueSpec{
								Type:    "value",
								EhrPath: "/nestedArrayItemName",
							},
						},
					},
					ValueSpec{
						Type:    "value",
						EhrPath: "/additionalData",
					},
				},
			},
			"twice nested array": ValueSpec{
				Type:    "array",
				Source:  "Data",
				EhrPath: "/arrayLevel1",
				IncludeItems: IncludeItemsStruct{
					Start: 1,
					End:   -1,
				},
				Format: "%s - %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "array",
						Source:  "Data",
						EhrPath: "/arrayLevel2",
						IncludeItems: IncludeItemsStruct{
							Start: 0,
							End:   -1,
						},
						Format: "%s/%s",
						Properties: []ValueSpec{
							ValueSpec{
								Type:    "value",
								EhrPath: "/nestedArrayItemID",
							},
							ValueSpec{
								Type:    "value",
								EhrPath: "/nestedArrayItemName",
							},
						},
					},
					ValueSpec{
						Type:    "value",
						EhrPath: "/additionalData",
					},
				},
			},
		},
	}

	testSpec2 = ReportSpec{
		Type:             "testReport",
		FileCategory:     "openehr::111|test|",
		GroupByPatientID: true,
		Columns:          []string{"patient ID", "createdAt", "multiple values from data"},
		ColumnsSpecs: map[string]ValueSpec{
			"patient ID": ValueSpec{
				Type:   "value",
				Source: "PatientID",
			},
			"createdAt": ValueSpec{
				Type:   "value",
				Source: "CreatedAt",
			},
			"multiple values from data": ValueSpec{
				Type:   "multipleValues",
				Source: "Data",
				Format: "%s %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "value",
						EhrPath: "/valueFromFile1",
					},
					ValueSpec{
						Type:    "value",
						EhrPath: "/valueFromFile2",
					},
				},
			},
		},
	}

	time1, _ = strfmt.ParseDateTime("2018-07-27T13:55:59.123Z")
	time2, _ = strfmt.ParseDateTime("2018-07-29T13:55:59.123Z")

	fileNoData = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{}",
	}
	fileNoDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "", "", ""}

	fileAllData = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{\"/userID\": \"ID\", \"/userName\": \"username\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemID\": \"ID0:0\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemName\": \"Name0:0\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemID\": \"ID0:1\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemName\": \"Name0:1\", \"/arrayLevel1:0/additionalData\": \"Data0\", \"/arrayLevel1:1/arrayLevel2:0/nestedArrayItemID\": \"ID1:0\", \"/arrayLevel1:1/arrayLevel2:0/nestedArrayItemName\": \"Name1:0\", \"/arrayLevel1:1/arrayLevel2:1/nestedArrayItemID\": \"ID1:1\", \"/arrayLevel1:1/arrayLevel2:1/nestedArrayItemName\": \"Name1:1\", \"/arrayLevel1:1/additionalData\": \"Data1\", \"/arrayLevel1:2/arrayLevel2:0/nestedArrayItemID\": \"ID2:0\", \"/arrayLevel1:2/arrayLevel2:0/nestedArrayItemName\": \"Name2:0\", \"/arrayLevel1:2/arrayLevel2:1/nestedArrayItemID\": \"ID2:1\", \"/arrayLevel1:2/arrayLevel2:1/nestedArrayItemName\": \"Name2:1\", \"/arrayLevel1:2/additionalData\": \"Data2\"}",
	}
	fileAllDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "ID username", "ID0:0/Name0:0, ID0:1/Name0:1 - Data0", "ID1:0/Name1:0, ID1:1/Name1:1 - Data1, ID2:0/Name2:0, ID2:1/Name2:1 - Data2"}

	fileMissingData = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{\"/userID\": \"ID\",  \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemID\": \"ID0:0\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemName\": \"Name0:0\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemID\": \"ID0:1\", \"/arrayLevel1:0/additionalData\": \"Data0\"}",
	}
	fileMissingDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "ID", "ID0:0/Name0:0, ID0:1/ - Data0", ""}

	invalidJSONDataFile = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{this is invalid json}",
	}

	file1_groupByPatientID = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{\"/valueFromFile1\": \"1\"}",
	}
	file2_groupByPatientID = reports.File{
		FileID:    "file_id_2",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time2,
		UpdatedAt: time2,
		Data:      "{\"/valueFromFile2\": \"2\"}",
	}
	groupedByPatientIDReportRow = []string{"patient_1", "2018-07-27T13:55:59.123Z, 2018-07-29T13:55:59.123Z", "1 2"}
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		description    string
		reportSpec     ReportSpec
		createdAtStart *strfmt.DateTime
		createdAtEnd   *strfmt.DateTime
		calls          func(*storageMock.MockStorage, *generatorMock.MockReportWriter)
		errorExpected  bool
		exactError     error
	}{
		{
			"Only empty file returned",
			testSpec1,
			nil,
			nil,
			func(storage *storageMock.MockStorage, reportWriter *generatorMock.MockReportWriter) {
				gomock.InOrder(
					storage.EXPECT().Find("", map[string]string{"/category": "openehr::111|test|"}, nil, nil).Return(
						&[]reports.File{fileNoData}, nil,
					),
					reportWriter.EXPECT().Write(testSpec1.Columns).Return(nil).Times(1),
					reportWriter.EXPECT().Write(fileNoDataReportRow).Return(nil).Times(1),
				)
			},
			noErrors,
			nil,
		},
		{
			"Full file and partial file returned, filter by date",
			testSpec1,
			&time1,
			&time2,
			func(storage *storageMock.MockStorage, reportWriter *generatorMock.MockReportWriter) {
				gomock.InOrder(
					storage.EXPECT().Find("", map[string]string{"/category": "openehr::111|test|"}, &time1, &time2).Return(
						&[]reports.File{fileAllData, fileMissingData}, nil,
					),
					reportWriter.EXPECT().Write(testSpec1.Columns).Return(nil).Times(1),
					reportWriter.EXPECT().Write(fileAllDataReportRow).Return(nil).Times(1),
					reportWriter.EXPECT().Write(fileMissingDataReportRow).Return(nil).Times(1),
				)
			},
			noErrors,
			nil,
		},
		{
			"Invalid JSON Data file returned",
			testSpec1,
			&time1,
			&time2,
			func(storage *storageMock.MockStorage, reportWriter *generatorMock.MockReportWriter) {
				gomock.InOrder(
					storage.EXPECT().Find("", map[string]string{"/category": "openehr::111|test|"}, &time1, &time2).Return(
						&[]reports.File{invalidJSONDataFile}, nil,
					),
					reportWriter.EXPECT().Write(testSpec1.Columns).Return(nil).Times(1),
				)
			},
			withErrors,
			nil,
		},
		{
			"Write failed",
			testSpec1,
			&time1,
			&time2,
			func(storage *storageMock.MockStorage, reportWriter *generatorMock.MockReportWriter) {
				gomock.InOrder(
					storage.EXPECT().Find("", map[string]string{"/category": "openehr::111|test|"}, &time1, &time2).Return(
						&[]reports.File{fileAllData, fileMissingData}, nil,
					),
					reportWriter.EXPECT().Write(testSpec1.Columns).Return(nil).Times(1),
					reportWriter.EXPECT().Write(fileAllDataReportRow).Return(fmt.Errorf("error")).Times(1),
				)
			},
			withErrors,
			nil,
		},
		{
			"Test grouping by patient ID",
			testSpec2,
			nil,
			nil,
			func(storage *storageMock.MockStorage, reportWriter *generatorMock.MockReportWriter) {
				gomock.InOrder(
					storage.EXPECT().Find("", map[string]string{"/category": "openehr::111|test|"}, nil, nil).Return(
						&[]reports.File{file1_groupByPatientID, file2_groupByPatientID}, nil,
					),
					reportWriter.EXPECT().Write(testSpec2.Columns).Return(nil).Times(1),
					reportWriter.EXPECT().Write(groupedByPatientIDReportRow).Return(nil).Times(1),
				)
			},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storageMock, writerMock, c := getTestService(t, test.reportSpec)
			defer c()

			// setup mock calls
			test.calls(storageMock, writerMock)

			// call genereate
			err := svc.Generate(context.Background(), writerMock, test.reportSpec, test.createdAtStart, test.createdAtEnd)

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func getTestService(t *testing.T, reportSpec ReportSpec) (*generator, *storageMock.MockStorage, *generatorMock.MockReportWriter, func()) {
	// setup mocks
	mockCtrl := gomock.NewController(t)
	storage := storageMock.NewMockStorage(mockCtrl)
	writer := generatorMock.NewMockReportWriter(mockCtrl)

	g := &generator{
		storage: storage,
		logger:  zerolog.New(os.Stdout),
	}

	cleanup := func() {
		mockCtrl.Finish()
	}

	return g, storage, writer, cleanup
}

func toJSON(in interface{}) string {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(in)
	return buf.String()
}
