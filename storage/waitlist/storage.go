package waitlist

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
)

type storage struct {
	db     *bolt.DB
	logger *zerolog.Logger
}

var bucketCurrent = []byte("current")
var bucketHistory = []byte("history")
var bucketListMetadata = []byte("listsMetadata")
var keyQueue = []byte("queue")

const priorityLevels = 4

var dbPermissions os.FileMode = 0666

// New returns a new instance of storage
func New(path string, key []byte, logger zerolog.Logger) (*storage, error) {
	logger = logger.With().Str("component", "storage/waitlist").Logger()
	logger.Debug().Msg("Initialize waitlist storage")

	if len(key) != 32 {
		return nil, fmt.Errorf("Encryption key must be 32 bytes long")
	}

	db, err := bolt.Open(key, path, dbPermissions, nil)
	if err != nil {
		return nil, err
	}

	s := &storage{
		db:     db,
		logger: &logger,
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

		_, err = tx.CreateBucketIfNotExists(bucketListMetadata)
		return err

	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes the database
func (s *storage) Close() error {
	return s.db.Close()
}
