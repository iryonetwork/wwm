package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

var (
	testRole = &models.Role{
		Name: swag.String("testrole"),
	}
	testRole2 = &models.Role{
		Name: swag.String("testrole2"),
	}
)

func TestAddRole(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user
	storage.AddUser(testUser2)
	testRole.Users = []string{testUser2.ID}

	// add role
	role, err := storage.AddRole(testRole)
	if role.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// add role with invalid user id
	testRole.Users = []string{"wrong user id"}
	_, err = storage.AddRole(testRole)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestGetRole(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and role
	storage.AddUser(testUser2)
	testRole.Users = []string{testUser2.ID}
	storage.AddRole(testRole)

	// get role
	role, err := storage.GetRole(testRole.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testRole, *role) {
		t.Fatalf("Expected returned role to be '%v', got '%v'", *testRole, *role)
	}

	// get role with wrong uuid
	_, err = storage.GetRole("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing role
	_, err = storage.GetRole("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetRoles(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and role
	storage.AddUser(testUser2)
	testRole.Users = []string{testUser2.ID}
	storage.AddRole(testRole)
	storage.AddRole(testRole2)

	// get roles
	roles, err := storage.GetRoles()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(roles) != 4 {
		t.Fatalf("Expected 4 roless; got %d", len(roles))
	}

	rolesMap := map[string]*models.Role{}
	for _, role := range roles {
		rolesMap[role.ID] = role
	}

	if !reflect.DeepEqual(*testRole, *rolesMap[testRole.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testRole, *rolesMap[testRole.ID])
	}

	if !reflect.DeepEqual(*testRole2, *rolesMap[testRole2.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testRole2, *rolesMap[testRole2.ID])
	}
}

func TestUpdateRole(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and role
	storage.AddUser(testUser2)
	testRole.Users = []string{testUser2.ID}
	storage.AddRole(testRole)

	// update role
	updateRole := &models.Role{
		ID:    testRole.ID,
		Users: []string{},
		Name:  swag.String("newname"),
	}
	role, err := storage.UpdateRole(updateRole)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*role, *updateRole) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *role, *updateRole)
	}

	// update role with invalid users
	updateRole.Users = []string{"wrong"}
	_, err = storage.UpdateRole(updateRole)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveRole(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and role
	storage.AddUser(testUser2)
	testRole.Users = []string{testUser2.ID}
	storage.AddRole(testRole)

	// remove role
	err := storage.RemoveRole(testRole.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// remove role again
	err = storage.RemoveRole(testRole.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
