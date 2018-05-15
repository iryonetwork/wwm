import _ from "lodash"
import store from "../store"

import api from "shared/modules/api"
import { loadDomainUserRoles, deleteUserRoleNoAlert, clearUserRoles } from "./userRoles"
import { deleteUserFromClinic, clearClinics } from "./clinics"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_ORGANIZATION = "organization/LOAD_ORGANIZATION"
const LOAD_ORGANIZATION_SUCCESS = "organization/LOAD_ORGANIZATION_SUCCESS"
const LOAD_ORGANIZATION_FAIL = "organization/LOAD_ORGANIZATION_FAIL"

const LOAD_ORGANIZATION_LOCATION_IDS = "organization/LOAD_ORGANIZATION_LOCATION_IDS"
const LOAD_ORGANIZATION_LOCATION_IDS_SUCCESS = "organization/LOAD_ORGANIZATION_LOCATION_IDS_SUCCESS"
const LOAD_ORGANIZATION_LOCATION_IDS_FAIL = "organization/LOAD_ORGANIZATION_LOCATION_IDS_FAIL"

const LOAD_ORGANIZATIONS = "organization/LOAD_ORGANIZATIONS"
const LOAD_ORGANIZATIONS_SUCCESS = "organization/LOAD_ORGANIZATIONS_SUCCESS"
const LOAD_ORGANIZATIONS_FAIL = "organization/LOAD_ORGANIZATIONS_FAIL"

const DELETE_ORGANIZATION_FAIL = "organization/DELETE_ORGANIZATION_FAIL"
const DELETE_ORGANIZATION_SUCCESS = "organization/DELETE_ORGANIZATION_SUCCESS"

const SAVE_ORGANIZATION_FAIL = "organization/SAVE_ORGANIZATION_FAIL"
const SAVE_ORGANIZATION_SUCCESS = "organization/SAVE_ORGANIZATION_SUCCESS"

const CLEAR_ORGANIZATIONS_STATE = "clinic/CLEAR_ORGANIZATIONS_STATE"

const initialState = {
    loading: false,
    allLoaded: false,
    forbidden: false
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_ORGANIZATION:
            return {
                ...state,
                loading: true
            }
        case LOAD_ORGANIZATION_SUCCESS:
            return {
                ...state,
                organizations: _.assign({}, state.organizations || {}, _.fromPairs([[action.organization.id, action.organization]])),
                loading: false
            }
        case LOAD_ORGANIZATION_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_ORGANIZATION_LOCATION_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_ORGANIZATION_LOCATION_IDS_SUCCESS:
            return {
                ...state,
                organizationsLocationIDs: _.assign({}, state.organizationsLocationIDs || {}, _.fromPairs([[action.organizationID, action.locationIDs]])),
                loading: false
            }
        case LOAD_ORGANIZATION_LOCATION_IDS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_ORGANIZATIONS:
            return {
                ...state,
                loading: true
            }
        case LOAD_ORGANIZATIONS_SUCCESS:
            return {
                ...state,
                organizations: _.keyBy(action.organizations, "id"),
                allLoaded: true,
                loading: false
            }
        case LOAD_ORGANIZATIONS_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_ORGANIZATION_SUCCESS:
            return {
                ...state,
                organizations: _.pickBy(state.organizations, organization => organization.id !== action.organizationID)
            }

        case SAVE_ORGANIZATION_SUCCESS:
            return {
                ...state,
                organizations: _.assign({}, state.organizations, _.fromPairs([[action.organization.id, action.organization]]))
            }

        case CLEAR_ORGANIZATIONS_STATE:
            return {
                organizations: undefined,
                allLoaded: false,
                loading: false
            }

        default:
            return state
    }
}

export const loadOrganization = organizationID => {
    return dispatch => {
        dispatch({
            type: LOAD_ORGANIZATION
        })

        return dispatch(api(`/auth/organizations/${organizationID}`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ORGANIZATION_SUCCESS,
                    organization: response
                })

                return response
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ORGANIZATION_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadOrganizationLocationIDs = organizationID => {
    return dispatch => {
        dispatch({
            type: LOAD_ORGANIZATION_LOCATION_IDS
        })

        return dispatch(api(`/auth/organizations/${organizationID}/locations`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ORGANIZATION_LOCATION_IDS_SUCCESS,
                    organizationID: organizationID,
                    locationIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ORGANIZATION_LOCATION_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadOrganizations = () => {
    return dispatch => {
        dispatch({
            type: LOAD_ORGANIZATIONS
        })

        return dispatch(api(`/auth/organizations`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ORGANIZATIONS_SUCCESS,
                    organizations: response
                })

                return response
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ORGANIZATIONS_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteOrganization = organizationID => {
    return dispatch => {
        dispatch(close())

        return dispatch(api(`/auth/organizations/${organizationID}`, "DELETE"))
            .then(response => {
                dispatch(clearUserRoles())
                dispatch(clearClinics())
                dispatch({
                    type: DELETE_ORGANIZATION_SUCCESS,
                    organizationID: organizationID
                })
                setTimeout(() => dispatch(open("Deleted organization", "", COLOR_SUCCESS, 5)), 100)
            })
            .catch(error => {
                dispatch({
                    type: DELETE_ORGANIZATION_FAIL,
                    organizationID: organizationID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveOrganization = organization => {
    return dispatch => {
        dispatch(close())

        let url = "/auth/organizations"
        let method = "POST"
        if (organization.id) {
            url += "/" + organization.id
            method = "PUT"
        }

        return dispatch(api(url, method, organization))
            .then(response => {
                if (organization.id) {
                    response = organization
                }
                dispatch({
                    type: SAVE_ORGANIZATION_SUCCESS,
                    organization: response
                })
                setTimeout(() => dispatch(open("Saved organization", "", COLOR_SUCCESS, 5)), 100)

                return response
            })
            .catch(error => {
                dispatch({
                    type: SAVE_ORGANIZATION_FAIL,
                    organization: organization
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUserFromOrganization = (organizationID, userID) => {
    return dispatch => {
        dispatch(close())

        let organization =
            store.getState().organizations.organizations && store.getState().organizations.organizations[organizationID]
                ? store.getState().organizations.organizations[organizationID]
                : undefined
        if (!organization) {
            return dispatch(loadOrganization(organizationID)).then(() => {
                return dispatch(deleteUserFromOrganization(organizationID, userID))
            })
        }

        // check for user roles to delete in store
        let userRolesToDelete = undefined
        let organizationUserRoles =
            store.getState().userRoles.domainUserRoles &&
            store.getState().userRoles.domainUserRoles["organization"] &&
            store.getState().userRoles.domainUserRoles["organization"][organizationID]
                ? store.getState().userRoles.domainUserRoles["organization"][organizationID]
                : undefined
        if (organizationUserRoles === undefined) {
            let userUserRoles =
                store.getState().userRoles.userUserRoles && store.getState().userRoles.userUserRoles[userID]
                    ? store.getState().userRoles.userUserRoles[userID]
                    : undefined
            if (userUserRoles !== undefined) {
                userRolesToDelete = _.pickBy(userUserRoles, userRole => userRole.domainType === "organization" && userRole.domainID === organizationID) || {}
            }
        } else {
            userRolesToDelete = _.pickBy(organizationUserRoles, userRole => userRole.userID === userID) || {}
        }

        // no user roles to delete in store, fetch
        if (userRolesToDelete === undefined) {
            return dispatch(loadDomainUserRoles("organization", organizationID)).then(() => {
                return dispatch(deleteUserFromOrganization(organizationID, userID))
            })
        }

        _.forEach(userRolesToDelete, userRole => {
            dispatch(deleteUserRoleNoAlert(userRole.id))
        })

        // remove user also from clinic
        _.forEach(organization.clinics, clinicID => {
            dispatch(deleteUserFromClinic(clinicID, userID))
        })

        return Promise.resolve(true)
    }
}

export const clearOrganizations = () => {
    return dispatch => {
        dispatch({ type: CLEAR_ORGANIZATIONS_STATE })

        return Promise.resolve()
    }
}
