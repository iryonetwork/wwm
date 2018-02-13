import React from "react"
import { Route, Link, NavLink } from "react-router-dom"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"

import Home from "../home"
import Rules from "../rules"
import Alert from "../alert"
import Users from "../users"
import UserDetail from "../users/detail"
import Roles from "../roles"
import { close } from "../../modules/alert"

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
            <div>
                <header>
                    <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
                        <div className="container">
                            <Link className="navbar-brand" to="/">
                                Home
                            </Link>
                            <button
                                className="navbar-toggler"
                                type="button"
                                data-toggle="collapse"
                                data-target="#navbarText"
                                aria-controls="navbarText"
                                aria-expanded="false"
                                aria-label="Toggle navigation"
                            >
                                <span className="navbar-toggler-icon" />
                            </button>

                            <div className="collapse navbar-collapse">
                                <ul className="navbar-nav mr-auto">
                                    <li className="nav-item">
                                        <NavLink className="nav-link" to="/users">
                                            Users
                                        </NavLink>
                                    </li>
                                    <li className="nav-item">
                                        <NavLink className="nav-link" to="/roles">
                                            Roles
                                        </NavLink>
                                    </li>
                                    <li className="nav-item">
                                        <NavLink className="nav-link" to="/rules">
                                            ACL
                                        </NavLink>
                                    </li>
                                </ul>

                                <ul className="navbar-nav">
                                    <li className="nav-item">
                                        <a href="/" className="nav-link" onClick={this.logout}>
                                            Logout
                                        </a>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </nav>
                </header>

                <main className="container">
                    <Alert />
                    <Route exact path="/" component={Home} />
                    <Route exact path="/users" component={Users} />
                    <Route path="/users/:id" component={UserDetail} />
                    <Route path="/roles" component={Roles} />
                    <Route exact path="/rules" component={Rules} />
                </main>
            </div>
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
