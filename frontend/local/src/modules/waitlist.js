import produce from "immer"
import { goBack, push } from "react-router-redux"
import _ from "lodash"
import { open, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"
import { read, BASE_URL, DEFAULT_WAITLIST_ID } from "shared/modules/config"
import { getToken } from "shared/modules/authentication"

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

export const REMOVE_ITEM = "waitlist/REMOVE_ITEM"

const initialState = {
    list: [],
    items: {}
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
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
                break

            case FAILED:
                draft.listing = draft.listed = false
                draft.adding = draft.added = false
                draft.failed = true
                break

            case UPDATE_ITEM:
                draft.items[action.itemID].updating = true
                break

            case UPDATE_ITEM_FAILED:
            case UPDATE_ITEM_DONE:
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

// export const newPatient = (formData) => (
//     (dispatch) => {
//         const url = `${dispatch(read(BASE_URL))}/discovery`

//         var data = {
//             connections: [
//                 {key: 'firstName', value: formData.firstName},
//                 {key: 'lastName', value: formData.lastName},
//                 {key: 'dateOfBirth', value: formData.dateOfBirth},
//                 {key: 'nationality', value: formData.nationality},
//                 {key: 'gender', value: formData.gender},
//                 {key: 'tent', value: formData.tent},
//                 {key: 'camp', value: formData.camp},
//             ],
//             locations: [dispatch(read(LOCATION_ID))],
//         };

//         (formData.documents || []).forEach(doc => {
//             data.connections.push({key: doc.type, value: doc.number})
//         });

//         return fetch(url, {
//             method: 'POST',
//             headers: {
//                 Authorization: dispatch(getToken()),
//                 "Content-Type": "application/json"
//             },
//             body: JSON.stringify(data),
//         })
//             .then(response => Promise.all([response.status === 201, response.json()]))
//             .then(([ok, data]) => {
//                 if (!ok) {
//                     throw new Error('Failed to load insert new card / patient')
//                 }
//                 return data
//             })
//     }
// )

export const add = (formData, patient) => dispatch => {
    const waitlistID = dispatch(read(DEFAULT_WAITLIST_ID))
    const url = `${dispatch(read(BASE_URL))}/waitlist/${waitlistID}`
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
                throw new Error(`Failed to add patient to waitlist (${status})`)
            }
            dispatch({ type: ADDED, result: data })
            dispatch(goBack())
            setTimeout(() => dispatch(open("Patient was added to waiting list", "", COLOR_SUCCESS, 5)), 100)
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({ type: FAILED })
        })
}

export const listAll = listID => dispatch => {
    const url = `${dispatch(read(BASE_URL))}/waitlist/${listID}`
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
            dispatch({ type: LISTED, results: data })
            return data
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
    // const url = `${dispatch(read(BASE_URL))}/waitlist/${waitlistID}/${itemID}`

    // return fetch(url, {
    //     method: "GET",
    //     headers: {
    //         Authorization: dispatch(getToken()),
    //         "Content-Type": "application/json"
    //     }
    // })
    //     .then(response => Promise.all([response.status === 200, response.json(), response.status]))
    //     .then(([ok, data, status]) => {
    //         if (!ok) {
    //             throw new Error(`Failed to fetch waitlist item (${status})`)
    //         }
    //         dispatch({ type: FETCHED, result: data })
    //         return data
    //     })
    //     .catch(ex => {
    //         dispatch(open(ex.message, "", COLOR_DANGER))
    //         dispatch({ type: FAILED })
    //     })
}

export const update = (listID, data) => dispatch => {
    const url = `${read(BASE_URL)}/waitlist/${listID}/${data.id}`
    dispatch({
        type: UPDATE_ITEM,
        itemID: data.id
    })

    return fetch(url, {
        method: "PUT",
        body: JSON.stringify(data),
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
                itemID: data.id
            })
            dispatch(goBack())
            setTimeout(() => dispatch(open("Waiting list was updated ", "", COLOR_SUCCESS, 5)), 100)
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({
                type: UPDATE_ITEM_FAILED,
                itemID: data.id
            })
        })
}

export const remove = (listID, itemID, reason) => dispatch => {
    const url = `${read(BASE_URL)}/waitlist/${listID}/${itemID}?reason=${reason}`

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

export const cardToObject = card => {
    return card.connections.reduce((acc, conn) => {
        acc[conn.key] = conn.value
        return acc
    }, {})
}
