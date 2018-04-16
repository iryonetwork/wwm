import _ from "lodash"

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
                codes: _.assign({}, state.codes || {}, _.fromPairs([[action.category, _.keyBy(action.codes, "code_id")]])),
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

        let codes = []
        // mocked codes response
        switch (category) {
            case CATEGORY_COUNTRIES:
                codes = [
                    {category_id: "codes", code_id: "DE", locale: "en", title: "Germany"},
                    {category_id: "codes", code_id: "FR", locale: "en", title: "France"},
                    {category_id: "codes", code_id: "PL", locale: "en", title: "Poland"},
                    {category_id: "codes", code_id: "SI", locale: "en", title: "Slovenia"},
                    {category_id: "codes", code_id: "GB", locale: "en", title: "United Kingdom"},
                    {category_id: "codes", code_id: "US", locale: "en", title: "United States of America"},
                ]
                break
            case CATEGORY_LANGUAGES:
                codes = [
                    {category_id: "languages", code_id: "AR", locale: "en", title: "Arabic"},
                    {category_id: "languages", code_id: "EN", locale: "en", title: "English"},
                    {category_id: "languages", code_id: "DE", locale: "en", title: "German"},
                    {category_id: "languages", code_id: "FR", locale: "en", title: "French"},
                    {category_id: "languages", code_id: "PL", locale: "en", title: "Polish"},
                    {category_id: "languages", code_id: "SI", locale: "en", title: "Slovenian"},
                ]
                break
            case CATEGORY_LICENSES:
                codes = [
                    {category_id: "licenses", code_id: "DL-A", locale: "en", title: "Driving license cat. A"},
                    {category_id: "licenses", code_id: "DL-A1", locale: "en", title: "Driving license cat. A1"},
                    {category_id: "licenses", code_id: "DL-B", locale: "en", title: "Driving license cat. B"},
                    {category_id: "licenses", code_id: "DL-C", locale: "en", title: "Driving license cat. C"},
                    {category_id: "licenses", code_id: "DL-C1", locale: "en", title: "Driving license cat. C1"},
                    {category_id: "licenses", code_id: "PL-P", locale: "en", title: "Private pilot license"},
                    {category_id: "licenses", code_id: "PL-C", locale: "en", title: "Commercial pilot license"},
                ]
                break
            default:
                codes = []
        }

        dispatch({
            type: LOAD_CODES_SUCCESS,
            category: category,
            codes: codes
        })

    }
}
