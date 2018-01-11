package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

var (
	testGroup = &models.Group{
		Name: swag.String("testgroup"),
	}
	testGroup2 = &models.Group{
		Name: swag.String("testgroup2"),
	}
)

func TestAddGroup(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user
	storage.AddUser(testUser2)
	testGroup.Users = []string{testUser2.ID}

	// add group
	group, err := storage.AddGroup(testGroup)
	if group.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// add group with invalid user id
	testGroup.Users = []string{"wrong user id"}
	_, err = storage.AddGroup(testGroup)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestGetGroup(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and group
	storage.AddUser(testUser2)
	testGroup.Users = []string{testUser2.ID}
	storage.AddGroup(testGroup)

	// get group
	group, err := storage.GetGroup(testGroup.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testGroup, *group) {
		t.Fatalf("Expected returned group to be '%v', got '%v'", *testGroup, *group)
	}

	// get group with wrong uuid
	_, err = storage.GetGroup("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing group
	_, err = storage.GetGroup("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetGroups(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and group
	storage.AddUser(testUser2)
	testGroup.Users = []string{testUser2.ID}
	storage.AddGroup(testGroup)
	storage.AddGroup(testGroup2)

	// get groups
	groups, err := storage.GetGroups()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(groups) != 2 {
		t.Fatalf("Expected 2 groupss; got %d", len(groups))
	}

	groupsMap := map[string]*models.Group{}
	for _, group := range groups {
		groupsMap[group.ID] = group
	}

	if !reflect.DeepEqual(*testGroup, *groupsMap[testGroup.ID]) {
		t.Fatalf("Expected group one to be '%v', got '%v'", *testGroup, *groupsMap[testGroup.ID])
	}

	if !reflect.DeepEqual(*testGroup2, *groupsMap[testGroup2.ID]) {
		t.Fatalf("Expected group one to be '%v', got '%v'", *testGroup2, *groupsMap[testGroup2.ID])
	}
}

func TestUpdateGroup(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and group
	storage.AddUser(testUser2)
	testGroup.Users = []string{testUser2.ID}
	storage.AddGroup(testGroup)

	// update group
	updateGroup := &models.Group{
		ID:    testGroup.ID,
		Users: []string{},
		Name:  swag.String("newname"),
	}
	group, err := storage.UpdateGroup(updateGroup)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*group, *updateGroup) {
		t.Fatalf("Expected group one to be '%v', got '%v'", *group, *updateGroup)
	}

	// update group with invalid users
	updateGroup.Users = []string{"wrong"}
	_, err = storage.UpdateGroup(updateGroup)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveGroup(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user and group
	storage.AddUser(testUser2)
	testGroup.Users = []string{testUser2.ID}
	storage.AddGroup(testGroup)

	// remove group
	err := storage.RemoveGroup(testGroup.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// remove group again
	err = storage.RemoveGroup(testGroup.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
