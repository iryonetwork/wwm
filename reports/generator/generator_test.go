package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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
		Columns:      []string{"file ID", "version", "patient ID", "createdAt", "updatedAt", "quantity", "code", "twice nested array first item", "twice nested array"},
		ColumnsSpecs: map[string]ValueSpec{
			"file ID": ValueSpec{
				Type:      "fileMeta",
				MetaField: "fileID",
			},
			"version": ValueSpec{
				Type:      "fileMeta",
				MetaField: "version",
			},
			"patient ID": ValueSpec{
				Type:      "fileMeta",
				MetaField: "patientID",
			},
			"createdAt": ValueSpec{
				Type:      "fileMeta",
				MetaField: "createdAt",
			},
			"updatedAt": ValueSpec{
				Type:      "fileMeta",
				MetaField: "updatedAt",
			},
			"quantity": ValueSpec{
				Type:    "quantity",
				Unit:    "unit",
				EhrPath: "/quantityValue",
			},
			"code": ValueSpec{
				Type:    "code",
				EhrPath: "/codeValue",
			},
			"twice nested array first item": ValueSpec{
				Type:    "array",
				EhrPath: "/arrayLevel1",
				IncludeItems: IncludeItemsStruct{
					Start: 0,
					End:   0,
				},
				Format: "%s - %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "array",
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
				EhrPath: "/arrayLevel1",
				IncludeItems: IncludeItemsStruct{
					Start: 1,
					End:   -1,
				},
				Format: "%s - %s",
				Properties: []ValueSpec{
					ValueSpec{
						Type:    "array",
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
		Columns:          []string{"patient ID", "createdAt", "updatedAt", "valueFromFile1", "valueFromFile2", "valueInBothFiles1", "valueInBothFiles2"},
		ColumnsSpecs: map[string]ValueSpec{
			"patient ID": ValueSpec{
				Type:      "fileMeta",
				MetaField: "patientID",
			},
			"createdAt": ValueSpec{
				Type:      "fileMeta",
				MetaField: "createdAt",
			},
			"updatedAt": ValueSpec{
				Type:      "fileMeta",
				MetaField: "updatedAt",
			},
			"valueFromFile1": ValueSpec{
				Type:    "value",
				EhrPath: "/valueFromFile1",
			},
			"valueFromFile2": ValueSpec{
				Type:    "value",
				EhrPath: "/valueFromFile2",
			},
			"valueInBothFiles1": ValueSpec{
				Type:    "value",
				EhrPath: "/valueInBothFiles1",
			},
			"valueInBothFiles2": ValueSpec{
				Type:    "value",
				EhrPath: "/valueInBothFiles2",
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
	fileNoDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "", "", "", ""}

	fileAllData = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{\"/quantityValue\": \"value,unit\", \"/codeValue\": \"category::code|codeTitle|\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemID\": \"ID0:0\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemName\": \"Name0:0\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemID\": \"ID0:1\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemName\": \"Name0:1\", \"/arrayLevel1:0/additionalData\": \"Data0\", \"/arrayLevel1:1/arrayLevel2:0/nestedArrayItemID\": \"ID1:0\", \"/arrayLevel1:1/arrayLevel2:0/nestedArrayItemName\": \"Name1:0\", \"/arrayLevel1:1/arrayLevel2:1/nestedArrayItemID\": \"ID1:1\", \"/arrayLevel1:1/arrayLevel2:1/nestedArrayItemName\": \"Name1:1\", \"/arrayLevel1:1/additionalData\": \"Data1\", \"/arrayLevel1:2/arrayLevel2:0/nestedArrayItemID\": \"ID2:0\", \"/arrayLevel1:2/arrayLevel2:0/nestedArrayItemName\": \"Name2:0\", \"/arrayLevel1:2/arrayLevel2:1/nestedArrayItemID\": \"ID2:1\", \"/arrayLevel1:2/arrayLevel2:1/nestedArrayItemName\": \"Name2:1\", \"/arrayLevel1:2/additionalData\": \"Data2\"}",
	}
	fileAllDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "value unit", "codeTitle", "ID0:0/Name0:0, ID0:1/Name0:1 - Data0", "ID1:0/Name1:0, ID1:1/Name1:1 - Data1, ID2:0/Name2:0, ID2:1/Name2:1 - Data2"}

	fileMissingData = reports.File{
		FileID:    "file_id_1",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      "{\"/codeValue\": \"invalidCodeValue\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemID\": \"ID0:0\", \"/arrayLevel1:0/arrayLevel2:0/nestedArrayItemName\": \"Name0:0\", \"/arrayLevel1:0/arrayLevel2:1/nestedArrayItemID\": \"ID0:1\", \"/arrayLevel1:0/additionalData\": \"Data0\"}",
	}
	fileMissingDataReportRow = []string{"file_id_1", "version_1", "patient_1", "2018-07-27T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "", "invalidCodeValue", "ID0:0/Name0:0, ID0:1/ - Data0", ""}

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
		Data:      "{\"/valueFromFile1\": \"1\", \"/valueInBothFiles1\": \"1\", \"/valueInBothFiles2\": \"1\"}",
	}
	file2_groupByPatientID = reports.File{
		FileID:    "file_id_2",
		Version:   "version_1",
		PatientID: "patient_1",
		CreatedAt: time2,
		UpdatedAt: time2,
		Data:      "{\"/valueFromFile2\": \"2\", \"/valueInBothFiles1\": \"2\", \"/valueInBothFiles2\": \"1\"}",
	}
	groupedByPatientIDReportRow = []string{"patient_1", "2018-07-27T13:55:59.123Z, 2018-07-29T13:55:59.123Z", "2018-07-29T13:55:59.123Z", "1", "2", "1, 2", "1"}
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

	codeRegexp, err := regexp.Compile(codeRe)
	if err != nil {
		t.Fatalf("failed to compile code regex")
	}

	g := &generator{
		storage:    storage,
		logger:     zerolog.New(os.Stdout),
		codeRegexp: codeRegexp,
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
