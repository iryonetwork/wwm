import produce from "immer"
import { defaultsDeep, get, findIndex } from "lodash"
import { open, COLOR_DANGER } from "./alert"

export const LOAD = "config/LOAD"
export const LOADED = "config/LOADED"
export const SAVED = "config/SAVED"
export const FAILED = "config/FAILED"

export const LOCALE = "locale"
export const BASE_URL = "baseUrl"
export const API_URL = "apiUrl"
export const CLINIC_ID = "clinicId"
export const LOCATION_ID = "locationId"
export const BABY_MAX_AGE = "babyMaxAge"
export const CHILD_MAX_AGE = "childMaxAge"
export const DEFAULT_WAITLIST_ID = "waitlistId"
export const ADVANCED_ROLE_IDS = "advancedRoleIDs"
export const REPORTS_STORAGE_BUCKET = "reportsStorageBucket"
export const LENGTH_UNIT = "lengthUnit"
export const WEIGHT_UNIT = "weightUnit"
export const TEMPERATURE_UNIT = "temperatureUnit"
export const BLOOD_PRESSURE_UNIT = "bloodPressureUnit"

const READ_ONLY_KEYS = [LOCALE, BASE_URL, API_URL, CLINIC_ID, LOCATION_ID, BABY_MAX_AGE, CHILD_MAX_AGE, ADVANCED_ROLE_IDS, REPORTS_STORAGE_BUCKET]

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOAD:
                draft.loading = true
                break

            case LOADED:
                draft = Object.assign(draft, action.result)
                draft.loading = false
                break

            case FAILED:
                draft.failed = true
                break

            case SAVED:
                draft[action.key] = action.value
                break

            default:
        }
    })
}

export const save = (key, value) => dispatch => {
    if (findIndex(READ_ONLY_KEYS, key) !== -1) {
        dispatch(open("Cannot save to config key " + key + " as it is read-only.", COLOR_DANGER))
        return
    }
    let data = JSON.parse(localStorage.getItem("config"))
    data[key] = value
    localStorage.setItem("config", JSON.stringify(data))

    dispatch({
        type: SAVED,
        key: key,
        value: value
    })
}

export const read = (key, defaultValue) => (dispatch, getState) => {
    return get(getState().config, key, defaultValue)
}

export const load = dispatch => {
    dispatch({ type: LOAD })

    fetch("/config.json", {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.ok, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to load config")
            }
            data = defaultsDeep(JSON.parse(localStorage.getItem("config")) || {}, data)
            localStorage.setItem("config", JSON.stringify(data))

            dispatch({ type: LOADED, result: data })
            return data
        })
        .catch(ex => {
            dispatch(open("Failed to load config :: " + ex.message, COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}
