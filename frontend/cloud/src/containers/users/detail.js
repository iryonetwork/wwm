import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import moment from "moment"

import { loadUser, saveUser } from "../../modules/users"
import { open, close, COLOR_DANGER } from "shared/modules/alert"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        if (props.user) {
            this.state = {
                username: "",
                email: props.user ? props.user.email : "",
                firstName: props.user.personalData ? (props.user.personalData.firstName ? props.user.personalData.firstName : "" ) : "",
                middleName: props.user.personalData ? (props.user.personalData.middleName ? props.user.personalData.middleName : "" ) : "",
                lastName: props.user.personalData ? (props.user.personalData.lastName ? props.user.personalData.lastName : "" ) : "",
                dateOfBirth: props.user.personalData ? (props.user.personalData.dateOfBirth ? moment(props.user.personalData.dateOfBirth).format('DD/MM/YYYY') : "" ) : "",
                password: "",
                password2: ""
            }
        } else {
            this.state = {
                username: "",
                email: "",
                firstName: "",
                middleName: "",
                lastName: "",
                dateOfBirth: "",
                password: "",
                password2: ""
            }
        }
    }

    componentDidMount() {
        if (!this.props.user && this.props.userID !== "new") {
            this.props.loadUser(this.props.userID)
        }
    }

    componentWillReceiveProps(props) {
        if (props.user) {
            this.setState({ email: props.user.email })
            this.setState({ firstName: props.user.personalData ? (props.user.personalData.firstName ? props.user.personalData.firstName : "" ) : "" })
            this.setState({ middleName: props.user.personalData ? (props.user.personalData.middleName ? props.user.personalData.middleName : "" ) : "" })
            this.setState({ lastName: props.user.personalData ? (props.user.personalData.lastName ? props.user.personalData.lastName : "" ) : "" })
            this.setState({ dateOfBirth: props.user.personalData ? (props.user.personalData.dateOfBirth ? moment(props.user.personalData.dateOfBirth).format('DD/MM/YYYY') : "" ) : "" })
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

    // updateDateOfBirth = e => {
    //     this.setState({ dateOfBirth: e.target.value })
    // }

    updateFirstName = e => {
        this.setState({ firstName: e.target.value })
    }

    updateMiddleName = e => {
        this.setState({ middleName: e.target.value })
    }

    updateLastName = e => {
        this.setState({ lastName: e.target.value })
    }

    updateDateOfBirth = e => {
        let previousValue = this.state.dateOfBirth ? this.state.dateOfBirth : ""
        let currentInput = e.target.value
        let dateOfBirth = ""
        let finalIndex = 0

        for (var i = 0; i < currentInput.length; i++) {
            if (finalIndex < 2 || (finalIndex > 2 && finalIndex < 5) || finalIndex > 5 ) {
                let digit = parseInt(currentInput.charAt(i))
                if (!isNaN(digit)) {
                    dateOfBirth += currentInput.charAt(i)
                    finalIndex++
                }
            } else {
                dateOfBirth += "/"
                finalIndex++
            }
        }

        if (previousValue.length === 3 && dateOfBirth.length === 2 ) {
            dateOfBirth = dateOfBirth.substring(0, 1)
        } else if (dateOfBirth.length === 2) {
            dateOfBirth += "/"
        } else if (previousValue.length === 6 && dateOfBirth.length === 5 ) {
            dateOfBirth = dateOfBirth.substring(0, 4)
        } else if (dateOfBirth.length === 5) {
            dateOfBirth += "/"
        }

        var caretLocation = e.target.selectionStart
        if (caretLocation === 2) {
            caretLocation = 3
        } else if (caretLocation === 5) {
            caretLocation = 6
        }

        this.setState(
            { dateOfBirth: dateOfBirth.substring(0, 10) },
            () => {
                this.refs.dateOfBirth.selectionStart = this.refs.dateOfBirth.selectionEnd = caretLocation
            }
        )
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
            user.personalData = {}
        }

        if (this.state.password !== "") {
            user.password = this.state.password
        }

        if (this.state.firstName !== "") {
            if (!user.personalData) {
                user.personalData = {}
            }
            user.personalData.firstName = this.state.firstName.trim()
        }
        if (this.state.middleName !== "") {
            if (!user.personalData) {
                user.personalData = {}
            }
            user.personalData.middleName = this.state.middleName.trim()
        }
        if (this.state.lastName !== "") {
            if (!user.personalData) {
                user.personalData = {}
            }
            user.personalData.lastName = this.state.lastName.trim()
        }
        if (this.state.dateOfBirth !== "") {
            let dateOfBirth = moment(this.state.dateOfBirth, "DD/MM/YYYY")
            console.log(dateOfBirth)
            if (!dateOfBirth.isValid()) {
                this.props.open("Invalid date of birth", "", COLOR_DANGER)
                return
            }
            user.personalData.dateOfBirth = dateOfBirth.local().format("YYYY-MM-DD")
        }

        this.props.saveUser(user)
    }

    render() {
        let props = this.props
        if (!props.user && props.userID !== "new") {
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

                <form onSubmit={this.submit}>
                    {props.user ? null : (
                        <div className="form-group">
                            <label htmlFor="username">Username</label>
                            <input className="form-control" id="username" value={this.state.username} onChange={this.updateUsername} />
                        </div>
                    )}
                    <div className="form-group">
                        <label htmlFor="firstName">First name</label>
                        <input className="form-control" id="firstName" value={this.state.firstName} onChange={this.updateFirstName} placeholder="First name" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="middleName">Middle name</label>
                        <input className="form-control" id="middleName" value={this.state.middleName} onChange={this.updateMiddleName} placeholder="Middle name" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="lastName">Last name</label>
                        <input className="form-control" id="lastName" value={this.state.lastName} onChange={this.updateLastName} placeholder="Last name" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="dateOfBirth">Date of birth</label>
                         <input ref="dateOfBirth" className="form-control" id="dateOfBirth" value={this.state.dateOfBirth} onChange={this.updateDateOfBirth} placeholder="DD/MM/YYYY" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input type="email" className="form-control" id="email" value={this.state.email} onChange={this.updateEmail} />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Change password</label>
                        <input type="password" className="form-control" id="paswword" value={this.state.password} onChange={this.updatePassword} />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password2">New password again</label>
                        <input type="password" className="form-control" id="paswword2" value={this.state.password2} onChange={this.updatePassword2} />
                    </div>
                    <button type="submit" className="btn btn-sm btn-outline-secondary">
                        Save
                    </button>
                </form>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.userID
    if (!id) {
        id = ownProps.match.params.id
    }

    return {
        user: state.users.users ? state.users.users[id] : undefined,
        loading: state.users.loading,
        userID: id,
        isHome: ownProps.home
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            saveUser,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserDetail))
