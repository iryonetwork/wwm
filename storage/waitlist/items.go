package waitlist

import (
	"bytes"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// ListItems returns all items in a waitlist
func (s *storage) ListItems(waitlistID []byte) ([]*models.Item, error) {
	var list []*models.Item

	err := s.db.View(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}

		var i byte
		for i = 1; i <= priorityLevels; i++ {
			q, err := s.getQueue(waitlistID, i)
			if err != nil {
				return err
			}

			for _, itemID := range q {
				var item models.Item
				err = item.UnmarshalBinary(bCurrent.Get(itemID))
				if err != nil {
					return err
				}

				list = append(list, &item)
			}

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// AddItem creates a new item in a waitlist
func (s *storage) AddItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	if *item.Priority < 1 || *item.Priority > priorityLevels {
		return nil, utils.NewError(utils.ErrBadRequest, "invalid priority level")
	}

	item.PriorityQueue = *item.Priority

	err := s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}

		// check with item with the same patientID already exists
		var i byte
		for i = 1; i <= priorityLevels; i++ {
			q, err := s.getQueue(waitlistID, i)
			if err != nil {
				return err
			}

			for _, currentItemID := range q {
				var currentItem models.Item
				err = currentItem.UnmarshalBinary(bCurrent.Get(currentItemID))
				if err != nil {
					return err
				}

				if *currentItem.PatientID == *item.PatientID {
					return utils.NewError(utils.ErrConflict, "waitlist already contains item with provided patientID")
				}
			}

		}

		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		item.ID = id.String()
		item.Added = strfmt.DateTime(time.Now())
		item.Status = models.ItemStatusWaiting

		data, err := item.MarshalBinary()
		if err != nil {
			return err
		}

		err = bCurrent.Put(id.Bytes(), data)
		if err != nil {
			return err
		}

		return s.addToQueue(bCurrent, byte(item.PriorityQueue), id.Bytes())
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem updates an item in a waitlist
func (s *storage) UpdateItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	if *item.Priority < 1 || *item.Priority > priorityLevels {
		return nil, utils.NewError(utils.ErrBadRequest, "invalid priority level")
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}

		id, err := uuid.FromString(item.ID)
		if err != nil {
			return err
		}

		currentItemData := bCurrent.Get(id.Bytes())
		if currentItemData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found")
		}
		var currentItem models.Item
		errorChecker.LogError(currentItem.UnmarshalBinary(currentItemData))

		item.PriorityQueue = currentItem.PriorityQueue

		if *currentItem.Priority != *item.Priority {
			if *currentItem.Priority == currentItem.PriorityQueue {
				item.PriorityQueue = *item.Priority
				err := s.removeFromQueue(bCurrent, byte(currentItem.PriorityQueue), id.Bytes())
				if err != nil {
					s.logger.Error().Err(err).Msg("failed to remove from queue")
					return err
				}

				err = s.addToQueue(bCurrent, byte(item.PriorityQueue), id.Bytes())
				if err != nil {
					s.logger.Error().Err(err).Msg("failed to add to queue")
					return err
				}
			}
		}

		data, err := item.MarshalBinary()
		if err != nil {
			return err
		}

		return bCurrent.Put(id.Bytes(), data)
	})

	return item, err
}

// UpdatePatient updates with new patient data all the items with specified patientID
func (s *storage) UpdatePatient(patientID []byte, patient models.Patient) ([]*models.Item, error) {
	patientUUID, err := uuid.FromBytes(patientID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "invalid patientID")
	}
	var list []*models.Item

	err = s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketCurrent).ForEach(func(waitlistID, _ []byte) error {
			bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
			if bCurrent == nil {
				return utils.NewError(utils.ErrServerError, "waitlist not found")
			}

			var i byte
			for i = 1; i <= priorityLevels; i++ {
				q, err := s.getQueue(waitlistID, i)
				if err != nil {
					return err
				}
				for _, itemID := range q {
					var item models.Item
					err = item.UnmarshalBinary(bCurrent.Get(itemID))
					if err != nil {
						return err
					}
					if *item.PatientID == patientUUID.String() {
						item.Patient = patient
						data, err := item.MarshalBinary()
						if err != nil {
							return err
						}
						err = bCurrent.Put(itemID, data)
						if err != nil {
							return err
						}

						list = append(list, &item)
					}
				}
			}
			return err
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// MoveItemToTop moves item to the top of the list diregarding priority
func (s *storage) MoveItemToTop(waitlistID, itemID []byte) (*models.Item, error) {
	var item models.Item

	err := s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}
		itemData := bCurrent.Get(itemID)
		if itemData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found")
		}

		errorChecker.LogError(item.UnmarshalBinary(itemData))
		if item.PriorityQueue == 0 {
			item.PriorityQueue = *item.Priority
		}

		err := s.removeFromQueue(bCurrent, byte(item.PriorityQueue), itemID)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to remove from queue")
			return err
		}

		item.PriorityQueue = 1
		err = s.addToQueueOnTop(bCurrent, byte(item.PriorityQueue), itemID)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to add to queue on top")
			return err
		}

		data, err := item.MarshalBinary()
		if err != nil {
			return err
		}

		return bCurrent.Put(itemID, data)
	})

	return &item, err
}

// DeleteItem removes an item from a waitlist and moves it to history
func (s *storage) DeleteItem(waitlistID, itemID []byte, reason string) error {
	if !(reason == models.ItemStatusCanceled || reason == models.ItemStatusFinished) {
		return utils.NewError(utils.ErrBadRequest, "delete reason must be '%s' or '%s'", models.ItemStatusCanceled, models.ItemStatusFinished)
	}

	errorChecker.LogError(s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}

		itemData := bCurrent.Get(itemID)
		if itemData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found")
		}
		var item models.Item
		err := item.UnmarshalBinary(itemData)
		if err != nil {
			return err
		}

		bHistory, err := tx.Bucket(bucketHistory).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}

		item.Status = reason
		err = s.moveToHistory(bHistory, itemID, &item)
		if err != nil {
			return err
		}

		return s.removeFromQueue(bCurrent, byte(item.PriorityQueue), itemID)
	}))

	return nil
}

// ListHistoryItems returns all items in waitlist's history
func (s *storage) ListHistoryItems(waitlistID []byte, reason *string) ([]*models.Item, error) {
	var list []*models.Item

	err := s.db.View(func(tx *bolt.Tx) error {
		bHistory := tx.Bucket(bucketHistory).Bucket(waitlistID)
		if bHistory == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist history not found")
		}

		err := bHistory.ForEach(func(k, v []byte) error {
			var item models.Item
			err := item.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			if reason == nil || item.Status == *reason {
				list = append(list, &item)
			}

			return nil
		})

		return err
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}

// ReopenHistoryItem puts item from history back to waitlist
func (s *storage) ReopenHistoryItem(waitlistID, itemID, newWaitlistID []byte) (*models.Item, error) {
	var item models.Item

	err := s.db.Update(func(tx *bolt.Tx) error {
		bHistory := tx.Bucket(bucketHistory).Bucket(waitlistID)
		if bHistory == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist history not found")
		}

		itemData := bHistory.Get(itemID)
		if itemData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found")
		}

		err := item.UnmarshalBinary(itemData)
		if err != nil {
			return utils.NewError(utils.ErrServerError, "failed to unmarshal the item")
		}
		item.Status = models.ItemStatusWaiting
		item.Finished = strfmt.NewDateTime()
		item.PriorityQueue = *item.Priority

		var bCurrent *bolt.Bucket
		if newWaitlistID != nil {
			bCurrent = tx.Bucket(bucketCurrent).Bucket(newWaitlistID)
			if bCurrent == nil {
				return utils.NewError(utils.ErrNotFound, "waitlist not found")
			}
		} else {
			bCurrent, err = tx.Bucket(bucketCurrent).CreateBucketIfNotExists(waitlistID)
			if err != nil || bCurrent == nil {
				return utils.NewError(utils.ErrServerError, "failed to access waitlist")
			}
		}

		newData, err := item.MarshalBinary()
		if err != nil {
			return utils.NewError(utils.ErrServerError, "failed to marshal the item")
		}

		err = bCurrent.Put(itemID, newData)
		if err != nil {
			return utils.NewError(utils.ErrServerError, "failed to put item in the waitlist")
		}

		//remove from history
		err = bHistory.Delete(itemID)
		if err != nil {
			return utils.NewError(utils.ErrServerError, "failed to remove item from history")
		}

		return s.addToQueue(bCurrent, byte(item.PriorityQueue), itemID)
	})

	return &item, err
}

// UpdateItem updates an item in a waitlist history
func (s *storage) UpdateHistoryItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	if *item.Priority < 1 || *item.Priority > priorityLevels {
		return nil, utils.NewError(utils.ErrBadRequest, "invalid priority level")
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return utils.NewError(utils.ErrNotFound, "waitlist not found")
		}

		id, err := uuid.FromString(item.ID)
		if err != nil {
			return err
		}

		currentItemData := bCurrent.Get(id.Bytes())
		if currentItemData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found")
		}

		bHistory, err := tx.Bucket(bucketHistory).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}

		currentItemHistoryData := bHistory.Get(id.Bytes())
		if currentItemHistoryData == nil {
			return utils.NewError(utils.ErrNotFound, "item not found in waitlist history")
		}

		var currentItem models.Item
		errorChecker.LogError(currentItem.UnmarshalBinary(currentItemData))

		data, err := item.MarshalBinary()
		if err != nil {
			return err
		}

		err = bHistory.Put(id.Bytes(), data)
		if err != nil {
			return err
		}
		return bCurrent.Put(id.Bytes(), data)
	})

	return item, err
}

func (s *storage) getQueue(waitlistID []byte, priority byte) ([][]byte, error) {
	var q [][]byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return nil
		}

		qKey := append(keyQueue, priority)

		data := bCurrent.Get(qKey)
		for i := 0; i < len(data); i += 16 {
			q = append(q, data[i:i+16])
		}
		return nil
	})

	return q, err
}

func (s *storage) addToQueue(bCurrent *bolt.Bucket, priority byte, id []byte) error {
	qKey := append(keyQueue, priority)

	currentQueue := bCurrent.Get(qKey)
	if currentQueue == nil {
		currentQueue = id
	} else {
		currentQueue = append(currentQueue, id...)
	}

	return errors.Wrap(bCurrent.Put(qKey, currentQueue), "addToQueue put failed")
}

func (s *storage) addToQueueOnTop(bCurrent *bolt.Bucket, priority byte, id []byte) error {
	qKey := append(keyQueue, priority)

	currentQueue := bCurrent.Get(qKey)
	if currentQueue == nil {
		currentQueue = id
	} else {
		currentQueue = append(id, currentQueue...)
	}

	return errors.Wrap(bCurrent.Put(qKey, currentQueue), "addToQueueOnTop put failed")
}

func (s *storage) removeFromQueue(bCurrent *bolt.Bucket, priority byte, id []byte) error {
	qKey := append(keyQueue, priority)

	q := bCurrent.Get(qKey)
	if q == nil {
		return errors.New("queue does not exist")
	}

	for i := 0; i < len(q); i += 16 {
		if bytes.Equal(id, q[i:i+16]) {
			q1 := q[:i]
			q2 := q[i+16:]

			// the following line causes 'fatal error: fault; signal SIGBUS: bus error code=0x2...'
			// I have no idea why...
			//q = append(q1, q2...)
			// so we'll append twice, I guess
			q = append(append([]byte{}, q1...), q2...)

			return errors.Wrap(bCurrent.Put(qKey, q), "removeFromQueue put failed")
		}
	}

	return errors.New("item does not exist in queue")
}

func (s *storage) moveToHistory(bHistory *bolt.Bucket, itemID []byte, item *models.Item) error {
	item.Finished = strfmt.DateTime(time.Now())
	itemData, err := item.MarshalBinary()
	if err != nil {
		return err
	}

	// maybe use time.Now as key?
	return bHistory.Put(itemID, itemData)
}
