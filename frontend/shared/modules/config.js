import produce from "immer"
import { defaultsDeep, get } from "lodash"

// export const SET_LOCALE = "config/SET_LOCALE"

export const LOCALE = "locale"
export const BASE_URL = "baseUrl"
export const CLINIC_ID = "clinicId"
export const LOCATION_ID = "locationId"
export const BABY_MAX_AGE = "babyMaxAge"
export const CHILD_MAX_AGE = "childMaxAge"
export const DEFAULT_WAITLIST_ID = "waitlistId"

let initialState = {
    [LOCALE]: "en",
    [BASE_URL]: "https://iryo.local",
    [CLINIC_ID]: "e4ebb41b-7c62-4db7-9e1c-f47058b96dd0",
    [LOCATION_ID]: "2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa",
    [BABY_MAX_AGE]: 1,
    [CHILD_MAX_AGE]: 7,
    [DEFAULT_WAITLIST_ID]: "22afd921-0630-49f4-89a8-d1ad7639ee83"
}

// @TODO: Fetch config from URL

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            default:
        }
    })
}

export const read = (key, defaultValue) => {
    let conf = defaultsDeep(localStorage.getItem("config") || {}, initialState)
    return get(conf, key, defaultValue)
}
