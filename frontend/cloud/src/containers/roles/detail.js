import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"

import { loadUsers } from "../../modules/users"
import { loadRules } from "../../modules/rules"
import { addUserToRole } from "../../modules/roles"
import { open, COLOR_DANGER } from "../../modules/alert"
import Rules from "../rules"
import Users from "../users/list"

class DetailRole extends React.Component {
    constructor(props) {
        super(props)

        this.state = {}
    }
    componentDidMount() {
        if (!this.props.users) {
            this.props.loadUsers()
        }
        this.props.loadRules()
    }

    addUser = () => e => {
        if (this.state.selectedUser) {
            this.props.addUserToRole(this.props.roleID, this.state.selectedUser)
        } else {
            this.props.open("You need to select user", "", COLOR_DANGER)
        }
    }

    changeSelectedUser = () => e => {
        this.setState({ selectedUser: e.target.value })
    }

    render() {
        let props = this.props
        return (
            <div>
                <header>
                    <h3>Users in {props.role.name}</h3>
                </header>
                <Users users={props.role.users} role={props.roleID} />
                <div className="input-group">
                    <select value={this.state.selectedUser} onChange={this.changeSelectedUser()} className="custom-select form-control-sm">
                        <option>Select user...</option>
                        {props.users
                            ? map(props.users, (user, userID) => (
                                  <option value={userID} key={userID}>
                                      {user.username} - {user.email}
                                  </option>
                              ))
                            : null}
                    </select>
                    <div className="input-group-append">
                        <button onClick={this.addUser()} className="btn btn-sm btn-outline-secondary" type="button">
                            Add user to role
                        </button>
                    </div>
                </div>

                <Rules rules={props.rules} subject={props.roleID} />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        role: state.roles.roles[ownProps.match.params.id],
        users: state.users.users,
        rules: state.rules.subjects ? state.rules.subjects[ownProps.match.params.id] || [] : [],
        roleID: ownProps.match.params.id
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRules,
            addUserToRole,
            open
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(DetailRole)
