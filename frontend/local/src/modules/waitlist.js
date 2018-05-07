import produce from "immer"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { read, BASE_URL, DEFAULT_WAITLIST_ID } from 'shared/modules/config'
import { getToken } from 'shared/modules/authentication'

export const LIST = "waitlist/LIST"
export const LISTED = "waitlist/LISTED"
// export const FETCH = "waitlist/FETCH"
// export const FETCHED = "waitlist/FETCHED"
export const ADD = "waitlist/ADD"
export const ADDED = "waitlist/ADDED"
export const FAILED = "waitlist/FAILED"

const initialState = {
    list: [],
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
                break

            // case FETCH:
            //     draft.fetching = true
            //     draft.fetched = draft.failed = false
            //     break

            // case FETCHED:
            //     draft.fetching = false
            //     draft.fetched = true
            //     draft.patient = action.result
            //     break

            case ADD:
                draft.adding = true
                draft.added = draft.failed = false
                break

            case ADDED:
                draft.adding = false
                draft.added = true
                draft.patient = action.result
                break

            case FAILED:
                draft.listing = draft.listed = false
                draft.adding = draft.added = false
                draft.failed = true
                break

            default:
        }
    })
}

// export const newPatient = (formData) => (
//     (dispatch) => {
//         const url = `${read(BASE_URL)}/discovery`

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
//             locations: [read(LOCATION_ID)],
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

export const add = (formData, patient) => (dispatch) => {
    const url = `${read(BASE_URL)}/waitlist/${read(DEFAULT_WAITLIST_ID)}`;
    const p = cardToObject(patient)
    dispatch({type: ADD});
    console.log(formData, patient)

    const data = {
        patient_id: patient.patientID,
        priority: parseInt(formData.priority, 10) || 1,
        complaint: formData.mainComplaint,
        patient: {
            name: `${p.lastName}, ${p.firstName}`,
            birthdate: p.dateOfBirth,
            gender: p.gender === 'CODED-at0310' ? 'M' : (p.gender === 'CODED-at0311' ? 'F' : '?')
        },
    }

    return fetch(url, {
        method: 'POST',
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data),
    })
        .then(response => Promise.all([response.status === 201, response.json(), response.status]))
        .then(([ok, data, status]) => {
            if (!ok) {
                throw new Error(`Failed to add patient to waitlist (${status})`)
            }
            dispatch({type: ADDED, results: data})
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({type: FAILED})
        })
}

export const listAll = (listID) => (dispatch) => {
    const url = `${read(BASE_URL)}/waitlist/${listID}`;
    dispatch({type: LIST});

    return fetch(url, {
        method: 'GET',
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        },
    })
        .then(response => Promise.all([response.status === 200, response.json(), response.status]))
        .then(([ok, data, status]) => {
            if (!ok) {
                throw new Error(`Failed to search for patients (${status})`)
            }
            dispatch({type: LISTED, results: data})
            return data
        })
        .catch(ex => {
            dispatch(open(ex.message, "", COLOR_DANGER))
            dispatch({type: FAILED})
        })
}

// export const get = (patientID) => (dispatch) => {
//     const url = `${read(BASE_URL)}/discovery/${patientID}`;
//     dispatch({type: FETCH});

//     return fetch(url, {
//         method: 'GET',
//         headers: {
//             Authorization: dispatch(getToken()),
//             "Content-Type": "application/json"
//         },
//     })
//         .then(response => Promise.all([response.status === 200, response.json(), response.status]))
//         .then(([ok, data, status]) => {
//             if (!ok) {
//                 throw new Error(`Failed to fetch patient's details (${status})`)
//             }
//             dispatch({type: FETCHED, result: data})
//             return data
//         })
//         .catch(ex => {
//             dispatch(open(ex.message, "", COLOR_DANGER))
//             dispatch({type: FAILED})
//         })
// }

export const cardToObject = (card) => {
    return card.connections.reduce(
        (acc, conn) => {
            acc[conn.key] = conn.value
            return acc
        },
        {}
    )
}
