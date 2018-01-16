package auth

import (
	"fmt"

	bolt "github.com/coreos/bbolt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/utils"
)

// GetUsers returns all users
func (s *Storage) GetUsers() ([]*models.User, error) {
	users := []*models.User{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketUsers)

		return b.ForEach(func(k, v []byte) error {
			user := &models.User{}
			err := user.UnmarshalBinary(v)
			if err != nil {
				return err
			}

			users = append(users, user)
			return nil
		})
	})

	return users, err
}

// GetUser returns user by the id
func (s *Storage) GetUser(id string) (*models.User, error) {
	userUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}
	user := &models.User{}

	// look up the user
	err = s.db.View(func(tx *bolt.Tx) error {
		// read user by id
		data := tx.Bucket(bucketUsers).Get(userUUID.Bytes())
		if data == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find user by id = '%s'", id)
		}

		// decode the user
		return user.UnmarshalBinary(data)
	})

	return user, err
}

// AddUser adds user to the database
func (s *Storage) AddUser(user *models.User) (*models.User, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// check for existing user
		if tx.Bucket(bucketUsernames).Get([]byte(*user.Username)) != nil {
			return utils.NewError(utils.ErrBadRequest, "User with username %s already exists", *user.Username)
		}

		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		user.ID = id.String()

		// hash the password
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
		if err != nil {
			return err
		}
		user.Password = string(password)

		data, err := user.MarshalBinary()
		if err != nil {
			return err
		}

		// insert user
		err = tx.Bucket(bucketUsers).Put(id.Bytes(), data)
		if err != nil {
			return err
		}

		// insert username
		return tx.Bucket(bucketUsernames).Put([]byte(*user.Username), id.Bytes())
	})
	if err != nil {
		return nil, err
	}

	// add every user to everyone role
	_, err = s.AddUserToRole(user.ID, everyoneRole.ID)

	return user, err
}

// UpdateUser updates the user
func (s *Storage) UpdateUser(user *models.User) (*models.User, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		// get buckets
		bUsers := tx.Bucket(bucketUsers)
		bUsernames := tx.Bucket(bucketUsernames)

		// get current user to check if username changed
		userUUID, err := uuid.FromString(user.ID)
		if err != nil {
			return utils.NewError(utils.ErrBadRequest, err.Error())
		}

		userData := bUsers.Get(userUUID.Bytes())
		if userData == nil {
			return utils.NewError(utils.ErrNotFound, "Failed to find user by id = '%s'", user.ID)
		}
		currentUser := &models.User{}
		err = currentUser.UnmarshalBinary(userData)
		if err != nil {
			return err
		}

		if *currentUser.Username != *user.Username {
			err := bUsernames.Delete([]byte(*currentUser.Username))
			if err != nil {
				return err
			}
			err = bUsernames.Put([]byte(*user.Username), userUUID.Bytes())
			if err != nil {
				return err
			}
		}

		if user.Password == "" {
			user.Password = currentUser.Password
		} else {
			// hash the password
			password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
			if err != nil {
				return err
			}
			user.Password = string(password)
		}

		data, err := user.MarshalBinary()
		if err != nil {
			return err
		}

		// update user
		return bUsers.Put(userUUID.Bytes(), data)
	})

	return user, err
}

// RemoveUser removes user by id
func (s *Storage) RemoveUser(id string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	userUUID, _ := uuid.FromString(id)

	return s.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(bucketUsers).Delete(userUUID.Bytes())
		if err != nil {
			return err
		}
		return tx.Bucket(bucketUsernames).Delete([]byte(*user.Username))
	})
}

// GetUserByUsername returns user by the username
func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	// look up the user
	err := s.db.View(func(tx *bolt.Tx) error {
		// find the user in the usernames bucket
		id := tx.Bucket(bucketUsernames).Get([]byte(username))
		if id == nil {
			return fmt.Errorf("Failed to find username %s", username)
		}

		// read user by id
		data := tx.Bucket(bucketUsers).Get(id)
		if data == nil {
			return fmt.Errorf("Failed to find user by username %s (id = %s)", username, id)
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
