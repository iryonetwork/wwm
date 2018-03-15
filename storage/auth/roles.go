package auth

import (
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetRoles returns all roles
func (s *Storage) GetRoles() ([]*models.Role, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	roles := []*models.Role{}

	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		roles, err = s.getRolesWithTx(tx)
		return err
	})

	return roles, err
}

// getRolesWithTx gets roles from the database within passed bolt transaction
func (s *Storage) getRolesWithTx(tx *bolt.Tx) ([]*models.Role, error) {
	roles := []*models.Role{}

	b := tx.Bucket(bucketRoles)

	err := b.ForEach(func(_, data []byte) error {
		role := &models.Role{}

		err := role.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		roles = append(roles, role)
		return nil
	})

	return roles, err
}

// GetRole returns role by the id
func (s *Storage) GetRole(id string) (*models.Role, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var role *models.Role
	// look up the role
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		role, err = s.getRoleWithTx(tx, id)
		return err
	})

	return role, err
}

// getRoleWithTx gets role from the database within passed bolt transaction
func (s *Storage) getRoleWithTx(tx *bolt.Tx, id string) (*models.Role, error) {
	roleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	role := &models.Role{}

	data := tx.Bucket(bucketRoles).Get(roleUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find role by id = '%s'", id)
	}

	// decode the role
	err = role.UnmarshalBinary(data)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}

	return role, err
}

// AddRole generates new UUID,  adds role to the database and updates related entities
func (s *Storage) AddRole(role *models.Role) (*models.Role, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	role.ID = id.String()

	return s.addRole(role)
}

func (s *Storage) addRole(role *models.Role) (*models.Role, error) {
	var addedRole *models.Role
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// insert role
		addedRole, err = s.insertRoleWithTx(tx, role)
		if err != nil {
			return err
		}

		return err
	})

	return addedRole, err
}

// UpdateRole updates the role and related entities in the database
func (s *Storage) UpdateRole(role *models.Role) (*models.Role, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedRole *models.Role
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error
		// get current role
		_, err = s.getRoleWithTx(tx, role.ID)
		if err != nil {
			return err
		}

		// insert role
		updatedRole, err = s.insertRoleWithTx(tx, role)
		if err != nil {
			return err
		}

		return err
	})

	return updatedRole, err
}

// insertRoleWithTx updates role in the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertRoleWithTx(tx *bolt.Tx, role *models.Role) (*models.Role, error) {
	// get ID as UUID
	roleUUID, err := uuid.FromString(role.ID)
	if err != nil {
		return nil, err
	}

	data, err := role.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// update role
	err = tx.Bucket(bucketRoles).Put(roleUUID.Bytes(), data)

	return role, err
}

// RemoveRole removes role by id and updates related entities
func (s *Storage) RemoveRole(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	if id == authCommon.EveryoneRole.ID || id == authCommon.SuperadminRole.ID {
		return utils.NewError(utils.ErrBadRequest, "You can't remove this protected role")
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := s.getRoleWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove userRoles
		err = s.removeUserRolesByRoleIDWithTx(tx, id)
		if err != nil {
			return err
		}

		err = s.removeRoleWithTx(tx, id)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// removeRoleWithTx removes role from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeRoleWithTx(tx *bolt.Tx, id string) error {
	roleUUID, _ := uuid.FromString(id)

	return tx.Bucket(bucketRoles).Delete(roleUUID.Bytes())
}
