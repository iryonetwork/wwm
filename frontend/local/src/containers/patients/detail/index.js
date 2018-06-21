import React from "react"
import { connect } from "react-redux"
import { Link, NavLink, Route, Redirect } from "react-router-dom"
import _ from "lodash"

import Consultation from "../../waitlist/detail"
import Data from "./data"
import History from "./history"
import Personal from "./personal"
import HealthRecord from "./record"

import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"

import {
    RESOURCE_PATIENT_IDENTIFICATION,
    RESOURCE_DEMOGRAPHIC_INFORMATION,
    RESOURCE_VITAL_SIGNS,
    RESOURCE_HEALTH_HISTORY,
    RESOURCE_EXAMINATION,
    READ,
    WRITE
} from "../../../modules/validations"
import { fetchPatient, fetchHealthRecords } from "../../../modules/patient"
import { get as getWaitlistItem } from "../../../modules/waitlist"

import { ReactComponent as InConsultationIcon } from "shared/icons/in-consultation-active.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/medical-data-active.svg"
import { ReactComponent as MedicalHistoryIcon } from "shared/icons/medical-history-active.svg"
import { ReactComponent as PersonalInfoIcon } from "shared/icons/personal-info-active.svg"
import { ReactComponent as HealthRecordIcon } from "shared/icons/health-record-active.svg"
import { ReactComponent as AddIcon } from "shared/icons/add.svg"

import "./style.css"

class PatientDetail extends React.Component {
    constructor(props) {
        super(props)

        if (props.match.params.waitlistID && props.match.params.itemID) {
            props
                .getWaitlistItem(props.match.params.waitlistID, props.match.params.itemID)
                .then(item => {
                    props.fetchPatient(item.patientID)
                    props.fetchHealthRecords(item.patientID)
                })
                .catch(ex => console.warn("Failed to load patient", ex))
        } else {
            props.fetchPatient(props.match.params.patientID).catch(ex => console.warn("Failed to load patient", ex))
            props.fetchHealthRecords(props.match.params.patientID)
        }
    }

    componentDidUpdate(prevProps) {
        if (!this.props.patientLoading) {
            if (this.props.match.params.waitlistID && this.props.match.params.itemID && this.props.match.params.itemID !== prevProps.match.params.itemID) {
                this.props
                    .getWaitlistItem(this.props.match.params.waitlistID, this.props.match.params.itemID)
                    .then(item => {
                        this.props.fetchPatient(item.patientID)
                        this.props.fetchHealthRecords(item.patientID)
                    })
                    .catch(ex => console.warn("Failed to load patient", ex))
            } else if (this.props.match.params.patientID !== prevProps.match.params.patientID) {
                this.props.fetchPatient(this.props.match.params.patientID).catch(ex => console.warn("Failed to load patient", ex))
                this.props.fetchHealthRecords(this.props.match.params.patientID)
            }
        }
    }

    render() {
        const { inConsultation, waitlistFetching, patientLoading, patientRecordsLoading, patient, bodyMeasurements, match } = this.props

        const waitlistID = match.params.waitlistID
        const waitlistItemID = match.params.itemID
        const patientID = match.params.patientID

        const baseURL = inConsultation ? `/waitlist/${waitlistID}/${waitlistItemID}/` : `/patients/${patientID}/`

        if (waitlistFetching || patientLoading || patientRecordsLoading) {
            return <Spinner />
        }

        return this.props.canSeePatientId ? (
            <div className="waitlist-detail">
                <div className="sidebar">
                    <Patient big={true} data={patient} />
                    {this.props.canSeeVitalSigns ? (
                        <div className="row measurements">
                            {bodyMeasurements && bodyMeasurements.height ? (
                                <div className="col-sm-4">
                                    <h5>Height</h5>
                                    {bodyMeasurements.height} cm
                                </div>
                            ) : (
                                patient.heightAtBirth && (
                                    <div className="col-sm-4">
                                        <h5>Height</h5>
                                        {patient.heightAtBirth} cm
                                    </div>
                                )
                            )}
                            {bodyMeasurements && bodyMeasurements.weight ? (
                                <div className="col-sm-4">
                                    <h5>Weight</h5>
                                    {bodyMeasurements.weight} kg
                                </div>
                            ) : (
                                patient.weightAtBirth && (
                                    <div className="col-sm-4">
                                        <h5>Weight</h5>
                                        {patient.weightAtBirth} grams
                                    </div>
                                )
                            )}
                            {bodyMeasurements &&
                                bodyMeasurements.bmi && (
                                    <div className="col-sm-4">
                                        <h5>BMI</h5>
                                        {bodyMeasurements.bmi}
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
                                                    In Consultation
                                                </NavLink>
                                            ) : (
                                                <span>
                                                    <InConsultationIcon />
                                                    In Consultation
                                                </span>
                                            )}
                                        </li>
                                    )}
                                    <li>
                                        {this.props.canSeeVitalSigns ? (
                                            <NavLink to={baseURL + "data"}>
                                                <MedicalDataIcon />
                                                Medical Data
                                            </NavLink>
                                        ) : null}
                                    </li>
                                    <li>
                                        {this.props.canSeeHealthHistory ? (
                                            <NavLink to={baseURL + "history"}>
                                                <MedicalHistoryIcon />
                                                Medical History
                                            </NavLink>
                                        ) : null}
                                    </li>
                                    {this.props.canSeeExamination && (
                                        <li>
                                            <NavLink to={baseURL + "record"}>
                                                <HealthRecordIcon />
                                                Medical Record
                                            </NavLink>
                                        </li>
                                    )}
                                    <li>
                                        {this.props.canSeeDemographicInformation ? (
                                            <NavLink to={baseURL + "personal"}>
                                                <PersonalInfoIcon />
                                                Personal Info
                                            </NavLink>
                                        ) : null}
                                    </li>
                                </ul>
                            </div>
                        </div>
                    ) : null}
                    {inConsultation && (this.props.canAddExamination || this.props.canAddVitalSigns) ? (
                        <div className="row menu">
                            <div className="col-sm-12">
                                <h4>Actions</h4>

                                <ul>
                                    {this.props.canAddExamination && (
                                        <li>
                                            <Link className="btn btn-link" to={baseURL + "consultation/add-diagnosis"}>
                                                <AddIcon />
                                                Add Diagnosis
                                            </Link>
                                        </li>
                                    )}
                                    {this.props.canAddVitalSigns && (
                                        <li>
                                            <Link className="btn btn-link" to={baseURL + "consultation/add-data"}>
                                                <AddIcon />
                                                Add Medical Data
                                            </Link>
                                        </li>
                                    )}
                                </ul>
                            </div>
                        </div>
                    ) : null}
                </div>
                <div className="container">
                    {this.props.canSeePatientId && <Route exact path="/patients/:patientID" render={() => <Redirect to={baseURL + "personal"} />} />}
                    {this.props.canSeeExamination && <Route path="/waitlist/:waitlistID/:itemID/consultation" component={Consultation} />}
                    {this.props.canSeeVitalSigns && <Route path={match.path + "/data"} component={Data} />}
                    {this.props.canSeeDemographicInformation && <Route path={match.path + "/personal"} component={Personal} />}
                    {this.props.canSeeExamination && <Route path={match.path + "/record"} component={HealthRecord} />}
                    {this.props.canSeeHealthHistory && <Route path={match.path + "/history"} component={History} />}
                </div>
            </div>
        ) : null
    }
}

PatientDetail = connect(
    (state, props) => {
        let inConsultation = props.match.params.waitlistID && props.match.params.itemID

        // get latest body measurements
        let bodyMeasurements = {}
        // if in consultation, fetch data from current consultation first
        if (inConsultation && state.waitlist.item) {
            bodyMeasurements = {
                height: state.waitlist.item.vitalSigns && state.waitlist.item.vitalSigns.height ? state.waitlist.item.vitalSigns.height.value : undefined,
                weight: state.waitlist.item.vitalSigns && state.waitlist.item.vitalSigns.weight ? state.waitlist.item.vitalSigns.weight.value : undefined,
                bmi: state.waitlist.item.vitalSigns && state.waitlist.item.vitalSigns.bmi ? state.waitlist.item.vitalSigns.bmi.value : undefined
            }
        }

        let patientRecordsNeeded = bodyMeasurements.height && bodyMeasurements.weight && bodyMeasurements.bmi ? false : true

        if (patientRecordsNeeded && state.patient.patientRecords.data) {
            let records = state.patient.patientRecords.data
            // sort records by creation time and reverse to have latest record as first
            records = _.reverse(
                _.sortBy(records, [
                    function(obj) {
                        return obj.meta.created
                    }
                ])
            )

            // collect latest data for each category
            _.forEach(records, ({ data, meta }) => {
                _.forEach(data.vitalSigns, (obj, key) => {
                    if (!bodyMeasurements.height && key === "height") {
                        bodyMeasurements.height = obj.value
                    }
                    if (!bodyMeasurements.weight && key === "weight") {
                        bodyMeasurements.weight = obj.value
                    }
                    if (!bodyMeasurements.bmi && key === "bmi") {
                        bodyMeasurements.bmi = obj.value
                    }
                    if (bodyMeasurements.height && bodyMeasurements.weight && bodyMeasurements.bmi) {
                        return false
                    }
                })
            })
        }

        return {
            inConsultation: inConsultation,
            waitlistFetching: state.waitlist.fetching,
            patientLoading: state.patient.loading || state.patient.saving,
            patient: state.patient.patient,
            patientID: props.match.params.patientID || state.patient.patient.ID,
            loadedPatientID: state.patient.patient.ID,
            patientRecordsNeeded: patientRecordsNeeded,
            patientRecordsLoading: state.patient.patientRecords.loading,
            patientRecords: state.patient.patientRecords.data,
            bodyMeasurements: bodyMeasurements,
            waitlistItem: state.waitlist.item || {},
            canSeePatientId: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
            canSeeDemographicInformation: ((state.validations.userRights || {})[RESOURCE_DEMOGRAPHIC_INFORMATION] || {})[READ],
            canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE],
            canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
            canSeeHealthHistory: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[READ],
            canAddExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[WRITE],
            canSeeExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ]
        }
    },
    {
        fetchPatient,
        getWaitlistItem,
        fetchHealthRecords
    }
)(PatientDetail)

export default PatientDetail
