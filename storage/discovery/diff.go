package discovery

import (
	"sort"

	"github.com/go-openapi/strfmt"
	"github.com/iryonetwork/wwm/gen/discovery/models"
	"github.com/iryonetwork/wwm/storage/discovery/db"
	"github.com/pkg/errors"
)

func processCardDiff(tx db.DB, current *models.Card, new *models.Card) error {
	// are we adding a new card
	if current.PatientID == "" {
		p := patient{PatientID: new.PatientID.String()}
		if err := tx.Create(&p).GetError(); err != nil {
			return errors.Wrap(err, "failed to insert patient")
		}
	}

	// are we removing a card
	if new.PatientID == "" {
		p := patient{PatientID: current.PatientID.String()}
		if err := tx.Delete(&p).GetError(); err != nil {
			return errors.Wrap(err, "failed to delete patient")
		}
	}

	// compare connections
	currentKeys, currentMap := connectionsToMap(current.Connections)
	newKeys, newMap := connectionsToMap(new.Connections)
	for _, k := range currentKeys {
		c := currentMap[k]
		if e, ok := newMap[k]; !ok {
			// connection removed
			conn := connection{
				PatientID: current.PatientID.String(),
				Key:       c.Key,
			}
			if err := tx.Delete(&conn).GetError(); err != nil {
				return errors.Wrap(err, "failed to delete connection")
			}
		} else if ok && e.Value != c.Value {
			// connection updated
			conn := connection{
				PatientID: current.PatientID.String(),
				Key:       c.Key,
				Value:     e.Value,
			}
			if err := tx.Save(&conn).GetError(); err != nil {
				return errors.Wrap(err, "failed to update connection")
			}
		}
	}
	for _, k := range newKeys {
		c := newMap[k]
		if _, ok := currentMap[k]; !ok {
			// connection added
			conn := connection{
				PatientID: new.PatientID.String(),
				Key:       c.Key,
				Value:     c.Value,
			}
			if err := tx.Create(&conn).GetError(); err != nil {
				return errors.Wrap(err, "failed to add a new connection")
			}
		}
	}

	// compare locations
	// any locations removed
	for _, l := range locationsDiff(current.Locations, new.Locations) {
		loc := location{PatientID: current.PatientID.String(), LocationID: l}
		if err := tx.Delete(&loc).GetError(); err != nil {
			return errors.Wrap(err, "failed to delete location")
		}
	}
	// any locations added
	for _, l := range locationsDiff(new.Locations, current.Locations) {
		loc := location{PatientID: new.PatientID.String(), LocationID: l}
		if err := tx.Create(&loc).GetError(); err != nil {
			return errors.Wrap(err, "failed to add a location")
		}
	}

	return nil
}

func connectionsToMap(conns models.Connections) ([]string, map[string]*models.Connection) {
	m := make(map[string]*models.Connection, len(conns))
	k := []string{}

	for _, c := range conns {
		m[c.Key] = c
		k = append(k, c.Key)
	}

	sort.Strings(k)
	return k, m
}

func locationsDiff(a, b models.Locations) []string {
	mb := map[strfmt.UUID]struct{}{}
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var ab []string
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x.String())
		}
	}

	sort.Strings(ab)
	return ab
}
