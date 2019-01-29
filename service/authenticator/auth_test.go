package authenticator

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/gobwas/glob"

	"github.com/rs/zerolog"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/service/authenticator/mock"
	"github.com/iryonetwork/wwm/storage/auth"
)

var (
	testClinicID = "d826c3f7-e9cf-4000-8783-4e1b938c87b2"
	sampleUser   = &models.User{ID: "8853C7BC-599A-4F43-8080-6D22B777433E", Username: swag.String("username"), Password: "$2a$10$USp/p1VpbjFETLEbtMkVseGu02NgXpaLDP4eYpZiNV5j/nY/qPviW"}
	domain       = fmt.Sprintf("%s.%s", authCommon.DomainTypeClinic, testClinicID)
	action       = strconv.FormatInt(int64(auth.Write), 10)
	resource     = "/auth/login"
)

func TestLogin(t *testing.T) {
	// prepare mocked authData
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authData := mock.NewMockAuthDataService(ctrl)
	enforcer := mock.NewMockEnforcer(ctrl)
	gomock.InOrder(
		authData.EXPECT().UserByUsername(gomock.Any(), "username").Times(1).Return(sampleUser, nil),
		enforcer.EXPECT().Enforce(sampleUser.ID, domain, resource, action).Times(1).Return(true),
		authData.EXPECT().UserByUsername(gomock.Any(), "username").Times(1).Return(sampleUser, nil),
		enforcer.EXPECT().Enforce(sampleUser.ID, domain, resource, action).Times(1).Return(true),
		authData.EXPECT().UserByUsername(gomock.Any(), "missing").Times(1).Return(nil, fmt.Errorf("Not found")),
		authData.EXPECT().UserByUsername(gomock.Any(), "username").Times(1).Return(sampleUser, nil),
		enforcer.EXPECT().Enforce(sampleUser.ID, domain, resource, action).Times(1).Return(false),
	)

	// initialize service
	svc := &service{domainType: authCommon.DomainTypeClinic, domainID: testClinicID, authData: authData, enforcer: enforcer}

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

	// #3 call with an error from authData
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
	authData := mock.NewMockAuthDataService(ctrl)
	enforcer := mock.NewMockEnforcer(ctrl)

	allowedServiceCertsAndPaths := map[string][]string{
		"testdata/testCert.pem": {
			"/auth/login",
			"/authData/*",
			"/something/other*",
		},
	}

	ss, err := New(authCommon.DomainTypeClinic, testClinicID, authData, enforcer, allowedServiceCertsAndPaths, zerolog.New(ioutil.Discard))
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
				glob: glob.MustCompile("{/api/auth/login,/api/storage/*,/api/something/other*,/api/storage/*/bucket}"),
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
