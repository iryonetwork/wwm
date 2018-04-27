import React from "react"
import { NavLink, Route, Redirect } from "react-router-dom"

import Consultation from "../../waitlist/detail"
import Data from "./data"
import History from "./history"
import Personal from "./personal"

import Patient from "shared/containers/patient"

import { ReactComponent as InConsultationIcon } from "shared/icons/in-consultation-active.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/medical-data-active.svg"
import { ReactComponent as MedicalHistoryIcon } from "shared/icons/medical-history-active.svg"
import { ReactComponent as PersonalInfoIcon } from "shared/icons/personal-info-active.svg"

import "./style.css"

export default ({ match }) => {
    const inConsultation = true
    const waitlistID = match.params.waitlistID
    const waitlistItemID = match.params.itemID
    const patientID = match.params.patientID

    const baseURL = inConsultation ? `/waitlist/${waitlistID}/${waitlistItemID}/` : `/patients/${patientID}/`

    return (
        <div className="waitlist-detail">
            <div className="sidebar">
                <Patient big={true} />
                <div className="row measurements">
                    <div className="col-sm-4">
                        <h5>Height</h5>
                        1,57m
                    </div>
                    <div className="col-sm-4">
                        <h5>Weight</h5>
                        54kg
                    </div>
                    <div className="col-sm-4">
                        <h5>BMI</h5>
                        22.2
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-12">
                        <h5>Allergies</h5>
                        <span className="danger">Peanuts</span>
                    </div>
                </div>

                <div className="row">
                    <div className="col-sm-12">
                        <h5>Chronic diseases</h5>
                        <span className="danger">Asthma</span>
                    </div>
                </div>

                <div className="row menu">
                    <div className="col-sm-12">
                        <h4>Menu</h4>

                        <ul>
                            {inConsultation && (
                                <li>
                                    <NavLink to={baseURL + "consultation"}>
                                        <InConsultationIcon />
                                        In consultation
                                    </NavLink>
                                </li>
                            )}
                            <li>
                                <NavLink to={baseURL + "data"}>
                                    <MedicalDataIcon />
                                    Medical data
                                </NavLink>
                            </li>
                            <li>
                                <NavLink to={baseURL + "history"}>
                                    <MedicalHistoryIcon />
                                    Medical history
                                </NavLink>
                            </li>
                            <li>
                                <NavLink to={baseURL + "personal"}>
                                    <PersonalInfoIcon />
                                    Personal info
                                </NavLink>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <div className="container">
                <Route exact path={match.url} render={() => <Redirect to={inConsultation ? baseURL + "consultation" : baseURL + "personal"} />} />
                {inConsultation && <Route path="/patients/:patientID" render={() => <Redirect to={baseURL + "consultation"} />} />}
                <Route path={match.url + "/consultation"} component={Consultation} />
                <Route path={match.url + "/data"} component={Data} />
                <Route path={match.url + "/personal"} component={Personal} />
                <Route path={match.url + "/history"} component={History} />
            </div>
        </div>
    )
}
