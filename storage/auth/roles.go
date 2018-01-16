package auth

import (
	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

// GetRoles returns all roles
func (s *Storage) GetRoles() ([]*models.Role, error) {
	roles := []*models.Role{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketRoles)

		return b.ForEach(func(k, v []byte) error {
			role := &models.Role{}
			err := role.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			roles = append(roles, role)
			return nil
		})
	})

	return roles, err
}

// GetRole returns role by the id
func (s *Storage) GetRole(id string) (*models.Role, error) {
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	role := &models.Role{}

	// look up the role
	err = s.db.View(func(tx *bolt.Tx) error {
		// read role by id
		data := tx.Bucket(bucketRoles).Get(roleUUID.Bytes())
		if data == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find role by id = '%s'", id)
		}

		// decode the role
		return role.UnmarshalBinary(data)
	})

	return role, err
}

// AddRole adds role to the database
func (s *Storage) AddRole(role *models.Role) (*models.Role, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// generatu uuid
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		role.ID = id.String()

		// check if users exist
		for _, userID := range role.Users {
			_, err := s.GetUser(userID)
			if err != nil {
				return err
			}
		}

		data, err := role.MarshalBinary()
		if err != nil {
			return err
		}

		// insert role
		return tx.Bucket(bucketRoles).Put(id.Bytes(), data)
	})

	if err != nil {
		go s.enforcer.LoadPolicy()
	}

	return role, err
}

// AddUserToRole adds user to role.
func (s *Storage) AddUserToRole(userID, roleID string) (*models.Role, error) {
	role, err := s.GetRole(roleID)
	if err != nil {
		return nil, err
	}
	role.Users = append(role.Users, userID)

	return s.UpdateRole(role)
}

// UpdateRole updates the role
func (s *Storage) UpdateRole(role *models.Role) (*models.Role, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// get buckets
		bRoles := tx.Bucket(bucketRoles)

		// check if role exists
		roleUUID, err := uuid.FromString(role.ID)
		if err != nil {
			return utils.NewError(utils.ErrBadRequest, err.Error())
		}

		if bRoles.Get(roleUUID.Bytes()) == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find role by id = '%s'", role.ID)
		}

		// check if users for role exist
		for _, userID := range role.Users {
			_, err := s.GetUser(userID)
			if err != nil {
				return err
			}
		}

		data, err := role.MarshalBinary()
		if err != nil {
			return err
		}

		// update role
		return bRoles.Put(roleUUID.Bytes(), data)
	})

	if err != nil {
		go s.enforcer.LoadPolicy()
	}

	return role, err
}

// RemoveRole removes role by id
func (s *Storage) RemoveRole(id string) error {
	_, err := s.GetRole(id)
	if err != nil {
		return err
	}

	roleUUID, _ := uuid.FromString(id)

	err = s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRoles).Delete(roleUUID.Bytes())
	})

	if err != nil {
		go s.enforcer.LoadPolicy()
	}

	return err
}
