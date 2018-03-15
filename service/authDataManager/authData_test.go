package authDataManager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	//uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/service/authDataManager/mock"
	//"github.com/iryonetwork/wwm/utils"
)

var (
	testUser1 = &models.User{
		ID:       "97b88406-c893-4c3e-aa1f-3a93a09b6d3d",
		Username: swag.String("testUser1"),
	}
	testUser2 = &models.User{
		ID:       "f78cb507-e8ea-4f1e-af13-4d3e148e8dc8",
		Username: swag.String("testUser2"),
	}
	testUser3 = &models.User{
		ID:       "d4bfca48-74ed-4381-808a-11a92d23fe55",
		Username: swag.String("testUser3"),
	}
	testOrganization1 = &models.Organization{
		ID:   "f8c5c01f-f782-4ac2-9d85-55c35ed81f34",
		Name: swag.String("testOrganization1"),
	}
	testOrganization2 = &models.Organization{
		ID:   "6cd71934-ada1-4035-a2b7-e9fa83553c66",
		Name: swag.String("testOrganization2"),
	}
	testLocation1 = &models.Location{
		ID:   "18cf23c2-f231-4d1a-b0af-1427af620430",
		Name: swag.String("testLocation1"),
	}
	testLocation2 = &models.Location{
		ID:   "bbba4fe0-837c-45b2-bdc1-b917b27fb52e",
		Name: swag.String("testLocation2"),
	}
	testClinic1 = &models.Clinic{
		ID:           "07be9ea9-1479-4c91-9470-336c74a4fd79",
		Organization: swag.String(testOrganization1.ID),
		Location:     swag.String(testLocation1.ID),
		Name:         swag.String("testClinic1"),
	}
	testClinic2 = &models.Clinic{
		ID:           "274e5857-3a7a-46cd-ac99-072d11a0b555",
		Organization: swag.String(testOrganization1.ID),
		Location:     swag.String(testLocation2.ID),
		Name:         swag.String("testClinic2"),
	}
	testClinic3 = &models.Clinic{
		ID:           "3740beaf-b4d2-4918-a07a-4dbc9b186276",
		Organization: swag.String(testOrganization2.ID),
		Location:     swag.String(testLocation1.ID),
		Name:         swag.String("testClinic3"),
	}
	testClinic4 = &models.Clinic{
		ID:           "1b952067-ab1f-40c4-82b8-cfda6d8adf72",
		Organization: swag.String(testOrganization2.ID),
		Location:     swag.String(testLocation2.ID),
		Name:         swag.String("tesClinic4"),
	}
	noErrors   = false
	withErrors = true
)

func TestUserRoleIDs(t *testing.T) {
	testCases := []struct {
		description   string
		userID        string
		domainType    *string
		domainID      *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all user's roles",
			testUser1.ID,
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's global everyone role
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.EveryoneRole.ID),
						DomainType: swag.String(authCommon.DomainTypeGlobal),
						DomainID:   swag.String(authCommon.DomainIDWildcard),
					},
					// user's member role at test organization 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization1.ID),
					},
					// user's member role at test organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
					// user's member role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, nil, nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{authCommon.EveryoneRole.ID, authCommon.MemberRole.ID},
		},
		{
			"Succesfully fetched all user's roles with domain filtering",
			testUser1.ID,
			&authCommon.DomainTypeOrganization,
			&testOrganization2.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's member role at test organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, &authCommon.DomainTypeOrganization, &testOrganization2.ID).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{authCommon.MemberRole.ID},
		},
		{
			"Error on finding user roles",
			testUser1.ID,
			&authCommon.DomainTypeOrganization,
			&testOrganization1.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, &authCommon.DomainTypeOrganization, &testOrganization1.ID).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.UserRoleIDs(context.TODO(), test.userID, test.domainType, test.domainID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func TestUserOrganizationIDs(t *testing.T) {
	testCases := []struct {
		description   string
		userID        string
		roleID        *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all user's organizations",
			testUser1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's member role at test organization 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization1.ID),
					},
					// user's member role at test organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
					// user's superadmin role at test organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, swag.String(authCommon.DomainTypeOrganization), nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testOrganization1.ID, testOrganization2.ID},
		},
		{
			"Succesfully fetched all user's organizations with role filtering",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's superadmin role at test organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeOrganization), nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testOrganization2.ID},
		},
		{
			"Error on finding user roles",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeOrganization), nil).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.UserOrganizationIDs(context.TODO(), test.userID, test.roleID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func TestUserClinicIDs(t *testing.T) {
	testCases := []struct {
		description   string
		userID        string
		roleID        *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all user's clinics",
			testUser1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's member role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
					// user's member role at test clinic 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic2.ID),
					},
					// user's superadmin role at test clinic 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, swag.String(authCommon.DomainTypeClinic), nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testClinic1.ID, testClinic2.ID},
		},
		{
			"Succesfully fetched all user's clinics with role filtering",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// user's superadmin role at test clinic 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeClinic), nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testClinic2.ID},
		},
		{
			"Error on finding user roles",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeClinic), nil).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.UserClinicIDs(context.TODO(), test.userID, test.roleID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})

	}
}

func TestUserLocationIDs(t *testing.T) {
	testCases := []struct {
		description   string
		userID        string
		roleID        *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all user's locations",
			testUser1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// user's superadming role at test location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}
				clinicUserRoles := []*models.UserRole{
					// user's member role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
					// user's superadmin role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
					// user's member role at test clinic 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, swag.String(authCommon.DomainTypeLocation), nil).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), nil, swag.String(authCommon.DomainTypeClinic), nil).Return(clinicUserRoles, nil).Times(1),
					s.EXPECT().GetClinic(testClinic1.ID).Return(testClinic1, nil).Times(1),
					s.EXPECT().GetClinic(testClinic2.ID).Return(testClinic2, nil).Times(1),
				}
			},
			noErrors,
			[]string{testLocation1.ID, testLocation2.ID},
		},
		{
			"Succesfully fetched all user's locations with role filtering",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// user's superadming role at test location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}
				clinicUserRoles := []*models.UserRole{
					// user's superadmin role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeLocation), nil).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeClinic), nil).Return(clinicUserRoles, nil).Times(1),
					s.EXPECT().GetClinic(testClinic1.ID).Return(testClinic1, nil).Times(1),
				}
			},
			noErrors,
			[]string{testLocation1.ID},
		},
		{
			"Error on finding location's user roles",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeLocation), nil).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
		{
			"Error on finding clinic's user roles",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeLocation), nil).Return([]*models.UserRole{}, nil).Times(1),
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeClinic), nil).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
		{
			"Error on fetchning clinic",
			testUser1.ID,
			&authCommon.SuperadminRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// user's superadming role at test location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}
				clinicUserRoles := []*models.UserRole{
					// user's superadmin role at test clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeLocation), nil).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().FindUserRoles(swag.String(testUser1.ID), swag.String(authCommon.SuperadminRole.ID), swag.String(authCommon.DomainTypeClinic), nil).Return(clinicUserRoles, nil).Times(1),
					s.EXPECT().GetClinic(testClinic1.ID).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.UserLocationIDs(context.TODO(), test.userID, test.roleID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})

	}
}

func TestRoleUserIDs(t *testing.T) {
	testCases := []struct {
		description   string
		roleID        string
		domainType    *string
		domainID      *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all role's users",
			authCommon.MemberRole.ID,
			nil,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// testUser1's member role at organization 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization1.ID),
					},
					// testUser1's member role at organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
					// testUser1's member role at clinic 1
					&models.UserRole{
						UserID:     swag.String(testUser2.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), nil, nil).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID, testUser2.ID},
		},
		{
			"Succesfully fetched all role's users with domain filtering",
			authCommon.MemberRole.ID,
			&authCommon.DomainTypeOrganization,
			&testOrganization2.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// testUser1's member role at organization 2
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization2.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization2.ID)).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID},
		},
		{
			"Error on finding role's user roles",
			authCommon.MemberRole.ID,
			&authCommon.DomainTypeOrganization,
			&testOrganization1.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization1.ID)).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.RoleUserIDs(context.TODO(), test.roleID, test.domainType, test.domainID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func TestLocationUserIDs(t *testing.T) {
	testCases := []struct {
		description   string
		locationID    string
		roleID        *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all location's users",
			testLocation1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// testUser1's superadmin role at location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}

				locationClinics := []*models.Clinic{testClinic1}

				clinicUserRoles := []*models.UserRole{
					// testUser1's member role at clinic1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
					// testUser2's superadmin role at clinic1
					&models.UserRole{
						UserID:     swag.String(testUser2.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation1.ID)).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().GetLocationClinics(testLocation1.ID).Return(locationClinics, nil).Times(1),
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeClinic), swag.String(testClinic1.ID)).Return(clinicUserRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID, testUser2.ID},
		},
		{
			"Succesfully fetched all location's users with role filtering",
			testLocation1.ID,
			&authCommon.MemberRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{}

				locationClinics := []*models.Clinic{testClinic1}

				clinicUserRoles := []*models.UserRole{
					// testUser1's member role at clinic1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), swag.String(authCommon.DomainTypeLocation), swag.String(testLocation1.ID)).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().GetLocationClinics(testLocation1.ID).Return(locationClinics, nil).Times(1),
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), swag.String(authCommon.DomainTypeClinic), swag.String(testClinic1.ID)).Return(clinicUserRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID},
		},
		{
			"Error on finding locations's user roles",
			testLocation1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation1.ID)).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
		{
			"Error on fetching location's clinics",
			testLocation1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// testUser1's superadmin role at location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation1.ID)).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().GetLocationClinics(testLocation1.ID).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
		{
			"Error on fetching clinic's user roles",
			testLocation1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				locationUserRoles := []*models.UserRole{
					// testUser1's superadmin role at location 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.SuperadminRole.ID),
						DomainType: swag.String(authCommon.DomainTypeLocation),
						DomainID:   swag.String(testLocation1.ID),
					},
				}

				locationClinics := []*models.Clinic{testClinic1}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeLocation), swag.String(testLocation1.ID)).Return(locationUserRoles, nil).Times(1),
					s.EXPECT().GetLocationClinics(testLocation1.ID).Return(locationClinics, nil).Times(1),
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeClinic), swag.String(testClinic1.ID)).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.LocationUserIDs(context.TODO(), test.locationID, test.roleID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func TestDomainUserIDs(t *testing.T) {
	testCases := []struct {
		description   string
		domainType    *string
		domainID      *string
		roleID        *string
		mockCalls     func(storage *mock.MockStorage) []*gomock.Call
		errorExpected bool
		expected      []string
	}{
		{
			"Succesfully fetched all domain's users",
			&authCommon.DomainTypeOrganization,
			&testOrganization1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// testUser1's member role at organization 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization1.ID),
					},
					// testUser2's member role at organization 1
					&models.UserRole{
						UserID:     swag.String(testUser2.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeOrganization),
						DomainID:   swag.String(testOrganization1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization1.ID)).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID, testUser2.ID},
		},
		{
			"Succesfully fetched all domain's users with role filtering",
			&authCommon.DomainTypeClinic,
			&testClinic1.ID,
			&authCommon.MemberRole.ID,
			func(s *mock.MockStorage) []*gomock.Call {
				userRoles := []*models.UserRole{
					// testUser1's member role at organization 1
					&models.UserRole{
						UserID:     swag.String(testUser1.ID),
						RoleID:     swag.String(authCommon.MemberRole.ID),
						DomainType: swag.String(authCommon.DomainTypeClinic),
						DomainID:   swag.String(testClinic1.ID),
					},
				}

				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, swag.String(authCommon.MemberRole.ID), swag.String(authCommon.DomainTypeClinic), swag.String(testClinic1.ID)).Return(userRoles, nil).Times(1),
				}
			},
			noErrors,
			[]string{testUser1.ID},
		},
		{
			"Error on finding role's users",
			&authCommon.DomainTypeOrganization,
			&testOrganization1.ID,
			nil,
			func(s *mock.MockStorage) []*gomock.Call {
				return []*gomock.Call{
					s.EXPECT().FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeOrganization), swag.String(testOrganization1.ID)).Return(nil, fmt.Errorf("error")).Times(1),
				}
			},
			withErrors,
			nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.description, func(t *testing.T) {
			svc, storage, cleanup := getTestService(t)
			defer cleanup()

			test.mockCalls(storage)

			out, err := svc.DomainUserIDs(context.TODO(), test.domainType, test.domainID, test.roleID)

			// check expected results
			if !reflect.DeepEqual(out, test.expected) {
				fmt.Println("Expected")
				printJson(test.expected)
				fmt.Println("Got")
				printJson(out)
				t.Errorf("Expected %v, got %v", test.expected, out)
			}

			// assert error
			if test.errorExpected && err == nil {
				t.Error("Expected error, got nil")
			} else if !test.errorExpected && err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}
		})
	}
}

func getTestService(t *testing.T) (Service, *mock.MockStorage, func()) {
	// setup storage
	storageCtrl := gomock.NewController(t)
	storage := mock.NewMockStorage(storageCtrl)

	svc := New(storage, zerolog.New(os.Stdout))

	cleanup := func() {
		storageCtrl.Finish()
	}

	return svc, storage, cleanup
}

func printJson(item interface{}) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(item)
}
