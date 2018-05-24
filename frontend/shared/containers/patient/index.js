import React from "react"
import classnames from "classnames"
import moment from "moment"

import PersonPlaceholder from "../../public/person.svg"

import "./style.css"

export default ({ data, style, big }) => {
    if (!data) {
        return null
    }

    const gender = data.gender === "CODED-at0310" ? "M" : data.gender === "CODED-at0311" ? "F" : "?"
    const dob = moment(data.dateOfBirth)
    const dobString = dob.toDate().toLocaleDateString()
    const age = moment().diff(dob, "years")

    return (
        <div className={classnames("patientCard", { big })}>
            <img src={PersonPlaceholder} alt="" />
            <div>
                <div className="name">
                    {data.lastName}, {data.firstName}
                </div>
                <div className="dob">
                    {dobString}
                    <span className="age">{age} y</span>
                    {gender}
                </div>
            </div>
        </div>
    )
}
