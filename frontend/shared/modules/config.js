import produce from "immer"
import { defaultsDeep, get, findIndex } from "lodash"
import { open, COLOR_DANGER } from "./alert"

export const LOAD = "config/LOAD"
export const LOADED = "config/LOADED"
export const SAVED = "config/SAVED"
export const FAILED = "config/FAILED"

export const READ_ONLY_KEYS = "readOnlyKeys"
export const LOCALE = "locale"
export const BASE_URL = "baseUrl"
export const API_URL = "apiUrl"

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

export const save = (key, value) => (dispatch, getState) => {
    let readOnlyKeys = get(getState().config, READ_ONLY_KEYS, undefined)

    if (readOnlyKeys === undefined) {
        dispatch(open("Cannot save configuration", COLOR_DANGER))
        return
    }

    if (findIndex(READ_ONLY_KEYS, key) !== -1) {
        dispatch(open("Cannot save configuration", COLOR_DANGER))
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
