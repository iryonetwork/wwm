import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import moment from "moment"
import _ from "lodash"

import { loadUser, saveUser } from "../../modules/users"
import { CATEGORY_COUNTRIES, CATEGORY_LANGUAGES, CATEGORY_LICENSES, loadCodes } from "../../modules/codes"
import { ADMIN_RIGHTS_RESOURCE, SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import OrganizationsList from "./organizationsList"
import ClinicsList from "./clinicsList"
import WildcardUserRolesList from "./wildcardUserRolesList"
import { processStateOnChange, processStateOnBlur } from "../../utils/formFieldsUpdate"

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
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
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
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            (!props.user && props.userID !== "new") ||
            props.usersLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
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

            // format dates
            if (personalData && personalData.dateOfBirth) {
                personalData.dateOfBirth = moment(personalData.dateOfBirth).format("DD/MM/YYYY")
            }
            if (personalData && personalData.passport && personalData.passport.expiryDate) {
                personalData.passport.expiryDate = moment(personalData.passport.expiryDate).format("DD/MM/YYYY")
            }
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

    processDateString = (previousStringValue, currentStringValue) => {
        let date = ""
        let finalIndex = 0

        for (var i = 0; i < currentStringValue.length; i++) {
            if (finalIndex < 2 || (finalIndex > 2 && finalIndex < 5) || finalIndex > 5) {
                let digit = parseInt(currentStringValue.charAt(i))
                if (!isNaN(digit)) {
                    date += currentStringValue.charAt(i)
                    finalIndex++
                }
            } else {
                date += "/"
                finalIndex++
            }
        }

        if (previousStringValue.length === 3 && date.length === 2) {
            date = date.substring(0, 1)
        } else if (date.length === 2) {
            date += "/"
        } else if (previousStringValue.length === 6 && date.length === 5) {
            date = date.substring(0, 4)
        } else if (date.length === 5) {
            date += "/"
        }

        return date.substring(0, 10)
    }

    updateDateOfBirth = e => {
        let dateOfBirth = this.processDateString(this.state.personalData.dateOfBirth ? this.state.personalData.dateOfBirth : "", e.target.value)

        var caretLocation = e.target.selectionStart
        if (caretLocation === 2) {
            caretLocation = 3
        } else if (caretLocation === 5) {
            caretLocation = 6
        }

        let validationErrors = this.state.validationErrors
        if (dateOfBirth === "") {
            validationErrors[e.target.id] = "Required"
        } else if (this.state.validationErrors[e.target.id]) {
            delete validationErrors[e.target.id]
        }

        this.setState(
            {
                personalData: _.assign({}, this.state.personalData, _.fromPairs([["dateOfBirth", dateOfBirth]])),
                validationErrors: validationErrors
            },
            () => {
                this.refs["personalData.dateOfBirth"].selectionStart = this.refs["personalData.dateOfBirth"].selectionEnd = caretLocation
            }
        )
    }

    updatePassportExpiryDate = e => {
        let expiryDate = this.processDateString(
            this.state.personalData.passport && this.state.personalData.passport.expiryDate ? this.state.personalData.passport.expiryDate : "",
            e.target.value
        )

        var caretLocation = e.target.selectionStart
        if (caretLocation === 2) {
            caretLocation = 3
        } else if (caretLocation === 5) {
            caretLocation = 6
        }

        let validationErrors = this.state.validationErrors
        if (expiryDate === "") {
            validationErrors[e.target.id] = "Required"
        } else if (this.state.validationErrors[e.target.id]) {
            delete validationErrors[e.target.id]
        }

        let passportData = this.state.personalData.passport ? this.state.personalData.passport : {}
        passportData = _.assign({}, passportData, _.fromPairs([["expiryDate", expiryDate]]))

        this.setState(
            {
                personalData: _.assign({}, this.state.personalData, _.fromPairs([["passport", passportData]])),
                validationErrors: validationErrors
            },
            () => {
                this.refs["personalData.passport.expiryDate"].selectionStart = this.refs["personalData.passport.expiryDate"].selectionEnd = caretLocation
            }
        )
    }

    newLanguage = () => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages = [...personalData.languages, { id: undefined, edit: true }]
        } else {
            personalData.languages = [{ id: undefined, edit: true }]
        }
        this.setState({ personalData: personalData })
    }

    updateLanguage = index => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages[index].id = e.target.value
        }
        this.setState({ personalData: personalData })
    }

    removeLanguage = index => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages.splice(index, 1)
        }
        this.setState({ personalData: personalData })
    }

    newLicense = () => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses = [...personalData.licenses, { id: undefined, edit: true }]
        } else {
            personalData.licenses = [{ id: undefined, edit: true }]
        }
        this.setState({ personalData: personalData })
    }

    updateLicense = index => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses[index].id = e.target.value
        }
        this.setState({ personalData: personalData })
    }

    removeLicense = index => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses.splice(index, 1)
        }
        this.setState({ personalData: personalData })
    }

    updateInput = e => {
        this.setState(processStateOnChange(this.state, e))
    }

    onBlurInput = e => {
        this.setState(processStateOnBlur(this.state, e))
    }

    submit = e => {
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

        // format dates
        if (user.personalData.dateOfBirth && user.personalData.dateOfBirth !== "") {
            let dateOfBirth = moment(user.personalData.dateOfBirth, "DD/MM/YYYY")
            if (!dateOfBirth.isValid()) {
                validationErrors["personalData.passport.expiryDate"] = "Invalid date of birth"
            } else {
                user.personalData.dateOfBirth = dateOfBirth.local().format("YYYY-MM-DD")
            }
        }
        if (this.state.personalData.passport) {
            user.personalData.passport = _.clone(this.state.personalData.passport)
            if (user.personalData.passport.expiryDate === "" || user.personalData.passport.expiryDate === undefined) {
                delete user.personalData.passport.expiryDate
            } else {
                if (user.personalData.passport.expiryDate.length !== 10) {
                    validationErrors["personalData.passport.expiryDate"] = "Invalid passport expiry date"
                } else {
                    let expiryDate = moment(user.personalData.passport.expiryDate, "DD/MM/YYYY")
                    if (!expiryDate.isValid()) {
                        validationErrors["personalData.passport.expiryDate"] = "Invalid passport expiry date"
                    } else {
                        user.personalData.passport.expiryDate = expiryDate.local().format("YYYY-MM-DD")
                    }
                }
            }
        }

        // format languages
        if (user.personalData.languages && user.personalData.languages.length !== 0) {
            user.personalData.languages = _.map(_.pickBy(user.personalData.languages, language => language.id && language.id !== ""), language => language.id)
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
                this.props.history.push(`/users/${response.id}`)
            }
        })
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <div>
                {props.home ? (
                    <h1>Hi, {props.user.username}</h1>
                ) : (
                    <div>
                        <h1>Users</h1>
                        <h2>{props.user ? props.user.username : "Add new user"}</h2>
                    </div>
                )}
                <div>
                    <form onSubmit={this.submit} className="needs-validation" noValidate>
                        {props.user ? null : (
                            <div className="form-group">
                                <label htmlFor="username">Username</label>
                                <input
                                    type="text"
                                    className={"form-control" + (this.state.validationErrors["username"] ? " is-invalid" : "")}
                                    id="username"
                                    value={this.state.username}
                                    onChange={this.updateInput}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="username"
                                    required="true"
                                />
                                {this.state.validationErrors["username"] ? (
                                    <div className="invalid-feedback">{this.state.validationErrors["username"]}</div>
                                ) : (
                                    <small className="form-text text-muted">Required</small>
                                )}
                            </div>
                        )}
                        {!props.user || props.canEditPassword ? (
                            <div>
                                <div className="form-group">
                                    <label htmlFor="password">{props.user ? "Enter new password" : "Enter password"}</label>
                                    <input
                                        type="password"
                                        className={"form-control" + (this.state.validationErrors["password"] ? " is-invalid" : "")}
                                        id="password"
                                        value={this.state.password}
                                        onChange={this.updateInput}
                                        onBlur={this.onBlurInput}
                                        disabled={!props.canEdit}
                                        placeholder={props.user ? "●●●●●" : "password"}
                                        required={props.user ? null : "true"}
                                    />
                                    <div className="invalid-feedback">{this.state.validationErrors["password"]}</div>
                                    {!props.user ? <small className="form-text text-muted">Required</small> : null}
                                </div>
                                <div className="form-group">
                                    <label htmlFor="password2">{props.user ? "Enter new password again" : "Enter password again"}</label>
                                    <input
                                        type="password"
                                        className={"form-control" + (this.state.validationErrors["password2"] ? " is-invalid" : "")}
                                        id="password2"
                                        value={this.state.password2}
                                        onChange={this.updateInput}
                                        disabled={!props.canEdit}
                                        placeholder={props.user ? "●●●●●" : "password"}
                                        required={props.user ? null : "true"}
                                    />
                                    <div className="invalid-feedback">{this.state.validationErrors["password2"]}</div>
                                    {!props.user ? <small className="form-text text-muted">Required</small> : null}
                                </div>
                            </div>
                        ) : null}
                        <div className="form-group">
                            <label htmlFor="email">Email address</label>
                            <input
                                type="email"
                                className={"form-control" + (this.state.validationErrors["email"] ? " is-invalid" : "")}
                                id="email"
                                value={this.state.email}
                                onChange={this.updateInput}
                                onBlur={this.onBlurInput}
                                disabled={!props.canEdit}
                                placeholder="user@email.com"
                                required="true"
                            />
                            {this.state.validationErrors["email"] ? (
                                <div className="invalid-feedback">{this.state.validationErrors["email"]}</div>
                            ) : (
                                <small className="form-text text-muted">Required</small>
                            )}
                        </div>
                        <div className="form-group">
                            <h3>Personal data</h3>
                            <div className="form-group">
                                <label htmlFor="firstName">First name</label>
                                <input
                                    type="text"
                                    className={"form-control" + (this.state.validationErrors["personalData.firstName"] ? " is-invalid" : "")}
                                    id="personalData.firstName"
                                    value={this.state.personalData.firstName}
                                    onChange={this.updateInput}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="First name"
                                    required="true"
                                />
                                {this.state.validationErrors["personalData.firstName"] ? (
                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.firstName"]}</div>
                                ) : (
                                    <small className="form-text text-muted">Required</small>
                                )}
                            </div>
                            <div className="form-group">
                                <label htmlFor="middleName">Middle name</label>
                                <input
                                    type="text"
                                    className="form-control"
                                    id="personalData.middleName"
                                    value={this.state.personalData.middleName}
                                    onChange={this.updateInput}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="Middle name"
                                />
                            </div>
                            <div className="form-group">
                                <label htmlFor="lastName">Last name</label>
                                <input
                                    type="text"
                                    className={"form-control" + (this.state.validationErrors["personalData.lastName"] ? " is-invalid" : "")}
                                    id="personalData.lastName"
                                    value={this.state.personalData.lastName}
                                    onChange={this.updateInput}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="Last name"
                                    required="true"
                                />
                                {this.state.validationErrors["personalData.lastName"] ? (
                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.lastName"]}</div>
                                ) : (
                                    <small className="form-text text-muted">Required</small>
                                )}
                            </div>
                            <div className="form-group">
                                <label htmlFor="dateOfBirth">Date of birth</label>
                                <input
                                    type="text"
                                    className={"form-control" + (this.state.validationErrors["personalData.dateOfBirth"] ? " is-invalid" : "")}
                                    id="personalData.dateOfBirth"
                                    ref="personalData.dateOfBirth"
                                    value={this.state.personalData.dateOfBirth}
                                    onChange={this.updateDateOfBirth}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="DD/MM/YYYY"
                                    required="true"
                                />
                                {this.state.validationErrors["personalData.dateOfBirth"] ? (
                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.dateOfBirth"]}</div>
                                ) : (
                                    <small className="form-text text-muted">Required</small>
                                )}
                            </div>
                            <div className="form-group">
                                <label htmlFor="specialisation">Specialisation</label>
                                <input
                                    type="text"
                                    className="form-control"
                                    id="personalData.specialisation"
                                    value={this.state.personalData.specialisation}
                                    onChange={this.updateInput}
                                    onBlur={this.onBlurInput}
                                    disabled={!props.canEdit}
                                    placeholder="Medical worker specialisation"
                                />
                            </div>
                            <div className="form-group">
                                <h4>Languages</h4>
                                <table className="table table-hover table-sm">
                                    <tbody>
                                        {_.map(this.state.personalData.languages ? this.state.personalData.languages : [], (language, i) => (
                                            <tr key={i}>
                                                <td className="col-6">
                                                    {language.edit ? (
                                                        <select
                                                            className="form-control form-control-sm"
                                                            id="personalData.languages"
                                                            value={language.id}
                                                            onChange={this.updateLanguage(i)}
                                                            disabled={!props.canEdit}
                                                        >
                                                            <option value="">Select language</option>
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
                                                    ) : props.languages[language.id] ? (
                                                        props.languages[language.id].title
                                                    ) : (
                                                        language.id
                                                    )}
                                                </td>
                                                <td className="text-right col-1">
                                                    {props.canEdit ? (
                                                        <button onClick={this.removeLanguage(i)} className="btn btn-sm btn-light" type="button">
                                                            {language.edit ? <span className="icon_close" /> : <span className="icon_trash" />}
                                                        </button>
                                                    ) : null}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                                {props.canEdit ? (
                                    <button type="button" className="btn btn-sm btn-outline-primary col" onClick={this.newLanguage()}>
                                        Add language
                                    </button>
                                ) : null}
                            </div>
                            <div className="form-group">
                                <h4>Licenses</h4>
                                <table className="table table-hover table-sm">
                                    <tbody>
                                        {_.map(this.state.personalData.licenses ? this.state.personalData.licenses : [], (license, i) => (
                                            <tr key={i}>
                                                <td className="col-6">
                                                    {license.edit ? (
                                                        <select
                                                            className="form-control form-control-sm"
                                                            id="personalData.licenses"
                                                            value={license.id}
                                                            onChange={this.updateLicense(i)}
                                                            disabled={!props.canEdit}
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
                                                    ) : props.licenses[license.id] ? (
                                                        props.licenses[license.id].title
                                                    ) : (
                                                        license.code.id
                                                    )}
                                                </td>
                                                <td className="text-right col-1">
                                                    {props.canEdit ? (
                                                        <button onClick={this.removeLicense(i)} className="btn btn-sm btn-light" type="button">
                                                            {license.edit ? <span className="icon_close" /> : <span className="icon_trash" />}
                                                        </button>
                                                    ) : null}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                                {props.canEdit ? (
                                    <button type="button" className="btn btn-sm btn-outline-primary col" onClick={this.newLicense()}>
                                        Add license
                                    </button>
                                ) : null}
                            </div>
                            <div className="form-group">
                                <label htmlFor="nationality">Nationality</label>
                                <select
                                    className="form-control form-control-sm"
                                    id="personalData.nationality"
                                    value={this.state.personalData.nationality}
                                    onChange={this.updateInput}
                                    disabled={!props.canEdit}
                                >
                                    <option value="">Select country</option>
                                    {_.map(props.countries, country => (
                                        <option key={country.id} value={country.id}>
                                            {country.title}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <div className="form-group">
                                <label htmlFor="residency">Residency</label>
                                <select
                                    className="form-control form-control-sm"
                                    id="personalData.residency"
                                    value={this.state.personalData.residency}
                                    onChange={this.updateInput}
                                    disabled={!props.canEdit}
                                >
                                    <option value="">Select country</option>
                                    {_.map(props.countries, country => (
                                        <option key={country.id} value={country.id}>
                                            {country.title}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <div className="form_group">
                                <h4>Passport</h4>
                                <div className="form-group">
                                    <label htmlFor="number">Number</label>
                                    <input
                                        type="text"
                                        className="form-control"
                                        id="personalData.passport.number"
                                        value={this.state.personalData.passport ? this.state.personalData.passport.number : undefined}
                                        onChange={this.updateInput}
                                        onBlur={this.onBlurInput}
                                        placeholder="Passport number"
                                        disabled={!props.canEdit}
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="issuingCountry">Issuing country</label>
                                    <select
                                        className="form-control form-control-sm"
                                        id="personalData.passport.issuingCountry"
                                        value={this.state.personalData.passport ? this.state.personalData.passport.issuingCountry : undefined}
                                        onChange={this.updateInput}
                                        disabled={!props.canEdit}
                                    >
                                        <option value="">Select country</option>
                                        {_.map(props.countries, country => (
                                            <option key={country.id} value={country.id}>
                                                {country.title}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                                <div className="form-group">
                                    <label htmlFor="expiryDate">Expiry date</label>
                                    <input
                                        type="text"
                                        className={"form-control" + (this.state.validationErrors["personalData.passport.expiryDate"] ? " is-invalid" : "")}
                                        id="personalData.passport.expiryDate"
                                        ref="personalData.passport.expiryDate"
                                        value={this.state.personalData.passport ? this.state.personalData.passport.expiryDate : undefined}
                                        onChange={this.updatePassportExpiryDate}
                                        onBlur={this.onBlurInput}
                                        placeholder="DD/MM/YYYY"
                                        disabled={!props.canEdit}
                                    />
                                    <div className="invalid-feedback">{this.state.validationErrors["personalData.passport.expiryDate"]}</div>
                                </div>
                            </div>
                        </div>
                        {props.canEdit ? (
                            <div className="form-group">
                                <button type="submit" className="btn btn-outline-primary col">
                                    Save
                                </button>
                            </div>
                        ) : null}
                    </form>
                </div>
                {props.user ? (
                    <div className="m-4">
                        <div className="m-4">
                            <OrganizationsList userID={props.userID} />
                        </div>
                        <div className="m-4">
                            <ClinicsList userID={props.userID} />
                        </div>
                        <div className="m-4">
                            <WildcardUserRolesList userID={props.userID} />
                        </div>
                    </div>
                ) : null}
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.userID
    if (!id) {
        id = ownProps.match.params.userID
    }

    return {
        userID: id,
        user: state.users.users ? state.users.users[id] : undefined,
        usersLoading: state.users.loading,
        countries: state.codes.codes[CATEGORY_COUNTRIES],
        languages: state.codes.codes[CATEGORY_LANGUAGES],
        licenses: state.codes.codes[CATEGORY_LICENSES],
        codesLoading: state.codes.loading,
        isHome: ownProps.home,
        canSee:
            state.authentication.token.sub === id || (state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined),
        canEdit:
            state.authentication.token.sub === id || (state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined),
        canEditPassword:
            state.authentication.token.sub === id || (state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined),
        validationsLoading: state.validations.loading,
        forbidden: state.users.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            saveUser,
            loadCodes,
            loadUserRights,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserDetail))
