import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"
import { reducer as reduxFormReducer } from "redux-form"

import authentication from "shared/modules/authentication"
import alert from "shared/modules/alert"
import clinics from "./clinics"
import codes from "shared/modules/codes"
import config from "shared/modules/config"
import discovery from "./discovery"
import locations from "./locations"
import patient from "./patient"
import users from "./users"
import status from "shared/modules/status"
import waitlist from "./waitlist"

export default combineReducers({
    router: routerReducer,
    form: reduxFormReducer,
    authentication,
    alert,
    clinics,
    codes,
    config,
    discovery,
    locations,
    patient,
    users,
    status,
    waitlist
})
