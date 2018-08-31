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

// method to ensure that locations used for tests are always fresh
func getTestLocations() (*models.Location, *models.Location) {
	locationName1 := "Location 1"
	locationName2 := "Location 2"
	testLocation1 := &models.Location{
		City:        "Beirut",
		Country:     "Lebanon",
		Electricity: true,
		Name:        &locationName1,
	}
	testLocation2 := &models.Location{
		City:        "Aleppo",
		Country:     "Syria",
		Electricity: true,
		Name:        &locationName2,
	}

	return testLocation1, testLocation2
}

func TestAddLocation(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testLocation, _ := getTestLocations()

	// add location
	location, err := storage.AddLocation(testLocation)
	if location.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// can't add location with the same name
	_, err = storage.AddLocation(testLocation)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	locations, err := storage.GetLocations()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(locations) != 1 {
		t.Fatalf("Expected 1 location; got %d", len(locations))
	}
}

func TestGetLocation(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testLocation, _ := getTestLocations()

	// add location
	_, err := storage.AddLocation(testLocation)
	errorChecker.FatalTesting(t, err)

	// get location
	location, err := storage.GetLocation(testLocation.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testLocation, *location) {
		t.Fatalf("Expected returned role to be '%v', got '%v'", *testLocation, *location)
	}

	// get location with wrong uuid
	_, err = storage.GetLocation("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing location
	_, err = storage.GetLocation("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetLocations(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testLocation1, testLocation2 := getTestLocations()

	// add locations
	_, err := storage.AddLocation(testLocation1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddLocation(testLocation2)
	errorChecker.FatalTesting(t, err)

	// get locations
	locations, err := storage.GetLocations()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(locations) != 2 {
		t.Fatalf("Expected 2 locations; got %d", len(locations))
	}

	locationsMap := map[string]*models.Location{}
	for _, location := range locations {
		locationsMap[location.ID] = location
	}

	if !reflect.DeepEqual(*testLocation1, *locationsMap[testLocation1.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testLocation1, *locationsMap[testLocation1.ID])
	}

	if !reflect.DeepEqual(*testLocation2, *locationsMap[testLocation2.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testLocation2, *locationsMap[testLocation2.ID])
	}
}

func TestGetLocationClinics(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testLocation, _ := getTestLocations()
	_, err := storage.AddLocation(testLocation)
	errorChecker.FatalTesting(t, err)

	testOrganization, _ := getTestOrganizations()
	_, err = storage.AddOrganization(testOrganization)
	errorChecker.FatalTesting(t, err)
	testClinic1, testClinic2 := getTestClinics()
	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation.ID
	testClinic2.Organization = &testOrganization.ID
	testClinic2.Location = &testLocation.ID
	_, err = storage.AddClinic(testClinic1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddClinic(testClinic2)
	errorChecker.FatalTesting(t, err)

	expectedClinicsMap := map[string]*models.Clinic{
		testClinic1.ID: testClinic1,
		testClinic2.ID: testClinic2,
	}

	// get location clinics
	locationClinics, err := storage.GetLocationClinics(testLocation.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(locationClinics) != 2 {
		t.Fatalf("Expected 2 clinics returned; got %d", len(locationClinics))
	}
	for _, returnedClinic := range locationClinics {
		expectedClinic, ok := expectedClinicsMap[returnedClinic.ID]
		if !ok {
			t.Fatalf("Clinic with ID %s was not expected", returnedClinic.ID)
		}
		if !reflect.DeepEqual(expectedClinic, returnedClinic) {
			t.Fatalf("Expected clinic to be '%v', got '%v'", expectedClinic, returnedClinic)
		}
	}

	// get location clinics with wrong uuid
	_, err = storage.GetLocationClinics("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get location clinics of non existing location
	_, err = storage.GetLocationClinics("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetLocationOrganizationIDs(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testLocation, _ := getTestLocations()
	_, err := storage.AddLocation(testLocation)
	errorChecker.FatalTesting(t, err)

	testOrganization1, testOrganization2 := getTestOrganizations()
	_, err = storage.AddOrganization(testOrganization1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddOrganization(testOrganization2)
	errorChecker.FatalTesting(t, err)
	testClinic1, testClinic2 := getTestClinics()
	testClinic1.Organization = &testOrganization1.ID
	testClinic1.Location = &testLocation.ID
	testClinic2.Organization = &testOrganization2.ID
	testClinic2.Location = &testLocation.ID
	_, err = storage.AddClinic(testClinic1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddClinic(testClinic2)
	errorChecker.FatalTesting(t, err)
	testClinic3Name := "testClinic3"
	testClinic3 := &models.Clinic{
		Name:         &testClinic3Name,
		Location:     &testLocation.ID,
		Organization: &testOrganization2.ID,
	}
	_, err = storage.AddClinic(testClinic3)
	errorChecker.FatalTesting(t, err)

	// each organization should be returned once
	expectedOrganizations := []string{testOrganization1.ID, testOrganization2.ID}

	// get location organizations
	locationOrganizations, err := storage.GetLocationOrganizationIDs(testLocation.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(locationOrganizations) != 2 {
		t.Fatalf("Expected 2 organizations returned; got %d", len(locationOrganizations))
	}
	if len(utils.DiffSlice(locationOrganizations, expectedOrganizations)) != 0 || len(utils.DiffSlice(expectedOrganizations, locationOrganizations)) != 0 {
		t.Fatalf("Expected '%v', got '%v'", expectedOrganizations, locationOrganizations)
	}

	// get location organizations with wrong uuid
	_, err = storage.GetLocationOrganizationIDs("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get location organizations of non existing location
	_, err = storage.GetLocationOrganizationIDs("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestUpdateLocation(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testLocation1, testLocation2 := getTestLocations()

	// add locations
	_, err := storage.AddLocation(testLocation1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddLocation(testLocation2)
	errorChecker.FatalTesting(t, err)

	// update location
	updateLocation := &models.Location{
		ID:          testLocation1.ID,
		Name:        testLocation1.Name,
		City:        testLocation1.City,
		Country:     testLocation1.Country,
		Electricity: testLocation1.Electricity,
		WaterSupply: true,
	}
	location, err := storage.UpdateLocation(updateLocation)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*location, *updateLocation) {
		t.Fatalf("Expected location one to be '%v', got '%v'", *location, *updateLocation)
	}

	// can't update location with the name of another location
	updateLocation = &models.Location{
		ID:          testLocation1.ID,
		Name:        testLocation2.Name,
		City:        testLocation1.City,
		Country:     testLocation1.Country,
		Electricity: testLocation1.Electricity,
	}
	_, err = storage.UpdateLocation(updateLocation)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveLocation(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testLocation, _ := getTestLocations()
	_, err := storage.AddLocation(testLocation)
	errorChecker.FatalTesting(t, err)

	// add organization & clinic to check if clinic gets removed on location removal
	testOrganization, _ := getTestOrganizations()
	_, err = storage.AddOrganization(testOrganization)
	errorChecker.FatalTesting(t, err)
	testClinic, _ := getTestClinics()
	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	_, err = storage.AddClinic(testClinic)
	errorChecker.FatalTesting(t, err)

	// add user roles to test if they are removed properly with location
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	testUser1, testUser2 := getTestUsers()
	_, err = storage.AddUser(testUser1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	testUserRole1 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(testUserRole1)
	errorChecker.FatalTesting(t, err)
	testUserRole2 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeOrganization, testOrganization.ID)
	_, err = storage.AddUserRole(testUserRole2)
	errorChecker.FatalTesting(t, err)
	testUserRole3 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeLocation, testLocation.ID)
	_, err = storage.AddUserRole(testUserRole3)
	errorChecker.FatalTesting(t, err)
	testUserRole4 := getTestUserRole(testUser1.ID, testRole.ID, authCommon.DomainTypeClinic, testClinic.ID)
	_, err = storage.AddUserRole(testUserRole4)
	errorChecker.FatalTesting(t, err)
	testUserRole5 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeClinic, testClinic.ID)
	_, err = storage.AddUserRole(testUserRole5)
	errorChecker.FatalTesting(t, err)

	// remove location
	err = storage.RemoveLocation(testLocation.ID)
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
	if len(userRoles) != 8 {
		if err == nil {
			t.Fatalf("Expected 8 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation.ID))
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
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization.ID))
	if len(userRoles) != 1 {
		if err == nil {
			t.Fatalf("Expected 1 user role; got %d", len(userRoles))
		}
	}

	// remove location again
	err = storage.RemoveLocation(testLocation.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
