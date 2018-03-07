package main

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/go-openapi/runtime/middleware"

	"github.com/iryonetwork/wwm/gen/auth/restapi/operations/database"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations/roles"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations/rules"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations/users"
	"github.com/iryonetwork/wwm/service/accountManager"
	"github.com/iryonetwork/wwm/storage/auth"
	"github.com/iryonetwork/wwm/utils"
)

func getUsers(service accountManager.Service) users.GetUsersHandler {
	return users.GetUsersHandlerFunc(func(params users.GetUsersParams, principal *string) middleware.Responder {
		u, err := service.Users(params.HTTPRequest.Context(), "")

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return users.NewGetUsersOK().WithPayload(u)
	})
}

func getUsersID(service accountManager.Service) users.GetUsersIDHandler {
	return users.GetUsersIDHandlerFunc(func(params users.GetUsersIDParams, principal *string) middleware.Responder {
		u, err := service.User(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return users.NewGetUsersIDOK().WithPayload(u)
	})
}

func postUsers(service accountManager.Service) users.PostUsersHandler {
	return users.PostUsersHandlerFunc(func(params users.PostUsersParams, principal *string) middleware.Responder {
		u, err := service.AddUser(params.HTTPRequest.Context(), params.User)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return users.NewPostUsersCreated().WithPayload(u)
	})
}

func putUsersID(service accountManager.Service) users.PutUsersIDHandler {
	return users.PutUsersIDHandlerFunc(func(params users.PutUsersIDParams, principal *string) middleware.Responder {
		_, err := service.UpdateUser(params.HTTPRequest.Context(), params.User)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return users.NewPutUsersIDNoContent()
	})
}

func deleteUsersID(service accountManager.Service) users.DeleteUsersIDHandler {
	return users.DeleteUsersIDHandlerFunc(func(params users.DeleteUsersIDParams, principal *string) middleware.Responder {
		err := service.RemoveUser(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return users.NewDeleteUsersIDNoContent()
	})
}

func getRoles(service accountManager.Service) roles.GetRolesHandler {
	return roles.GetRolesHandlerFunc(func(params roles.GetRolesParams, principal *string) middleware.Responder {
		r, err := service.Roles(params.HTTPRequest.Context(), "")

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return roles.NewGetRolesOK().WithPayload(r)
	})
}

func getRolesID(service accountManager.Service) roles.GetRolesIDHandler {
	return roles.GetRolesIDHandlerFunc(func(params roles.GetRolesIDParams, principal *string) middleware.Responder {
		r, err := service.Role(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return roles.NewGetRolesIDOK().WithPayload(r)
	})
}

func postRoles(service accountManager.Service) roles.PostRolesHandler {
	return roles.PostRolesHandlerFunc(func(params roles.PostRolesParams, principal *string) middleware.Responder {
		r, err := service.AddRole(params.HTTPRequest.Context(), params.Role)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return roles.NewPostRolesCreated().WithPayload(r)
	})
}

func putRolesID(service accountManager.Service) roles.PutRolesIDHandler {
	return roles.PutRolesIDHandlerFunc(func(params roles.PutRolesIDParams, principal *string) middleware.Responder {
		_, err := service.UpdateRole(params.HTTPRequest.Context(), params.Role)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return roles.NewPutRolesIDNoContent()
	})
}

func deleteRolesID(service accountManager.Service) roles.DeleteRolesIDHandler {
	return roles.DeleteRolesIDHandlerFunc(func(params roles.DeleteRolesIDParams, principal *string) middleware.Responder {
		err := service.RemoveRole(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return roles.NewDeleteRolesIDNoContent()
	})
}

func getRules(service accountManager.Service) rules.GetRulesHandler {
	return rules.GetRulesHandlerFunc(func(params rules.GetRulesParams, principal *string) middleware.Responder {
		r, err := service.Rules(params.HTTPRequest.Context(), "")

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return rules.NewGetRulesOK().WithPayload(r)
	})
}

func getRulesID(service accountManager.Service) rules.GetRulesIDHandler {
	return rules.GetRulesIDHandlerFunc(func(params rules.GetRulesIDParams, principal *string) middleware.Responder {
		r, err := service.Rule(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return rules.NewGetRulesIDOK().WithPayload(r)
	})
}

func postRules(service accountManager.Service) rules.PostRulesHandler {
	return rules.PostRulesHandlerFunc(func(params rules.PostRulesParams, principal *string) middleware.Responder {
		r, err := service.AddRule(params.HTTPRequest.Context(), params.Rule)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return rules.NewPostRulesCreated().WithPayload(r)
	})
}

func putRulesID(service accountManager.Service) rules.PutRulesIDHandler {
	return rules.PutRulesIDHandlerFunc(func(params rules.PutRulesIDParams, principal *string) middleware.Responder {
		_, err := service.UpdateRule(params.HTTPRequest.Context(), params.Rule)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return rules.NewPutRulesIDNoContent()
	})
}

func deleteRulesID(service accountManager.Service) rules.DeleteRulesIDHandler {
	return rules.DeleteRulesIDHandlerFunc(func(params rules.DeleteRulesIDParams, principal *string) middleware.Responder {
		err := service.RemoveRule(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return rules.NewDeleteRulesIDNoContent()
	})
}

func getDatabase(storage *auth.Storage) database.GetDatabaseHandler {
	return database.GetDatabaseHandlerFunc(func(params database.GetDatabaseParams, principal *string) middleware.Responder {
		etag := strings.Trim(params.HTTPRequest.Header.Get("Etag"), `"`)

		checksum, err := storage.GetChecksum()
		if err != nil {
			return utils.UseProducer(utils.NewServerError(err), utils.JSONProducer)
		}

		currentEtag := base64.RawURLEncoding.EncodeToString(checksum)

		if etag == currentEtag {
			return database.NewGetDatabaseNotModified()
		}

		reader, writer := io.Pipe()

		go func() {
			_, err := storage.WriteTo(writer)
			writer.CloseWithError(err)
		}()

		return utils.UseProducer(
			database.NewGetDatabaseOK().
				WithPayload(reader).
				WithEtag(`"`+currentEtag+`"`),
			utils.BinProducer)
	})
}
