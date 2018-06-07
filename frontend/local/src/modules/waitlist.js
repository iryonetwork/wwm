import produce from "immer"
import { push } from "react-router-redux"
import _ from "lodash"
import moment from "moment"

import { open, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"
import { read, API_URL } from "shared/modules/config"
import { getToken } from "shared/modules/authentication"
import { fetchCode } from "shared/modules/codes"
import { round } from "shared/utils"

export const LIST = "waitlist/LIST"
export const LISTED = "waitlist/LISTED"
export const FETCH = "waitlist/FETCH"
export const FETCHED = "waitlist/FETCHED"
export const ADD = "waitlist/ADD"
export const ADDED = "waitlist/ADDED"
export const FAILED = "waitlist/FAILED"

export const UPDATE_ITEM = "waitlist/UPDATE_ITEM"
export const UPDATE_ITEM_DONE = "waitlist/UPDATE_ITEM_DONE"
export const UPDATE_ITEM_FAILED = "waitlist/UPDATE_ITEM_FAILED"

export const MOVE_TO_TOP_ITEM = "waitlist/MOVE_TO_TOP_ITEM"
export const MOVE_TO_TOP_ITEM_DONE = "waitlist/MOVE_TO_TOP_ITEM_DONE"
export const MOVE_TO_TOP_ITEM_FAILED = "waitlist/ MOVE_TO_TOP_ITEM_FAILED"

export const REMOVE_ITEM = "waitlist/REMOVE_ITEM"

export const RESET_INDICATORS = "waitlist/RESET_INDICATORS"

const initialState = {
    list: [],
    items: {}
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case RESET_INDICATORS:
                draft.listing = draft.listed = draft.fetching = draft.fetched = draft.adding = draft.added = draft.updating = draft.updated = false
                break

            case LIST:
                draft.listing = true
                draft.listed = draft.failed = false
                break

            case LISTED:
                draft.listing = false
                draft.listed = true
                draft.list = action.results
                draft.items = _.keyBy(action.results, "id")
                break

            case FETCH:
                draft.fetching = true
                draft.fetched = draft.failed = false
                break

            case FETCHED:
                draft.fetching = false
                draft.fetched = true
                draft.item = action.result
                break

            case ADD:
                draft.adding = true
                draft.added = draft.failed = false
                break

            case ADDED:
                draft.adding = false
                draft.added = true
                draft.item = action.result
                draft.items[action.result.id] = action.result
                break

            case FAILED:
                draft.listing = draft.listed = false
                draft.adding = draft.added = false
                draft.updating = draft.updated = false
                draft.failed = true
                break

            case UPDATE_ITEM:
                draft.updating = true
                draft.updated = draft.failed = false
                draft.items[action.itemID].updating = true
                break

            case UPDATE_ITEM_FAILED:
                draft.updating = draft.updated = false
                draft.failed = true
                draft.items[action.itemID].updating = false
                break
            case UPDATE_ITEM_DONE:
                draft.updating = false
                draft.updated = true
                draft.items[action.itemID] = action.updated
                draft.item = action.updated
                draft.items[action.itemID].updating = false
                break

            case MOVE_TO_TOP_ITEM:
                draft.items[action.itemID].updating = true
                break

            case MOVE_TO_TOP_ITEM_FAILED:
                draft.items[action.itemID].updating = false
                break
            case MOVE_TO_TOP_ITEM_DONE:
                draft.item = draft.items[action.itemID]
                draft.items[action.itemID].updating = false
                break

            case REMOVE_ITEM:
                delete draft.items[action.itemID]
                draft.list = _.filter(state.list, item => item.id !== action.itemID)
                break

            default:
        }
    })
}

export const add = (waitlistID, formData, patient) => dispatch => {
    const url = `${dispatch(read(API_URL))}/waitlist/${waitlistID}`
    dispatch({ type: ADD })

    const data = {
        patientID: patient.patientID,
        patient: patient.connections,
        priority: parseInt(formData.priority, 10) || 1,
        mainComplaint: {
            complaint: formData.mainComplaint,
            comment: formData.mainComplaintDetails
        }
    }

    return fetch(url, {
        method: "POST",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
        .then(response => Promise.all([response.status === 201, response.json(), response.status]))
        .then(([ok, data, status]) => {
            if (!ok) {
                if (status === 409) {
                    dispatch(listAll(waitlistID)).then(() => {
                        dispatch({ type: RESET_INDICATORS })
                    })
                    setTimeout(() => dispatch(open("Patient has been already added to the Waiting List", "", COLOR_DANGER, 5)), 100)
                    return undefined
                }
                throw new Error(`Failed to add patient to waitlist (${status})`)
            }
            dispatch({ type: ADDED, result: data })
            setTimeout(() => dispatch(open("Patient was added to waiting list", "", COLOR_SUCCESS, 5)), 100)
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}

export const listAll = listID => dispatch => {
    const url = `${dispatch(read(API_URL))}/waitlist/${listID}`
    dispatch({ type: LIST })

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

            let items = []
            ;(data || []).forEach((item, i) => {
                items.push(dispatch(migrateItem(item)))
            })

            return Promise.all(items)
        })
        .then(items => {
            dispatch({ type: LISTED, results: items })
            return items
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}

export const get = (waitlistID, itemID) => (dispatch, getState) => {
    dispatch({ type: FETCH })
    return dispatch(listAll(waitlistID)).then(list => {
        const items = (list || []).filter(item => item.id === itemID)

        if (items.length === 1) {
            dispatch({ type: FETCHED, result: items[0] })
            return items[0]
        }

        dispatch({ type: FAILED })
        dispatch(open("Waitlist item not found", "", COLOR_DANGER))
        throw new Error("waitlist item not found")
    })
}

export const update = (listID, data) => dispatch => {
    const url = `${dispatch(read(API_URL))}/waitlist/${listID}/${data.id}`
    dispatch({
        type: UPDATE_ITEM,
        itemID: data.id
    })

    let item = {
        ...data,
        priority: parseInt(data.priority, 10) || 1
    }

    return fetch(url, {
        method: "PUT",
        body: JSON.stringify(item),
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.status === 204, response.status]))
        .then(([ok, status]) => {
            if (!ok) {
                throw new Error(`Failed to update waiting list item (${status})`)
            }
            dispatch({
                type: UPDATE_ITEM_DONE,
                itemID: data.id,
                updated: data
            })
            return ok
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({
                type: UPDATE_ITEM_FAILED,
                itemID: data.id
            })
        })
}

export const moveToTop = (listID, itemID) => dispatch => {
    const url = `${dispatch(read(API_URL))}/waitlist/${listID}/${itemID}/top`
    dispatch({
        type: MOVE_TO_TOP_ITEM,
        itemID: itemID
    })

    return fetch(url, {
        method: "PUT",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.status === 204, response.status]))
        .then(([ok, status]) => {
            if (!ok) {
                throw new Error(`Failed to update waiting list item (${status})`)
            }
            dispatch({
                type: MOVE_TO_TOP_ITEM_DONE,
                itemID: itemID
            })
            dispatch(listAll(listID))
            setTimeout(() => dispatch(open("Waiting list was updated ", "", COLOR_SUCCESS, 5)), 100)
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({
                type: MOVE_TO_TOP_ITEM_FAILED,
                itemID: itemID
            })
        })
}

export const resetIndicators = () => dispatch => {
    dispatch({
        type: RESET_INDICATORS
    })
}

export const remove = (listID, itemID, reason) => dispatch => {
    const url = `${dispatch(read(API_URL))}/waitlist/${listID}/${itemID}?reason=${reason}`
    return fetch(url, {
        method: "DELETE",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.status === 204, response.status]))
        .then(([ok, status]) => {
            if (!ok) {
                throw new Error(`Failed to remove waiting list item (${status})`)
            }

            dispatch({
                type: REMOVE_ITEM,
                itemID: itemID
            })

            dispatch(push(`/waitlist/${listID}`))
            setTimeout(() => dispatch(open("Patient was removed from Waiting list", "", COLOR_SUCCESS, 5)), 100)
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
        })
}

// function to migrate old items to new format
const migrateItem = item => dispatch => {
    let diagnoses = []
    ;(item.diagnoses || []).forEach((diagnosis, i) => {
        diagnoses.push(dispatch(migrateDiagnosis(diagnosis)))
    })
    let vitalSigns = dispatch(migrateVitalSigns(item.vitalSigns))

    return Promise.all([Promise.all(diagnoses), vitalSigns]).then(([diagnoses, vitalSigns]) => {
        if (diagnoses.length > 0) {
            item.diagnoses = diagnoses
        }
        item.vitalSigns = vitalSigns
        return Promise.resolve(item)
    })
}

// migrating diagnosis with only SNOMED code and without label
const migrateDiagnosis = diagnosis => dispatch => {
    if (diagnosis.label) {
        return Promise.resolve(diagnosis)
    }

    return dispatch(fetchCode("diagnosis", diagnosis.diagnosis))
        .then(data => {
            diagnosis.label = data ? data.title : diagnosis.diagnosis
            return Promise.resolve(diagnosis)
        })
        .catch(ex => {
            diagnosis.label = diagnosis.diagnosis
            return Promise.resolve(diagnosis)
        })
}

// migrating vital signs without separately saved BMI
const migrateVitalSigns = vitalSigns => dispatch => {
    // migrate BMI (to be removed)
    if (vitalSigns && !vitalSigns.bmi && vitalSigns.height && vitalSigns.weight) {
        vitalSigns.bmi = {}
        vitalSigns.bmi.value = round(vitalSigns.weight.value / vitalSigns.height.value / vitalSigns.height.value * 10000, 2)
        vitalSigns.bmi.timestamp = moment.max(moment(vitalSigns.height.timestamp), moment(vitalSigns.weight.timestamp)).format()
    }

    return Promise.resolve(vitalSigns)
}

export const cardToObject = card => {
    return card.connections.reduce((acc, conn) => {
        acc[conn.key] = conn.value
        return acc
    }, {})
}
