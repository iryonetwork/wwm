import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import { reduxForm } from "redux-form"

import { RenderForm } from "../new/step3"
import { joinPaths } from "shared/utils"
import { RESOURCE_HEALTH_HISTORY, READ, WRITE } from "../../../modules/validations"
import { getCodesAsOptions, loadCategories as loadCategoriesImport } from "shared/modules/codes"
import { updatePatient } from "../../../modules/patient"

let History = ({ match, patient, canSee, canEdit }) => {
    return canSee ? (
        <div className="history">
            <header>
                <h1>Medical History</h1>
                {canEdit ? (
                    <Link to={joinPaths(match.url, "edit")} className="btn btn-secondary btn-wide">
                        Edit
                    </Link>
                ) : null}
            </header>

            {/* @TODO confirm if needed <div className="section">
                <div className="name">Blood type</div>
                <div className="values">A+</div>
            </div> */}

            <div className="section">
                <div className="name">Allergies</div>
                <div className="values">
                    <dl>
                        {(patient.allergies || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.allergy}</dt>
                                {(item.critical === "true" || item.comment) && (
                                    <dd>
                                        {item.critical === "true" && <div>High risk</div>}
                                        {item.comment && <div>{item.comment}</div>}
                                    </dd>
                                )}
                            </React.Fragment>
                        ))}
                        {(patient.allergies || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Immunization</div>
                <div className="values">
                    <dl>
                        {(patient.immunizations || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.immunization}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{item.date}</dd>}
                            </React.Fragment>
                        ))}
                        {(patient.immunizations || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Chronic diseases</div>
                <div className="values">
                    <dl>
                        {(patient.chronicDiseases || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.disease}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{item.date}</dd>}
                            </React.Fragment>
                        ))}
                        {(patient.chronicDiseases || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Injuries &amp; handicaps</div>
                <div className="values">
                    <dl>
                        {(patient.injuries || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.injury}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{item.date}</dd>}
                            </React.Fragment>
                        ))}
                        {(patient.injuries || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Surgeries</div>
                <div className="values">
                    <dl>
                        {(patient.surgeries || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.injury}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{item.date}</dd>}
                            </React.Fragment>
                        ))}
                        {(patient.surgeries || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Additional medications</div>
                <div className="values">
                    <dl>
                        {(patient.medications || []).map((item, i) => (
                            <React.Fragment key={i}>
                                <dt>{item.medication}</dt>
                                {item.comment && <dd>{item.comment}</dd>}
                            </React.Fragment>
                        ))}
                        {(patient.medications || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>
        </div>
    ) : null
}

History = connect(
    state => ({
        patient: state.patient.patient,
        canSee: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[READ],
        canEdit: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[WRITE]
    }),
    {}
)(History)

class EditHistory extends React.Component {
    constructor(props) {
        super(props)

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    handleSubmit(form) {
        this.props.updatePatient(form).then(() => {
            this.props.history.push(".")
        })
    }

    componentWillMount() {
        this.props.loadCategories("babyFood", "childCommunication", "deliveryType")
    }

    render() {
        const { handleSubmit, dateOfBirth, codesLoading, getCodes, updating } = this.props
        return this.props.canEdit ? (
            <div className="edit-history">
                <header>
                    <h1>Edit Medical History</h1>
                    <Link to="." className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={handleSubmit(this.handleSubmit)} className="btn btn-primary btn-wide">
                        {updating ? "Saving..." : "Save"}
                    </button>
                </header>
                <div className="formContainer">
                    <form onSubmit={handleSubmit(this.handleSubmit)} className="patient-form">
                        <RenderForm
                            dateOfBirth={dateOfBirth}
                            babyFoods={getCodes("babyFood")}
                            communicationTypes={getCodes("childCommunication")}
                            deliveryTypes={getCodes("deliveryType")}
                            codesLoading={codesLoading && !updating}
                        />

                        <div className="row buttons">
                            <div className="col-sm-4">
                                <Link to="." className="btn btn-secondary btn-block">
                                    Cancel
                                </Link>
                            </div>
                            <div className="col-sm-4">
                                <button type="submit" className="btn btn-primary btn-block" disabled={updating}>
                                    {updating ? "Saving..." : "Save"}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        ) : null
    }
}

EditHistory = reduxForm({
    form: "editMedicalHistory"
})(EditHistory)

EditHistory = connect(
    state => ({
        dateOfBirth: state.patient.patient.dateOfBirth,
        codesLoading: state.codes.loading,
        initialValues: state.patient.patient,
        updating: state.patient.updating,
        canEdit: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[WRITE]
    }),
    {
        getCodes: getCodesAsOptions,
        loadCategories: loadCategoriesImport,
        updatePatient
    }
)(EditHistory)

let HistoryRoutes = ({ match, canSee, canEdit }) => (
    <div>
        {canSee ? <Route exact path={match.url} component={History} /> : null}
        {canEdit ? <Route exact path={match.url + "/edit"} component={EditHistory} /> : null}
    </div>
)

export default connect(
    state => ({
        canSee: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[READ],
        canEdit: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[WRITE]
    }),
    {}
)(HistoryRoutes)
