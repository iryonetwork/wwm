import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter, Route, Link, NavLink, Switch } from "react-router-dom"

import { push } from "react-router-redux"
import _ from "lodash"
import classnames from "classnames"

import { joinPaths } from "shared/utils"
import { loadUser, saveUser } from "../../modules/users"
import { CATEGORY_COUNTRIES, CATEGORY_LANGUAGES, CATEGORY_LICENSES, loadCodes } from "../../modules/codes"
import { SELF_RIGHTS_RESOURCE, ADMIN_RIGHTS_RESOURCE, SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import OrganizationsList from "./organizationsList"
import ClinicsList from "./clinicsList"
import WildcardUserRolesList from "./wildcardUserRolesList"
import { processStateOnChange, processStateOnBlur } from "../../utils/formFieldsUpdate"
import Spinner from "shared/containers/spinner"

import { ReactComponent as RemoveIcon } from "shared/icons/negative.svg"

import "../../styles/style.css"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            username: "",
            email: "",
            password: "",
            password2: "",
            personalData: {
                passport: {}
            },
            loading: true,
            validationErrors: {}
        }
    }

    componentDidMount() {
        if (!this.props.user && this.props.userID !== "new") {
            this.props.loadUser(this.props.userID)
        }
        if (!this.props.countries) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }
        if (!this.props.languages) {
            this.props.loadCodes(CATEGORY_LANGUAGES)
        }
        if (!this.props.licenses) {
            this.props.loadCodes(CATEGORY_LICENSES)
        }
        if (
            this.props.canSeePersonal === undefined ||
            this.props.canSeeOrganizations === undefined ||
            this.props.canSeeClinics === undefined ||
            this.props.canSeeWildcardUserRoles === undefined ||
            this.props.canEditPersonal === undefined
        ) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.user && nextProps.userID !== "new" && !nextProps.usersLoading) {
            this.props.loadUser(nextProps.userID)
        }
        if (!nextProps.countries && !nextProps.codesLoading) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }
        if (!nextProps.languages && !nextProps.codesLoading) {
            this.props.loadCodes(CATEGORY_LANGUAGES)
        }
        if (!nextProps.licenses && !nextProps.codesLoading) {
            this.props.loadCodes(CATEGORY_LICENSES)
        }
        if (
            (nextProps.canSeePersonal === undefined ||
                nextProps.canSeeOrganizations === undefined ||
                nextProps.canSeeClinics === undefined ||
                nextProps.canSeeWildcardUserRoles === undefined ||
                nextProps.canEditPersonal === undefined) &&
            !nextProps.validationsLoading
        ) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            (!props.user && props.userID !== "new") ||
            props.usersLoading ||
            props.canEditPersonal === undefined ||
            props.canSeePersonal === undefined ||
            props.validationsLoading ||
            !props.countries ||
            !props.languages ||
            !props.licenses ||
            props.codesLoading
        this.setState({ loading: loading })

        if (props.user) {
            let personalData = _.clone(props.user.personalData) || {}

            personalData.passport = _.clone(personalData.passport) || {}
            personalData.languages = _.clone(personalData.languages) || []
            personalData.licenses = _.clone(personalData.licenses) || []

            // format languages
            if (personalData && personalData.languages) {
                personalData.languages = _.map(personalData.languages, languageCodeID => {
                    return { id: languageCodeID }
                })
            }
            // format licenses
            if (personalData && personalData.licenses) {
                personalData.licenses = _.map(personalData.licenses, licenseCodeID => {
                    return { id: licenseCodeID }
                })
            }

            this.setState({
                email: props.user.email,
                personalData: personalData || {}
            })
        }
    }

    newLanguage() {
        return e => {
            let personalData = this.state.personalData
            if (personalData.languages) {
                personalData.languages = [...personalData.languages, { id: undefined, edit: true }]
            } else {
                personalData.languages = [{ id: undefined, edit: true }]
            }
            this.setState({ personalData: personalData })
        }
    }

    updateLanguage(index) {
        return e => {
            let personalData = this.state.personalData
            if (personalData.languages) {
                personalData.languages[index].id = e.target.value
            }
            this.setState({ personalData: personalData })
        }
    }

    removeLanguage(index) {
        return e => {
            let personalData = this.state.personalData
            if (personalData.languages) {
                personalData.languages.splice(index, 1)
            }
            this.setState({ personalData: personalData })
        }
    }

    newLicense() {
        return e => {
            let personalData = this.state.personalData
            if (personalData.licenses) {
                personalData.licenses = [...personalData.licenses, { id: undefined, edit: true }]
            } else {
                personalData.licenses = [{ id: undefined, edit: true }]
            }
            this.setState({ personalData: personalData })
        }
    }

    updateLicense(index) {
        return e => {
            let personalData = this.state.personalData
            if (personalData.licenses) {
                personalData.licenses[index].id = e.target.value
            }
            this.setState({ personalData: personalData })
        }
    }

    removeLicense(index) {
        return e => {
            let personalData = this.state.personalData
            if (personalData.licenses) {
                personalData.licenses.splice(index, 1)
            }
            this.setState({ personalData: personalData })
        }
    }

    updateInput(e) {
        return e => {
            this.setState(processStateOnChange(this.state, e))
        }
    }

    onBlurInput(e) {
        return e => {
            this.setState(processStateOnBlur(this.state, e))
        }
    }

    submit(e) {
        return e => {
            e.preventDefault()
            this.props.close()

            let validationErrors = {}

            if (!this.props.user && (!this.state.username || this.state.username === "")) {
                validationErrors["username"] = "Required"
            }

            if (!this.state.email || this.state.email === "") {
                validationErrors["email"] = "Required"
            }

            if (!this.props.user && this.state.password === "") {
                validationErrors["password"] = "Required"
            }

            if (!this.props.user && this.state.password2 === "") {
                validationErrors["password2"] = "Required"
            }

            if (this.state.password !== this.state.password2) {
                validationErrors["password"] = "Passwords don't match"
                validationErrors["password2"] = "Passwords don't match"
            }

            if (!this.state.personalData.firstName || this.state.personalData.firstName === "") {
                validationErrors["personalData.firstName"] = "Required"
            }
            if (!this.state.personalData.lastName || this.state.personalData.lastName === "") {
                validationErrors["personalData.lastName"] = "Required"
            }
            if (!this.state.personalData.dateOfBirth || this.state.personalData.dateOfBirth === "" || this.state.personalData.dateOfBirth.length !== 10) {
                validationErrors["personalData.dateOfBirth"] = "Required"
            }

            let user = {
                email: this.state.email
            }

            if (this.props.user) {
                user.id = this.props.user.id
                user.username = this.props.user.username
            } else {
                user.username = this.state.username
            }
            if (this.state.password && this.state.password.trim() !== "") {
                user.password = this.state.password
            }

            user.personalData = _.clone(this.state.personalData)

            // format languages
            if (user.personalData.languages && user.personalData.languages.length !== 0) {
                user.personalData.languages = _.map(
                    _.pickBy(user.personalData.languages, language => language.id && language.id !== ""),
                    language => language.id
                )
            }

            // format licenses
            if (user.personalData.licenses && user.personalData.licenses.length !== 0) {
                user.personalData.licenses = _.map(_.pickBy(user.personalData.licenses, license => license.id && license.id !== ""), license => license.id)
            }

            if (!_.isEmpty(validationErrors)) {
                this.props.open("There are errors in the data submitted", "", COLOR_DANGER)
                this.setState({ validationErrors: validationErrors })
                return
            }

            this.props.saveUser(user).then(response => {
                if (!user.id && response && response.id) {
                    this.props.push(`/users/${response.id}`)
                }
            })
        }
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <Spinner />
        }
        if (!props.canSeePersonal || props.forbidden) {
            return null
        }

        let basePath = props.home ? "/me" : `/users/${props.userID}`

        return (
            <div>
                <header>
                    {props.isSelf ? <h1>My Profile</h1> : <h1>Users</h1>}
                    <Link to="/users" className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={this.submit()} className="btn btn-primary btn-wide">
                        {props.usersUpdating ? "Saving..." : "Save"}
                    </button>
                </header>
                <h2>{props.user ? props.user.username : "New User"}</h2>
                {props.user ? (
                    <div className="navigation">
                        {props.canSeePersonal ? (
                            <NavLink exact to={basePath}>
                                Personal Info
                            </NavLink>
                        ) : null}
                        {props.canSeeOrganizations ? <NavLink to={joinPaths(basePath, "organizations")}>Organizations</NavLink> : null}
                        {props.canSeeClinics ? <NavLink to={joinPaths(basePath, "clinics")}>Clinics</NavLink> : null}
                        {props.canSeeWildcardUserRoles ? <NavLink to={joinPaths(basePath, "userroles")}>Wildcard Roles</NavLink> : null}
                    </div>
                ) : null}
                <div className="user-form">
                    <form onSubmit={this.submit()} className="needs-validation" noValidate>
                        <div>
                            {props.user ? null : (
                                <div className="section">
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className={"form-control" + (this.state.validationErrors["username"] ? " is-invalid" : "")}
                                                    id="username"
                                                    value={this.state.username || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="Username"
                                                    required="true"
                                                />
                                                <span>Username</span>
                                                {this.state.validationErrors["username"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["username"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )}
                                            </label>
                                        </div>
                                    </div>
                                </div>
                            )}
                            <div className="section">
                                <h3>Personal Data</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className={"form-control" + (this.state.validationErrors["personalData.firstName"] ? " is-invalid" : "")}
                                                    id="personalData.firstName"
                                                    value={this.state.personalData.firstName || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="First Name"
                                                    required="true"
                                                />
                                                <span>First name</span>
                                                {this.state.validationErrors["personalData.firstName"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.firstName"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )}
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="personalData.middleName"
                                                    value={this.state.personalData.middleName || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="Middle Name"
                                                />
                                                <span>Middle Name</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className={"form-control" + (this.state.validationErrors["personalData.lastName"] ? " is-invalid" : "")}
                                                    id="personalData.lastName"
                                                    value={this.state.personalData.lastName || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="Last Name"
                                                    required="true"
                                                />
                                                <span>Last Name</span>
                                                {this.state.validationErrors["personalData.lastName"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.lastName"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )}
                                            </label>
                                        </div>
                                    </div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="date"
                                                    id="personalData.dateOfBirth"
                                                    className={"form-control" + (this.state.validationErrors["personalData.dateOfBirth"] ? " is-invalid" : "")}
                                                    value={this.state.personalData.dateOfBirth || ""}
                                                    placeholder="Date of birth"
                                                    disabled={!props.canEditPersonal}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    required="true"
                                                />
                                                <span>Date of birth</span>
                                                {this.state.validationErrors["personalData.dateOfBirth"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.dateOfBirth"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )}
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <select
                                                    className={classnames("form-control", { selected: this.state.personalData.nationality })}
                                                    id="personalData.nationality"
                                                    value={this.state.personalData.nationality || ""}
                                                    onChange={this.updateInput()}
                                                    disabled={!props.canEditPersonal}
                                                >
                                                    <option value="" disabled>
                                                        Nationality
                                                    </option>
                                                    {_.map(props.countries, country => (
                                                        <option key={country.id} value={country.id}>
                                                            {country.title}
                                                        </option>
                                                    ))}
                                                </select>
                                                <span>Nationality</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <select
                                                    className={classnames("form-control", { selected: this.state.personalData.residency })}
                                                    id="personalData.residency"
                                                    value={this.state.personalData.residency || ""}
                                                    onChange={this.updateInput()}
                                                    disabled={!props.canEditPersonal}
                                                >
                                                    <option value="" disabled>
                                                        Residency
                                                    </option>
                                                    {_.map(props.countries, country => (
                                                        <option key={country.id} value={country.id}>
                                                            {country.title}
                                                        </option>
                                                    ))}
                                                </select>
                                                <span>Residency</span>
                                            </label>
                                        </div>
                                    </div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-12">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="personalData.specialisation"
                                                    value={this.state.personalData.specialisation || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="Medical Worker Specialisation"
                                                />
                                                <span>Medical Worker Specialisation</span>
                                            </label>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="section">
                                <h3>Languages</h3>
                                <div>
                                    {_.map(this.state.personalData.languages ? this.state.personalData.languages : [], (language, i) => (
                                        <div className="form-row no-label" key={`personalData.languages[${i}]`}>
                                            <div className="form-group col-sm-6">
                                                <select
                                                    className={classnames("form-control", { disabled: !props.canEditPersonal, selected: language.id })}
                                                    id={`personalData.languages[${i}]`}
                                                    value={language.id || ""}
                                                    onChange={this.updateLanguage(i)}
                                                    disabled={!props.canEditPersonal}
                                                >
                                                    <option value="" disabled>
                                                        Select language
                                                    </option>
                                                    {_.map(
                                                        _.difference(
                                                            _.map(props.languages, language => language.id),
                                                            _.without(_.map(this.state.personalData.languages, language => language.id), language.id)
                                                        ),
                                                        languageCodeID => (
                                                            <option key={languageCodeID} value={languageCodeID}>
                                                                {props.languages[languageCodeID].title}
                                                            </option>
                                                        )
                                                    )}
                                                </select>
                                            </div>
                                            <div className="form-group col-sm-6">
                                                {props.canEditPersonal ? (
                                                    <button type="button" onClick={this.removeLanguage(i)} className="btn remove">
                                                        <RemoveIcon />
                                                    </button>
                                                ) : null}
                                            </div>
                                        </div>
                                    ))}
                                    {props.canEditPersonal ? (
                                        <button type="button" className="btn btn-link" onClick={this.newLanguage()}>
                                            Add language
                                        </button>
                                    ) : null}
                                </div>
                            </div>
                            <div className="section">
                                <h3>Licenses</h3>
                                <div>
                                    {_.map(this.state.personalData.licenses ? this.state.personalData.licenses : [], (license, i) => (
                                        <div className="form-row no-label" key={`personalData.licenses[${i}]`}>
                                            <div className="form-group col-sm-6">
                                                <select
                                                    className={classnames("form-control", { disabled: !props.canEditPersonal, selected: license.id })}
                                                    id={`personalData.licenses[${i}]`}
                                                    value={license.id || ""}
                                                    onChange={this.updateLicense(i)}
                                                    disabled={!props.canEditPersonal}
                                                >
                                                    <option value="">Select license</option>
                                                    {_.map(
                                                        _.difference(
                                                            _.map(props.licenses, license => license.id),
                                                            _.without(_.map(this.state.personalData.licenses, license => license.id), license.id)
                                                        ),
                                                        licenseCodeID => (
                                                            <option key={licenseCodeID} value={licenseCodeID}>
                                                                {props.licenses[licenseCodeID].title}
                                                            </option>
                                                        )
                                                    )}
                                                </select>
                                            </div>
                                            {props.canEditPersonal ? (
                                                <div className="form-group col-sm-6">
                                                    <button type="button" onClick={this.removeLicense(i)} className="btn remove">
                                                        <RemoveIcon />
                                                    </button>
                                                </div>
                                            ) : null}
                                        </div>
                                    ))}
                                    {props.canEditPersonal ? (
                                        <button type="button" className="btn btn-link" onClick={this.newLicense()}>
                                            Add License
                                        </button>
                                    ) : null}
                                </div>
                            </div>
                            {!props.user || props.canEditPassword ? (
                                <div className="section">
                                    <h3>Security</h3>
                                    <div>
                                        <div className="form-row">
                                            <div className="form-group col-sm-4">
                                                <label>
                                                    <input
                                                        type="password"
                                                        className={"form-control" + (this.state.validationErrors["password"] ? " is-invalid" : "")}
                                                        id="password"
                                                        value={this.state.password || ""}
                                                        onChange={this.updateInput()}
                                                        onBlur={this.onBlurInput()}
                                                        disabled={!props.canEditPersonal}
                                                        placeholder={props.user ? "●●●●●" : "Enter password"}
                                                        required={props.user ? null : "true"}
                                                    />
                                                    <span>{props.user ? "Enter new password" : "Enter password"}</span>
                                                    <div className="invalid-feedback">{this.state.validationErrors["password"]}</div>
                                                    {!props.user ? <small className="form-text text-muted">Required</small> : null}
                                                </label>
                                            </div>
                                        </div>
                                        <div className="form-row">
                                            <div className="form-group col-sm-4">
                                                <label>
                                                    <input
                                                        type="password"
                                                        className={"form-control" + (this.state.validationErrors["password2"] ? " is-invalid" : "")}
                                                        id="password2"
                                                        value={this.state.password2 || ""}
                                                        onChange={this.updateInput()}
                                                        disabled={!props.canEditPersonal}
                                                        placeholder={props.user ? "●●●●●" : "Enter password again"}
                                                        required={props.user ? null : "true"}
                                                    />
                                                    <span>{props.user ? "Enter new password again" : "Enter password again"}</span>
                                                    <div className="invalid-feedback">{this.state.validationErrors["password2"]}</div>
                                                    {!props.user ? <small className="form-text text-muted">Required</small> : null}
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            ) : null}
                            <div className="section">
                                <h3>Contact Information</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="email"
                                                    className={"form-control" + (this.state.validationErrors["email"] ? " is-invalid" : "")}
                                                    id="email"
                                                    value={this.state.email || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="E-mail address"
                                                    required="true"
                                                />
                                                <span>E-mail address</span>
                                                {this.state.validationErrors["email"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["email"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )}
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="personalData.phoneNumber"
                                                    value={this.state.personalData.phoneNumber || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="Phone Number"
                                                />
                                                <span>Phone Number</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="personalData.whatsApp"
                                                    value={this.state.personalData.whatsApp || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEditPersonal}
                                                    placeholder="WhatsApp"
                                                />
                                                <span>WhatsApp</span>
                                            </label>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="section">
                                <h3>Passport</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="personalData.passport.number"
                                                    value={this.state.personalData.passport.number || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    placeholder="Number"
                                                    disabled={!props.canEditPersonal}
                                                />
                                                <span>Number</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <select
                                                    className={classnames("form-control", { selected: this.state.personalData.passport.issuingCountry })}
                                                    id="personalData.passport.issuingCountry"
                                                    value={this.state.personalData.passport.issuingCountry || ""}
                                                    onChange={this.updateInput()}
                                                    disabled={!props.canEditPersonal}
                                                >
                                                    <option value="" disabled>
                                                        Issuing Country
                                                    </option>
                                                    {_.map(props.countries, country => (
                                                        <option key={country.id} value={country.id}>
                                                            {country.title}
                                                        </option>
                                                    ))}
                                                </select>
                                                <span>Issuing Country</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="date"
                                                    className={
                                                        "form-control" + (this.state.validationErrors["personalData.passport.expiryDate"] ? " is-invalid" : "")
                                                    }
                                                    id="personalData.passport.expiryDate"
                                                    value={this.state.personalData.passport.expiryDate || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    placeholder="Expiry Date"
                                                    disabled={!props.canEditPersonal}
                                                />
                                                <div className="invalid-feedback">{this.state.validationErrors["personalData.passport.expiryDate"]}</div>
                                                <span>Expiry Date</span>
                                            </label>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        {props.canEditPersonal ? (
                            <div className="row buttons">
                                <div className="col-sm-4">
                                    <Link to="/users" className="btn btn-secondary btn-block">
                                        Cancel
                                    </Link>
                                </div>
                                <div className="col-sm-4">
                                    <button type="submit" className="btn btn-primary btn-block">
                                        {props.usersUpdating ? "Saving..." : "Save"}
                                    </button>
                                </div>
                            </div>
                        ) : null}
                    </form>
                </div>
            </div>
        )
    }
}

const mapUserDetailStateToProps = (state, ownProps) => {
    let userID = ownProps.userID
    if (!userID) {
        userID = ownProps.match.params.userID
    }
    let isSelf = state.authentication.token.sub === userID

    return {
        basePath: ownProps.home ? "/me" : `/users/${userID}`,
        isSelf: isSelf,
        userID: userID,
        user: state.users.users ? state.users.users[userID] : undefined,
        usersLoading: state.users.loading,
        usersUpdating: state.users.updating,
        countries: state.codes.codes[CATEGORY_COUNTRIES],
        languages: state.codes.codes[CATEGORY_LANGUAGES],
        licenses: state.codes.codes[CATEGORY_LICENSES],
        codesLoading: state.codes.loading,
        isHome: ownProps.isSelf,
        canSeePersonal: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeOrganizations: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeWildcardUserRoles: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
        canEditPersonal: isSelf || (state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined),
        canEditPassword: isSelf || (state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined),
        validationsLoading: state.validations.loading,
        forbidden: state.users.forbidden
    }
}

const mapUserDetailDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            saveUser,
            loadCodes,
            loadUserRights,
            open,
            close,
            push
        },
        dispatch
    )

UserDetail = withRouter(connect(mapUserDetailStateToProps, mapUserDetailDispatchToProps)(UserDetail))

class UserRoutes extends React.Component {
    componentDidMount() {
        if (
            this.props.canSeePersonal === undefined ||
            this.props.canSeeOrganizations === undefined ||
            this.props.canSeeClinics === undefined ||
            this.props.canSeeWildcardUserRoles === undefined
        ) {
            this.props.loadUserRights()
        }
    }

    componentWillReceiveProps(nextProps) {
        if (
            (nextProps.canSeePersonal === undefined ||
                nextProps.canSeeOrganizations === undefined ||
                nextProps.canSeeClinics === undefined ||
                nextProps.canSeeWildcardUserRoles === undefined) &&
            !nextProps.validationsLoading
        ) {
            this.props.loadUserRights()
        }
    }

    render() {
        let { match, home, userID, canSeePersonal, canSeeOrganizations, canSeeClinics, canSeeWildcardUserRoles } = this.props
        return (
            <Switch>
                {canSeePersonal && <Route exact path={match.url} component={() => <UserDetail home={home} userID={userID} />} />}
                {canSeeOrganizations && (
                    <Route exact path={joinPaths(match.url, "organizations")} component={() => <OrganizationsList home={home} userID={userID} />} />
                )}
                {canSeeOrganizations && (
                    <Route
                        exact
                        path={joinPaths(match.url, "organizations", ":organizationID")}
                        component={() => <OrganizationsList home={home} userID={userID} />}
                    />
                )}
                {canSeeClinics && <Route exact path={joinPaths(match.url, "clinics")} component={() => <ClinicsList home={home} userID={userID} />} />}
                {canSeeClinics && (
                    <Route exact path={joinPaths(match.url, "clinics", ":clinicID")} component={() => <ClinicsList home={home} userID={userID} />} />
                )}
                {canSeeWildcardUserRoles && (
                    <Route path={joinPaths(match.url, "userroles")} component={() => <WildcardUserRolesList home={home} userID={userID} />} />
                )}
            </Switch>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let userID = ownProps.userID
    if (!userID) {
        userID = ownProps.match.params.userID
    }

    return {
        userID: userID,
        canSeePersonal: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeOrganizations: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeWildcardUserRoles: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserRoutes))
