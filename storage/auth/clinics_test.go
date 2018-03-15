package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/utils"
	authCommon "github.com/iryonetwork/wwm/auth"
)

// method to ensure that clinics used for tests are always fresh
func getTestClinics() (*models.Clinic, *models.Clinic) {
	clinicName1 := "Clinic 1"
	clinicName2 := "Clinic 2"
	testClinic1 := &models.Clinic{
		Name: &clinicName1,
	}
	testClinic2 := &models.Clinic{
		Name: &clinicName2,
	}

	return testClinic1, testClinic2
}

func TestAddClinic(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	testLocation1, testLocation2 := getTestLocations()
	testClinic1, _ := getTestClinics()

	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation1)
	storage.AddLocation(testLocation2)

	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation1.ID

	// add clinic
	clinic, err := storage.AddClinic(testClinic1)
	if clinic.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// check if clinic was added to organization
	organization, err := storage.GetOrganization(testOrganization.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(organization.Clinics) != 1 {
		t.Fatalf("Expected number of clinics of the organization to be 1; got %d", len(organization.Clinics))
	}
	if organization.Clinics[0] != clinic.ID {
		t.Fatalf("Expected clinic added to the organization to have ID '%s'; got '%s'", clinic.ID, organization.Clinics[0])
	}
	// check if clinic was added to location
	location, err := storage.GetLocation(testLocation1.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(location.Clinics) != 1 {
		t.Fatalf("Expected number of clinics of the location to be 1; got %d", len(location.Clinics))
	}
	if location.Clinics[0] != clinic.ID {
		t.Fatalf("Expected clinic added to the location to have ID '%s'; got '%s'", clinic.ID, location.Clinics[0])
	}

	// can't add clinic with the same name
	_, err = storage.AddClinic(testClinic1)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	clinics, err := storage.GetClinics()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(clinics) != 1 {
		t.Fatalf("Expected 1 clinics; got %d", len(clinics))
	}
}

func TestGetClinic(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	testLocation, _ := getTestLocations()
	testClinic, _ := getTestClinics()

	// populate DB
	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation)

	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	storage.AddClinic(testClinic)

	// get clinic
	clinic, err := storage.GetClinic(testClinic.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testClinic, *clinic) {
		t.Fatalf("Expected returned role to be '%v', got '%v'", *testClinic, *clinic)
	}

	// get clinic with wrong uuid
	_, err = storage.GetClinic("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing clinic
	_, err = storage.GetClinic("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetClinics(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	testLocation1, testLocation2 := getTestLocations()
	testClinic1, testClinic2 := getTestClinics()

	// populate DB
	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation1)
	storage.AddLocation(testLocation2)

	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation1.ID
	testClinic2.Organization = &testOrganization.ID
	testClinic2.Location = &testLocation2.ID

	// add clinics
	storage.AddClinic(testClinic1)
	storage.AddClinic(testClinic2)

	// get clinics
	clinics, err := storage.GetClinics()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(clinics) != 2 {
		t.Fatalf("Expected 2 clinics; got %d", len(clinics))
	}

	clinicsMap := map[string]*models.Clinic{}
	for _, clinic := range clinics {
		clinicsMap[clinic.ID] = clinic
	}

	if !reflect.DeepEqual(*testClinic1, *clinicsMap[testClinic1.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testClinic1, *clinicsMap[testClinic1.ID])
	}

	if !reflect.DeepEqual(*testClinic2, *clinicsMap[testClinic2.ID]) {
		t.Fatalf("Expected role one to be '%v', got '%v'", *testClinic2, *clinicsMap[testClinic2.ID])
	}
}

func TestGetClinicOrganization(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testLocation, _ := getTestLocations()
	storage.AddLocation(testLocation)
	testOrganization, _ := getTestOrganizations()
	storage.AddOrganization(testOrganization)

	testClinic, _ := getTestClinics()
	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	storage.AddClinic(testClinic)

	expectedOrganization, _ := storage.GetOrganization(testOrganization.ID)

	// get clinic organization
	clinicOrganization, err := storage.GetClinicOrganization(testClinic.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(expectedOrganization, clinicOrganization) {
		t.Fatalf("Expected user to be '%v', got '%v'", expectedOrganization, clinicOrganization)
	}

	// get clinic organization with wrong uuid
	_, err = storage.GetClinicOrganization("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get clinic organization of non existing clinic
	_, err = storage.GetClinicLocation("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetClinicLocation(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// populate DB
	testLocation, _ := getTestLocations()
	storage.AddLocation(testLocation)
	testOrganization, _ := getTestOrganizations()
	storage.AddOrganization(testOrganization)

	testClinic, _ := getTestClinics()
	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	storage.AddClinic(testClinic)

	expectedLocation, _ := storage.GetLocation(testLocation.ID)

	// get clinic location
	clinicLocation, err := storage.GetClinicLocation(testClinic.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(expectedLocation, clinicLocation) {
		t.Fatalf("Expected user to be '%v', got '%v'", expectedLocation, clinicLocation)
	}

	// get clinic location with wrong uuid
	_, err = storage.GetClinicLocation("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get clinic location of non existing clinic
	_, err = storage.GetClinicLocation("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestUpdateClinic(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	testLocation1, testLocation2 := getTestLocations()
	testClinic1, testClinic2 := getTestClinics()

	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation1)
	storage.AddLocation(testLocation2)

	testClinic1.Organization = &testOrganization.ID
	testClinic1.Location = &testLocation1.ID
	testClinic2.Organization = &testOrganization.ID
	testClinic2.Location = &testLocation2.ID
	storage.AddClinic(testClinic1)
	storage.AddClinic(testClinic2)

	// update clinic
	updateClinic := &models.Clinic{
		Name:         testClinic1.Name,
		ID:           testClinic1.ID,
		Location:     &testLocation2.ID,
		Organization: testClinic1.Organization,
	}

	clinic, err := storage.UpdateClinic(updateClinic)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*clinic, *updateClinic) {
		t.Fatalf("Expected clinic one to be '%v', got '%v'", *clinic, *updateClinic)
	}

	// check if clinic was removed from testLocation1 and added to testLocation2
	location, _ := storage.GetLocation(testLocation1.ID)
	if len(location.Clinics) != 0 {
		t.Fatalf("Expected 0 clinics, got %d", len(location.Clinics))
	}
	location, _ = storage.GetLocation(testLocation2.ID)
	if len(location.Clinics) != 2 {
		t.Fatalf("Expected 2 clinics, got %d", len(location.Clinics))
	}

	// can update clinic with the name of another clinic if they are at different location & organization
	updateClinic = &models.Clinic{
		Name:         testClinic2.Name,
		ID:           testClinic1.ID,
		Location:     testClinic1.Location,
		Organization: testClinic1.Organization,
	}
	_, err = storage.UpdateClinic(updateClinic)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// can't update clinic with the name of another clinic if they are at the same time location & organization
	updateClinic = &models.Clinic{
		Name:         testClinic2.Name,
		ID:           testClinic1.ID,
		Location:     testClinic2.Location,
		Organization: testClinic2.Organization,
	}
	_, err = storage.UpdateClinic(updateClinic)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveClinic(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testOrganization, _ := getTestOrganizations()
	testLocation, _ := getTestLocations()
	testClinic, _ := getTestClinics()

	storage.AddOrganization(testOrganization)
	storage.AddLocation(testLocation)

	testClinic.Organization = &testOrganization.ID
	testClinic.Location = &testLocation.ID
	storage.AddClinic(testClinic)

	// add user roles to test if they are removed properly with clinic
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
	testUserRole5 := getTestUserRole(testUser2.ID, authCommon.EveryoneRole.ID, authCommon.DomainTypeClinic, testClinic.ID)
	storage.AddUserRole(testUserRole5)

	// remove clinic
	err := storage.RemoveClinic(testClinic.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	// check if clinic was removed from test organization
	organization, _ := storage.GetOrganization(testOrganization.ID)
	if len(organization.Clinics) != 0 {
		t.Fatalf("Expected 0 clinics, got %d", len(organization.Clinics))
	}
	// check if clinic was removed from test location
	location, _ := storage.GetLocation(testLocation.ID)
	if len(location.Clinics) != 0 {
		t.Fatalf("Expected 0 clinics, got %d", len(location.Clinics))
	}
	// check if user roles were removed
	userRoles, _ := storage.GetUserRoles()
	if len(userRoles) != 7 {
		if err == nil {
			t.Fatalf("Expected 7 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeClinic), swag.String(testClinic.ID))
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}

	// remove clinic again
	err = storage.RemoveClinic(testClinic.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
