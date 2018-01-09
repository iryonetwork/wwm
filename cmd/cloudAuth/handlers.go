package main

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/iryonetwork/wwm/gen/models"
	"github.com/iryonetwork/wwm/gen/restapi/operations/auth"
	"github.com/iryonetwork/wwm/service/authenticator"
)

func getAuthRenew(svc authenticator.Service) auth.GetAuthRenewHandler {
	return auth.GetAuthRenewHandlerFunc(func(params auth.GetAuthRenewParams, principal *models.User) middleware.Responder {
		return middleware.NotImplemented("aoperation auth.GetAuthRenew has not yet been implemented")
	})
}

func postAuthLogin(svc authenticator.Service) auth.PostAuthLoginHandler {
	return auth.PostAuthLoginHandlerFunc(func(params auth.PostAuthLoginParams) middleware.Responder {
		token, err := svc.Login(params.HTTPRequest.Context(), *params.Login.Username, *params.Login.Password)
		if err != nil {
			return auth.NewPostAuthLoginUnauthorized().WithPayload(&models.Error{
				Code:    "unauthorized",
				Message: err.Error(),
			})
		}

		return auth.NewPostAuthLoginOK().WithPayload(token)
	})
}

func postAuthValidateHandler(svc authenticator.Service) auth.PostAuthValidateHandler {
	return auth.PostAuthValidateHandlerFunc(func(params auth.PostAuthValidateParams, principal *models.User) middleware.Responder {
		result, err := svc.Validate(params.HTTPRequest.Context(), params.Validate)

		if err != nil {
			return auth.NewPostAuthValidateInternalServerError().WithPayload(&models.Error{
				Code:    "server_error",
				Message: err.Error(),
			})
		}

		return auth.NewPostAuthValidateOK().WithPayload(result)
	})
}
