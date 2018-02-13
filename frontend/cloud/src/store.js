import { createStore, applyMiddleware, compose } from "redux"
import { routerMiddleware } from "react-router-redux"
import thunk from "redux-thunk"
import createHistory from "history/createBrowserHistory"
import rootReducer from "./modules"
import jwtDecode from "jwt-decode"

import { renewToken } from "./modules/authentication"

export const history = createHistory()

let store

let initialState = {
    authentication: {
        form: {},
        retries: 0
    }
}
try {
    const token = localStorage.getItem("token")
    initialState.authentication.tokenString = token
    initialState.authentication.token = jwtDecode(token)
    setTimeout(() => {
        store.dispatch(renewToken())
    }, 1000)
} catch (e) {}

const enhancers = []
const middleware = [thunk, routerMiddleware(history)]

if (process.env.NODE_ENV === "development") {
    const devToolsExtension = window.devToolsExtension

    if (typeof devToolsExtension === "function") {
        enhancers.push(devToolsExtension())
    }
}

const composedEnhancers = compose(applyMiddleware(...middleware), ...enhancers)

store = createStore(rootReducer, initialState, composedEnhancers)

export default store
