import _ from "lodash"
import store from '../store'

import api from "./api"
import { open, close, COLOR_DANGER } from "shared/modules/alert"

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
    console.log(action)
    switch (action.type) {
        case LOAD_CODES:
            return {
                ...state,
                loading: true
            }
        case LOAD_CODES_SUCCESS:
            return {
                loading: false,
                codes: _.assign({}, state.codes || {}, _.fromPairs([[action.category, _.keyBy(action.codes, "id")]])),
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
        let codes = []

        // mocked codes response
        switch (category) {
            case CATEGORY_LANGUAGES:
                codes = [
                    {category: "languages", id: "AR", locale: "en", title: "Arabic"},
                    {category: "languages", id: "EN", locale: "en", title: "English"},
                    {category: "languages", id: "DE", locale: "en", title: "German"},
                    {category: "languages", id: "FR", locale: "en", title: "French"},
                    {category: "languages", id: "PL", locale: "en", title: "Polish"},
                    {category: "languages", id: "SI", locale: "en", title: "Slovenian"},
                ]
                dispatch({
                    type: LOAD_CODES_SUCCESS,
                    category: category,
                    codes: codes
                })
                break
            case CATEGORY_LICENSES:
                codes = [
                    {category: "licenses", id: "DL-A", locale: "en", title: "Driving license cat. A"},
                    {category: "licenses", id: "DL-A1", locale: "en", title: "Driving license cat. A1"},
                    {category: "licenses", id: "DL-B", locale: "en", title: "Driving license cat. B"},
                    {category: "licenses", id: "DL-C", locale: "en", title: "Driving license cat. C"},
                    {category: "licenses", id: "DL-C1", locale: "en", title: "Driving license cat. C1"},
                    {category: "licenses", id: "PL-P", locale: "en", title: "Private pilot license"},
                    {category: "licenses", id: "PL-C", locale: "en", title: "Commercial pilot license"},
                ]
                dispatch({
                    type: LOAD_CODES_SUCCESS,
                    category: category,
                    codes: codes
                })
                break
            default:
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
}
