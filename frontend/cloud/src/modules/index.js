import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"

import authentication from "shared/modules/authentication"
import alert from "shared/modules/alert"
import config from "shared/modules/config"
import status from "shared/modules/status"
import users from "./users"
import roles from "./roles"
import rules from "./rules"
import userRoles from "./userRoles"
import locations from "./locations"
import organizations from "./organizations"
import clinics from "./clinics"
import codes from "./codes"
import validations from "./validations"
import reports from "./reports"

export default combineReducers({
    router: routerReducer,
    authentication,
    alert,
    config,
    status,
    users,
    roles,
    rules,
    userRoles,
    locations,
    organizations,
    clinics,
    codes,
    validations,
    reports
})
