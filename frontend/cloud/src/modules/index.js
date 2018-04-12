import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"

import authentication from "shared/modules/authentication"
import alert from "shared/modules/alert"
import users from "./users"
import roles from "./roles"
import rules from "./rules"
import userRoles from "./userRoles"
import locations from "./locations"
import organizations from "./organizations"
import clinics from "./clinics"

export default combineReducers({
    router: routerReducer,
    authentication,
    alert,
    users,
    roles,
    rules,
    userRoles,
    locations,
    organizations,
    clinics
})
