package publisher

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/storageSync"
	"github.com/iryonetwork/wwm/storageSync/publisher/mock"
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
			s, conn, cleanup := getTestPublisher(t)
			defer cleanup()

			test.mockCalls(conn)

			// call publish
			err := s.Publish(storageSync.FileNew, file)

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
			s, conn, cleanup := getTestPublisher(t)
			defer cleanup()

			test.mockCalls(conn)

			// call publish async with retries and wait
			err := s.PublishAsyncWithRetries(storageSync.FileNew, file)
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

func TestClose(t *testing.T) {
	s, conn, cleanup := getTestPublisher(t)
	defer cleanup()

	// First complete 5 retries, then close connection
	msg, _ := file.Marshal()
	pubCall := conn.EXPECT().Publish(string(storageSync.FileNew), msg).Return(fmt.Errorf("error")).Times(5)
	conn.EXPECT().Close().After(pubCall).Return(nil)

	s.PublishAsyncWithRetries(storageSync.FileNew, file)
	s.Close()
}

func getTestPublisher(t *testing.T) (*stanPublisher, *mock.MockStanConnection, func()) {
	mockCtrl := gomock.NewController(t)
	mockConn := mock.NewMockStanConnection(mockCtrl)

	publisher := &stanPublisher{
		conn:            mockConn,
		logger:          zerolog.New(os.Stdout),
		retries:         5,
		startRetryWait:  time.Millisecond,
		retryWaitFactor: 1.0,
	}

	cleanup := func() {
		mockConn.EXPECT().Close().Times(1)
		publisher.Close()
		mockCtrl.Finish()
	}

	return publisher, mockConn, cleanup
}
