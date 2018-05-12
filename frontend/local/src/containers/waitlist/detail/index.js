import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"

import "./style.css"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as LaboratoryIcon } from "shared/icons/laboratory.svg"
import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"
import { ReactComponent as PositiveIcon } from "shared/icons/positive.svg"

import { RESOURCE_PATIENT_IDENTIFICATION, RESOURCE_VITAL_SIGNS, RESOURCE_EXAMINATION, RESOURCE_LABORATORY_TEST, RESOURCE_WAITLIST, READ, WRITE, UPDATE, DELETE } from "../../../modules/validations"
import { joinPaths } from "shared/utils"

import MedicalData from "./add-data"
import LaboratoryTest from "./add-lab-test"
import EditComplaint from "./edit-complaint"
import AddDiagnosis from "./add-diagnosis"
import Remove from "./remove"

class InConsultation extends React.Component {
    render() {
        const { match, waitlistItem, waitlistFetching, patient } = this.props
        if (waitlistFetching) {
            return null
        }
        return (
            <div>
                <header>
                    <h1>In consultation</h1>
                    {this.props.canAddDiagnosis && (
                        <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-primary btn-wide">
                            Add diagnosis
                        </Link>
                    )}
                </header>

                {this.props.canSeeMainComplaint && (
                    <div className="section">
                        <header>
                            <h2>
                                <ComplaintIcon />Main Complaint
                            </h2>
                            {this.props.canEditMainComplaint && (
                                <Link to={joinPaths(match.url, "edit-complaint")} className="btn btn-link">
                                    Edit main complaint
                                </Link>
                            )}
                        </header>

                        <h3>{waitlistItem.mainComplaint.complaint}</h3>
                        {waitlistItem.mainComplaint.comment && <p>{waitlistItem.mainComplaint.comment}</p>}
                    </div>
                )}

                {this.props.canSeeVitalSigns && (
                    <div className="section">
                        <header>
                            <h2>
                                <MedicalDataIcon />
                                Medical Data
                            </h2>
                            {this.props.canAddVitalSigns && (
                                <Link to={joinPaths(match.url, "add-data")} className="btn btn-link">
                                    Add medical data
                                </Link>
                            )}
                        </header>
                        <div className="card-group">
                            <div className="col-md-5 col-lg-4 col-xl-3">
                                <div className="card">
                                    <div className="card-header">Height</div>
                                    <div className="card-body">
                                        <div className="card-text">
                                            <p>
                                                <span className="big">1.56</span>m
                                            </p>
                                            <p>5ft 1in</p>
                                        </div>
                                    </div>
                                    <div className="card-footer">5 feb 2018</div>
                                </div>
                            </div>

                            <div className="col-md-5 col-lg-4 col-xl-3">
                                <div className="card">
                                    <div className="card-header">Body mass</div>
                                    <div className="card-body">
                                        <div className="card-text">
                                            <p>
                                                <span className="big">54.4</span>kg
                                            </p>
                                            <p>1008.8 lb</p>
                                        </div>
                                    </div>
                                    <div className="card-footer">5 feb 2018</div>
                                </div>
                            </div>

                            <div className="col-md-5 col-lg-4 col-xl-3">
                                <div className="card">
                                    <div className="card-header">BMI</div>
                                    <div className="card-body">
                                        <div className="card-text">
                                            <p>
                                                <span className="big">22.2</span>
                                            </p>
                                        </div>
                                    </div>
                                    <div className="card-footer">5 feb 2018</div>
                                </div>
                            </div>

                            <div className="col-md-5 col-lg-4 col-xl-3">
                                <div className="card">
                                    <div className="card-header">BMI</div>
                                    <div className="card-body">
                                        <div className="card-text">
                                            <p>
                                                <span className="big">22.2</span>
                                            </p>
                                        </div>
                                    </div>
                                    <div className="card-footer">5 feb 2018</div>
                                </div>
                            </div>

                            <div className="col-md-5 col-lg-4 col-xl-3">
                                <div className="card">
                                    <div className="card-header">BMI</div>
                                    <div className="card-body">
                                        <div className="card-text">
                                            <p>
                                                <span className="big">22.2</span>
                                            </p>
                                        </div>
                                    </div>
                                    <div className="card-footer">5 feb 2018</div>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {this.props.canSeeLaboratoryTests && (
                    <div className="section">
                        <header>
                            <h2>
                                <LaboratoryIcon />
                                Laboratory Tests
                            </h2>
                            {this.props.canAddLaboratoryTests && (
                                <Link to={joinPaths(match.url, "add-lab-test")} className="btn btn-link">
                                    Add laboratory test
                                </Link>
                            )}
                        </header>

                        <dl className="lab">
                            <dt>Bladder or kidney infections</dt>
                            <dd>
                                <PositiveIcon />
                                white or red blood cells or bacteria in the urine
                            </dd>
                            <dt>Pregnancy</dt>
                            <dd>
                                <PositiveIcon />
                                hCG in urine after 2 weeks post-conception
                            </dd>
                            <dt>Preeclampsia</dt>
                            <dd>
                                <NegativeIcon />
                                high blood presure plus protein in the urine
                            </dd>
                        </dl>
                    </div>
                )}

                {this.props.canAddDiagnosis && (<Route path={match.path + "/add-diagnosis"} component={AddDiagnosis} />)}
                {this.props.canAddVitalSigns && (<Route path={match.path + "/add-data"} component={MedicalData} />)}
                {this.props.canAddLaboratoryTests && (<Route path={match.path + "/add-lab-test"} component={LaboratoryTest} />)}
                {this.props.canEditMainComplaint && (<Route path={match.path + "/edit-complaint"} component={EditComplaint} />)}
                {this.props.canRemoveFromWaitlist && (<Route path={match.path + "/remove"} component={Remove} />)}
            </div>
        )
    }
}

InConsultation = connect(
    state => ({
        patient: state.patient.patient,
        waitlistItem: state.waitlist.item,
        waitlistFetching: state.waitlist.fetching,
        canSeeDiagnosis: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ],
        canAddDiagnosis: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[WRITE],
        canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
        canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE],
        canSeeLaboratoryTests: ((state.validations.userRights || {})[RESOURCE_LABORATORY_TEST] || {})[READ],
        canAddLaboratoryTests: ((state.validations.userRights || {})[RESOURCE_LABORATORY_TEST] || {})[WRITE],
        canSeeMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canEditMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[UPDATE],
        canRemoveFromWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[DELETE],
    }),
    {}
)(InConsultation)

export default InConsultation
