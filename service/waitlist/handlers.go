package waitlist

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the discovery handlers
type Handlers interface {
	// GetWaitlists returns all active lists
	GetWaitlists() operations.GetHandler
	// CreateWaitlist creates new list
	CreateWaitlist() operations.PostHandler
	// UpdateWaitlist updates list metadata
	UpdateWaitlist() operations.PutListIDHandler
	// DeleteWaitlist removes list from active lists and move its items to history
	DeleteWaitlist() operations.DeleteListIDHandler
	// GetWaitlist returns all items in a waitlist
	GetWaitlist() operations.GetListIDHandler
	// GetWaitlistHistory returns all items in waitlist's history
	GetWaitlistHistory() operations.GetListIDHistoryHandler
	// DeleteItem removes an item from a waitlist and moves it to history
	DeleteItem() operations.DeleteListIDItemIDHandler
	// CreateItem creates a new item in a waitlist
	CreateItem() operations.PostListIDHandler
	// UpdateItem updates an item in a waitlist
	UpdateItem() operations.PutListIDItemIDHandler
	// UpdatePatient updates with new patient data all the items with specified patientID
	UpdatePatient() operations.PutPatientPatientIDHandler
	// UpdatePatient updates with new patient data all the items with specified patientID
	MoveItemToTop() operations.PutListIDItemIDTopHandler
	// ReopenHistoryItem puts item from history back to waitlist
	ReopenHistoryItem() operations.PutListIDItemIDReopenHandler
}

type handlers struct {
	s      Service
	logger zerolog.Logger
}

// GetWaitlists returns all active lists
func (h *handlers) GetWaitlists() operations.GetHandler {
	return operations.GetHandlerFunc(func(params operations.GetParams, principal *string) middleware.Responder {
		lists, err := h.s.GetWaitlists()
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetOK().WithPayload(lists)
	})
}

// CreateWaitlist creates new list
func (h *handlers) CreateWaitlist() operations.PostHandler {
	return operations.PostHandlerFunc(func(params operations.PostParams, principal *string) middleware.Responder {
		list, err := h.s.CreateWaitlist(*params.List.Name)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostCreated().WithPayload(list)
	})
}

// UpdateWaitlist updates list metadata
func (h *handlers) UpdateWaitlist() operations.PutListIDHandler {
	return operations.PutListIDHandlerFunc(func(params operations.PutListIDParams, principal *string) middleware.Responder {
		if params.ListID.String() != params.List.ID {
			return utils.NewError(utils.ErrBadRequest, "URL list ID and body list ID do not match")
		}

		_, err := h.s.UpdateWaitlist(params.List)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutListIDNoContent()
	})
}

// DeleteWaitlist removes list from active lists and move its items to history
func (h *handlers) DeleteWaitlist() operations.DeleteListIDHandler {
	return operations.DeleteListIDHandlerFunc(func(params operations.DeleteListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		err := h.s.DeleteWaitlist(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteListIDNoContent()
	})
}

// GetWaitlist returns all items in a waitlist
func (h *handlers) GetWaitlist() operations.GetListIDHandler {
	return operations.GetListIDHandlerFunc(func(params operations.GetListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.GetWaitlist(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetListIDOK().WithPayload(items)
	})
}

// DeleteItem removes an item from a waitlist and moves it to history
func (h *handlers) DeleteItem() operations.DeleteListIDItemIDHandler {
	return operations.DeleteListIDItemIDHandlerFunc(func(params operations.DeleteListIDItemIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)
		itemID, _ := utils.UUIDToBytes(params.ItemID)

		err := h.s.DeleteItem(listID, itemID, params.Reason)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteListIDItemIDNoContent()
	})
}

// CreateItem creates a new item in a waitlist
func (h *handlers) CreateItem() operations.PostListIDHandler {
	return operations.PostListIDHandlerFunc(func(params operations.PostListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		newItem, err := h.s.CreateItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostListIDCreated().WithPayload(newItem)
	})
}

// UpdateItem updates an item in a waitlist
func (h *handlers) UpdateItem() operations.PutListIDItemIDHandler {
	return operations.PutListIDItemIDHandlerFunc(func(params operations.PutListIDItemIDParams, principal *string) middleware.Responder {
		if params.ItemID.String() != params.Item.ID {
			return utils.NewError(utils.ErrBadRequest, "URL item ID and body item ID do not match")
		}

		listID, _ := utils.UUIDToBytes(params.ListID)
		_, err := h.s.UpdateItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutListIDItemIDNoContent()
	})
}

// UpdatePatient updates with new patient data all the items with specified patientID
func (h *handlers) UpdatePatient() operations.PutPatientPatientIDHandler {
	return operations.PutPatientPatientIDHandlerFunc(func(params operations.PutPatientPatientIDParams, principal *string) middleware.Responder {
		patientID, _ := utils.UUIDToBytes(params.PatientID)
		_, err := h.s.UpdatePatient(patientID, params.Patient)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutPatientPatientIDNoContent()
	})
}

// UpdatePatient updates with new patient data all the items with specified patientID
func (h *handlers) MoveItemToTop() operations.PutListIDItemIDTopHandler {
	return operations.PutListIDItemIDTopHandlerFunc(func(params operations.PutListIDItemIDTopParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)
		itemID, _ := utils.UUIDToBytes(params.ItemID)

		_, err := h.s.MoveItemToTop(listID, itemID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutListIDItemIDTopNoContent()
	})
}

// GetWaitlistHistory returns all items in waitlist's history
func (h *handlers) GetWaitlistHistory() operations.GetListIDHistoryHandler {
	return operations.GetListIDHistoryHandlerFunc(func(params operations.GetListIDHistoryParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.GetWaitlistHistory(listID, params.Reason)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetListIDHistoryOK().WithPayload(items)
	})
}

// ReopenHistoryItem puts item from history back to waitlist
func (h *handlers) ReopenHistoryItem() operations.PutListIDItemIDReopenHandler {
	return operations.PutListIDItemIDReopenHandlerFunc(func(params operations.PutListIDItemIDReopenParams, principal *string) middleware.Responder {
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

		return operations.NewPutListIDItemIDReopenNoContent()
	})
}

// NewHandlers returns a new instance of waitlist handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	return &handlers{
		s:      service,
		logger: logger.With().Str("component", "service/waitlist/handlers").Logger(),
	}
}
