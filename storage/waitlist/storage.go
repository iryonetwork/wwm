package waitlist

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// Storage provides an interface for waitlist public functions
type Storage interface {
	// EnsureDefaultList ensures that default list exists
	EnsureDefaultList(id, name string) (*models.List, error)

	// Lists returns all active lists
	Lists() ([]*models.List, error)

	// AddList adds new list
	AddList(name string) (*models.List, error)

	// UpdateList updates list metadata
	UpdateList(list *models.List) (*models.List, error)

	// DeleteList removes list from active lists and move its items to history
	DeleteList(waitlistID []byte) error

	// ListItems returns all items in a waitlist
	ListItems(waitlistID []byte) ([]*models.Item, error)

	// AddItem creates a new item in a waitlist
	AddItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// UpdateItem updates an item in a waitlist
	UpdateItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// MoveItemToTop moves item to the top of the list diregarding priority
	MoveItemToTop(waitlistID, itemID []byte) (*models.Item, error)

	// DeleteItem removes an item from a waitlist and moves it to history
	DeleteItem(waitlistID, itemID []byte, reason string) error

	// ListHistoryItems returns all items in waitlist's history
	ListHistoryItems(waitlistID []byte, reason *string) ([]*models.Item, error)

	// ReopenHistoryItem puts item from history back to waitlist
	ReopenHistoryItem(waitlistID, itemID, newWaitlistID []byte) (*models.Item, error)

	// Close closes the database
	Close() error

	// MigrateVitalSigns migrates vital signs from old format to new format, to be removed
	MigrateVitalSigns() error
}

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
func New(path string, key []byte, logger zerolog.Logger) (Storage, error) {
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

// MigrateVitalSigns migrates vital signs from old format to new format, to be removed
func (s *storage) MigrateVitalSigns() error {
	waitlists, err := s.Lists()
	if err != nil {
		return err
	}

	for _, waitlist := range waitlists {
		waitlistID, _ := utils.UUIDStringToBytes(waitlist.ID)
		items, err := s.ListItems(waitlistID)
		if err != nil {
			s.logger.Error().Err(err).Msgf("failed to fetch items from waitlist %s", waitlist.ID)
		} else {
			for _, item := range items {
				if item.VitalSigns != nil {
					vitalSigns := make(map[string]map[string]interface{})

					originalVitalSings, ok := item.VitalSigns.(map[string]interface{})
					if !ok {
						s.logger.Error().Msgf("couldn't migrate vitalSigns &v", item.VitalSigns)
					} else {
						for key, vitalSign := range originalVitalSings {
							_, ok := vitalSign.(string)
							if !ok {
								vitalSignMap, ok := vitalSign.(map[string]interface{})
								if !ok {
									s.logger.Error().Msgf("invaid vital sign value %v", vitalSign)
								} else {
									_, ok := vitalSignMap["timestamp"]
									if ok {
										vitalSigns[key] = vitalSignMap
									} else {
										vitalSigns[key] = map[string]interface{}{
											"timestamp": time.Now().Unix(),
											"value":     vitalSign,
										}
									}
								}
							} else {
								vitalSigns[key] = map[string]interface{}{
									"timestamp": time.Now().Unix(),
									"value":     vitalSign,
								}
							}
						}
					}
					item.VitalSigns = vitalSigns
					_, err = s.UpdateItem(waitlistID, item)
					if err != nil {
						s.logger.Error().Msgf("failed to migrate waitlist item %v", item)
					}
				}
			}
		}
	}
	return nil
}
