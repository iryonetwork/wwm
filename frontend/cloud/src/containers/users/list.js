import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadUsers, deleteUser } from "../../modules/users"
import { removeUserFromRole } from "../../modules/roles"

class Users extends React.Component {
    componentDidMount() {
        this.props.loadUsers()
    }

    removeUser = userID => e => {
        if (this.props.roleID) {
            this.props.removeUserFromRole(this.props.roleID, userID)
        } else {
            this.props.deleteUser(userID)
        }
    }

    render() {
        let props = this.props
        if (props.loading) {
            return <div>Loading...</div>
        }
        let i = 0
        return (
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
                    {_.map(_.filter(props.users, user => user), user => (
                        <tr key={user.id}>
                            <th scope="row">{++i}</th>
                            <td>
                                <Link to={`/users/${user.id}`}>{user.username}</Link>
                            </td>
                            <td>{user.email}</td>
                            <td className="text-right">
                                <button onClick={this.removeUser(user.id)} className="btn btn-sm btn-light" type="button">
                                    <span className="icon_trash" />
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    users:
        (ownProps.users ? (state.users.users ? _.fromPairs(_.map(ownProps.users, userID => [userID, state.users.users[userID]])) : {}) : state.users.users) ||
        {},
    roleID: ownProps.role,
    loading: state.users.loading
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            removeUserFromRole,
            deleteUser
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Users)
