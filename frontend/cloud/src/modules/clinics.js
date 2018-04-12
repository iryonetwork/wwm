import _ from "lodash"

import api from "./api"
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

const initialState = {
    loading: true
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
                if (clinic.id) {
                    response = clinic
                }
                dispatch({
                    type: SAVE_CLINIC_SUCCESS,
                    clinic: response
                })
                dispatch(open("Saved clinic", "", COLOR_SUCCESS, 5))
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
