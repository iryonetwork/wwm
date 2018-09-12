import React from "react"
import classnames from "classnames"
import moment from "moment"

import PersonPlaceholder from "../../public/person.svg"

import "./style.css"

const PatientImage = ({ data, style, big }) => {
    if (!data) {
        return null
    }

    return (
        <div className={classnames("patientImage", { big })}>
            <img src={PersonPlaceholder} alt="" />
        </div>
    )
}

const PatientCard = ({ data, style, big, withoutImage }) => {
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
            {withoutImage ? null : (
                <div className={classnames("patientImage", { big })}>
                    <img src={PersonPlaceholder} alt="" />
                </div>
            )}
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

export { PatientImage, PatientCard }
export default PatientCard
