import React from "react"
import { Link } from "react-router-dom"
import { connect } from "react-redux"

import Users from "./list"

const UsersIndex = props => (
    <div>
        <h1>Users</h1>
        <Users />
        {props.forbidden ? null : (
            <Link to="/users/new" className="btn btn-sm btn-outline-secondary">
                Add new user
            </Link>
        )}
    </div>
)

const mapStateToProps = state => ({
    forbidden: state.users.forbidden
})

export default connect(mapStateToProps)(UsersIndex)
