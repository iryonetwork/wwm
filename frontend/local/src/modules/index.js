import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"

import authentication from "shared/modules/authentication"
import alert from "shared/modules/alert"

export default combineReducers({
    router: routerReducer,
    authentication,
    alert
})
