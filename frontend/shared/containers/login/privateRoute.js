import React from "react"
import { connect } from "react-redux"
import { Redirect, Route } from "react-router-dom"
import get from "lodash/get"

const PrivateRoute = ({ component: Component, isAuthenticated, ...rest }) => (
    <Route
        {...rest}
        render={props =>
            isAuthenticated ? (
                <Component {...props} />
            ) : (
                <Redirect
                    to={{
                        pathname: "/login",
                        state: { from: props.location }
                    }}
                />
            )
        }
    />
)

const mapStateToProps = state => ({
    isAuthenticated:
        get(state, "authentication.token.exp", 0) > Date.now() / 1000
})

export default connect(mapStateToProps)(PrivateRoute)
