package authenticator

import (
	"context"
	"fmt"
	"testing"

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
