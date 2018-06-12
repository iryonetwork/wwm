package waitlist

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/item"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/waitlist"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the discovery handlers
type Handlers interface {
	// GetWaitlists returns all active lists
	GetWaitlists() waitlist.GetHandler
	// CreateWaitlist creates new list
	CreateWaitlist() waitlist.PostHandler
	// UpdateWaitlist updates list metadata
	UpdateWaitlist() waitlist.PutListIDHandler
	// DeleteWaitlist removes list from active lists and move its items to history
	DeleteWaitlist() waitlist.DeleteListIDHandler
	// GetWaitlist returns all items in a waitlist
	GetWaitlist() item.GetListIDHandler
	// GetWaitlistHistory returns all items in waitlist's history
	GetWaitlistHistory() item.GetListIDHistoryHandler
	// DeleteItem removes an item from a waitlist and moves it to history
	DeleteItem() item.DeleteListIDItemIDHandler
	// CreateItem creates a new item in a waitlist
	CreateItem() item.PostListIDHandler
	// UpdateItem updates an item in a waitlist
	UpdateItem() item.PutListIDItemIDHandler
	// UpdatePatient updates with new patient data all the items with specified patientID
	UpdatePatient() item.PutPatientPatientIDHandler
	// UpdatePatient updates with new patient data all the items with specified patientID
	MoveItemToTop() item.PutListIDItemIDTopHandler
	// ReopenHistoryItem puts item from history back to waitlist
	ReopenHistoryItem() item.PutListIDItemIDReopenHandler
}

type handlers struct {
	s      Service
	logger zerolog.Logger
}

// GetWaitlists returns all active lists
func (h *handlers) GetWaitlists() waitlist.GetHandler {
	return waitlist.GetHandlerFunc(func(params waitlist.GetParams, principal *string) middleware.Responder {
		lists, err := h.s.GetWaitlists()
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewGetOK().WithPayload(lists)
	})
}

// CreateWaitlist creates new list
func (h *handlers) CreateWaitlist() waitlist.PostHandler {
	return waitlist.PostHandlerFunc(func(params waitlist.PostParams, principal *string) middleware.Responder {
		list, err := h.s.CreateWaitlist(*params.List.Name)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewPostCreated().WithPayload(list)
	})
}

// UpdateWaitlist updates list metadata
func (h *handlers) UpdateWaitlist() waitlist.PutListIDHandler {
	return waitlist.PutListIDHandlerFunc(func(params waitlist.PutListIDParams, principal *string) middleware.Responder {
		if params.ListID.String() != params.List.ID {
			return utils.NewError(utils.ErrBadRequest, "URL list ID and body list ID do not match")
		}

		_, err := h.s.UpdateWaitlist(params.List)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewPutListIDNoContent()
	})
}

// DeleteWaitlist removes list from active lists and move its items to history
func (h *handlers) DeleteWaitlist() waitlist.DeleteListIDHandler {
	return waitlist.DeleteListIDHandlerFunc(func(params waitlist.DeleteListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		err := h.s.DeleteWaitlist(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewDeleteListIDNoContent()
	})
}

// GetWaitlist returns all items in a waitlist
func (h *handlers) GetWaitlist() item.GetListIDHandler {
	return item.GetListIDHandlerFunc(func(params item.GetListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.GetWaitlist(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewGetListIDOK().WithPayload(items)
	})
}

// DeleteItem removes an item from a waitlist and moves it to history
func (h *handlers) DeleteItem() item.DeleteListIDItemIDHandler {
	return item.DeleteListIDItemIDHandlerFunc(func(params item.DeleteListIDItemIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)
		itemID, _ := utils.UUIDToBytes(params.ItemID)

		err := h.s.DeleteItem(listID, itemID, params.Reason)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewDeleteListIDItemIDNoContent()
	})
}

// CreateItem creates a new item in a waitlist
func (h *handlers) CreateItem() item.PostListIDHandler {
	return item.PostListIDHandlerFunc(func(params item.PostListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		newItem, err := h.s.CreateItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPostListIDCreated().WithPayload(newItem)
	})
}

// UpdateItem updates an item in a waitlist
func (h *handlers) UpdateItem() item.PutListIDItemIDHandler {
	return item.PutListIDItemIDHandlerFunc(func(params item.PutListIDItemIDParams, principal *string) middleware.Responder {
		if params.ItemID.String() != params.Item.ID {
			return utils.NewError(utils.ErrBadRequest, "URL item ID and body item ID do not match")
		}

		listID, _ := utils.UUIDToBytes(params.ListID)
		_, err := h.s.UpdateItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPutListIDItemIDNoContent()
	})
}

// UpdatePatient updates with new patient data all the items with specified patientID
func (h *handlers) UpdatePatient() item.PutPatientPatientIDHandler {
	return item.PutPatientPatientIDHandlerFunc(func(params item.PutPatientPatientIDParams, principal *string) middleware.Responder {
		patientID, _ := utils.UUIDToBytes(params.PatientID)
		_, err := h.s.UpdatePatient(patientID, params.Patient)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPutPatientPatientIDNoContent()
	})
}

// UpdatePatient updates with new patient data all the items with specified patientID
func (h *handlers) MoveItemToTop() item.PutListIDItemIDTopHandler {
	return item.PutListIDItemIDTopHandlerFunc(func(params item.PutListIDItemIDTopParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)
		itemID, _ := utils.UUIDToBytes(params.ItemID)

		_, err := h.s.MoveItemToTop(listID, itemID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPutListIDItemIDTopNoContent()
	})
}

// GetWaitlistHistory returns all items in waitlist's history
func (h *handlers) GetWaitlistHistory() item.GetListIDHistoryHandler {
	return item.GetListIDHistoryHandlerFunc(func(params item.GetListIDHistoryParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.GetWaitlistHistory(listID, params.Reason)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewGetListIDHistoryOK().WithPayload(items)
	})
}

// ReopenHistoryItem puts item from history back to waitlist
func (h *handlers) ReopenHistoryItem() item.PutListIDItemIDReopenHandler {
	return item.PutListIDItemIDReopenHandlerFunc(func(params item.PutListIDItemIDReopenParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)
		itemID, _ := utils.UUIDToBytes(params.ItemID)

		var newListID []byte
		if params.NewListID != nil {
			newListID, _ = utils.UUIDToBytes(*params.NewListID)
		}

		_, err := h.s.ReopenHistoryItem(listID, itemID, newListID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPutListIDItemIDReopenNoContent()
	})
}

// NewHandlers returns a new instance of waitlist handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	return &handlers{
		s:      service,
		logger: logger.With().Str("component", "service/waitlist/handlers").Logger(),
	}
}
