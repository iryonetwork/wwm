package auth

import (
	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

// GetGroups returns all groups
func (s *Storage) GetGroups() ([]*models.Group, error) {
	groups := []*models.Group{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketGroups)

		return b.ForEach(func(k, v []byte) error {
			group := &models.Group{}
			err := group.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			groups = append(groups, group)
			return nil
		})
	})

	return groups, err
}

// GetGroup returns group by the id
func (s *Storage) GetGroup(id string) (*models.Group, error) {
	groupUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	group := &models.Group{}

	// look up the group
	err = s.db.View(func(tx *bolt.Tx) error {
		// read group by id
		data := tx.Bucket(bucketGroups).Get(groupUUID.Bytes())
		if data == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find group by id = '%s'", id)
		}

		// decode the group
		return group.UnmarshalBinary(data)
	})

	return group, err
}

// AddGroup adds group to the database
func (s *Storage) AddGroup(group *models.Group) (*models.Group, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// generatu uuid
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		group.ID = id.String()

		// check if users exist
		for _, userID := range group.Users {
			_, err := s.GetUser(userID)
			if err != nil {
				return err
			}
		}

		data, err := group.MarshalBinary()
		if err != nil {
			return err
		}

		// insert group
		return tx.Bucket(bucketGroups).Put(id.Bytes(), data)
	})

	return group, err
}

// UpdateGroup updates the group
func (s *Storage) UpdateGroup(group *models.Group) (*models.Group, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// get buckets
		bGroups := tx.Bucket(bucketGroups)

		// check if group exists
		groupUUID, err := uuid.FromString(group.ID)
		if err != nil {
			return utils.NewError(utils.ErrBadRequest, err.Error())
		}

		if bGroups.Get(groupUUID.Bytes()) == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find group by id = '%s'", group.ID)
		}

		// check if users for group exist
		for _, userID := range group.Users {
			_, err := s.GetUser(userID)
			if err != nil {
				return err
			}
		}

		data, err := group.MarshalBinary()
		if err != nil {
			return err
		}

		// update group
		return bGroups.Put(groupUUID.Bytes(), data)
	})

	return group, err
}

// RemoveGroup removes group by id
func (s *Storage) RemoveGroup(id string) error {
	_, err := s.GetGroup(id)
	if err != nil {
		return err
	}

	groupUUID, _ := uuid.FromString(id)

	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketGroups).Delete(groupUUID.Bytes())
	})
}
