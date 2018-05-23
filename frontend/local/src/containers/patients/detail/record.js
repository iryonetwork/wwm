import React from "react"
import _ from "lodash"
import { connect } from "react-redux"

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
    }

    componentDidUpdate(prevProps) {
        if (prevProps.patientID !== this.props.patientID) {
            this.props.fetchHealthRecords(this.props.patientID)
        }
    }

    render() {
        let { records, canSeeExamination } = this.props

        if (records.loading) {
            return <Spinner />
        }

        return (
            <div className="records">
                <header>
                    <h1>Health Record</h1>
                </header>

                {canSeeExamination && (
                    <div className="section">
                        {records.data && records.data.length ? (
                            <div>
                                {records.data.map(({ data, meta }) => (
                                    <div key={meta.name}>
                                        <h3>
                                            <ComplaintIcon /> {data.mainComplaint.complaint}
                                        </h3>
                                        <div className="small">{new Date(meta.created).toLocaleString()}</div>
                                        {data.mainComplaint.comment && <div className="comment">{data.mainComplaint.comment}</div>}

                                        {!_.isEmpty(data.diagnoses) && (
                                            <div className="part">
                                                <h4>
                                                    <DiagnosisIcon /> Diagnoses
                                                </h4>
                                                <dl>
                                                    {data.diagnoses.map(diagnosis => (
                                                        <React.Fragment key={diagnosis.comment}>
                                                            <dt>{diagnosis.name || "NAME"}</dt>
                                                            <dd>{diagnosis.comment}</dd>
                                                        </React.Fragment>
                                                    ))}
                                                </dl>
                                            </div>
                                        )}

                                        {!_.isEmpty(data.therapies) && (
                                            <div className="part">
                                                <h4>
                                                    <TherapyIcon /> Therapies
                                                </h4>
                                                <dl>
                                                    {data.therapies.map(therapy => (
                                                        <React.Fragment key={therapy.medication}>
                                                            <dt>{therapy.medication}</dt>
                                                            <dd>{therapy.instructions}</dd>
                                                        </React.Fragment>
                                                    ))}
                                                </dl>
                                            </div>
                                        )}

                                        {!_.isEmpty(data.vitalSigns) && (
                                            <div className="part">
                                                <h4>
                                                    <MedicalDataIcon /> Vital Signs
                                                </h4>
                                                <dl>
                                                    {_.map(data.vitalSigns, (value, key) => (
                                                        <React.Fragment key={key}>
                                                            <dt>{_.upperFirst(key)}</dt>
                                                            <dd>{value}</dd>
                                                        </React.Fragment>
                                                    ))}
                                                </dl>
                                            </div>
                                        )}
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
    state => ({
        patientID: state.patient.patient.ID,
        records: state.patient.patientRecords,
        canSeeExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ]
    }),
    { fetchHealthRecords }
)(HealthRecord)

export default HealthRecord
