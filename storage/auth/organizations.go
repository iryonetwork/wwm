package auth

import (
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/encrypted-bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetOrganizations returns all organizations
func (s *Storage) GetOrganizations() ([]*models.Organization, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var organizations []*models.Organization
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		organizations, err = s.getOrganizationsWithTx(tx)
		return err
	})

	return organizations, err
}

// getOrganizationsWithTx gets organizations from the database within passed bolt transaction
func (s *Storage) getOrganizationsWithTx(tx *bolt.Tx) ([]*models.Organization, error) {
	organizations := []*models.Organization{}

	b := tx.Bucket(bucketOrganizations)

	err := b.ForEach(func(_, data []byte) error {
		organization := &models.Organization{}
		err := organization.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		organizations = append(organizations, organization)
		return nil
	})

	return organizations, err
}

// GetOrganization returns organization by the id
func (s *Storage) GetOrganization(id string) (*models.Organization, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var organization *models.Organization
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		organization, err = s.getOrganizationWithTx(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return organization, nil
}

// getOrganizationWithTx gets organization from the database within passed bolt transaction
func (s *Storage) getOrganizationWithTx(tx *bolt.Tx, id string) (*models.Organization, error) {
	organizationUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}

	organization := &models.Organization{}
	// read organization by id
	data := tx.Bucket(bucketOrganizations).Get(organizationUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find organization by id = '%s'", id)
	}

	// decode the organization
	err = organization.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}
	return organization, nil
}

// GetOrganizationUsers returns list of clinics associated with this organization
func (s *Storage) GetOrganizationClinics(id string) ([]*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	clinics := []*models.Clinic{}
	err := s.db.View(func(tx *bolt.Tx) error {
		// get organization
		organization, err := s.getOrganizationWithTx(tx, id)
		if err != nil {
			return err
		}

		// get all the clinics
		for _, clinicID := range organization.Clinics {
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
				// error to fetch clinic is ignored but logged
				s.logger.Error().Err(err).Msg("failed to fetch clinic")
			} else {
				clinics = append(clinics, clinic)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return clinics, nil
}

// GetOrganizationLocationIDs returns list of IDs of locations associated with this organization
func (s *Storage) GetOrganizationLocationIDs(id string) ([]string, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// collect data in map to have unique locations
	locationIDsMap := make(map[string]bool)

	err := s.db.View(func(tx *bolt.Tx) error {
		// get organization
		organization, err := s.getOrganizationWithTx(tx, id)
		if err != nil {
			return err
		}

		// get all clinics at the location
		for _, clinicID := range organization.Clinics {
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
				// error to fetch clinic is ignored but logged
				s.logger.Error().Err(err).Msg("failed to fetch clinic")
			} else {
				locationIDsMap[*clinic.Location] = true
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// convert map to list
	locationIDs := []string{}
	for locationID := range locationIDsMap {
		locationIDs = append(locationIDs, locationID)
	}

	return locationIDs, nil
}

// AddOrganization generates new UUID and, adds organization to the database and updates related entities
func (s *Storage) AddOrganization(organization *models.Organization) (*models.Organization, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	organization.ID = id.String()

	return s.addOrganization(organization)
}

// AddOrganization adds organization to the database and updates related entities
func (s *Storage) addOrganization(organization *models.Organization) (*models.Organization, error) {
	var addedOrganization *models.Organization
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// get ID as UUID
		id, err := uuid.FromString(organization.ID)
		if err != nil {
			return err
		}

		// check of organization name is not already taken
		if tx.Bucket(bucketOrganizationNames).Get([]byte(*organization.Name)) != nil {
			return utils.NewError(utils.ErrBadRequest, "Organization with name %s already exists", *organization.Name)
		}

		// insert organization
		addedOrganization, err = s.insertOrganizationWithTx(tx, organization)
		if err != nil {
			return err
		}

		// insert organizationName
		return tx.Bucket(bucketOrganizationNames).Put([]byte(*addedOrganization.Name), id.Bytes())
	})

	if err != nil {
		return nil, err
	}

	if s.refreshRules {
		go s.loadPolicy()
	}

	return addedOrganization, nil
}

// UpdateOrganization updates the organization and related entities in the database
func (s *Storage) UpdateOrganization(organization *models.Organization) (*models.Organization, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedOrganization *models.Organization
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error
		// get current organization to check later if related entities changed and update them if needed
		oldOrganization, err := s.getOrganizationWithTx(tx, organization.ID)
		if err != nil {
			return err
		}

		// copy over clinics as they are read only
		organization.Clinics = oldOrganization.Clinics

		// insert organization
		updatedOrganization, err = s.insertOrganizationWithTx(tx, organization)
		if err != nil {
			return err
		}

		// update organization name if needed
		if *oldOrganization.Name != *updatedOrganization.Name {
			bOrganizationNames := tx.Bucket(bucketOrganizationNames)

			if bOrganizationNames.Get([]byte(*updatedOrganization.Name)) != nil {
				return utils.NewError(utils.ErrBadRequest, "Organization with name %s already exists", organization.Name)
			}

			err := bOrganizationNames.Delete([]byte(*oldOrganization.Name))
			if err != nil {
				return err
			}

			organizationUUID, err := uuid.FromString(updatedOrganization.ID)
			if err != nil {
				return err
			}
			err = bOrganizationNames.Put([]byte(*updatedOrganization.Name), organizationUUID.Bytes())
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
		go s.loadPolicy()
	}

	return updatedOrganization, nil
}

// insertOrganizationWithTx updates organization in the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertOrganizationWithTx(tx *bolt.Tx, organization *models.Organization) (*models.Organization, error) {
	// check if clinics for organization exist and remove non existing clinics form organization
	clinics := make([]string, len(organization.Clinics))
	length := 0
	for _, clinicID := range organization.Clinics {
		_, err := s.getClinicWithTx(tx, clinicID)
		if err == nil {
			clinics[length] = clinicID
			length++
		}
	}
	organization.Clinics = clinics[:length]

	// get ID as UUID
	organizationUUID, err := uuid.FromString(organization.ID)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}

	data, err := organization.MarshalBinary()
	if err != nil {
		return nil, err
	}

	// update organization
	err = tx.Bucket(bucketOrganizations).Put(organizationUUID.Bytes(), data)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

// addClinicToOrganizationWithTx adds clinic to the organization and updates the organization in the database within passed bolt transaction
func (s *Storage) addClinicToOrganizationWithTx(tx *bolt.Tx, organizationID, clinicID string) (*models.Organization, error) {
	organization, err := s.getOrganizationWithTx(tx, organizationID)
	if err != nil {
		return nil, err
	}
	organization.Clinics = append(organization.Clinics, clinicID)

	return s.insertOrganizationWithTx(tx, organization)
}

// removeClinicFromOrganizationWithTx removes clinic from the organization and updates the organization in the database within passed bolt transaction
func (s *Storage) removeClinicFromOrganizationWithTx(tx *bolt.Tx, organizationID, clinicID string) (*models.Organization, error) {
	organization, err := s.getOrganizationWithTx(tx, organizationID)
	if err != nil {
		return nil, err
	}
	organization.Clinics = utils.DiffSlice(organization.Clinics, []string{clinicID})

	return s.insertOrganizationWithTx(tx, organization)
}

// RemoveOrganization removes location by id and updates related entities
func (s *Storage) RemoveOrganization(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		// get organization
		organization, err := s.getOrganizationWithTx(tx, id)
		if err != nil {
			return err
		}

		err = s.removeOrganizationWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove clinics belonging to the organization and update related entries
		for _, clinicID := range organization.Clinics {
			// get clinic
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
				return err
			}
			// remove clinic from location
			_, err = s.removeClinicFromLocationWithTx(tx, *clinic.Location, clinic.ID)
			if err != nil {
				return err
			}

			// remove userRoles
			err = s.removeUserRolesByDomainWithTx(tx, authCommon.DomainTypeClinic, clinicID)
			if err != nil {
				return err
			}

			// remove clinic
			err = s.removeClinicWithTx(tx, clinicID)
			if err != nil {
				return err
			}
		}

		// remove userRoles
		err = s.removeUserRolesByDomainWithTx(tx, authCommon.DomainTypeOrganization, id)
		if err != nil {
			return err
		}

		// remove organization anme
		return tx.Bucket(bucketOrganizationNames).Delete([]byte(*organization.Name))
	})

	if err == nil && s.refreshRules {
		go s.loadPolicy()
	}

	return err
}

// removeOrganizationWithTx removes organization from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeOrganizationWithTx(tx *bolt.Tx, id string) error {
	organizationUUID, err := uuid.FromString(id)
	if err != nil {
		return utils.NewError(utils.ErrBadRequest, err.Error())
	}

	return tx.Bucket(bucketOrganizations).Delete(organizationUUID.Bytes())
}
