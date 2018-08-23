import _ from "lodash"

import api from "shared/modules/api"
import { clearUserRoles } from "./userRoles"
import { open, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_USER = "user/LOAD_USER"
const LOAD_USER_SUCCESS = "user/LOAD_USER_SUCCESS"
const LOAD_USER_FAIL = "user/LOAD_USER_FAIL"

const LOAD_USERS = "user/LOAD_USERS"
const LOAD_USERS_SUCCESS = "user/LOAD_USERS_SUCCESS"
const LOAD_USERS_FAIL = "user/LOAD_USERS_FAIL"

const DELETE_USER = "user/DELETE_USER"
const DELETE_USER_FAIL = "user/DELETE_USER_FAIL"
const DELETE_USER_SUCCESS = "user/DELETE_USER_SUCCESS"

const SAVE_USER = "user/SAVE_USER"
const SAVE_USER_FAIL = "user/SAVE_USER_FAIL"
const SAVE_USER_SUCCESS = "user/SAVE_USER_SUCCESS"

const initialState = {
    loading: false,
    updating: false,
    allLoaded: false,
    forbidden: false
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

        case LOAD_USERS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USERS_SUCCESS:
            return {
                ...state,
                users: _.keyBy(action.users, "id"),
                allLoaded: true,
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

        case DELETE_USER:
            return {
                ...state,
                updating: true
            }

        case DELETE_USER_SUCCESS:
            return {
                ...state,
                users: _.pickBy(state.users, user => user.id !== action.userID)
            }

        case SAVE_USER:
            return {
                ...state,
                updating: true
            }

        case SAVE_USER_SUCCESS:
            return {
                ...state,
                updating: false,
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

        return dispatch(api(`/auth/users/${userID}`, "GET"))
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

export const loadUsers = () => {
    return dispatch => {
        dispatch({
            type: LOAD_USERS
        })

        return dispatch(api(`/auth/users`, "GET"))
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
        dispatch({
            type: DELETE_USER
        })

        return dispatch(api(`/auth/users/${userID}`, "DELETE"))
            .then(response => {
                dispatch(clearUserRoles())
                dispatch({
                    type: DELETE_USER_SUCCESS,
                    userID: userID
                })
                setTimeout(() => dispatch(open("Deleted User", "", COLOR_SUCCESS, 5)), 100)
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
        dispatch({
            type: SAVE_USER
        })

        let url = "/auth/users"
        let method = "POST"
        if (user.id) {
            url += "/" + user.id
            method = "PUT"
        }

        return dispatch(api(url, method, user))
            .then(response => {
                if (user.id) {
                    response = user
                }
                dispatch({
                    type: SAVE_USER_SUCCESS,
                    user: response
                })
                setTimeout(() => dispatch(open("Saved User", "", COLOR_SUCCESS, 5)), 100)

                return response
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
