package main

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/gen/restapi/operations/auth"
	"github.com/iryonetwork/wwm/service/authenticator"
)

func getRenew(svc authenticator.Service) auth.GetRenewHandler {
	return auth.GetRenewHandlerFunc(func(params auth.GetRenewParams, principal *models.User) middleware.Responder {
		token, err := svc.CreateTokenForUser(params.HTTPRequest.Context(), principal)
		if err != nil {
			return auth.NewGetRenewInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return auth.NewGetRenewOK().WithPayload(token)
	})
}

func postLogin(svc authenticator.Service) auth.PostLoginHandler {
	return auth.PostLoginHandlerFunc(func(params auth.PostLoginParams) middleware.Responder {
		token, err := svc.Login(params.HTTPRequest.Context(), *params.Login.Username, *params.Login.Password)
		if err != nil {
			return auth.NewPostLoginUnauthorized().WithPayload(&models.Error{
				Code:    "unauthorized",
				Message: err.Error(),
			})
		}

		return auth.NewPostLoginOK().WithPayload(token)
	})
}

func postValidateHandler(svc authenticator.Service) auth.PostValidateHandler {
	return auth.PostValidateHandlerFunc(func(params auth.PostValidateParams, principal *models.User) middleware.Responder {
		result, err := svc.Validate(params.HTTPRequest.Context(), params.Validate)

		if err != nil {
			return auth.NewPostValidateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return auth.NewPostValidateOK().WithPayload(result)
	})
}
