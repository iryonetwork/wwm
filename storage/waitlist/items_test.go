package waitlist

import (
	"bytes"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

var (
	patient1ID, _ = uuid.NewV4()
	patient2ID, _ = uuid.NewV4()
	patient3ID, _ = uuid.NewV4()
)

func initWaitlist(name string) ([]byte, *testStorage) {
	storage := newTestStorage(nil)
	list, err := storage.AddList(name)
	if err != nil {
		panic(err)
	}

	id, _ := uuid.FromString(list.ID)

	return id.Bytes(), storage
}

func TestAddItem(t *testing.T) {
	waitlistID, storage := initWaitlist("room 1")
	defer storage.Close()

	// create items to test
	item1 := &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient1ID.String()),
	}
	item2 := &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient2ID.String()),
	}
	item3 := &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient1ID.String()),
	}

	// #1 Succesfully add first item
	item1, err := storage.AddItem(waitlistID, item1)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if item1.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}

	// #2 Succesfully add second item
	item2, err = storage.AddItem(waitlistID, item2)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if item2.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}

	// #3 Fail to add item with items that has patientID already present in waitlist
	_, err = storage.AddItem(waitlistID, item3)
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrConflict {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrConflict, uErr.Code())
	}

	storage.db.View(func(tx *bolt.Tx) error {
		var q [32]byte
		copy(q[:], tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(1))))

		id1, _ := uuid.FromString(item1.ID)
		id2, _ := uuid.FromString(item2.ID)

		var expectedQ [32]byte

		copy(expectedQ[:16], id1.Bytes())
		copy(expectedQ[16:], id2.Bytes())

		if q != expectedQ {
			t.Fatalf("Expected queue to be '%v'; got '%v'", expectedQ, q)
		}

		if tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(id1.Bytes()) == nil {
			t.Fatalf("Expected database to have %s stored; got nil", item1.ID)
		}

		if tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(id2.Bytes()) == nil {
			t.Fatalf("Expected database to have %s stored; got nil", item2.ID)
		}

		return nil
	})

	_, err = storage.AddItem(waitlistID, &models.Item{Priority: swag.Int64(0)})
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestListItem(t *testing.T) {
	waitlistID, storage := initWaitlist("room 1")
	defer storage.Close()

	item1, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient1ID.String()),
	})

	list, err := storage.ListItems(waitlistID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if list[0].ID != item1.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item1.ID, list[0].ID)
	}

	// add high priority item
	item2, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient2ID.String()),
	})

	list, err = storage.ListItems(waitlistID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if list[0].ID != item2.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item2.ID, list[0].ID)
	}
	if list[1].ID != item1.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item1.ID, list[1].ID)
	}
}

func TestUpdateItem(t *testing.T) {
	waitlistID, storage := initWaitlist("room 1")
	defer storage.Close()

	item1, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient1ID.String()),
	})
	id1, _ := uuid.FromString(item1.ID)

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		if !bytes.Equal(q, id1.Bytes()) {
			t.Fatalf("Expected queue 4 to be have '%v'; got '%v'", id1.Bytes(), q)
		}

		return nil
	})

	item1.Priority = swag.Int64(1)
	updatedItem, err := storage.UpdateItem(waitlistID, item1)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if *updatedItem.Priority != 1 {
		t.Fatalf("Expected item priority to be 1, got %d", *updatedItem.Priority)
	}

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		if len(q) != 0 {
			t.Fatalf("Expected queue 4 to be empty; got '%v'", q)
		}

		q = tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(1)))
		if !bytes.Equal(q[:16], id1.Bytes()) {
			t.Fatalf("Expected queue 1 to have '%v' on top; got '%v'", id1.Bytes(), q[:16])
		}
		if !bytes.Equal(q, id1.Bytes()) {
			t.Fatalf("Expected queue 1 to have '%v'; got '%v'", id1.Bytes(), q)
		}

		return nil
	})
}

func TestUpdatePatient(t *testing.T) {
	waitlist1ID, storage := initWaitlist("waitlist 1")
	defer storage.Close()
	waitlist2, err := storage.AddList("waitlist 2")
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	waitlist2UUID, err := uuid.FromString(waitlist2.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	waitlist2ID := waitlist2UUID.Bytes()

	// add item for patient1 to waitlist1
	storage.AddItem(waitlist1ID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient1ID.String()),
		Patient: models.Patient{
			&models.PatientData{
				Key:   "status",
				Value: "added",
			},
		},
	})
	// add item for patient2 to waitlist1
	storage.AddItem(waitlist2ID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient1ID.String()),
		Patient: models.Patient{
			&models.PatientData{
				Key:   "status",
				Value: "added",
			},
		},
	})

	updatedItems, err := storage.UpdatePatient(
		patient1ID.Bytes(),
		models.Patient{
			&models.PatientData{
				Key:   "status",
				Value: "updated",
			},
		})

	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if updatedItems[0].Patient[0].Value != "updated" {
		t.Fatalf("Patient data returned on update were not updated")
	}
	if updatedItems[1].Patient[0].Value != "updated" {
		t.Fatalf("Patient data returned on update were not updated")
	}

	// check if update actually happen by fetching content of waitlists
	items, err := storage.ListItems(waitlist1ID)
	if err != nil {
		t.Fatalf("Failed to fetch waitlist 1")
	}
	if items[0].Patient[0].Value != "updated" {
		t.Fatalf("Patient data were not updated")
	}
	// check if update actually happen by fetching content of waitlists
	items, err = storage.ListItems(waitlist2ID)
	if err != nil {
		t.Fatalf("Failed to fetch waitlist 1")
	}
	if items[0].Patient[0].Value != "updated" {
		t.Fatalf("Patient data were not updated")
	}
}

func TestMoveItemToTop(t *testing.T) {
	waitlistID, storage := initWaitlist("room 1")
	defer storage.Close()

	// add priority 1 item
	_, _ = storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(1),
		PatientID:     swag.String(patient1ID.String()),
	})

	item1, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient2ID.String()),
	})
	id1, _ := uuid.FromString(item1.ID)

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		if !bytes.Equal(q, id1.Bytes()) {
			t.Fatalf("Expected queue 4 to be have '%v'; got '%v'", id1.Bytes(), q)
		}

		return nil
	})

	updatedItem, err := storage.MoveItemToTop(waitlistID, id1.Bytes())
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if updatedItem.PriorityQueue != 1 {
		t.Fatalf("Expected item priorityQueue to be 1, got %d", updatedItem.PriorityQueue)
	}

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		if len(q) != 0 {
			t.Fatalf("Expected queue 4 to be empty; got '%v'", q)
		}

		q = tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(1)))
		if !bytes.Equal(q[:16], id1.Bytes()) {
			t.Fatalf("Expected queue 1 to have '%v' on top; got '%v'", id1.Bytes(), q[:16])
		}

		return nil
	})

	id2, _ := uuid.NewV4()
	_, err = storage.UpdateItem(waitlistID, &models.Item{ID: id2.String(), Priority: swag.Int64(1)})
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestDeleteItem(t *testing.T) {
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

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		expectedQ := append(id1.Bytes(), append(id2.Bytes(), id3.Bytes()...)...)

		if !bytes.Equal(q, expectedQ) {
			t.Fatalf("Expected queue to be '%v'; got '%v'", q, expectedQ)
		}

		return nil
	})

	err := storage.DeleteItem(waitlistID, id2.Bytes(), models.ItemStatusCanceled)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		expectedQ := append(id1.Bytes(), id3.Bytes()...)

		if !bytes.Equal(q, expectedQ) {
			t.Fatalf("Expected queue to be '%v'; got '%v'", q, expectedQ)
		}

		return nil
	})

	err = storage.DeleteItem(waitlistID, id1.Bytes(), models.ItemStatusCanceled)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		expectedQ := id3.Bytes()

		if !bytes.Equal(q, expectedQ) {
			t.Fatalf("Expected queue to be '%v'; got '%v'", q, expectedQ)
		}

		return nil
	})

	err = storage.DeleteItem(waitlistID, id3.Bytes(), models.ItemStatusCanceled)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	storage.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket(bucketCurrent).Bucket(waitlistID).Get(append(keyQueue, byte(4)))
		expectedQ := []byte{}

		if !bytes.Equal(q, expectedQ) {
			t.Fatalf("Expected queue to be '%v'; got '%v'", q, expectedQ)
		}

		return nil
	})
}

func TestListHistoryItems(t *testing.T) {
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
	id1, _ := uuid.FromString(item1.ID)
	id2, _ := uuid.FromString(item2.ID)

	list, err := storage.ListHistoryItems(waitlistID, nil)
	if err == nil {
		t.Fatal("Expected error; got nil")
	}
	utilsErr, ok := err.(utils.Error)
	if !ok || utilsErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error to be `not_found`; got '%v'", err)
	}
	if len(list) != 0 {
		t.Fatalf("Expected list length to be 0, got %d", len(list))
	}

	err = storage.DeleteItem(waitlistID, id1.Bytes(), models.ItemStatusCanceled)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	err = storage.DeleteItem(waitlistID, id2.Bytes(), models.ItemStatusFinished)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	list, err = storage.ListHistoryItems(waitlistID, nil)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(list) != 2 {
		t.Fatalf("Expected list length to be 2, got %d", len(list))
	}

	list, err = storage.ListHistoryItems(waitlistID, swag.String(models.ItemStatusCanceled))
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}
	if list[0].ID != item1.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item1.ID, list[0].ID)
	}

	list, err = storage.ListHistoryItems(waitlistID, swag.String(models.ItemStatusFinished))
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}
	if list[0].ID != item2.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item2.ID, list[0].ID)
	}
}

func TestReopenHistoryItem(t *testing.T) {
	waitlist1ID, storage := initWaitlist("waitlist1")
	waitlist2, err := storage.AddList("waitlist2")
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	waitlist2UUID, err := uuid.FromString(waitlist2.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	waitlist2ID := waitlist2UUID.Bytes()

	item1, _ := storage.AddItem(waitlist1ID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient1ID.String()),
	})
	item2, _ := storage.AddItem(waitlist1ID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient2ID.String()),
	})
	id1, _ := uuid.FromString(item1.ID)
	id2, _ := uuid.FromString(item2.ID)

	err = storage.DeleteItem(waitlist1ID, id1.Bytes(), models.ItemStatusCanceled)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	err = storage.DeleteItem(waitlist1ID, id2.Bytes(), models.ItemStatusFinished)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	list, err := storage.ListItems(waitlist1ID)
	if err != nil {
		t.Fatalf("Expected nil; got nil; got '%v'", err)
	}
	if len(list) != 0 {
		t.Fatalf("Expected list length to be 0, got %d", len(list))
	}

	reopenedItem, err := storage.ReopenHistoryItem(waitlist1ID, id1.Bytes(), waitlist2ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if reopenedItem.Status != models.ItemStatusWaiting {
		t.Fatalf("Expected item priority to be `%s`, got `%s`", models.ItemStatusWaiting, reopenedItem.Status)
	}
	list, err = storage.ListItems(waitlist2ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}
	if list[0].ID != item1.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item1.ID, list[0].ID)
	}

	reopenedItem, err = storage.ReopenHistoryItem(waitlist1ID, id2.Bytes(), nil)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if reopenedItem.Status != models.ItemStatusWaiting {
		t.Fatalf("Expected item priority to be `%s`, got `%s`", models.ItemStatusWaiting, reopenedItem.Status)
	}
	list, err = storage.ListItems(waitlist1ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}
	if list[0].ID != item2.ID {
		t.Fatalf("Expected list item 0 ID to be '%s', got '%s'", item2.ID, list[0].ID)
	}
}

func TestUpdateHistoryItem(t *testing.T) {
	waitlistID, storage := initWaitlist("waitlist")

	item1, _ := storage.AddItem(waitlistID, &models.Item{
		MainComplaint: &models.Complaint{"something", "comment"},
		Priority:      swag.Int64(4),
		PatientID:     swag.String(patient1ID.String()),
	})
	id1, _ := uuid.FromString(item1.ID)

	err := storage.DeleteItem(waitlistID, id1.Bytes(), models.ItemStatusFinished)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	list, err := storage.ListItems(waitlistID)
	if err != nil {
		t.Fatalf("Expected nil; got nil; got '%v'", err)
	}
	if len(list) != 0 {
		t.Fatalf("Expected list length to be 0, got %d", len(list))
	}

	item1.Status = models.ItemStatusCanceled

	updatedItem, err := storage.UpdateHistoryItem(waitlistID, item1)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if updatedItem.Status != models.ItemStatusCanceled {
		t.Fatalf("Expected item status to be `canceled`, got %s", updatedItem.Status)
	}

	var item models.Item
	storage.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketHistory).Bucket(waitlistID).Get(id1.Bytes())
		if data == nil {
			t.Fatal("Expected item to be in history bucket, got nil")
		}

		err := item.UnmarshalBinary(data)
		if err != nil {
			t.Fatal("Failed to unmarshal binary")
		}
		return nil
	})

	if item.Status != models.ItemStatusCanceled {
		t.Fatalf("Expected item status to be `canceled`, got %s", item.Status)
	}
}

var items []*models.Item

func benchmarkListItems(i int, b *testing.B) {
	id, storage := initWaitlist("test")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < i; j++ {
		item := &models.Item{
			MainComplaint: &models.Complaint{"something", "comment"},
			Priority:      swag.Int64(int64(r.Intn(3) + 1)),
			PatientID:     swag.String(patient1ID.String()),
		}

		storage.AddItem(id, item)
	}
	b.ResetTimer()

	var res []*models.Item
	for n := 0; n < b.N; n++ {
		res, _ = storage.ListItems(id)
	}
	items = res
}

type sorted []*models.Item

func (a sorted) Len() int      { return len(a) }
func (a sorted) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sorted) Less(i, j int) bool {
	if *(a[i].Priority) < *(a[j].Priority) {
		return true
	}
	if *(a[i].Priority) > *(a[j].Priority) {
		return false
	}
	return a[i].Added.String() < a[j].Added.String()
}

func benchmarkListItemsSort(i int, b *testing.B) {
	id, storage := initWaitlist("test")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < i; j++ {
		item := &models.Item{
			MainComplaint: &models.Complaint{"something", "comment"},
			Priority:      swag.Int64(int64(r.Intn(3) + 1)),
			PatientID:     swag.String(patient1ID.String()),
		}

		storage.db.Update(func(tx *bolt.Tx) error {
			bCurrent := tx.Bucket(bucketCurrent).Bucket(id)

			id, err := uuid.NewV4()
			if err != nil {
				return err
			}
			item.ID = id.String()
			item.Added = strfmt.DateTime(time.Now())

			data, err := item.MarshalBinary()
			if err != nil {
				return err
			}

			return bCurrent.Put(id.Bytes(), data)
		})
	}

	b.ResetTimer()

	var res []*models.Item
	for n := 0; n < b.N; n++ {
		itms := []*models.Item{}
		storage.db.View(func(tx *bolt.Tx) error {
			bCurrent := tx.Bucket(bucketCurrent).Bucket(id)

			bCurrent.ForEach(func(_, v []byte) error {
				var currentItem models.Item
				currentItem.UnmarshalBinary(v)

				itms = append(itms, &currentItem)
				return nil
			})

			return nil
		})

		sort.Sort(sorted(itms))

		res = itms
	}
	items = res
}

func BenchmarkListItems10(b *testing.B)  { benchmarkListItems(10, b) }
func BenchmarkListItems20(b *testing.B)  { benchmarkListItems(20, b) }
func BenchmarkListItems50(b *testing.B)  { benchmarkListItems(50, b) }
func BenchmarkListItems100(b *testing.B) { benchmarkListItems(100, b) }

func BenchmarkListItemsSort10(b *testing.B)  { benchmarkListItemsSort(10, b) }
func BenchmarkListItemsSort20(b *testing.B)  { benchmarkListItemsSort(20, b) }
func BenchmarkListItemsSort50(b *testing.B)  { benchmarkListItemsSort(50, b) }
func BenchmarkListItemsSort100(b *testing.B) { benchmarkListItemsSort(100, b) }
