package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	blake2b "github.com/minio/blake2b-simd"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/encrypted-bolt"
	"github.com/iryonetwork/wwm/log/errorChecker"
)

// GetChecksum calculates and returns checksum of the whole database
func (s *Storage) GetChecksum() ([]byte, error) {
	info := s.db.Info()
	reader, writer := io.Pipe()
	hash := blake2b.New256()

	go func() {
		err := s.db.View(func(tx *bolt.Tx) error {
			_, err := tx.WriteTo(writer)
			return err
		})
		errorChecker.LogError(writer.CloseWithError(err))
	}()

	// ignore metadata
	_, err := io.CopyN(ioutil.Discard, reader, int64(info.PageSize*2))
	errorChecker.LogError(err)
	_, err = io.Copy(hash, reader)
	errorChecker.LogError(err)

	return hash.Sum([]byte{}), nil
}

// WriteTo writes the whole database to a writer
func (s *Storage) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	return n, s.db.View(func(tx *bolt.Tx) error {
		var err error
		n, err = tx.WriteTo(writer)
		return err
	})
}

// ReplaceDB reads db from reader and replaces it if the checksum matches
func (s *Storage) ReplaceDB(src io.ReadCloser, checksum []byte) error {
	// save new db to temp file
	tmpFileName := s.db.Path() + base64.RawURLEncoding.EncodeToString(checksum)
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}
	src.Close()

	// create new instance from received db and check checksum
	testStorage, _, err := New(tmpFileName, s.encryptionKey, true, false, NewEnforcer, zerolog.New(ioutil.Discard))
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}

	testChecksum, err := testStorage.GetChecksum()
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}
	testStorage.Close()

	if !bytes.Equal(testChecksum, checksum) {
		os.Remove(tmpFileName)
		return fmt.Errorf("Checksums don't match")
	}

	readOnly := s.db.IsReadOnly()
	path := s.db.Path()

	s.dbSync.Lock()

	err = s.db.Close()
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}

	// replace old db file with new
	err = os.Rename(tmpFileName, path)
	if err != nil {
		return err
	}

	d, err := bolt.Open(s.encryptionKey, path, dbPermissions, &bolt.Options{ReadOnly: readOnly})
	if err != nil {
		return err
	}
	s.db = d
	s.dbSync.Unlock()
	errorChecker.LogError(s.enforcer.LoadPolicy())

	return nil
}
