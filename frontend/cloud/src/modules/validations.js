import _ from "lodash"

import api from "shared/modules/api"
import { open, COLOR_DANGER } from "shared/modules/alert"

const LOAD_USER_RIGHTS = "user/LOAD_USER_RIGHTS"
const LOAD_USER_RIGHTS_SUCCESS = "user/LOAD_USER_RIGHTS_SUCCESS"
const LOAD_USER_RIGHTS_FAIL = "user/LOAD_USER_RIGHTS_FAIL"

export const SUPERADMIN_RIGHTS_RESOURCE = "/frontend/dashboard/superadmin"
export const ADMIN_RIGHTS_RESOURCE = "/frontend/dashboard/admin"
export const SELF_RIGHTS_RESOURCE = "/frontend/dashboard/self"

const initialState = {
    loading: false,
    forbidden: false
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_USER_RIGHTS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_RIGHTS_SUCCESS:
            return {
                ...state,
                userRights: _.reduce(
                    action.userRights,
                    function(result, value, key) {
                        result[value.query.resource] = value.result
                        return result
                    },
                    {}
                ),
                loading: false
            }
        case LOAD_USER_RIGHTS_FAIL:
            return {
                ...state,
                loading: false
            }
        default:
            return state
    }
}

export const loadUserRights = userID => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_RIGHTS
        })

        let validations = [
            {
                resource: SELF_RIGHTS_RESOURCE,
                actions: 1,
                domainType: "cloud",
                domainID: "*"
            },
            {
                resource: ADMIN_RIGHTS_RESOURCE,
                actions: 1,
                domainType: "cloud",
                domainID: "*"
            },
            {
                resource: SUPERADMIN_RIGHTS_RESOURCE,
                actions: 1,
                domainType: "cloud",
                domainID: "*"
            }
        ]

        return dispatch(api("/auth/validate", "POST", validations))
            .then(response => {
                dispatch({
                    type: LOAD_USER_RIGHTS_SUCCESS,
                    userRights: response
                })

                return response
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_RIGHTS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
