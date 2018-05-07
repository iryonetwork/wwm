import { newPatient } from "./discovery"
// DEV ONLY MODULE
import produce from "immer"
// insert into storage
import { open, COLOR_DANGER } from "shared/modules/alert"
import { createPatient as createPatientInStorage } from "./storage"

export const CREATE = "patient/CREATE"
export const CREATED = "patient/CREATED"
export const LOADING = "patient/LOADING"
export const LOAD = "patient/LOAD"
export const FAILED = "patient/FAILED"

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
    dateOfBirth: "1983-05-21",
    gender: "CODED-at0310",
    maritalStatus: "SNOMED-125681006",
    numberOfKids: "1",
    nationality: "SI",
    countryOfOrigin: "SI",
    education: "gimnasium,",
    profession: "developer",
    country: "SI",
    camp: "19",
    tent: "83",
    clinic: "ZD Kranj",
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

            case FAILED:
                draft.creating = draft.created = false
                draft.failed = true
                draft.reason = action.reason
                break

            default:
        }
    })
}

export const createPatient = formData => (dispatch, getState) => {
    dispatch({ type: CREATE })
    // insert into discovery
    return dispatch(newPatient(formData))
        .then(newPatientCard => {
            // insert into storage
            return dispatch(createPatientInStorage(newPatientCard.patientID, formData)).then(result => {
                dispatch({ type: CREATED })
                return newPatientCard.patientID
            })
        })
        .catch(ex => {
            dispatch(open(`Failed to create a new patient (${ex.message})`, "", COLOR_DANGER))
            // throw ex
        })
}
