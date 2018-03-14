import React from "react"
import { render } from "react-dom"
import { Provider } from "react-redux"
import { ConnectedRouter } from "react-router-redux"
import { Route, Switch } from "react-router-dom"

import store, { history } from "./store"
import App from "./containers/app"
import Login from "shared/containers/login"
import registerServiceWorker from "shared/registerServiceWorker"
import PrivateRoute from "shared/containers/login/privateRoute"

import "shared/styles/index.css"

const target = document.querySelector("#root")

render(
    <Provider store={store}>
        <ConnectedRouter history={history}>
            <Switch>
                <Route exact path="/login" component={Login} />
                <PrivateRoute path="/" component={App} />
            </Switch>
        </ConnectedRouter>
    </Provider>,
    target
)

registerServiceWorker()
