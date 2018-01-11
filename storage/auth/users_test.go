package auth

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-openapi/swag"
	"golang.org/x/crypto/bcrypt"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

var (
	testUser = &models.User{
		Email:    swag.String("test@iryo.io"),
		Username: swag.String("testuser"),
		Password: "pass",
	}
	testUser2 = &models.User{
		Email:    swag.String("test2@iryo.io"),
		Username: swag.String("testuser2"),
		Password: "password",
	}
)

type testStorage struct {
	*Storage
}

func newTestStorage() *testStorage {
	// retrieve a temporary path
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	path := file.Name()
	file.Close()
	os.Remove(path)

	// open the database
	db, err := New(path)
	if err != nil {
		panic(err)
	}

	// return wrapped type
	return &testStorage{db}
}

func (storage *testStorage) Close() {
	defer os.Remove(storage.db.Path())
	storage.Storage.Close()
}

func TestAddUser(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

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
	storage := newTestStorage()
	defer storage.Close()

	// add user
	storage.AddUser(testUser)

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

func TestGetUsers(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add users
	storage.AddUser(testUser)
	storage.AddUser(testUser2)

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
	storage := newTestStorage()
	defer storage.Close()

	// add user
	storage.AddUser(testUser)

	password := testUser.Password
	updateUser := &models.User{
		ID:       testUser.ID,
		Username: testUser.Username,
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

	users, err := storage.GetUsers()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(users) != 1 {
		t.Fatalf("Expected length of user to be 1; got %d", len(users))
	}
}

func TestRemoveUser(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user
	storage.AddUser(testUser)

	// remove user
	err := storage.RemoveUser(testUser.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// remove user again
	err = storage.RemoveUser(testUser.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
