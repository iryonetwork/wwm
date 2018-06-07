package waitlist

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
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

// MigrateVitalSigns migrates vital signs from old format to new format, to be removed
func (s *storage) MigrateVitalSigns() error {
	waitlists, err := s.Lists()
	if err != nil {
		return err
	}

	for _, waitlist := range waitlists {
		waitlistID, _ := utils.UUIDStringToBytes(waitlist.ID)

		// migrate current items
		items, err := s.ListItems(waitlistID)
		if err != nil {
			s.logger.Error().Err(err).Msgf("failed to fetch items from waitlist %s", waitlist.ID)
		} else {
			for _, item := range items {
				if item.VitalSigns != nil {
					vitalSigns := make(map[string]map[string]interface{})

					originalVitalSings, ok := item.VitalSigns.(map[string]interface{})
					if !ok {
						s.logger.Error().Msgf("couldn't migrate vitalSigns %v", item.VitalSigns)
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
										// if timestamp is float64, convert, treatingas unix epoch
										epoch, ok := vitalSignMap["timestamp"].(float64)
										if ok {
											vitalSignMap["timestamp"] = time.Unix(int64(epoch), 0).UTC().Format(time.RFC3339)
										}
									} else {
										vitalSigns[key] = map[string]interface{}{
											"timestamp": time.Now().UTC().Format(time.RFC3339),
											"value":     vitalSign,
										}
									}
								}
							} else {
								vitalSigns[key] = map[string]interface{}{
									"timestamp": time.Now().UTC().Format(time.RFC3339),
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

		// migrate history items
		items, err = s.ListHistoryItems(waitlistID, nil)
		if err != nil {
			s.logger.Error().Err(err).Msgf("failed to fetch items from waitlist %s", waitlist.ID)
		} else {
			for _, item := range items {
				if item.VitalSigns != nil {
					vitalSigns := make(map[string]map[string]interface{})

					originalVitalSings, ok := item.VitalSigns.(map[string]interface{})
					if !ok {
						s.logger.Error().Msgf("couldn't migrate vitalSigns %v", item.VitalSigns)
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
										// if timestamp is float64, convert, treatingas unix epoch
										epoch, ok := vitalSignMap["timestamp"].(float64)
										if ok {
											vitalSignMap["timestamp"] = time.Unix(int64(epoch), 0).UTC().Format(time.RFC3339)
										}
									} else {
										vitalSigns[key] = map[string]interface{}{
											"timestamp": time.Now().UTC().Format(time.RFC3339),
											"value":     vitalSign,
										}
									}
								}
							} else {
								vitalSigns[key] = map[string]interface{}{
									"timestamp": time.Now().UTC().Format(time.RFC3339),
									"value":     vitalSign,
								}
							}
						}
					}
					item.VitalSigns = vitalSigns
					_, err = s.UpdateHistoryItem(waitlistID, item)
					if err != nil {
						s.logger.Error().Msgf("failed to migrate waitlist item %v", item)
					}
				}
			}
		}
	}

	return nil
}
