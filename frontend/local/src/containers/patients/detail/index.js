import React from "react"
import { connect } from "react-redux"
import { NavLink, Route, Redirect } from "react-router-dom"

import Consultation from "../../waitlist/detail"
import Data from "./data"
import History from "./history"
import Personal from "./personal"

import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"

import {
    RESOURCE_PATIENT_IDENTIFICATION,
    RESOURCE_DEMOGRAPHIC_INFORMATION,
    RESOURCE_VITAL_SIGNS,
    RESOURCE_HEALTH_HISTORY,
    RESOURCE_WAITLIST,
    READ
} from "../../../modules/validations"
import { fetchPatient } from "../../../modules/patient"
import { get as getWaitlistItem } from "../../../modules/waitlist"

import { ReactComponent as InConsultationIcon } from "shared/icons/in-consultation-active.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/medical-data-active.svg"
import { ReactComponent as MedicalHistoryIcon } from "shared/icons/medical-history-active.svg"
import { ReactComponent as PersonalInfoIcon } from "shared/icons/personal-info-active.svg"

import "./style.css"

class PatientDetail extends React.Component {
    constructor(props) {
        super(props)

        if (props.match.params.waitlistID && props.match.params.itemID) {
            props
                .getWaitlistItem(props.match.params.waitlistID, props.match.params.itemID)
                .then(item => props.fetchPatient(item.patientID))
                .catch(ex => console.warn("Failed to load patient", ex))
        } else {
            props.fetchPatient(props.match.params.patientID).catch(ex => console.warn("Failed to load patient", ex))
        }
    }

    render() {
        const { waitlistFetching, patientLoading, waitlistItem, patient, match } = this.props

        const waitlistID = match.params.waitlistID
        const waitlistItemID = match.params.itemID
        const patientID = match.params.patientID
        const inConsultation = waitlistID && waitlistItemID

        const baseURL = inConsultation ? `/waitlist/${waitlistID}/${waitlistItemID}/` : `/patients/${patientID}/`

        if (waitlistFetching || patientLoading) {
            return <Spinner />
        }

        return this.props.canSeePatientId ? (
            <div className="waitlist-detail">
                <div className="sidebar">
                    <Patient big={true} data={patient} />
                    {this.props.canSeeVitalSigns ? (
                        <div className="row measurements">
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.height && (
                                    <div className="col-sm-4">
                                        <h5>Height</h5>
                                        {waitlistItem.vitalSigns.height}cm
                                    </div>
                                )}
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.weight && (
                                    <div className="col-sm-4">
                                        <h5>Weight</h5>
                                        {waitlistItem.vitalSigns.weight}kg
                                    </div>
                                )}
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.height &&
                                waitlistItem.vitalSigns.weight && (
                                    <div className="col-sm-4">
                                        <h5>BMI</h5>
                                        {round(waitlistItem.vitalSigns.weight / waitlistItem.vitalSigns.height / waitlistItem.vitalSigns.height * 10000, 2)}
                                    </div>
                                )}
                        </div>
                    ) : null}

                    {this.props.canSeeHealthHistory ? (
                        <div>
                            {patient.allergies &&
                                patient.allergies.length > 0 && (
                                    <div className="row">
                                        <div className="col-sm-12">
                                            <h5>Allergies</h5>
                                            {patient.allergies.map((item, i) => (
                                                <div key={i}>
                                                    <span className="danger">{item.allergy}</span>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                            {patient.chronicDiseases &&
                                patient.chronicDiseases.length > 0 && (
                                    <div className="row">
                                        <div className="col-sm-12">
                                            <h5>Chronic diseases</h5>
                                            {patient.chronicDiseases.map((item, i) => (
                                                <div key={i}>
                                                    <span className="danger">{item.disease}</span>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                        </div>
                    ) : null}
                    {this.props.canSeeDemographicInformation ||
                    this.props.canSeeVitalSigns ||
                    this.props.canSeeHealthHistory ||
                    this.props.canSeeExamination ? (
                        <div className="row menu">
                            <div className="col-sm-12">
                                <h4>Menu</h4>

                                <ul>
                                    {inConsultation && (
                                        <li>
                                            {this.props.canSeeVitalSigns ? (
                                                <NavLink to={baseURL + "consultation"}>
                                                    <InConsultationIcon />
                                                    In consultation
                                                </NavLink>
                                            ) : (
                                                <span>
                                                    <InConsultationIcon />
                                                    In consultation
                                                </span>
                                            )}
                                        </li>
                                    )}
                                    {/* <li>
                                        {this.props.canSeeVitalSigns ? (
                                            <NavLink to={baseURL + "data"}>
                                                <MedicalDataIcon />
                                                Medical data
                                            </NavLink>
                                        ) : null}
                                    </li> */}
                                    <li>
                                        {this.props.canSeeHealthHistory ? (
                                            <NavLink to={baseURL + "history"}>
                                                <MedicalHistoryIcon />
                                                Medical history
                                            </NavLink>
                                        ) : null}
                                    </li>
                                    <li>
                                        {this.props.canSeeDemographicInformation ? (
                                            <NavLink to={baseURL + "personal"}>
                                                <PersonalInfoIcon />
                                                Personal info
                                            </NavLink>
                                        ) : null}
                                    </li>
                                </ul>
                            </div>
                        </div>
                    ) : null}
                </div>
                <div className="container">
                    {this.props.canSeePatientId && <Route exact path="/patients/:patientID" render={() => <Redirect to={baseURL + "personal"} />} />}
                    {this.props.canSeeWaitlist && <Route path="/waitlist/:waitlistID/:itemID/consultation" component={Consultation} />}
                    {this.props.canSeeVitalSigns && <Route path={match.path + "/data"} component={Data} />}
                    {this.props.canSeeDemographicInformation && <Route path={match.path + "/personal"} component={Personal} />}
                    {this.props.canSeeHealthHistory && <Route path={match.path + "/history"} component={History} />}
                </div>
            </div>
        ) : null
    }
}

PatientDetail = connect(
    (state, props) => ({
        waitlistFetching: state.waitlist.fetching,
        patientLoading: state.patient.loading,
        patient: state.patient.patient,
        waitlistItem: state.waitlist.item || {},
        canSeePatientId: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ],
        canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
        canSeeHealthHistory: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[READ],
        canSeeWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[READ]
    }),
    {
        fetchPatient,
        getWaitlistItem
    }
)(PatientDetail)

export default PatientDetail

const round = (number, precision) => {
    var shift = function(number, precision) {
        var numArray = ("" + number).split("e")
        return +(numArray[0] + "e" + (numArray[1] ? +numArray[1] + precision : precision))
    }
    return shift(Math.round(shift(number, +precision)), -precision)
}
