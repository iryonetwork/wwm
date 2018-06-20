// DEV ONLY MODULE
import produce from "immer"
import moment from "moment"

// insert into storage
import { open, COLOR_DANGER } from "shared/modules/alert"
import { newPatient, updatePatient as updateDiscoveryPatient, get, cardToObject } from "./discovery"
import { extractPatientData, composePatientData, composeEncounterData, extractEncounterDataWithContext } from "./ehr"
import { createPatient as createPatientInStorage, readFileByLabel, updateFile, uploadFile, readFilesByLabel } from "./storage"
import { get as waitlistGet, remove as waitlistRemove } from "./waitlist"

export const CREATE = "patient/CREATE"
export const CREATED = "patient/CREATED"
export const LOADING = "patient/LOADING"
export const LOADED = "patient/LOADED"
export const SAVING_CONSULTATION = "patient/SAVING_CONSULTATION"
export const SAVED_CONSULTATION = "patient/SAVED_CONSULTATION"
export const FAILED = "patient/FAILED"

export const UPDATE = "patient/UPDATE"
export const UPDATE_DONE = "patient/UPDATE_DONE"
export const UPDATE_FAILED = "patient/UPDATE_FAILED"

export const FETCH_RECORDS = "patient/FETCH_RECORDS"
export const FETCH_RECORDS_DONE = "patient/FETCH_RECORDS_DONE"
export const FETCH_RECORDS_FAILED = "patient/FETCH_RECORDS_FAILED"

export const CLEAR_RECORDS_CACHE = "patient/CLEAR_RECORDS_CACHE"

const RECORDS_CACHE_DURATION = moment.duration(180, "seconds")

const newPatientFormData = {
    documents: [
        {
            type: "syrian-id",
            number: "1234567"
        },
        {
            type: "un-id",
            number: "987654"
        }
    ],
    firstName: "Patient",
    lastName: "Num1",
    dateOfBirth: "2018-05-21",
    gender: "CODED-at0310",
    maritalStatus: "SNOMED-125681006",
    numberOfKids: "1",
    nationality: "SI",
    countryOfOrigin: "SI",
    education: "gimnasium,",
    profession: "developer",
    country: "SI",
    region: "Kranj",
    phone: "040123456",
    email: "email@test.com",
    whatsapp: "987987654",
    dateOfLeaving: "2016-01-01",
    transitCountries: "Slovenia",
    dateOfArrival: "2016-12-31",
    people_in_family: "1",
    people_living_together: "2",
    familyMembers: [
        {
            firstName: "Test",
            lastName: "Test",
            dateOfBirth: "1993-05-04",
            relation: "child",
            livingTogether: "true",
            documentType: "syrian_id",
            documentNumber: "1234567"
        }
    ],
    allergies: [
        {
            allergy: "allergy 1",
            comment: "allergy comment 1",
            critical: "false"
        },
        {
            allergy: "allergy 2",
            comment: "allergy comment 2",
            critical: "true"
        }
    ],
    immunizations: [
        {
            immunization: "Immunization 1",
            date: "2014-02-01"
        },
        {
            immunization: "Immunuzation 2",
            date: "2004-03-02"
        }
    ],
    chronicDiseases: [
        {
            disease: "Chronic 1",
            date: "2004-03-12",
            medication: "Medication 1"
        },
        {
            disease: "Chronic 2 ",
            date: "2005-04-23",
            medication: "Medication 2"
        }
    ],
    injuries: [
        {
            injury: "Injury 1",
            date: "2007-06-05",
            medication: "Aids 1"
        },
        {
            injury: "Injury 2",
            date: "2009-08-07",
            medication: "Aids 2"
        }
    ],
    surgeries: [
        {
            injury: "Surgery 1",
            date: "2004-03-12",
            medication: "Comment 1"
        },
        {
            injury: "Surgery 2 ",
            date: "2004-03-12",
            medication: "Comment 2"
        }
    ],
    medications: [
        {
            medication: "Medication 1",
            comment: "Comment 1"
        },
        {
            medication: "Medication 2",
            comment: "Comment 2"
        }
    ],
    habits_smoking: "true",
    habits_smoking_comment: "10 boxes a day",
    habits_drugs: "true",
    habits_drugs_comment: "All of them",
    conditions_basic_hygiene: "true",
    conditions_heating: "true",
    conditions_good_appetite: "true",
    conditions_food_supply: "true",
    conditions_clean_water: "true",
    conditions_electricity: "true"
}

const initialState = {
    newData: process.env.NODE_ENV === "development" ? newPatientFormData : { documents: [] },
    patient: {},
    patientRecords: {
        cache: {}
    },
    patients: {}
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case CREATE:
                draft.creating = true
                draft.created = false
                break

            case CREATED:
                draft.creating = false
                draft.created = true
                break

            case LOADING:
                draft.loading = true
                draft.loaded = false
                break

            case LOADED:
                draft.loading = false
                draft.loaded = true
                draft.patient = action.result
                break

            case SAVING_CONSULTATION:
                draft.saving = true
                draft.saved = false
                break

            case SAVED_CONSULTATION:
                draft.saving = false
                draft.saved = true
                break

            case FAILED:
                draft.creating = draft.created = false
                draft.failed = true
                draft.reason = action.reason
                break

            case UPDATE:
                draft.updating = true
                break

            case UPDATE_DONE:
                draft.updating = false
                draft.patient = action.data
                break

            case UPDATE_FAILED:
                draft.updating = false
                break

            case FETCH_RECORDS:
                draft.patientRecords.loading = true
                break

            case FETCH_RECORDS_DONE:
                draft.patientRecords.loading = false
                draft.patientRecords.cache[action.patientID] = {
                    records: action.data,
                    fetchTimestamp: action.fetchTimestamp
                }
                draft.patientRecords.data = action.data
                break

            case FETCH_RECORDS_FAILED:
                draft.patientRecords.loading = false
                draft.patientRecords.data = undefined
                break

            case CLEAR_RECORDS_CACHE:
                if (action.patientID) {
                    delete draft.patientRecords.cache[action.patientID]
                } else {
                    draft.patientRecords.cache = {}
                }
                break

            default:
        }
    })
}

const getNewMembers = formData => {
    return (formData.familyMembers || []).filter(member => member.patientID === undefined).map(member => {
        if (member.livingTogether === "true") {
            member.tent = formData.tent
            member.camp = formData.camp
        }
        if (member.documentType) {
            member.documents = [
                {
                    type: member.documentType,
                    number: member.documentNumber
                }
            ]
        }
        return member
    })
}

export const createPatient = formData => (dispatch, getState) => {
    dispatch({ type: CREATE })
    // insert into discovery

    let familyMembers = getNewMembers(formData)

    return Promise.all(familyMembers.map(member => Promise.all([member, dispatch(newPatient(member))])))
        .then(members => {
            return Promise.all(
                members.map(([member, newPatient]) => {
                    return dispatch(composePatientData(member)).then(({ person }) => {
                        member.patientID = newPatient.patientID
                        return Promise.all([
                            dispatch(uploadFile(newPatient.patientID, person, "person", "openEHR-DEMOGRAPHIC-PERSON.person.v1")),
                            dispatch(uploadFile(newPatient.patientID, {}, "info", "openEHR-EHR-ITEM_TREE.patient_info.v0"))
                        ])
                    })
                })
            )
        })
        .then(() => {
            return dispatch(newPatient(formData)).then(newPatientCard => {
                // insert into storage
                return dispatch(createPatientInStorage(newPatientCard.patientID, formData)).then(result => {
                    dispatch({ type: CREATED })
                    return newPatientCard.patientID
                })
            })
        })
        .catch(ex => {
            dispatch(open(`Failed to create a new patient (${ex.message})`, "", COLOR_DANGER))
            // throw ex
        })
}

export const updatePatient = formData => (dispatch, getState) => {
    dispatch({ type: UPDATE })

    let patient = getState().patient.patient
    let familyMembers = getNewMembers(formData)

    return Promise.all(familyMembers.map(member => Promise.all([member, dispatch(newPatient(member))])))
        .then(members => {
            return Promise.all(
                members.map(([member, newPatient]) => {
                    return dispatch(composePatientData(member)).then(({ person }) => {
                        member.patientID = newPatient.patientID
                        return Promise.all([
                            dispatch(uploadFile(newPatient.patientID, person, "person", "openEHR-DEMOGRAPHIC-PERSON.person.v1")),
                            dispatch(uploadFile(newPatient.patientID, {}, "info", "openEHR-EHR-ITEM_TREE.patient_info.v0"))
                        ])
                    })
                })
            )
        })
        .then(() => dispatch(updateDiscoveryPatient(patient.ID, formData)))
        .then(() => {
            return dispatch(composePatientData(formData))
                .then(({ person, info }) => {
                    // upload user data to storage
                    return Promise.all([
                        dispatch(updateFile(patient.ID, patient.personFileID, person, "person", "openEHR-DEMOGRAPHIC-PERSON.person.v1")),
                        dispatch(updateFile(patient.ID, patient.infoFileID, info, "info", "openEHR-EHR-ITEM_TREE.patient_info.v0"))
                    ])
                })
                .then(result => {
                    dispatch({ type: UPDATE_DONE, data: formData })
                    return result
                })
        })
        .catch(ex => {
            dispatch({ type: UPDATE_FAILED })
            dispatch(open(`Failed to update patient (${ex.message})`, "", COLOR_DANGER))
        })
}

export const fetchPatient = patientID => dispatch => {
    dispatch({ type: LOADING })

    return Promise.all([dispatch(readFileByLabel(patientID, "person")), dispatch(readFileByLabel(patientID, "info")).catch(err => ({}))])
        .then(([person, info]) => Promise.all([dispatch(extractPatientData(person, info)), info.fileID, person.fileID]))
        .then(([patient, infoFileID, personFileID]) => {
            // add patientID
            patient.ID = patientID
            patient.infoFileID = infoFileID
            patient.personFileID = personFileID

            return Promise.all(
                (patient.familyMembers || []).map((member, index) => {
                    return dispatch(get(member.patientID)).then(card => {
                        let obj = cardToObject(card)

                        obj.documents = []
                        if (obj["syrian-id"]) {
                            obj.documents.push({
                                type: "syrian-id",
                                number: obj["syrian-id"]
                            })
                        }
                        if (obj["un-id"]) {
                            obj.documents.push({
                                type: "un-id",
                                number: obj["un-id"]
                            })
                        }

                        patient.familyMembers[index] = {
                            ...member,
                            ...obj
                        }
                    })
                })
            ).then(() => {
                dispatch({ type: LOADED, result: patient })
                return patient
            })
        })
}

export const saveConsultation = (waitlistID, itemID) => dispatch => {
    dispatch({ type: SAVING_CONSULTATION })

    // get waitlist item
    return (
        dispatch(waitlistGet(waitlistID, itemID))
            .then(item => {
                let data = {
                    patientID: item.patientID,
                    vitalSigns: item.vitalSigns || {},
                    mainComplaint: item.mainComplaint || {},
                    diagnoses: [],
                    therapies: []
                }
                ;(item.diagnoses || []).forEach((el, i) => {
                    data.diagnoses.push({
                        diagnosis: {
                            label: el.label,
                            id: el.diagnosis
                        },
                        comment: el.comment
                    })
                    // extract therapies
                    ;(el.therapies || []).forEach(therapy =>
                        data.therapies.push({
                            medication: therapy.medicine,
                            instructions: therapy.instructions,
                            diagnosis: i
                        })
                    )
                })
                return data
            })
            // create the document
            .then(data => Promise.all([Promise.resolve(data.patientID), dispatch(composeEncounterData(data))]))
            // upload the file
            .then(([patientID, doc]) =>
                Promise.all([Promise.resolve(patientID), dispatch(uploadFile(patientID, doc, "encounter", "openEHR-EHR-COMPOSITION.encounter.v1"))])
            )
            // remove from waitlist
            .then(([patientID, fileMeta]) => {
                dispatch(waitlistRemove(waitlistID, itemID, "finished"))
                dispatch({
                    type: CLEAR_RECORDS_CACHE,
                    patientID: patientID
                })
                dispatch({ type: SAVED_CONSULTATION })
            })
            .catch(ex => {
                console.log("failed to close", ex)
                dispatch({ type: FAILED })
                dispatch(open(`Failed to close consultation: {ex.message}`, "", COLOR_DANGER))
                throw ex
            })
    )
}

export const fetchHealthRecords = patientID => dispatch => {
    dispatch({ type: FETCH_RECORDS })

    let cachedRecords = dispatch(fetchCachedRecord(patientID))
    if (cachedRecords) {
        dispatch({
            type: FETCH_RECORDS_DONE,
            patientID: patientID,
            data: cachedRecords.records,
            fetchTimestamp: cachedRecords.fetchTimestamp
        })
        return
    }

    return dispatch(readFilesByLabel(patientID, "encounter"))
        .then(documents => Promise.all(documents.map(document => Promise.all([dispatch(extractEncounterDataWithContext(document.data)), document.meta]))))
        .then(documents => documents.map(([data, meta]) => ({ data, meta })))
        .then(documents => {
            dispatch({
                type: FETCH_RECORDS_DONE,
                patientID: patientID,
                data: documents,
                fetchTimestamp: moment()
            })
        })
        .catch(ex => {
            console.log(ex)
            dispatch({ type: FETCH_RECORDS_FAILED })
            dispatch(open("Failed to fetch health records", "", COLOR_DANGER))
        })
}

const fetchCachedRecord = patientID => (dispatch, getState) => {
    let cachedRecords = getState().patient.patientRecords.cache ? getState().patient.patientRecords.cache[patientID] : undefined
    if (cachedRecords) {
        if (
            cachedRecords.fetchTimestamp
                .clone()
                .add(RECORDS_CACHE_DURATION)
                .isAfter(moment())
        ) {
            return cachedRecords
        }
        dispatch({
            type: CLEAR_RECORDS_CACHE,
            patientID: patientID
        })
    }

    return undefined
}
