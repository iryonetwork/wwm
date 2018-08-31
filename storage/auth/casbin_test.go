package auth

import (
	"fmt"
	"testing"

	"github.com/go-openapi/swag"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
)

func TestRules(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	// add locations, organization, clinics
	testLocation, _ := storage.AddLocation(&models.Location{Name: swag.String("Test location")})
	testOrganization1, _ := storage.AddOrganization(&models.Organization{Name: swag.String("Test organization")})
	testOrganization2, _ := storage.AddOrganization(&models.Organization{Name: swag.String("Test organization 2")})
	testClinic1, _ := storage.AddClinic(&models.Clinic{Name: swag.String("Test clinic 1"), Location: &testLocation.ID, Organization: &testOrganization1.ID})
	testClinic2, _ := storage.AddClinic(&models.Clinic{Name: swag.String("Test clinic 2"), Location: &testLocation.ID, Organization: &testOrganization2.ID})

	// add users
	u1, _ := storage.AddUser(&models.User{Username: swag.String("user1")})
	u2, _ := storage.AddUser(&models.User{Username: swag.String("user2")})
	u3, _ := storage.AddUser(&models.User{Username: swag.String("user3")})
	u4, _ := storage.AddUser(&models.User{Username: swag.String("user4")})

	// add roles
	adminRole, _ := storage.AddRole(&models.Role{Name: swag.String("adminRole")})
	doctorRole, _ := storage.AddRole(&models.Role{Name: swag.String("doctorRole")})
	nurseRole, _ := storage.AddRole(&models.Role{Name: swag.String("nurseRole")})

	// add rules
	_, err := storage.AddRule(&models.Rule{
		Subject:  &authCommon.EveryoneRole.ID,
		Action:   swag.Int64(Write),
		Resource: swag.String("/auth/login"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Write),
		Subject:  swag.String(adminRole.ID),
		Resource: swag.String("/clinic/login"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Write),
		Subject:  swag.String(doctorRole.ID),
		Resource: swag.String("/clinic/login"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Write),
		Subject:  swag.String(nurseRole.ID),
		Resource: swag.String("/clinic/login"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(adminRole.ID),
		Resource: swag.String("/frontend/admin*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(adminRole.ID),
		Resource: swag.String("/storage/file*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(adminRole.ID),
		Resource: swag.String("/storage/file/basicInfo"),
		Deny:     true,
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read),
		Subject:  swag.String(doctorRole.ID),
		Resource: swag.String("/frontend/doctor*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(doctorRole.ID),
		Resource: swag.String("/frontend/diagnosis*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(doctorRole.ID),
		Resource: swag.String("/storage/file*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read),
		Subject:  swag.String(nurseRole.ID),
		Resource: swag.String("/frontend/nurse*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read),
		Subject:  swag.String(nurseRole.ID),
		Resource: swag.String("/frontend/diagnosis*"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Read | Write),
		Subject:  swag.String(nurseRole.ID),
		Resource: swag.String("/storage/file/basicInfo"),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(&models.Rule{
		Action:   swag.Int64(Update | Delete),
		Subject:  swag.String(authCommon.AuthorRole.ID),
		Resource: swag.String("/storage/file*"),
	})
	errorChecker.FatalTesting(t, err)

	// add user roles
	// add user 1 to both organizations
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u1.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization1.ID),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u1.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization2.ID),
	})
	// add user 2 to both organizations
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u2.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization1.ID),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u2.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization2.ID),
	})
	// add user 3 to both organizations
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u3.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization1.ID),
	})
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u3.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization2.ID),
	})
	// add user 4 to both test organization 2
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u4.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeOrganization),
		DomainID:   swag.String(testOrganization2.ID),
	})
	// give user1 global admin role
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u1.ID),
		RoleID:     swag.String(adminRole.ID),
		DomainType: swag.String(authCommon.DomainTypeGlobal),
		DomainID:   swag.String(authCommon.DomainIDWildcard),
	})
	// give user2 global doctor role
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u2.ID),
		RoleID:     swag.String(doctorRole.ID),
		DomainType: swag.String(authCommon.DomainTypeGlobal),
		DomainID:   swag.String(authCommon.DomainIDWildcard),
	})
	// give user3 global nurse role
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u3.ID),
		RoleID:     swag.String(nurseRole.ID),
		DomainType: swag.String(authCommon.DomainTypeGlobal),
		DomainID:   swag.String(authCommon.DomainIDWildcard),
	})
	//give user1 doctor role at clinic 1
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u1.ID),
		RoleID:     swag.String(doctorRole.ID),
		DomainType: swag.String(authCommon.DomainTypeClinic),
		DomainID:   swag.String(testClinic1.ID),
	})
	// give user2 admin role at clinic 1
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u2.ID),
		RoleID:     swag.String(adminRole.ID),
		DomainType: swag.String(authCommon.DomainTypeClinic),
		DomainID:   swag.String(testClinic1.ID),
	})
	// give user4 doctorRole at clinic 2
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUserRole(&models.UserRole{
		UserID:     swag.String(u4.ID),
		RoleID:     swag.String(doctorRole.ID),
		DomainType: swag.String(authCommon.DomainTypeClinic),
		DomainID:   swag.String(testClinic2.ID),
	})
	errorChecker.FatalTesting(t, err)

	// validations to be checked
	commonValidations := []*models.ValidationPair{
		// available for everyone
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/auth/login"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// holding a role at clinic needed
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/clinic/login"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// holding a role at clinic needed
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/clinic/login"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// admin role at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/admin/dashboard"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// admin role at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/admin/dashboard"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// doctor role at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/doctor/dashboard"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// doctor role at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/doctor/dashboard"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// nurse at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/nurse"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// nurse at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/nurse"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// doctor or nurse at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/diagnosis"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// doctor at clinic needed
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/frontend/diagnosis"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic1.ID),
		},
		// doctor or nurse at clinic needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/frontend/diagnosis"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// doctor at clinic needed
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/frontend/diagnosis"),
			DomainType: swag.String(authCommon.DomainTypeClinic),
			DomainID:   swag.String(testClinic2.ID),
		},
		// doctor or admin role at location needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/storage/file/something"),
			DomainType: swag.String(authCommon.DomainTypeLocation),
			DomainID:   swag.String(testLocation.ID),
		},
		// doctor role at location needed
		{
			Actions:    swag.Int64(Write),
			Resource:   swag.String("/storage/file/vitalSign"),
			DomainType: swag.String(authCommon.DomainTypeLocation),
			DomainID:   swag.String(testLocation.ID),
		},
		// doctor or nurse location needed
		{
			Actions:    swag.Int64(Read),
			Resource:   swag.String("/storage/file/basicInfo"),
			DomainType: swag.String(authCommon.DomainTypeLocation),
			DomainID:   swag.String(testLocation.ID),
		},
		// author role over user ID domain needed
		{
			Actions:    swag.Int64(Delete),
			Resource:   swag.String("/storage/file/mine"),
			DomainType: swag.String(authCommon.DomainTypeUser),
			DomainID:   swag.String(u1.ID),
		},
		// author role over user ID domain needed
		{
			Actions:    swag.Int64(Update),
			Resource:   swag.String("/storage/file/notmine"),
			DomainType: swag.String(authCommon.DomainTypeUser),
			DomainID:   swag.String(u2.ID),
		},
		// author role over user ID domain needed
		{
			Actions:    swag.Int64(Update),
			Resource:   swag.String("/storage/file/notmine"),
			DomainType: swag.String(authCommon.DomainTypeUser),
			DomainID:   swag.String(u3.ID),
		},
		// author role over user ID domain needed
		{
			Actions:    swag.Int64(Update),
			Resource:   swag.String("/storage/file/notmine"),
			DomainType: swag.String(authCommon.DomainTypeUser),
			DomainID:   swag.String(u4.ID),
		},
	}

	tests := []struct {
		userID      string
		validations []*models.ValidationPair
		results     []bool
	}{
		{
			userID:      u1.ID,
			validations: commonValidations,
			results: []bool{
				true,
				true,
				true,
				true,
				true,
				true,
				false,
				false,
				false,
				true,
				true,
				false,
				false,
				true,
				true,
				false,
				true,
				false,
				false,
				false,
			},
		},
		{
			userID:      u2.ID,
			validations: commonValidations,
			results: []bool{
				true,
				true,
				true,
				true,
				false,
				true,
				true,
				false,
				false,
				true,
				true,
				true,
				true,
				true,
				true,
				false,
				false,
				true,
				false,
				false,
			},
		},
		{
			userID:      u3.ID,
			validations: commonValidations,
			results: []bool{
				true,
				true,
				true,
				false,
				false,
				false,
				false,
				true,
				true,
				true,
				false,
				true,
				false,
				false,
				false,
				true,
				false,
				false,
				true,
				false,
			},
		},
		{
			userID:      u4.ID,
			validations: commonValidations,
			results: []bool{
				true,
				false,
				true,
				false,
				false,
				false,
				true,
				false,
				false,
				false,
				false,
				true,
				true,
				true,
				true,
				true,
				false,
				false,
				false,
				true,
			},
		},
	}

	for testIndex, test := range tests {
		errorChecker.FatalTesting(t, storage.enforcer.LoadPolicy())

		results := storage.FindACL(test.userID, test.validations)

		for i, res := range results {
			if *res.Result != test.results[i] {
				fmt.Println(test.userID)
				printJson(test.validations[i])
				t.Fatalf("Test %d; validation %d: Expected validation '%v' to be %t; got %t", testIndex, i, test.validations[i], test.results[i], *res.Result)
			}
		}
	}
}
