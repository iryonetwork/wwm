package waitlist

import (
	"fmt"

	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
)

func (s *Storage) ListItems(waitlistID []byte) ([]*models.Item, error) {
	var list []*models.Item

	err := s.db.View(func(tx *bolt.Tx) error {
		bCurrent, err := tx.Bucket(bucketCurrent).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}

		var i byte
		for i = 1; i <= priorityLevels; i++ {
			q, err := s.getQueue(waitlistID, i)
			if err != nil {
				return err
			}

			for _, itemID := range q {
				var item *models.Item
				err = item.UnmarshalBinary(bCurrent.Get(itemID))
				if err != nil {
					return err
				}

				list = append(list, item)
			}

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Storage) AddItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	if *item.Priority < 1 || *item.Priority > priorityLevels {
		return nil, fmt.Errorf("Invalid priority level") // TODO: return bad_request error
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		item.ID = id.String()

		bCurrent, err := tx.Bucket(bucketCurrent).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}

		data, err := item.MarshalBinary()
		if err != nil {
			return err
		}

		err = bCurrent.Put(id.Bytes(), data)
		if err != nil {
			return err
		}

		return s.addToQueue(waitlistID, byte(*item.Priority), id.Bytes())
	})

	return item, err
}

func (s *Storage) getQueue(waitlistID []byte, priority byte) ([][]byte, error) {
	if priority < 1 || priority > priorityLevels {
		return nil, fmt.Errorf("Invalid priority level") // TODO: return bad_request error
	}

	var q [][]byte
	err := s.db.View(func(tx *bolt.Tx) error {
		bCurrent, err := tx.Bucket(bucketCurrent).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
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

func (s *Storage) addToQueue(waitlistID []byte, priority byte, id []byte) error {
	if priority < 1 || priority > priorityLevels {
		return fmt.Errorf("Invalid priority level") // TODO: return bad_request error
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bCurrent, err := tx.Bucket(bucketCurrent).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}
		qKey := append(keyQueue, priority)

		currentQueue := bCurrent.Get(qKey)
		if currentQueue == nil {
			currentQueue = id
		} else {
			currentQueue = append(currentQueue, id...)
		}

		return bCurrent.Put(qKey, currentQueue)
	})
}
