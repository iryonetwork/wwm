import _ from "lodash"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_LOCATION = "location/LOAD_LOCATION"
const LOAD_LOCATION_SUCCESS = "location/LOAD_LOCATION_SUCCESS"
const LOAD_LOCATION_FAIL = "location/LOAD_LOCATION_FAIL"

const LOAD_LOCATION_ORGANIZATION_IDS = "location/LOAD_LOCATION_ORGANIZATION_IDS"
const LOAD_LOCATION_ORGANIZATION_IDS_SUCCESS = "location/LOAD_LOCATION_ORGANIZATION_IDS_SUCCESS"
const LOAD_LOCATION_ORGANIZATION_IDS_FAIL = "location/LOAD_LOCATION_ORGANIZATION_IDS_FAIL"

const LOAD_LOCATION_USER_IDS = "location/LOAD_LOCATION_USER_IDS"
const LOAD_LOCATION_USER_IDS_SUCCESS = "location/LOAD_LOCATION_USER_IDS_SUCCESS"
const LOAD_LOCATION_USER_IDS_FAIL = "location/LOAD_LOCATION_USER_IDS_FAIL"

const LOAD_LOCATIONS = "location/LOAD_LOCATIONS"
const LOAD_LOCATIONS_SUCCESS = "location/LOAD_LOCATIONS_SUCCESS"
const LOAD_LOCATIONS_FAIL = "location/LOAD_LOCATIONS_FAIL"

const DELETE_LOCATION_FAIL = "location/DELETE_LOCATION_FAIL"
const DELETE_LOCATION_SUCCESS = "location/DELETE_LOCATION_SUCCESS"

const SAVE_LOCATION_FAIL = "location/SAVE_LOCATION_FAIL"
const SAVE_LOCATION_SUCCESS = "location/SAVE_LOCATION_SUCCESS"

const ROLE_ID_ALL = "all"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_LOCATION:
            return {
                ...state,
                loading: true
            }
        case LOAD_LOCATION_SUCCESS:
            return {
                ...state,
                locations: _.assign({}, state.locations || {}, _.fromPairs([[action.location.id, action.location]])),
                loading: false
            }
        case LOAD_LOCATION_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_LOCATION_ORGANIZATION_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_LOCATION_ORGANIZATION_IDS_SUCCESS:
            return {
                ...state,
                locationsOrganizationIDs: _.assign({}, state.locationsOrganizationIDs || {}, _fromPairs([[action.locationID, action.organizationIDs]])),
                loading: false
            }
        case LOAD_LOCATION_ORGANIZATION_IDS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_LOCATION_USER_IDS:
            return {
                ...state,
                loading: true
            }
        case LOAD_LOCATION_USER_IDS_SUCCESS:
            return {
                ...state,
                locationsUserIDs: _.assign({}, state.locationsUserIDs || {}, _fromPairs([[action.locationID, _.fromPairs([[action.roleID, action.userIDs]])]])),
                loading: false
            }
        case LOAD_LOCATION_USER_IDS_FAIL:
            return {
                ...state,
                loading: false
            }

        case LOAD_LOCATIONS:
            return {
                ...state,
                loading: true
            }
        case LOAD_LOCATIONS_SUCCESS:
            return {
                ...state,
                locations: _.keyBy(action.locations, "id"),
                loading: false
            }
        case LOAD_LOCATIONS_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case DELETE_LOCATION_SUCCESS:
            return {
                ...state,
                locations: _.pickBy(state.locations, location => location.id !== action.locationID)
            }

        case SAVE_LOCATION_SUCCESS:
            return {
                ...state,
                locations: _.assign({}, state.locations, _.fromPairs([[action.location.id, action.location]]))
            }
        default:
            return state
    }
}

export const loadLocation = locationID => {
    return dispatch => {
        dispatch({
            type: LOAD_LOCATION
        })

        return api(`/auth/locations/${locationID}`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_LOCATION_SUCCESS,
                    location: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_LOCATION_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadLocationOrganizationIDs = locationID => {
    return dispatch => {
        dispatch({
            type: LOAD_LOCATION_ORGANIZATION_IDS
        })

        return api(`/auth/locations/${locationID}/organizations`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_LOCATION_ORGANIZATION_IDS_SUCCESS,
                    locationID: locationID,
                    organizationIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_LOCATION_ORGANIZATION_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadLocationUserIDs = (locationID, roleID) => {
    return dispatch => {
        dispatch({
            type: LOAD_LOCATION_USER_IDS
        })

        var url = `/auth/locations/${locationID}/users`
        if (roleID && roleID !== ROLE_ID_ALL) {
            url += `?roleID=${roleID}`
        }

        return api(url, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_LOCATION_USER_IDS_SUCCESS,
                    locationID: locationID,
                    roleID: roleID ? roleID : ROLE_ID_ALL,
                    userIDs: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_LOCATION_USER_IDS_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const loadLocations = () => {
    return dispatch => {
        dispatch({
            type: LOAD_LOCATIONS
        })

        return api(`/auth/locations`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_LOCATIONS_SUCCESS,
                    locations: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_LOCATIONS_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteLocation = locationID => {
    return dispatch => {
        dispatch(close())

        return api(`/auth/locations/${locationID}`, "DELETE")
            .then(response => {
                dispatch({
                    type: DELETE_LOCATION_SUCCESS,
                    locationID: locationID
                })
                dispatch(open("Deleted location", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: DELETE_LOCATION_FAIL,
                    locationID: locationID
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveLocation = location => {
    return dispatch => {
        dispatch(close())
        let url = "/auth/locations"
        let method = "POST"
        if (location.id) {
            url += "/" + location.id
            method = "PUT"
        }

        return api(url, method, location)
            .then(response => {
                if (location.id) {
                    response = location
                }
                dispatch({
                    type: SAVE_LOCATION_SUCCESS,
                    location: response
                })
                dispatch(open("Saved location", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: SAVE_LOCATION_FAIL,
                    location: location
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
