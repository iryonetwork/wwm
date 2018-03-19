import React from "react"
import { Route, NavLink } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import Status from "shared/containers/status"

import Patients from "../patients"
import Waitlist from "../waitlist"

import { ReactComponent as PatientsIcon } from "shared/icons/patients.svg"
import { ReactComponent as WaitlistIcon } from "shared/icons/waiting-list.svg"
import { ReactComponent as LogoutIcon } from "shared/icons/logout.svg"

import "./style.css"

class App extends React.Component {
    componentWillReceiveProps(nextProps) {
        if (nextProps.location.pathname !== this.props.location.pathname) {
            this.props.close()
        }
    }

    logout() {
        localStorage.removeItem("token")
    }

    render() {
        return (
            <React.Fragment>
                <nav>
                    <div className="logo">
                        <Logo style={{ width: "100px" }} />
                        <Status />
                    </div>

                    <NavLink exact className="navigation" to="/">
                        <PatientsIcon />
                        Patients
                    </NavLink>

                    <NavLink exact className="navigation" to="/waitlist">
                        <WaitlistIcon />
                        Waiting list
                    </NavLink>

                    <div className="bottom">
                        <a className="navigation" href="/" onClick={this.logout}>
                            <LogoutIcon />
                            Logout
                        </a>
                    </div>
                </nav>
                <main>
                    <div className="container">
                        <Route exact path="/" component={Patients} />
                        <Route exact path="/waitlist" component={Waitlist} />
                    </div>
                </main>
            </React.Fragment>
        )
    }
}

const mapStateToProps = state => ({})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            close
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(App)
