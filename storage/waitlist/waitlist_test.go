package waitlist

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/encrypted-bolt"
)

type testStorage struct {
	*storage
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
	storage.storage.Close()
}

func TestEnsureDefaultList(t *testing.T) {
	s := newTestStorage(nil)
	defer s.Close()

	defaultListID := "22afd921-0630-49f4-89a8-d1ad7639ee83"
	defaultListName := "default"

	list, err := s.EnsureDefaultList(defaultListID, defaultListName)
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
		t.Fatalf("Expected list ID to be '%s', got '%s'", defaultListID, lists[0].ID)
	}
	if *(lists[0].Name) != *list.Name {
		t.Fatalf("Expected list name to be '%s', got '%s'", defaultListName, *(lists[0].Name))
	}

	// trying to ensure again with the same ID, should return the list from storage and do not add another one
	list, err = s.EnsureDefaultList(defaultListID, "different name")
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if list.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}

	lists, err = s.Lists()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if lists[0].ID != list.ID {
		t.Fatalf("Expected list ID to be '%s', got '%s'", defaultListID, lists[0].ID)
	}
	if *(lists[0].Name) != *list.Name {
		t.Fatalf("Expected list name to be '%s', got '%s'", defaultListName, *(lists[0].Name))
	}

	// trying to add with invalid ID should return an error
	_, err = s.EnsureDefaultList("invalid_id", "list with invalid id")
	if err == nil {
		t.Fatalf("Expected error to error, got '%v'", err)
	}
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
	cs := &storage{db: &db}
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
	cs := &storage{db: &db}
	l, err := cs.UpdateList(list)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
	if l != nil {
		t.Fatalf("Expected list to be nil; got '%v'", *list)
	}
}

func TestDeleteList(t *testing.T) {
	waitlistID, storage := initWaitlist("room 1")
	defer storage.Close()

	item1, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient1ID.String()),
	})
	item2, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient2ID.String()),
	})
	item3, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient3ID.String()),
	})

	id1, _ := uuid.FromString(item1.ID)
	id2, _ := uuid.FromString(item2.ID)
	id3, _ := uuid.FromString(item3.ID)

	errorChecker.FatalTesting(t, storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		expectedQ := append(id1.Bytes(), append(id2.Bytes(), id3.Bytes()...)...)

		if !bytes.Equal(q, expectedQ) {
			t.Fatalf("Expected queue to be '%v'; got '%v'", q, expectedQ)
		}

		return nil
	}))

	err := storage.DeleteList(waitlistID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	errorChecker.FatalTesting(t, storage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if b != nil {
			t.Fatalf("Exepected bucket to be deleted")
		}

		for _, id := range []uuid.UUID{id1, id2, id3} {
			if tx.Bucket(bucketHistory).Bucket(waitlistID).Get(id.Bytes()) == nil {
				t.Fatalf("Exepected item %s to be in history", id.String())
			}
		}

		return nil
	}))
}
