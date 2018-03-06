package waitlist

import (
	"bytes"
	"fmt"

	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
)

// Lists returns all active lists
func (s *storage) Lists() ([]*models.List, error) {
	var lists []*models.List

	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketCurrent).ForEach(func(waitlistID, _ []byte) error {
			name := tx.Bucket(bucketListNames).Get(waitlistID)
			if name == nil {
				return fmt.Errorf("List name not found")
			}

			id, err := uuid.FromBytes(waitlistID)
			if err != nil {
				return err
			}

			list := &models.List{
				ID:   id.String(),
				Name: swag.String(string(name)),
			}

			lists = append(lists, list)

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return lists, nil
}

// AddList adds new list
func (s *storage) AddList(name string) (*models.List, error) {
	list := &models.List{
		Name: &name,
	}
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	list.ID = id.String()

	err = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.Bucket(bucketCurrent).CreateBucket(id.Bytes())
		if err != nil {
			return err
		}
		return tx.Bucket(bucketListNames).Put(id.Bytes(), []byte(name))
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// UpdateList updates list metadata
func (s *storage) UpdateList(list *models.List) (*models.List, error) {
	id, err := uuid.FromString(list.ID)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketListNames).Put(id.Bytes(), []byte(*list.Name))
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteList removes list from active lists and move its items to history
func (s *storage) DeleteList(waitlistID []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bCurrent := tx.Bucket(bucketCurrent).Bucket(waitlistID)
		if bCurrent == nil {
			return nil
		}

		bHistory, err := tx.Bucket(bucketHistory).CreateBucketIfNotExists(waitlistID)
		if err != nil {
			return err
		}

		// mark all items in waitlist as canceled and move them to history
		err = bCurrent.ForEach(func(k, v []byte) error {
			// skip queue keys
			if bytes.HasPrefix(k, keyQueue) {
				return nil
			}

			var item models.Item
			err := item.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			item.Status = models.ItemStatusCanceled

			return s.moveToHistory(bHistory, k, &item)
		})
		if err != nil {
			return err
		}

		return tx.Bucket(bucketCurrent).DeleteBucket(waitlistID)
	})
}
