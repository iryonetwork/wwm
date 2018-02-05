import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import _ from "lodash"

import { loadUser } from "../../modules/users"
import { loadRoles } from "../../modules/roles"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            email: "",
            password: "",
            password2: ""
        }
    }

    componentDidMount() {
        this.props.loadUser(this.props.userID)
        this.props.loadRoles()
    }

    componentWillReceiveProps(props) {
        if (props.user) {
            this.setState({ email: props.user.email })
        }
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

    updateRole = roleID => e => {
        let roles = { ...this.state.roles }
        roles[roleID] = !roles[roleID]
        this.setState({
            roles: roles
        })
    }

    render() {
        let props = this.props
        if (!props.user) {
            return <div>Loading...</div>
        }
        return (
            <div>
                <h1>Users</h1>

                <h2>{props.user.username}</h2>

                <form>
                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input
                            type="email"
                            className="form-control"
                            id="email"
                            value={this.state.email}
                            onChange={this.updateEmail}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Change password</label>
                        <input
                            type="password"
                            className="form-control"
                            id="paswword"
                            value={this.state.password}
                            onChange={this.updatePassword}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password2">New password again</label>
                        <input
                            type="password"
                            className="form-control"
                            id="paswword2"
                            value={this.state.password2}
                            onChange={this.updatePassword2}
                        />
                    </div>
                </form>

                <h2>Roles</h2>
                {this.state.roles
                    ? _.map(this.state.roles, (selected, roleID) => {
                          return (
                              <div className="form-check" key={roleID}>
                                  <input
                                      className="form-check-input"
                                      type="checkbox"
                                      checked={selected}
                                      id={`role-checkbox-${roleID}`}
                                      onChange={this.updateRole(roleID)}
                                  />
                                  <label
                                      className="form-check-label"
                                      htmlFor={`role-checkbox-${roleID}`}
                                  >
                                      {props.allRoles[roleID].name}
                                  </label>
                              </div>
                          )
                      })
                    : null}
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        user: state.users.user,
        loading: state.users.loading,
        userID: ownProps.match.params.id,
        roles: _.get(state, `roles.users['${ownProps.match.params.id}']`, []),
        allRoles: state.roles.roles
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            loadRoles
        },
        dispatch
    )

export default withRouter(
    connect(mapStateToProps, mapDispatchToProps)(UserDetail)
)
