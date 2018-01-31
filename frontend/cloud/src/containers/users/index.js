import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"

import { loadUsers } from "../../modules/users"

class Users extends React.Component {
    componentDidMount() {
        this.props.loadUsers()
    }

    render() {
        let props = this.props
        if (props.loading) {
            return <div>Loading...</div>
        }
        return (
            <div>
                <h1>Users</h1>
                <table className="table">
                    <thead>
                        <tr>
                            <th scope="col">#</th>
                            <th scope="col">Username</th>
                            <th scope="col">Email</th>
                            <th scope="col" />
                        </tr>
                    </thead>
                    <tbody>
                        {props.users.map((user, i) => (
                            <tr key={user.id}>
                                <th scope="row">{i + 1}</th>
                                <td>
                                    <Link to={`/users/${user.id}`}>
                                        {user.username}
                                    </Link>
                                </td>
                                <td>{user.email}</td>
                                <td />
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        )
    }
}

const mapStateToProps = state => ({
    users: state.users.users || [],
    loading: state.users.loading
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Users)
