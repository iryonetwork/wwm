import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import { reduxForm } from "redux-form"
import moment from "moment"
import classnames from "classnames"

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
                            <div className="item" key={i}>
                                <dt className={classnames({ danger: item.critical === "true" })}>{item.allergy}</dt>
                                {(item.critical === "true" || item.comment) && (
                                    <dd className={classnames({ withIconComment: item.critical === "true" })}>
                                        {item.critical === "true" && <div>High risk</div>}
                                        {item.comment && <div>{item.comment}</div>}
                                    </dd>
                                )}
                            </div>
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
                            <div className="item" key={i}>
                                <dt>{item.immunization}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{moment(item.date).format("Do MMMM Y")}</dd>}
                            </div>
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
                            <div className="item" key={i}>
                                <dt className="danger">{item.disease}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd className="withIconComment">{moment(item.date).format("Do MMMM Y")}</dd>}
                            </div>
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
                            <div className="item" key={i}>
                                <dt>{item.injury}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{moment(item.date).format("Do MMMM Y")}</dd>}
                            </div>
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
                            <div className="item" key={i}>
                                <dt>{item.injury}</dt>
                                {/* @TODO format date */}
                                {item.date && <dd>{moment(item.date).format("Do MMMM Y")}</dd>}
                            </div>
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
                            <div className="item" key={i}>
                                <dt>{item.medication}</dt>
                                {item.comment && <dd>{item.comment}</dd>}
                            </div>
                        ))}
                        {(patient.medications || []).length === 0 && <dd>None</dd>}
                    </dl>
                </div>
            </div>

            <div className="section">
                <div className="name">Habits</div>
                <div className="values">
                    {patient.habits_smoking === "true" || patient.habits_drugs === "true" ? (
                        <dl>
                            {patient.habits_smoking === "true" && (
                                <div className="item" key="habits_smoking">
                                    <dt className="danger">Smoker</dt>
                                    {patient.habits_smoking_comment && <dd className="withIconComment">{patient.habits_smoking_comment}</dd>}
                                </div>
                            )}
                            {patient.habits_drugs === "true" && (
                                <div className="item" key="habits_drugs">
                                    <dt className="danger">Taking drugs</dt>
                                    {patient.habits_drugs_comment && <dd className="withIconComment">{patient.habits_drugs_comment}</dd>}
                                </div>
                            )}
                        </dl>
                    ) : (
                        <dl>
                            <dd>None</dd>
                        </dl>
                    )}
                </div>
            </div>
            <div className="section">
                <div className="name">Conditions</div>
                <div className="values">
                    {patient.conditions_basic_hygiene ||
                    patient.conditions_clean_water ||
                    patient.conditions_electricity ||
                    patient.conditions_food_supply ||
                    patient.conditions_good_appetite ||
                    patient.conditions_heating ? (
                        <dl>
                            {patient.conditions_basic_hygiene && (
                                <div className="item" key="conditions_basic_hygiene">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_basic_hygiene === "true",
                                            negative: patient.conditions_basic_hygiene === "false"
                                        })}
                                    >
                                        {patient.conditions_basic_hygiene === "true"
                                            ? "Has resources for basic hygiene."
                                            : "Does not have resources for basic hygiene."}
                                    </dt>
                                    {patient.conditions_basic_hygiene_comment && (
                                        <dd className="withIconComment">{patient.conditions_basic_hygiene_comment}</dd>
                                    )}
                                </div>
                            )}
                            {patient.conditions_clean_water && (
                                <div className="item" key="conditions_clean_water">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_clean_water === "true",
                                            negative: patient.conditions_clean_water === "false"
                                        })}
                                    >
                                        {patient.conditions_clean_water === "true" ? "Has access to clean water." : "No access to clean water."}
                                    </dt>
                                    {patient.conditions_clean_water_comment && <dd className="withIconComment">{patient.conditions_clean_water_comment}</dd>}
                                </div>
                            )}
                            {patient.conditions_food_supply && (
                                <div className="item" key="conditions_food_supply">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_food_supply === "true",
                                            negative: patient.conditions_food_supply === "false"
                                        })}
                                    >
                                        {patient.conditions_food_supply === "true" ? "Has sufficient food supply." : "Does not have sufficient food supply."}
                                    </dt>
                                    {patient.conditions_food_supply_comment && <dd className="withIconComment">{patient.conditions_food_supply_comment}</dd>}
                                </div>
                            )}
                            {patient.conditions_good_appetite && (
                                <div className="item" key="conditions_good_appetite">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_good_appetite === "true",
                                            negative: patient.conditions_good_appetite === "false"
                                        })}
                                    >
                                        {patient.conditions_good_appetite === "true" ? "Has good appetite." : "Does not have good appetite."}
                                    </dt>
                                    {patient.conditions_good_appetite_comment && (
                                        <dd className="withIconComment">{patient.conditions_good_appetite_comment}</dd>
                                    )}
                                </div>
                            )}
                            {patient.conditions_heating && (
                                <div className="item" key="conditions_heating">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_heating === "true",
                                            negative: patient.conditions_heating === "false"
                                        })}
                                    >
                                        {patient.conditions_heating === "true" ? "Accomodation has heating." : "Accomodation does not have heating."}
                                    </dt>
                                    {patient.conditions_heating_comment && <dd className="withIconComment">{patient.conditions_heating_comment}</dd>}
                                </div>
                            )}
                            {patient.conditions_electricity && (
                                <div className="item" key="conditions_electricity">
                                    <dt
                                        className={classnames({
                                            positive: patient.conditions_electricity === "true",
                                            negative: patient.conditions_electricity === "false"
                                        })}
                                    >
                                        {patient.conditions_electricity === "true"
                                            ? "Accomodation has electricity."
                                            : "Accomodation does not have electricity."}
                                    </dt>
                                    {patient.conditions_electricity_comment && <dd className="withIconComment">{patient.conditions_electricity_comment}</dd>}
                                </div>
                            )}
                        </dl>
                    ) : (
                        <dl>
                            <dd>None</dd>
                        </dl>
                    )}
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
