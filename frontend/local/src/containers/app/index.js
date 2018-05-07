import React from "react"
import { Route, NavLink, Redirect } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import { read, DEFAULT_WAITLIST_ID } from "shared/modules/config"
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import Status from "shared/containers/status"

import Patients from "../patients"
import NewPatient from "../patients/new"
import Waitlist from "../waitlist"
import AddToWaitlist from "../waitlist/add"
import WaitlistDetail from "../patients/detail"

import { ReactComponent as PatientsIcon } from "shared/icons/patients.svg"
import { ReactComponent as WaitlistIcon } from "shared/icons/waiting-list.svg"
import { ReactComponent as LogoutIcon } from "shared/icons/logout.svg"

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

                    <NavLink className="navigation" to="/patients">
                        <PatientsIcon />
                        Patients
                    </NavLink>

                    <NavLink className="navigation" to={`/waitlist/${read(DEFAULT_WAITLIST_ID)}`}>
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
                        <Route exact path="/" render={() => <Redirect to="/patients" />} />
                        <Route exact path="/patients" component={Patients} />
                        <Route exact path="/new-patient" component={NewPatient} />
                        <Route path="/to-waitlist/:patientID" component={AddToWaitlist} meta={{ modal: true }} />
                        <Route exact path="/waitlist/:waitlistID" component={Waitlist} />
                    </div>
                    <Route path="/waitlist/:waitlistID/:itemID" component={WaitlistDetail} />
                    <Route path="/patients/:patientID" component={WaitlistDetail} />
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
