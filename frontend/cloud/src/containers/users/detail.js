import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import _ from "lodash"

import { loadUsers, saveUser } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadRules } from "../../modules/rules"
import { open, close, COLOR_DANGER } from "../../modules/alert"
//import Rules from "../rules"

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
        if (!this.props.user) {
            this.props.loadUsers()
        }
        //this.props.loadRoles()
        if (!this.props.allRules) {
            this.props.loadRules()
        }
    }

    componentWillReceiveProps(props) {
        if (props.user) {
            this.setState({ email: props.user.email })
        }
        /*
        if (props.allRoles) {
            let selected = _.reduce(
                props.roles,
                (obj, role) => {
                    obj[role] = true
                    return obj
                },
                {}
            )
            let all = _.mapValues(props.allRoles, () => false)

            this.setState({
                roles: _.defaults(selected, all)
            })
        }
        */
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
                <h1>Users</h1>

                <h2>{props.user ? props.user.username : "Add new user"}</h2>

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
                {/*
                <h2>Roles</h2>
                {this.state.roles
                    ? _.map(this.state.roles, (selected, roleID) => {
                          return (
                              <div className="form-check form-check-inline" key={roleID}>
                                  <input
                                      className="form-check-input"
                                      type="checkbox"
                                      checked={selected}
                                      id={`role-checkbox-${roleID}`}
                                      onChange={this.updateRole(roleID)}
                                  />
                                  <label className="form-check-label" htmlFor={`role-checkbox-${roleID}`}>
                                      {props.allRoles[roleID].name}
                                  </label>
                              </div>
                          )
                      })
                    : null}

                <Rules rules={props.rules} subject={props.userID} />
                    */}
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        user: state.users.users ? state.users.users[ownProps.match.params.id] : undefined,
        loading: state.users.loading,
        userID: ownProps.match.params.id,
        roles: _.get(state, `roles.users['${ownProps.match.params.id}']`, []),
        allRoles: state.roles.roles,
        allRules: state.rules.rules,
        rules: _.get(state, `rules.subjects['${ownProps.match.params.id}']`, [])
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRoles,
            loadRules,
            saveUser,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(UserDetail))
