import React from "react"
import classnames from "classnames"

import PersonPlaceholder from "../../public/person.svg"

import "./style.css"

const PatientImage = ({ data, big }) => {
    if (!data) {
        return null
    }

    return (
        <div className={classnames("patientImage", { big })}>
            <img src={PersonPlaceholder} alt="" />
        </div>
    )
}

export default PatientImage
