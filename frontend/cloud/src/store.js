import { createStore, applyMiddleware, compose } from "redux"
import { routerMiddleware } from "react-router-redux"
import thunk from "redux-thunk"
import createHistory from "history/createBrowserHistory"
import rootReducer from "./modules"

import { getInitialState } from "shared/modules/authentication"
import { load as loadConfig } from "shared/modules/config"

export const history = createHistory()

let store

const dispatch = a => {
    store.dispatch(a)
}

const initialState = {
    authentication: getInitialState(dispatch)
}

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

loadConfig(dispatch)

export default store
