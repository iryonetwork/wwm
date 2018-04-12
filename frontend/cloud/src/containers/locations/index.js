import React from "react"
import { Link } from "react-router-dom"
import { connect } from "react-redux"

import Locations from "./list"

const LocationsIndex = props => (
    <div>
        <h1>Locations</h1>
        <Locations />
    </div>
)

const mapStateToProps = state => ({
    forbidden: false
})

export default connect(mapStateToProps)(LocationsIndex)
