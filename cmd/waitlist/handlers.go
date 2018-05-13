package main

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/item"
	"github.com/iryonetwork/wwm/gen/waitlist/restapi/operations/waitlist"
	storage "github.com/iryonetwork/wwm/storage/waitlist"
	"github.com/iryonetwork/wwm/utils"
)

type handlers struct {
	s storage.Storage
}

func (h *handlers) WaitlistGet() waitlist.GetHandler {
	return waitlist.GetHandlerFunc(func(params waitlist.GetParams, principal *string) middleware.Responder {
		lists, err := h.s.Lists()
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewGetOK().WithPayload(lists)
	})
}

func (h *handlers) WaitlistPost() waitlist.PostHandler {
	return waitlist.PostHandlerFunc(func(params waitlist.PostParams, principal *string) middleware.Responder {
		list, err := h.s.AddList(*params.List.Name)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewPostCreated().WithPayload(list)
	})
}

func (h *handlers) WaitlistPutListID() waitlist.PutListIDHandler {
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

func (h *handlers) WaitlistDeleteListID() waitlist.DeleteListIDHandler {
	return waitlist.DeleteListIDHandlerFunc(func(params waitlist.DeleteListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		err := h.s.DeleteList(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return waitlist.NewDeleteListIDNoContent()
	})
}

func (h *handlers) ItemGetListID() item.GetListIDHandler {
	return item.GetListIDHandlerFunc(func(params item.GetListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		items, err := h.s.ListItems(listID)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewGetListIDOK().WithPayload(items)
	})
}

func (h *handlers) ItemDeleteListIDItemID() item.DeleteListIDItemIDHandler {
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

func (h *handlers) ItemPostListID() item.PostListIDHandler {
	return item.PostListIDHandlerFunc(func(params item.PostListIDParams, principal *string) middleware.Responder {
		listID, _ := utils.UUIDToBytes(params.ListID)

		newItem, err := h.s.AddItem(listID, params.Item)
		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return item.NewPostListIDCreated().WithPayload(newItem)
	})
}

func (h *handlers) ItemPutListIDItemID() item.PutListIDItemIDHandler {
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

func (h *handlers) ItemPutListIDItemIDTop() item.PutListIDItemIDTopHandler {
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
