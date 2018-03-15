package auth

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
	authCommon "github.com/iryonetwork/wwm/auth"
)

// GetClinics returns all clinics
func (s *Storage) GetClinics() ([]*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	clinics := []*models.Clinic{}
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		clinics, err = s.getClinicsWithTx(tx)
		return err
	})

	return clinics, err
}

// getClinicsWithTx gets clinics from the database within passed bolt transaction
func (s *Storage) getClinicsWithTx(tx *bolt.Tx) ([]*models.Clinic, error) {
	clinics := []*models.Clinic{}

	b := tx.Bucket(bucketClinics)
	err := b.ForEach(func(_, data []byte) error {
		clinic := &models.Clinic{}
		err := clinic.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		clinics = append(clinics, clinic)
		return nil
	})

	return clinics, err
}

// GetClinic returns clinic by the id
func (s *Storage) GetClinic(id string) (*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	clinic := &models.Clinic{}

	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		clinic, err = s.getClinicWithTx(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return clinic, nil
}

// getClinicWithTx gets clinic from the database within passed bolt transaction
func (s *Storage) getClinicWithTx(tx *bolt.Tx, id string) (*models.Clinic, error) {
	clinicUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}

	clinic := &models.Clinic{}

	data := tx.Bucket(bucketClinics).Get(clinicUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find clinic by id = '%s'", id)
	}

	err = clinic.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}
	return clinic, nil
}

// GetClinicOrganization returns clinic's organization
func (s *Storage) GetClinicOrganization(id string) (*models.Organization, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var organization *models.Organization
	err := s.db.View(func(tx *bolt.Tx) error {
		// get clinic
		clinic, err := s.getClinicWithTx(tx, id)
		if err != nil {
			return err
		}

		// get organization
		organization, err = s.getOrganizationWithTx(tx, *clinic.Organization)
		return err
	})

	if err != nil {
		return nil, err
	}
	return organization, nil
}

// GetClinicLocation returns clinic's location
func (s *Storage) GetClinicLocation(id string) (*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var location *models.Location
	err := s.db.View(func(tx *bolt.Tx) error {
		// get clinic
		clinic, err := s.getClinicWithTx(tx, id)
		if err != nil {
			return err
		}

		// get location
		location, err = s.getLocationWithTx(tx, *clinic.Location)
		return err
	})

	if err != nil {
		return nil, err
	}
	return location, nil
}

// AddClinic generates new UUID, adds clinic to the database and updates related entities
func (s *Storage) AddClinic(clinic *models.Clinic) (*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID for clinic
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	clinic.ID = id.String()

	return s.addClinic(clinic)
}

func (s *Storage) addClinic(clinic *models.Clinic) (*models.Clinic, error) {
	var addedClinic *models.Clinic
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// get ID as UUID
		id, err := uuid.FromString(clinic.ID)
		if err != nil {
			return err
		}

		// check of clinic name is not already taken
		if tx.Bucket(bucketClinicNames).Get([]byte(getFullClinicName(clinic))) != nil {
			return utils.NewError(utils.ErrBadRequest, "Clinic with name %s already exists in the specified location", clinic.Name)
		}

		// insert clinic
		addedClinic, err = s.insertClinicWithTx(tx, clinic)
		if err != nil {
			return err
		}

		// insert clinic name
		err = tx.Bucket(bucketClinicNames).Put([]byte(getFullClinicName(addedClinic)), id.Bytes())
		if err != nil {
			return err
		}

		// update organization of the clinic
		_, err = s.addClinicToOrganizationWithTx(tx, *addedClinic.Organization, addedClinic.ID)
		if err != nil {
			return err
		}
		// update location of the clinic
		_, err = s.addClinicToLocationWithTx(tx, *addedClinic.Location, addedClinic.ID)
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

	return addedClinic, nil
}

// UpdateClinic updates clinic and related entities in the database
func (s *Storage) UpdateClinic(clinic *models.Clinic) (*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedClinic *models.Clinic
	err := s.db.Update(func(tx *bolt.Tx) error {
		// get current clinic to check later if related entities changed and update them if needed
		oldClinic, err := s.getClinicWithTx(tx, clinic.ID)
		if err != nil {
			return err
		}

		// insert clinic
		updatedClinic, err = s.insertClinicWithTx(tx, clinic)
		if err != nil {
			return err
		}

		// update clinic name if needed
		if getFullClinicName(oldClinic) != getFullClinicName(updatedClinic) {
			bClinicNames := tx.Bucket(bucketClinicNames)

			// check if new clinic name is not taken
			if bClinicNames.Get([]byte(getFullClinicName(updatedClinic))) != nil {
				return utils.NewError(utils.ErrBadRequest, "Clinic with name %s already exists in the specified location", clinic.Name)
			}

			// delete old clinic name
			err := bClinicNames.Delete([]byte(getFullClinicName(oldClinic)))
			if err != nil {
				return err
			}

			clinicUUID, err := uuid.FromString(updatedClinic.ID)
			if err != nil {
				return err
			}
			// insert new clinic name
			err = bClinicNames.Put([]byte(getFullClinicName(updatedClinic)), clinicUUID.Bytes())
			if err != nil {
				return err
			}
		}

		// update organization(s) if needed
		if *oldClinic.Organization != *updatedClinic.Organization {
			_, err := s.removeClinicFromOrganizationWithTx(tx, *oldClinic.Organization, updatedClinic.ID)
			if err != nil {
				return err
			}

			_, err = s.addClinicToOrganizationWithTx(tx, *updatedClinic.Organization, updatedClinic.ID)
			if err != nil {
				return err
			}
		}
		// update location(s) if needed
		if *oldClinic.Location != *updatedClinic.Location {
			_, err := s.removeClinicFromLocationWithTx(tx, *oldClinic.Location, updatedClinic.ID)
			if err != nil {
				return err
			}

			_, err = s.addClinicToLocationWithTx(tx, *updatedClinic.Location, updatedClinic.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return updatedClinic, nil
}

// insertClinicWithTx updates clinic in the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertClinicWithTx(tx *bolt.Tx, clinic *models.Clinic) (*models.Clinic, error) {
	// check if clinic location exists, if not - return an error
	_, err := s.getLocationWithTx(tx, *clinic.Location)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Location with id %s does not exist", *clinic.Location)
	}

	// check if clinic organization exists, if not - return an error
	_, err = s.getOrganizationWithTx(tx, *clinic.Organization)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, "Organization with id %s does not exist", *clinic.Organization)
	}

	// get ID as UUID
	clinicUUID, err := uuid.FromString(clinic.ID)
	if err != nil {
		return nil, err
	}

	data, err := clinic.MarshalBinary()
	if err != nil {
		return nil, err
	}

	err = tx.Bucket(bucketClinics).Put(clinicUUID.Bytes(), data)
	if err != nil {
		return nil, err
	}

	return clinic, nil
}

// RemoveClinic removes clinic from the database by id and updates related entities
func (s *Storage) RemoveClinic(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		// get current clinic
		clinic, err := s.getClinicWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove clinic
		err = s.removeClinicWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove clinic from organization
		_, err = s.removeClinicFromOrganizationWithTx(tx, *clinic.Organization, clinic.ID)
		if err != nil {
			return err
		}
		// remove clinic from location
		_, err = s.removeClinicFromLocationWithTx(tx, *clinic.Location, clinic.ID)
		if err != nil {
			return err
		}
		// remove userRoles
		err = s.removeUserRolesByDomainWithTx(tx, authCommon.DomainTypeClinic, clinic.ID)
		if err != nil {
			return err
		}

		// remove clinic name
		return tx.Bucket(bucketClinicNames).Delete([]byte(getFullClinicName(clinic)))
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return err
}

// removeClinicWithTx removes clinic from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeClinicWithTx(tx *bolt.Tx, id string) error {
	clinicUUID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return tx.Bucket(bucketClinics).Delete(clinicUUID.Bytes())
}

// getFullClinicName returns clinic name prefixed with 'locationID.organizationID.'
func getFullClinicName(clinic *models.Clinic) string {
	return fmt.Sprintf("%s.%s.%s", *clinic.Location, *clinic.Organization, *clinic.Name)
}
