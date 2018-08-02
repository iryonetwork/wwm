package reportsStorage

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
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/reportsStorage/models"
	storageModels "github.com/iryonetwork/wwm/gen/storage/models"

	"github.com/iryonetwork/wwm/storage/s3/mock"
	"github.com/iryonetwork/wwm/storage/s3/object"
)

var (
	time1, _ = strfmt.ParseDateTime("2018-01-09T13:10:07.123Z")
	time2, _ = strfmt.ParseDateTime("2018-01-26T15:16:15.123Z")
	time3, _ = strfmt.ParseDateTime("2018-01-29T11:06:51.223Z")
	rfd1v1   = &models.ReportFileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time1,
		Name:        "MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==",
		Path:        "REPORT_TYPE/MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
		ReportType:  "REPORT_TYPE",
		DataSince:   time1,
		DataUntil:   time2,
	}
	fd1v1 = &storageModels.FileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time1,
		Name:        "MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==",
		Path:        "REPORT_TYPE/MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
	}
	rfd1v2 = &models.ReportFileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time2,
		Name:        "MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==",
		Path:        "REPORT_TYPE/MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==/V2",
		Version:     "V2",
		Size:        8,
		Operation:   "w",
		ReportType:  "REPORT_TYPE",
		DataSince:   time1,
		DataUntil:   time2,
	}
	fd1v2 = &storageModels.FileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time2,
		Name:        "MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==",
		Path:        "REPORT_TYPE/MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==/V2",
		Version:     "V2",
		Size:        8,
		Operation:   "w",
	}
	rfd2v1 = &models.ReportFileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time1,
		Name:        "MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa",
		Path:        "REPORT_TYPE/MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa/V1",
		Version:     "V1",
		Size:        15698,
		Operation:   "w",
		ReportType:  "REPORT_TYPE",
		DataUntil:   time3,
	}
	fd2v1 = &storageModels.FileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time1,
		Name:        "MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa",
		Path:        "REPORT_TYPE/MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa/V1",
		Version:     "V1",
		Size:        15698,
		Operation:   "w",
	}
	rfd2v2 = &models.ReportFileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time3,
		Name:        "MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa",
		Path:        "REPORT_TYPE/MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa/V2",
		Version:     "UUID",
		Size:        0,
		Operation:   "d",
		ReportType:  "REPORT_TYPE",
		DataUntil:   time3,
	}
	fd2v2 = &storageModels.FileDescriptor{
		Checksum:    "CHS",
		ContentType: "text/csv",
		Created:     time3,
		Name:        "MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa",
		Path:        "REPORT_TYPE/MjAxOC0wMS0yOVQxMTowNjo1MS4yMjNa/V2",
		Version:     "UUID",
		Size:        0,
		Operation:   "d",
	}
	noErrors   = false
	withErrors = true
)

func TestChecksum(t *testing.T) {
	expected := "7XACtDnprIRfIjV9giusFERzD722AW0-yUMil7nsn3M="
	svc := service{s3: nil, logger: zerolog.New(os.Stdout)}
	out, err := svc.Checksum(bytes.NewBuffer([]byte("content")))
	if out != expected {
		t.Errorf("Expected %s to equal %s", out, expected)
	}
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
}

func TestReportList(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      []*models.ReportFileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"BucketExists fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(false, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"List fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "").Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Bucket does not exist",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(false, nil),
				}
			},
			[]*models.ReportFileDescriptor{},
			noErrors,
			nil,
		},
		{
			"Successful call",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "").Return([]*storageModels.FileDescriptor{fd1v2, fd2v2, fd1v1, fd2v1}, nil),
				}
			},
			[]*models.ReportFileDescriptor{rfd1v2},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, c := getTestService(t)
			defer c()

			// setup calls
			test.calls(s)

			// call the MakeBucket
			out, err := svc.ReportList(context.TODO(), "REPORT_TYPE")

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected list to equal\n%+v\ngot\n%+v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && test.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestReportNew(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      *models.ReportFileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"MakeBucket fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket(gomock.Any(), "REPORT_TYPE").Return(fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket(gomock.Any(), "REPORT_TYPE").Return(nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage) []*gomock.Call {
				no := &object.NewObjectInfo{
					Size:        int64(8),
					Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
					Created:     strfmt.DateTime(time1),
					ContentType: "CONT/TYPE",
					Version:     "UUID",
					Name:        "MjAxOC0wMS0wOVQxMzoxMDowNy4xMjNaLzIwMTgtMDEtMjZUMTU6MTY6MTUuMTIzWg==",
					Operation:   "w",
				}
				return []*gomock.Call{
					s.EXPECT().MakeBucket(gomock.Any(), "REPORT_TYPE").Return(nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", no, gomock.Any()).Return(fd1v1, nil).Times(1),
				}
			},
			rfd1v1,
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time1) }

			// setup calls
			test.calls(s)

			// prepare the reader
			r := bytes.NewReader([]byte("contents"))

			out, err := svc.ReportNew(context.TODO(), "REPORT_TYPE", r, "CONT/TYPE", &time1, time2)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected report file descriptor to equal\n%+v\ngot\n%+v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && test.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestReportUpdate(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      *models.ReportFileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"Read fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, fd1v1, nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage) []*gomock.Call {
				no := &object.NewObjectInfo{
					Size:        int64(8),
					Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
					Created:     strfmt.DateTime(time2),
					ContentType: "CONT/TYPE",
					Version:     "UUID",
					Name:        "FILE",
					Operation:   "w",
				}

				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, fd1v1, nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", no, gomock.Any()).Return(fd1v2, nil),
				}
			},
			rfd1v2,
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time2) }

			// setup calls
			test.calls(s)

			// prepare the reader
			r := bytes.NewReader([]byte("contents"))

			out, err := svc.ReportUpdate(context.TODO(), "REPORT_TYPE", "FILE", r, "CONT/TYPE")

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected report file descriptor to equal\n%+v\ngot\n%+v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && test.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestReportDelete(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"Read fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, fd2v1, nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage) []*gomock.Call {
				no := &object.NewObjectInfo{
					Size:        int64(0),
					Created:     strfmt.DateTime(time3),
					ContentType: "text/csv",
					Version:     "UUID",
					Name:        "FILE",
					Operation:   "d",
				}

				return []*gomock.Call{
					s.EXPECT().Read(gomock.Any(), "REPORT_TYPE", "FILE", "").Return(nil, fd2v1, nil),
					s.EXPECT().Write(gomock.Any(), "REPORT_TYPE", no, gomock.Any()).Return(fd2v2, nil),
				}
			},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time3) }

			// setup calls
			test.calls(s)

			// call the MakeBucket
			err := svc.ReportDelete(context.TODO(), "REPORT_TYPE", "FILE")

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && test.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestReportListVersions(t *testing.T) {
	testCases := []struct {
		description    string
		createdAtSince *strfmt.DateTime
		createdAtUntil *strfmt.DateTime
		calls          func(*mock.MockStorage) []*gomock.Call
		expected       []*models.ReportFileDescriptor
		errorExpected  bool
		exactError     error
	}{
		{
			"BucketExsits fails",
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(false, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"List fails",
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "Report1").Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Bucket does not exist",
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(false, nil),
				}
			},
			[]*models.ReportFileDescriptor{},
			noErrors,
			nil,
		},
		{
			"Successful call",
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "Report1").Return([]*storageModels.FileDescriptor{fd1v2, fd1v1}, nil),
				}
			},
			[]*models.ReportFileDescriptor{rfd1v2, rfd1v1},
			noErrors,
			nil,
		},
		{
			"Successful call with createdAtSince filtering",
			&time1,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "Report1").Return([]*storageModels.FileDescriptor{fd1v2, fd1v1}, nil),
				}
			},
			[]*models.ReportFileDescriptor{rfd1v2},
			noErrors,
			nil,
		},
		{
			"Successful call with createdAtUntil filtering",
			nil,
			&time2,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "Report1").Return([]*storageModels.FileDescriptor{fd1v2, fd1v1}, nil),
				}
			},
			[]*models.ReportFileDescriptor{rfd1v1},
			noErrors,
			nil,
		},
		{
			"Successful call with createdAtSince and createdAtUntil filtering",
			&time1,
			&time2,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists(gomock.Any(), "REPORT_TYPE").Return(true, nil),
					s.EXPECT().List(gomock.Any(), "REPORT_TYPE", "Report1").Return([]*storageModels.FileDescriptor{fd1v2, fd1v1}, nil),
				}
			},
			[]*models.ReportFileDescriptor{},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, c := getTestService(t)
			defer c()

			// setup calls
			test.calls(s)

			// call SyncReportList
			out, err := svc.ReportListVersions(context.TODO(), "REPORT_TYPE", "Report1", test.createdAtSince, test.createdAtUntil)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected list to equal\n%+v\ngot\n%+v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && test.exactError != err {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func getTestService(t *testing.T) (*service, *mock.MockStorage, func()) {
	// setup s3 mock
	storageCtrl := gomock.NewController(t)
	s3storage := mock.NewMockStorage(storageCtrl)

	svc := &service{
		s3:     s3storage,
		logger: zerolog.New(os.Stdout),
	}

	cleanup := func() {
		storageCtrl.Finish()
	}

	return svc, s3storage, cleanup
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
