import { combineReducers } from "redux"
import { routerReducer } from "react-router-redux"
import { reducer as reduxFormReducer } from "redux-form"

import authentication from "shared/modules/authentication"
import alert from "shared/modules/alert"

export default combineReducers({
    router: routerReducer,
    form: reduxFormReducer,
    authentication,
    alert
})
