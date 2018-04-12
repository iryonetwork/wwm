import _ from "lodash"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_USERROLE = "userRole/LOAD_USERROLE"
const LOAD_USERROLE_SUCCESS = "userRole/LOAD_USERROLE_SUCCESS"
const LOAD_USERROLE_FAIL = "userRole/LOAD_USERROLE_FAIL"

const LOAD_USERROLES = "userRole/LOAD_USERROLES"
const LOAD_USERROLES_SUCCESS = "userRole/LOAD_USERROLES_SUCCESS"
const LOAD_USERROLES_FAIL = "userRole/LOAD_USERROLES_FAIL"

const DELETE_USERROLE_FAIL = "userRole/DELETE_USERROLE_FAIL"
const DELETE_USERROLE_SUCCESS = "userRole/DELETE_USERROLE_SUCCESS"

const SAVE_USERROLE_FAIL = "userRole/SAVE_USERROLE_FAIL"
const SAVE_USERROLE_SUCCESS = "userRole/SAVE_USERROLE_SUCCESS"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_USERROLE:
            return {
                ...state,
                loading: true
            }
        case LOAD_USERROLE_SUCCESS:
            return {
                ...state,
                userRoles: _.assign({}, state.userRoles || {}, _.fromPairs([[action.userRole.id, action.userRole]])),
                loading: false
            }
        case LOAD_USERROLE_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_USERROLES:
            return {
                ...state,
                loading: true
            }
        case LOAD_USERROLES_SUCCESS:
            return {
                ...state,
                userRoles: _.keyBy(action.userRoles, "id"),
                loading: false
            }
        case LOAD_USERROLES_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_USERROLE_SUCCESS:
            return {
                ...state,
                userRoles: _.pickBy(state.userRoles, userRole => userRole.id !== action.userRoleID)
            }

        case SAVE_USERROLE_SUCCESS:
            return {
                ...state,
                userRoles: _.assign({}, state.userRoles, _.fromPairs([[action.userRole.id, action.userRole]]))
            }
        default:
            return state
    }
}

export const loadUserRole = userRoleID => {
    return dispatch => {
        dispatch({
            type: LOAD_USERROLE
        })

        return api(`/auth/userRoles/${userRoleID}`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USERROLE_SUCCESS,
                    userRole: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USERROLE_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadUserRoles = () => {
    return dispatch => {
        dispatch({
            type: LOAD_USERROLES
        })

        return api(`/auth/userRoles`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_USERROLES_SUCCESS,
                    userRoles: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USERROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUserRole = userRoleID => {
    return dispatch => {
        dispatch(close())

        return api(`/auth/userRoles/${userRoleID}`, "DELETE")
            .then(response => {
                dispatch({
                    type: DELETE_USERROLE_SUCCESS,
                    userRoleID: userRoleID
                })
                dispatch(open("Deleted userRole", "", COLOR_SUCCESS, 5))
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

        return api(url, method, userRole)
            .then(response => {
                if (userRole.id) {
                    response = userRole
                }
                dispatch({
                    type: SAVE_USERROLE_SUCCESS,
                    userRole: response
                })
                dispatch(open("Saved userRole", "", COLOR_SUCCESS, 5))
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
