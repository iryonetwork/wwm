import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"

import "./style.css"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as LaboratoryIcon } from "shared/icons/laboratory.svg"
import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"
import { ReactComponent as PositiveIcon } from "shared/icons/positive.svg"

import { joinPaths } from "shared/utils"

import MedicalData from "./add-data"
import LaboratoryTest from "./add-lab-test"
import EditComplaint from "./edit-complaint"
import AddDiagnosis from "./add-diagnosis"

let InConsultation = ({ match, waitlistItem, waitlistFetching, patient }) => {
    if (waitlistFetching) {
        return null
    }

    return (
        <div>
            <header>
                <h1>In consultation</h1>
                <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-primary btn-wide">
                    Add diagnosis
                </Link>
            </header>

            <div className="section">
                <header>
                    <h2>
                        <ComplaintIcon />Main Complaint
                    </h2>
                    <Link to={joinPaths(match.url, "edit-complaint")} className="btn btn-link">
                        Edit main complaint
                    </Link>
                </header>

                <h3>{waitlistItem.complaint}</h3>
                {waitlistItem.comment && <p>{waitlistItem.comment}</p>}
            </div>

            <div className="section">
                <header>
                    <h2>
                        <MedicalDataIcon />
                        Medical Data
                    </h2>
                    <Link to={joinPaths(match.url, "add-data")} className="btn btn-link">
                        Add medical data
                    </Link>
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

            <div className="section">
                <header>
                    <h2>
                        <LaboratoryIcon />
                        Laboratory Tests
                    </h2>
                    <Link to={joinPaths(match.url, "add-lab-test")} className="btn btn-link">
                        Add laboratory test
                    </Link>
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

            <Route path={match.url + "/add-diagnosis"} component={AddDiagnosis} />
            <Route path={match.url + "/add-data"} component={MedicalData} />
            <Route path={match.url + "/add-lab-test"} component={LaboratoryTest} />
            <Route path={match.url + "/edit-complaint"} component={EditComplaint} />
        </div>
    )
}

InConsultation = connect(
    store => ({
        patient: store.patient.patient,
        waitlistItem: store.waitlist.item,
        waitlistFetching: store.waitlist.fetching
    }),
    {}
)(InConsultation)

export default InConsultation
