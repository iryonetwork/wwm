package auth

import (
	"bytes"
	"fmt"

	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetUserRoles returns all userRoles
func (s *Storage) GetUserRoles() ([]*models.UserRole, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	userRoles := []*models.UserRole{}

	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		userRoles, err = s.getUserRolesWithTx(tx)

		return err
	})

	return userRoles, err
}

// getUserRolesWithTx gets userRoles from the database within passed bolt transaction
func (s *Storage) getUserRolesWithTx(tx *bolt.Tx) ([]*models.UserRole, error) {
	userRoles := []*models.UserRole{}

	b := tx.Bucket(bucketUserRoles)

	err := b.ForEach(func(_, data []byte) error {
		userRole := &models.UserRole{}

		err := userRole.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		userRoles = append(userRoles, userRole)
		return nil
	})

	return userRoles, err
}

// GetUserRole returns userRole by the id
func (s *Storage) GetUserRole(id string) (*models.UserRole, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var userRole *models.UserRole
	// look up the userRole
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		userRole, err = s.getUserRoleWithTx(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return userRole, nil
}

// getUserRoleWithTx gets userRole from the database within passed bolt transaction
func (s *Storage) getUserRoleWithTx(tx *bolt.Tx, id string) (*models.UserRole, error) {
	userRoleUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	userRole := &models.UserRole{}

	data := tx.Bucket(bucketUserRoles).Get(userRoleUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find userRole by id = '%s'", id)
	}

	// decode the userRole
	err = userRole.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}
	return userRole, nil
}

// GetUserRoleByContent returns userRole based on full content
func (s *Storage) GetUserRoleByContent(userID string, roleID string, domainType string, domainID string) (*models.UserRole, error) {
	userID, err := utils.NormalizeUUIDString(userID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Invalid userID")
	}
	roleID, err = utils.NormalizeUUIDString(roleID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Invalid roleID")
	}
	if domainID != authCommon.DomainIDWildcard {
		domainID, err = utils.NormalizeUUIDString(domainID)
		if err != nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Invalid domainID")
		}
	}

	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var userRole *models.UserRole
	// look up the userRole
	err = s.db.View(func(tx *bolt.Tx) error {
		var err error
		userRole, err = s.getUserRoleByContentWithTx(tx, userID, roleID, domainType, domainID)
		return err
	})

	if err != nil {
		return nil, err
	}
	return userRole, nil
}

// getUserRoleWithTx gets userRole from the database within passed bolt transaction
func (s *Storage) getUserRoleByContentWithTx(tx *bolt.Tx, userID string, roleID string, domainType string, domainID string) (*models.UserRole, error) {
	indexID := fmt.Sprintf("%s.%s.%s.%s", domainType, domainID, userID, roleID)

	userRoleIDBytes := tx.Bucket(bucketDomainUserRolesIndex).Get([]byte(indexID))
	if userRoleIDBytes == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find userRole with userID = %s, roleID = %s, domainType = %s, domainID = %s", userID, roleID, domainType, domainID)
	}

	userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
	if err != nil {
		return nil, err
	}

	return s.getUserRoleWithTx(tx, userRoleUUID.String())
}

// FindUserRole finds userRoles based on its parameters
func (s *Storage) FindUserRoles(userID *string, roleID *string, domainType *string, domainID *string) ([]*models.UserRole, error) {
	if userID != nil {
		id, err := utils.NormalizeUUIDString(*userID)
		if err != nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Invalid userID")
		}
		userID = swag.String(id)
	}
	if roleID != nil {
		id, err := utils.NormalizeUUIDString(*roleID)
		if err != nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Invalid roleID")
		}
		roleID = swag.String(id)
	}
	if domainID != nil && *domainID != authCommon.DomainIDWildcard {
		id, err := utils.NormalizeUUIDString(*domainID)
		if err != nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Invalid domainID")
		}
		domainID = swag.String(id)
	}

	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	userRoles := []*models.UserRole{}

	// look up the userRoles
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		userRoles, err = s.findUserRolesWithTx(tx, userID, roleID, domainType, domainID)

		return err
	})

	return userRoles, err
}

// findUserRolesWithTx finds userRoles from the database within passed bolt transaction
func (s *Storage) findUserRolesWithTx(tx *bolt.Tx, userID *string, roleID *string, domainType *string, domainID *string) ([]*models.UserRole, error) {
	// special cases in which it makes more sense to fetch user roles directly
	switch {
	case userID == nil && roleID == nil && domainType == nil && domainID == nil:
		return s.GetUserRoles()
	case userID != nil && roleID != nil && domainType != nil && *domainType == authCommon.DomainTypeGlobal:
		userRole, err := s.getUserRoleByContentWithTx(tx, *userID, *roleID, *domainType, authCommon.DomainIDWildcard)
		if err != nil {
			return nil, err
		}
		return []*models.UserRole{userRole}, nil
	case userID != nil && roleID != nil && domainType != nil && domainID != nil:
		userRole, err := s.getUserRoleByContentWithTx(tx, *userID, *roleID, *domainType, *domainID)
		if err != nil {
			return nil, err
		}
		return []*models.UserRole{userRole}, nil
	}

	// get IDs of matching user roles
	userRoleIDs, err := s.findUserRoleIDsWithTx(tx, userID, roleID, domainType, domainID)
	if err != nil {
		return nil, err
	}

	// get actual user roles
	userRoles := []*models.UserRole{}
	for _, userRoleID := range userRoleIDs {
		userRole, err := s.getUserRoleWithTx(tx, userRoleID)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// findUserRolesWithTx finds userRoles from the database within passed bolt transaction
func (s *Storage) findUserRoleIDsWithTx(tx *bolt.Tx, userID *string, roleID *string, domainType *string, domainID *string) ([]string, error) {
	userRoleIDs := []string{}

	switch {
	case userID == nil && roleID == nil && domainType == nil && domainID == nil:
		userRoles, err := s.getUserRolesWithTx(tx)
		if err != nil {
			return userRoleIDs, err
		}
		for _, userRole := range userRoles {
			userRoleIDs = append(userRoleIDs, userRole.ID)
		}
	case userID != nil:
		switch {
		case roleID != nil:
			// prefix for search in user ID index
			var userIDIndexPrefix []byte
			// prefix for search in role ID index
			var roleIDIndexPrefix []byte

			switch {
			case domainType != nil && *domainType == authCommon.DomainTypeGlobal:
				userRole, err := s.getUserRoleByContentWithTx(tx, *userID, *roleID, authCommon.DomainTypeGlobal, authCommon.DomainIDWildcard)
				if err != nil {
					return userRoleIDs, err
				}
				return []string{userRole.ID}, nil
			case domainType != nil && domainID != nil:
				userRole, err := s.getUserRoleByContentWithTx(tx, *userID, *roleID, *domainType, *domainID)
				if err != nil {
					return userRoleIDs, err
				}
				return []string{userRole.ID}, nil
			case domainType != nil && domainID == nil:
				// set prefixes
				userIDIndexPrefix = []byte(fmt.Sprintf("%s.%s", *userID, *domainType))
				roleIDIndexPrefix = []byte(fmt.Sprintf("%s.%s", *roleID, *domainType))
			case domainType == nil && domainID == nil:
				// set prefixes
				userIDIndexPrefix = []byte(*userID)
				roleIDIndexPrefix = []byte(*roleID)
			}
			// search based on set prefixes
			userIDIndexIDs := []string{}
			c := tx.Bucket(bucketUserIDUserRolesIndex).Cursor()
			for index, userRoleIDBytes := c.Seek(userIDIndexPrefix); index != nil && bytes.HasPrefix(index, userIDIndexPrefix); index, userRoleIDBytes = c.Next() {
				userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
				if err != nil {
					return nil, err
				}
				userIDIndexIDs = append(userIDIndexIDs, userRoleUUID.String())
			}
			roleIDIndexIDs := []string{}
			c = tx.Bucket(bucketRoleIDUserRolesIndex).Cursor()
			for index, userRoleIDBytes := c.Seek(roleIDIndexPrefix); index != nil && bytes.HasPrefix(index, roleIDIndexPrefix); index, userRoleIDBytes = c.Next() {
				userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
				if err != nil {
					return nil, err
				}
				roleIDIndexIDs = append(roleIDIndexIDs, userRoleUUID.String())
			}
			// get intersection of the two slices
			userRoleIDs = utils.IntersectionSlice(userIDIndexIDs, roleIDIndexIDs)
			return userRoleIDs, nil
		case roleID == nil:
			// prefix for search in user ID index
			var prefix []byte
			switch {
			case (domainType != nil && domainID == nil) || (domainType != nil && *domainType == authCommon.DomainTypeGlobal):
				prefix = []byte(fmt.Sprintf("%s.%s", *userID, *domainType))
			case domainType != nil && domainID != nil:
				prefix = []byte(fmt.Sprintf("%s.%s.%s", *userID, *domainType, *domainID))
			case domainType == nil:
				prefix = []byte(*userID)
			}
			c := tx.Bucket(bucketUserIDUserRolesIndex).Cursor()
			for index, userRoleIDBytes := c.Seek(prefix); index != nil && bytes.HasPrefix(index, prefix); index, userRoleIDBytes = c.Next() {
				userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
				if err != nil {
					return nil, err
				}
				userRoleIDs = append(userRoleIDs, userRoleUUID.String())
			}
		}
	case userID == nil && roleID != nil:
		// prefix for search in role ID index
		var prefix []byte
		switch {
		case (domainType != nil && domainID == nil) || (domainType != nil && *domainType == authCommon.DomainTypeGlobal):
			prefix = []byte(fmt.Sprintf("%s.%s", *roleID, *domainType))
		case domainType != nil && domainID != nil:
			prefix = []byte(fmt.Sprintf("%s.%s.%s", *roleID, *domainType, *domainID))
		case domainType == nil:
			prefix = []byte(*roleID)
		}
		c := tx.Bucket(bucketRoleIDUserRolesIndex).Cursor()
		for index, userRoleIDBytes := c.Seek(prefix); index != nil && bytes.HasPrefix(index, prefix); index, userRoleIDBytes = c.Next() {
			userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
			if err != nil {
				return nil, err
			}
			userRoleIDs = append(userRoleIDs, userRoleUUID.String())
		}
	case userID == nil && roleID == nil && domainType != nil:
		// prefix for search in domain index
		var prefix []byte
		switch {
		case domainID == nil || *domainType == authCommon.DomainTypeGlobal:
			prefix = []byte(*domainType)
		case domainID != nil:
			prefix = []byte(fmt.Sprintf("%s.%s", *domainType, *domainID))
		}
		c := tx.Bucket(bucketDomainUserRolesIndex).Cursor()
		for index, userRoleIDBytes := c.Seek(prefix); index != nil && bytes.HasPrefix(index, prefix); index, userRoleIDBytes = c.Next() {
			userRoleUUID, err := uuid.FromBytes(userRoleIDBytes)
			if err != nil {
				return nil, err
			}
			userRoleIDs = append(userRoleIDs, userRoleUUID.String())
		}
	}

	return userRoleIDs, nil
}

// AddUserRole generates new UUID, adds userRole to the database and updates related entities
func (s *Storage) AddUserRole(userRole *models.UserRole) (*models.UserRole, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// normalize UUIDs
	userID, err := utils.NormalizeUUIDString(*userRole.UserID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Invalid userID")
	}
	*userRole.UserID = userID
	roleID, err := utils.NormalizeUUIDString(*userRole.RoleID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Invalid roleID")
	}
	*userRole.RoleID = roleID
	if *userRole.DomainID != authCommon.DomainIDWildcard {
		domainID, err := utils.NormalizeUUIDString(*userRole.DomainID)
		if err != nil {
			return nil, utils.NewError(utils.ErrBadRequest, "Invalid domainID")
		}
		*userRole.DomainID = domainID
	}

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	userRole.ID = id.String()

	return s.addUserRole(userRole)
}

// AddGlobalSuperadminUserRole creates global superadmin role for user with provided user ID
func (s *Storage) AddGlobalSuperadminUserRole(userID string) (*models.UserRole, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// normalize UUIDs
	userID, err := utils.NormalizeUUIDString(userID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Invalid userID")
	}

	userRole := &models.UserRole{
		RoleID:     &authCommon.SuperadminRole.ID,
		DomainType: &authCommon.DomainTypeGlobal,
		DomainID:   &authCommon.DomainIDWildcard,
	}

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	userRole.ID = id.String()

	return s.addUserRole(userRole)
}

func (s *Storage) addUserRole(userRole *models.UserRole) (*models.UserRole, error) {
	var addedUserRole *models.UserRole
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// sanitize domain ID for domain type global domain ID should always be wildcard
		if *userRole.DomainType == authCommon.DomainTypeGlobal {
			userRole.DomainID = &authCommon.DomainIDWildcard
		}

		// check if user exists
		_, err = s.getUserWithTx(tx, *userRole.UserID)
		if err != nil {
			return utils.NewError(
				utils.ErrBadRequest,
				"User with userId = %s does not exist",
				*userRole.UserID,
			)
		}

		// check if role exists
		_, err = s.getRoleWithTx(tx, *userRole.RoleID)
		if err != nil {
			return utils.NewError(
				utils.ErrBadRequest,
				"Role with roleId = %s does not exist",
				*userRole.RoleID,
			)
		}

		// check if domain exists
		switch *userRole.DomainType {
		case authCommon.DomainTypeClinic:
			if *userRole.DomainID != authCommon.DomainIDWildcard {
				_, err = s.getClinicWithTx(tx, *userRole.DomainID)
			}
		case authCommon.DomainTypeOrganization:
			if *userRole.DomainID != authCommon.DomainIDWildcard {
				_, err = s.getOrganizationWithTx(tx, *userRole.DomainID)
			}
		case authCommon.DomainTypeLocation:
			if *userRole.DomainID != authCommon.DomainIDWildcard {
				_, err = s.getLocationWithTx(tx, *userRole.DomainID)
			}
		case authCommon.DomainTypeUser:
			if *userRole.DomainID != authCommon.DomainIDWildcard {
				_, err = s.getUserWithTx(tx, *userRole.DomainID)
			}
		case authCommon.DomainTypeGlobal:
			// do nothing
		case authCommon.DomainTypeCloud:
			// do nothing
		default:
			return utils.NewError(
				utils.ErrBadRequest,
				"Invalid domainType: %s",
				*userRole.DomainType,
			)
		}
		if err != nil {
			return utils.NewError(
				utils.ErrBadRequest,
				"Domain with domainType = %s, domainId = %s does not exist",
				*userRole.DomainType,
				*userRole.DomainID,
			)
		}

		// check if role with same content does exitst
		_, err = s.getUserRoleByContentWithTx(tx, *userRole.UserID, *userRole.RoleID, *userRole.DomainType, *userRole.DomainID)
		if err == nil {
			return utils.NewError(
				utils.ErrBadRequest,
				"UserRole with parameters (userId = %s, roleId = %s, domainType = %s, domainID = %s) already exists",
				*userRole.UserID,
				*userRole.RoleID,
				*userRole.DomainType,
				*userRole.DomainID,
			)
		}

		// insert userRole
		addedUserRole, err = s.insertUserRoleWithTx(tx, userRole)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return addedUserRole, err
}

// insertUserRoleWithTx inserts userRole to all database buckets within passed bolt transaction
func (s *Storage) insertUserRoleWithTx(tx *bolt.Tx, userRole *models.UserRole) (*models.UserRole, error) {
	// get IDs as UUIDs
	userRoleUUID, err := uuid.FromString(userRole.ID)
	if err != nil {
		return nil, err
	}

	data, err := userRole.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// insert user role
	err = tx.Bucket(bucketUserRoles).Put(userRoleUUID.Bytes(), data)
	if err != nil {
		return nil, err
	}

	// insert into domain index
	err = s.insertDomainIndexWithTx(tx, userRole)
	if err != nil {
		return nil, err
	}

	// insert into user ID
	err = s.insertUserIDIndexWithTx(tx, userRole)
	if err != nil {
		return nil, err
	}

	// insert into role ID index
	err = s.insertRoleIDIndexWithTx(tx, userRole)
	if err != nil {
		return nil, err
	}

	return userRole, nil
}

// insertDomainIndexWithTx inserts userRole into domain index
func (s *Storage) insertDomainIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	// get IDs as UUIDs
	userRoleUUID, err := uuid.FromString(userRole.ID)
	if err != nil {
		return err
	}

	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.DomainType, *userRole.DomainID, *userRole.UserID, *userRole.RoleID)

	return tx.Bucket(bucketDomainUserRolesIndex).Put([]byte(indexID), userRoleUUID.Bytes())
}

// insertUserIDIndexWithTx inserts userRole into user ID index
func (s *Storage) insertUserIDIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	// get IDs as UUIDs
	userRoleUUID, err := uuid.FromString(userRole.ID)
	if err != nil {
		return err
	}

	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.UserID, *userRole.DomainType, *userRole.DomainID, *userRole.RoleID)

	return tx.Bucket(bucketUserIDUserRolesIndex).Put([]byte(indexID), userRoleUUID.Bytes())
}

// insertUserIDIndexWithTx inserts userRole into user ID index
func (s *Storage) insertRoleIDIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	// get IDs as UUIDs
	userRoleUUID, err := uuid.FromString(userRole.ID)
	if err != nil {
		return err
	}

	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.RoleID, *userRole.DomainType, *userRole.DomainID, *userRole.UserID)

	return tx.Bucket(bucketRoleIDUserRolesIndex).Put([]byte(indexID), userRoleUUID.Bytes())
}

// RemoveUserRole removes userRole by id
func (s *Storage) RemoveUserRole(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		// remove user role
		return s.removeUserRoleWithTx(tx, id)
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return err
}

// removeUserRoleWithTx removes userRole from the main database bucket within passed bolt transaction
func (s *Storage) removeUserRoleWithTx(tx *bolt.Tx, id string) error {
	// get IDs as UUIDs
	userRoleUUID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	userRole, err := s.getUserRoleWithTx(tx, id)
	if err != nil {
		return err
	}

	err = s.removeUserRoleFromDomainIndexWithTx(tx, userRole)
	if err != nil {
		return err
	}
	err = s.removeUserRoleFromUserIDIndexWithTx(tx, userRole)
	if err != nil {
		return err
	}
	err = s.removeUserRoleFromRoleIDIndexWithTx(tx, userRole)
	if err != nil {
		return err
	}

	// delete from main bucket
	return tx.Bucket(bucketUserRoles).Delete(userRoleUUID.Bytes())
}

// removeUserRoleFromDomainIndexWithTx removes userRole from the domain index bucket within passed bolt transaction
func (s *Storage) removeUserRoleFromDomainIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.DomainType, *userRole.DomainID, *userRole.UserID, *userRole.RoleID)

	return tx.Bucket(bucketDomainUserRolesIndex).Delete([]byte(indexID))
}

// removeUserRoleFromUserIDIndexWithTx removes userRole from the user ID index bucket within passed bolt transaction
func (s *Storage) removeUserRoleFromUserIDIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.UserID, *userRole.DomainType, *userRole.DomainID, *userRole.RoleID)

	return tx.Bucket(bucketUserIDUserRolesIndex).Delete([]byte(indexID))
}

// removeUserRoleFromRoleIDIndexWithTx removes userRole from the role ID index bucket within passed bolt transaction
func (s *Storage) removeUserRoleFromRoleIDIndexWithTx(tx *bolt.Tx, userRole *models.UserRole) error {
	indexID := fmt.Sprintf("%s.%s.%s.%s", *userRole.RoleID, *userRole.DomainType, *userRole.DomainID, *userRole.UserID)

	return tx.Bucket(bucketRoleIDUserRolesIndex).Delete([]byte(indexID))
}

// removeUserRolesByDomainWithTx removes all userRoles from the database by user ID within passed bolt transaction
func (s *Storage) removeUserRolesByDomainWithTx(tx *bolt.Tx, domainType string, domainID string) error {
	// get IDs for removal
	userRoleIDs, err := s.findUserRoleIDsWithTx(tx, nil, nil, &domainType, &domainID)
	if err != nil {
		return err
	}
	for _, id := range userRoleIDs {
		err = s.removeUserRoleWithTx(tx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// removeUserRolesByUserIDWithTx removes all userRoles from the database by user ID within passed bolt transaction
func (s *Storage) removeUserRolesByUserIDWithTx(tx *bolt.Tx, userID string) error {
	id, err := utils.NormalizeUUIDString(userID)
	if err != nil {
		return err
	}

	// get IDs for removal
	userRoleIDs, err := s.findUserRoleIDsWithTx(tx, &id, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, id := range userRoleIDs {
		err = s.removeUserRoleWithTx(tx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

// removeUserRolesByRoleIDWithTx removes all userRoles from the database by user ID within passed bolt transaction
func (s *Storage) removeUserRolesByRoleIDWithTx(tx *bolt.Tx, roleID string) error {
	id, err := utils.NormalizeUUIDString(roleID)
	if err != nil {
		return err
	}

	// get IDs for removal
	userRoleIDs, err := s.findUserRoleIDsWithTx(tx, nil, &id, nil, nil)
	if err != nil {
		return err
	}
	for _, id := range userRoleIDs {
		err = s.removeUserRoleWithTx(tx, id)
		if err != nil {
			return err
		}
	}

	return nil
}
