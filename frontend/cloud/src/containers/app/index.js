import React from "react"
import { Route, Link, NavLink } from "react-router-dom"
import Home from "../home"
import About from "../about"

const App = () => (
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
                                <NavLink exact className="nav-link" to="/">
                                    Home
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink
                                    exact
                                    className="nav-link"
                                    to="/about-us"
                                >
                                    About Us
                                </NavLink>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
        </header>

        <main className="container">
            <Route exact path="/" component={Home} />
            <Route exact path="/about-us" component={About} />
        </main>
    </div>
)

export default App
