package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3"
	"github.com/iryonetwork/wwm/storage/s3/mock"
	"github.com/iryonetwork/wwm/storage/s3/object"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	mockStorageSync "github.com/iryonetwork/wwm/sync/storage/mock"
)

var (
	time1, _ = strfmt.ParseDateTime("2018-01-18T15:22:46.123Z")
	time2, _ = strfmt.ParseDateTime("2018-01-26T15:16:15.123Z")
	time3, _ = strfmt.ParseDateTime("2018-01-29T11:06:51.223Z")
	file1V1  = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time1,
		Name:        "File1",
		Path:        "BUCKET/File1/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
	}
	file1V2 = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time2,
		Name:        "File1",
		Path:        "BUCKET/File1/V2",
		Version:     "V2",
		Size:        8,
		Operation:   "w",
	}
	file2V1 = &models.FileDescriptor{
		Archetype:   "",
		Checksum:    "CHS",
		ContentType: "image/jpeg",
		Created:     time1,
		Name:        "Image",
		Path:        "BUCKET/Image/V1",
		Version:     "V1",
		Size:        15698,
		Operation:   "w",
	}
	file2V2 = &models.FileDescriptor{
		Archetype:   "",
		Checksum:    "CHS",
		ContentType: "image/jpeg",
		Created:     time3,
		Name:        "Image",
		Path:        "BUCKET/Image/V2",
		Version:     "V2",
		Size:        0,
		Operation:   "d",
	}
	file3V1 = &models.FileDescriptor{
		Archetype:   "ARCH",
		Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
		ContentType: "text/openEhrXml",
		Created:     time2,
		Name:        "FILE3",
		Path:        "BUCKET/FILE3/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
	}
	file3V1ALT = &models.FileDescriptor{
		Archetype:   "ARCH",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time2,
		Name:        "FILE3",
		Path:        "BUCKET/FILE3/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
	}
	bucket1 = &models.BucketDescriptor{
		Name:    "BUCKET1",
		Created: time1,
	}
	bucket2 = &models.BucketDescriptor{
		Name:    "BUCKET2",
		Created: time2,
	}
	noErrors   = false
	withErrors = true
)

func TestChecksum(t *testing.T) {
	expected := "7XACtDnprIRfIjV9giusFERzD722AW0-yUMil7nsn3M="
	svc := service{s3: nil, keyProvider: nil, logger: zerolog.New(os.Stdout)}
	out, err := svc.Checksum(bytes.NewBuffer([]byte("content")))
	if out != expected {
		t.Errorf("Expected %s to equal %s", out, expected)
	}
	if err != nil {
		t.Errorf("Expected err to be nil; got %v", err)
	}
}

func TestBucketList(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      []*models.BucketDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"List fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().ListBuckets().Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			fmt.Errorf("Error"),
		},
		{
			"Successful call",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().ListBuckets().Return([]*models.BucketDescriptor{bucket1, bucket2}, nil),
				}
			},
			[]*models.BucketDescriptor{bucket1, bucket2},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, _, c := getTestService(t)
			defer c()

			// setup calls
			test.calls(s)

			// call BucketList
			out, err := svc.BucketList()

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
			if !reflect.DeepEqual(err, test.exactError) {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestFileList(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      []*models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"BucketExists fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists("BUCKET").Return(false, fmt.Errorf("Error")),
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
					s.EXPECT().BucketExists("BUCKET").Return(true, nil),
					s.EXPECT().List("BUCKET", "").Return(nil, fmt.Errorf("Error")),
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
					s.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			[]*models.FileDescriptor{},
			noErrors,
			nil,
		},
		{
			"Successful call",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists("BUCKET").Return(true, nil),
					s.EXPECT().List("BUCKET", "").Return([]*models.FileDescriptor{file1V2, file2V2, file1V1, file2V1}, nil),
				}
			},
			[]*models.FileDescriptor{file1V2},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, _, c := getTestService(t)
			defer c()

			// setup calls
			test.calls(s)

			// call the MakeBucket
			out, err := svc.FileList("BUCKET")

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

func TestFileNew(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage, *mockStorageSync.MockPublisher) []*gomock.Call
		expected      *models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"MakeBucket fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Write("BUCKET", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				no := &object.NewObjectInfo{
					Archetype:   "ARCH",
					Size:        int64(8),
					Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
					Created:     strfmt.DateTime(time2),
					ContentType: "CONT/TYPE",
					Version:     "UUID",
					Name:        "UUID",
					Operation:   "w",
				}

				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Write("BUCKET", no, gomock.Any()).Return(file1V2, nil),
					p.EXPECT().PublishAsyncWithRetries(gomock.Any(), storageSync.FileNew, gomock.Eq(&storageSync.FileInfo{"BUCKET", "UUID", "UUID", time2})),
				}
			},
			file1V2,
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, p, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time2) }

			// setup calls
			test.calls(s, p)

			// prepare the reader
			r := bytes.NewReader([]byte("contents"))

			// call the MakeBucket
			out, err := svc.FileNew("BUCKET", r, "CONT/TYPE", "ARCH")

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected file descriptor to equal\n%+v\ngot\n%+v", test.expected, out)
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

func TestFileUpdate(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage, *mockStorageSync.MockPublisher) []*gomock.Call
		expected      *models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"Read fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				no := &object.NewObjectInfo{
					Archetype:   "ARCH",
					Size:        int64(8),
					Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
					Created:     strfmt.DateTime(time2),
					ContentType: "CONT/TYPE",
					Version:     "UUID",
					Name:        "FILE",
					Operation:   "w",
				}

				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", no, gomock.Any()).Return(file1V2, nil),
					p.EXPECT().PublishAsyncWithRetries(gomock.Any(), storageSync.FileUpdate, gomock.Eq(&storageSync.FileInfo{"BUCKET", "FILE", "UUID", time2})),
				}
			},
			file1V2,
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, p, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time2) }

			// setup calls
			test.calls(s, p)

			// prepare the reader
			r := bytes.NewReader([]byte("contents"))

			// call the MakeBucket
			out, err := svc.FileUpdate("BUCKET", "FILE", r, "CONT/TYPE", "ARCH")

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected file descriptor to equal\n%+v\ngot\n%+v", test.expected, out)
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

func TestFileDelete(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage, *mockStorageSync.MockPublisher) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"Read fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage, p *mockStorageSync.MockPublisher) []*gomock.Call {
				no := &object.NewObjectInfo{
					Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
					Size:        int64(0),
					Checksum:    "",
					Created:     strfmt.DateTime(time2),
					ContentType: "text/openEhrXml",
					Version:     "UUID",
					Name:        "FILE",
					Operation:   "d",
				}

				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", no, gomock.Any()).Return(file1V2, nil),
					p.EXPECT().PublishAsyncWithRetries(
						gomock.Any(),
						storageSync.FileDelete,
						gomock.Eq(&storageSync.FileInfo{"BUCKET", "FILE", "UUID", strfmt.DateTime(time2)}),
					),
				}
			},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, p, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time2) }

			// setup calls
			test.calls(s, p)

			// call the MakeBucket
			err := svc.FileDelete("BUCKET", "FILE")

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

func TestSyncFileList(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      []*models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"BucketExsits fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists("BUCKET").Return(false, fmt.Errorf("Error")),
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
					s.EXPECT().BucketExists("BUCKET").Return(true, nil),
					s.EXPECT().List("BUCKET", "").Return(nil, fmt.Errorf("Error")),
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
					s.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			[]*models.FileDescriptor{},
			noErrors,
			nil,
		},
		{
			"Successful call",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().BucketExists("BUCKET").Return(true, nil),
					s.EXPECT().List("BUCKET", "").Return([]*models.FileDescriptor{file1V2, file2V2, file1V1, file2V1}, nil),
				}
			},
			[]*models.FileDescriptor{file1V2, file2V2},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, _, c := getTestService(t)
			defer c()

			// setup calls
			test.calls(s)

			// call SyncFileList
			out, err := svc.SyncFileList("BUCKET")

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

func TestSyncFile(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockStorage) []*gomock.Call
		expected      *models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"MakeBucket fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Read fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Read("BUCKET", "FILE3", "V1").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"Already exists - matching checksum",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Read("BUCKET", "FILE3", "V1").Return(nil, file3V1, nil),
				}
			},
			file3V1,
			withErrors,
			ErrAlreadyExists,
		},
		{
			"Already exists - conflict",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Read("BUCKET", "FILE3", "V1").Return(nil, file3V1ALT, nil),
				}
			},
			nil,
			withErrors,
			ErrAlreadyExistsConflict,
		},
		{
			"Write fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Read("BUCKET", "FILE3", "V1").Return(nil, nil, s3.ErrNotFound),
					s.EXPECT().Write("BUCKET", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
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
					Archetype:   "ARCH",
					Size:        int64(8),
					Checksum:    "0bKln76n4gB3r5-Rsn6V6GUGGycL4D_1Oas7c1h4gug=",
					Created:     strfmt.DateTime(time2),
					ContentType: "text/openEhrXml",
					Version:     "V1",
					Name:        "FILE3",
					Operation:   "w",
				}

				return []*gomock.Call{
					s.EXPECT().MakeBucket("BUCKET").Return(nil),
					s.EXPECT().Read("BUCKET", "FILE3", "V1").Return(nil, nil, s3.ErrNotFound),
					s.EXPECT().Write("BUCKET", no, gomock.Any()).Return(file3V1, nil),
				}
			},
			file3V1,
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, _, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }
			getTime = func() strfmt.DateTime { return strfmt.DateTime(time2) }

			// setup calls
			test.calls(s)

			// prepare the reader
			r := bytes.NewReader([]byte("contents"))

			// call the SyncFile
			out, err := svc.SyncFile("BUCKET", "FILE3", "V1", r, "text/openEhrXml", time2, "ARCH")

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected file descriptor to equal\n%+v\ngot\n%+v", test.expected, out)
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

func TestSyncFileDelete(t *testing.T) {
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
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write fails",
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"Write successfull",
			func(s *mock.MockStorage) []*gomock.Call {
				no := &object.NewObjectInfo{
					Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
					Size:        int64(0),
					Checksum:    "",
					Created:     strfmt.DateTime(time2),
					ContentType: "text/openEhrXml",
					Version:     "DEL_VERSION",
					Name:        "FILE",
					Operation:   "d",
				}

				return []*gomock.Call{
					s.EXPECT().Read("BUCKET", "FILE", "").Return(nil, file1V1, nil),
					s.EXPECT().Write("BUCKET", no, gomock.Any()).Return(file1V2, nil),
				}
			},
			noErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init service
			svc, s, _, _, c := getTestService(t)
			defer c()

			// mock getUUID and getTime
			getUUID = func() string { return "UUID" }

			// setup calls
			test.calls(s)

			// call the MakeBucket
			err := svc.SyncFileDelete("BUCKET", "FILE", "DEL_VERSION", strfmt.DateTime(time2))

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

func getTestService(t *testing.T) (*service, *mock.MockStorage, *mock.MockKeyProvider, *mockStorageSync.MockPublisher, func()) {
	// setup s3 mock
	storageCtrl := gomock.NewController(t)
	s3storage := mock.NewMockStorage(storageCtrl)

	// setup key provider mock
	keyProviderCtrl := gomock.NewController(t)
	keyProvider := mock.NewMockKeyProvider(keyProviderCtrl)

	// setup publisher mock
	publisherCtrl := gomock.NewController(t)
	publisher := mockStorageSync.NewMockPublisher(publisherCtrl)

	svc := &service{
		s3:          s3storage,
		keyProvider: keyProvider,
		publisher:   publisher,
		logger:      zerolog.New(os.Stdout),
	}

	cleanup := func() {
		storageCtrl.Finish()
		keyProviderCtrl.Finish()
		publisherCtrl.Finish()
	}

	return svc, s3storage, keyProvider, publisher, cleanup
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
