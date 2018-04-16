import React from "react"
import { Link } from "react-router-dom"
import { connect } from "react-redux"

import Organizations from "./list"

const OrganizationsIndex = props => (
    <div>
        <h1>Organizations</h1>
        <Organizations />
        {props.forbidden ? null : (
            <Link to="/organizations/new" className="btn btn-sm btn-outline-primary col">
                Add new organization
            </Link>
        )}
    </div>
)

const mapStateToProps = state => ({
    forbidden: false
})

export default connect(mapStateToProps)(OrganizationsIndex)
