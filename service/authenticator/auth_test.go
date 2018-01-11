package authenticator

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/service/authenticator/mock"
)

var (
	sampleUser = &models.User{Username: swag.String("username"), Password: "$2a$10$USp/p1VpbjFETLEbtMkVseGu02NgXpaLDP4eYpZiNV5j/nY/qPviW"}
)

func TestLogin(t *testing.T) {
	// prepare mocked storage
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		storage.EXPECT().GetUserByUsername("username").Times(2).Return(sampleUser, nil),
		storage.EXPECT().GetUserByUsername("missing").Times(1).Return(nil, fmt.Errorf("Not found")))

	// initialize service
	svc := &auth{storage: storage}

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
}
