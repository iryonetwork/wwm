import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import  moment  from "moment"
import { Collapse } from 'reactstrap';

import Spinner from "shared/containers/spinner"
import { RESOURCE_EXAMINATION, READ } from "../../../modules/validations"
import { fetchHealthRecords } from "../../../modules/patient"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import { ReactComponent as DiagnosisIcon } from "shared/icons/diagnosis.svg"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as TherapyIcon } from "shared/icons/therapy.svg"

class HealthRecord extends React.Component {
    constructor(props) {
        super(props)
        if (props.patientID) {
            props.fetchHealthRecords(props.patientID)
        }

        this.togglePart = this.togglePart.bind(this)

        this.state = {}
    }

    componentWillReceiveProps(nextProps) {
        if ((nextProps.records === undefined || nextProps.patientID !== nextProps.loadedPatientID) && !nextProps.recordsLoading) {
            this.props.fetchHealthRecords(nextProps.patientID)
        }
    }


    togglePart = (fileName, part) => () => {
        let fileState = _.clone(this.state[fileName]) || {}
        fileState[part] = fileState[part] ? !fileState[part] : true
        this.setState({[fileName]: fileState})
    }

    render() {
        let { records, recordsLoading, canSeeExamination } = this.props

        if (recordsLoading) {
            return <Spinner />
        }

        return (
            <div className="records">
                <header>
                    <h1>Health Record</h1>
                </header>

                {canSeeExamination && (
                    <div>
                        {records && records.length ? (
                            <div>
                                {records.map(({ data, meta }) => (
                                    <div key={meta.name}>
                                        <div className="dateheader">{moment(meta.created).format("dddd, Do MMM YYYY")}</div>
                                        <div className="section" key={meta.name}>
                                            <div key="mainDiagnosis">
                                                {!_.isEmpty(data.diagnoses) ? (
                                                    <React.Fragment key="0">
                                                        <h3>
                                                            <DiagnosisIcon /> {data.diagnoses[0].diagnosis ? (data.diagnoses[0].diagnosis.label && !_.isObject(data.diagnoses[0].diagnosis.label) ? data.diagnoses[0].diagnosis.label : "Diagnosis") : "Diagnosis"}
                                                        </h3>
                                                        <div className="comment">{data.diagnoses[0].comment ? data.diagnoses[0].comment : ""}</div>
                                                    </React.Fragment>
                                                ) : (
                                                    <React.Fragment key="missingDiagnosis">
                                                        <h3 className="missing">
                                                            <DiagnosisIcon /> Diagnosis was not set
                                                        </h3>
                                                    </React.Fragment>
                                                )}
                                            </div>


                                            {!_.isEmpty(data.mainComplaint) ? (
                                                <div className="part" key="mainComplaint">
                                                    <div className="partHeader" onClick={this.togglePart(meta.name, "mainComplaint")}>
                                                        <h4>
                                                            <ComplaintIcon /> Main complaint
                                                        </h4>
                                                    </div>
                                                    <Collapse isOpen={this.state[meta.name] ? this.state[meta.name]["mainComplaint"] : false}>
                                                        <dl>
                                                            <dt>{data.mainComplaint ? data.mainComplaint.complaint : null}</dt>
                                                            {data.mainComplaint && data.mainComplaint.comment && <dd>{data.mainComplaint.comment}</dd>}
                                                        </dl>
                                                    </Collapse>
                                                </div>
                                            ) : (
                                                <div className="part">
                                                    <h4 className="missing">
                                                        <ComplaintIcon /> Main complaint was not set
                                                    </h4>
                                                </div>
                                            )}

                                            {!_.isEmpty(_.filter(data.diagnoses, (diagnosis, i) => (i !== 0))) && (
                                                <div className="part" key="complementaryDiagnoses">
                                                    <div className="partHeader" onClick={this.togglePart(meta.name, "complementaryDiagnoses")}>
                                                        <h4>
                                                            <DiagnosisIcon /> Complementary diagnoses
                                                        </h4>
                                                    </div>
                                                    <Collapse isOpen={this.state[meta.name] ? this.state[meta.name]["complementaryDiagnoses"] : false}>
                                                        <dl>
                                                            {data.diagnoses.map((diagnosis, i) => {
                                                                return (i !== 0) && (
                                                                    <React.Fragment key={"diagnosis" + i}>
                                                                        <dt>{diagnosis.diagnosis ? (diagnosis.diagnosis.label || "Diagnosis") : "Diagnosis"}</dt>
                                                                        <dd>{diagnosis.comment}</dd>
                                                                    </React.Fragment>
                                                            )})}
                                                        </dl>
                                                    </Collapse>
                                                </div>
                                            )}

                                            {!_.isEmpty(data.therapies) && (
                                                <div className="part" key="therapies">
                                                    <div className="partHeader" onClick={this.togglePart(meta.name, "therapies")}>
                                                        <h4>
                                                            <TherapyIcon /> Therapy
                                                        </h4>
                                                    </div>

                                                    <Collapse isOpen={this.state[meta.name] ? this.state[meta.name]["therapies"] : false}>
                                                        <dl>
                                                            {data.therapies.map((therapy, i) => (
                                                                <React.Fragment key={i}>
                                                                    <dt>{therapy.medication}</dt>
                                                                    {data.diagnoses.length > 1 && data.diagnoses[therapy.diagnosis].diagnosis.label && <aside className="diagnosisReference">{data.diagnoses[therapy.diagnosis].diagnosis.label}</aside>}
                                                                    <dd>{therapy.instructions}</dd>
                                                                </React.Fragment>
                                                            ))}
                                                        </dl>
                                                    </Collapse>
                                                </div>
                                            )}


                                            {!_.isEmpty(data.vitalSigns) && (
                                                <div className="part">
                                                    <div className="partHeader" onClick={this.togglePart(meta.name, "medicalData")}>
                                                        <h4>
                                                            <MedicalDataIcon /> Medical data
                                                        </h4>
                                                    </div>
                                                    <Collapse isOpen={this.state[meta.name] ? this.state[meta.name]["medicalData"] : false}>
                                                        <dl>
                                                            <div className="card-group">
                                                                {data.vitalSigns.height && (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Height</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.height.value}</span>cm
                                                                                        </p>
                                                                                        {/* <p>5ft 1in</p> */}
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.height.timestamp ? moment(data.vitalSigns.height.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                {data.vitalSigns.weight && (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Body mass</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.weight.value}</span>kg
                                                                                        </p>
                                                                                        {/* <p>1008.8 lb</p> */}
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.weight.timestamp ? moment(data.vitalSigns.weight.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                {data.vitalSigns.bmi && (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">BMI</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.bmi.value}</span>
                                                                                        </p>
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.bmi.timestamp ? moment(data.vitalSigns.bmi.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                {data.vitalSigns.temperature &&
                                                                    (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Body temperature</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.temperature.value}</span>°C
                                                                                        </p>
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.temperature.timestamp ? moment(data.vitalSigns.temperature.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                {data.vitalSigns.heart_rate &&
                                                                    (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Heart rate</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.heart_rate.value}</span>bpm
                                                                                        </p>
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.heart_rate.timestamp ? moment(data.vitalSigns.heart_rate.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}
                                                                {data.vitalSigns.pressure && data.vitalSigns.pressure.value &&
                                                                    data.vitalSigns.pressure.value.systolic &&
                                                                    data.vitalSigns.pressure.value.diastolic &&
                                                                    (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Blood pressure</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.pressure.value.systolic}/{data.vitalSigns.pressure.value.diastolic}</span>mmHg
                                                                                        </p>
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.pressure.timestamp ? moment(data.vitalSigns.pressure.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}

                                                                {data.vitalSigns.oxygen_saturation &&
                                                                    (
                                                                        <div className="col-md-5 col-lg-4 col-xl-3">
                                                                            <div className="card">
                                                                                <div className="card-header">Oxygen saturation</div>
                                                                                <div className="card-body">
                                                                                    <div className="card-text">
                                                                                        <p>
                                                                                            <span className="big">{data.vitalSigns.oxygen_saturation.value}</span>%
                                                                                        </p>
                                                                                    </div>
                                                                                </div>
                                                                                <div className="card-footer">{data.vitalSigns.oxygen_saturation.timestamp ? moment(data.vitalSigns.oxygen_saturation.timestamp).format("Do MMM Y") : moment(meta.created).format("Do MMM Y")}</div>
                                                                            </div>
                                                                        </div>
                                                                    )}
                                                                </div>
                                                        </dl>
                                                    </Collapse>
                                                </div>
                                            )}

                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <h3>No records found</h3>
                        )}
                    </div>
                )}
            </div>
        )
    }
}

HealthRecord = connect(
    (state, props) => {
        let records = state.patient.patientRecords.data
        // sort records by creation time and reverse to have latest record as first
        records = _.reverse(_.sortBy(records, [function(obj) { return obj.meta.created; }]))

        return {
            patientID: props.match.params.patientID || state.patient.patient.ID,
            loadedPatientID: state.patient.patient.ID,
            records: records,
            recordsLoading: state.patient.patientRecords.loading,
            canSeeExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ]
        }
    },
    { fetchHealthRecords }
)(HealthRecord)

export default HealthRecord
