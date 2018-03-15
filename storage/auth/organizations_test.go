package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/utils"
)

// method to ensure that organizations used for tests are always fresh
func getTestOrganizations() (*models.Organization, *models.Organization) {
	orgName1 := "Organization1"
	orgName2 := "Organization2"
	testOrganization1 := &models.Organization{
		Address: &models.Address{
			Country: "USA",
			City:    "New Jersey",
		},
		LegalStatus: "NGO",
		Name:        &orgName1,
	}
	testOrganization2 := &models.Organization{
		Address: &models.Address{
			Country: "Slovenia",
			City:    "Ljubljana",
		},
		LegalStatus: "NGO",
		Name:        &orgName2,
	}
	return testOrganization1, testOrganization2
}

func TestAddOrganization(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization1, _ := getTestOrganizations()

	// add organization
	organization, err := storage.AddOrganization(testOrganization1)
	if organization.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// can't add organization with the same name
	_, err = storage.AddOrganization(testOrganization1)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	organizations, err := storage.GetOrganizations()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(organizations) != 1 {
		t.Fatalf("Expected 1 organizations; got %d", len(organizations))
	}
}

func TestGetOrganization(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization1, _ := getTestOrganizations()
	storage.AddOrganization(testOrganization1)

	// get organization
	organization, err := storage.GetOrganization(testOrganization1.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testOrganization1, *organization) {
		t.Fatalf("Expected returned role to be '%v', got '%v'", *testOrganization1, *organization)
	}

	// get organization with wrong uuid
	_, err = storage.GetOrganization("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing organization
	_, err = storage.GetOrganization("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetOrganizations(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization1, testOrganization2 := getTestOrganizations()
	storage.AddOrganization(testOrganization1)
	storage.AddOrganization(testOrganization2)

	// get organizations
	organizations, err := storage.GetOrganizations()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(organizations) != 2 {
		t.Fatalf("Expected 2 organizations; got %d", len(organizations))
	}

	organizationsMap := map[string]*models.Organization{}
	for _, organization := range organizations {
		organizationsMap[organization.ID] = organization
	}

	if !reflect.DeepEqual(*testOrganization1, *organizationsMap[testOrganization1.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testOrganization1, *organizationsMap[testOrganization1.ID])
	}

	if !reflect.DeepEqual(*testOrganization2, *organizationsMap[testOrganization2.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testOrganization2, *organizationsMap[testOrganization2.ID])
	}
}

func TestGetOrganizationClinics(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testOrganization, _ := getTestOrganizations()
	testLocation, _ := getTestLocations()
	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation)

	testClinic1, testClinic2 := getTestClinics()
	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation.ID
	testClinic2.Organization = &testOrganization.ID
	testClinic2.Location = &testLocation.ID
	storage.AddClinic(testClinic1)
	storage.AddClinic(testClinic2)

	expectedClinicsMap := map[string]*models.Clinic{
		testClinic1.ID: testClinic1,
		testClinic2.ID: testClinic2,
	}

	// get organization clinics
	organizationClinics, err := storage.GetOrganizationClinics(testOrganization.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(organizationClinics) != 2 {
		t.Fatalf("Expected 2 clinics returned; got %d", len(organizationClinics))
	}
	for _, returnedClinic := range organizationClinics {
		expectedClinic, ok := expectedClinicsMap[returnedClinic.ID]
		if !ok {
			t.Fatalf("Clinic with ID %s was not expected", returnedClinic.ID)
		}
		if !reflect.DeepEqual(expectedClinic, returnedClinic) {
			t.Fatalf("Expected clinic to be '%v', got '%v'", expectedClinic, returnedClinic)
		}
	}

	// get organization clinics with wrong uuid
	_, err = storage.GetOrganizationClinics("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get organization clinics of non existing organization
	organizationClinics, err = storage.GetOrganizationClinics("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetOrganizationLocationIDs(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testOrganization, _ := getTestOrganizations()
	testLocation1, testLocation2 := getTestLocations()
	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation1)
	storage.AddLocation(testLocation2)

	testClinic1, testClinic2 := getTestClinics()
	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation1.ID
	testClinic2.Organization = &testOrganization.ID
	testClinic2.Location = &testLocation1.ID
	storage.AddClinic(testClinic1)
	storage.AddClinic(testClinic2)
	testClinic3 := &models.Clinic{
		Name:         testClinic1.Name,
		Location:     &testLocation2.ID,
		Organization: &testOrganization.ID,
	}
	storage.AddClinic(testClinic3)

	// each location should be returned once
	expectedLocations := []string{testLocation1.ID, testLocation2.ID}

	// get organization locations
	organizationLocations, err := storage.GetOrganizationLocationIDs(testOrganization.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(organizationLocations) != 2 {
		t.Fatalf("Expected 2 locations returned; got %d", len(organizationLocations))
	}
	if len(utils.DiffSlice(organizationLocations, expectedLocations)) != 0 || len(utils.DiffSlice(expectedLocations, organizationLocations)) != 0 {
		t.Fatalf("Expected '%v', got '%v'", expectedLocations, organizationLocations)
	}

	// get organization locations with wrong uuid
	_, err = storage.GetOrganizationLocationIDs("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get organization locations of non existing organization
	organizationLocations, err = storage.GetOrganizationLocationIDs("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestUpdateOrganization(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization1, testOrganization2 := getTestOrganizations()
	storage.AddOrganization(testOrganization1)
	storage.AddOrganization(testOrganization2)

	// update organization
	updatedName := "updatedName"
	updateOrganization := &models.Organization{
		Name: &updatedName,
		ID:   testOrganization1.ID,
	}
	organization, err := storage.UpdateOrganization(updateOrganization)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*organization, *updateOrganization) {
		t.Fatalf("Expected organization one to be '%v', got '%v'", *organization, *updateOrganization)
	}

	// can't update organization with the name of another organization
	updateOrganization = &models.Organization{
		Name: testOrganization2.Name,
		ID:   testOrganization1.ID,
	}
	_, err = storage.UpdateOrganization(updateOrganization)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveOrganization(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	storage.AddOrganization(testOrganization)

	// add location & clinic to check if clinic gets removed on organization
	testLocation, _ := getTestLocations()
	storage.AddLocation(testLocation)
	testClinic, _ := getTestClinics()
	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	storage.AddClinic(testClinic)

	// add user roles to test if they are removed properly with organization
	testRole, _ := getTestRoles()
	storage.AddRole(testRole)
	testUser1, testUser2 := getTestUsers()
	storage.AddUser(testUser1)
	storage.AddUser(testUser2)
	testUserRole1 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	storage.AddUserRole(testUserRole1)
	testUserRole2 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeOrganization, testOrganization.ID)
	storage.AddUserRole(testUserRole2)
	testUserRole3 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeOrganization, testOrganization.ID)
	storage.AddUserRole(testUserRole3)
	testUserRole4 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeClinic, testClinic.ID)
	storage.AddUserRole(testUserRole4)
	testUserRole5 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeLocation, testLocation.ID)
	storage.AddUserRole(testUserRole5)

	// remove organization
	err := storage.RemoveOrganization(testOrganization.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	// check if clinic was removed
	clinics, _ := storage.GetClinics()
	if len(clinics) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 clinics; got %d", len(clinics))
		}
	}
	// check if user roles were removed
	userRoles, _ := storage.GetUserRoles()
	if len(userRoles) != 6 {
		if err == nil {
			t.Fatalf("Expected 6 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization.ID))
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeClinic), swag.String(testClinic.ID))
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation.ID))
	if len(userRoles) != 1 {
		if err == nil {
			t.Fatalf("Expected 1 user role; got %d", len(userRoles))
		}
	}

	// remove organization again
	err = storage.RemoveOrganization(testOrganization.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
