import { composePatientData } from "./ehr"
import { read, BASE_URL } from "shared/modules/config"
import { getToken } from "shared/modules/authentication"

export const createPatient = (patientId, formData) => dispatch => {
    // compose files
    return dispatch(composePatientData(formData)).then(({ person, info }) => {
        // upload user data to storage
        return Promise.all([
            dispatch(uploadFile(patientId, person, "person", "openEHR-DEMOGRAPHIC-PERSON.person.v1")),
            dispatch(uploadFile(patientId, person, "info", "openEHR-EHR-ITEM_TREE.patient_info.v0"))
        ])
    })
}

export const uploadFile = (patientId, data, labels, archetype) => dispatch => {
    const url = `${read(BASE_URL)}/storage/${patientId}`

    try {
        var formData = new FormData()
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
    } catch (ex) {
        console.log(ex)
    }
}
