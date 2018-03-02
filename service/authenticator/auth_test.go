package authenticator

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gobwas/glob"

	"github.com/rs/zerolog"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/service/authenticator/mock"
	"github.com/iryonetwork/wwm/storage/auth"
)

var (
	sampleUser = &models.User{ID: "8853C7BC-599A-4F43-8080-6D22B777433E", Username: swag.String("username"), Password: "$2a$10$USp/p1VpbjFETLEbtMkVseGu02NgXpaLDP4eYpZiNV5j/nY/qPviW"}
	aclRequest = &models.ValidationPair{Actions: swag.Int64(auth.Write), Resource: swag.String("/auth/login")}
)

func TestLogin(t *testing.T) {
	// prepare mocked storage
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		storage.EXPECT().GetUserByUsername("username").Times(1).Return(sampleUser, nil),
		storage.EXPECT().FindACL(sampleUser.ID, gomock.Any()).Times(1).Return([]*models.ValidationResult{{Query: aclRequest, Result: true}}),
		storage.EXPECT().GetUserByUsername("username").Times(1).Return(sampleUser, nil),
		storage.EXPECT().FindACL(sampleUser.ID, gomock.Any()).Times(1).Return([]*models.ValidationResult{{Query: aclRequest, Result: true}}),
		storage.EXPECT().GetUserByUsername("missing").Times(1).Return(nil, fmt.Errorf("Not found")),
		storage.EXPECT().GetUserByUsername("username").Times(1).Return(sampleUser, nil),
		storage.EXPECT().FindACL(sampleUser.ID, gomock.Any()).Times(1).Return([]*models.ValidationResult{{Query: aclRequest, Result: false}}))

	// initialize service
	svc := &service{storage: storage}

	// #1 call with a valid username and password
	out, err := svc.Login(context.Background(), "username", "password")
	if out == "" {
		t.Errorf("Expected login to return a token, got an empty string")
	}
	if err != nil {
		t.Errorf("Expected error to be nil; got '%v'", err)
	}

	// #2 call with an invalid password
	out, err = svc.Login(context.Background(), "username", "wrongPassword")
	if out != "" {
		t.Errorf("Expected login to return empty token, got %v", out)
	}
	if err == nil {
		t.Errorf("Expected error to be nil; got %v", err)
	}

	// #3 call with an error from storage
	out, err = svc.Login(context.Background(), "missing", "password")
	if out != "" {
		t.Errorf("Expected login to return an empty string, got %v", out)
	}
	if err == nil {
		t.Errorf("Expected error; got nil")
	}

	// #4 call with invalid login permissions
	out, err = svc.Login(context.Background(), "username", "password")
	if out != "" {
		t.Errorf("Expected login to return an empty string, got %v", out)
	}
	if err == nil {
		t.Errorf("Expected error; got nil")
	}
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mock.NewMockStorage(ctrl)

	allowedServiceCertsAndPaths := map[string][]string{
		"testdata/testCert.pem": {
			"/auth/login",
			"/storage/*",
			"/something/other*",
		},
	}

	ss, err := New(storage, allowedServiceCertsAndPaths, zerolog.New(ioutil.Discard))
	if err != nil {
		t.Fatalf("Expected error to be nil; got %v", err)
	}
	s := ss.(*service)

	_, ok := s.syncServices["YcS9Uj_ddqPxsJc9ISJYPLhJTRgIZPqE3T8fX3s9Q6I"]
	if !ok {
		t.Fatalf("Expected service to have syncService with key id 'YcS9Uj_ddqPxsJc9ISJYPLhJTRgIZPqE3T8fX3s9Q6I'; got sync services: %v", s.syncServices)
	}
}

func TestAuthorizerForSyncPaths(t *testing.T) {
	s := &service{
		syncServices: map[string]syncService{
			"test_cert_key_id": {
				glob: glob.MustCompile("{/auth/login,/storage/*,/something/other*,/storage/*/bucket}"),
			},
		},
	}

	authorizer := s.Authorizer()

	tests := []struct {
		path    string
		success bool
	}{
		{"/auth/login", true},
		{"/auth/logindsa", false},
		{"/storage/", true},
		{"/storage/dsa", true},
		{"storage", false},
		{"/something/other", true},
		{"/something/otherdsadsa", true},
		{"/something/othe", false},
		{"/storage/abc/bucket", true},
	}

	for _, test := range tests {
		r, _ := http.NewRequest("GET", test.path, nil)
		err := authorizer.Authorize(r, swag.String(servicePrincipal+"test_cert_key_id"))
		if err == nil && !test.success {
			t.Errorf("Authorizing path '%s' got success; expected fail", test.path)
		}
		if err != nil && test.success {
			t.Errorf("Authorizing path '%s' got fail; expected success", test.path)
		}
	}

}
