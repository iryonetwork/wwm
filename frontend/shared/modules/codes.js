import produce from "immer"
import _ from "lodash"

import { read, LOCALE, API_URL } from "./config"
import { open, COLOR_DANGER } from "./alert"
import { getToken } from "./authentication"

export const LOADING = "codes/LOADING"
export const LOADED = "codes/LOADED"
export const LOAD = "codes/LOAD"
export const SEARCH = "codes/SEARCH"
export const SEARCHED = "codes/SEARCHED"
export const FETCH = "codes/FETCH"
export const FETCHED = "codes/FETCHED"
export const FAILED = "codes/FAILED"

const initialState = {
    loading: false,
    cache: {},
    fetching: [],
    codeFetching: {},
    fetchedCache: {},
    failed: []
}

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
                    draft.cache[action.category] = _.reduce(
                        action.data,
                        (result, data) => {
                            result[data.id] = data
                            return result
                        },
                        {}
                    )
                    draft.fetching = draft.fetching.filter(cat => cat !== action.category)
                }

                // mark everything as loaded if the
                if (draft.fetching.length === 0) {
                    draft.loading = false
                }
                break

            case SEARCH:
                draft.searching = true
                break

            case SEARCHED:
                draft.searching = false
                draft.searchResults = action.data
                break

            case FETCH:
                draft.codeFetching[action.category] = _.assign(draft.codeFetching[action.category] || {}, _.fromPairs([[action.id, true]]))
                draft.isFetching = true
                break

            case FETCHED:
                draft.codeFetching[action.category] = _.assign(draft.codeFetching[action.category] || {}, _.fromPairs([[action.id, false]]))
                draft.fetchedCache[action.category] = _.assign(draft.fetchedCache[action.category] || {}, _.fromPairs([[action.id, action.data]]))
                draft.isFetching = false
                draft.fetchResults = action.data
                break

            case FAILED:
                draft.fetching = draft.fetching.filter(cat => cat !== action.category)

                if (draft.failed.indexOf(action.category) === -1) {
                    draft.failed.push(action.category)
                }

                if (draft.fetching.length === 0) {
                    draft.loading = false
                }
                break

            default:
        }
    })
}

export const getCodes = category => (dispatch, getState) => {
    return getState().codes.cache[category] || {}
}

export const getCodesAsOptions = category => {
    return (dispatch, getState) => {
        return _.map(getState().codes.cache[category] || {}, code => ({ label: code.title, value: code.id }))
    }
}

export const searchCodes = (category, query, parentId) => dispatch => {
    dispatch({ type: SEARCH })
    const locale = dispatch(read(LOCALE))
    const url = `${dispatch(read(API_URL))}/discovery/codes/${category}?locale=${locale}`
    return fetch(url, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.ok, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to load codes")
            }
            dispatch({ type: SEARCHED, data })
            return data
        })
        .catch(ex => {
            dispatch(open("Failed to fetch codes :: " + ex.message, COLOR_DANGER))
            dispatch({ type: FAILED, category })
        })
}

export const loadCategories = (...categories) => {
    return (dispatch, getState) => {
        const state = getState().codes
        const requestedCategories = (categories || []).length
        let preloadedCategories = 0
        dispatch({ type: LOADING })

        if (typeof getState !== "function") {
            return
        }

        // iterate over categories an start loading missing categories
        ;(categories || []).forEach(category => {
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
        })

        // try to mark as done if all categories were skipped
        if (requestedCategories === preloadedCategories) {
            dispatch({ type: LOADED })
        }
    }
}

export const fetchCode = (category, id) => (dispatch, getState) => {
    const locale = dispatch(read(LOCALE))
    const url = `${dispatch(read(API_URL))}/discovery/codes/${category}/${id}?locale=${locale}`
    dispatch({ type: FETCH, category, id })

    // check in cache of whole categories
    let cachedCode = _.get(getState().codes.cache, `${category}.${id}`)
    if (cachedCode) {
        dispatch({ type: FETCHED, data: cachedCode, category, id })

        return Promise.resolve(cachedCode)
    }

    // check in cache of individually fetched codes
    cachedCode = _.get(getState().codes.fetchedCache, `${category}.${id}`)
    if (cachedCode) {
        dispatch({ type: FETCHED, data: cachedCode, category, id })

        return Promise.resolve(cachedCode)
    }

    return fetch(url, {
        method: "GET",
        headers: {
            Authorization: dispatch(getToken()),
            "Content-Type": "application/json"
        }
    })
        .then(response => Promise.all([response.ok, response.json()]))
        .then(([ok, data]) => {
            if (!ok) {
                throw new Error("Failed to fetch code")
            }
            dispatch({ type: FETCHED, category, id, data })
            return data
        })
        .catch(ex => {
            dispatch(open("Failed to fetch codes :: " + ex.message, COLOR_DANGER))
            dispatch({ type: FAILED, category, id })
        })
}

export const load = category => {
    return dispatch => {
        dispatch({ type: LOADING, category })
        const locale = dispatch(read(LOCALE))
        const url = `${dispatch(read(API_URL))}/discovery/codes/${category}?locale=${locale}`

        return fetch(url, {
            method: "GET",
            headers: {
                Authorization: dispatch(getToken()),
                "Content-Type": "application/json"
            }
        })
            .then(response => Promise.all([response.ok, response.json()]))
            .then(([ok, data]) => {
                if (!ok) {
                    throw new Error("Failed to load codes")
                }
                dispatch({ type: LOADED, category, data })
                return data
            })
            .catch(ex => {
                dispatch(open("Failed to fetch codes :: " + ex.message, COLOR_DANGER))
                dispatch({ type: FAILED, category })
            })
    }
}
