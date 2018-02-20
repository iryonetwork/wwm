package publisher

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/publisher/mock"
)

var (
	time1, _   = strfmt.ParseDateTime("2018-02-05T15:16:15.123Z")
	file       = &storageSync.FileInfo{"bucket", "file", "version"}
	noErrors   = false
	withErrors = true
)

func TestPublish(t *testing.T) {
	expectedError := fmt.Errorf("error")

	testCases := []struct {
		description   string
		mockCalls     func(*mock.MockStanConnection) []*gomock.Call
		errorExpected bool
		exactError    error
	}{
		{
			"Publish succeeds",
			func(c *mock.MockStanConnection) []*gomock.Call {
				msg, _ := file.Marshal()
				return []*gomock.Call{
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(nil).Times(1),
				}
			},
			noErrors,
			nil,
		},
		{
			"Publish fails",
			func(c *mock.MockStanConnection) []*gomock.Call {
				msg, _ := file.Marshal()
				return []*gomock.Call{
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(expectedError).Times(1),
				}
			},
			withErrors,
			expectedError,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			s, conn, cleanup := getTestPublisher(t, context.Background())
			defer cleanup()

			test.mockCalls(conn)

			// call publish
			err := s.Publish(context.Background(), storageSync.FileNew, file)

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

func TestPublishAsyncWithRetries(t *testing.T) {
	testCases := []struct {
		description   string
		mockCalls     func(*mock.MockStanConnection) []*gomock.Call
		errorExpected bool
	}{
		{
			"Publish succeeds without retries",
			func(c *mock.MockStanConnection) []*gomock.Call {
				msg, _ := file.Marshal()
				return []*gomock.Call{
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(nil).Times(1),
				}
			},
			noErrors,
		},
		{
			"Publish succeeds on second retry",
			func(c *mock.MockStanConnection) []*gomock.Call {
				msg, _ := file.Marshal()
				return []*gomock.Call{
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(fmt.Errorf("error")).Times(2),
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(nil).Times(1),
				}
			},
			noErrors,
		},
		{
			"Publish fails after retry limit",
			func(c *mock.MockStanConnection) []*gomock.Call {
				msg, _ := file.Marshal()
				return []*gomock.Call{
					c.EXPECT().Publish(string(storageSync.FileNew), msg).Return(fmt.Errorf("error")).Times(5),
				}
			},
			noErrors,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			s, conn, cleanup := getTestPublisher(t, context.Background())
			defer cleanup()

			test.mockCalls(conn)

			// call publish async with retries and wait
			err := s.PublishAsyncWithRetries(context.Background(), storageSync.FileNew, file)
			s.wg.Wait()

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func TestPublishContextCancelled(t *testing.T) {
	ctx := context.Background()
	s, conn, cleanup := getTestPublisher(t, ctx)
	defer cleanup()

	// Try to publish once and do not complete retries due to context cancellation
	pubCtx, cancel := context.WithCancel(ctx)
	msg, _ := file.Marshal()
	publishCallReceived := make(chan bool)

	conn.EXPECT().
		Publish(string(storageSync.FileNew), msg).
		Do(func(subject string, data []byte) error {
			publishCallReceived <- true
			return fmt.Errorf("error")
		}).Times(1)

	s.PublishAsyncWithRetries(pubCtx, storageSync.FileNew, file)

	<-publishCallReceived
	cancel()
}

func TestClose(t *testing.T) {
	s, conn, cleanup := getTestPublisher(t, context.Background())
	defer cleanup()

	// First complete 5 retries, then close connection
	msg, _ := file.Marshal()
	pubCall := conn.EXPECT().Publish(string(storageSync.FileNew), msg).Return(fmt.Errorf("error")).Times(5)
	conn.EXPECT().Close().After(pubCall).Return(nil)

	s.PublishAsyncWithRetries(context.Background(), storageSync.FileNew, file)
	s.Close()
}

func TestGeneralContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s, conn, cleanup := getTestPublisher(t, ctx)
	defer cleanup()

	// First complete 5 retries, then close connection
	msg, _ := file.Marshal()
	pubCall := conn.EXPECT().Publish(string(storageSync.FileNew), msg).Return(fmt.Errorf("error")).Times(5)
	conn.EXPECT().Close().After(pubCall).Return(nil)

	s.PublishAsyncWithRetries(context.Background(), storageSync.FileNew, file)
	cancel()
}

func getTestPublisher(t *testing.T, ctx context.Context) (*stanPublisher, *mock.MockStanConnection, func()) {
	mockCtrl := gomock.NewController(t)
	mockConn := mock.NewMockStanConnection(mockCtrl)

	publisher := New(ctx, mockConn, 5, time.Millisecond, 1.0, zerolog.New(os.Stdout))

	cleanup := func() {
		mockConn.EXPECT().Close().Times(1)
		publisher.Close()
		mockCtrl.Finish()
	}

	return publisher.(*stanPublisher), mockConn, cleanup
}
