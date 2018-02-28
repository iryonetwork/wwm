package waitlist

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog"
)

type Storage struct {
	db *bolt.DB
}

var bucketCurrent = []byte("current")
var bucketHistory = []byte("history")
var bucketListNames = []byte("listsNames")
var keyQueue = []byte("queue")

const priorityLevels = 4

// New returns a new instance of storage
func New(path string, key []byte, logger zerolog.Logger) (*Storage, error) {
	logger.Debug().Msg("Initialize waitlist storage")
	if len(key) != 32 {
		return nil, fmt.Errorf("Encryption key must be 32 bytes long")
	}

	db, err := bolt.Open(path, 0x600, nil)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: db,
	}

	// add initial buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketCurrent)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bucketHistory)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bucketListNames)
		return err

	})
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
