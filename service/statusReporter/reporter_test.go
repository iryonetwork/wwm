package statusReporter

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/service/statusReporter/polling"
	"github.com/iryonetwork/wwm/service/statusReporter/polling/mock"
	"github.com/iryonetwork/wwm/status"
)

var (
	// interval for status calls in tests
	interval = time.Duration(2 * time.Millisecond)
)

func TestStatus(t *testing.T) {
	testCases := []struct {
		description    string
		cfg            *polling.Cfg
		localN         int
		cloudN         int
		externalN      int
		mockCalls      func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call
		expectedCallsN int
		expected       *Response
	}{
		{
			"status.OK: no components added",
			&polling.Cfg{Interval: &interval},
			0,
			0,
			0,
			func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call {
				return []*gomock.Call{}
			},
			0,
			&Response{Status: status.OK},
		},
		{
			"status.OK: all components passing",
			&polling.Cfg{Interval: &interval},
			2,
			2,
			1,
			func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call {
				var mockCalls []*gomock.Call
				for _, l := range local {
					mockCalls = append(mockCalls, l.EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil))
				}
				for _, c := range cloud {
					mockCalls = append(mockCalls, c.EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil))
				}
				for _, e := range external {
					mockCalls = append(mockCalls, e.EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil))
				}

				return mockCalls
			},
			5,
			&Response{
				Status: status.OK,
				Local: &status.Response{
					Status: status.OK,
					Components: map[string]*status.Response{
						"local0": &status.Response{Status: status.OK},
						"local1": &status.Response{Status: status.OK},
					},
				},
				Cloud: &status.Response{
					Status: status.OK,
					Components: map[string]*status.Response{
						"cloud0": &status.Response{Status: status.OK},
						"cloud1": &status.Response{Status: status.OK},
					},
				},
				External: &status.Response{
					Status: status.OK,
					Components: map[string]*status.Response{
						"external0": &status.Response{Status: status.OK},
					},
				},
			},
		},
		{
			"status.OK: warning & error returned twice",
			&polling.Cfg{Interval: &interval},
			1,
			1,
			0,
			func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call {
				var mockCalls []*gomock.Call
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(2))
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Warning}, nil).Times(2))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(2))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Error}, nil).Times(2))

				return mockCalls
			},
			8,
			&Response{
				Status: status.OK,
				Local: &status.Response{
					Status: status.OK,
					Components: map[string]*status.Response{
						"local0": &status.Response{Status: status.OK},
					},
				},
				Cloud: &status.Response{
					Status: status.OK,
					Components: map[string]*status.Response{
						"cloud0": &status.Response{Status: status.OK},
					},
				},
			},
		},
		{
			"status.Warning, status.Error, status.OK: last statuses count threshold rule",
			&polling.Cfg{Interval: &interval},
			1,
			2,
			0,
			func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call {
				var mockCalls []*gomock.Call
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(5))
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Warning}, nil).Times(3))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(5))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Error}, nil).Times(3))
				mockCalls = append(mockCalls, cloud[1].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Warning}, nil).Times(5))
				mockCalls = append(mockCalls, cloud[1].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(3))

				return mockCalls
			},
			24,
			&Response{
				Status: status.Warning,
				Local: &status.Response{
					Status: status.Warning,
					Components: map[string]*status.Response{
						"local0": &status.Response{Status: status.Warning},
					},
				},
				Cloud: &status.Response{
					Status: status.Error,
					Components: map[string]*status.Response{
						"cloud0": &status.Response{Status: status.Error},
						"cloud1": &status.Response{Status: status.OK},
					},
				},
			},
		},
		{
			"status.Warning, status.Error, status.OK: majority rule",
			&polling.Cfg{Interval: &interval},
			1,
			2,
			0,
			func(local, cloud, external []*mock.MockURLStatusEndpoint, called chan bool) []*gomock.Call {
				var mockCalls []*gomock.Call
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Warning}, nil).Times(5))
				mockCalls = append(mockCalls, local[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(2))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Error}, nil).Times(5))
				mockCalls = append(mockCalls, cloud[0].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Warning}, nil).Times(2))
				mockCalls = append(mockCalls, cloud[1].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.OK}, nil).Times(5))
				mockCalls = append(mockCalls, cloud[1].EXPECT().FetchStatus().Do(func() { called <- true }).Return(&status.Response{Status: status.Error}, nil).Times(2))

				return mockCalls
			},
			21,
			&Response{
				Status: status.Warning,
				Local: &status.Response{
					Status: status.Warning,
					Components: map[string]*status.Response{
						"local0": &status.Response{Status: status.Warning},
					},
				},
				Cloud: &status.Response{
					Status: status.Error,
					Components: map[string]*status.Response{
						"cloud0": &status.Response{Status: status.Error},
						"cloud1": &status.Response{Status: status.OK},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// get status reporter
			s := New(zerolog.New(os.Stdout))

			// get mock url status endpoints
			local, cloud, external, cleanup := getMockURLStatusEndpoints(t, test.localN, test.cloudN, test.externalN)
			defer cleanup()

			// mock calls
			ch := make(chan bool)
			test.mockCalls(local, cloud, external, ch)

			// add components to service, run mock calls and start polling
			ctx, cancel := context.WithCancel(context.Background())
			for i, url := range local {
				p := polling.New(url, test.cfg, zerolog.New(os.Stdout))
				s.AddComponent(Local, fmt.Sprintf("local%d", i), p)
				p.Start(ctx)
			}
			for i, url := range cloud {
				p := polling.New(url, test.cfg, zerolog.New(os.Stdout))
				s.AddComponent(Cloud, fmt.Sprintf("cloud%d", i), p)
				p.Start(ctx)
			}
			for i, url := range external {
				p := polling.New(url, test.cfg, zerolog.New(os.Stdout))
				s.AddComponent(External, fmt.Sprintf("external%d", i), p)
				p.Start(ctx)
			}

			// wait for all the mock calls to be completed and then cancel context to stop polling
			for i := 0; i < test.expectedCallsN; i++ {
				select {
				case <-ch:
				// do nothing
				case <-time.After(time.Second):
					t.Fatal("waiting too long for expected mock call")
				}
			}
			// wait just a bit longer to make sure all statuses were recorded
			<-time.After(50 * time.Microsecond)
			cancel()
			// wait a bit to test context cancellation
			<-time.After(interval)

			resp := s.Status()

			// check expected results
			if !reflect.DeepEqual(resp, test.expected) {
				t.Errorf("Expected list to equal\n%+v\ngot\n%+v", test.expected, resp)
			}
		})
	}
}

func getMockURLStatusEndpoints(t *testing.T, localN, cloudN, externalN int) ([]*mock.MockURLStatusEndpoint, []*mock.MockURLStatusEndpoint, []*mock.MockURLStatusEndpoint, func()) {
	ctrl := gomock.NewController(t)
	cleanup := func() {
		ctrl.Finish()
	}

	var local, cloud, external []*mock.MockURLStatusEndpoint

	for i := 0; i < localN; i++ {
		url := mock.NewMockURLStatusEndpoint(ctrl)
		local = append(local, url)
	}
	for i := 0; i < cloudN; i++ {
		url := mock.NewMockURLStatusEndpoint(ctrl)
		cloud = append(cloud, url)
	}
	for i := 0; i < externalN; i++ {
		url := mock.NewMockURLStatusEndpoint(ctrl)
		external = append(external, url)
	}

	return local, cloud, external, cleanup
}
