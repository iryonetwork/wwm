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
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
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
                        <Link to="/">
                            <Logo style={{ width: "100px" }} />
                        </Link>
                    </div>

                    <NavLink className="navigation" to="/users">
                        Users
                    </NavLink>

                    <NavLink className="navigation" to="/roles">
                        Roles
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

                    <NavLink className="navigation" to="/userRoles">
                        User roles
                    </NavLink>

                    <NavLink className="navigation" to="/rules">
                        ACL
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
                        <Route exact path="/" component={Home} />
                        <Route exact path="/users" component={Users} />
                        <Route path="/users/:id" component={UserDetail} />
                        <Route path="/roles" component={Roles} />
                        <Route exact path="/locations" component={Locations} />
                        <Route path="/locations/:id" component={LocationDetail} />
                        <Route exact path="/organizations" component={Organizations} />
                        <Route path="/organizations/:id" component={OrganizationDetail} />
                        <Route exact path="/clinics" component={Clinics} />
                        <Route path="/userRoles" component={UserRoles} />
                        <Route exact path="/rules" component={Rules} />
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
