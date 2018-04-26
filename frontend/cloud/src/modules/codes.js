import _ from "lodash"
import store from "../store"

import api from "./api"
import { open, COLOR_DANGER } from "shared/modules/alert"

const LOAD_CODES = "rules/LOAD_CODES"
const LOAD_CODES_SUCCESS = "rules/LOAD_CODES_SUCCESS"
const LOAD_CODES_FAIL = "rules/LOAD_CODES_FAIL"

const initialState = {
    loading: false,
    forbidden: false,
    codes: {}
}

export const CATEGORY_COUNTRIES = "countries"
export const CATEGORY_LANGUAGES = "languages"
export const CATEGORY_LICENSES = "licenses"

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_CODES:
            return {
                ...state,
                loading: true
            }
        case LOAD_CODES_SUCCESS:
            return {
                loading: false,
                codes: _.assign({}, state.codes || {}, _.fromPairs([[action.category, _.keyBy(action.codes, "id")]]))
            }
        case LOAD_CODES_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        default:
            return state
    }
}

export const loadCodes = category => {
    return dispatch => {
        dispatch({
            type: LOAD_CODES
        })

        let locale = store.getState().locale || "en"

        let url = "/discovery/codes/" + category + "?locale=" + locale
        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_CODES_SUCCESS,
                    category: category,
                    codes: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_CODES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
