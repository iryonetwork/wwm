import { composePatientData } from "./ehr"
import { read } from "shared/modules/config"
import { API_URL } from "./config"
import { getToken } from "shared/modules/authentication"

export const createPatient = (patientId, formData) => dispatch => {
    // compose files
    return dispatch(composePatientData(formData)).then(({ person, info }) => {
        // upload user data to storage
        return Promise.all([
            dispatch(uploadFile(patientId, person, "person", "openEHR-DEMOGRAPHIC-PERSON.person.v1")),
            dispatch(uploadFile(patientId, info, "info", "openEHR-EHR-ITEM_TREE.patient_info.v0"))
        ])
    })
}

export const uploadFile = (patientId, data, labels, archetype) => dispatch => {
    const url = `${dispatch(read(API_URL))}/storage/${patientId}`

    let formData = new FormData()
    formData.append("file", new Blob([JSON.stringify(data)], { type: "application/json" }))
    formData.append("contentType", "application/json")
    formData.append("archetype", archetype)
    formData.append("labels", labels)

    return fetch(url, {
        method: "POST",
        headers: {
            Authorization: dispatch(getToken())
        },
        body: formData
    })
        .then(response => Promise.all([response.status === 201, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to upload file to storage")
            }
            return data
        })
}

export const updateFile = (patientID, fileID, data, labels, archetype) => dispatch => {
    const url = `${dispatch(read(API_URL))}/storage/${patientID}/${fileID}`

    let formData = new FormData()
    formData.append("file", new Blob([JSON.stringify(data)], { type: "application/json" }))
    formData.append("contentType", "application/json")
    formData.append("archetype", archetype)
    formData.append("labels", labels)

    return fetch(url, {
        method: "PUT",
        headers: {
            Authorization: dispatch(getToken())
        },
        body: formData
    })
        .then(response => Promise.all([response.status === 201, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to upload file to storage")
            }
            return data
        })
}

export const readFile = (patientID, fileID) => dispatch => {
    const url = `${dispatch(read(API_URL))}/storage/${patientID}/${fileID}`

    return fetch(url, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken())
        }
    })
        .then(response => Promise.all([response.status === 200, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to read file from storage")
            }
            data.fileID = fileID
            return data
        })
}

export const readFileByLabel = (patientID, label) => dispatch => {
    const labelUrl = `${dispatch(read(API_URL))}/storage/${patientID}/${label}`

    return fetch(labelUrl, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken())
        }
    })
        .then(response => Promise.all([response.status === 200, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to read label from storage")
            }
            if (data && data[0]) {
                return data[0]
            }
            throw new Error('no files found matching label "' + label + '"')
        })

        .then(fileData => {
            return dispatch(readFile(patientID, fileData.name))
        })
}

export const readFilesByLabel = (patientID, label) => dispatch => {
    const labelUrl = `${dispatch(read(API_URL))}/storage/${patientID}/${label}`

    return fetch(labelUrl, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken())
        }
    })
        .then(response => {
            if (response.status === 404) {
                return []
            } else if (response.status !== 200) {
                throw new Error("Failed to read label from storage")
            }
            return response.json()
        })

        .then(files => {
            return Promise.all(files.map(fileData => Promise.all([dispatch(readFile(patientID, fileData.name)), fileData])))
        })
        .then(files => files.map(([data, meta]) => ({ data: data, meta: meta })))
}
