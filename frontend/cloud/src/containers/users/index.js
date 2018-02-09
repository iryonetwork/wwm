import React from "react"
import { Link } from "react-router-dom"

import Users from "./list"

export default props => (
    <div>
        <h1>Users</h1>
        <Users />
        <Link to="/users/new" className="btn btn-sm btn-outline-secondary">
            Add new users
        </Link>
    </div>
)
