import React from "react"
import classnames from "classnames"

import PersonPlaceholder from "../../public/person.svg"

import "./style.css"

export default ({ style, big }) => (
    <div className={classnames("patientCard", { big })}>
        <img src={PersonPlaceholder} alt="" />
        <div>
            <div className="name">Graves, Alma</div>
            <div className="dob">
                3 Jun 1994
                <span className="age">24 y</span>
                F
            </div>
        </div>
    </div>
)
