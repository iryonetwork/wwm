import _ from "lodash"

import api from "shared/modules/api"
import { open, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_ROLES = "roles/LOAD_ROLES"
const LOAD_ROLES_SUCCESS = "roles/LOAD_ROLES_SUCCESS"
const LOAD_ROLES_FAIL = "roles/LOAD_ROLES_FAIL"

const UPDATE_ROLE_SUCCESS = "roles/UPDATE_ROLES_SUCCESS"
const CREATE_ROLE_SUCCESS = "roles/CREATE_ROLE_SUCCESS"
const DELETE_ROLE_SUCCESS = "roles/DELETE_ROLE_SUCCESS"

const initialState = {
    loading: false,
    allLoaded: false,
    forbidden: false
}

export default (state = initialState, action) => {
    let roles
    let forbidden = false

    switch (action.type) {
        case LOAD_ROLES:
            return {
                ...state,
                loading: true
            }

        case LOAD_ROLES_SUCCESS:
            return {
                ...state,
                loading: false,
                allLoaded: true,
                roles: _.keyBy(action.roles, "id")
            }

        case LOAD_ROLES_FAIL:
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case UPDATE_ROLE_SUCCESS:
            roles = { ...state.roles }
            roles[action.role.id] = action.role

            return {
                loading: false,
                roles
            }

        case CREATE_ROLE_SUCCESS:
            roles = { ...state.roles }
            roles[action.role.id] = action.role
            return {
                ...state,
                roles
            }

        case DELETE_ROLE_SUCCESS:
            roles = { ...state.roles }
            delete roles[action.roleID]
            return {
                ...state,
                roles
            }

        default:
            return state
    }
}

export const loadRoles = () => {
    return dispatch => {
        dispatch({
            type: LOAD_ROLES
        })

        return dispatch(api(`/auth/roles`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_ROLES_SUCCESS,
                    roles: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ROLES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const addRole = name => {
    return dispatch => {
        let role = {
            name
        }

        return dispatch(api(`/auth/roles`, "POST", role))
            .then(response => {
                dispatch({
                    type: CREATE_ROLE_SUCCESS,
                    role: response
                })
                dispatch(open(`Created role ${name}`, "", COLOR_SUCCESS, 5))

                return response
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteRole = roleID => {
    return dispatch => {
        return dispatch(api(`/auth/roles/${roleID}`, "DELETE"))
            .then(response => {
                dispatch({
                    type: DELETE_ROLE_SUCCESS,
                    roleID: roleID
                })
                dispatch(open(`Deleted role`, "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
