import React from "react"
import { Route, Link, NavLink, Switch } from "react-router-dom"
import { reduxForm } from "redux-form"

//import "./style.css"

import { joinPaths } from "shared/utils"
import { Form as PatientForm } from "../new/step1"
import { Form as FamilyForm } from "../new/step2"

const View = ({ match, location }) => (
    <div>
        <header>
            <h1>Personal Info</h1>
            <Link to={joinPaths(match.url, "edit", location.pathname.indexOf("family") !== -1 ? "family" : "")} className="btn btn-secondary btn-wide">
                Edit
            </Link>
        </header>

        <div className="navigation">
            <NavLink exact to={match.url}>
                Patient
            </NavLink>
            <NavLink to={joinPaths(match.url, "family")}>Family details</NavLink>
        </div>

        <Switch>
            <Route path={match.url + "/family"} component={ViewFamily} />
            <Route path={match.url} component={ViewPersonal} />
        </Switch>
    </div>
)

const ViewPersonal = () => (
    <div>
        <div className="section">
            <h3>Identification</h3>
            <div className="content">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">First name</div>
                        <div className="value">Alma</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Middle name</div>
                        <div className="value">Tina</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Last name</div>
                        <div className="value">Graves</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Date of birth</div>
                        <div className="value">3 June 1994</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Gender</div>
                        <div className="value">Female</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Marital status</div>
                        <div className="value">Married</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Number of kids</div>
                        <div className="value">2</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Nationality</div>
                        <div className="value">Syrian</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Country of origin</div>
                        <div className="value">Syria</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Education</div>
                        <div className="value">Secondary school</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Occupation</div>
                        <div className="value">Computer scientist</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">UN ID</div>
                        <div className="value">453ds4a56w4d8</div>
                    </div>
                </div>
            </div>
        </div>

        <div className="section">
            <h3>Contact</h3>
            <div className="content">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Country</div>
                        <div className="value">Lebanon</div>
                    </div>
                    <div className="col-sm-2">
                        <div className="label">Camp</div>
                        <div className="value">017</div>
                    </div>
                    <div className="col-sm-2">
                        <div className="label">Tent</div>
                        <div className="value">12</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Clinic</div>
                        <div className="value">CareHealth</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Phone number</div>
                        <div className="value">+963 29 2939 2919</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Email address</div>
                        <div className="value">alma@gmail.com</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Whatsapp</div>
                        <div className="value">+963 29 2939 2919</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Date of leaving home country</div>
                        <div className="value">13 November 2017</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Date of arrival</div>
                        <div className="value">2 February 2018</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
)

const ViewFamily = () => (
    <div>
        <div className="section">
            <h3>Summary</h3>
            <div className="content">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">No. of people in the family</div>
                        <div className="value">3</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">No. of people living together</div>
                        <div className="value">5</div>
                    </div>
                </div>
            </div>
        </div>

        <div className="section">
            <h3>Husband</h3>
            <div className="content">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Name</div>
                        <div className="value">
                            <Link to={`/patients/asddsa`}>Michael Graves &middot; A-</Link>
                        </div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Date of birth</div>
                        <div className="value">21 May 1986</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Living together</div>
                        <div className="value">Yes</div>
                    </div>
                </div>
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Phone number</div>
                        <div className="value">+963 29 2939 2919</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Email address</div>
                        <div className="value">alma@gmail.com</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Whatsapp</div>
                        <div className="value">+963 29 2939 2919</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">UN ID</div>
                        <div className="value">453ds4a56w4d8</div>
                    </div>
                </div>
            </div>
        </div>

        <div className="section">
            <h3>Child</h3>
            <div className="content">
                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">Name</div>
                        <div className="value">
                            <Link to={`/patients/asdsaefw`}>Michael Graves &middot; A-</Link>
                        </div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Date of birth</div>
                        <div className="value">21 May 1986</div>
                    </div>
                    <div className="col-sm-4">
                        <div className="label">Living together</div>
                        <div className="value">Yes</div>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-4">
                        <div className="label">UN ID</div>
                        <div className="value">453ds4a56w4d8</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
)

const Edit = ({ match, location }) => (
    <div>
        <header>
            <h1>Personal Info</h1>
            <Link to={location.pathname.replace("/edit", "")} className="btn btn-secondary btn-wide">
                Close
            </Link>
        </header>

        <div className="navigation">
            <NavLink exact to={match.url}>
                Patient
            </NavLink>
            <NavLink to={joinPaths(match.url, "family")}>Family details</NavLink>
        </div>

        <Switch>
            <Route path={match.url + "/family"} component={EditFamily} />
            <Route path={match.url} component={EditPersonal} />
        </Switch>
    </div>
)

let EditPersonal = () => (
    <div>
        <form>
            <PatientForm />
            <div className="section">
                <div className="row buttons">
                    <div className="col-sm-4">
                        <button type="button" className="btn btn-secondary btn-block">
                            Close
                        </button>
                    </div>
                    <div className="col-sm-4">
                        <button type="submit" className="btn btn-primary btn-block">
                            Save
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
)

EditPersonal = reduxForm({
    form: "personal",
    initialValues: {
        documents: [{}]
    }
})(EditPersonal)

let EditFamily = () => (
    <div>
        <form>
            <FamilyForm />
            <div className="section">
                <div className="row buttons">
                    <div className="col-sm-4">
                        <button type="button" className="btn btn-secondary btn-block">
                            Close
                        </button>
                    </div>
                    <div className="col-sm-4">
                        <button type="submit" className="btn btn-primary btn-block">
                            Save
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
)

EditFamily = reduxForm({
    form: "family"
})(EditFamily)

export default ({ match }) => (
    <div className="personal">
        <Switch>
            <Route path={match.url + "/edit"} render={props => <Edit {...props} closeUrl={match.url} />} />
            <Route path={match.url} component={View} />
        </Switch>
    </div>
)
