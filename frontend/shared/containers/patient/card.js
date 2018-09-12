import React from "react"
import classnames from "classnames"
import moment from "moment"

import PatientImage from "./image"

import "./style.css"

const PatientCard = ({ data, big, withoutImage }) => {
    if (!data) {
        return null
    }

    const gender = data.gender === "CODED-at0310" ? "M" : data.gender === "CODED-at0311" ? "F" : "?"
    const dob = moment(data.dateOfBirth)
    const dobString = dob.format("Do MMM Y")
    const ageYears = moment().diff(dob, "years")
    const ageMonths = moment().diff(dob, "months")
    const ageWeeks = moment().diff(dob, "weeks")

    return (
        <div className={classnames("patientCard", { big })}>
            {withoutImage ? null : <PatientImage data={data} big={big} />}
            <div>
                <div className="name">
                    {data.lastName}, {data.firstName}
                </div>
                <div className="dob">
                    {dobString}
                    <span className="age">{ageYears < 2 ? (ageMonths < 3 ? `${ageWeeks} w` : `${ageMonths} m`) : `${ageYears} y`}</span>
                    {gender}
                </div>
            </div>
        </div>
    )
}

export default PatientCard
