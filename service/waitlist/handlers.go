package waitlist

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/item"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/waitlist"
	storage "github.com/iryonetwork/wwm/storage/waitlist"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the discovery handlers
type Handlers interface {
	GetWaitlists() waitlist.GetHandler
	CreateWaitlist() waitlist.PostHandler
	UpdateWaitlist() waitlist.PutListIDHandler
	DeleteWaitlist() waitlist.DeleteListIDHandler
	GetWaitlist() item.GetListIDHandler
	GetWaitlistHistory() item.GetListIDHistoryHandler
	DeleteItem() item.DeleteListIDItemIDHandler
	CreateItem() item.PostListIDHandler
	UpdateItem() item.PutListIDItemIDHandler
	MoveItemToTop() item.PutListIDItemIDTopHandler
	ReopenHistoryItem() item.PutListIDItemIDReopenHandler
}

type handlers struct {
	s      storage.Storage
	logger zerolog.Logger
}

func (h *handlers) GetWaitlists() waitlist.GetHandler {
	return waitlist.GetHandlerFunc(func(params waitlist.GetParams, principal *string) middleware.Responder {
		lists, err := h.s.Lists()
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewGetOK().WithPayload(lists)
	})
}

func (h *handlers) CreateWaitlist() waitlist.PostHandler {
	return waitlist.PostHandlerFunc(func(params waitlist.PostParams, principal *string) middleware.Responder {
		list, err := h.s.AddList(*params.List.Name)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewPostCreated().WithPayload(list)
	})
}

func (h *handlers) UpdateWaitlist() waitlist.PutListIDHandler {
	return waitlist.PutListIDHandlerFunc(func(params waitlist.PutListIDParams, principal *string) middleware.Responder {
		if params.ListID.String() != params.List.ID {
			return utils.NewError(utils.ErrBadRequest, "URL list ID and body list ID do not match")
		}

		_, err := h.s.UpdateList(params.List)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewPutListIDNoContent()
	})
}

func (h *handlers) DeleteWaitlist() waitlist.DeleteListIDHandler {
	return waitlist.DeleteListIDHandlerFunc(func(params waitlist.DeleteListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		err := h.s.DeleteList(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewDeleteListIDNoContent()
	})
}

func (h *handlers) GetWaitlist() item.GetListIDHandler {
	return item.GetListIDHandlerFunc(func(params item.GetListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.ListItems(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewGetListIDOK().WithPayload(items)
	})
}

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

func (h *handlers) CreateItem() item.PostListIDHandler {
	return item.PostListIDHandlerFunc(func(params item.PostListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		newItem, err := h.s.AddItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPostListIDCreated().WithPayload(newItem)
	})
}

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

func (h *handlers) GetWaitlistHistory() item.GetListIDHistoryHandler {
	return item.GetListIDHistoryHandlerFunc(func(params item.GetListIDHistoryParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.ListHistoryItems(listID, params.Reason)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewGetListIDHistoryOK().WithPayload(items)
	})
}

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
func NewHandlers(storage storage.Storage, logger zerolog.Logger) Handlers {
	return &handlers{
		s:      storage,
		logger: logger.With().Str("component", "service/waitlist/handlers").Logger(),
	}
}
