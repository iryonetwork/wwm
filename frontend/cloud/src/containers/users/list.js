import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadUsers, deleteUser } from "../../modules/users"

class Users extends React.Component {
    componentDidMount() {
        this.props.loadUsers()
    }

    removeUser = userID => e => {
        this.props.deleteUser(userID)
    }

    getName(user) {
        if (user.personalData !== undefined) {
            var name = ""
            if (user.personalData.firstName !== undefined && user.personalData.firstName !== "") {
                name += user.personalData.firstName
            }
            if (user.personalData.middleName !== undefined && user.personalData.middleName !== "") {
                name = name + " " + user.personalData.middleName
            }
            if (user.personalData.lastName !== undefined && user.personalData.lastName !== "") {
                name = name + " " + user.personalData.lastName
            }
            return name
        }

        return "Unknown"
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
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
                        <th scope="col">Name</th>
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
                            <td>{this.getName(user)}</td>
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
    loading: state.users.loading,
    forbidden: state.users.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            deleteUser
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Users)
