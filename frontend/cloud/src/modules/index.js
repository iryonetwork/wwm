import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"

import counter from "./counter"
import authentication from "./authentication"
import alert from "./alert"
import users from "./users"

export default combineReducers({
    router: routerReducer,
    authentication,
    counter,
    alert,
    users
})
