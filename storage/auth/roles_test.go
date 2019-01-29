package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/utils"
)

// method to ensure that roles used for tests are always fresh
func getTestRoles() (*models.Role, *models.Role) {
	testRole := &models.Role{
		Name: swag.String("testrole"),
	}
	testRole2 := &models.Role{
		Name: swag.String("testrole2"),
	}

	return testRole, testRole2
}

func TestAddRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()

	// add role
	role, err := storage.AddRole(testRole)
	if role.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
}

func TestGetRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

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
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, testRole2 := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRole(testRole2)
	errorChecker.FatalTesting(t, err)

	// get roles
	roles, err := storage.GetRoles()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(roles) != 6 {
		t.Fatalf("Expected 6 roles; got %d", len(roles))
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
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, testRole2 := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRole(testRole2)
	errorChecker.FatalTesting(t, err)

	// update role
	updateRole := &models.Role{
		ID:   testRole.ID,
		Name: swag.String("newname"),
	}
	role, err := storage.UpdateRole(updateRole)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*role, *updateRole) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *role, *updateRole)
	}
}

func TestRemoveRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	// add user roles to test if they are removed with role
	testUser1, testUser2 := getTestUsers()
	_, err = storage.AddUser(testUser1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	testOrganization, _ := getTestOrganizations()
	_, err = storage.AddOrganization(testOrganization)
	errorChecker.FatalTesting(t, err)
	testUserRole1 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(testUserRole1)
	errorChecker.FatalTesting(t, err)
	testUserRole2 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeOrganization, testOrganization.ID)
	_, err = storage.AddUserRole(testUserRole2)
	errorChecker.FatalTesting(t, err)
	testUserRole3 := getTestUserRole(testUser2.ID, testRole.ID, authCommon.DomainTypeUser, testUser1.ID)
	_, err = storage.AddUserRole(testUserRole3)
	errorChecker.FatalTesting(t, err)

	// remove role
	err = storage.RemoveRole(testRole.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	// check if user roles were removed
	userRoles, _ := storage.GetUserRoles()
	if len(userRoles) != 6 {
		if err == nil {
			t.Fatalf("Expected 6 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, swag.String(testRole.ID), nil, nil)
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}

	// remove role again
	err = storage.RemoveRole(testRole.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
