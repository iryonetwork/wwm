import produce from "immer"
import { read, LOCALE, BASE_URL } from "./config"
import { open, close, COLOR_DANGER } from "./alert"
import { getToken } from "./authentication"

export const LOADING = "codes/LOADING"
export const LOADED = "codes/LOADED"
export const LOAD = "codes/LOAD"
export const FAILED = "codes/FAILED"

const initialState = {
    loading: false,
    cache: {},
    fetching: [],
    failed: [],
}

const locale = 'en'
const onlyUnique = (value, index, self) => self.indexOf(value) === index

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOADING:
                draft.loading = true
                // add category to list of fetching categories
                if (action.category && draft.fetching.indexOf(action.category) === -1) {
                    draft.fetching.push(action.category)
                }
                break

            case LOADED:
                if (action.category) {
                    draft.cache[action.category] = action.data
                    draft.fetching = draft.fetching.filter(cat => cat !== action.category)
                }

                // mark everything as loaded if the
                if (draft.fetching.length === 0) {
                    draft.loading = false
                }
                break

            case FAILED:
                if (draft.failed.indexOf(action.category) === -1) {
                    draft.failed.push(action.category)
                }

            default:
        }
    })
}

export const getCodes = (category) => (dispatch, getState) => {
    return []
}

export const getCodesAsOptions = (category) => {
    return (dispatch, getState) => {
        return (getState().codes.cache[category] || []).map((code => (
            {label: code.title, value: code.id}
        )))
    }
}

export const loadCategories = (...categories) => {
    return (dispatch, getState) => {
        const state = getState().codes
        const requestedCategories = (categories || []).length
        let preloadedCategories = 0
        dispatch({type: LOADING})

        if (typeof getState !== 'function') {
            return
        }

        // iterate over categories an start loading missing categories
        (categories || []).forEach(category => {
            // skip if category is
            if (state.cache.hasOwnProperty(category)) {
                preloadedCategories++
                return
            }

            // skip if category is already being loaded
            if ((state.fetching || []).indexOf(category) !== -1) {
                return
            }

            // load the category
            dispatch(load(category))
        });

        // try to mark as done if all categories were skipped
        if (requestedCategories === preloadedCategories) {
            dispatch({type: LOADED})
        }
    }
}

const load = (category) => {
    return dispatch => {
        dispatch({type: LOADING, category})
        let locale = read(LOCALE)
        const url = `${read(BASE_URL)}/discovery/codes/${category}?locale=${locale}`
        fetch(url, {
            method: 'GET',
            headers: {
                Authorization: dispatch(getToken()),
                "Content-Type": "application/json"
            }
        })
            .then(response => Promise.all([response.ok, response.json()]))
            .then(([ok, data]) => {
                if (ok) {
                    return dispatch({type: LOADED, category, data})
                }
                throw new Error('Failed to load codes')
            })
            .catch(ex => {
                dispatch(open('Failed to fetch codes :: '+ex.message, COLOR_DANGER))
                dispatch({type: FAILED, category})
            })
    }
}

/*
export const open = (message, code, color, closeIn) => {
    return dispatch => {
        dispatch({
            type: SHOW_ALERT,
            message,
            code,
            color
        })

        if (closeIn) {
            setTimeout(() => {
                dispatch({ type: HIDE_ALERT })
            }, closeIn * 1000)
        }
    }
}

export const close = () => {
    return dispatch => {
        dispatch({ type: HIDE_ALERT })
    }
}
*/
