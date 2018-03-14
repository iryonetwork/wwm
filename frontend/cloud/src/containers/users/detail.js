import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"

import { loadUser, saveUser } from "../../modules/users"
import { open, close, COLOR_DANGER } from "shared/modules/alert"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            username: "",
            email: props.user ? props.user.email : "",
            password: "",
            password2: ""
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

    updateRole = roleID => e => {
        let roles = { ...this.state.roles }
        roles[roleID] = !roles[roleID]
        this.setState({
            roles: roles
        })
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

        if (this.state.password !== "") {
            user.password = this.state.password
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
