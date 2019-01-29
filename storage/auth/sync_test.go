package auth

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

var testDBChecksum = []byte{0xda, 0x7a, 0x70, 0xd6, 0xd4, 0x4c, 0x9b, 0x60, 0x68, 0x67, 0x14, 0x22, 0xc1, 0x51, 0x77, 0xe1, 0xa3, 0x45, 0x9a, 0xf4, 0x42, 0xe1, 0xd5, 0xb1, 0xd7, 0x96, 0xc7, 0xd3, 0x10, 0xd1, 0xa2, 0x84}

func TestReplaceDB(t *testing.T) {
	f, _ := os.Open("testdata/test.db")
	defer f.Close()

	storage, _ := newTestStorage([]byte{0xe9, 0xf8, 0x2d, 0xf9, 0xc4, 0x14, 0xc1, 0x41, 0xdb, 0x87, 0x31, 0x1a, 0x95, 0x79, 0x5, 0xbf, 0x71, 0x12, 0x30, 0xd3, 0x2d, 0x8b, 0x59, 0x9d, 0x27, 0x13, 0xfa, 0x84, 0x55, 0x63, 0x64, 0x64})
	defer storage.Close()

	var wg sync.WaitGroup
	wg.Add(1000)

	for j := 0; j < 10; j++ {
		errCh := make(chan error)
		for i := 0; i < 100; i++ {
			go func() {
				_, err := storage.GetRoles()
				wg.Done()
				errCh <- err
			}()
			err := <-errCh
			if err != nil {
				t.Fatalf("Error reading from database: %s", err)
			}
		}
		close(errCh)
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
