import React from "react"
import classnames from "classnames"

import PersonPlaceholder from "../../public/person.svg"

import "./style.css"

export default ({ data, style, big }) => {
    if (!data) {
        return null
    }

    const gender = data.gender === 'CODED-at0310' ? 'M' : (data.gender === 'CODED-at0311' ? 'F' : '?')
    const age = Math.floor((Date.now() - (new Date(data.dateOfBirth)).getTime()) / (1000 * 60 * 60 * 24 * 365))
    const dob = (new Date(data.dateOfBirth)).toLocaleDateString()
    return (
        <div className={classnames("patientCard", { big })}>
            <img src={PersonPlaceholder} alt="" />
            <div>
                <div className="name">{data.lastName}, {data.firstName}</div>
                <div className="dob">
                    {dob}
                    <span className="age">{age} y</span>
                    {gender}
                </div>
            </div>
        </div>
    )
}
