import _ from "lodash"
import store from "../store"

import api from "./api"
import { loadDomainUserRoles, deleteUserRoleNoAlert, clearUserRoles } from "./userRoles"
import { clearLocations } from "./locations"
import { clearOrganizations } from "./organizations"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_CLINIC = "clinic/LOAD_CLINIC"
const LOAD_CLINIC_SUCCESS = "clinic/LOAD_CLINIC_SUCCESS"
const LOAD_CLINIC_FAIL = "clinic/LOAD_CLINIC_FAIL"

const LOAD_CLINICS = "clinic/LOAD_CLINICS"
const LOAD_CLINICS_SUCCESS = "clinic/LOAD_CLINICS_SUCCESS"
const LOAD_CLINICS_FAIL = "clinic/LOAD_CLINICS_FAIL"

const DELETE_CLINIC_FAIL = "clinic/DELETE_CLINIC_FAIL"
const DELETE_CLINIC_SUCCESS = "clinic/DELETE_CLINIC_SUCCESS"

const SAVE_CLINIC_FAIL = "clinic/SAVE_CLINIC_FAIL"
const SAVE_CLINIC_SUCCESS = "clinic/SAVE_CLINIC_SUCCESS"

const CLEAR_CLINICS_STATE = "clinic/CLEAR_CLINICS_STATE"

const initialState = {
    loading: false,
    allLoaded: false,
    forbidden: false
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_CLINIC:
            return {
                ...state,
                loading: true
            }
        case LOAD_CLINIC_SUCCESS:
            return {
                ...state,
                clinics: _.assign({}, state.clinics || {}, _.fromPairs([[action.clinic.id, action.clinic]])),
                loading: false
            }
        case LOAD_CLINIC_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_CLINICS:
            return {
                ...state,
                loading: true
            }
        case LOAD_CLINICS_SUCCESS:
            return {
                ...state,
                clinics: _.keyBy(action.clinics, "id"),
                allLoaded: true,
                loading: false
            }
        case LOAD_CLINICS_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_CLINIC_SUCCESS:
            return {
                ...state,
                clinics: _.pickBy(state.clinics, clinic => clinic.id !== action.clinicID)
            }

        case SAVE_CLINIC_SUCCESS:
            return {
                ...state,
                clinics: _.assign({}, state.clinics, _.fromPairs([[action.clinic.id, action.clinic]]))
            }

        case CLEAR_CLINICS_STATE:
            return {
                loading: false,
                clinics: undefined,
                allLoaded: false
            }

        default:
            return state
    }
}

export const loadClinic = clinicID => {
    return dispatch => {
        dispatch({
            type: LOAD_CLINIC
        })

        return api(`/auth/clinics/${clinicID}`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_CLINIC_SUCCESS,
                    clinic: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_CLINIC_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadClinics = () => {
    return dispatch => {
        dispatch({
            type: LOAD_CLINICS
        })

        return api(`/auth/clinics`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_CLINICS_SUCCESS,
                    clinics: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_CLINICS_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteClinic = clinicID => {
    return dispatch => {
        dispatch(close())

        return api(`/auth/clinics/${clinicID}`, "DELETE")
            .then(response => {
                dispatch(clearOrganizations())
                dispatch(clearUserRoles())
                dispatch(clearLocations())
                dispatch({
                    type: DELETE_CLINIC_SUCCESS,
                    clinicID: clinicID
                })
                dispatch(open("Deleted clinic", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: DELETE_CLINIC_FAIL,
                    clinicID: clinicID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveClinic = clinic => {
    return dispatch => {
        dispatch(close())

        let url = "/auth/clinics"
        let method = "POST"
        if (clinic.id) {
            url += "/" + clinic.id
            method = "PUT"
        }

        return api(url, method, clinic)
            .then(response => {
                dispatch(clearOrganizations())
                if (clinic.id) {
                    response = clinic
                }
                dispatch({
                    type: SAVE_CLINIC_SUCCESS,
                    clinic: response
                })
                dispatch(open("Saved clinic", "", COLOR_SUCCESS, 5))

                return response
            })
            .catch(error => {
                dispatch({
                    type: SAVE_CLINIC_FAIL,
                    clinic: clinic
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteUserFromClinic = (clinicID, userID) => {
    return dispatch => {
        dispatch(close())

        // check for user roles to delete in store
        let userRolesToDelete = undefined
        let clinicUserRoles =
            store.getState().userRoles.domainUserRoles &&
            store.getState().userRoles.domainUserRoles["clinic"] &&
            store.getState().userRoles.domainUserRoles["clinic"][clinicID]
                ? store.getState().userRoles.domainUserRoles["clinic"][clinicID]
                : undefined
        if (clinicUserRoles === undefined) {
            let userUserRoles =
                store.getState().userRoles.userUserRoles && store.getState().userRoles.userUserRoles[userID]
                    ? store.getState().userRoles.userUserRoles[userID]
                    : undefined
            if (userUserRoles !== undefined) {
                userRolesToDelete = _.pickBy(userUserRoles, userRole => userRole.domainType === "clinic" && userRole.domainID === clinicID) || {}
            }
        } else {
            userRolesToDelete = _.pickBy(clinicUserRoles, userRole => userRole.userID === userID) || {}
        }

        // no user roles to delete in store, fetch
        if (userRolesToDelete === undefined) {
            return dispatch(loadDomainUserRoles("clinic", clinicID)).then(() => {
                return dispatch(deleteUserFromClinic(clinicID, userID))
            })
        }

        _.forEach(userRolesToDelete, userRole => {
            dispatch(deleteUserRoleNoAlert(userRole.id))
        })

        return Promise.resolve()
    }
}

export const clearClinics = () => {
    return dispatch => {
        dispatch({
            type: CLEAR_CLINICS_STATE
        })

        return Promise.resolve()
    }
}
