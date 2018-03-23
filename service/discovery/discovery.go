package discovery

//go:generate swagger generate client --quiet -A discovery -t ../../gen/discovery/ -f ../../docs/api/discovery.yml --principal string
//go:generate ../../bin/mockgen.sh service/discovery Service $GOFILE

import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"

	"github.com/iryonetwork/wwm/gen/discovery/client"
	"github.com/iryonetwork/wwm/gen/discovery/client/operations"
	"github.com/iryonetwork/wwm/gen/discovery/models"
	"github.com/iryonetwork/wwm/storage/discovery"
	"github.com/rs/zerolog"
)

type (
	// Service exposes external API
	Service interface {
		// Query searches for a card with a connection matching a given query
		// string
		Query(query string) (models.Cards, error)

		// QueryProxy calls Query on the cloud instance
		ProxyQuery(query, authToken string) (models.Cards, error)

		// Create creates a new connection
		Create(card *models.NewCard) (*models.Card, error)

		// Update updates an existing card
		Update(patientID strfmt.UUID, card *models.NewCard) (*models.Card, error)

		// Fetch looks up a card by patientID
		Fetch(patientID strfmt.UUID) (*models.Card, error)

		// Link creates a connection between a patient and a location
		Link(patientID, locationID strfmt.UUID) (models.Locations, error)

		// ProxyLink calls Link on cloud instance
		ProxyLink(patientID, locationID strfmt.UUID, authToken string) (models.Locations, error)

		// Unlink removes a connection between a patient and a location
		Unlink(patientID, locationID strfmt.UUID) error

		// ProxyUnlink calls Unlink on cloud instance
		ProxyUnlink(patientID, locationID strfmt.UUID, authToken string) error

		// Delete removes patient's card
		Delete(patientID strfmt.UUID) error
	}

	service struct {
		ctx     context.Context
		logger  zerolog.Logger
		storage discovery.Storage
		client  *client.Discovery
	}
)

var (
	// ErrNotFound marks that the requested resource was not found
	ErrNotFound = errors.New("resource not found")
)

// New returns a new instance of Service
func New(ctx context.Context, storage discovery.Storage, client *client.Discovery, log zerolog.Logger) Service {
	return &service{
		ctx:     ctx,
		logger:  log.With().Str("component", "service/discovery").Logger(),
		storage: storage,
		client:  client,
	}
}

func (svc *service) Query(query string) (models.Cards, error) {
	return svc.storage.Find(query)
}

func (svc *service) ProxyQuery(query, authToken string) (models.Cards, error) {
	if svc.client == nil {
		return nil, errors.New("client not available")
	}

	params := &operations.QueryParams{Query: &query}
	params.WithContext(svc.ctx).WithTimeout(5 * time.Second)
	res, err := svc.client.Operations.Query(params, newAuthWriter(authToken))

	if err != nil {
		return nil, errors.Wrap(err, "failed to proxy query call")
	}
	return res.Payload, nil
}

func (svc *service) Create(card *models.NewCard) (*models.Card, error) {
	return svc.storage.Create(card.Connections, card.Locations)
}

func (svc *service) Update(patientID strfmt.UUID, card *models.NewCard) (*models.Card, error) {
	return svc.storage.Update(patientID, card.Connections, card.Locations)
}

func (svc *service) Fetch(patientID strfmt.UUID) (*models.Card, error) {
	return svc.storage.Get(patientID)
}

func (svc *service) Delete(patientID strfmt.UUID) error {
	return svc.storage.Delete(patientID)
}

func (svc *service) Link(patientID, locationID strfmt.UUID) (models.Locations, error) {
	return svc.storage.Link(patientID, locationID)
}

func (svc *service) ProxyLink(patientID, locationID strfmt.UUID, authToken string) (models.Locations, error) {
	if svc.client == nil {
		return nil, errors.New("client not available")
	}

	params := &operations.LinkParams{PatientID: patientID, LocationID: locationID}
	params.WithContext(svc.ctx).WithTimeout(5 * time.Second)
	res, err := svc.client.Operations.Link(params, newAuthWriter(authToken))

	if _, ok := err.(*operations.LinkNotFound); ok {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to proxy link call")
	}

	return res.Payload, nil
}

func (svc *service) Unlink(patientID, locationID strfmt.UUID) error {
	return svc.storage.Unlink(patientID, locationID)
}

func (svc *service) ProxyUnlink(patientID, locationID strfmt.UUID, authToken string) error {
	if svc.client == nil {
		return errors.New("client not available")
	}

	params := &operations.UnlinkParams{
		PatientID:  patientID,
		LocationID: locationID,
	}
	params.WithContext(svc.ctx).WithTimeout(5 * time.Second)
	_, err := svc.client.Operations.Unlink(params, newAuthWriter(authToken))

	if _, ok := err.(*operations.LinkNotFound); ok {
		return ErrNotFound
	} else if err != nil {
		return errors.Wrap(err, "failed to proxy Unlink call")
	}
	return nil
}
