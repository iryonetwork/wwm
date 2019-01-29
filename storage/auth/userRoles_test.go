package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/utils"
)

func getTestUserRole(userID, roleID, domainType, domainID string) *models.UserRole {
	return &models.UserRole{
		UserID:     swag.String(userID),
		RoleID:     swag.String(roleID),
		DomainType: swag.String(domainType),
		DomainID:   swag.String(domainID),
	}
}

func TestAddUserRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	testUser, _ := getTestUsers()
	_, err = storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	userRole, err := storage.AddUserRole(testUserRole)
	if userRole.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// add same userRole again
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestAddUserRoleUserDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	nonExistingUserID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(nonExistingUserID.String(), testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleRoleDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)

	nonExistingRoleID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(testUser.ID, nonExistingRoleID.String(), authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleDomainClinicDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	nonExistingClinicID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeClinic, nonExistingClinicID.String())

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleDomainOrganizationDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	nonExistingOrganizationID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeOrganization, nonExistingOrganizationID.String())

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleDomainLocationDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	nonExistingLocationID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeLocation, nonExistingLocationID.String())

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleDomainUserDoesNotExist(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	nonExistingUserID, _ := uuid.NewV4()

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeUser, nonExistingUserID.String())

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestAddUserRoleInvalidDomainType(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, "something", "someID")

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	if err == nil {
		t.Fatalf("Expected error; got '%v'", err)
	}
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestGetUserRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	errorChecker.FatalTesting(t, err)

	// get userRole
	userRole, err := storage.GetUserRole(testUserRole.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testUserRole, *userRole) {
		t.Fatalf("Expected returned user to be '%v', got '%v'", *testUserRole, *userRole)
	}

	// get userRole with wrong uuid
	_, err = storage.GetUserRole("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing userRole
	_, err = storage.GetUserRole("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetUserRoles(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testRole, _ := getTestRoles()
	_, err := storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	testUser1, testUser2 := getTestUsers()
	_, err = storage.AddUser(testUser1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)

	testUserRole1 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	testUserRole2 := getTestUserRole(testUser2.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(testUserRole1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(testUserRole2)
	errorChecker.FatalTesting(t, err)

	// get userRoles
	userRoles, err := storage.GetUserRoles()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if len(userRoles) != 8 {
		t.Fatalf("Expected 8 user roles; got %d", len(userRoles))
	}

	userRolesMap := map[string]*models.UserRole{}
	for _, userRole := range userRoles {
		userRolesMap[userRole.ID] = userRole
	}

	if !reflect.DeepEqual(*testUserRole1, *userRolesMap[testUserRole1.ID]) {
		t.Fatalf("Expected user role one to be '%v', got '%v'", *testUserRole1, *userRolesMap[testUserRole1.ID])
	}

	if !reflect.DeepEqual(*testUserRole2, *userRolesMap[testUserRole2.ID]) {
		t.Fatalf("Expected user role two to be '%v', got '%v'", *testUserRole2, *userRolesMap[testUserRole2.ID])
	}
}

func TestGetUserRoleByContent(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	errorChecker.FatalTesting(t, err)

	// get userRole
	userRole, err := storage.GetUserRoleByContent(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testUserRole, *userRole) {
		t.Fatalf("Expected returned user to be '%v', got '%v'", *testUserRole, *userRole)
	}

	// get user role by content with wrong user ID
	_, err = storage.GetUserRoleByContent("E4363A8D-4041-4B17-A43E-17705C96C1CD", testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}
}

func TestFindUserRoles(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	// populate DB varied set of user roles
	// add users
	testUser1, testUser2 := getTestUsers()
	_, err := storage.AddUser(testUser1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	// add roles
	testRole1, testRole2 := getTestRoles()
	_, err = storage.AddRole(testRole1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRole(testRole2)
	errorChecker.FatalTesting(t, err)
	// add organizations
	testOrganization1, testOrganization2 := getTestOrganizations()
	_, err = storage.AddOrganization(testOrganization1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddOrganization(testOrganization2)
	errorChecker.FatalTesting(t, err)
	// add locations
	testLocation1, testLocation2 := getTestLocations()
	_, err = storage.AddLocation(testLocation1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddLocation(testLocation2)
	errorChecker.FatalTesting(t, err)
	// add clinics
	testClinic1, testClinic2 := getTestClinics()
	testClinic1.Organization = swag.String(testOrganization1.ID)
	testClinic1.Location = swag.String(testLocation1.ID)
	_, err = storage.AddClinic(testClinic1)
	errorChecker.FatalTesting(t, err)
	testClinic2.Organization = swag.String(testOrganization2.ID)
	testClinic2.Location = swag.String(testLocation2.ID)
	_, err = storage.AddClinic(testClinic2)
	errorChecker.FatalTesting(t, err)

	// user 1 admin role on global domain
	userRole1 := getTestUserRole(testUser1.ID, authCommon.SuperadminRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(userRole1)
	errorChecker.FatalTesting(t, err)
	// user 1 everyone role on organization 1
	userRole2 := getTestUserRole(testUser1.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeOrganization, testOrganization1.ID)
	_, err = storage.AddUserRole(userRole2)
	errorChecker.FatalTesting(t, err)
	// user 2 everyone role on organization 2
	userRole3 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeOrganization, testOrganization2.ID)
	_, err = storage.AddUserRole(userRole3)
	errorChecker.FatalTesting(t, err)
	// user 1 everyone role on clinic 1
	userRole4 := getTestUserRole(testUser1.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeClinic, testClinic1.ID)
	_, err = storage.AddUserRole(userRole4)
	errorChecker.FatalTesting(t, err)
	// user 1 testRole1 on clinic 1
	userRole5 := getTestUserRole(testUser1.ID, testRole1.ID, authCommon.DomainTypeClinic, testClinic1.ID)
	_, err = storage.AddUserRole(userRole5)
	errorChecker.FatalTesting(t, err)
	// user 2 everyone role on clinic 2
	userRole6 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeClinic, testClinic2.ID)
	_, err = storage.AddUserRole(userRole6)
	errorChecker.FatalTesting(t, err)
	// user 2 testRole1 on clinic 2
	userRole7 := getTestUserRole(testUser2.ID, testRole1.ID, authCommon.DomainTypeClinic, testClinic2.ID)
	_, err = storage.AddUserRole(userRole7)
	errorChecker.FatalTesting(t, err)
	// user 2 testRole2 on clinic 2
	userRole8 := getTestUserRole(testUser2.ID, testRole2.ID, authCommon.DomainTypeClinic, testClinic2.ID)
	_, err = storage.AddUserRole(userRole8)
	errorChecker.FatalTesting(t, err)
	// user 2 testRole2 on location 1 (freely assigned role)
	userRole9 := getTestUserRole(testUser2.ID, testRole2.ID, authCommon.DomainTypeLocation, testLocation1.ID)
	_, err = storage.AddUserRole(userRole9)
	errorChecker.FatalTesting(t, err)
	// user 2 everyone role on organization 1
	userRole10 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeOrganization, testOrganization1.ID)
	_, err = storage.AddUserRole(userRole10)
	errorChecker.FatalTesting(t, err)

	testCases := []struct {
		description           string
		UserID                *string
		RoleID                *string
		DomainType            *string
		DomainID              *string
		expectedNumberOfRoles int
		expectedResultMap     map[string]*models.UserRole
		errorExpected         bool
		exactError            error
	}{
		{
			"Find all userRoles",
			nil,
			nil,
			nil,
			nil,
			16,
			map[string]*models.UserRole{
				userRole1.ID:  userRole1,
				userRole2.ID:  userRole2,
				userRole3.ID:  userRole3,
				userRole4.ID:  userRole4,
				userRole5.ID:  userRole5,
				userRole6.ID:  userRole6,
				userRole7.ID:  userRole7,
				userRole8.ID:  userRole8,
				userRole9.ID:  userRole9,
				userRole10.ID: userRole10,
			},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser1, role: testRole1, domainType: clinic; domainID: testClinic1",
			swag.String(testUser1.ID),
			swag.String(testRole1.ID),
			swag.String(authCommon.DomainTypeClinic),
			swag.String(testClinic1.ID),
			1,
			map[string]*models.UserRole{
				userRole5.ID: userRole5,
			},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser2, domainType: clinic; domainID: testClinic1",
			swag.String(testUser2.ID),
			nil,
			swag.String(authCommon.DomainTypeClinic),
			swag.String(testClinic1.ID),
			0,
			map[string]*models.UserRole{},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser1",
			swag.String(testUser1.ID),
			nil,
			nil,
			nil,
			7,
			map[string]*models.UserRole{
				userRole1.ID: userRole1,
				userRole2.ID: userRole2,
				userRole4.ID: userRole4,
				userRole5.ID: userRole5,
			},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser1, role: authCommon.EveryoneRole",
			swag.String(testUser1.ID),
			swag.String(authCommon.EveryoneRole.ID),
			nil,
			nil,
			3,
			map[string]*models.UserRole{
				userRole2.ID: userRole2,
				userRole4.ID: userRole4,
			},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser2, role: testRole2",
			swag.String(testUser2.ID),
			swag.String(testRole2.ID),
			nil,
			nil,
			2,
			map[string]*models.UserRole{
				userRole8.ID: userRole8,
				userRole9.ID: userRole9,
			},
			false,
			nil,
		},
		{
			"Find userRoles; user: testUser2, role: testRole2; domainType: clinic",
			swag.String(testUser2.ID),
			swag.String(testRole2.ID),
			swag.String(authCommon.DomainTypeClinic),
			nil,
			1,
			map[string]*models.UserRole{
				userRole8.ID: userRole8,
			},
			false,
			nil,
		},
		{
			"Find userRoles for; user: testUser1; domainType: clinic",
			swag.String(testUser1.ID),
			nil,
			swag.String(authCommon.DomainTypeClinic),
			nil,
			2,
			map[string]*models.UserRole{
				userRole4.ID: userRole4,
				userRole5.ID: userRole5,
			},
			false,
			nil,
		},
		{
			"Find userRoles for; user: testUser2; domainType: organization",
			swag.String(testUser2.ID),
			nil,
			swag.String(authCommon.DomainTypeOrganization),
			nil,
			2,
			map[string]*models.UserRole{
				userRole3.ID:  userRole3,
				userRole10.ID: userRole10,
			},
			false,
			nil,
		},
		{
			"Find userRoles for; user: testUser2; domainType: organization; domainId: testOrganization1",
			swag.String(testUser2.ID),
			nil,
			swag.String(authCommon.DomainTypeOrganization),
			swag.String(testOrganization1.ID),
			1,
			map[string]*models.UserRole{
				userRole10.ID: userRole10,
			},
			false,
			nil,
		},
		{
			"Find userRoles for; role: authCommon.EveryoneRole; domainType: clinic; domainId: testClinic1",
			nil,
			swag.String(authCommon.EveryoneRole.ID),
			swag.String(authCommon.DomainTypeClinic),
			swag.String(testClinic1.ID),
			1,
			map[string]*models.UserRole{
				userRole4.ID: userRole4,
			},
			false,
			nil,
		},
		{
			"Find userRoles for; domainType: clinic; domainId: testClinic2",
			nil,
			nil,
			swag.String(authCommon.DomainTypeClinic),
			swag.String(testClinic2.ID),
			3,
			map[string]*models.UserRole{
				userRole6.ID: userRole6,
				userRole7.ID: userRole7,
				userRole8.ID: userRole8,
			},
			false,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			// call FindUserRoles
			out, err := storage.FindUserRoles(test.UserID, test.RoleID, test.DomainType, test.DomainID)

			// create map of the result for comparison
			resultMap := make(map[string]*models.UserRole)
			for _, userRole := range out {
				resultMap[userRole.ID] = userRole
			}

			// check expected results
			// expect 4 more user roles as 2 are added per user on adding user
			if len(resultMap) != test.expectedNumberOfRoles {
				t.Errorf("Expected %d user roles, got %d", test.expectedNumberOfRoles, len(resultMap))
			}
			for id, expectedUserRole := range test.expectedResultMap {
				outUserRole, ok := resultMap[id]
				if !ok {
					fmt.Println("Expected")
					printJson(expectedUserRole)
					t.Errorf("Expected user role with id %s, not present in result", id)
				}
				if !reflect.DeepEqual(outUserRole, expectedUserRole) {
					fmt.Println("Expected")
					printJson(expectedUserRole)
					fmt.Println("Got")
					printJson(outUserRole)
					t.Errorf("Expected user role with id %s to equal\n%+v\ngot\n%+v", id, expectedUserRole, outUserRole)
				}
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			// assert actual error
			if !reflect.DeepEqual(err, test.exactError) {
				t.Errorf("Expected error to equal '%v'; got %v", test.exactError, err)
			}
		})
	}
}

func TestRemoveUserRole(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)

	testUserRole := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)

	// add userRole
	_, err = storage.AddUserRole(testUserRole)
	errorChecker.FatalTesting(t, err)

	// remove user
	err = storage.RemoveUserRole(testUserRole.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// try to get user role by content; should not be there
	_, err = storage.GetUserRoleByContent(*testUserRole.UserID, *testUserRole.RoleID, *testUserRole.DomainType, *testUserRole.DomainID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	// remove user again
	err = storage.RemoveUser(testUserRole.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
