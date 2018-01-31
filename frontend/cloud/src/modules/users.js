import produce from "immer"

import api from "./api"
import { open, COLOR_DANGER } from "./alert"

const LOAD_USER = "user/LOAD_USER"
const LOAD_USER_SUCCESS = "user/LOAD_USER_SUCCESS"
const LOAD_USER_FAIL = "user/LOAD_USER_FAIL"

const LOAD_USERS = "user/LOAD_USERS"
const LOAD_USERS_SUCCESS = "user/LOAD_USERS_SUCCESS"
const LOAD_USERS_FAIL = "user/LOAD_USERS_FAIL"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOAD_USER:
                draft.loading = true
                break
            case LOAD_USER_SUCCESS:
                draft.loading = false
                draft.user = action.user
                break
            case LOAD_USER_FAIL:
                draft.loading = false
                break

            case LOAD_USERS:
                draft.loading = true
                break
            case LOAD_USERS_SUCCESS:
                draft.loading = false
                draft.users = action.users
                break
            case LOAD_USERS_FAIL:
                draft.loading = false
                break
            default:
        }
    })
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
                    type: LOAD_USERS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
