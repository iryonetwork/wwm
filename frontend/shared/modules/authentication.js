import produce from "immer"
import { push } from "react-router-redux"
import jwtDecode from "jwt-decode"
import { read, API_URL } from "./config"

import { open, close, COLOR_DANGER } from "./alert"
import { load as loadConfig } from "./config"

export const LOGIN = "authentication/LOGIN"
export const ERROR = "authentication/ERROR"
export const LOGGEDIN = "authentication/LOGGEDIN"
export const SET_USERNAME = "authentication/SET_USERNAME"
export const SET_PASSWORD = "authentication/SET_PASSWORD"
export const RENEW_RETRY = "authentication/RENEW_RETRY"
export const RENEW_FAILED = "authentication/RENEW_FAILED"

const initialState = {
    form: {},
    retries: 0
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOGIN:
                draft.pending = true
                draft.error = false
                break
            case ERROR:
                draft.pending = false
                draft.error = true
                break
            case LOGGEDIN:
                draft.error = null
                draft.pending = false
                draft.redirectToReferrer = true
                draft.tokenString = action.token
                draft.token = jwtDecode(action.token)
                break
            case SET_USERNAME:
                draft.form.username = action.username
                break
            case SET_PASSWORD:
                draft.form.password = action.password
                break
            case RENEW_RETRY:
                draft.redirectToReferrer = false
                draft.retries = state.retries + 1
                break
            case RENEW_FAILED:
                draft.token = null
                draft.redirectToReferrer = false
                draft.retries = 0
                draft.form = {}
                break
            default:
        }
    })
}

// 10 minutes
const renewInterval = 10 * 60 * 1000

export const renewToken = () => {
    return (dispatch, getState) => {
        return fetch(`${dispatch(read(API_URL))}/auth/renew`, {
            headers: {
                "Content-Type": "application/json",
                Authorization: getState().authentication.tokenString
            }
        })
            .catch(ex => {
                throw ex
            })
            .then(response => {
                if (response.status === 200) {
                    response.text().then(token => {
                        dispatch({
                            type: LOGGEDIN,
                            token
                        })

                        localStorage.setItem("token", token)

                        dispatch({ type: RENEW_RETRY })

                        return setTimeout(() => {
                            dispatch(renewToken())
                        }, renewInterval)
                    })
                } else {
                    return dispatch(renewRetry(getState().authentication.retries))
                }
            })
            .catch(ex => {
                return dispatch(renewRetry(getState().authentication.retries))
            })
    }
}

let renewRetry = tries => {
    return dispatch => {
        if (tries === 5) {
            localStorage.removeItem("token")
            dispatch({ type: RENEW_FAILED })
            dispatch(open("Failed to renew authentication token. Please sign in again", "", COLOR_DANGER))
            return dispatch(push("/login"))
        } else {
            return setTimeout(() => {
                dispatch({ type: RENEW_RETRY })
                dispatch(renewToken())
            }, tries * 1000)
        }
    }
}

export const login = () => {
    return (dispatch, getState) => {
        dispatch(close())
        dispatch({
            type: LOGIN
        })

        return fetch(`${dispatch(read(API_URL))}/auth/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(getState().authentication.form)
        })
            .then(response => {
                if (response.status === 200) {
                    response.text().then(token => {
                        dispatch({
                            type: LOGGEDIN,
                            token
                        })

                        localStorage.setItem("token", token)

                        dispatch(loadConfig(dispatch(getUserId())))

                        return setTimeout(() => {
                            dispatch(renewToken())
                        }, renewInterval)
                    })
                } else {
                    response
                        .json()
                        .then(data => {
                            dispatch({ type: ERROR })
                            dispatch(open(data.message, data.code, COLOR_DANGER))
                        })
                        .catch(ex => {
                            dispatch({ type: ERROR })
                            dispatch(open(response.statusText, response.status, COLOR_DANGER))
                        })
                }
            })
            .catch(ex => {
                dispatch({ type: ERROR })
                dispatch(open("Failed to connect to server", null, COLOR_DANGER))
            })
    }
}

export const setUsername = username => {
    return dispatch => {
        dispatch({
            type: SET_USERNAME,
            username
        })
    }
}

export const setPassword = password => {
    return dispatch => {
        dispatch({
            type: SET_PASSWORD,
            password
        })
    }
}

export const getToken = () => {
    return (dispatch, getState) => {
        return getState().authentication.tokenString
    }
}

export const getUserId = () => {
    return (dispatch, getState) => {
        return getState().authentication.token.sub
    }
}

export const getInitialState = dispatch => {
    try {
        const tokenString = localStorage.getItem("token")
        const token = jwtDecode(tokenString)
        let state = { ...initialState }

        if (token.exp - 30 > Date.now() / 1000) {
            setTimeout(() => {
                dispatch(renewToken())
            }, 5000)

            state.tokenString = tokenString
            state.token = token
        }

        return state
    } catch (e) {
        return initialState
    }
}
