import React from "react"
import { Link } from "react-router-dom"
import { connect } from "react-redux"

import Clinics from "./list"

const ClinicsIndex = props => (
    <div>
        <h1>Clinics</h1>
        <Clinics />
    </div>
)

const mapStateToProps = state => ({
    forbidden: false
})

export default connect(mapStateToProps)(ClinicsIndex)
