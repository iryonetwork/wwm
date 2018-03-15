import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import get from "lodash/get"
import { Redirect } from "react-router-dom"

import Alert from "../../containers/alert"
import { setUsername, setPassword, login } from "../../modules/authentication"

import Logo from "../logo"
import Status from "../status"

import "./style.css"

const Login = props => {
    if (props.redirectToReferrer) {
        return <Redirect to={props.from} />
    }
    return (
        <div className="login-container">
            <div className="logo">
                <Logo style={{ width: "200px" }} />
                <Status />
            </div>
            <form className="login" onSubmit={props.handleSubmit}>
                <div className="form-group">
                    <label htmlFor="inputEmail" className="sr-only">
                        Email address
                    </label>
                    <input
                        type="text"
                        className={"form-control" + (props.error ? " is-invalid" : "")}
                        placeholder="Username"
                        onChange={props.setUsername}
                        value={props.username}
                        required
                        autoFocus
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="inputPassword" className="sr-only">
                        Password
                    </label>
                    <input
                        type="password"
                        className={"form-control" + (props.error ? " is-invalid" : "")}
                        placeholder="Password"
                        onChange={props.setPassword}
                        value={props.password}
                        required
                    />
                </div>
                <Alert disableClose={true} />
                <button className="btn btn-primary btn-block" type="submit">
                    Sign in
                </button>
            </form>
            <div className="footer" />
        </div>
    )
}

const mapStateToProps = state => ({
    username: get(state, "authentication.form.username", ""),
    password: get(state, "authentication.form.password", ""),
    pending: state.authentication.pending,
    error: state.authentication.error,
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
