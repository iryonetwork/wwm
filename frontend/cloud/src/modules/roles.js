import { push } from "react-router-redux"
import _ from "lodash"

import api from "./api"
import { open, COLOR_DANGER, COLOR_SUCCESS } from "./alert"

const LOAD_ROLES = "roles/LOAD_ROLES"
const LOAD_ROLES_SUCCESS = "roles/LOAD_ROLES_SUCCESS"
const LOAD_ROLES_FAIL = "roles/LOAD_ROLES_FAIL"

const UPDATE_ROLE_SUCCESS = "roles/UPDATE_ROLES_SUCCESS"
const CREATE_ROLE_SUCCESS = "roles/CREATE_ROLE_SUCCESS"
const DELETE_ROLE_SUCCESS = "roles/DELETE_ROLE_SUCCESS"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_ROLES:
            return {
                ...state,
                loading: true
            }

        case LOAD_ROLES_SUCCESS:
            let users = {}
            _.forEach(action.roles, role => {
                _.forEach(role.users, user => {
                    if (!users[user]) {
                        users[user] = []
                    }
                    users[user].push(role.id)
                })
            })
            return {
                ...state,
                loading: false,
                roles: _.keyBy(action.roles, "id"),
                users: users
            }

        case LOAD_ROLES_FAIL:
            return {
                ...state,
                loading: false
            }

        case UPDATE_ROLE_SUCCESS:
            let roles = { ...state.roles }
            users = { ...state.users }
            if (_.indexOf(roles[action.role.id].users, action.userID) === -1) {
                users[action.userID].push(action.role.id)
            } else {
                users[action.userID] = _.without(users[action.userID], action.role.id)
            }
            roles[action.role.id] = action.role

            return {
                loading: false,
                roles,
                users
            }

        case CREATE_ROLE_SUCCESS:
            roles = { ...state.roles }
            roles[action.role.id] = action.role
            return {
                ...state,
                roles
            }

        case DELETE_ROLE_SUCCESS:
            let checkUsers = state.roles[action.roleID].users
            roles = { ...state.roles }
            users = { ...state.users }
            delete roles[action.roleID]
            _.forEach(checkUsers, user => {
                users[user] = _.without(users[user], action.roleID)
            })
            return {
                ...state,
                roles,
                users
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

        return api(`/auth/roles`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_ROLES_SUCCESS,
                    roles: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_ROLES_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const removeUserFromRole = (roleID, userID) => {
    return (dispatch, getState) => {
        let role = getState().roles.roles[roleID]
        role.users = _.without(role.users, userID)

        return api(`/auth/roles/${roleID}`, "PUT", role)
            .then(response => {
                dispatch({
                    type: UPDATE_ROLE_SUCCESS,
                    role: role,
                    userID: userID
                })
                dispatch(open("User removed from role", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const addUserToRole = (roleID, userID) => {
    return (dispatch, getState) => {
        let role = getState().roles.roles[roleID]
        role.users.push(userID)

        return api(`/auth/roles/${roleID}`, "PUT", role)
            .then(response => {
                dispatch({
                    type: UPDATE_ROLE_SUCCESS,
                    role: role,
                    userID: userID
                })
                dispatch(open("Added user to role", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const addRole = name => {
    return dispatch => {
        let role = {
            users: [],
            name
        }

        return api(`/auth/roles`, "POST", role)
            .then(response => {
                dispatch({
                    type: CREATE_ROLE_SUCCESS,
                    role: response
                })
                dispatch(open(`Created role ${name}`, "", COLOR_SUCCESS, 5))
                dispatch(push(`/roles/${response.id}`))
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteRole = roleID => {
    return dispatch => {
        return api(`/auth/roles/${roleID}`, "DELETE")
            .then(response => {
                dispatch(push(`/roles/`))
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
