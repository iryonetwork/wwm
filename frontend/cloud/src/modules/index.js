import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"

import counter from "./counter"
import authentication from "./authentication"

export default combineReducers({
    router: routerReducer,
    authentication,
    counter
})
