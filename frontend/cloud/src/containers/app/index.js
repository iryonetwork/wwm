import React from "react"
import { Route, Link, NavLink } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import Home from "../home"
import Rules from "../rules"
import Alert from "shared/containers/alert"
import Users from "../users"
import UserDetail from "../users/detail"
import Roles from "../roles"
import UserRoles from "../userRoles"
import Locations from "../locations"
import LocationDetail from "../locations/detail"
import Organizations from "../organizations"
import OrganizationDetail from "../organizations/detail"
import Clinics from "../clinics"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import { ReactComponent as LogoutIcon } from "shared/icons/logout.svg"
import { ReactComponent as MoreIcon } from "shared/icons/more.svg"

class App extends React.Component {
    componentDidMount() {
        if (this.props.isAdmin === undefined) {
            this.props.loadUserRights()
        }
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.location.pathname !== this.props.location.pathname) {
            this.props.close()
        }
        if (nextProps.isAdmin === undefined && !nextProps.validationsLoading) {
            this.props.loadUserRights()
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
                        <Link to="/">
                            <Logo style={{ width: "100px" }} />
                        </Link>
                    </div>

                    {this.props.isAdmin ? (
                        <div>
                            <NavLink className="navigation" to="/users">
                                Users
                            </NavLink>

                            <NavLink className="navigation" to="/locations">
                                Locations
                            </NavLink>

                            <NavLink className="navigation" to="/organizations">
                                Organizations
                            </NavLink>

                            <NavLink className="navigation" to="/clinics">
                                Clinics
                            </NavLink>

                            <NavLink className="navigation" to="/roles">
                                Roles
                            </NavLink>

                            <NavLink className="navigation" to="/rules">
                                ACL
                            </NavLink>

                            <NavLink className="navigation" to="/userRoles">
                                User roles
                            </NavLink>
                        </div>
                    ) : null}

                    <div className="bottom">
                        <NavLink className="navigation" to="/me">
                            <MoreIcon />
                            My profile
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
                        <Route exact path="/" component={Home} />
                        <Route exact path="/me" component={Home} />
                        <Route exact path="/users" component={Users} />
                        <Route exact path="/users/:userID" component={UserDetail} />
                        <Route exact path="/users/:userID/organizations/:organizationID" component={UserDetail} />
                        <Route exact path="/users/:userID/clinics/:clinicID" component={UserDetail} />
                        <Route path="/roles" component={Roles} />
                        <Route exact path="/locations" component={Locations} />
                        <Route path="/locations/:locationID" component={LocationDetail} />
                        <Route exact path="/organizations" component={Organizations} />
                        <Route exact path="/organizations/:organizationID" component={OrganizationDetail} />
                        <Route exact path="/organizations/:organizationID/users/:userID" component={OrganizationDetail} />
                        <Route exact path="/clinics" component={Clinics} />
                        <Route exact path="/clinics/:clinicID" component={Clinics} />
                        <Route exact path="/clinics/:clinicID/users/:userID" component={Clinics} />
                        <Route path="/userRoles" component={UserRoles} />
                        <Route exact path="/rules" component={Rules} />
                    </div>
                </main>
            </React.Fragment>
        )
    }
}

const mapStateToProps = state => {
    return {
        isAdmin: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading
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
