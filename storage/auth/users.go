package auth

import (
	"github.com/go-openapi/swag"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetUsers returns all users
func (s *Storage) GetUsers() ([]*models.User, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	users := []*models.User{}

	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		users, err = s.getUsersWithTx(tx)

		return err
	})

	return users, err
}

// getUsersWithTx gets users from the database within passed bolt transaction
func (s *Storage) getUsersWithTx(tx *bolt.Tx) ([]*models.User, error) {
	users := []*models.User{}

	b := tx.Bucket(bucketUsers)

	err := b.ForEach(func(_, data []byte) error {
		user := &models.User{}

		err := user.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		users = append(users, user)
		return nil
	})

	return users, err
}

// GetUser returns user by the id
func (s *Storage) GetUser(id string) (*models.User, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var user *models.User
	// look up the user
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		user, err = s.getUserWithTx(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// getUserWithTx gets user from the database within passed bolt transaction
func (s *Storage) getUserWithTx(tx *bolt.Tx, id string) (*models.User, error) {
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	user := &models.User{}

	data := tx.Bucket(bucketUsers).Get(userUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find user by id = '%s'", id)
	}

	// decode the user
	err = user.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUser generates new UUID, adds user to the database and updates related entities
func (s *Storage) AddUser(user *models.User) (*models.User, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	user.ID = id.String()

	return s.addUser(user)
}

func (s *Storage) addUser(user *models.User) (*models.User, error) {
	var addedUser *models.User
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// get ID as UUID
		id, err := uuid.FromString(user.ID)
		if err != nil {
			return err
		}

		// check if username is not already taken
		if tx.Bucket(bucketUsernames).Get([]byte(*user.Username)) != nil {
			return utils.NewError(utils.ErrBadRequest, "User with username %s already exists", *user.Username)
		}

		// hash the password
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
		if err != nil {
			return err
		}
		user.Password = string(password)

		// insert user
		addedUser, err = s.insertUserWithTx(tx, user)
		if err != nil {
			return err
		}

		// insert username
		err = tx.Bucket(bucketUsernames).Put([]byte(*addedUser.Username), id.Bytes())

		return nil
	})

	if err != nil {
		return nil, err
	}

	// give every new user everyone role globally
	_, err = s.AddUserRole(&models.UserRole{
		UserID:     swag.String(user.ID),
		RoleID:     swag.String(authCommon.EveryoneRole.ID),
		DomainType: swag.String(authCommon.DomainTypeGlobal),
		DomainID:   swag.String(authCommon.DomainIDWildcard),
	})
	if err != nil {
		return addedUser, err
	}

	// give every new user member role for cloud
	_, err = s.AddUserRole(&models.UserRole{
		UserID:     swag.String(user.ID),
		RoleID:     swag.String(authCommon.MemberRole.ID),
		DomainType: swag.String(authCommon.DomainTypeCloud),
		DomainID:   swag.String(authCommon.DomainIDWildcard),
	})
	if err != nil {
		return addedUser, err
	}

	// give every new user author over own user domain
	_, err = s.AddUserRole(&models.UserRole{
		UserID:     swag.String(user.ID),
		RoleID:     swag.String(authCommon.AuthorRole.ID),
		DomainType: swag.String(authCommon.DomainTypeUser),
		DomainID:   swag.String(user.ID),
	})
	if err != nil {
		return addedUser, err
	}

	return addedUser, nil
}

// UpdateUser updates the user and related entities
func (s *Storage) UpdateUser(user *models.User) (*models.User, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedUser *models.User
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// get currentUser
		oldUser, err := s.getUserWithTx(tx, user.ID)
		if err != nil {
			return err
		}

		// check if password is changing
		if user.Password == "" {
			user.Password = oldUser.Password
		} else {
			// hash the password
			password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
			if err != nil {
				return err
			}
			user.Password = string(password)
		}

		// insert updated user
		updatedUser, err = s.insertUserWithTx(tx, user)
		if err != nil {
			return err
		}

		// update username if needed
		if *oldUser.Username != *updatedUser.Username {
			bUsernames := tx.Bucket(bucketUsernames)

			// check if new username is not already taken
			if bUsernames.Get([]byte(*updatedUser.Username)) != nil {
				return utils.NewError(utils.ErrBadRequest, "User with username %s already exists", updatedUser.Username)
			}

			// delete old user name
			err := bUsernames.Delete([]byte(*oldUser.Username))
			if err != nil {
				return err
			}

			userUUID, err := uuid.FromString(updatedUser.ID)
			if err != nil {
				return err
			}
			// insert new username
			err = bUsernames.Put([]byte(*updatedUser.Username), userUUID.Bytes())
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

// insertUserWithTx updates user to the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertUserWithTx(tx *bolt.Tx, user *models.User) (*models.User, error) {
	// get ID as UUID
	userUUID, err := uuid.FromString(user.ID)
	if err != nil {
		return nil, err
	}

	data, err := user.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// update user
	err = tx.Bucket(bucketUsers).Put(userUUID.Bytes(), data)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// RemoveUser removes user by id
func (s *Storage) RemoveUser(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		// fetch user
		user, err := s.getUserWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove userRoles
		err = s.removeUserRolesByUserIDWithTx(tx, id)
		if err != nil {
			return err
		}
		// remove userRoles on this user domain
		err = s.removeUserRolesByDomainWithTx(tx, authCommon.DomainTypeUser, id)
		if err != nil {
			return err
		}

		// remove user
		err = s.removeUserWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove username
		return tx.Bucket(bucketUsernames).Delete([]byte(*user.Username))
	})

	return err
}

// removeUserWithTx removes user from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeUserWithTx(tx *bolt.Tx, id string) error {
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return tx.Bucket(bucketUsers).Delete(userUUID.Bytes())
}

// GetUserByUsername returns user by the username
func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	user := &models.User{}

	// look up the user
	err := s.db.View(func(tx *bolt.Tx) error {
		// find the user in the usernames bucket
		id := tx.Bucket(bucketUsernames).Get([]byte(username))
		if id == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find username %s", username)
		}

		// read user by id
		data := tx.Bucket(bucketUsers).Get(id)
		if data == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find user by username %s (id = %s)", username, id)
		}

		// decode the user
		err := user.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}
