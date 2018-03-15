package auth

import (
	uuid "github.com/satori/go.uuid"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/storage/encrypted_bolt"
	"github.com/iryonetwork/wwm/utils"
)

// GetLocations returns all locations
func (s *Storage) GetLocations() ([]*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	locations := []*models.Location{}

	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		locations, err = s.getLocationsWithTx(tx)
		return err
	})

	return locations, err
}

// getLocationsWithTx gets locations from the database within passed bolt transaction
func (s *Storage) getLocationsWithTx(tx *bolt.Tx) ([]*models.Location, error) {
	locations := []*models.Location{}

	b := tx.Bucket(bucketLocations)

	err := b.ForEach(func(_, data []byte) error {
		location := &models.Location{}
		err := location.UnmarshalBinary(data)
		if err != nil {
			return err
		}

		locations = append(locations, location)
		return nil
	})

	return locations, err
}

// GetLocation returns location by the id
func (s *Storage) GetLocation(id string) (*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var location *models.Location
	// look up the location
	err := s.db.View(func(tx *bolt.Tx) error {
		var err error
		location, err = s.getLocationWithTx(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	return location, nil
}

// getLocationWithTx gets location from the database within passed bolt transaction
func (s *Storage) getLocationWithTx(tx *bolt.Tx, id string) (*models.Location, error) {
	locationUUID, err := uuid.FromString(id)
	if err != nil {
		return nil, utils.NewError(utils.ErrBadRequest, err.Error())
	}

	location := &models.Location{}

	// read location by id
	data := tx.Bucket(bucketLocations).Get(locationUUID.Bytes())
	if data == nil {
		return nil, utils.NewError(utils.ErrNotFound, "Failed to find location by id = '%s'", id)
	}

	// decode the location
	err = location.UnmarshalBinary(data)

	if err != nil {
		return nil, err
	}
	return location, nil
}

// GetLocationUsers returns list of clinics associated with this location
func (s *Storage) GetLocationClinics(id string) ([]*models.Clinic, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	clinics := []*models.Clinic{}
	err := s.db.View(func(tx *bolt.Tx) error {
		// get locations
		location, err := s.getLocationWithTx(tx, id)
		if err != nil {
			return err
		}

		// get all clinics
		for _, clinicID := range location.Clinics {
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
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

// GetLocationOrganizationIDs returns list of IDs of organizations associated with this location
func (s *Storage) GetLocationOrganizationIDs(id string) ([]string, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// collect data in map to have unique organizations
	organizationIDsMap := make(map[string]bool)
	err := s.db.View(func(tx *bolt.Tx) error {
		// ger location
		location, err := s.getLocationWithTx(tx, id)
		if err != nil {
			return err
		}

		// get all clinics from location
		for _, clinicID := range location.Clinics {
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
				// error to fetch clinic is ignored but logged
				s.logger.Error().Err(err).Msg("failed to fetch clinic")
			} else {
				organizationIDsMap[*clinic.Organization] = true
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// convert map to list
	organizationIDs := []string{}
	for organizationID := range organizationIDsMap {
		organizationIDs = append(organizationIDs, organizationID)
	}

	return organizationIDs, nil
}

// AddLocation generates new UUID, adds locationto the database and updates related entities
func (s *Storage) AddLocation(location *models.Location) (*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	// generate ID
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	location.ID = id.String()

	return s.addLocation(location)
}

func (s *Storage) addLocation(location *models.Location) (*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var addedLocation *models.Location
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error

		// get ID as UUID
		id, err := uuid.FromString(location.ID)
		if err != nil {
			return err
		}

		// check if location is not already taken
		if tx.Bucket(bucketLocationNames).Get([]byte(*location.Name)) != nil {
			return utils.NewError(utils.ErrBadRequest, "Location with name %s already exists", *location.Name)
		}

		// insert location
		addedLocation, err = s.insertLocationWithTx(tx, location)
		if err != nil {
			return err
		}

		// insert locationName
		err = tx.Bucket(bucketLocationNames).Put([]byte(*addedLocation.Name), id.Bytes())
		if err != nil {
			return err
		}

		return err
	})

	if err != nil {
		return nil, err
	}

	if s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return addedLocation, nil
}

// UpdateLocation updates location and related entities in the database
func (s *Storage) UpdateLocation(location *models.Location) (*models.Location, error) {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	var updatedLocation *models.Location
	err := s.db.Update(func(tx *bolt.Tx) error {
		var err error
		// get current location to check if locationName changed
		oldLocation, err := s.getLocationWithTx(tx, location.ID)
		if err != nil {
			return err
		}

		// copy over clinics that are read only
		location.Clinics = oldLocation.Clinics

		// insert location
		updatedLocation, err = s.insertLocationWithTx(tx, location)
		if err != nil {
			return err
		}

		// update updatedLocation name if needed
		if *oldLocation.Name != *updatedLocation.Name {
			bLocationNames := tx.Bucket(bucketLocationNames)
			if bLocationNames.Get([]byte(*updatedLocation.Name)) != nil {
				return utils.NewError(utils.ErrBadRequest, "Location with name %s already exists", *updatedLocation.Name)
			}

			err := bLocationNames.Delete([]byte(*oldLocation.Name))
			if err != nil {
				return err
			}

			locationUUID, err := uuid.FromString(updatedLocation.ID)
			if err != nil {
				return err
			}
			err = bLocationNames.Put([]byte(*updatedLocation.Name), locationUUID.Bytes())
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

	return updatedLocation, nil
}

// insertLocationWithTx updates location in the database within passed bolt transaction (does not update related entities)
func (s *Storage) insertLocationWithTx(tx *bolt.Tx, location *models.Location) (*models.Location, error) {
	// check if clinics for location exist and remove non-existing clinics from the list
	clinics := make([]string, len(location.Clinics))
	length := 0
	for _, clinicID := range location.Clinics {
		_, err := s.getClinicWithTx(tx, clinicID)
		if err == nil {
			clinics[length] = clinicID
			length++
		}
	}
	location.Clinics = clinics[:length]

	// get ID as UUID
	locationUUID, err := uuid.FromString(location.ID)
	if err != nil {
		return nil, err
	}

	data, err := location.MarshalBinary()
	if err != nil {
		return nil, err
	}

	err = tx.Bucket(bucketLocations).Put(locationUUID.Bytes(), data)
	if err != nil {
		return nil, err
	}

	return location, err
}

// addClinicToLocationWithTx adds clinic to location and updates location in the database within passed bolt transaction
func (s *Storage) addClinicToLocationWithTx(tx *bolt.Tx, locationID, clinicID string) (*models.Location, error) {
	location, err := s.getLocationWithTx(tx, locationID)
	if err != nil {
		return nil, err
	}
	location.Clinics = append(location.Clinics, clinicID)

	return s.insertLocationWithTx(tx, location)
}

// removeClinicToLocationWithTx removes clinic from location and updates location in the database within passed bolt transaction
func (s *Storage) removeClinicFromLocationWithTx(tx *bolt.Tx, locationID, clinicID string) (*models.Location, error) {
	location, err := s.getLocationWithTx(tx, locationID)
	if err != nil {
		return nil, err
	}
	location.Clinics = utils.DiffSlice(location.Clinics, []string{clinicID})

	return s.insertLocationWithTx(tx, location)
}

// RemoveLocation removes location by id and updates related entities
func (s *Storage) RemoveLocation(id string) error {
	s.dbSync.RLock()
	defer s.dbSync.RUnlock()

	err := s.db.Update(func(tx *bolt.Tx) error {
		// get location
		location, err := s.getLocationWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove location
		err = s.removeLocationWithTx(tx, id)
		if err != nil {
			return err
		}

		// remove clinics belonging to the location and update related entries
		for _, clinicID := range location.Clinics {
			// get clinic
			clinic, err := s.getClinicWithTx(tx, clinicID)
			if err != nil {
				return err
			}
			// remove clinic from location
			_, err = s.removeClinicFromOrganizationWithTx(tx, *clinic.Organization, clinic.ID)
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

		// remove userRoles explicitly set for location
		err = s.removeUserRolesByDomainWithTx(tx, authCommon.DomainTypeLocation, id)
		if err != nil {
			return err
		}

		// remove location name
		return tx.Bucket(bucketLocationNames).Delete([]byte(*location.Name))
	})

	if err == nil && s.refreshRules {
		go s.enforcer.LoadPolicy()
	}

	return err
}

// removeLocationWithTx removes location from the database within passed bolt transaction (does not update related entities)
func (s *Storage) removeLocationWithTx(tx *bolt.Tx, id string) error {
	locationUUID, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	return tx.Bucket(bucketLocations).Delete(locationUUID.Bytes())
}
