package auth

import (
	"os"
	"reflect"
	"sync"
	"testing"

	bolt "github.com/coreos/bbolt"
)

var testDBChecksum = []byte{0xbf, 0x6, 0x5c, 0x82, 0xd3, 0x7a, 0x57, 0xe9, 0x2c, 0x74, 0x95, 0x8e, 0x53, 0xc8, 0x5, 0xb4, 0x9b, 0x6b, 0x0, 0xaf, 0x15, 0xdf, 0x33, 0x90, 0x89, 0x4b, 0xa9, 0xd7, 0x65, 0xd2, 0xbc, 0x2d}

func TestGetChecksum(t *testing.T) {
	db, _ := bolt.Open("testdata/test.db", 0x600, &bolt.Options{ReadOnly: true})
	defer db.Close()

	storage := &Storage{
		db: db,
	}

	checksum, err := storage.GetChecksum()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(checksum, testDBChecksum) {
		t.Fatalf("Expected checksum to be '%v'; got '%v'", testDBChecksum, checksum)
	}
}

func TestReplaceDB(t *testing.T) {
	f, _ := os.Open("testdata/test.db")
	defer f.Close()

	storage := newTestStorage([]byte{0xe9, 0xf8, 0x2d, 0xf9, 0xc4, 0x14, 0xc1, 0x41, 0xdb, 0x87, 0x31, 0x1a, 0x95, 0x79, 0x5, 0xbf, 0x71, 0x12, 0x30, 0xd3, 0x2d, 0x8b, 0x59, 0x9d, 0x27, 0x13, 0xfa, 0x84, 0x55, 0x63, 0x64, 0x64})
	defer storage.Close()

	var wg sync.WaitGroup
	wg.Add(1000)

	for j := 0; j < 10; j++ {
		go func() {
			for i := 0; i < 100; i++ {
				_, err := storage.GetRoles()
				wg.Done()

				if err != nil {
					t.Fatalf("Error reading from database: %s", err)
				}
			}
		}()
	}

	// replace db with test.db
	err := storage.ReplaceDB(f, testDBChecksum)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// checksum should now be tha same as for test.db
	checksum, err := storage.GetChecksum()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(checksum, testDBChecksum) {
		t.Fatalf("Expected checksum to be '%v'; got '%v'", testDBChecksum, checksum)
	}
	wg.Wait()
}
