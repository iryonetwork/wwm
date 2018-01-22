package authenticator

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations/auth"
	"github.com/iryonetwork/wwm/utils"
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
			return utils.UseProducer(auth.NewGetRenewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			}), utils.JSONProducer)
		}

		return utils.UseProducer(auth.NewGetRenewOK().WithPayload(token), utils.TextProducer)
	})
}

func (h *handlers) PostLogin() auth.PostLoginHandler {
	return auth.PostLoginHandlerFunc(func(params auth.PostLoginParams) middleware.Responder {
		token, err := h.service.Login(params.HTTPRequest.Context(), *params.Login.Username, *params.Login.Password)
		if err != nil {
			return utils.UseProducer(auth.NewPostLoginUnauthorized().WithPayload(&models.Error{
				Code:    "unauthorized",
				Message: err.Error(),
			}), utils.JSONProducer)
		}

		return utils.UseProducer(auth.NewPostLoginOK().WithPayload(token), utils.TextProducer)
	})
}

func (h *handlers) PostValidate() auth.PostValidateHandler {
	return auth.PostValidateHandlerFunc(func(params auth.PostValidateParams, principal *string) middleware.Responder {
		result, err := h.service.Validate(params.HTTPRequest.Context(), principal, params.Validate)

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
