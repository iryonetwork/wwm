import React from "react"
import { Route, Link, NavLink, Redirect } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import Home from "../home"
import Rules from "../rules"
import Alert from "shared/containers/alert"
import Users from "../users"
import User from "../users/detail"
import Roles from "../roles"
import UserRoles from "../userRoles"
import Locations from "../locations"
import LocationDetail from "../locations/detail"
import Organizations from "../organizations"
import Organization from "../organizations/detail"
import Clinics from "../clinics"
import { SUPERADMIN_RIGHTS_RESOURCE, ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import Reports from "../reports"
import Spinner from "shared/containers/spinner"
import Status from "shared/containers/status"

import { ReactComponent as UsersIcon } from "shared/icons/patients.svg"
import { ReactComponent as LocationsIcon } from "shared/icons/locations-active.svg"
import { ReactComponent as OrganizationsIcon } from "shared/icons/organization-active.svg"
import { ReactComponent as ClinicsIcon } from "shared/icons/doctor-active.svg"
import { ReactComponent as ReportsIcon } from "shared/icons/reports-active.svg"
import { ReactComponent as RolesIcon } from "shared/icons/role-active.svg"
import { ReactComponent as AclIcon } from "shared/icons/acl-active.svg"
import { ReactComponent as UserRolesIcon } from "shared/icons/userroles-active.svg"
import { ReactComponent as MyProfileIcon } from "shared/icons/personal-info-active.svg"
import { ReactComponent as LogoutIcon } from "shared/icons/logout.svg"

class App extends React.Component {
    componentDidMount() {
        if (!this.props.configLoading) {
            if (this.props.isAdmin === undefined || this.props.isSuperadmin === undefined) {
                this.props.loadUserRights()
            }
        }
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.location.pathname !== this.props.location.pathname) {
            this.props.close()
        }

        if (!this.props.configLoading) {
            if ((nextProps.isAdmin === undefined || nextProps.isSuperadmin === undefined) && !nextProps.validationsLoading) {
                this.props.loadUserRights()
            }
        }
    }

    logout() {
        localStorage.removeItem("token")
    }

    render() {
        if (this.props.configLoading) {
            return <Spinner />
        }

        return (
            <React.Fragment>
                <nav>
                    <div className="logo">
                        <Link to="/">
                            <Logo style={{ width: "100px" }} />
                        </Link>
                        <Status />
                    </div>

                    {this.props.isAdmin ? (
                        <div>
                            <NavLink className="navigation" to="/users">
                                <UsersIcon />
                                Users
                            </NavLink>

                            <NavLink className="navigation" to="/locations">
                                <LocationsIcon />
                                Locations
                            </NavLink>

                            <NavLink className="navigation" to="/organizations">
                                <OrganizationsIcon />
                                Organizations
                            </NavLink>

                            <NavLink className="navigation" to="/clinics">
                                <ClinicsIcon />
                                Clinics
                            </NavLink>
                            <NavLink className="navigation" to="/reports">
                                <ReportsIcon />
                                Reports
                            </NavLink>
                        </div>
                    ) : null}

                    {this.props.isSuperadmin ? (
                        <div>
                            <NavLink className="navigation" to="/roles">
                                <RolesIcon />
                                Roles
                            </NavLink>

                            <NavLink className="navigation" to="/rules">
                                <AclIcon />
                                Access Control List
                            </NavLink>

                            <NavLink className="navigation" to="/userRoles">
                                <UserRolesIcon />
                                User Roles
                            </NavLink>
                        </div>
                    ) : null}

                    <div className="bottom">
                        <NavLink className="navigation" to="/me">
                            <MyProfileIcon />
                            My Profile
                        </NavLink>
                        <a className="navigation" href="/" onClick={this.logout}>
                            <LogoutIcon />
                            Logout
                        </a>
                    </div>
                </nav>
                <main>
                    <div className="container">
                        <Alert />
                        <Route exact path="/" render={() => <Redirect to="/me" />} />
                        <Route path="/me" component={Home} />
                        <Route exact path="/users" component={Users} />
                        <Route path="/users/:userID" component={User} />
                        <Route exact path="/roles" component={Roles} />
                        <Route exact path="/roles/:roleID" component={Roles} />
                        <Route exact path="/locations" component={Locations} />
                        <Route path="/locations/:locationID" component={LocationDetail} />
                        <Route exact path="/organizations" component={Organizations} />
                        <Route path="/organizations/:organizationID" component={Organization} />
                        <Route exact path="/clinics" component={Clinics} />
                        <Route exact path="/clinics/:clinicID" component={Clinics} />
                        <Route exact path="/clinics/:clinicID/users/:userID" component={Clinics} />
                        <Route path="/userRoles" component={UserRoles} />
                        <Route exact path="/rules" component={Rules} />
                        <Route exact path="/reports" component={Reports} />
                    </div>
                </main>
            </React.Fragment>
        )
    }
}

const mapStateToProps = state => {
    return {
        isAdmin: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        isSuperadmin: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading,
        configLoading: state.config.loading
    }
}
const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUserRights,
            close
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(App)
