package discovery

//go:generate ../../bin/mockgen.sh storage/discovery Storage $GOFILE

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/agext/uuid"
	"github.com/go-openapi/strfmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/discovery/models"
	"github.com/iryonetwork/wwm/storage/discovery/db"
	"github.com/iryonetwork/wwm/utils"
)

type (
	// Storage describes state's public API
	Storage interface {
		// Create adds a new patient
		Create(conns models.Connections, locs models.Locations) (*models.Card, error)

		// Update write new patient data
		Update(patientID strfmt.UUID, conns models.Connections, locs models.Locations) (*models.Card, error)

		// Delete removes a patient
		Delete(patientID strfmt.UUID) error

		// Find looks up a patient by connection's value
		Find(q string) (models.Cards, error)

		// Get returns patient's card
		Get(patientID strfmt.UUID) (*models.Card, error)

		// Link creates a connection between a patient and a location
		Link(patientID, locationID strfmt.UUID) (models.Locations, error)

		// Unlink removes a connection between a patient and a location
		Unlink(patientID, locationID strfmt.UUID) error

		// CodesGet fetches matching codes
		CodesGet(category, query, parentID, locale string) (models.Codes, error)

		// CodeGet fetches code by ID
		CodeGet(category, id, locale string) (*models.Code, error)

		// Close closes the DB connection
		Close() error
	}

	storage struct {
		ctx        context.Context
		logger     zerolog.Logger
		db         db.DB
		gdb        *gorm.DB
		locationID string
	}

	patient struct {
		PatientID   string       `gorm:"primary_key"`
		Connections []connection `gorm:"foreignkey:PatientID"`
		Locations   []location   `gorm:"foreignkey:PatientID"`
	}

	connection struct {
		PatientID string `gorm:"primary_key"`
		Key       string `gorm:"primary_key"`
		Value     string
	}

	location struct {
		PatientID  string `gorm:"primary_key"`
		LocationID string `gorm:"primary_key"`
	}

	code struct {
		CategoryID string `gorm:"primary_key"`
		CodeID     string `gorm:"primary_key"`
		ParentID   string
		Titles     []codeTitle `gorm:"foreignkey:CategoryID,CodeID;association_foreignkey:CategoryID,CodeID"`
	}

	codeTitle struct {
		CategoryID string `gorm:"primary_key"`
		CodeID     string `gorm:"primary_key"`
		Locale     string
		Title      string
	}
)

// ErrNotFound indicates the item was not found
var ErrNotFound = errors.New("Item not found")

// New initializes a new instance of Storage
func New(ctx context.Context, gdb *gorm.DB, locID string, logger zerolog.Logger) (Storage, error) {
	s := &storage{
		ctx:        ctx,
		locationID: locID,
		logger:     logger.With().Str("component", "storage/discovery").Logger(),
		db:         db.New(gdb),
		gdb:        gdb,
	}

	return s, nil
}

func (s *storage) Create(conns models.Connections, locs models.Locations) (*models.Card, error) {
	card := &models.Card{
		PatientID:   getNewUUID(),
		Connections: conns,
		Locations:   locs,
	}

	tx := s.db.Begin()
	if err := tx.GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to start a transaction")
	}

	if err := processCardDiff(tx, &models.Card{}, card); err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to write diff to db")
	}

	tx.Commit()
	if err := tx.GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to commit new card to database")
	}

	return card, nil
}

func (s *storage) Update(patientID strfmt.UUID, conns models.Connections, locs models.Locations) (*models.Card, error) {
	tx := s.db.Begin()
	if err := tx.GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to start a transaction")
	}

	existingCard, err := getCard(tx, patientID)
	if err == ErrNotFound {
		tx.Rollback()
		return nil, err
	} else if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to fetch card to update")
	}

	newCard := &models.Card{
		PatientID:   patientID,
		Connections: conns,
		Locations:   locs,
	}

	if err := processCardDiff(tx, existingCard, newCard); err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to write diff to db")
	}

	tx.Commit()
	if err := tx.GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to commit new card to database")
	}

	return newCard, nil
}

func (s *storage) Delete(patientID strfmt.UUID) error {
	tx := s.db.Begin()
	if err := tx.GetError(); err != nil {
		return errors.Wrap(err, "failed to start a transaction")
	}

	p := patient{PatientID: patientID.String()}
	if err := tx.First(&p).GetError(); err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		return errors.Wrap(err, "failed to red patient from db")
	}

	if err := tx.Delete(connection{}, "patient_id = ?", patientID.String()).GetError(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete patient's connections")
	}

	if err := tx.Delete(location{}, "patient_id = ?", patientID.String()).GetError(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete patient's locations")
	}

	if err := tx.Delete(patient{}, "patient_id = ?", patientID.String()).GetError(); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to delete patient")
	}

	if err := tx.Commit().GetError(); err != nil {
		return errors.Wrap(err, "failed to commit new card to database")
	}

	return nil
}

func (s *storage) Get(patientID strfmt.UUID) (*models.Card, error) {
	return getCard(s.db, patientID)
}

var getCard = func(tx db.DB, patientID strfmt.UUID) (*models.Card, error) {
	card := &models.Card{
		PatientID:   patientID,
		Connections: models.Connections{},
		Locations:   models.Locations{},
	}

	p := patient{PatientID: patientID.String()}
	if err := tx.Preload("Connections").Preload("Locations").First(&p).GetError(); err == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read patient and its relations")
	}

	// process connections
	for _, conn := range p.Connections {
		card.Connections = append(card.Connections, &models.Connection{
			Key:   conn.Key,
			Value: conn.Value,
		})
	}

	// process locations
	for _, loc := range p.Locations {
		card.Locations = append(card.Locations, strfmt.UUID(loc.LocationID))
	}

	return card, nil
}

// Find looks up a patient by connection's value
func (s *storage) Find(q string) (models.Cards, error) {
	// search for each separate token
	sqlWhere := []string{}
	sqlAttrs := []interface{}{}
	tokens := strings.Split(q, " ")
	for _, token := range tokens {
		sqlWhere = append(sqlWhere, "value ILIKE ?")
		sqlAttrs = append(sqlAttrs, fmt.Sprintf("%%%s%%", token))
	}

	rows, err := s.db.Model(&connection{}).
		Where(strings.Join(sqlWhere, " OR "), sqlAttrs...).
		Select("DISTINCT patient_id").Rows()
	if err != nil {
		return nil, errors.Wrap(err, "failed to look up matching connections")
	}
	defer rows.Close()

	// iterate connections and fetch cards
	results := models.Cards{}
	for rows.Next() {
		var patientID string
		rows.Scan(&patientID)
		c, err := getCard(s.db, strfmt.UUID(patientID))
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch card for search results")
		}

		results = append(results, c)
	}

	return results, nil
}

func (s *storage) Link(patientID, locationID strfmt.UUID) (models.Locations, error) {
	tx := s.db.Begin()
	if err := tx.GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to start a transaction")
	}

	// load the current patient
	existingCard, err := getCard(tx, patientID)
	if err == ErrNotFound {
		tx.Rollback()
		return nil, err
	} else if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to fetch card to update")
	}

	// stop if the location already exists
	for _, loc := range existingCard.Locations {
		if loc == locationID {
			// already exists, return current set of locations
			tx.Rollback()
			return existingCard.Locations, nil
		}
	}

	// append the new location
	loc := location{
		PatientID:  patientID.String(),
		LocationID: locationID.String(),
	}
	if err := tx.Create(&loc).GetError(); err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "failed to insert location")
	}

	if err := tx.Commit().GetError(); err != nil {
		return nil, errors.Wrap(err, "failed to commit new location")
	}

	return append(existingCard.Locations, locationID), nil
}

func (s *storage) Unlink(patientID, locationID strfmt.UUID) error {
	tx := s.db.Begin()
	if err := tx.GetError(); err != nil {
		return errors.Wrap(err, "failed to start a transaction")
	}

	// load the current patient
	existingCard, err := getCard(tx, patientID)
	if err == ErrNotFound {
		tx.Rollback()
		return err
	} else if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to fetch card to update")
	}

	// check if the location exists
	for _, loc := range existingCard.Locations {
		if loc == locationID {
			// remove it
			loc := location{
				PatientID:  patientID.String(),
				LocationID: locationID.String(),
			}
			if err := tx.Delete(&loc).GetError(); err != nil {
				tx.Rollback()
				return errors.Wrap(err, "failed to delete location")
			}

			if err := tx.Commit().GetError(); err != nil {
				return errors.Wrap(err, "failed to commit new location")
			}

			return nil
		}
	}

	tx.Rollback()
	return ErrNotFound
}

func (s *storage) CodesGet(category, query, parentID, locale string) (models.Codes, error) {
	if locale == "" {
		locale = "en"
	}

	// build the query

	q := `SELECT ct.category_id, ct.code_id, ct.title, ct.locale, c.parent_id
		  FROM codes AS c INNER JOIN code_titles AS ct ON ct.category_id = c.category_id AND ct.code_id = c.code_id
		  WHERE c.category_id = ? AND ct.locale = ?`
	fmt.Println(q)
	values := []interface{}{category, locale}

	if query != "" {
		q += " AND ct.title ILIKE ?"
		values = append(values, fmt.Sprintf("%%%s%%", query))
	}

	if parentID != "" {
		q += " AND c.parent_id = ?"
		values = append(values, parentID)
	}

	// add order by title
	q += " ORDER BY ct.title ASC"

	rows, err := s.gdb.Raw(q, values...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to query codes")
	}
	defer rows.Close()

	res := models.Codes{}
	for rows.Next() {
		var categoryID, codeID, title, locale, parentID string
		rows.Scan(&categoryID, &codeID, &title, &locale, &parentID)

		res = append(res, &models.Code{
			Category: &categoryID,
			ID:       &codeID,
			ParentID: parentID,
			Locale:   locale,
			Title:    &title,
		})
	}

	return res, nil
}

func (s *storage) CodeGet(category, id, locale string) (*models.Code, error) {
	if locale == "" {
		locale = "en"
	}

	// build the query

	q := `SELECT ct.category_id, ct.code_id, ct.title, ct.locale, c.parent_id
		  FROM codes AS c INNER JOIN code_titles AS ct ON ct.category_id = c.category_id AND ct.code_id = c.code_id
		  WHERE c.category_id = ? AND ct.code_id = ? AND ct.locale = ?`
	fmt.Println(q)
	values := []interface{}{category, id, locale}

	rows, err := s.gdb.Raw(q, values...).Rows()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get code")
	}
	defer rows.Close()

	if rows.Next() {
		var categoryID, codeID, title, locale, parentID string
		rows.Scan(&categoryID, &codeID, &title, &locale, &parentID)

		return &models.Code{
			Category: &categoryID,
			ID:       &codeID,
			ParentID: parentID,
			Locale:   locale,
			Title:    &title,
		}, nil
	}

	return nil, utils.NewError(utils.ErrNotFound, "code_does_not_exist")
}

func (s *storage) Close() error {
	return s.db.Close()
}

// getNewUUID returns a new UUID
var getNewUUID = func() strfmt.UUID {
	return strfmt.UUID(uuid.NewCrypto().String())
}

// getCurrentTime returns the current time
var getCurrentTime = func() time.Time {
	return time.Now()
}
