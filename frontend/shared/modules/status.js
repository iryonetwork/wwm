import { read, BASE_URL } from "./config"

const LOAD_STATUS = "user/LOAD_STATUS"
const LOAD_STATUS_SUCCESS = "user/LOAD_STATUS_SUCCESS"
const LOAD_STATUS_FAIL = "user/LOAD_STATUS_FAIL"

const initialState = {
    status: {status: "warning"},
    loading: false,
    forbidden: false
}

export default (state = initialState, action) => {
    switch (action.type) {
        case LOAD_STATUS:
            return {
                ...state,
                loading: true
            }
        case LOAD_STATUS_SUCCESS:
            return {
                ...state,
                status: action.status,
                loading: false
            }
        case LOAD_STATUS_FAIL:
            return {
                ...state,
                loading: false
            }
        default:
            return state
    }
}

export const loadStatus = () => {
    return dispatch => {
        dispatch({
            type: LOAD_STATUS
        })

        const url = `${dispatch(read(BASE_URL))}/status`

        return fetch(url, {
            method: "GET",
            headers: {"Content-Type": "application/json"}
        })
            .then(response => Promise.all([response.ok, response.json()]))
            .then(([ok, data]) => {
                if (!ok) {
                    dispatch({
                        type: LOAD_STATUS_FAIL,
                        status: {status: "error"},
                    })
                }
                dispatch({
                    type: LOAD_STATUS_SUCCESS,
                    status: data
                })
                return data
            })
    }
}
