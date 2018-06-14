import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { Route, Link, NavLink, Switch } from "react-router-dom"
import { reduxForm } from "redux-form"

import { RESOURCE_DEMOGRAPHIC_INFORMATION, READ, UPDATE } from "../../../modules/validations"
import { loadCategories, getCodes, getCodesAsOptions } from "shared/modules/codes"
import { updatePatient } from "../../../modules/patient"
import validate from "../shared/validatePersonal"
import { relationOptions, livingTogetherOptions } from "shared/forms/options"

import { joinPaths } from "shared/utils"
import Spinner from "shared/containers/spinner"
import { Form as PatientForm } from "../new/step1"
import { Form as FamilyForm } from "../new/step2"

let View = ({ patient, match, location, canSeeDemographicInformation, canEditDemographicInformation }) => {
    return canSeeDemographicInformation ? (
        <div>
            <header>
                <h1>Personal Info</h1>
                {canEditDemographicInformation && (
                    <div>
                        <Link
                            to={joinPaths(match.url, "edit", location.pathname.indexOf("family") !== -1 ? "family" : "")}
                            className="btn btn-secondary btn-wide"
                        >
                            Edit
                        </Link>
                    </div>
                )}
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
    ) : null
}

View = connect(
    state => ({
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ],
        canEditDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[UPDATE]
    }),
    {}
)(View)

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
        const { codesLoading, patient, fetchCodes, canSeeDemographicInformation } = this.props

        if (codesLoading) {
            return <Spinner />
        }

        return canSeeDemographicInformation ? (
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
                            <Column width="4" label="Country" value={patient.country} key="country" codes={fetchCodes("countries")} />
                            <Column width="4" label="Region" value={patient.region} key="region" />
                            <Column width="4" label="Address" value={patient.address} key="address" />
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
        ) : null
    }
}

ViewPersonal = connect(
    state => ({
        patient: state.patient.patient,
        codesLoading: state.codes.loading,
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ]
    }),
    {
        loadCategories,
        fetchCodes: getCodes
    }
)(ViewPersonal)

let ViewFamily = ({ patient, canSeeDemographicInformation }) => {
    return canSeeDemographicInformation ? (
        patient.peopleInFamily || patient.peopleLivingTogether || (patient.familyMembers && patient.familyMembers.length > 0) ? (
            <div>
                <div className="section">
                    <h3>Summary</h3>
                    <div className="content">
                        <div className="row">
                            <Column width="4" label="No. of people in the family" value={patient.peopleInFamily} key="peopleInFamily" />
                            <Column width="4" label="No. of people living together" value={patient.peopleLivingTogether} key="peopleLivingTogether" />
                        </div>
                    </div>
                </div>

                {(patient.familyMembers || []).map(member => (
                    <div className="section" key={member.patientID}>
                        <h3>{(_.find(relationOptions, { value: member.relation }) || { label: member.relation }).label}</h3>
                        <div className="content">
                            <div className="row">
                                <div className="col-sm-4">
                                    <div className="label">Name</div>
                                    <div className="value">
                                        <Link to={`/patients/${member.patientID}/personal`}>
                                            {member.lastName}, {member.firstName}
                                        </Link>
                                    </div>
                                </div>
                                <Column width="4" label="Date of birth" value={member.dateOfBirth} />
                                <Column
                                    width="4"
                                    label="Living together"
                                    value={(_.find(livingTogetherOptions, { value: member.livingTogether }) || { label: member.livingTogether }).label}
                                />
                            </div>

                            {member.documents &&
                                member.documents.length > 0 && (
                                    <div className="row">{member.documents.map((el, i) => <Column width="4" label={el.type} value={el.number} key={i} />)}</div>
                                )}
                        </div>
                    </div>
                ))}
            </div>
        ) : (
            <h3>No information about family was found.</h3>
        )
    ) : null
}

ViewFamily = connect(
    state => ({
        patient: state.patient.patient,
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ]
    }),
    {}
)(ViewFamily)

class EditPersonal extends React.Component {
    constructor(props) {
        super(props)

        this.update = this.update.bind(this)
    }

    update(form) {
        this.props.updatePatient(form).then(() => {
            this.props.history.push(this.props.location.pathname.replace("/edit", ""))
        })
    }

    componentWillMount() {
        this.props.loadCategories("gender", "maritalStatus", "countries", "documentTypes")
    }

    render() {
        let { match, codesLoading, getCodes, handleSubmit, updating, location } = this.props

        if (codesLoading && !updating) {
            return <Spinner />
        }

        return this.props.canEditDemographicInformation ? (
            <div>
                <header>
                    <h1>Personal Info</h1>
                    <Link to={location.pathname.replace("/edit", "")} className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={handleSubmit(this.update)} className="btn btn-primary btn-wide">
                        {updating ? "Saving..." : "Save"}
                    </button>
                </header>

                <div className="navigation">
                    <NavLink exact to={match.url}>
                        Patient
                    </NavLink>
                    <NavLink to={joinPaths(match.url, "family")}>Family details</NavLink>
                </div>

                <div>
                    <form onSubmit={handleSubmit(this.update)}>
                        <PatientForm
                            countries={getCodes("countries")}
                            maritalStatus={getCodes("maritalStatus")}
                            genders={getCodes("gender")}
                            documentTypes={getCodes("documentTypes")}
                        />
                        <div className="section">
                            <div className="row buttons">
                                <div className="col-sm-4">
                                    <Link to={location.pathname.replace("/edit", "")} className="btn btn-secondary btn-block">
                                        Cancel
                                    </Link>
                                </div>
                                <div className="col-sm-4">
                                    <button type="submit" className="btn btn-primary btn-block" disabled={updating}>
                                        {updating ? "Saving..." : "Save"}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        ) : null
    }
}

EditPersonal = reduxForm({
    form: "personal",
    initialValues: {
        documents: [{}]
    },
    validate
})(EditPersonal)

EditPersonal = connect(
    state => ({
        codesLoading: state.codes.loading,
        initialValues: state.patient.patient,
        updating: state.patient.updating,
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ],
        canEditDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[UPDATE]
    }),
    {
        getCodes: getCodesAsOptions,
        loadCategories,
        updatePatient
    }
)(EditPersonal)

class EditFamily extends React.Component {
    constructor(props) {
        super(props)

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    handleSubmit(form) {
        this.props.updatePatient(form).then(() => {
            this.props.history.push(this.props.location.pathname.replace("/edit", ""))
        })
    }

    render() {
        let { match, location, updating, handleSubmit } = this.props
        return this.props.canEditDemographicInformation ? (
            <div>
                <header>
                    <h1>Personal Info</h1>
                    <Link to={location.pathname.replace("/edit", "")} className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={handleSubmit(this.handleSubmit)} className="btn btn-primary btn-wide">
                        {updating ? "Saving..." : "Save"}
                    </button>
                </header>

                <div className="navigation">
                    <NavLink exact to={location.pathname.replace("/family", "")}>
                        Patient
                    </NavLink>
                    <NavLink exact to={match.url}>
                        Family details
                    </NavLink>
                </div>

                <div>
                    <form onSubmit={handleSubmit(this.handleSubmit)}>
                        <FamilyForm />
                        <div className="section">
                            <div className="row buttons">
                                <div className="col-sm-4">
                                    <Link to={location.pathname.replace("/edit", "")} className="btn btn-secondary btn-block">
                                        Cancel
                                    </Link>
                                </div>
                                <div className="col-sm-4">
                                    <button type="submit" className="btn btn-primary btn-block" disabled={updating}>
                                        {updating ? "Saving..." : "Save"}
                                    </button>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        ) : null
    }
}

EditFamily = reduxForm({
    form: "family"
})(EditFamily)

EditFamily = connect(
    state => ({
        initialValues: state.patient.patient,
        updating: state.patient.updating,
        canEditDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[UPDATE]
    }),
    {
        updatePatient
    }
)(EditFamily)

let PersonalInfoRoutes = ({ match, canSeeDemographicInformation, canEditDemographicInformation }) => (
    <div className="personal">
        <Switch>
            {canEditDemographicInformation && <Route exact path={match.url + "/edit"} component={EditPersonal} />}
            {canEditDemographicInformation && <Route exact path={match.url + "/edit/family"} component={EditFamily} />}
            {canSeeDemographicInformation && <Route path={match.url} component={View} />}
        </Switch>
    </div>
)

export default connect(
    state => ({
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ],
        canEditDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[UPDATE]
    }),
    {}
)(PersonalInfoRoutes)
