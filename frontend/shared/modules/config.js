import produce from "immer"
import { defaultsDeep, get } from "lodash"

// export const SET_LOCALE = "config/SET_LOCALE"

export const LOCALE = 'locale'
export const BASE_URL = 'baseUrl'

let initialState = {}
initialState[LOCALE] = 'en'
initialState[BASE_URL] = 'https://iryo.local'

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            default:
        }
    })
}

export const read = (key, defaultValue) => {
    let conf = defaultsDeep(localStorage.getItem('config') || {}, initialState)
    return get(conf, key, defaultValue)
}
