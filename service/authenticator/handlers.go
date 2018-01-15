package authenticator

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/gen/restapi/operations/auth"
)

// Handlers describes the actions supported by the authenticator handlers
type Handlers interface {
	GetRenew() auth.GetRenewHandler
	PostLogin() auth.PostLoginHandler
	PostValidate() auth.PostValidateHandler
}

type handlers struct {
	service Service
}

func (h *handlers) GetRenew() auth.GetRenewHandler {
	return auth.GetRenewHandlerFunc(func(params auth.GetRenewParams, principal *string) middleware.Responder {
		token, err := h.service.CreateTokenForUserID(params.HTTPRequest.Context(), principal)
		if err != nil {
			return auth.NewGetRenewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return auth.NewGetRenewOK().WithPayload(token)
	})
}

func (h *handlers) PostLogin() auth.PostLoginHandler {
	return auth.PostLoginHandlerFunc(func(params auth.PostLoginParams) middleware.Responder {
		token, err := h.service.Login(params.HTTPRequest.Context(), *params.Login.Username, *params.Login.Password)
		if err != nil {
			return auth.NewPostLoginUnauthorized().WithPayload(&models.Error{
				Code:    "unauthorized",
				Message: err.Error(),
			})
		}

		return auth.NewPostLoginOK().WithPayload(token)
	})
}

func (h *handlers) PostValidate() auth.PostValidateHandler {
	return auth.PostValidateHandlerFunc(func(params auth.PostValidateParams, principal *string) middleware.Responder {
		result, err := h.service.Validate(params.HTTPRequest.Context(), params.Validate)

		if err != nil {
			return auth.NewPostValidateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return auth.NewPostValidateOK().WithPayload(result)
	})
}

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service) Handlers {
	return &handlers{service: service}
}
