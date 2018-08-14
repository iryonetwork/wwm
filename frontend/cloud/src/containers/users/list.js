import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadUsers, deleteUser } from "../../modules/users"
import { getName } from "../../utils/user"

class Users extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.users) {
            this.props.loadUsers()
        }
        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.users && !nextProps.usersLoading) {
            this.props.loadUsers()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.users || props.usersLoading
        this.setState({ loading: loading })
    }

    removeUser(userID) {
        return e => {
            this.props.deleteUser(userID)
        }
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        let i = 0
        return (
            <table className="table">
                <thead>
                    <tr>
                        <th className="w-7" scope="col">
                            #
                        </th>
                        <th scope="col">Username</th>
                        <th scope="col">Name</th>
                        <th scope="col">Email</th>
                        <th />
                    </tr>
                </thead>
                <tbody>
                    {_.map(_.filter(props.users, user => user), user => (
                        <tr key={user.id}>
                            <th className="w-7" scope="row">
                                {++i}
                            </th>
                            <td>
                                <Link to={`/users/${user.id}`}>{user.username}</Link>
                            </td>
                            <td>{getName(user)}</td>
                            <td>{user.email}</td>
                            <td className="text-right">
                                <button onClick={this.removeUser(user.id)} className="btn btn-link" type="button">
                                    <span className="remove-link">Remove</span>
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
    users: ownProps.users
        ? state.users.allLoaded
            ? _.fromPairs(_.map(ownProps.users, userID => [userID, state.users.users[userID]]))
            : undefined
        : state.users.allLoaded
            ? state.users.users
            : undefined,
    usersLoading: state.users.loading,
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
