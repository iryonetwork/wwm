import _ from "lodash"
import store from "../store"

import api from "shared/modules/api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_ALL_USERROLES = "userRole/LOAD_ALL_USERROLES"
const LOAD_ALL_USERROLES_SUCCESS = "userRole/LOAD_ALL_USERROLES_SUCCESS"
const LOAD_ALL_USERROLES_FAIL = "userRole/LOAD_ALL_USERROLES_FAIL"

const LOAD_USER_USERROLES = "userRole/LOAD_USER_USERROLES"
const LOAD_USER_USERROLES_SUCCESS = "userRole/LOAD_USER_USERROLES_SUCCESS"
const LOAD_USER_USERROLES_FAIL = "userRole/LOAD_USER_USERROLES_FAIL"

const LOAD_ROLE_USERROLES = "userRole/LOAD_ROLE_USERROLES"
const LOAD_ROLE_USERROLES_SUCCESS = "userRole/LOAD_ROLE_USERROLES_SUCCESS"
const LOAD_ROLE_USERROLES_FAIL = "userRole/LOAD_ROLE_USERROLES_FAIL"

const LOAD_DOMAIN_USERROLES = "userRole/LOAD_DOMAIN_USERROLES"
const LOAD_DOMAIN_USERROLES_SUCCESS = "userRole/LOAD_DOMAIN_USERROLES_SUCCESS"
const LOAD_DOMAIN_USERROLES_FAIL = "userRole/LOAD_DOMAIN_USERROLES_FAIL"

const DELETE_USERROLE_FAIL = "userRole/DELETE_USERROLE_FAIL"
const DELETE_USERROLE_SUCCESS = "userRole/DELETE_USERROLE_SUCCESS"

const SAVE_USERROLE_FAIL = "userRole/SAVE_USERROLE_FAIL"
const SAVE_USERROLE_SUCCESS = "userRole/SAVE_USERROLE_SUCCESS"

const CLEAR_USERROLES_STATE = "userRole/CLEAR_USERROLES_STATE"

const initialState = {
    loading: false,
    allLoaded: false,
    userRoles: undefined,
    userUserRoles: undefined,
    roleUserRoles: undefined,
    domainUserRoles: undefined,
    forbidden: false
}

export default (state = initialState, action) => {
    let forbidden = false
    let userRoles = {}
    let userUserRoles = {}
    let roleUserRoles = {}
    let domainUserRoles = {}

    switch (action.type) {
        case LOAD_ALL_USERROLES:
            return {
                ...state,
                loading: true
            }
        case LOAD_ALL_USERROLES_SUCCESS:
            userRoles = _.keyBy(action.userRoles, "id")
            userUserRoles = _.mapValues(_.groupBy(action.userRoles, "userID"), function(o) {
                return _.keyBy(o, "id")
            })
            roleUserRoles = _.mapValues(_.groupBy(action.userRoles, "roleID"), function(o) {
                return _.keyBy(o, "id")
            })
            domainUserRoles = _.groupBy(action.userRoles, "domainType")
            _.forEach(domainUserRoles, function(domainTypeUserRoles, domainType) {
                domainUserRoles[domainType] = _.mapValues(_.groupBy(domainTypeUserRoles, "domainID"), function(o) {
                    return _.keyBy(o, "id")
                })
            })

            return {
                ...state,
                userRoles: userRoles,
                userUserRoles: userUserRoles,
                roleUserRoles: roleUserRoles,
                domainUserRoles: domainUserRoles,
                allLoaded: true,
                loading: false
            }
        case LOAD_ALL_USERROLES_FAIL:
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case LOAD_USER_USERROLES:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_USERROLES_SUCCESS:
            return {
                ...state,
                userUserRoles: _.assign({}, state.userUserRoles, _.fromPairs([[action.userID, _.keyBy(action.userRoles, "id")]])),
                loading: false
            }
        case LOAD_USER_USERROLES_FAIL:
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case LOAD_DOMAIN_USERROLES:
            return {
                ...state,
                loading: true
            }
        case LOAD_DOMAIN_USERROLES_SUCCESS:
            domainUserRoles = state.domainUserRoles ? state.domainUserRoles : {}
            domainUserRoles[action.domainType] = domainUserRoles[action.domainType] ? domainUserRoles[action.domainType] : {}
            domainUserRoles[action.domainType][action.domainID] = _.keyBy(action.userRoles, "id")
            return {
                ...state,
                domainUserRoles: domainUserRoles,
                loading: false
            }
        case LOAD_DOMAIN_USERROLES_FAIL:
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case LOAD_ROLE_USERROLES:
            return {
                ...state,
                loading: true
            }
        case LOAD_ROLE_USERROLES_SUCCESS:
            return {
                ...state,
                roleUserRoles: _.assign({}, state.roleUserRoles, _.fromPairs([[action.roleID, _.keyBy(action.userRoles, "id")]])),
                loading: false
            }
        case LOAD_ROLE_USERROLES_FAIL:
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_USERROLE_SUCCESS:
            userRoles = _.pickBy(state.userRoles, userRole => userRole.id !== action.userRoleID)
            userUserRoles = {}
            _.forEach(state.userUserRoles, function(userRoles, userID) {
                userUserRoles[userID] = _.pickBy(userRoles, userRole => userRole.id !== action.userRoleID)
            })
            roleUserRoles = {}
            _.forEach(state.roleUserRoles, function(userRoles, roleID) {
                roleUserRoles[roleID] = _.pickBy(userRoles, userRole => userRole.id !== action.userRoleID)
            })
            domainUserRoles = {}
            _.forEach(state.domainUserRoles, function(domainTypeUserRoles, domainType) {
                domainUserRoles[domainType] = {}
                _.forEach(domainTypeUserRoles, function(userRoles, domainID) {
                    domainUserRoles[domainType][domainID] = _.pickBy(userRoles, userRole => userRole.id !== action.userRoleID)
                })
            })

            return {
                ...state,
                userRoles: userRoles,
                userUserRoles: userUserRoles,
                roleUserRoles: roleUserRoles,
                domainUserRoles: domainUserRoles
            }

        case SAVE_USERROLE_SUCCESS:
            let id = action.userRole.id
            let userID = action.userRole.userID
            let roleID = action.userRole.roleID
            let domainType = action.userRole.domainType
            let domainID = action.userRole.domainID

            userRoles = state.userRoles || {}
            if (userRoles) {
                userRoles = _.assign({}, state.userRoles, _.fromPairs([[id, action.userRole]]))
            }

            userUserRoles = state.userUserRoles || {}
            if (userUserRoles[userID]) {
                userUserRoles[userID] = _.assign({}, userUserRoles[userID] || {}, _.fromPairs([[id, action.userRole]]))
            }

            roleUserRoles = state.roleUserRoles || {}
            if (roleUserRoles[roleID]) {
                roleUserRoles[roleID] = _.assign({}, roleUserRoles[roleID] || {}, _.fromPairs([[id, action.userRole]]))
            }

            domainUserRoles = state.domainUserRoles || {}
            if (domainUserRoles[domainType]) {
                if (domainUserRoles[domainType][domainID]) {
                    domainUserRoles[domainType][domainID] = _.assign({}, domainUserRoles[domainType][domainID], _.fromPairs([[id, action.userRole]]))
                }
            }

            return {
                ...state,
                userRoles: userRoles,
                userUserRoles: userUserRoles,
                roleUserRoles: roleUserRoles,
                domainUserRoles: domainUserRoles
            }

        case CLEAR_USERROLES_STATE:
            return {
                loading: false,
                allLoaded: false,
                userRoles: undefined,
                userUserRoles: undefined,
                roleUserRoles: undefined,
                domainUserRoles: undefined
            }

        default:
            return state
    }
}

export const loadAllUserRoles = () => {
    return dispatch => {
        if (store.getState().userRoles.loading) {
            return Promise.resolve()
        }

        dispatch({
            type: LOAD_ALL_USERROLES
        })

        return dispatch(api(`/auth/userRoles`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ALL_USERROLES_SUCCESS,
                    userRoles: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ALL_USERROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserUserRoles = userID => {
    return dispatch => {
        if (store.getState().userRoles.loading) {
            return Promise.resolve()
        }

        dispatch({
            type: LOAD_USER_USERROLES
        })

        let url = "/auth/userRoles?userID=" + userID

        return dispatch(api(url, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_USER_USERROLES_SUCCESS,
                    userRoles: response,
                    userID: userID
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_USERROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadRoleUserRoles = roleID => {
    return dispatch => {
        dispatch({
            type: LOAD_ROLE_USERROLES
        })

        let url = "/auth/userRoles?roleID=" + roleID

        return dispatch(api(url, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ROLE_USERROLES_SUCCESS,
                    userRoles: response,
                    roleID: roleID
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ROLE_USERROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadDomainUserRoles = (domainType, domainID) => {
    return dispatch => {
        dispatch({
            type: LOAD_DOMAIN_USERROLES
        })

        let url = "/auth/userRoles?domainType=" + domainType + "&domainID=" + domainID

        return dispatch(api(url, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_DOMAIN_USERROLES_SUCCESS,
                    userRoles: response,
                    domainType: domainType,
                    domainID: domainID
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_DOMAIN_USERROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUserRole = userRoleID => {
    return dispatch => {
        dispatch(close())

        return dispatch(api(`/auth/userRoles/${userRoleID}`, "DELETE"))
            .then(response => {
                dispatch({
                    type: DELETE_USERROLE_SUCCESS,
                    userRoleID: userRoleID
                })
                dispatch(open("Deleted user role", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: DELETE_USERROLE_FAIL,
                    userRoleID: userRoleID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUserRoleNoAlert = userRoleID => {
    return dispatch => {
        dispatch(close())

        return dispatch(api(`/auth/userRoles/${userRoleID}`, "DELETE"))
            .then(response => {
                dispatch({
                    type: DELETE_USERROLE_SUCCESS,
                    userRoleID: userRoleID
                })
            })
            .catch(error => {
                dispatch({
                    type: DELETE_USERROLE_FAIL,
                    userRoleID: userRoleID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveUserRole = userRole => {
    return dispatch => {
        dispatch(close())

        let url = "/auth/userRoles"
        let method = "POST"
        if (userRole.id) {
            url += "/" + userRole.id
            method = "PUT"
        }

        return dispatch(api(url, method, userRole))
            .then(response => {
                if (userRole.id) {
                    response = userRole
                }
                dispatch({
                    type: SAVE_USERROLE_SUCCESS,
                    userRole: response
                })
                dispatch(open("Saved user role", "", COLOR_SUCCESS, 5))

                return response
            })
            .catch(error => {
                dispatch({
                    type: SAVE_USERROLE_FAIL,
                    userRole: userRole
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveUserRoleCustomMsg = (userRole, msg) => {
    return dispatch => {
        dispatch(close())

        let url = "/auth/userRoles"
        let method = "POST"
        if (userRole.id) {
            url += "/" + userRole.id
            method = "PUT"
        }

        return dispatch(api(url, method, userRole))
            .then(response => {
                if (userRole.id) {
                    response = userRole
                }
                dispatch({
                    type: SAVE_USERROLE_SUCCESS,
                    userRole: response
                })
                dispatch(open(msg, "", COLOR_SUCCESS, 5))

                return response
            })
            .catch(error => {
                dispatch({
                    type: SAVE_USERROLE_FAIL,
                    userRole: userRole
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const clearUserRoles = () => {
    return dispatch => {
        dispatch({
            type: CLEAR_USERROLES_STATE
        })

        return Promise.resolve()
    }
}
