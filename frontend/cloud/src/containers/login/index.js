import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import get from "lodash/get"
import { Redirect } from "react-router-dom"

import Alert from "../../containers/alert"
import { setUsername, setPassword, login } from "../../modules/authentication"

import "./style.css"

const Login = props => {
    if (props.redirectToReferrer) {
        return <Redirect to={props.from} />
    }
    return (
        <form className="login" onSubmit={props.handleSubmit}>
            <h2>Please sign in</h2>
            <Alert />
            <label htmlFor="inputEmail" className="sr-only">
                Email address
            </label>
            <input
                type="text"
                className="form-control"
                placeholder="Username"
                onChange={props.setUsername}
                value={props.username}
                required
                autoFocus
            />
            <label htmlFor="inputPassword" className="sr-only">
                Password
            </label>
            <input
                type="password"
                className="form-control"
                placeholder="Password"
                onChange={props.setPassword}
                value={props.password}
                required
            />

            <button className="btn btn-lg btn-primary btn-block" type="submit">
                Sign in
            </button>
        </form>
    )
}

const mapStateToProps = state => ({
    username: get(state, "authentication.form.username", ""),
    password: get(state, "authentication.form.password", ""),
    pending: state.authentication.pending,
    redirectToReferrer: state.authentication.redirectToReferrer,
    from: get(state, "router.location.state.from", { pathname: "/" })
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            setUsername: e => setUsername(e.target.value),
            setPassword: e => setPassword(e.target.value),
            handleSubmit: e => {
                e.preventDefault()
                return login()
            }
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Login)
