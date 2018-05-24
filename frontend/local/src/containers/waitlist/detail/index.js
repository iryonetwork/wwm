import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import moment from "moment"

import "./style.css"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as DiagnosisIcon } from "shared/icons/diagnosis.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
//import { ReactComponent as LaboratoryIcon } from "shared/icons/laboratory.svg"
//import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"
//import { ReactComponent as PositiveIcon } from "shared/icons/positive.svg"

import { joinPaths } from "shared/utils"
import { fetchCode } from "shared/modules/codes"
import {
    RESOURCE_PATIENT_IDENTIFICATION,
    RESOURCE_VITAL_SIGNS,
    RESOURCE_EXAMINATION,
    RESOURCE_LABORATORY_TEST,
    RESOURCE_WAITLIST,
    READ,
    WRITE,
    UPDATE,
    DELETE
} from "../../../modules/validations"
import { saveConsultation } from "../../../modules/patient"

import MedicalData from "./add-data"
import LaboratoryTest from "./add-lab-test"
import EditComplaint from "./edit-complaint"
import AddDiagnosis from "./add-diagnosis"
import Remove from "./remove"
import Spinner from "shared/containers/spinner"

class CodeTitle extends React.Component {
    constructor(props) {
        super(props)
        this.loadCode = this.loadCode.bind(this)
        if (props.codeId) {
            this.loadCode(props.categoryId, props.codeId)
        }

        this.componentWillUnmount = this.componentWillUnmount.bind(this)
        this.state = { loading: true, title: "" }
    }

    loadCode(categoryId, codeId) {
        this.props
            .fetchCode(categoryId, codeId)
            .then(code => {
                if (!this.unmounted) {
                    this.setState({ loading: false, failed: false, title: code.title })
                }
            })
            .catch(ex => {
                this.setState({ loading: false, failed: true })
            })
    }

    componentWillUnmount() {
        this.unmounted = true
    }

    render() {
        if (this.state.loading) {
            return <span>...</span>
        }

        if (this.state.failed) {
            return <span>Failed to fetch title!</span>
        }

        return <span>{this.state.title}</span>
    }
}

CodeTitle = connect(state => ({}), { fetchCode })(CodeTitle)

class InConsultation extends React.Component {
    constructor(props) {
        super(props)
        this.closeConsultation = this.closeConsultation.bind(this)
        this.state = {saving: false}
    }

    closeConsultation(ev) {
        ev.preventDefault()
        this.setState({saving: true})
        this.props.saveConsultation(this.props.match.params.waitlistID, this.props.match.params.itemID)
    }

    render() {
        const { match, waitlistItem, waitlistFetching } = this.props

        if (this.state.saving || waitlistFetching) {
            return <Spinner />
        }
        return waitlistItem ? (
            <div>
                <header>
                    <h1>In consultation</h1>
                    {this.props.canAddDiagnosis && (
                        <React.Fragment>
                            {(!waitlistItem.diagnoses && (
                                <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-success btn-wide">
                                    Add diagnosis
                                </Link>
                            )) || (
                                <a href={`/waitlist/${match.params.waitlistID}`} onClick={this.closeConsultation} className="btn btn-success btn-wide btn-close">
                                    Close
                                </a>
                            )}
                        </React.Fragment>
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

                        <h3>{(waitlistItem.mainComplaint || {}).complaint}</h3>
                        {(waitlistItem.mainComplaint || {}).comment && <p>{(waitlistItem.mainComplaint || {}).comment}</p>}
                    </div>
                )}

                {this.props.canSeeDiagnosis &&
                    (waitlistItem.diagnoses || []).length > 0 && (
                        <div className="section">
                            <header>
                                <h2>
                                    <DiagnosisIcon />
                                    Diagnoses
                                </h2>
                                {this.props.canAddDiagnosis && (
                                    <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-link">
                                        Add complementary diagnosis
                                    </Link>
                                )}
                            </header>

                            {waitlistItem.diagnoses.map((el, key) => (
                                <React.Fragment key={key}>
                                    <h3>
                                        {el.label ? el.label : (<CodeTitle categoryId="diagnosis" codeId={el.diagnosis} />)}
                                    </h3>
                                    {el.comment && <p>{el.comment}</p>}

                                    {el.therapies && <h4>Therapies</h4>}
                                    {(el.therapies || []).map((tel, tkey) => (
                                        <React.Fragment key={tkey}>
                                            <h5>{tel.medicine}</h5>
                                            {tel.instructions && <p>{tel.instructions}</p>}
                                        </React.Fragment>
                                    ))}
                                    {this.props.canAddDiagnosis && (
                                        <Link to={joinPaths(match.url, `diagnoses/${key}/edit`)} className="btn btn-link">
                                            Edit diagnosis
                                        </Link>
                                    )}
                                </React.Fragment>
                            ))}
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
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.height && (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Height</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.height.value}</span>cm
                                                    </p>
                                                    {/* <p>5ft 1in</p> */}
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.height.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}

                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.weight && (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Body mass</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.weight.value}</span>kg
                                                    </p>
                                                    {/* <p>1008.8 lb</p> */}
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.weight.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}

                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.height &&
                                waitlistItem.vitalSigns.weight && (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">BMI</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">
                                                            {round(
                                                                waitlistItem.vitalSigns.weight.value /
                                                                    waitlistItem.vitalSigns.height.value /
                                                                    waitlistItem.vitalSigns.height.value *
                                                                    10000,
                                                                2
                                                            )}
                                                        </span>
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.height.timestamp, "X").isAfter(moment(waitlistItem.vitalSigns.weight.timestamp, "X")) ? moment(waitlistItem.vitalSigns.height.timestamp, "X").format("Do MMM Y") : moment(waitlistItem.vitalSigns.weight.timestamp, "X").format("Do MMM Y") }</div>
                                        </div>
                                    </div>
                                )}

                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.temperature &&
                                (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Body temperature</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.temperature.value}</span>°C
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.temperature.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.heart_rate &&
                                (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Heart rate</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.heart_rate.value}</span>bpm
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.heart_rate.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.pressure &&
                                waitlistItem.vitalSigns.pressure.value.systolic &&
                                waitlistItem.vitalSigns.pressure.value.diastolic &&
                                (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Blood pressure</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.pressure.value.systolic}/{waitlistItem.vitalSigns.pressure.value.diastolic}</span>mmHg
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.pressure.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}
                            {waitlistItem.vitalSigns &&
                                waitlistItem.vitalSigns.oxygen_saturation &&
                                (
                                    <div className="col-md-5 col-lg-4 col-xl-3">
                                        <div className="card">
                                            <div className="card-header">Oxygen saturation</div>
                                            <div className="card-body">
                                                <div className="card-text">
                                                    <p>
                                                        <span className="big">{waitlistItem.vitalSigns.oxygen_saturation.value}</span>%
                                                    </p>
                                                </div>
                                            </div>
                                            <div className="card-footer">{moment(waitlistItem.vitalSigns.oxygen_saturation.timestamp, "X").format("Do MMM Y")}</div>
                                        </div>
                                    </div>
                                )}
                        </div>
                    </div>
                )}

                {/*this.props.canSeeLaboratoryTests && (
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
                )*/}

                {this.props.canAddDiagnosis && <Route path={match.path + "/add-diagnosis"} component={AddDiagnosis} />}
                {this.props.canAddDiagnosis && <Route path={match.path + "/diagnoses/:diagnosisIndex/edit"} component={AddDiagnosis} />}
                {this.props.canAddVitalSigns && <Route path={match.path + "/add-data"} component={MedicalData} />}
                {this.props.canAddLaboratoryTests && <Route path={match.path + "/add-lab-test"} component={LaboratoryTest} />}
                {this.props.canEditMainComplaint && <Route path={match.path + "/edit-complaint"} component={EditComplaint} />}
                {this.props.canRemoveFromWaitlist && <Route path={match.path + "/remove"} component={Remove} />}
            </div>
        ) : (null)
    }
}

InConsultation = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]

        return {
            patient: state.patient.patient,
            waitlistItem: item,
            waitlistFetching: state.waitlist.fetching || state.waitlist.listing,
            canSeeDiagnosis: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ],
            canAddDiagnosis: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[WRITE],
            canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
            canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE],
            canSeeLaboratoryTests: ((state.validations.userRights || {})[RESOURCE_LABORATORY_TEST] || {})[READ],
            canAddLaboratoryTests: ((state.validations.userRights || {})[RESOURCE_LABORATORY_TEST] || {})[WRITE],
            canSeeMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
            canEditMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[UPDATE],
            canRemoveFromWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[DELETE]
        }
    },
    {
        saveConsultation
    }
)(InConsultation)

export default InConsultation

const round = (number, precision) => {
    var shift = function(number, precision) {
        var numArray = ("" + number).split("e")
        return +(numArray[0] + "e" + (numArray[1] ? +numArray[1] + precision : precision))
    }
    return shift(Math.round(shift(number, +precision)), -precision)
}
