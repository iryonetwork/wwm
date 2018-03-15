package authDataManager

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/go-openapi/runtime/middleware"

	authCommon "github.com/iryonetwork/wwm/auth"
	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/gen/auth/restapi/operations"
	"github.com/iryonetwork/wwm/utils"
)

// Handlers describes the actions supported by the authDataManager handlers.
type Handlers interface {
	// GetUsers is a handler for HTTP GET request that fetches list of all users.
	GetUsers() operations.GetUsersHandler

	// GetUsersID is a handler for HTTP GET request that fetches the user based on user ID.
	GetUsersID() operations.GetUsersIDHandler

	// GetUsersIDRoles is a handler for HTTP GET request that fetches IDs of roles that the user has been assigned (with optional domain filtering).
	GetUsersIDRoles() operations.GetUsersIDRolesHandler

	// GetUsersIDOrganizations is a handler for HTTP GET request that fetches IDs of organizations at which the user has been assigned a role (with optional role ID filtering).
	GetUsersIDOrganizations() operations.GetUsersIDOrganizationsHandler

	// GetUsersIDClinics is a handler for HTTP GET request that fetches IDs of clinics at which the user has been assigned a role (with optional role ID filtering).
	GetUsersIDClinics() operations.GetUsersIDClinicsHandler

	// GetUsersIDLocations is a handler for HTTP GET request that fetches IDs of locations at which the user has been assigned a role (with optional role ID filtering); both locations of clinics and locations at which user has been assigned a role manually are returned.
	GetUsersIDLocations() operations.GetUsersIDLocationsHandler

	// GetUsersMe is a handler for HTTP GET request that fetches currently logged-in user.
	GetUsersMe() operations.GetUsersMeHandler

	// GetUsersMeRoles is a handler for HTTP GET request that fetches roles that currently logged-in user has been assigned (with optional domain filtering).
	GetUsersMeRoles() operations.GetUsersMeRolesHandler

	// GetUsersMeOrganizations is a handler for HTTP GET request that fetches IDs of organizations at which currently logged-in user has been assigned a role (with optional role ID filtering).
	GetUsersMeOrganizations() operations.GetUsersMeOrganizationsHandler

	// GetUsersMeClinics is a handler for HTTP GET request that fetches IDs of clinics at which currently logged-in user has been assigned a role (with optional role ID filtering).
	GetUsersMeClinics() operations.GetUsersMeClinicsHandler

	// GetUsersMeLocations is a handler for HTTP GET request that fetches IDs of locations at which currently logged-in user has been assigned a role (with optional role ID filtering); both locations of clinics and locations at which user has been assigned a role manually are returned.
	GetUsersMeLocations() operations.GetUsersMeLocationsHandler

	// PostValidate is a handler for HTTP POST request that creates a new user.
	PostUsers() operations.PostUsersHandler

	// PutUserID is a handler for HTTP PUT request that updates user identified by user ID.
	PutUsersID() operations.PutUsersIDHandler

	// DeleteUserID is a handler for HTTP DELETE request that deletes user identified by user ID.
	DeleteUsersID() operations.DeleteUsersIDHandler

	// GetRoles is a handler for HTTP GET request that fetches list of all roles.
	GetRoles() operations.GetRolesHandler

	// GetRolesID is a handler for HTTP GET request that fetches the role based on role ID.
	GetRolesID() operations.GetRolesIDHandler

	// GetRolesIDUsers is a handler for HTTP GET request that fetches list of IDs of users that have been assigned the role (with optional domain filtering).
	GetRolesIDUsers() operations.GetRolesIDUsersHandler

	// PostRoles is a handler for HTTP POST request that creates a new role.
	PostRoles() operations.PostRolesHandler

	// PutUserID is a handler for HTTP PUT request that updates role identified by role ID.
	PutRolesID() operations.PutRolesIDHandler

	// DeleteUserID is a handler for HTTP DELETE request that deletes role identified by role ID.
	DeleteRolesID() operations.DeleteRolesIDHandler

	// GetRules is a handler for HTTP GET request that fetches list of all rules.
	GetRules() operations.GetRulesHandler

	// GetRulesID is a handler for HTTP GET request that fetches the rule based on rule ID.
	GetRulesID() operations.GetRulesIDHandler

	// PostValidate is a handler for HTTP POST request that creates a new rule.
	PostRules() operations.PostRulesHandler

	// PutUserID is a handler for HTTP PUT request that updates rule identified by rule ID.
	PutRulesID() operations.PutRulesIDHandler

	// DeleteUserID is a handler for HTTP DELETE request that deletes rule identified by rule ID.
	DeleteRulesID() operations.DeleteRulesIDHandler

	// GetClinics is a handler for HTTP GET request that fetches list of all clinics.
	GetClinics() operations.GetClinicsHandler

	// GetClinicsID is a handler for HTTP GET request that fetches the clinic based on clinic ID.
	GetClinicsID() operations.GetClinicsIDHandler

	// GetClinicsIDUsers is a handler for HTTP GET request that fetches list of IDs of users that have been assigned a role at the clinic (with optional role ID filtering).
	GetClinicsIDUsers() operations.GetClinicsIDUsersHandler

	// PostValidate is a handler for HTTP POST request that creates a new clinic.
	PostClinics() operations.PostClinicsHandler

	// PutClinicID is a handler for HTTP PUT request that updates clinic identified by clinic ID.
	PutClinicsID() operations.PutClinicsIDHandler

	// DeleteClinicID is a handler for HTTP DELETE request that deletes clinic identified by clinic ID.
	DeleteClinicsID() operations.DeleteClinicsIDHandler

	// GetLocations is a handler for HTTP GET request that fetches list of all locations.
	GetLocations() operations.GetLocationsHandler

	// GetLocationsID is a handler for HTTP GET request that fetches the location based on location ID.
	GetLocationsID() operations.GetLocationsIDHandler

	// GetLocationsIDOrganizations is a handler for HTTP GET request that fetches IDs of location's organizations based on location ID.
	GetLocationsIDOrganizations() operations.GetLocationsIDOrganizationsHandler

	// GetLocationsIDUsers is a handler for HTTP GET request that fetches list of IDs of users that have been assigned a role at the location (with optional role ID filtering); both users of clinics associated with the locations and users that have been assigned a role at the location manually are returned.
	GetLocationsIDUsers() operations.GetLocationsIDUsersHandler

	// PostValidate is a handler for HTTP POST request that creates a new location.
	PostLocations() operations.PostLocationsHandler

	// PutLocationID is a handler for HTTP PUT request that updates location identified by location ID.
	PutLocationsID() operations.PutLocationsIDHandler

	// DeleteLocationID is a handler for HTTP DELETE request that deletes location identified by location ID.
	DeleteLocationsID() operations.DeleteLocationsIDHandler

	// GetOrganizations is a handler for HTTP GET request that fetches list of all organizations.
	GetOrganizations() operations.GetOrganizationsHandler

	// GetOrganizationsID is a handler for HTTP GET request that fetches the organization based on organization ID.
	GetOrganizationsID() operations.GetOrganizationsIDHandler

	// GetOrganizationsIDLocations is a handler for HTTP GET request that fetches IDs of organization's locations based on organization ID.
	GetOrganizationsIDLocations() operations.GetOrganizationsIDLocationsHandler

	// GetOrganizationsIDUsers is a handler for HTTP GET request that fetches list of IDs of users that have been assigned a role at the organization (with optional role ID filtering).
	GetOrganizationsIDUsers() operations.GetOrganizationsIDUsersHandler

	// PostValidate is a handler for HTTP POST request that creates a new organization.
	PostOrganizations() operations.PostOrganizationsHandler

	// PutOrganizationID is a handler for HTTP PUT request that updates organization identified by organization ID.
	PutOrganizationsID() operations.PutOrganizationsIDHandler

	// DeleteOrganizationID is a handler for HTTP DELETE request that deletes organization identified by organization ID.
	DeleteOrganizationsID() operations.DeleteOrganizationsIDHandler

	// GetUserRoles is a handler for HTTP GET request that fetches list of user roles based on filtering query parameters.
	GetUserRoles() operations.GetUserRolesHandler

	// GetUserRolesID is a handler for HTTP GET request that fetches the user role based on user role ID.
	GetUserRolesID() operations.GetUserRolesIDHandler

	// PostUserRoles is a handler for HTTP POST request that creates a new user role.
	PostUserRoles() operations.PostUserRolesHandler

	// DeleteUserRolesID is a handler for HTTP DELETE request that deletes user role identified by user role ID.
	DeleteUserRolesID() operations.DeleteUserRolesIDHandler

	// GetDatabase is a handler for HTTP GET request that fetches whole database.
	GetDatabase() operations.GetDatabaseHandler
}

type handlers struct {
	service Service
}

func (h *handlers) GetUsers() operations.GetUsersHandler {
	return operations.GetUsersHandlerFunc(func(params operations.GetUsersParams, principal *string) middleware.Responder {
		u, err := h.service.Users(params.HTTPRequest.Context())

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersID() operations.GetUsersIDHandler {
	return operations.GetUsersIDHandlerFunc(func(params operations.GetUsersIDParams, principal *string) middleware.Responder {
		u, err := h.service.User(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersIDOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersIDRoles() operations.GetUsersIDRolesHandler {
	return operations.GetUsersIDRolesHandlerFunc(func(params operations.GetUsersIDRolesParams, principal *string) middleware.Responder {
		u, err := h.service.UserRoleIDs(params.HTTPRequest.Context(), params.ID, params.DomainType, params.DomainID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersIDRolesOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersIDOrganizations() operations.GetUsersIDOrganizationsHandler {
	return operations.GetUsersIDOrganizationsHandlerFunc(func(params operations.GetUsersIDOrganizationsParams, principal *string) middleware.Responder {
		u, err := h.service.UserOrganizationIDs(params.HTTPRequest.Context(), params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersIDOrganizationsOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersIDClinics() operations.GetUsersIDClinicsHandler {
	return operations.GetUsersIDClinicsHandlerFunc(func(params operations.GetUsersIDClinicsParams, principal *string) middleware.Responder {
		u, err := h.service.UserClinicIDs(params.HTTPRequest.Context(), params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersIDClinicsOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersIDLocations() operations.GetUsersIDLocationsHandler {
	return operations.GetUsersIDLocationsHandlerFunc(func(params operations.GetUsersIDLocationsParams, principal *string) middleware.Responder {
		u, err := h.service.UserLocationIDs(params.HTTPRequest.Context(), params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersIDLocationsOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersMe() operations.GetUsersMeHandler {
	return operations.GetUsersMeHandlerFunc(func(params operations.GetUsersMeParams, principal *string) middleware.Responder {
		u, err := h.service.User(params.HTTPRequest.Context(), *principal)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersMeOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersMeRoles() operations.GetUsersMeRolesHandler {
	return operations.GetUsersMeRolesHandlerFunc(func(params operations.GetUsersMeRolesParams, principal *string) middleware.Responder {
		u, err := h.service.UserRoleIDs(params.HTTPRequest.Context(), *principal, params.DomainType, params.DomainID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersMeRolesOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersMeOrganizations() operations.GetUsersMeOrganizationsHandler {
	return operations.GetUsersMeOrganizationsHandlerFunc(func(params operations.GetUsersMeOrganizationsParams, principal *string) middleware.Responder {
		u, err := h.service.UserOrganizationIDs(params.HTTPRequest.Context(), *principal, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersMeOrganizationsOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersMeClinics() operations.GetUsersMeClinicsHandler {
	return operations.GetUsersMeClinicsHandlerFunc(func(params operations.GetUsersMeClinicsParams, principal *string) middleware.Responder {
		u, err := h.service.UserClinicIDs(params.HTTPRequest.Context(), *principal, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersMeClinicsOK().WithPayload(u)
	})
}

func (h *handlers) GetUsersMeLocations() operations.GetUsersMeLocationsHandler {
	return operations.GetUsersMeLocationsHandlerFunc(func(params operations.GetUsersMeLocationsParams, principal *string) middleware.Responder {
		u, err := h.service.UserLocationIDs(params.HTTPRequest.Context(), *principal, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUsersMeLocationsOK().WithPayload(u)
	})
}

func (h *handlers) PostUsers() operations.PostUsersHandler {
	return operations.PostUsersHandlerFunc(func(params operations.PostUsersParams, principal *string) middleware.Responder {
		u, err := h.service.AddUser(params.HTTPRequest.Context(), params.User)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostUsersCreated().WithPayload(u)
	})
}

func (h *handlers) PutUsersID() operations.PutUsersIDHandler {
	return operations.PutUsersIDHandlerFunc(func(params operations.PutUsersIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateUser(params.HTTPRequest.Context(), params.User)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutUsersIDNoContent()
	})
}

func (h *handlers) DeleteUsersID() operations.DeleteUsersIDHandler {
	return operations.DeleteUsersIDHandlerFunc(func(params operations.DeleteUsersIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveUser(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteUsersIDNoContent()
	})
}

func (h *handlers) GetRoles() operations.GetRolesHandler {
	return operations.GetRolesHandlerFunc(func(params operations.GetRolesParams, principal *string) middleware.Responder {
		r, err := h.service.Roles(params.HTTPRequest.Context())

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetRolesOK().WithPayload(r)
	})
}

func (h *handlers) GetRolesID() operations.GetRolesIDHandler {
	return operations.GetRolesIDHandlerFunc(func(params operations.GetRolesIDParams, principal *string) middleware.Responder {
		r, err := h.service.Role(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetRolesIDOK().WithPayload(r)
	})
}

func (h *handlers) GetRolesIDUsers() operations.GetRolesIDUsersHandler {
	return operations.GetRolesIDUsersHandlerFunc(func(params operations.GetRolesIDUsersParams, principal *string) middleware.Responder {
		r, err := h.service.RoleUserIDs(params.HTTPRequest.Context(), params.ID, params.DomainType, params.DomainID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetRolesIDUsersOK().WithPayload(r)
	})
}

func (h *handlers) PostRoles() operations.PostRolesHandler {
	return operations.PostRolesHandlerFunc(func(params operations.PostRolesParams, principal *string) middleware.Responder {
		r, err := h.service.AddRole(params.HTTPRequest.Context(), params.Role)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostRolesCreated().WithPayload(r)
	})
}

func (h *handlers) PutRolesID() operations.PutRolesIDHandler {
	return operations.PutRolesIDHandlerFunc(func(params operations.PutRolesIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateRole(params.HTTPRequest.Context(), params.Role)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutRolesIDNoContent()
	})
}

func (h *handlers) DeleteRolesID() operations.DeleteRolesIDHandler {
	return operations.DeleteRolesIDHandlerFunc(func(params operations.DeleteRolesIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveRole(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteRolesIDNoContent()
	})
}

func (h *handlers) GetRules() operations.GetRulesHandler {
	return operations.GetRulesHandlerFunc(func(params operations.GetRulesParams, principal *string) middleware.Responder {
		r, err := h.service.Rules(params.HTTPRequest.Context())

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetRulesOK().WithPayload(r)
	})
}

func (h *handlers) GetRulesID() operations.GetRulesIDHandler {
	return operations.GetRulesIDHandlerFunc(func(params operations.GetRulesIDParams, principal *string) middleware.Responder {
		r, err := h.service.Rule(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetRulesIDOK().WithPayload(r)
	})
}

func (h *handlers) PostRules() operations.PostRulesHandler {
	return operations.PostRulesHandlerFunc(func(params operations.PostRulesParams, principal *string) middleware.Responder {
		r, err := h.service.AddRule(params.HTTPRequest.Context(), params.Rule)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostRulesCreated().WithPayload(r)
	})
}

func (h *handlers) PutRulesID() operations.PutRulesIDHandler {
	return operations.PutRulesIDHandlerFunc(func(params operations.PutRulesIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateRule(params.HTTPRequest.Context(), params.Rule)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutRulesIDNoContent()
	})
}

func (h *handlers) DeleteRulesID() operations.DeleteRulesIDHandler {
	return operations.DeleteRulesIDHandlerFunc(func(params operations.DeleteRulesIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveRule(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteRulesIDNoContent()
	})
}

func (h *handlers) GetClinics() operations.GetClinicsHandler {
	return operations.GetClinicsHandlerFunc(func(params operations.GetClinicsParams, principal *string) middleware.Responder {
		u, err := h.service.Clinics(params.HTTPRequest.Context())

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetClinicsOK().WithPayload(u)
	})
}

func (h *handlers) GetClinicsID() operations.GetClinicsIDHandler {
	return operations.GetClinicsIDHandlerFunc(func(params operations.GetClinicsIDParams, principal *string) middleware.Responder {
		u, err := h.service.Clinic(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetClinicsIDOK().WithPayload(u)
	})
}

func (h *handlers) GetClinicsIDUsers() operations.GetClinicsIDUsersHandler {
	return operations.GetClinicsIDUsersHandlerFunc(func(params operations.GetClinicsIDUsersParams, principal *string) middleware.Responder {
		u, err := h.service.DomainUserIDs(params.HTTPRequest.Context(), &authCommon.DomainTypeClinic, &params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetClinicsIDUsersOK().WithPayload(u)
	})
}

func (h *handlers) PostClinics() operations.PostClinicsHandler {
	return operations.PostClinicsHandlerFunc(func(params operations.PostClinicsParams, principal *string) middleware.Responder {
		u, err := h.service.AddClinic(params.HTTPRequest.Context(), params.Clinic)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostClinicsCreated().WithPayload(u)
	})
}

func (h *handlers) PutClinicsID() operations.PutClinicsIDHandler {
	return operations.PutClinicsIDHandlerFunc(func(params operations.PutClinicsIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateClinic(params.HTTPRequest.Context(), params.Clinic)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutClinicsIDNoContent()
	})
}

func (h *handlers) DeleteClinicsID() operations.DeleteClinicsIDHandler {
	return operations.DeleteClinicsIDHandlerFunc(func(params operations.DeleteClinicsIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveClinic(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteClinicsIDNoContent()
	})
}

func (h *handlers) GetLocations() operations.GetLocationsHandler {
	return operations.GetLocationsHandlerFunc(func(params operations.GetLocationsParams, principal *string) middleware.Responder {
		u, err := h.service.Locations(params.HTTPRequest.Context())

		if err != nil {
			return operations.NewGetLocationsInternalServerError().WithPayload(&models.Error{
				Code:    utils.ErrServerError,
				Message: err.Error(),
			})
		}

		return operations.NewGetLocationsOK().WithPayload(u)
	})
}

func (h *handlers) GetLocationsID() operations.GetLocationsIDHandler {
	return operations.GetLocationsIDHandlerFunc(func(params operations.GetLocationsIDParams, principal *string) middleware.Responder {
		u, err := h.service.Location(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetLocationsIDOK().WithPayload(u)
	})
}

func (h *handlers) GetLocationsIDOrganizations() operations.GetLocationsIDOrganizationsHandler {
	return operations.GetLocationsIDOrganizationsHandlerFunc(func(params operations.GetLocationsIDOrganizationsParams, principal *string) middleware.Responder {
		u, err := h.service.LocationOrganizationIDs(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetLocationsIDOrganizationsOK().WithPayload(u)
	})
}

func (h *handlers) GetLocationsIDUsers() operations.GetLocationsIDUsersHandler {
	return operations.GetLocationsIDUsersHandlerFunc(func(params operations.GetLocationsIDUsersParams, principal *string) middleware.Responder {
		r, err := h.service.LocationUserIDs(params.HTTPRequest.Context(), params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetLocationsIDUsersOK().WithPayload(r)
	})
}

func (h *handlers) PostLocations() operations.PostLocationsHandler {
	return operations.PostLocationsHandlerFunc(func(params operations.PostLocationsParams, principal *string) middleware.Responder {
		u, err := h.service.AddLocation(params.HTTPRequest.Context(), params.Location)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostLocationsCreated().WithPayload(u)
	})
}

func (h *handlers) PutLocationsID() operations.PutLocationsIDHandler {
	return operations.PutLocationsIDHandlerFunc(func(params operations.PutLocationsIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateLocation(params.HTTPRequest.Context(), params.Location)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutLocationsIDNoContent()
	})
}

func (h *handlers) DeleteLocationsID() operations.DeleteLocationsIDHandler {
	return operations.DeleteLocationsIDHandlerFunc(func(params operations.DeleteLocationsIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveLocation(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteLocationsIDNoContent()
	})
}

func (h *handlers) GetOrganizations() operations.GetOrganizationsHandler {
	return operations.GetOrganizationsHandlerFunc(func(params operations.GetOrganizationsParams, principal *string) middleware.Responder {
		u, err := h.service.Organizations(params.HTTPRequest.Context())

		if err != nil {
			return operations.NewGetOrganizationsInternalServerError().WithPayload(&models.Error{
				Code:    utils.ErrServerError,
				Message: err.Error(),
			})
		}

		return operations.NewGetOrganizationsOK().WithPayload(u)
	})
}

func (h *handlers) GetOrganizationsID() operations.GetOrganizationsIDHandler {
	return operations.GetOrganizationsIDHandlerFunc(func(params operations.GetOrganizationsIDParams, principal *string) middleware.Responder {
		u, err := h.service.Organization(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetOrganizationsIDOK().WithPayload(u)
	})
}

func (h *handlers) GetOrganizationsIDLocations() operations.GetOrganizationsIDLocationsHandler {
	return operations.GetOrganizationsIDLocationsHandlerFunc(func(params operations.GetOrganizationsIDLocationsParams, principal *string) middleware.Responder {
		u, err := h.service.OrganizationLocationIDs(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetOrganizationsIDLocationsOK().WithPayload(u)
	})
}

func (h *handlers) GetOrganizationsIDUsers() operations.GetOrganizationsIDUsersHandler {
	return operations.GetOrganizationsIDUsersHandlerFunc(func(params operations.GetOrganizationsIDUsersParams, principal *string) middleware.Responder {
		u, err := h.service.DomainUserIDs(params.HTTPRequest.Context(), &authCommon.DomainTypeOrganization, &params.ID, params.RoleID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetOrganizationsIDUsersOK().WithPayload(u)
	})
}

func (h *handlers) PostOrganizations() operations.PostOrganizationsHandler {
	return operations.PostOrganizationsHandlerFunc(func(params operations.PostOrganizationsParams, principal *string) middleware.Responder {
		u, err := h.service.AddOrganization(params.HTTPRequest.Context(), params.Organization)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostOrganizationsCreated().WithPayload(u)
	})
}

func (h *handlers) PutOrganizationsID() operations.PutOrganizationsIDHandler {
	return operations.PutOrganizationsIDHandlerFunc(func(params operations.PutOrganizationsIDParams, principal *string) middleware.Responder {
		_, err := h.service.UpdateOrganization(params.HTTPRequest.Context(), params.Organization)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPutOrganizationsIDNoContent()
	})
}

func (h *handlers) DeleteOrganizationsID() operations.DeleteOrganizationsIDHandler {
	return operations.DeleteOrganizationsIDHandlerFunc(func(params operations.DeleteOrganizationsIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveOrganization(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteOrganizationsIDNoContent()
	})
}

func (h *handlers) GetUserRoles() operations.GetUserRolesHandler {
	return operations.GetUserRolesHandlerFunc(func(params operations.GetUserRolesParams, principal *string) middleware.Responder {
		r, err := h.service.FindUserRoles(params.HTTPRequest.Context(), params.UserID, params.RoleID, params.DomainType, params.DomainID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUserRolesOK().WithPayload(r)
	})
}

func (h *handlers) GetUserRolesID() operations.GetUserRolesIDHandler {
	return operations.GetUserRolesIDHandlerFunc(func(params operations.GetUserRolesIDParams, principal *string) middleware.Responder {
		r, err := h.service.UserRole(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewGetUserRolesIDOK().WithPayload(r)
	})
}

func (h *handlers) PostUserRoles() operations.PostUserRolesHandler {
	return operations.PostUserRolesHandlerFunc(func(params operations.PostUserRolesParams, principal *string) middleware.Responder {
		r, err := h.service.AddUserRole(params.HTTPRequest.Context(), params.UserRole)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewPostUserRolesCreated().WithPayload(r)
	})
}

func (h *handlers) DeleteUserRolesID() operations.DeleteUserRolesIDHandler {
	return operations.DeleteUserRolesIDHandlerFunc(func(params operations.DeleteUserRolesIDParams, principal *string) middleware.Responder {
		err := h.service.RemoveUserRole(params.HTTPRequest.Context(), params.ID)

		if err != nil {
			return utils.NewErrorResponse(err)
		}

		return operations.NewDeleteUserRolesIDNoContent()
	})
}

func (h *handlers) GetDatabase() operations.GetDatabaseHandler {
	return operations.GetDatabaseHandlerFunc(func(params operations.GetDatabaseParams, principal *string) middleware.Responder {
		etag := strings.Trim(params.HTTPRequest.Header.Get("Etag"), `"`)

		checksum, err := h.service.DBChecksum()
		if err != nil {
			return utils.UseProducer(utils.NewServerError(err), utils.JSONProducer)
		}

		currentEtag := base64.RawURLEncoding.EncodeToString(checksum)

		if etag == currentEtag {
			return operations.NewGetDatabaseNotModified()
		}

		reader, writer := io.Pipe()

		go func() {
			_, err := h.service.WriteDBTo(writer)
			writer.CloseWithError(err)
		}()

		return utils.UseProducer(
			operations.NewGetDatabaseOK().
				WithPayload(reader).
				WithEtag(`"`+currentEtag+`"`),
			utils.BinProducer)
	})
}

// NewHandlers returns a new instance of authDataManager handlers
func NewHandlers(service Service) Handlers {
	return &handlers{service: service}
}
