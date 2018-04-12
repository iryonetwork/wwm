import React from "react"
import { Link } from "react-router-dom"
import { connect } from "react-redux"

import UserRoles from "./list"

const UserRolesIndex = props => (
    <div>
        <h1>User roles</h1>
        <UserRoles />
    </div>
)

const mapStateToProps = state => ({
    forbidden: false
})

export default connect(mapStateToProps)(UserRolesIndex)
