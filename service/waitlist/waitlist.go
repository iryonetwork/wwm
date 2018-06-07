package waitlist

import (
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/models"
)

// Service describes the actions supported by the discovery service
type Service interface {
	// EnsureDefaultList ensures that default list exists
	EnsureDefaultWaitlist(id, name string) (*models.List, error)

	// GetWaitlists returns all active lists
	GetWaitlists() ([]*models.List, error)

	// CreateWaitlist creates new list
	CreateWaitlist(name string) (*models.List, error)

	// UpdateWaitlist updates list metadata
	UpdateWaitlist(list *models.List) (*models.List, error)

	// DeleteWaitlist removes list from active lists and move its items to history
	DeleteWaitlist(waitlistID []byte) error

	// GetWaitlist returns all items in a waitlist
	GetWaitlist(waitlistID []byte) ([]*models.Item, error)

	// DeleteItem removes an item from a waitlist and moves it to history
	DeleteItem(waitlistID, itemID []byte, reason string) error

	// CreateItem creates a new item in a waitlist
	CreateItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// UpdateItem updates an item in a waitlist
	UpdateItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// MoveItemToTop moves item to the top of the list diregarding priority
	MoveItemToTop(waitlistID, itemID []byte) (*models.Item, error)

	// GetWaitlistHistory returns all items in waitlist's history
	GetWaitlistHistory(waitlistID []byte, reason *string) ([]*models.Item, error)

	// ReopenHistoryItem puts item from history back to waitlist
	ReopenHistoryItem(waitlistID, itemID, newWaitlistID []byte) (*models.Item, error)
}

// Storage provides an interface for waitlist public functions
type Storage interface {
	// EnsureDefaultList ensures that default list exists
	EnsureDefaultList(id, name string) (*models.List, error)

	// Lists returns all active lists
	Lists() ([]*models.List, error)

	// AddList adds new list
	AddList(name string) (*models.List, error)

	// UpdateList updates list metadata
	UpdateList(list *models.List) (*models.List, error)

	// DeleteList removes list from active lists and move its items to history
	DeleteList(waitlistID []byte) error

	// ListItems returns all items in a waitlist
	ListItems(waitlistID []byte) ([]*models.Item, error)

	// AddItem creates a new item in a waitlist
	AddItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// UpdateItem updates an item in a waitlist
	UpdateItem(waitlistID []byte, item *models.Item) (*models.Item, error)

	// MoveItemToTop moves item to the top of the list diregarding priority
	MoveItemToTop(waitlistID, itemID []byte) (*models.Item, error)

	// DeleteItem removes an item from a waitlist and moves it to history
	DeleteItem(waitlistID, itemID []byte, reason string) error

	// ListHistoryItems returns all items in waitlist's history
	ListHistoryItems(waitlistID []byte, reason *string) ([]*models.Item, error)

	// ReopenHistoryItem puts item from history back to waitlist
	ReopenHistoryItem(waitlistID, itemID, newWaitlistID []byte) (*models.Item, error)
}

type service struct {
	storage Storage
	logger  zerolog.Logger
}

// EnsureDefaultList ensures that default list exists
func (s *service) EnsureDefaultWaitlist(id, name string) (*models.List, error) {
	return s.storage.EnsureDefaultList(id, name)
}

// GetWaitlists returns all active lists
func (s *service) GetWaitlists() ([]*models.List, error) {
	return s.storage.Lists()
}

// CreateWaitlist creates new list
func (s *service) CreateWaitlist(name string) (*models.List, error) {
	return s.storage.AddList(name)
}

// UpdateWaitlist updates list metadata
func (s *service) UpdateWaitlist(list *models.List) (*models.List, error) {
	return s.storage.UpdateList(list)
}

// DeleteWaitlist removes list from active lists and move its items to history
func (s *service) DeleteWaitlist(waitlistID []byte) error {
	return s.storage.DeleteList(waitlistID)
}

// GetWaitlist returns all items in a waitlist
func (s *service) GetWaitlist(waitlistID []byte) ([]*models.Item, error) {
	return s.storage.ListItems(waitlistID)
}

// DeleteItem removes an item from a waitlist and moves it to history
func (s *service) DeleteItem(waitlistID, itemID []byte, reason string) error {
	return s.storage.DeleteItem(waitlistID, itemID, reason)
}

// CreateItem creates a new item in a waitlist
func (s *service) CreateItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	return s.storage.AddItem(waitlistID, item)
}

// UpdateItem updates an item in a waitlist
func (s *service) UpdateItem(waitlistID []byte, item *models.Item) (*models.Item, error) {
	return s.storage.UpdateItem(waitlistID, item)
}

// MoveItemToTop moves item to the top of the list diregarding priority
func (s *service) MoveItemToTop(waitlistID, itemID []byte) (*models.Item, error) {
	return s.storage.MoveItemToTop(waitlistID, itemID)
}

// GetWaitlistHistory returns all items in waitlist's history
func (s *service) GetWaitlistHistory(waitlistID []byte, reason *string) ([]*models.Item, error) {
	return s.storage.ListHistoryItems(waitlistID, reason)
}

// ReopenHistoryItem puts item from history back to waitlist
func (s *service) ReopenHistoryItem(waitlistID, itemID, newWaitlistID []byte) (*models.Item, error) {
	return s.storage.ReopenHistoryItem(waitlistID, itemID, newWaitlistID)
}

// New returns a new instance of waitlist service
func New(storage Storage, logger zerolog.Logger) Service {
	return &service{
		storage: storage,
		logger:  logger.With().Str("component", "service/waitlist").Logger(),
	}
}
