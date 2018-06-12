import produce from "immer"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { read, API_URL, LOCATION_ID } from "shared/modules/config"
import { getToken } from "shared/modules/authentication"
import { updatePatient as updateWaitlistPatient } from "./waitlist"

export const SEARCH = "patient/SEARCH"
export const SEARCHED = "patient/SEARCHED"
export const FETCH = "patient/FETCH"
export const FETCHED = "patient/FETCHED"
export const UPDATE = "patient/UPDATE"
export const UPDATED = "patient/UPDATED"
export const FAILED = "patient/FAILED"

const initialState = {}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case SEARCH:
                draft.searching = true
                draft.searched = draft.failed = false
                break

            case SEARCHED:
                draft.searching = false
                draft.searched = true
                draft.patients = action.results
                break

            case FETCH:
                draft.fetching = true
                draft.fetched = draft.failed = false
                break

            case FETCHED:
                draft.fetching = false
                draft.fetched = true
                draft.patient = action.result
                break

            case UPDATE:
                draft.updating = true
                draft.updated = draft.failed = false
                break

            case UPDATED:
                draft.updating = false
                draft.updated = true
                draft.patient = action.result
                break

            case FAILED:
                draft.searching = draft.searched = false
                draft.fetching = draft.fetched = false
                draft.failed = true
                break

            default:
        }
    })
}

export const newPatient = formData => dispatch => {
    const url = `${dispatch(read(API_URL))}/discovery`

    var data = {
        connections: [
            { key: "firstName", value: formData.firstName },
            { key: "lastName", value: formData.lastName },
            { key: "dateOfBirth", value: formData.dateOfBirth },
            { key: "nationality", value: formData.nationality },
            { key: "gender", value: formData.gender },
            { key: "region", value: formData.region },
            { key: "address", value: formData.address || "" }
        ],
        locations: [dispatch(read(LOCATION_ID))]
    }
    ;(formData.documents || []).forEach(doc => {
        data.connections.push({ key: doc.type, value: doc.number })
    })

    return fetch(url, {
        method: "POST",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => Promise.all([response.status === 201, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to load insert new card / patient")
            }
            return data
        })
}

export const updatePatient = (patientID, formData) => dispatch => {
    const url = `${dispatch(read(API_URL))}/discovery/${patientID}`
    dispatch({ type: UPDATE })

    var data = {
        connections: [
            { key: "firstName", value: formData.firstName },
            { key: "lastName", value: formData.lastName },
            { key: "dateOfBirth", value: formData.dateOfBirth },
            { key: "nationality", value: formData.nationality },
            { key: "gender", value: formData.gender },
            { key: "region", value: formData.region },
            { key: "address", value: formData.address || "" }
        ],
        locations: [dispatch(read(LOCATION_ID))]
    }
    ;(formData.documents || []).forEach(doc => {
        data.connections.push({ key: doc.type, value: doc.number })
    })

    return fetch(url, {
        method: "PUT",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => Promise.all([response.status === 200, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to update patient")
            }

            dispatch(updateWaitlistPatient(patientID, data.connections))

            dispatch({ type: UPDATED, result: data })
            return data
        })
}

export const search = query => dispatch => {
    const url = `${dispatch(read(API_URL))}/discovery?query=${query}`
    dispatch({ type: SEARCH })

    return fetch(url, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.status === 200, response.json(), response.status]))
        .then(([ok, data, status]) => {
            if (!ok) {
                throw new Error(`Failed to search for patients (${status})`)
            }
            dispatch({ type: SEARCHED, results: data })
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}

export const get = patientID => dispatch => {
    const url = `${dispatch(read(API_URL))}/discovery/${patientID}`
    dispatch({ type: FETCH })

    return fetch(url, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.status === 200, response.json(), response.status]))
        .then(([ok, data, status]) => {
            if (!ok) {
                throw new Error(`Failed to fetch patient's details (${status})`)
            }
            dispatch({ type: FETCHED, result: data })
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}

export const cardToObject = card => {
    return card.connections.reduce((acc, conn) => {
        acc[conn.key] = conn.value
        return acc
    }, {})
}
