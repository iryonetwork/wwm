package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/reports"
	"github.com/iryonetwork/wwm/storage/reports/db/mock"
)

var (
	uuid1        = "45a1b7dd-dc9a-44b6-9c52-59a1489e3792"
	uuid2        = "9f6e2d98-fde6-401e-8ef2-c33cd1038c51"
	uuid3        = "afc1a438-0443-4cef-b99b-3d6f44e8f979"
	uuid4        = "f6f4700e-bc7f-4ff7-9bcd-c139ec13fa2e"
	uuid5        = "7cd54cf3-c01c-49eb-be23-8aae397701a3"
	version1UUID = "74dc1c52-2993-4985-8940-faf61b329788"
	version2UUID = "ff76b73e-03b5-415f-9123-5a0133513e33"
	noErrors     = false
	withErrors   = true
	time1, _     = strfmt.ParseDateTime("2018-01-18T15:22:46.123Z")
	time2, _     = strfmt.ParseDateTime("2018-01-19T18:33:57.123Z")
	data         = "{\"it_is\": \"just_some_json_string\"}"

	file1 = reports.File{
		FileID:    uuid1,
		Version:   version1UUID,
		PatientID: uuid2,
		CreatedAt: time1,
		UpdatedAt: time1,
		Data:      data,
	}
	file2 = reports.File{
		FileID:    uuid3,
		Version:   version1UUID,
		PatientID: uuid2,
		CreatedAt: time2,
		UpdatedAt: time2,
		Data:      data,
	}
	file3 = reports.File{
		FileID:    uuid4,
		Version:   version1UUID,
		PatientID: uuid5,
		CreatedAt: time1,
		UpdatedAt: time2,
		Data:      data,
	}
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		title         string
		fileID        string
		version       string
		patientID     string
		timestamp     strfmt.DateTime
		data          string
		calls         func(*mock.MockDB)
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull insert",
			uuid1,
			version1UUID,
			uuid2,
			time2,
			data,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Exec(
						"INSERT INTO \"files\" (file_id, version, patient_id, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (file_id) DO UPDATE SET version = ?, updated_at = ?, data = ?",
						uuid1,
						version1UUID,
						uuid2,
						time2.String(),
						time2.String(),
						data,
						version1UUID,
						time2.String(),
						data,
					).Return(db),
					db.EXPECT().GetError().Return(nil))
			},
			noErrors,
			nil,
		},
		{
			"Error on insert",
			uuid1,
			version1UUID,
			uuid2,
			time2,
			data,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Exec(
						"INSERT INTO \"files\" (file_id, version, patient_id, created_at, updated_at, data) VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (file_id) DO UPDATE SET version = ?, updated_at = ?, data = ?",
						uuid1,
						version1UUID,
						uuid2,
						time1.String(),
						time1.String(),
						data,
						version1UUID,
						time1.String(),
						data,
					).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")),
				)
			},
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			err := s.Insert(tc.fileID, tc.version, tc.patientID, tc.timestamp, tc.data)

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		title         string
		fileID        string
		version       string
		calls         func(*mock.MockDB)
		expected      *reports.File
		errorExpected bool
		exactError    error
	}{
		{
			"Successfull get, call without version",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&file1,
			noErrors,
			nil,
		},
		{
			"File not found, call without version",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound),
				)
			},
			nil,
			noErrors,
			nil,
		},
		{
			"Successfull get with version",
			uuid1,
			version1UUID,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&file1,
			noErrors,
			nil,
		},
		{
			"Different version in DB",
			uuid1,
			version2UUID,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			nil,
			noErrors,
			nil,
		},
		{
			"Error on query",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")),
				)
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Get(tc.fileID, tc.version)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	testCases := []struct {
		title         string
		fileID        string
		version       string
		calls         func(*mock.MockDB)
		expected      bool
		errorExpected bool
		exactError    error
	}{
		{
			"File exists in DB, call without version",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().Select("version").Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			true,
			noErrors,
			nil,
		},
		{
			"File does not exist in DB, call without version",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().Select("version").Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound),
				)
			},
			false,
			noErrors,
			nil,
		},
		{
			"File exists in DB, call with version",
			uuid1,
			version1UUID,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().Select("version").Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			true,
			noErrors,
			nil,
		},
		{
			"Different version in DB",
			uuid1,
			version2UUID,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().Select("version").Return(db),
					db.EXPECT().First(gomock.Any()).Do(func(f *reports.File) {
						f.FileID = file1.FileID
						f.Version = file1.Version
						f.PatientID = file1.PatientID
						f.CreatedAt = file1.CreatedAt
						f.UpdatedAt = file1.UpdatedAt
						f.Data = file1.Data
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			false,
			noErrors,
			nil,
		},
		{
			"Error on query",
			uuid1,
			"",
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"file_id = ?",
						uuid1,
					).Return(db),
					db.EXPECT().Select("version").Return(db),
					db.EXPECT().First(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")),
				)
			},
			false,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Exists(tc.fileID, tc.version)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func TestFind(t *testing.T) {
	testCases := []struct {
		title          string
		patientID      string
		dataKeyValues  map[string]string
		createdAtStart *strfmt.DateTime
		createdAtEnd   *strfmt.DateTime
		calls          func(*mock.MockDB)
		expected       *[]reports.File
		errorExpected  bool
		exactError     error
	}{
		{
			"Succesful search by patientID only",
			uuid2,
			map[string]string{},
			nil,
			nil,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"patient_id = ?",
						[]interface{}{uuid2},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Do(func(files *[]reports.File) {
						*files = append(*files, file1, file2)
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&[]reports.File{file1, file2},
			noErrors,
			nil,
		},
		{
			"Succesful search by patientID and a key",
			uuid2,
			map[string]string{"it_is": "just_some_json_string"},
			nil,
			nil,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"patient_id = ? AND data->>'it_is' = ?",
						[]interface{}{uuid2, "just_some_json_string"},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Do(func(files *[]reports.File) {
						*files = append(*files, file1, file2)
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&[]reports.File{file1, file2},
			noErrors,
			nil,
		},
		{
			"Succesful search by key and dates",
			"",
			map[string]string{"it_is": "just_some_json_string"},
			&time1,
			&time1,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"created_at >= ? AND created_at <= ? AND data->>'it_is' = ?",
						[]interface{}{time1.String(), time1.String(), "just_some_json_string"},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Do(func(files *[]reports.File) {
						*files = append(*files, file1, file3)
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&[]reports.File{file1, file3},
			noErrors,
			nil,
		},
		{
			"Succesful search by patientID and createdAtStart",
			uuid2,
			nil,
			&time2,
			nil,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"patient_id = ? AND created_at >= ?",
						[]interface{}{uuid2, time2.String()},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Do(func(files *[]reports.File) {
						*files = append(*files, file2)
					}).Return(db),
					db.EXPECT().GetError().Return(nil),
				)
			},
			&[]reports.File{file2},
			noErrors,
			nil,
		},
		{
			"Record not found",
			uuid2,
			nil,
			&time2,
			nil,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"patient_id = ? AND created_at >= ?",
						[]interface{}{uuid2, time2.String()},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(gorm.ErrRecordNotFound),
				)
			},
			&[]reports.File{},
			noErrors,
			nil,
		},
		{
			"Error on query",
			uuid2,
			nil,
			&time2,
			nil,
			func(db *mock.MockDB) {
				gomock.InOrder(
					db.EXPECT().Where(
						"patient_id = ? AND created_at >= ?",
						[]interface{}{uuid2, time2.String()},
					).Return(db),
					db.EXPECT().Order("created_at asc").Return(db),
					db.EXPECT().Find(gomock.Any()).Return(db),
					db.EXPECT().GetError().Return(fmt.Errorf("error")),
				)
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// init storage
			s, db, c := getTestStorage(t)
			defer c()

			// collect mocked calls
			tc.calls(db)

			// call the method
			out, err := s.Find(tc.patientID, tc.dataKeyValues, tc.createdAtStart, tc.createdAtEnd)

			// check expected results
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected\n%s\nto equal\n\t\t%s", toJSON(out), toJSON(tc.expected))
			}

			// assert error
			if tc.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if tc.exactError != nil && tc.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", tc.exactError, err)
			}
		})
	}
}

func getTestStorage(t *testing.T) (*storage, *mock.MockDB, func()) {
	// setup db mock
	dbCtrl := gomock.NewController(t)
	db := mock.NewMockDB(dbCtrl)

	s := &storage{
		ctx:    context.Background(),
		db:     db,
		logger: zerolog.New(os.Stdout),
	}

	cleanup := func() {
		dbCtrl.Finish()
	}

	return s, db, cleanup
}

func toJSON(in interface{}) string {
	buf := bytes.NewBuffer(nil)
	errorChecker.FatalError(json.NewEncoder(buf).Encode(in))
	return buf.String()
}
