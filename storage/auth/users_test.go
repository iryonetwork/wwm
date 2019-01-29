package auth

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/utils"
)

type testStorage struct {
	*Storage
}

// method to ensure that users used for tests are always fresh
func getTestUsers() (*models.User, *models.User) {
	testUser := &models.User{
		Email:    swag.String("test@iryo.io"),
		Username: swag.String("testuser"),
		Password: "pass",
	}

	testUser2 := &models.User{
		Email:    swag.String("test2@iryo.io"),
		Username: swag.String("testuser2"),
		Password: "password",
	}

	return testUser, testUser2
}

func newTestStorage(key []byte) (*testStorage, Enforcer) {
	// retrieve a temporary path
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	path := file.Name()
	file.Close()

	if key == nil {
		key = make([]byte, 32)
		_, err = rand.Read(key)
		if err != nil {
			panic(err)
		}
	}

	// open the database
	db, enforcer, err := New(path, key, false, false, NewEnforcer, zerolog.New(ioutil.Discard))
	if err != nil {
		panic(err)
	}

	// return wrapped type
	return &testStorage{db}, enforcer
}

func (storage *testStorage) Close() {
	defer os.Remove(storage.db.Path())
	storage.Storage.Close()
}

func TestAddUser(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	// add user
	user, err := storage.AddUser(testUser)
	if user.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte("pass"))
	if err != nil {
		t.Fatalf("Expected correct password hash; got error '%v'", err)
	}

	// add same user again
	_, err = storage.AddUser(testUser)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestGetUser(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)

	// get user
	user, err := storage.GetUser(testUser.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testUser, *user) {
		t.Fatalf("Expected returned user to be '%v', got '%v'", *testUser, *user)
	}

	// get user with wrong uuid
	_, err = storage.GetUser("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing user
	_, err = storage.GetUser("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetUserByUsername(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, _ := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)

	// get user
	user, err := storage.GetUserByUsername(*testUser.Username)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testUser, *user) {
		t.Fatalf("Expected returned user to be '%v', got '%v'", *testUser, *user)
	}

	// get non existing username
	_, err = storage.GetUserByUsername("no")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetUsers(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, testUser2 := getTestUsers()

	// add users
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)

	// get users
	users, err := storage.GetUsers()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(users) != 2 {
		t.Fatalf("Expected 2 users; got %d", len(users))
	}

	usersMap := map[string]*models.User{}
	for _, user := range users {
		usersMap[user.ID] = user
	}

	if !reflect.DeepEqual(*testUser, *usersMap[testUser.ID]) {
		t.Fatalf("Expected user one to be '%v', got '%v'", *testUser, *usersMap[testUser.ID])
	}

	if !reflect.DeepEqual(*testUser2, *usersMap[testUser2.ID]) {
		t.Fatalf("Expected user two to be '%v', got '%v'", *testUser2, *usersMap[testUser2.ID])
	}
}

func TestUpdateUser(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser1, testUser2 := getTestUsers()
	_, err := storage.AddUser(testUser1)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)

	password := testUser1.Password
	updateUser := &models.User{
		ID:       testUser1.ID,
		Username: testUser1.Username,
		Email:    swag.String("new@iryo.io"),
	}

	// update user
	user, err := storage.UpdateUser(updateUser)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if user.Password != password {
		t.Fatalf("Expected password to stay the same")
	}
	if *user.Email != *updateUser.Email {
		t.Fatalf("Expected email to stay the same")
	}

	// update user with username and password change
	updateUser.Username = swag.String("newusername")
	updateUser.Password = "newpassword"
	user, err = storage.UpdateUser(updateUser)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("newpassword"))
	if err != nil {
		t.Fatalf("Expected correct password hash; got error '%v'", err)
	}
	if *user.Email != *updateUser.Email {
		t.Fatalf("Expected email to stay the same")
	}

	userByUsername, err := storage.GetUserByUsername("newusername")
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if userByUsername.ID != user.ID {
		t.Fatalf("Expected to get user with id '%s'; got '%s'", user.ID, userByUsername.ID)
	}

	// cannot update user with username of other user
	updateUser.Username = testUser2.Username
	_, err = storage.UpdateUser(updateUser)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	users, err := storage.GetUsers()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(users) != 2 {
		t.Fatalf("Expected number of users to be 2; got %d", len(users))
	}
}

func TestRemoveUser(t *testing.T) {
	storage, _ := newTestStorage(nil)
	defer storage.Close()

	testUser, testUser2 := getTestUsers()
	_, err := storage.AddUser(testUser)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)

	// add user roles to test if they are removed with user
	testRole, _ := getTestRoles()
	_, err = storage.AddRole(testRole)
	errorChecker.FatalTesting(t, err)
	testUserRole1 := getTestUserRole(testUser.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(testUserRole1)
	errorChecker.FatalTesting(t, err)
	testUserRole2 := getTestUserRole(testUser2.ID, testRole.ID, authCommon.DomainTypeUser, testUser.ID)
	_, err = storage.AddUserRole(testUserRole2)
	errorChecker.FatalTesting(t, err)
	testUserRole3 := getTestUserRole(testUser2.ID, testRole.ID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
	_, err = storage.AddUserRole(testUserRole3)
	errorChecker.FatalTesting(t, err)

	// remove user
	err = storage.RemoveUser(testUser.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	// check if user roles were properly removed
	userRoles, _ := storage.GetUserRoles()
	if len(userRoles) != 4 {
		if err == nil {
			t.Fatalf("Expected 4 user role; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(swag.String(testUser.ID), nil, nil, nil)
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(nil, nil, swag.String(authCommon.DomainTypeUser), swag.String(testUser.ID))
	if len(userRoles) != 0 {
		if err == nil {
			t.Fatalf("Expected 0 user roles; got %d", len(userRoles))
		}
	}
	userRoles, _ = storage.FindUserRoles(swag.String(testUser2.ID), nil, nil, nil)
	if len(userRoles) != 4 {
		if err == nil {
			t.Fatalf("Expected 4 user role; got %d", len(userRoles))
		}
	}

	// remove user again
	err = storage.RemoveUser(testUser.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
