import React from "react"

//import "./style.css"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as LaboratoryIcon } from "shared/icons/laboratory.svg"

import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"
import { ReactComponent as PositiveIcon } from "shared/icons/positive.svg"

export default () => (
    <div>
        <header>
            <h1>In consultation</h1>
            <button className="btn btn-primary btn-wide">Add diagnosis</button>
        </header>

        <div className="section">
            <header>
                <h2>
                    <ComplaintIcon />Main Complaint
                </h2>
                <button className="btn btn-link">Edit main complaint</button>
            </header>

            <h3>Knee pain</h3>
            <p>
                Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since
                the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries,
                but also the leap into electronic typesetting, remaining essentially unchanged.
            </p>
        </div>

        <div className="section">
            <header>
                <h2>
                    <MedicalDataIcon />
                    Medical Data
                </h2>
                <button className="btn btn-link">Add medical data</button>
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
                <button className="btn btn-link">Add laboratory test</button>
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
    </div>
)
