import React from "react"
import { Switch, Route, NavLink, Redirect } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import { DEFAULT_WAITLIST_ID } from "shared/modules/config"
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import Spinner from "shared/containers/spinner"
import Status from "shared/containers/status"
import Alert from "shared/containers/alert"
import { RESOURCE_PATIENT_IDENTIFICATION, RESOURCE_WAITLIST, READ, WRITE, loadUserRights } from "../../modules/validations"

import Patients from "../patients"
import NewPatient from "../patients/new"
import Waitlist from "../waitlist"
import AddToWaitlist from "../waitlist/add"
import WaitlistDetail from "../patients/detail"

import { ReactComponent as PatientsIcon } from "shared/icons/patients.svg"
import { ReactComponent as WaitlistIcon } from "shared/icons/waiting-list.svg"
import { ReactComponent as LogoutIcon } from "shared/icons/logout.svg"

class App extends React.Component {
    componentDidMount() {
        if (!this.props.configLoading) {
            if (!this.props.userRights && !this.props.userRightsLoading) {
                this.props.loadUserRights()
            }
        }
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.location.pathname !== this.props.location.pathname) {
            this.props.close()
        }

        if (!nextProps.configLoading) {
            if (!nextProps.userRights && !nextProps.userRightsLoading) {
                this.props.loadUserRights()
            }
        }
    }

    logout() {
        localStorage.removeItem("token")
    }

    render() {
        if (this.props.configLoading || this.props.userRightsLoading) {
            return <Spinner />
        }

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

                    <NavLink className="navigation" to={`/waitlist/${this.props.defaultWaitlist}`}>
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
                        <Alert />
                        {this.props.canSeePatients && (<Route exact path="/" render={() => <Redirect to="/patients" />} />)}
                        {this.props.canSeePatients && (<Route exact path="/patients" component={Patients} />)}
                        {this.props.canSeePatients && (<Route path="/patients/:patientID" component={WaitlistDetail} />)}
                        {this.props.canAddPatient && (<Route exact path="/new-patient" component={NewPatient} />)}
                        {this.props.canAddToWaitlist && (<Route path="/to-waitlist/:patientID" component={AddToWaitlist} meta={{ modal: true }} />)}
                        {this.props.canSeeWaitlist && (
                            <div className="container">
                                <Route exact path="/waitlist/:waitlistID" component={Waitlist} />
                                <Switch>
                                    <Route path="/waitlist/:waitlistID/:itemID/edit-complaint" component={Waitlist} />
                                    <Route path="/waitlist/:waitlistID/:itemID/add-data" component={Waitlist} />
                                    <Route path="/waitlist/:waitlistID/:itemID" component={WaitlistDetail} />
                                </Switch>
                            </div>
                        )}
                    </div>
                </main>
            </React.Fragment>
        )
    }
}

const mapStateToProps = state => ({
    configLoading: state.config.loading,
    defaultWaitlist: state.config[DEFAULT_WAITLIST_ID],
    userRights: state.validations.userRights ? state.validations.userRights : undefined,
    userRightsLoading: state.validations.loading,
    canSeePatients: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
    canAddPatient: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[WRITE],
    canSeeWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[READ],
    canAddToWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[WRITE],
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            close,
            loadUserRights
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(App)
