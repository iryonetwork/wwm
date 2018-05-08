import React from "react"
import { connect } from "react-redux"
import { NavLink, Route, Redirect } from "react-router-dom"

import Consultation from "../../waitlist/detail"
import Data from "./data"
import History from "./history"
import Personal from "./personal"

import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"

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
        props
            .getWaitlistItem(props.match.params.waitlistID, props.match.params.itemID)
            .then(item => props.fetchPatient(item.patient_id))
            .catch(ex => console.warn("Failed to load patient", ex))
    }

    render() {
        const { waitlistFetching, patientLoading, patient, match } = this.props

        const inConsultation = true
        const waitlistID = match.params.waitlistID
        const waitlistItemID = match.params.itemID
        const patientID = match.params.patientID

        const baseURL = inConsultation ? `/waitlist/${waitlistID}/${waitlistItemID}/` : `/patients/${patientID}/`

        if (waitlistFetching || patientLoading) {
            return <Spinner />
        }

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
                    <Route path={match.url + "/personal"} component={Personal} patient={patient} />
                    <Route path={match.url + "/history"} component={History} />
                </div>
            </div>
        )
    }
}

PatientDetail = connect(
    state => ({
        waitlistFetching: state.waitlist.fetching,
        patientLoading: state.patient.loading,
        patient: state.patient.patient
    }),
    {
        fetchPatient,
        getWaitlistItem
    }
)(PatientDetail)

export default PatientDetail
