package authenticator

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the authenticator handlers
type Handlers interface {
	// GetRenew is a handler for HTTP GET request that renews auth token
	GetRenew() operations.GetRenewHandler

	// PostLogin is a handler for HTTP POST request that logs in user and returns auth token
	PostLogin() operations.PostLoginHandler

	// PostValidate is a handler for HTTP POST request that checks if logged in user
	// has permissions to do specified queries
	PostValidate() operations.PostValidateHandler
}

type handlers struct {
	service Service
}

func (h *handlers) GetRenew() operations.GetRenewHandler {
	return operations.GetRenewHandlerFunc(func(params operations.GetRenewParams, principal *string) middleware.Responder {
		token, err := h.service.CreateTokenForUserID(params.HTTPRequest.Context(), principal)
		if err != nil {
			return utils.UseProducer(operations.NewGetRenewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			}), utils.JSONProducer)
		}

		return utils.UseProducer(operations.NewGetRenewOK().WithPayload(token), utils.TextProducer)
	})
}

func (h *handlers) PostLogin() operations.PostLoginHandler {
	return operations.PostLoginHandlerFunc(func(params operations.PostLoginParams) middleware.Responder {
		token, err := h.service.Login(params.HTTPRequest.Context(), *params.Login.Username, *params.Login.Password)
		if err != nil {
			return utils.UseProducer(operations.NewPostLoginUnauthorized().WithPayload(&models.Error{
				Code:    "unauthorized",
				Message: err.Error(),
			}), utils.JSONProducer)
		}

		return utils.UseProducer(operations.NewPostLoginOK().WithPayload(token), utils.TextProducer)
	})
}

func (h *handlers) PostValidate() operations.PostValidateHandler {
	return operations.PostValidateHandlerFunc(func(params operations.PostValidateParams, principal *string) middleware.Responder {
		result, err := h.service.Validate(params.HTTPRequest.Context(), principal, params.Validate)

		if err != nil {
			return operations.NewPostValidateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return operations.NewPostValidateOK().WithPayload(result)
	})
}

// NewHandlers returns a new instance of authenticator handlers
func NewHandlers(service Service) Handlers {
	return &handlers{service: service}
}
