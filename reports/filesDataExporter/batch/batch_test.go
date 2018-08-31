package batch

import (
	"context"
	"os"
	"testing"
	"time"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/models"
	"github.com/iryonetwork/wwm/reports/filesDataExporter"
	"github.com/iryonetwork/wwm/reports/filesDataExporter/mock"
)

var (
	time1, _ = strfmt.ParseDateTime("2018-02-18T12:36:12.143Z")
	time2, _ = strfmt.ParseDateTime("2018-02-19T12:36:12.143Z")
	time3, _ = strfmt.ParseDateTime("2018-02-20T12:36:12.143Z")
	time4, _ = strfmt.ParseDateTime("2018-02-21T12:36:12.143Z")
	time5, _ = strfmt.ParseDateTime("2018-02-22T12:36:12.143Z")
	bucket1  = &models.BucketDescriptor{
		Name:    "Bucket1",
		Created: time1,
	}
	bucket2 = &models.BucketDescriptor{
		Name:    "Bucket2",
		Created: time1,
	}
	bucketToSkip = &models.BucketDescriptor{
		Name:    "BUCKET_TO_SKIP",
		Created: time1,
	}
	file1V3 = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time5,
		Name:        "File1",
		Path:        "Bucket1/File1/V3",
		Version:     "V3",
		Size:        8,
		Operation:   "d",
	}
	file3V3 = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.body_mass_index.v2",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time4,
		Name:        "File3",
		Path:        "Bucket3/File3/V3",
		Version:     "V3",
		Size:        8,
		Operation:   "w",
	}
	fileToSkip = &models.FileDescriptor{
		Archetype:   "openEHR-EHR-OBSERVATION.blood_pressure.v1",
		Checksum:    "CHS",
		ContentType: "text/openEhrXml",
		Created:     time1,
		Name:        "FileToSkip",
		Path:        "Bucket1/FileToSkip/V1",
		Version:     "V1",
		Size:        8,
		Operation:   "w",
		Labels:      []string{"LABEL_TO_SKIP"},
	}
	noErrors   = false
	withErrors = true
)

func TestExport(t *testing.T) {
	testCases := []struct {
		description   string
		lastRun       time.Time
		mockCalls     func(*mock.MockHandlers) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"Export succesful",
			time.Time(time3),
			func(c *mock.MockHandlers) []*gomock.Call {
				return []*gomock.Call{
					c.EXPECT().
						ListSourceBuckets(gomock.Any()).
						Return([]*models.BucketDescriptor{bucket1, bucket2, bucketToSkip}, nil).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket1.Name, time3).
						Return([]*models.FileDescriptor{file1V3, fileToSkip}, nil).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket2.Name, time3).
						Return([]*models.FileDescriptor{file3V3}, nil).
						Times(1),
					c.EXPECT().
						ExportFileDelete(gomock.Any(), bucket1.Name, file1V3.Name, file1V3.Version, file1V3.Created).
						Return(filesDataExporter.ResultExported, nil).
						Times(1),
					c.EXPECT().
						ExportFile(gomock.Any(), bucket2.Name, file3V3.Name, file3V3.Version, file3V3.Created).
						Return(filesDataExporter.ResultExported, nil).
						Times(1),
				}
			},
			noErrors,
			nil,
		},
		{
			"Failed export of one of the files",
			time.Time(time4),
			func(c *mock.MockHandlers) []*gomock.Call {
				return []*gomock.Call{
					c.EXPECT().
						ListSourceBuckets(gomock.Any()).
						Return([]*models.BucketDescriptor{bucket1, bucket2, bucketToSkip}, nil).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket1.Name, time4).
						Return([]*models.FileDescriptor{file1V3, fileToSkip}, nil).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket2.Name, time4).
						Return([]*models.FileDescriptor{}, nil).
						Times(1),
					c.EXPECT().
						ExportFileDelete(gomock.Any(), bucket1.Name, file1V3.Name, file1V3.Version, file1V3.Created).
						Return(filesDataExporter.ResultError, errors.Errorf("fail")).
						Times(1),
				}
			},
			withErrors,
			errors.Errorf("1 failure(s) out of 2 bucket(s) to export"),
		},
		{
			"Failed to list source files for one of the buckets",
			time.Time(time1),
			func(c *mock.MockHandlers) []*gomock.Call {
				return []*gomock.Call{
					c.EXPECT().
						ListSourceBuckets(gomock.Any()).
						Return([]*models.BucketDescriptor{bucket1, bucket2, bucketToSkip}, nil).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket1.Name, time1).
						Return(nil, errors.Errorf("fail")).
						Times(1),
					c.EXPECT().
						ListSourceFilesAsc(gomock.Any(), bucket2.Name, time1).
						Return([]*models.FileDescriptor{file3V3}, nil).
						Times(1),
					c.EXPECT().
						ExportFile(gomock.Any(), bucket2.Name, file3V3.Name, file3V3.Version, file3V3.Created).
						Return(filesDataExporter.ResultExported, nil).
						Times(1),
				}
			},
			withErrors,
			errors.Errorf("1 failure(s) out of 2 bucket(s) to export"),
		},
		{
			"Failed to list buckets",
			time.Time(time1),
			func(c *mock.MockHandlers) []*gomock.Call {
				return []*gomock.Call{
					c.EXPECT().
						ListSourceBuckets(gomock.Any()).
						Return(nil, errors.Errorf("fail")).
						Times(1),
				}
			},
			withErrors,
			errors.Errorf("failed to list source buckets: fail"),
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			h, cleanup := getMockHandlers(t)
			defer cleanup()
			s := getTestService(t, h)

			test.mockCalls(h)

			// call sync
			err := s.Export(context.Background(), test.lastRun)

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if test.exactError != nil && err.Error() != test.exactError.Error() {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h, cleanup := getMockHandlers(t)
	defer cleanup()
	s := getTestService(t, h)

	called := make(chan bool)
	contextCancelled := make(chan bool)

	h.EXPECT().
		ListSourceBuckets(gomock.Any()).
		Return([]*models.BucketDescriptor{bucket1}, nil).
		Times(1)
	h.EXPECT().
		ListSourceFilesAsc(gomock.Any(), bucket1.Name, time2).
		Return([]*models.FileDescriptor{file1V3, file1V3}, nil).
		Times(1)
	h.EXPECT().
		ExportFileDelete(gomock.Any(), bucket1.Name, file1V3.Name, file1V3.Version, file1V3.Created).
		Return(filesDataExporter.ResultExported, nil).
		Do(func(_ context.Context, _, _, _ string, _ strfmt.DateTime) {
			called <- true
			<-contextCancelled
		}).
		Times(1)
	//cancel()

	errCh := make(chan error)
	go func() {
		errCh <- s.Export(ctx, time.Time(time2))
	}()

	<-called
	cancel()
	contextCancelled <- true
	// If context cancellation failed there will be missing mock expectations as there were files to sync
	time.Sleep(time.Duration(50 * time.Millisecond))

	if err := <-errCh; err == nil {
		t.Fatal("Got no error. Expected an error")
	}
}

func getMockHandlers(t *testing.T) (*mock.MockHandlers, func()) {
	mockHandlersCtrl := gomock.NewController(t)
	mockHandlers := mock.NewMockHandlers(mockHandlersCtrl)

	cleanup := func() {
		mockHandlersCtrl.Finish()
	}

	return mockHandlers, cleanup
}

func getTestService(t *testing.T, handlers filesDataExporter.Handlers) filesDataExporter.BatchFilesDataExporter {
	return New(
		handlers,
		1,
		[]string{"BUCKET_TO_SKIP"},
		[]string{"LABEL_TO_SKIP"},
		zerolog.New(os.Stdout),
	)
}
