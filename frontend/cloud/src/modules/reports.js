import api from "shared/modules/api"
import { read, API_URL, REPORTS_STORAGE_BUCKET } from "shared/modules/config"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { getToken } from "shared/modules/authentication"

const READ_FILE = "reports/READ_FILE"
const READ_FILE_SUCCESS = "reports/READ_FILE_SUCCESS"
const READ_FILE_FAIL = "reports/READ_FILE_FAIL"

const LOAD_REPORTS_BY_TYPE = "reports/LOAD_REPORTS_BY_TYPE"
const LOAD_REPORTS_BY_TYPE_SUCCESS = "reports/LOAD_REPORTS_BY_TYPE_SUCCESS"
const LOAD_REPORTS_BY_TYPE_FAIL = "reports/LOAD_REPORTS_BY_TYPE_FAIL"

const initialState = {
    loading: false,
    reading: false,
    forbidden: false
}

export default (state = initialState, action) => {
    let reports, files, loadedFile
    switch (action.type) {
        case READ_FILE:
            return {
                ...state,
                reading: true
            }

        case READ_FILE_SUCCESS:
            loadedFile = {
                fileName: action.fileName,
                reportType: action.reportType,
                blob: action.blob
            }

            return {
                ...state,
                reading: false,
                loadedFile: loadedFile
            }

        case READ_FILE_FAIL:
            files = state.files || {}
            delete files[action.fileName]
            return {
                ...state,
                reading: false,
                files: files
            }

        case LOAD_REPORTS_BY_TYPE:
            return {
                ...state,
                loading: true
            }

        case LOAD_REPORTS_BY_TYPE_FAIL:
            reports = state.reports || {}
            reports[action.reportType] = {}

            return {
                ...state,
                loading: false,
                reports: reports
            }

        case LOAD_REPORTS_BY_TYPE_SUCCESS:
            reports = state.reports || {}
            reports[action.reportType] = action.reports

            return {
                ...state,
                loading: false,
                reports: reports
            }

        default:
            return state
    }
}

export const readFile = (reportType, fileName) => {
    return dispatch => {
        const url = `${dispatch(read(API_URL))}/storage/${dispatch(read(REPORTS_STORAGE_BUCKET))}/${fileName}`
        dispatch({
            type: READ_FILE
        })

        return fetch(url, {
            method: "GET",
            headers: {
                Authorization: dispatch(getToken()),
                "Content-Type": "application/json"
            }
        })
            .catch(error => {
                console.log(error)
                throw new Error("Failed to connect to server")
            })
            .then(response => Promise.all([response.status === 200, response.blob()]))
            .then(([ok, data]) => {
                if (!ok) {
                    throw new Error("Failed to read file from reports storage")
                }
                dispatch({
                    type: READ_FILE_SUCCESS,
                    reportType: reportType,
                    fileName: fileName,
                    blob: data
                })
                return data
            })
            .catch(error => {
                dispatch({
                    type: READ_FILE_FAIL,
                    reportType: reportType,
                    fileName: fileName
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadReportsByType = reportType => {
    return dispatch => {
        dispatch({
            type: LOAD_REPORTS_BY_TYPE
        })

        return dispatch(api(`/storage/${dispatch(read(REPORTS_STORAGE_BUCKET))}/${reportType}`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_REPORTS_BY_TYPE_SUCCESS,
                    reportType: reportType,
                    reports: response
                })

                return response
            })
            .catch(error => {
                dispatch({
                    type: LOAD_REPORTS_BY_TYPE_FAIL,
                    reportType: reportType
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
