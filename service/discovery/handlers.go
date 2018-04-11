package discovery

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/iryonetwork/wwm/gen/discovery/models"
	"github.com/iryonetwork/wwm/gen/discovery/restapi/operations"
	"github.com/iryonetwork/wwm/storage/discovery"
	"github.com/rs/zerolog"
)

// Handlers describes the actions supported by the discovery handlers
type Handlers interface {
	Query() operations.QueryHandler
	Create() operations.CreateHandler
	Update() operations.UpdateHandler
	Delete() operations.DeleteHandler
	Fetch() operations.FetchHandler
	Link() operations.LinkHandler
	ProxyLink() operations.LinkHandler
	Unlink() operations.UnlinkHandler
	ProxyUnlink() operations.UnlinkHandler
	Authorizer() runtime.Authorizer
	GetUserIDFromToken(token string) (*string, error)
}

type handlers struct {
	service Service
	logger  zerolog.Logger
}

func (h *handlers) Query() operations.QueryHandler {
	return operations.QueryHandlerFunc(func(params operations.QueryParams, principal *string) middleware.Responder {
		q := swag.StringValue(params.Query)

		var (
			res models.Cards
			err error
		)

		if params.OnCloud != nil && *params.OnCloud == true {
			res, err = h.service.ProxyQuery(q, params.HTTPRequest.Header.Get("Authorization"))
		} else {
			res, err = h.service.Query(q)
		}

		if err != nil {
			return operations.NewQueryInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewQueryOK().WithPayload(res)
	})
}

func (h *handlers) Create() operations.CreateHandler {
	return operations.CreateHandlerFunc(func(params operations.CreateParams, principal *string) middleware.Responder {
		c, err := h.service.Create(params.NewCard)
		if err != nil {
			return operations.NewCreateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewCreateCreated().WithPayload(c)
	})
}

func (h *handlers) Update() operations.UpdateHandler {
	return operations.UpdateHandlerFunc(func(params operations.UpdateParams, principal *string) middleware.Responder {
		c, err := h.service.Update(params.PatientID, params.Card)
		if err == discovery.ErrNotFound {
			return operations.NewUpdateNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewUpdateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewUpdateOK().WithPayload(c)
	})
}

func (h *handlers) Delete() operations.DeleteHandler {
	return operations.DeleteHandlerFunc(func(params operations.DeleteParams, principal *string) middleware.Responder {
		err := h.service.Delete(params.PatientID)
		if err == discovery.ErrNotFound {
			return operations.NewDeleteNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewDeleteInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewDeleteNoContent()
	})
}

func (h *handlers) Fetch() operations.FetchHandler {
	return operations.FetchHandlerFunc(func(params operations.FetchParams, principal *string) middleware.Responder {
		c, err := h.service.Fetch(params.PatientID)
		if err == discovery.ErrNotFound {
			return operations.NewFetchNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewFetchInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewFetchOK().WithPayload(c)
	})
}

func (h *handlers) Link() operations.LinkHandler {
	return operations.LinkHandlerFunc(func(params operations.LinkParams, principal *string) middleware.Responder {
		l, err := h.service.Link(params.PatientID, params.LocationID)

		if err == discovery.ErrNotFound {
			return operations.NewLinkNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewLinkInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewLinkCreated().WithPayload(l)
	})
}

func (h *handlers) ProxyLink() operations.LinkHandler {
	return operations.LinkHandlerFunc(func(params operations.LinkParams, principal *string) middleware.Responder {
		authToken := params.HTTPRequest.Header.Get("Authorization")
		l, err := h.service.ProxyLink(params.PatientID, params.LocationID, authToken)
		h.logger.Debug().Msgf("Proxy link result: %+v, %+v", l, err)

		if err == ErrNotFound {
			return operations.NewLinkNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewLinkInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewLinkCreated().WithPayload(l)
	})
}

func (h *handlers) Unlink() operations.UnlinkHandler {
	return operations.UnlinkHandlerFunc(func(params operations.UnlinkParams, principal *string) middleware.Responder {
		err := h.service.Unlink(params.PatientID, params.LocationID)

		if err == discovery.ErrNotFound {
			return operations.NewUnlinkNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewUnlinkInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewUnlinkNoContent()
	})
}

func (h *handlers) ProxyUnlink() operations.UnlinkHandler {
	return operations.UnlinkHandlerFunc(func(params operations.UnlinkParams, principal *string) middleware.Responder {
		authToken := params.HTTPRequest.Header.Get("Authorization")
		err := h.service.ProxyUnlink(params.PatientID, params.LocationID, authToken)

		if err == ErrNotFound {
			return operations.NewUnlinkNotFound().WithPayload(&models.Error{
				Code:    "not_found",
				Message: err.Error(),
			})
		} else if err != nil {
			return operations.NewUnlinkInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}
		return operations.NewUnlinkNoContent()
	})
}

func (h *handlers) Authorizer() runtime.Authorizer {
	return runtime.AuthorizerFunc(func(r *http.Request, principal interface{}) error {
		// @TODO
		return nil
	})
}

func (h *handlers) GetUserIDFromToken(token string) (*string, error) {
	userID := "USER"
	return &userID, nil
}

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service, logger zerolog.Logger) Handlers {
	return &handlers{
		service: service,
		logger:  logger.With().Str("component", "service/discovery/handlers").Logger(),
	}
}
