import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import moment from "moment"
import _ from "lodash"

import { loadUser, saveUser } from "../../modules/users"
import { CATEGORY_COUNTRIES, CATEGORY_LANGUAGES, CATEGORY_LICENSES, loadCodes } from "../../modules/codes"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import OrganizationsList from "./organizationsList"
import ClinicsList from "./clinicsList"
import WildcardUserRolesList from "./wildcardUserRolesList"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            username: "",
            email: "",
            password: "",
            password2: "",
            personalData: {},
            loading: true
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

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = (!props.user && props.userID !== "new") || props.usersLoading
        this.setState({ loading: loading })

        if (props.user) {
            let personalData = _.clone(props.user.personalData)
            if (props.user.personalData.passport) {
                personalData.passport = _.clone(props.user.personalData.passport)
            }
            if (props.user.personalData.languages) {
                personalData.languages = _.clone(props.user.personalData.languages)
            }
            if (props.user.personalData.licenses) {
                personalData.licenses = _.clone(props.user.personalData.licenses)
            }

            // format dates
            if (personalData && personalData.dateOfBirth) {
                personalData.dateOfBirth = moment(personalData.dateOfBirth).format('DD/MM/YYYY')
            }
            if (personalData && personalData.passport && personalData.passport.expiryDate) {
                personalData.passport.expiryDate = moment(personalData.passport.expiryDate).format('DD/MM/YYYY')
            }
            // format languages
            if (personalData && personalData.languages) {
                personalData.languages = _.map(personalData.languages, languageCodeID => {return {"code_id": languageCodeID}})
            }
            // format licenses
            if (personalData && personalData.licenses) {
                personalData.licenses = _.map(personalData.licenses, licenseCodeID => {return {"code_id": licenseCodeID}})
            }

            this.setState({ email: props.user.email })
            this.setState({ personalData: personalData ? personalData : {}})
        }
    }

    updateEmail = e => {
        this.setState({ email: e.target.value })
    }

    updatePassword = e => {
        this.setState({ password: e.target.value })
    }

    updatePassword2 = e => {
        this.setState({ password2: e.target.value })
    }

    updateUsername = e => {
        this.setState({ username: e.target.value })
    }

    updatePersonalData = e => {
        const target = e.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const id = target.id;

        this.setState({ personalData: _.assign({}, this.state.personalData, _.fromPairs([[id, value]])) });
    }

    updatePassportData = e => {
        const target = e.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const id = target.id;

        let passportData = this.state.personalData.passport ? this.state.personalData.passport : {}
        passportData = _.assign({}, passportData, _.fromPairs([[id, value]]))

        this.setState({ personalData: _.assign({}, this.state.personalData, _.fromPairs([["passport", passportData]])) });
    }

    processDateString = (previousStringValue, currentStringValue) => {
        let date = ""
        let finalIndex = 0

        for (var i = 0; i < currentStringValue.length; i++) {
            if (finalIndex < 2 || (finalIndex > 2 && finalIndex < 5) || finalIndex > 5 ) {
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

        if (previousStringValue.length === 3 && date.length === 2 ) {
            date = date.substring(0, 1)
        } else if (date.length === 2) {
            date += "/"
        } else if (previousStringValue.length === 6 && date.length === 5 ) {
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

        this.setState(
            { personalData: _.assign({}, this.state.personalData, _.fromPairs([["dateOfBirth", dateOfBirth]]))},
            () => {
                this.refs.dateOfBirth.selectionStart = this.refs.dateOfBirth.selectionEnd = caretLocation
            }
        )
    }

    updatePassportExpiryDate = e => {
        let expiryDate = this.processDateString((this.state.personalData.passport && this.state.personalData.passport.expiryDate) ? this.state.personalData.passport.expiryDate : "", e.target.value)

        var caretLocation = e.target.selectionStart
        if (caretLocation === 2) {
            caretLocation = 3
        } else if (caretLocation === 5) {
            caretLocation = 6
        }

        let passportData = this.state.personalData.passport ? this.state.personalData.passport : {}
        passportData = _.assign({}, passportData, _.fromPairs([["expiryDate", expiryDate]]))
        this.setState(
            { personalData: _.assign({}, this.state.personalData, _.fromPairs([["passport", passportData]]))},
            () => {
                this.refs.expiryDate.selectionStart = this.refs.expiryDate.selectionEnd = caretLocation
            }
        )
    }

   newLanguage = () => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages = [...personalData.languages, { code_id: undefined, edit: true }]
        } else {
            personalData.languages = [{ code_id: undefined, edit: true }]
        }
        this.setState({ personalData: personalData })
    }

    updateLanguage = index => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages[index].code_id = e.target.value
        }
        this.setState({ personalData: personalData})
    }

    removeLanguage = index => e => {
        let personalData = this.state.personalData
        if (personalData.languages) {
            personalData.languages.splice(index, 1)
        }
        this.setState({ personalData: personalData})
    }

   newLicense = () => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses = [...personalData.licenses, { code_id: undefined, edit: true }]
        } else {
            personalData.licenses = [{ code_id: undefined, edit: true }]
        }
        this.setState({ personalData: personalData })
    }

    updateLicense = index => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses[index].code_id = e.target.value
        }
        this.setState({ personalData: personalData})
    }

    removeLicense = index => e => {
        let personalData = this.state.personalData
        if (personalData.licenses) {
            personalData.licenses.splice(index, 1)
        }
        this.setState({ personalData: personalData})
    }

    submit = e => {
        e.preventDefault()
        this.props.close()
        if (!this.props.user && this.state.username === "") {
            this.props.open("Username is required", "", COLOR_DANGER)
            return
        }

        if (this.state.email === "") {
            this.props.open("Email is required", "", COLOR_DANGER)
            return
        }

        if (this.state.password !== this.state.password2) {
            this.props.open("Passwords aren't the same", "", COLOR_DANGER)
            return
        }

        if (!this.props.user && this.state.password === "") {
            this.props.open("Password is required", "", COLOR_DANGER)
            return
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
        if (this.state.password && this.state.password !== "") {
            user.password = this.state.password
        }

        user.personalData = _.clone(this.state.personalData)

        // format dates
        if (user.personalData.dateOfBirth && user.personalData.dateOfBirth !== "") {
            let dateOfBirth = moment(user.personalData.dateOfBirth, "DD/MM/YYYY")
            if (!dateOfBirth.isValid()) {
                this.props.open("Invalid date of birth", "", COLOR_DANGER)
                return
            }
            user.personalData.dateOfBirth = dateOfBirth.local().format("YYYY-MM-DD")
        }
        if (user.personalData.passport && user.personalData.passport.expiryDate && user.personalData.passport.expiryDate !== "") {
            let expiryDate = moment(user.personalData.passport.expiryDate, "DD/MM/YYYY")
            if (!expiryDate.isValid()) {
                this.props.open("Invalid passport expiry date", "", COLOR_DANGER)
                return
            }
            user.personalData.passport.expiryDate = expiryDate.local().format("YYYY-MM-DD")
        }

        // format languages
        if (user.personalData.languages && user.personalData.languages.length !== 0) {
            user.personalData.languages = _.map(_.pickBy(user.personalData.languages, language => (language.code_id && language.code_id !== "")), language => language.code_id)
        }

        // format licenses
        if (user.personalData.licenses && user.personalData.licenses.length !== 0) {
            user.personalData.licenses = _.map(_.pickBy(user.personalData.licenses, license => (license.code_id && license.code_id !== "")), license => license.code_id)
        }

        this.props.saveUser(user)
            .then(response => {
                if (!user.id && response.id) {
                    this.props.history.push(`/users/${response.id}`)
                }
            })
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
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
                <form onSubmit={this.submit}>
                    {props.user ? null : (
                        <div className="form-group">
                            <label htmlFor="username">Username</label>
                            <input className="form-control" id="username" value={this.state.username} onChange={this.updateUsername} placeholder="username"/>
                        </div>
                    )}
                    <div className="form-group">
                        <label htmlFor="password">{props.user ? "Enter new password" : "Enter password"}</label>
                        <input type="password" className="form-control" id="paswword" value={this.state.password} onChange={this.updatePassword} placeholder={props.user ? "●●●●●" : "password"}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="password2">{props.user ? "Enter new password again" : "Enter password again"}</label>
                        <input type="password" className="form-control" id="paswword2" value={this.state.password2} onChange={this.updatePassword2} placeholder={props.user ? "●●●●●" : "password"}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input type="email" className="form-control" id="email" value={this.state.email} onChange={this.updateEmail} placeholder="user@email.com"/>
                    </div>
                    <div className="form-group">
                        <h3>Personal data</h3>
                        <div className="form-group">
                            <label htmlFor="firstName">First name</label>
                            <input className="form-control" id="firstName" value={this.state.personalData.firstName} onChange={this.updatePersonalData} placeholder="First name" />
                        </div>
                        <div className="form-group">
                            <label htmlFor="middleName">Middle name</label>
                            <input className="form-control" id="middleName" value={this.state.personalData.middleName} onChange={this.updatePersonalData} placeholder="Middle name" />
                        </div>
                        <div className="form-group">
                            <label htmlFor="lastName">Last name</label>
                            <input className="form-control" id="lastName" value={this.state.personalData.lastName} onChange={this.updatePersonalData} placeholder="Last name" />
                        </div>
                        <div className="form-group">
                            <label htmlFor="dateOfBirth">Date of birth</label>
                             <input className="form-control" id="dateOfBirth" ref="dateOfBirth" value={this.state.personalData.dateOfBirth} onChange={this.updateDateOfBirth} placeholder="DD/MM/YYYY" />
                        </div>
                        <div className="form-group">
                            <label htmlFor="specialisation">Specialisation</label>
                            <input className="form-control" id="specialisation" value={this.state.personalData.specialisation} onChange={this.updatePersonalData} placeholder="Medical worker specialisation" />
                        </div>
                        <div className="form-group">
                            <h4>Languages</h4>
                            <table className="table table-hover table-sm">
                                <tbody>
                                    {_.map(this.state.personalData.languages ? this.state.personalData.languages : [], (language, i) => (
                                        <tr key={i}>
                                            <td>
                                                {language.edit ? (
                                                    <select className="form-control form-control-sm" id="residency" value={language.code_id} onChange={this.updateLanguage(i)}>
                                                        <option value="">Select language</option>
                                                        {_.map(_.difference(_.map(props.languages, language => language.code_id), _.without(_.map(this.state.personalData.languages, language => language.code_id), language.code_id)), languageCodeID => (
                                                            <option key={languageCodeID} value={languageCodeID}>
                                                                {props.languages[languageCodeID].title}
                                                            </option>
                                                        ))}
                                                    </select>
                                                ) : (
                                                    props.languages[language.code_id] ? props.languages[language.code_id].title : language.code_id
                                                )}
                                            </td>
                                            <td className="text-right">
                                                <button onClick={this.removeLanguage(i)} className="btn btn-sm btn-light" type="button">
                                                    {language.edit ? (
                                                        <span className="icon_close" />
                                                    ) : (
                                                        <span className="icon_trash" />
                                                    )}
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                            <button type="button" className="btn btn-sm btn-outline-secondary col" onClick={this.newLanguage()}>
                                Add language
                            </button>
                        </div>
                        <div className="form-group">
                            <h4>Licenses</h4>
                            <table className="table table-hover table-sm">
                                <tbody>
                                    {_.map(this.state.personalData.licenses ? this.state.personalData.licenses : [], (license, i) => (
                                        <tr key={i}>
                                            <td>
                                                {license.edit ? (
                                                    <select className="form-control form-control-sm" id="residency" value={license.code_id} onChange={this.updateLicense(i)}>
                                                        <option value="">Select license</option>
                                                        {_.map(_.difference(_.map(props.licenses, license => license.code_id), _.without(_.map(this.state.personalData.licenses, license => license.code_id), license.code_id)), licenseCodeID => (
                                                            <option key={licenseCodeID} value={licenseCodeID}>
                                                                {props.licenses[licenseCodeID].title}
                                                            </option>
                                                        ))}
                                                    </select>
                                                ) : (
                                                    props.licenses[license.code_id] ? props.licenses[license.code_id].title : license.code.id
                                                )}
                                            </td>
                                            <td className="text-right">
                                                <button onClick={this.removeLicense(i)} className="btn btn-sm btn-light" type="button">
                                                    {license.edit ? (
                                                        <span className="icon_close" />
                                                    ) : (
                                                        <span className="icon_trash" />
                                                    )}
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                            <button type="button" className="btn btn-sm btn-outline-secondary col" onClick={this.newLicense()}>
                                Add license
                            </button>
                        </div>
                        <div className="form-group">
                            <label htmlFor="nationality">Nationality</label>
                            <select className="form-control form-control-sm" id="nationality" value={this.state.personalData.nationality} onChange={this.updatePersonalData}>
                                <option value="">Select country</option>
                                {_.map(props.countries, country => (
                                    <option key={country.code_id} value={country.code_id}>
                                        {country.title}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div className="form-group">
                            <label htmlFor="residency">Residency</label>
                            <select className="form-control form-control-sm" id="residency" value={this.state.personalData.residency} onChange={this.updatePersonalData}>
                                <option value="">Select country</option>
                                {_.map(props.countries, country => (
                                    <option key={country.code_id} value={country.code_id}>
                                        {country.title}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div className="form_group">
                            <h4>Passport</h4>
                            <div className="form-group">
                                <label htmlFor="number">Number</label>
                                <input className="form-control" id="number" value={this.state.personalData.passport ? this.state.personalData.passport.number : undefined} onChange={this.updatePassportData} placeholder="Passport number" />
                            </div>
                            <div className="form-group">
                                <label htmlFor="issuingCountry">Issuing country</label>
                                <select className="form-control form-control-sm" id="issuingCountry" value={this.state.personalData.passport ? this.state.personalData.passport.issuingCountry : undefined} onChange={this.updatePassportData}>
                                    <option value="">Select country</option>
                                    {_.map(props.countries, country => (
                                        <option key={country.code_id} value={country.title}>
                                            {country.title}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <div className="form-group">
                                <label htmlFor="expiryDate">Expiry date</label>
                                 <input className="form-control" id="expiryDate" ref="expiryDate" value={this.state.personalData.passport ? this.state.personalData.passport.expiryDate : undefined} onChange={this.updatePassportExpiryDate} placeholder="DD/MM/YYYY" />
                            </div>
                        </div>
                    </div>
                    <div className="form-group">
                        <button type="submit" className="btn btn-outline-primary col">
                            Save user
                        </button>
                    </div>
                </form>
                </div>
                {props.user ? (
                    <div className="m-4">
                        <div className="m-4">
                            <h2>User's organizations</h2>
                            <OrganizationsList userID={props.userID} />
                        </div>
                        <div className="m-4">
                            <h2>User's clinics</h2>
                            <ClinicsList userID={props.userID} />
                        </div>
                        <div className="m-4">
                            <h2>User's wildcard roles</h2>
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
        isHome: ownProps.home
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            saveUser,
            loadCodes,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserDetail))
