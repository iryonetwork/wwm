import produce from "immer"
import keyBy from "lodash/keyBy"
import forEach from "lodash/forEach"

import api from "./api"
import { open, COLOR_DANGER } from "./alert"

const LOAD_ROLES = "roles/LOAD_ROLES"
const LOAD_ROLES_SUCCESS = "roles/LOAD_ROLES_SUCCESS"
const LOAD_ROLES_FAIL = "roles/LOAD_ROLES_FAIL"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOAD_ROLES:
                draft.loading = true
                break
            case LOAD_ROLES_SUCCESS:
                draft.loading = false
                draft.roles = keyBy(action.roles, "id")
                draft.users = {}
                forEach(action.roles, role => {
                    forEach(role.users, user => {
                        if (!draft.users[user]) {
                            draft.users[user] = []
                        }
                        draft.users[user].push(role.id)
                    })
                })
                break
            case LOAD_ROLES_FAIL:
                draft.loading = false
                break
            default:
        }
    })
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
