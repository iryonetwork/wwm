import React from "react"
import { Route, NavLink } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import Alert from "shared/containers/alert"
import { close } from "shared/modules/alert"
import Logo from "shared/containers/logo"
import Status from "shared/containers/status"

import Patients from "../patients"
import Waitlist from "../waitlist"

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
                        <span className="icon_menu" />
                        Patients
                    </NavLink>

                    <NavLink exact className="navigation" to="/waitlist">
                        <span className="icon_adjust-vert" />
                        Waiting list
                    </NavLink>

                    <div className="bottom">
                        <a className="navigation" href="/" onClick={this.logout}>
                            <span className="icon_compass_alt" />
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
