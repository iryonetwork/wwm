import React from "react"
import { connect } from "react-redux"
import { push } from "react-router-redux"
import { Route, Link } from "react-router-dom"
import _ from "lodash"

import "./style.css"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as DiagnosisIcon } from "shared/icons/diagnosis.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as TherapyIcon } from "shared/icons/therapy.svg"

import { joinPaths } from "shared/utils"
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
import VitalSignCard from "shared/containers/vitalSign"
import CodeTitle from "shared/containers/codes/title"
import { open, COLOR_SUCCESS } from "shared/modules/alert"

class InConsultation extends React.Component {
    constructor(props) {
        super(props)
        this.closeConsultation = this.closeConsultation.bind(this)
    }

    closeConsultation(ev) {
        ev.preventDefault()
        this.props.saveConsultation(this.props.match.params.waitlistID, this.props.match.params.itemID).then(() => {
            this.props.push(`/waitlist/${this.props.match.params.waitlistID}`)
            this.props.open("Consultation was closed", "", COLOR_SUCCESS, 5)
        })
    }

    render() {
        const { match, saving, waitlistItem, waitlistFetching } = this.props

        // collect all therapies
        let therapies = []
        waitlistItem &&
            _.forEach(waitlistItem.diagnoses, (diagnosis, i) => {
                _.forEach(diagnosis.therapies, therapy => {
                    let t = _.clone(therapy)
                    t.diagnosis = i
                    therapies.push(t)
                })
            })

        if (saving || waitlistFetching) {
            return <Spinner />
        }
        return waitlistItem ? (
            <div className="consultation">
                <header>
                    <h1>In Consultation</h1>
                    {this.props.canAddDiagnosis && (
                        <React.Fragment>
                            {(!waitlistItem.diagnoses && (
                                <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-success btn-wide">
                                    Add diagnosis
                                </Link>
                            )) || (
                                <a
                                    href={`/waitlist/${match.params.waitlistID}`}
                                    onClick={this.closeConsultation}
                                    className="btn btn-success btn-wide btn-close"
                                >
                                    Close
                                </a>
                            )}
                        </React.Fragment>
                    )}
                </header>

                {this.props.canSeeDiagnosis &&
                    (waitlistItem.diagnoses || []).length > 0 && (
                        <React.Fragment key="diagnoses">
                            <div className="section" key="diangosis0">
                                <header>
                                    <h2>
                                        <DiagnosisIcon />Diagnosis
                                    </h2>
                                    {this.props.canAddDiagnosis && (
                                        <Link to={joinPaths(match.url, "add-diagnosis")} className="btn btn-link">
                                            Add Complementary Diagnosis
                                        </Link>
                                    )}
                                </header>
                                <div className="part diagnosis" key="diagnosis0body">
                                    <header>
                                        <h3>
                                            {waitlistItem.diagnoses[0].label ? (
                                                waitlistItem.diagnoses[0].label
                                            ) : (
                                                <CodeTitle categoryId="diagnosis" codeId={waitlistItem.diagnoses[0].diagnosis} />
                                            )}
                                        </h3>
                                        {this.props.canAddDiagnosis && (
                                            <Link to={joinPaths(match.url, `diagnoses/0/edit`)} className="btn btn-link">
                                                Edit
                                            </Link>
                                        )}
                                    </header>
                                    <React.Fragment key="diagnosis0comment">
                                        {waitlistItem.diagnoses[0].comment && <p>{waitlistItem.diagnoses[0].comment}</p>}
                                    </React.Fragment>
                                </div>
                                {waitlistItem.diagnoses.length > 1 && (
                                    <div className="subsection part" key="complementaryDiagnoses">
                                        <header>
                                            <h3>Complementary diagnoses</h3>
                                        </header>
                                        <dl>
                                            {waitlistItem.diagnoses.map((el, key) => {
                                                return (
                                                    key !== 0 && (
                                                        <React.Fragment key={`diagnosis${key}`}>
                                                            <dt className="with-btn">
                                                                <h4>{el.label ? el.label : <CodeTitle categoryId="diagnosis" codeId={el.diagnosis} />}</h4>
                                                                {this.props.canAddDiagnosis && (
                                                                    <Link to={joinPaths(match.url, `diagnoses/${key}/edit`)} className="btn btn-link">
                                                                        Edit
                                                                    </Link>
                                                                )}
                                                            </dt>
                                                            <dd>{el.comment}</dd>
                                                        </React.Fragment>
                                                    )
                                                )
                                            })}
                                        </dl>
                                    </div>
                                )}
                            </div>

                            {therapies.length > 0 && (
                                <div className="section" key="therapy">
                                    <header>
                                        <h2>
                                            <TherapyIcon />Therapy
                                        </h2>
                                    </header>
                                    <div className="part">
                                        <dl>
                                            {therapies.map((therapy, i) => (
                                                <React.Fragment key={i}>
                                                    <dt>
                                                        <h3>{therapy.medicine}</h3>

                                                        {waitlistItem.diagnoses.length > 1 &&
                                                            waitlistItem.diagnoses[therapy.diagnosis].label && (
                                                                <aside className="diagnosisReference">{waitlistItem.diagnoses[therapy.diagnosis].label}</aside>
                                                            )}
                                                    </dt>
                                                    <dd>{therapy.instructions}</dd>
                                                </React.Fragment>
                                            ))}
                                        </dl>
                                    </div>
                                </div>
                            )}
                        </React.Fragment>
                    )}

                {this.props.canSeeMainComplaint && (
                    <div className="section">
                        <header>
                            <h2>
                                <ComplaintIcon />Main complaint
                            </h2>
                            {this.props.canEditMainComplaint && (
                                <Link to={joinPaths(match.url, "edit-complaint")} className="btn btn-link">
                                    Edit Main Complaint
                                </Link>
                            )}
                        </header>
                        <div className="part" key="mainComplaint">
                            {!_.isEmpty(waitlistItem.mainComplaint) ? (
                                <dl>
                                    <dt>{waitlistItem.mainComplaint ? waitlistItem.mainComplaint.complaint : null}</dt>
                                    {waitlistItem.mainComplaint && waitlistItem.mainComplaint.comment && <dd>{waitlistItem.mainComplaint.comment}</dd>}
                                </dl>
                            ) : (
                                <dl>
                                    <dd className="missing">Main complaint was not set</dd>
                                </dl>
                            )}
                        </div>
                    </div>
                )}

                {this.props.canSeeVitalSigns && (
                    <div className="section">
                        <header>
                            <h2>
                                <MedicalDataIcon />
                                Medical data
                            </h2>
                            {this.props.canAddVitalSigns && (
                                <Link to={joinPaths(match.url, "add-data")} className="btn btn-link">
                                    Add Medical Data
                                </Link>
                            )}
                        </header>
                        <div className="part" key="vitalSigns">
                            <div className="card-group">
                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.height && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="height"
                                                name="Height"
                                                value={waitlistItem.vitalSigns.height.value}
                                                unit="cm"
                                                timestamp={waitlistItem.vitalSigns.height.timestamp}
                                            />
                                        </div>
                                    )}

                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.weight && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="weight"
                                                name="Body mass"
                                                value={waitlistItem.vitalSigns.weight.value}
                                                unit="kg"
                                                timestamp={waitlistItem.vitalSigns.weight.timestamp}
                                            />
                                        </div>
                                    )}

                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.bmi && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="bmi"
                                                name="BMI"
                                                value={waitlistItem.vitalSigns.bmi.value}
                                                unit=""
                                                timestamp={waitlistItem.vitalSigns.bmi.timestamp}
                                            />
                                        </div>
                                    )}

                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.temperature && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="temperature"
                                                name="Body temperature"
                                                value={waitlistItem.vitalSigns.temperature.value}
                                                unit="°C"
                                                timestamp={waitlistItem.vitalSigns.temperature.timestamp}
                                            />
                                        </div>
                                    )}
                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.heart_rate && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="heart_rate"
                                                name="Heart rate"
                                                value={waitlistItem.vitalSigns.heart_rate.value}
                                                unit="bpm"
                                                timestamp={waitlistItem.vitalSigns.heart_rate.timestamp}
                                            />
                                        </div>
                                    )}
                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.pressure &&
                                    waitlistItem.vitalSigns.pressure.value.systolic &&
                                    waitlistItem.vitalSigns.pressure.value.diastolic && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="pressure"
                                                name="Blood pressure"
                                                value={`${waitlistItem.vitalSigns.pressure.value.systolic}/${waitlistItem.vitalSigns.pressure.value.diastolic}`}
                                                unit="mmHg"
                                                timestamp={waitlistItem.vitalSigns.pressure.timestamp}
                                            />
                                        </div>
                                    )}
                                {waitlistItem.vitalSigns &&
                                    waitlistItem.vitalSigns.oxygen_saturation && (
                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                            <VitalSignCard
                                                id="oxygen_saturation"
                                                name="Oxygen saturation"
                                                value={waitlistItem.vitalSigns.oxygen_saturation.value}
                                                unit="%"
                                                timestamp={waitlistItem.vitalSigns.oxygen_saturation.timestamp}
                                            />
                                        </div>
                                    )}
                            </div>
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
                {this.props.canAddVitalSigns && <Route exact path={match.path + "/add-data/:sign"} component={MedicalData} />}
                {this.props.canAddVitalSigns && <Route exact path={match.path + "/add-data"} component={MedicalData} />}
                {this.props.canAddLaboratoryTests && <Route path={match.path + "/add-lab-test"} component={LaboratoryTest} />}
                {this.props.canEditMainComplaint && <Route path={match.path + "/edit-complaint"} component={EditComplaint} />}
                {this.props.canRemoveFromWaitlist && <Route path={match.path + "/remove"} component={Remove} />}
            </div>
        ) : null
    }
}

InConsultation = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]

        return {
            patient: state.patient.patient,
            saving: state.patient.saving,
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
        saveConsultation,
        open,
        push
    }
)(InConsultation)

export default InConsultation
