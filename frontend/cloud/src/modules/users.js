import _ from "lodash"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_USER = "user/LOAD_USER"
const LOAD_USER_SUCCESS = "user/LOAD_USER_SUCCESS"
const LOAD_USER_FAIL = "user/LOAD_USER_FAIL"

const LOAD_USER_LOCATION_IDS = "user/LOAD_USER_LOCATION_IDS"
const LOAD_USER_LOCATION_IDS_SUCCESS = "user/LOAD_USER_LOCATION_IDS_SUCCESS"
const LOAD_USER_LOCATION_IDS_FAIL = "user/LOAD_USER_LOCATION_IDS_FAIL"

const LOAD_USER_ORGANIZATION_IDS = "user/LOAD_USER_ORGANIZATION_IDS"
const LOAD_USER_ORGANIZATION_IDS_SUCCESS = "user/LOAD_USER_ORGANIZATION_IDS_SUCCESS"
const LOAD_USER_ORGANIZATION_IDS_FAIL = "user/LOAD_USER_ORGANIZATION_IDS_FAIL"

const LOAD_USER_CLINIC_IDS = "user/LOAD_USER_CLINIC_IDS"
const LOAD_USER_CLINIC_IDS_SUCCESS = "user/LOAD_USER_CLINIC_IDS_SUCCESS"
const LOAD_USER_CLINIC_IDS_FAIL = "user/LOAD_USER_CLINIC_IDS_FAIL"

const LOAD_USER_ROLE_IDS = "user/LOAD_USER_ROLE_IDS"
const LOAD_USER_ROLE_IDS_SUCCESS = "user/LOAD_USER_ROLE_IDS_SUCCESS"
const LOAD_USER_ROLE_IDS_FAIL = "user/LOAD_USER_ROLE_IDS_FAIL"

const LOAD_USERS = "user/LOAD_USERS"
const LOAD_USERS_SUCCESS = "user/LOAD_USERS_SUCCESS"
const LOAD_USERS_FAIL = "user/LOAD_USERS_FAIL"

const DELETE_USER_FAIL = "user/DELETE_USER_FAIL"
const DELETE_USER_SUCCESS = "user/DELETE_USER_SUCCESS"

const SAVE_USER_FAIL = "user/SAVE_USER_FAIL"
const SAVE_USER_SUCCESS = "user/SAVE_USER_SUCCESS"

const ROLE_ID_ALL = "all"
const DOMAIN_TYPE_ALL = "all"
const DOMAIN_ID_ALL = "all"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_USER:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_SUCCESS:
            return {
                ...state,
                users: _.assign({}, state.users || {}, _.fromPairs([[action.user.id, action.user]])),
                loading: false
            }
        case LOAD_USER_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USER_LOCATION_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_LOCATION_IDS_SUCCESS:
            return {
                ...state,
                usersLocationIDs: _.assign({}, state.usersLocationIDs || {}, _fromPairs([[action.userID, _.fromPairs([[action.roleID, action.locationIDs]])]])),
                loading: false
            }
        case LOAD_USER_LOCATION_IDS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USER_ORGANIZATION_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_ORGANIZATION_IDS_SUCCESS:
            return {
                ...state,
                usersOrganizationIDs: _.assign({}, state.usersOrganizationIDs || {}, _fromPairs([[action.userID, _.fromPairs([[action.roleID, action.organizationIDs]])]])),
                loading: false
            }
        case LOAD_USER_ORGANIZATION_IDS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USER_CLINIC_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_CLINIC_IDS_SUCCESS:
            return {
                ...state,
                usersClinicIDs: _.assign({}, state.usersClinicIDs || {}, _fromPairs([[action.userID, _.fromPairs([[action.roleID, action.clinicIDs]])]])),
                loading: false
            }
        case LOAD_USER_CLINICS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USER_ROLE_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_ROLE_IDS_SUCCESS:
            return {
                ...state,
                usersRoleIDs: _.assign({}, state.usersRoleIDs || {}, _fromPairs([[action.userID, _.fromPairs([[action.domainType, _.fromPairs([[action.domainID, action.roleIDs]])]])]])),
                loading: false
            }
        case LOAD_USER_ROLES_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USERS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USERS_SUCCESS:
            return {
                ...state,
                users: _.keyBy(action.users, "id"),
                loading: false
            }
        case LOAD_USERS_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_USER_SUCCESS:
            return {
                ...state,
                users: _.pickBy(state.users, user => user.id !== action.userID)
            }

        case SAVE_USER_SUCCESS:
            return {
                ...state,
                users: _.assign({}, state.users, _.fromPairs([[action.user.id, action.user]]))
            }
        default:
            return state
    }
}

export const loadUser = userID => {
    return dispatch => {
        dispatch({
            type: LOAD_USER
        })

        return api(`/auth/users/${userID}`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USER_SUCCESS,
                    user: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserOrganizationIDs = (userID, roleID) => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_ORGANIZATION_IDS
        })

        var url = `/auth/users/${userID}/organizations`
        if (roleID && roleID !== ROLE_ID_ALL) {
            url += `?roleID=${roleID}`
        }

        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USER_ORGANIZATION_IDS_SUCCESS,
                    userID: userID,
                    roleID: roleID ? roleID : ROLE_ID_ALL,
                    organizationIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_ORGANIZATION_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserLocationIDs = (userID, roleID) => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_LOCATION_IDS
        })

        var url = `/auth/users/${userID}/locations`
        if (roleID && roleID !== ROLE_ID_ALL) {
            url += `?roleID=${roleID}`
        }

        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USER_LOCATION_IDS_SUCCESS,
                    userID: userID,
                    roleID: roleID ? roleID : ROLE_ID_ALL,
                    locationIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_LOCATION_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserClinicIDs = (userID, roleID) => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_CLINIC_IDS
        })

        var url = `/auth/users/${userID}/clinics`
        if (roleID && roleID !== ROLE_ID_ALL) {
            url += `?roleID=${roleID}`
        }

        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USER_CLINIC_IDS_SUCCESS,
                    userID: userID,
                    roleID: roleID ? roleID : ROLE_ID_ALL,
                    clinicIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_CLINIC_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserRoleIDs = (userID, domainType, domainID) => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_ROLE_IDS
        })

        var url = `/auth/users/${userID}/roles`
        if (domainType && domainType !== DOMAIN_TYPE_ALL) {
            url += `?domainType=${domainType}`
            if (domainID && domainID !== DOMAIN_ID_ALL) {
                url += `&domainID=${domainID}`
            }
        }

        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USER_ROLE_IDS_SUCCESS,
                    userID: userID,
                    domainType: domainType ? domainType : DOMAIN_TYPE_ALL,
                    domainID: domainType ? (domainID ? domainID : DOMAIN_ID_ALL) : DOMAIN_ID_ALL,
                    roleIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_ROLE_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUsers = () => {
    return dispatch => {
        dispatch({
            type: LOAD_USERS
        })

        return api(`/auth/users`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USERS_SUCCESS,
                    users: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USERS_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUser = userID => {
    return dispatch => {
        dispatch(close())

        return api(`/auth/users/${userID}`, "DELETE")
            .then(response => {
                dispatch({
                    type: DELETE_USER_SUCCESS,
                    userID: userID
                })
                dispatch(open("Deleted user", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: DELETE_USER_FAIL,
                    userID: userID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveUser = user => {
    return dispatch => {
        dispatch(close())
        let url = "/auth/users"
        let method = "POST"
        if (user.id) {
            url += "/" + user.id
            method = "PUT"
        }

        return api(url, method, user)
            .then(response => {
                if (user.id) {
                    response = user
                }
                dispatch({
                    type: SAVE_USER_SUCCESS,
                    user: response
                })
                dispatch(open("Saved user", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: SAVE_USER_FAIL,
                    user: user
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
