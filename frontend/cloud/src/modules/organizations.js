import _ from "lodash"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_ORGANIZATION = "organization/LOAD_ORGANIZATION"
const LOAD_ORGANIZATION_SUCCESS = "organization/LOAD_ORGANIZATION_SUCCESS"
const LOAD_ORGANIZATION_FAIL = "organization/LOAD_ORGANIZATION_FAIL"

const LOAD_ORGANIZATIONS = "organization/LOAD_ORGANIZATIONS"
const LOAD_ORGANIZATIONS_SUCCESS = "organization/LOAD_ORGANIZATIONS_SUCCESS"
const LOAD_ORGANIZATIONS_FAIL = "organization/LOAD_ORGANIZATIONS_FAIL"

const DELETE_ORGANIZATION_FAIL = "organization/DELETE_ORGANIZATION_FAIL"
const DELETE_ORGANIZATION_SUCCESS = "organization/DELETE_ORGANIZATION_SUCCESS"

const SAVE_ORGANIZATION_FAIL = "organization/SAVE_ORGANIZATION_FAIL"
const SAVE_ORGANIZATION_SUCCESS = "organization/SAVE_ORGANIZATION_SUCCESS"

const initialState = {
    loading: true
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

        case LOAD_ORGANIZATIONS:
            return {
                ...state,
                loading: true
            }
        case LOAD_ORGANIZATIONS_SUCCESS:
            return {
                ...state,
                organizations: _.keyBy(action.organizations, "id"),
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
        default:
            return state
    }
}

export const loadOrganization = organizationID => {
    return dispatch => {
        dispatch({
            type: LOAD_ORGANIZATION
        })

        return api(`/auth/organizations/${organizationID}`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_ORGANIZATION_SUCCESS,
                    organization: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ORGANIZATION_FAIL
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

        return api(`/auth/organizations`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_ORGANIZATIONS_SUCCESS,
                    organizations: response
                })
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

        return api(`/auth/organizations/${organizationID}`, "DELETE")
            .then(response => {
                dispatch({
                    type: DELETE_ORGANIZATION_SUCCESS,
                    organizationID: organizationID
                })
                dispatch(open("Deleted organization", "", COLOR_SUCCESS, 5))
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

        return api(url, method, organization)
            .then(response => {
                if (organization.id) {
                    response = organization
                }
                dispatch({
                    type: SAVE_ORGANIZATION_SUCCESS,
                    organization: response
                })
                dispatch(open("Saved organization", "", COLOR_SUCCESS, 5))
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
