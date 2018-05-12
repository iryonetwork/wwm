import _ from "lodash"

import api from "shared/modules/api"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { read, CLINIC_ID } from "shared/modules/config"

const LOAD_USER_RIGHTS = "user/LOAD_USER_RIGHTS"
const LOAD_USER_RIGHTS_SUCCESS = "user/LOAD_USER_RIGHTS_SUCCESS"
const LOAD_USER_RIGHTS_FAIL = "user/LOAD_USER_RIGHTS_FAIL"

export const READ = 1
export const WRITE = 2
export const DELETE = 4
export const UPDATE = 8

export const ROLE_ADMIN = "/frontend/role/admin"
export const ROLE_DOCTOR = "/frontend/role/doctor"
export const ROLE_NURSE = "/frontend/role/nurse"
export const ROLE_PHARMACY = "/frontend/role/pharmacy"

export const RESOURCE_WAITLIST = "/frontend/waitlist"

export const RESOURCE_PATIENT_IDENTIFICATION = "/frontend/patientID/id"
export const RESOURCE_DEMOGRAPHIC_INFORMATION = "/frontend/patientID/demographic"
export const RESOURCE_BABY_SCREENING = "/frontend/patientID/babyScreening"

export const RESOURCE_VITAL_SIGNS = "/frontend/testsAndHistory/vitalSigns"
export const RESOURCE_HEALTH_HISTORY = "/frontend/testsAndHistory/healthHistory"
export const RESOURCE_LABORATORY_TEST = "/frontend/testsAndHistory/laboratoryTest"

export const RESOURCE_EXAMINATION = "/frontend/examination/examination"
export const RESOURCE_MEDICATION = "/frontend/examination/medication"


const initialState = {
    loading: false,
    forbidden: false
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_USER_RIGHTS:
            return {
                ...state,
                loading: true
            }
        case LOAD_USER_RIGHTS_SUCCESS:
            return {
                ...state,
                userRights: _.reduce(
                    action.userRights,
                    function(result, value, key) {
                        result[value.query.resource] || (result[value.query.resource] = {})
                        result[value.query.resource][value.query.actions] = value.result
                        return result
                    },
                    {}
                ),
                loading: false
            }
        case LOAD_USER_RIGHTS_FAIL:
            return {
                ...state,
                loading: false
            }
        default:
            return state
    }
}

export const loadUserRights = userID => {
    return dispatch => {
        dispatch({
            type: LOAD_USER_RIGHTS
        })

        let clinicID = dispatch(read(CLINIC_ID))

        let validations = [
            {
                resource: ROLE_ADMIN,
                actions: READ,
                domainType: "clinic",
                domainID: clinicID
            },
            {
                resource: ROLE_DOCTOR,
                actions: READ,
                domainType: "clinic",
                domainID: clinicID
            },
            {
                resource: ROLE_NURSE,
                actions: READ,
                domainType: "clinic",
                domainID: clinicID
            },
            {
                resource: ROLE_PHARMACY,
                actions: READ,
                domainType: "clinic",
                domainID: clinicID
            },
        ]

        let resources = [
            RESOURCE_WAITLIST,
            RESOURCE_PATIENT_IDENTIFICATION,
            RESOURCE_DEMOGRAPHIC_INFORMATION,
            RESOURCE_BABY_SCREENING,
            RESOURCE_VITAL_SIGNS,
            RESOURCE_HEALTH_HISTORY,
            RESOURCE_LABORATORY_TEST,
            RESOURCE_EXAMINATION,
            RESOURCE_MEDICATION
        ]

        for (var i in resources) {
            let resourceValidations = [
                {
                    resource: resources[i],
                    actions: READ,
                    domainType: "clinic",
                    domainID: clinicID
                },
                {
                    resource: resources[i],
                    actions: WRITE,
                    domainType: "clinic",
                    domainID: clinicID
                },
                {
                    resource: resources[i],
                    actions: DELETE,
                    domainType: "clinic",
                    domainID: clinicID
                },
                {
                    resource: resources[i],
                    actions: UPDATE,
                    domainType: "clinic",
                    domainID: clinicID
                },
            ]
            validations.push(...resourceValidations)
        }

        return dispatch(api("/auth/validate", "POST", validations))
            .then(response => {
                dispatch({
                    type: LOAD_USER_RIGHTS_SUCCESS,
                    userRights: response
                })
                return response
            })
            .catch(error => {
                dispatch({
                    type: LOAD_USER_RIGHTS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
