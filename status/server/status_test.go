package server

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/status"
	"github.com/iryonetwork/wwm/status/mock"
)

func TestStatus(t *testing.T) {
	testCases := []struct {
		description string
		mockCalls   func(s1, s2, s3 *mock.MockComponent) []*gomock.Call
		expected    *status.Response
	}{
		{
			"status.OK: all services are ok",
			func(s1, s2, s3 *mock.MockComponent) []*gomock.Call {
				return []*gomock.Call{
					s1.EXPECT().Status().Return(&status.Response{Status: status.OK}),
					s2.EXPECT().Status().Return(&status.Response{Status: status.OK}),
					s3.EXPECT().Status().Return(&status.Response{Status: status.OK}),
				}
			},
			&status.Response{
				Status: status.OK,
				Components: map[string]*status.Response{
					"s1": &status.Response{Status: status.OK},
					"s2": &status.Response{Status: status.OK},
					"s3": &status.Response{Status: status.OK},
				},
			},
		},
		{
			"status.Warning: one service has warning status",
			func(s1, s2, s3 *mock.MockComponent) []*gomock.Call {
				return []*gomock.Call{
					s1.EXPECT().Status().Return(&status.Response{Status: status.OK}),
					s2.EXPECT().Status().Return(&status.Response{Status: status.Warning}),
					s3.EXPECT().Status().Return(&status.Response{Status: status.OK}),
				}
			},
			&status.Response{
				Status: status.Warning,
				Components: map[string]*status.Response{
					"s1": &status.Response{Status: status.OK},
					"s2": &status.Response{Status: status.Warning},
					"s3": &status.Response{Status: status.OK},
				},
			},
		},
		{
			"status.Error: one service has error status",
			func(s1, s2, s3 *mock.MockComponent) []*gomock.Call {
				return []*gomock.Call{
					s1.EXPECT().Status().Return(&status.Response{Status: status.Error}),
					s2.EXPECT().Status().Return(&status.Response{Status: status.Warning}),
					s3.EXPECT().Status().Return(&status.Response{Status: status.OK}),
				}
			},
			&status.Response{
				Status: status.Error,
				Components: map[string]*status.Response{
					"s1": &status.Response{Status: status.Error},
					"s2": &status.Response{Status: status.Warning},
					"s3": &status.Response{Status: status.OK},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// get mock status services
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s1 := mock.NewMockComponent(ctrl)
			s2 := mock.NewMockComponent(ctrl)
			s3 := mock.NewMockComponent(ctrl)

			// get status service
			s := New(zerolog.New(os.Stdout))
			s.AddComponent("s1", s1)
			s.AddComponent("s2", s2)
			s.AddComponent("s3", s3)

			test.mockCalls(s1, s2, s3)

			// call Status()
			status := s.Status()

			// check expected results
			if !reflect.DeepEqual(status, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(status)
				t.Errorf("Expected list to equal\n%+v\ngot\n%+v", test.expected, status)
			}

		})
	}
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
