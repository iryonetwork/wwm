package waitlist

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog"
)

type testStorage struct {
	*Storage
}

func newTestStorage(key []byte) *testStorage {
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
	db, err := New(path, key, zerolog.New(ioutil.Discard))
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

func TestAddList(t *testing.T) {
	s := newTestStorage(nil)
	defer s.Close()

	list, err := s.AddList("room 1")
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if list.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}

	lists, err := s.Lists()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if lists[0].ID != list.ID {
		t.Fatalf("Expected list ID to be '%s', got '%s'", list.ID, lists[0].ID)
	}
	if *(lists[0].Name) != *list.Name {
		t.Fatalf("Expected list name to be '%s', got '%s'", *list.Name, *(lists[0].Name))
	}

	var db bolt.DB
	cs := &Storage{&db}
	list, err = cs.AddList("room 2")
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
	if list != nil {
		t.Fatalf("Expected list to be nil; got '%v'", *list)
	}
}

func TestUpdateList(t *testing.T) {
	s := newTestStorage(nil)
	defer s.Close()

	list, _ := s.AddList("room 1")
	list.Name = swag.String("room 2")

	updatedList, err := s.UpdateList(list)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if *updatedList.Name != "room 2" {
		t.Fatalf("Expected updated list name to be 'room 2'; got '%s'", *updatedList.Name)
	}

	lists, err := s.Lists()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if *(lists[0]).Name != "room 2" {
		t.Fatalf("Expected updated list name to be 'room 2'; got '%s'", *(lists[0]).Name)
	}

	var db bolt.DB
	cs := &Storage{&db}
	l, err := cs.UpdateList(list)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
	if l != nil {
		t.Fatalf("Expected list to be nil; got '%v'", *list)
	}
}
