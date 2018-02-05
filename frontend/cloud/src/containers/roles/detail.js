import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"

import { loadUsers } from "../../modules/users"
import { loadRules } from "../../modules/rules"
import Rules from "../rules"

class DetailRole extends React.Component {
    componentDidMount() {
        this.props.loadUsers()
        this.props.loadRules()
    }

    render() {
        let props = this.props
        return (
            <div>
                <header>
                    <h3>Users in {props.role.name}</h3>
                </header>
                <table className="table table-hover">
                    <thead>
                        <tr>
                            <th scope="col">#</th>
                            <th scope="col">Username</th>
                            <th scope="col">Email</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {props.users
                            ? map(props.role.users, (userID, i) => (
                                  <tr key={userID}>
                                      <th scope="row">{i + 1}</th>
                                      <td>{props.users[userID].username}</td>
                                      <td>{props.users[userID].email}</td>
                                      <td className="text-right">
                                          <button
                                              className="btn btn-sm btn-light"
                                              type="button"
                                          >
                                              <span className="oi oi-trash" />
                                          </button>
                                      </td>
                                  </tr>
                              ))
                            : null}
                    </tbody>
                </table>
                <div className="input-group">
                    <select className="custom-select">
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
                        <button
                            className="btn btn-outline-secondary"
                            type="button"
                        >
                            Add user to role
                        </button>
                    </div>
                </div>

                <Rules rules={props.rules} />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        role: state.roles.roles[ownProps.match.params.id],
        users: state.users.users,
        rules: state.rules.subjects
            ? state.rules.subjects[ownProps.match.params.id] || []
            : []
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRules
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(DetailRole)
