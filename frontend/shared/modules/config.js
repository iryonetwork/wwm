import produce from "immer"
import { defaultsDeep, get } from "lodash"
import { open, COLOR_DANGER } from "./alert"

export const LOAD = "config/LOAD"
export const LOADED = "config/LOADED"
export const FAILED = "config/FAILED"

export const LOCALE = "locale"
export const BASE_URL = "baseUrl"
export const CLINIC_ID = "clinicId"
export const LOCATION_ID = "locationId"
export const BABY_MAX_AGE = "babyMaxAge"
export const CHILD_MAX_AGE = "childMaxAge"
export const DEFAULT_WAITLIST_ID = "waitlistId"

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

            default:
        }
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
            data = defaultsDeep(localStorage.getItem("config") || {}, data)
            dispatch({ type: LOADED, result: data })
            return data
        })
        .catch(ex => {
            dispatch(open("Failed to load config :: " + ex.message, COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}
