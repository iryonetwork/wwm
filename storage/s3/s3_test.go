package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/iryonetwork/wwm/storage/s3/object"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/minio/minio-go"

	"github.com/go-openapi/strfmt"

	"github.com/golang/mock/gomock"
	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/storage/s3/mock"
	objectmock "github.com/iryonetwork/wwm/storage/s3/object/mock"
)

// simple wrapper for turning a Reader to ReadCloser
type noopCloser struct {
	io.Reader
}

func (n noopCloser) Close() error {
	return nil
}

var (
	time1, _ = strfmt.ParseDateTime("2018-01-18T15:22:46.123Z")
	time2, _ = strfmt.ParseDateTime("2018-01-26T15:16:15.123Z")
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
		Labels:      []string{"vitalSign"},
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
		Labels:      []string{"vitalSign", "basicPatientInfo"},
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
		Created:     time2,
		Name:        "Image",
		Path:        "BUCKET/Image/V2",
		Version:     "V2",
		Size:        0,
		Operation:   "d",
	}
	info1V1 = minio.ObjectInfo{
		Key:  "File1.V1.w.1516288966123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWdu",
		Size: 8,
	}
	info1V2 = minio.ObjectInfo{
		Key:  "File1.V2.w.1516979775123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWduLGJhc2ljUGF0aWVudEluZm8=",
		Size: 8,
	}
	info2V1 = minio.ObjectInfo{
		Key:  "Image.V1.w.1516288966123.CHS.aW1hZ2UvanBlZw==..",
		Size: 15698,
	}
	info2V2 = minio.ObjectInfo{
		Key:  "Image.V2.d.1516979775123.CHS.aW1hZ2UvanBlZw==..",
		Size: 0,
	} //Fri Jan 26 16:16:26
	infoErr = minio.ObjectInfo{
		Err: errors.New("error occurred"),
	}
	infoBrokenFD = minio.ObjectInfo{
		ContentType: "text/err",
		Key:         "Error.V1.w.invalid.CHS...",
		Size:        15698,
	}
	removeErr = minio.RemoveObjectError{
		ObjectName: "File1.V2.w.1516979775123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWduLGJhc2ljUGF0aWVudEluZm8=",
		Err:        errors.New("error occurred"),
	}
	newInfo = &object.NewObjectInfo{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time1,
		Name:        "File1",
		Version:     "version",
		Size:        8,
		Labels:      []string{string("vitalSign")},
	}
	noErrors   = false
	withErrors = true
)

func getTestStorage(t *testing.T) (*s3storage, *mock.MockMinio, *mock.MockKeyProvider, func()) {
	// setup minio mock
	minioCtrl := gomock.NewController(t)
	minio := mock.NewMockMinio(minioCtrl)

	// setup key provider mock
	keyProviderCtrl := gomock.NewController(t)
	keyProvider := mock.NewMockKeyProvider(keyProviderCtrl)

	s := &s3storage{
		cfg:    &Config{Region: "REGION"},
		client: minio,
		keys:   keyProvider,
		logger: zerolog.New(os.Stdout),
	}

	cleanup := func() {
		minioCtrl.Finish()
		keyProviderCtrl.Finish()
	}

	return s, minio, keyProvider, cleanup
}

func getTestObject(t *testing.T) (*objectmock.MockObject, func()) {
	objectCtrl := gomock.NewController(t)
	obj := objectmock.NewMockObject(objectCtrl)

	f := func() {
		objectCtrl.Finish()
	}

	return obj, f
}

func TestBucketExists(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockMinio) []*gomock.Call
		errorExpected bool
		exactError    error
		expected      bool
	}{
		{
			"bucket exists",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
				}
			},
			noErrors,
			nil,
			true,
		},
		{
			"bucket does not exist",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			noErrors,
			nil,
			false,
		},
		{
			"bucketExists fails",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
			false,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, _, c := getTestStorage(t)
			defer c()

			// setup calls
			test.calls(m)

			// call the MakeBucket
			exists, err := s.BucketExists(context.TODO(), "BUCKET")

			// check expected results
			if !reflect.DeepEqual(exists, test.expected) {
				t.Errorf("Expected result to equal '%t'; got %t", test.expected, exists)
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

func TestMakeBucket(t *testing.T) {
	testCases := []struct {
		description   string
		calls         func(*mock.MockMinio) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"bucket already exists",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
				}
			},
			withErrors,
			ErrAlreadyExists,
		},
		{
			"bucket does not exist and is created",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
					m.EXPECT().MakeBucket("BUCKET", "REGION").Return(nil),
				}
			},
			noErrors,
			nil,
		},
		{
			"bucketExists fails",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
		{
			"makeBucket fails",
			func(m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
					m.EXPECT().MakeBucket("BUCKET", "REGION").Return(fmt.Errorf("Error")),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, _, c := getTestStorage(t)
			defer c()

			// setup calls
			test.calls(m)

			// call the MakeBucket
			err := s.MakeBucket(context.TODO(), "BUCKET")

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

func TestS3List(t *testing.T) {
	testCases := []struct {
		description   string
		infos         []minio.ObjectInfo
		calls         func(chan minio.ObjectInfo, *mock.MockMinio) []*gomock.Call
		expected      []*models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"valid call",
			[]minio.ObjectInfo{info1V1, info1V2, info2V2, info2V1},
			func(i chan minio.ObjectInfo, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX", false, gomock.Any()).Return(i),
				}
			},
			[]*models.FileDescriptor{file1V2, file2V2, file1V1, file2V1},
			noErrors,
			nil,
		},
		{
			"file contains error",
			[]minio.ObjectInfo{infoErr},
			func(i chan minio.ObjectInfo, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX", false, gomock.Any()).Return(i),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"file has broken filename",
			[]minio.ObjectInfo{infoBrokenFD},
			func(i chan minio.ObjectInfo, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX", false, gomock.Any()).Return(i),
				}
			},
			nil,
			withErrors,
			nil,
		},
		{
			"bucket does not exist",
			[]minio.ObjectInfo{infoBrokenFD},
			func(i chan minio.ObjectInfo, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			[]*models.FileDescriptor{},
			noErrors,
			nil,
		},
		{
			"failed to check if bucket exsits",
			[]minio.ObjectInfo{infoBrokenFD},
			func(i chan minio.ObjectInfo, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, fmt.Errorf("Failed to check if bucket exists")),
				}
			},
			nil,
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, _, c := getTestStorage(t)
			defer c()

			// prepare ObjectInfos channel
			infos := make(chan minio.ObjectInfo, len(test.infos))
			for _, info := range test.infos {
				infos <- info
			}
			close(infos)

			// setup calls
			test.calls(infos, m)

			// call List method
			list, err := s.List(context.TODO(), "BUCKET", "PREFIX")

			// check expected results
			if !reflect.DeepEqual(list, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(list)
				t.Errorf("Expected list to equal\n%+v\ngot\n%+v", test.expected, list)
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

func TestS3Read(t *testing.T) {
	expectedFileName := "File1.V2.w.1516979775123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWduLGJhc2ljUGF0aWVudEluZm8="

	testCases := []struct {
		description   string
		version       string
		listInfos     []minio.ObjectInfo
		calls         func(chan minio.ObjectInfo, *mock.MockMinio, *mock.MockKeyProvider) []*gomock.Call
		reader        []byte
		descriptor    *models.FileDescriptor
		errorExpected bool
		exactError    error
	}{
		{
			"valid call without a version",
			"",
			[]minio.ObjectInfo{info1V2},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				rc := noopCloser{bytes.NewBuffer([]byte("contents"))}
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.", false, gomock.Any()).Return(i),
					k.EXPECT().Get("BUCKET").Return("SECRET-KEY", nil),
					m.EXPECT().GetObjectWithContext(gomock.Any(), "BUCKET", expectedFileName, gomock.Any()).Return(rc, nil),
				}
			},
			[]byte("contents"),
			file1V2,
			noErrors,
			nil,
		},
		{
			"valid call with a version",
			"VERSION",
			[]minio.ObjectInfo{info1V2},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				rc := noopCloser{bytes.NewBuffer([]byte("contents"))}
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					k.EXPECT().Get("BUCKET").Return("SECRET-KEY", nil),
					m.EXPECT().GetObjectWithContext(gomock.Any(), "BUCKET", expectedFileName, gomock.Any()).Return(rc, nil),
				}
			},
			[]byte("contents"),
			file1V2,
			noErrors,
			nil,
		},
		{
			"list fails",
			"",
			[]minio.ObjectInfo{infoErr},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.", false, gomock.Any()).Return(i),
				}
			},
			nil,
			nil,
			withErrors,
			nil,
		},
		{
			"list is empty",
			"",
			[]minio.ObjectInfo{},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.", false, gomock.Any()).Return(i),
				}
			},
			nil,
			nil,
			withErrors,
			ErrNotFound,
		},
		{
			"list is empty (bucket does not exist)",
			"",
			[]minio.ObjectInfo{},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			nil,
			nil,
			withErrors,
			ErrNotFound,
		},
		{
			"key server fails",
			"VERSION",
			[]minio.ObjectInfo{info1V2},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					k.EXPECT().Get("BUCKET").Return("", errors.New("Error")),
				}
			},
			nil,
			nil,
			withErrors,
			nil,
		},
		{
			"GetEncryptedObject fails",
			"VERSION",
			[]minio.ObjectInfo{info1V2},
			func(i chan minio.ObjectInfo, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					k.EXPECT().Get("BUCKET").Return("SECRET-KEY", nil),
					m.EXPECT().GetObjectWithContext(gomock.Any(), "BUCKET", expectedFileName, gomock.Any()).Return(nil, errors.New("Error")),
				}
			},
			nil,
			nil,
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, k, c := getTestStorage(t)
			defer c()

			// prepare ObjectInfos channel
			infos := make(chan minio.ObjectInfo, len(test.listInfos))
			for _, info := range test.listInfos {
				infos <- info
			}
			close(infos)

			// setup calls
			test.calls(infos, m, k)

			// call List method
			reader, fd, err := s.Read(context.TODO(), "BUCKET", "PREFIX", test.version)

			// check expected results
			if !reflect.DeepEqual(fd, test.descriptor) {
				t.Errorf("Expected FileDescriptor to equal\n%+v\ngot\n%+v", test.descriptor, fd)
			}

			if test.reader == nil && reader != nil {
				b, err := ioutil.ReadAll(reader)
				t.Errorf("Expected reader to be nil; got b=%s, err=%v", b, err)
			} else if test.reader != nil {
				if b, err := ioutil.ReadAll(reader); !bytes.Equal(test.reader, b) {
					t.Errorf("Expected '%s' to equal '%s'", b, test.reader)
				} else if err != nil {
					t.Errorf("Expected err from ioutil.ReadAll to be nil; got %v", err)
				}
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

func TestS3Write(t *testing.T) {
	testCases := []struct {
		description   string
		newObject     *object.NewObjectInfo
		calls         func(io.Reader, *mock.MockMinio, *mock.MockKeyProvider) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"valid call",
			&object.NewObjectInfo{
				Name:        "File1",
				Version:     "V1",
				Operation:   "w",
				Created:     time1,
				ContentType: "text/openEhrXml",
				Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
				Checksum:    "CHS",
				Labels:      []string{"vitalSign", "basicPatientInfo"},
			},
			func(r io.Reader, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					k.EXPECT().Get("BUCKET").Return("SECRET", nil),
					m.EXPECT().PutObjectWithContext(
						gomock.Any(),
						"BUCKET",
						"File1.V1.w.1516288966123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWduLGJhc2ljUGF0aWVudEluZm8=",
						r,
						int64(-1),
						gomock.Any(),
					).Return(int64(8), nil),
				}
			},
			noErrors,
			nil,
		},
		{
			"PutObjectWithContext returns error",
			&object.NewObjectInfo{
				Name:        "File1",
				Version:     "V1",
				Operation:   "w",
				Created:     time1,
				ContentType: "text/openEhrXml",
				Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
				Checksum:    "CHS",
				Labels:      []string{"vitalSign", "basicPatientInfo"},
			},
			func(r io.Reader, m *mock.MockMinio, k *mock.MockKeyProvider) []*gomock.Call {
				return []*gomock.Call{
					k.EXPECT().Get("BUCKET").Return("SECRET", nil),
					m.EXPECT().PutObjectWithContext(
						gomock.Any(),
						"BUCKET",
						"File1.V1.w.1516288966123.CHS.dGV4dC9vcGVuRWhyWG1s.b3BlbkVIUi1FSFItT0JTRVJWQVRJT04uYmxvb2RfcHJlc3N1cmUudjE=.dml0YWxTaWduLGJhc2ljUGF0aWVudEluZm8=",
						r,
						int64(-1),
						gomock.Any(),
					).Return(int64(0), errors.New("Error")),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, k, c := getTestStorage(t)
			defer c()

			// setup reader
			reader := bytes.NewBuffer([]byte("contents"))

			// setup calls
			test.calls(reader, m, k)

			// call List method
			_, err := s.Write(context.TODO(), "BUCKET", test.newObject, reader)

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

func TestS3Delete(t *testing.T) {
	testCases := []struct {
		description            string
		version                string
		listInfos              []minio.ObjectInfo
		listRemoveObjectErrors []minio.RemoveObjectError
		calls                  func(chan minio.ObjectInfo, chan minio.RemoveObjectError, *mock.MockMinio) []*gomock.Call
		errorExpected          bool
		exactError             error
	}{
		{
			"valid call without a version",
			"",
			[]minio.ObjectInfo{info1V2},
			[]minio.RemoveObjectError{},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.", false, gomock.Any()).Return(i),
					m.EXPECT().RemoveObjects("BUCKET", gomock.Any()).Return(errCh),
				}
			},
			noErrors,
			nil,
		},
		{
			"valid call with a version",
			"VERSION",
			[]minio.ObjectInfo{info1V2},
			[]minio.RemoveObjectError{},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					m.EXPECT().RemoveObjects("BUCKET", gomock.Any()).Return(errCh),
				}
			},
			noErrors,
			nil,
		},
		{
			"list fails",
			"VERSION",
			[]minio.ObjectInfo{infoErr},
			[]minio.RemoveObjectError{},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
				}
			},
			withErrors,
			nil,
		},
		{
			"list is empty",
			"VERSION",
			[]minio.ObjectInfo{},
			[]minio.RemoveObjectError{},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					m.EXPECT().RemoveObjects("BUCKET", gomock.Any()).Return(errCh),
				}
			},
			noErrors,
			nil,
		},
		{
			"list is empty (bucket does not exist)",
			"VERSION",
			[]minio.ObjectInfo{},
			[]minio.RemoveObjectError{},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(false, nil),
				}
			},
			noErrors,
			nil,
		},
		{
			"failed to remove object",
			"VERSION",
			[]minio.ObjectInfo{info1V2},
			[]minio.RemoveObjectError{removeErr},
			func(i chan minio.ObjectInfo, errCh chan minio.RemoveObjectError, m *mock.MockMinio) []*gomock.Call {
				return []*gomock.Call{
					m.EXPECT().BucketExists("BUCKET").Return(true, nil),
					m.EXPECT().ListObjectsV2("BUCKET", "PREFIX.VERSION.", false, gomock.Any()).Return(i),
					m.EXPECT().RemoveObjects("BUCKET", gomock.Any()).Return(errCh),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// init storage
			s, m, _, c := getTestStorage(t)
			defer c()

			// prepare ObjectInfos channel
			infos := make(chan minio.ObjectInfo, len(test.listInfos))
			for _, info := range test.listInfos {
				infos <- info
			}
			close(infos)

			// prepare RemoveObjectErrors channel
			removeObjectErrors := make(chan minio.RemoveObjectError, len(test.listRemoveObjectErrors))
			for _, removeObjectError := range test.listRemoveObjectErrors {
				removeObjectErrors <- removeObjectError
			}
			close(removeObjectErrors)

			// setup calls
			test.calls(infos, removeObjectErrors, m)

			// call List method
			err := s.Delete(context.TODO(), "BUCKET", "PREFIX", test.version)

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

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
