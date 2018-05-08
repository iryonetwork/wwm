import React from "react"
import { connect } from "react-redux"
import { Route, Link, NavLink, Switch } from "react-router-dom"
import { reduxForm } from "redux-form"

import { loadCategories, getCodes } from "shared/modules/codes"

import { joinPaths } from "shared/utils"
import Spinner from "shared/containers/spinner"
import { Form as PatientForm } from "../new/step1"
import { Form as FamilyForm } from "../new/step2"

const View = ({ patient, match, location }) => (
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

const Column = ({ value, label, codes, width }) => {
    // don't render if empty
    if (value === undefined) {
        return null
    }

    // convert for a code
    if (codes && codes.length > 0) {
        value = codes.reduce((acc, code) => {
            if (code.id === value) {
                return code.title
            }
            return acc
        }, undefined)
        // don't render if code is not found
        if (value === undefined) {
            return null
        }
    }

    return (
        <div className={`col-sm-${width}`}>
            <div className="label" key="label">
                {label}
            </div>
            <div className="value" key="value">
                {value}
            </div>
        </div>
    )
}

class ViewPersonal extends React.Component {
    constructor(props) {
        super(props)
        props.loadCategories("countries", "maritalStatus", "gender")
    }

    render() {
        const { codesLoading, patient, fetchCodes } = this.props

        if (codesLoading) {
            return <Spinner />
        }

        return (
            <div>
                <div className="section">
                    <h3>Identification</h3>
                    <div className="content">
                        <div className="row">
                            <Column width="4" label="First name" value={patient.firstName} key="firstName" />
                            <Column width="4" label="Middle name" value={patient.middleName} key="middleName" />
                            <Column width="4" label="Last name" value={patient.lastName} key="lastName" />
                        </div>

                        <div className="row">
                            <Column width="4" label="Date of birth" value={patient.dateOfBirth} key="dateOfBirth" /> {/* @TODO format date */}
                            <Column width="4" label="Gender" value={patient.gender} key="gender" codes={fetchCodes("gender")} />
                        </div>

                        <div className="row">
                            <Column width="4" label="Marital status" value={patient.maritalStatus} key="maritalStatus" codes={fetchCodes("maritalStatus")} />
                            <Column width="4" label="Number of kids" value={patient.numberOfKids} key="numberOfKids" />
                        </div>

                        <div className="row">
                            <Column width="4" label="Nationality" value={patient.nationality} key="nationality" codes={fetchCodes("countries")} />
                            <Column width="4" label="Country of origin" value={patient.countryOfOrigin} key="countryOfOrigin" codes={fetchCodes("countries")} />
                        </div>

                        <div className="row">
                            <Column width="4" label="Education" value={patient.education} key="education" codes={[]} /> {/* @TODO codes */}
                            <Column width="4" label="Occupation" value={patient.profession} key="profession" />
                        </div>

                        {patient.documents &&
                            patient.documents.length > 0 && (
                                <div className="row">{patient.documents.map((el, i) => <Column width="4" label={el.type} value={el.number} key={i} />)}</div>
                            )}
                    </div>
                </div>

                <div className="section">
                    <h3>Contact</h3>
                    <div className="content">
                        <div className="row">
                            <Column width="2" label="Country" value={patient.country} key="country" codes={fetchCodes("countries")} />
                            <Column width="2" label="Camp" value={patient.camp} key="camp" />
                            <Column width="2" label="Tent" value={patient.tent} key="tent" />
                        </div>

                        <div className="row">
                            <Column width="4" label="Phone number" value={patient.phone} key="phone" />
                            <Column width="4" label="Email address" value={patient.email} key="email" />
                            <Column width="4" label="Whatsapp" value={patient.whatsapp} key="whatsapp" />
                        </div>

                        <div className="row">
                            <Column width="4" label="Date of leaving home country" value={patient.dateOfLeaving} key="dateOfLeaving" />
                            {/* @TODO format date */}
                            <Column width="4" label="Date of arrival" value={patient.dateOfArrival} key="dateOfArrival" /> {/* @TODO format date */}
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

ViewPersonal = connect(
    state => ({
        patient: state.patient.patient,
        codesLoading: state.codes.loading
    }),
    {
        loadCategories,
        fetchCodes: getCodes
    }
)(ViewPersonal)

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
