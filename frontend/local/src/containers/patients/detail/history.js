import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import { reduxForm, Field, FieldArray, Fields } from "redux-form"

import { renderSelect, renderHabitFields } from "shared/forms/renderField"
import { renderMedications, renderSurgeries, renderInjuries, renderChronicDiseases, renderImmunizations, renderAllergies } from "../new/step3"
import { joinPaths } from "shared/utils"

//import "./style.css"

let History = ({ match, patient }) => (
    <div className="history">
        <header>
            <h1>Medical History</h1>
            <Link to={joinPaths(match.url, "edit")} className="btn btn-secondary btn-wide">
                Edit
            </Link>
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
                        <React.Fragment>
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
                        <React.Fragment>
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
                        <React.Fragment>
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
                        <React.Fragment>
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
                        <React.Fragment>
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
                        <React.Fragment>
                            <dt>{item.medication}</dt>
                            {item.comment && <dd>{item.comment}</dd>}
                        </React.Fragment>
                    ))}
                    {(patient.medications || []).length === 0 && <dd>None</dd>}
                </dl>
            </div>
        </div>

        <div className="section">
            <div className="name">Additional medications</div>
            <div className="values">
                <dl>
                    <dt>medication</dt>
                    <dd>comment</dd>
                </dl>
            </div>
        </div>
    </div>
)

History = connect(
    state => ({
        patient: state.patient.patient
    }),
    {}
)(History)

const bloodTypeOptions = [
    {
        label: "A+",
        value: "A+"
    },
    {
        label: "A-",
        value: "A-"
    },
    {
        label: "B+",
        value: "B+"
    },
    {
        label: "B-",
        value: "B-"
    },
    {
        label: "O+",
        value: "O+"
    },
    {
        label: "O-",
        value: "O-"
    },
    {
        label: "AB+",
        value: "AB+"
    },
    {
        label: "AB-",
        value: "AB-"
    }
]

let EditHistory = ({ match }) => (
    <div className="edit-history">
        <header>
            <h1>Edit Medical History</h1>
            <Link to="." className="btn btn-secondary btn-wide">
                Close
            </Link>
        </header>

        <div className="section blood-type">
            <h3>Blood type</h3>
            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="bloodType" options={bloodTypeOptions} component={renderSelect} label="Blood type" />
                    <p className="warning">Warning: Be very careful when entering blood type</p>
                </div>
            </div>
        </div>

        <div className="section">
            <h3>Allergies</h3>
            <FieldArray name="allergies" component={renderAllergies} />
        </div>

        <div className="section">
            <h3>Immunization</h3>
            <FieldArray name="immunizations" component={renderImmunizations} />
        </div>

        <div className="section">
            <h3>Chronic diseases</h3>
            <FieldArray name="chronicDiseases" component={renderChronicDiseases} />
        </div>

        <div className="section">
            <h3>Injuries &amp; handicaps</h3>
            <FieldArray name="injuries" component={renderInjuries} />
        </div>

        <div className="section">
            <h3>Surgeries</h3>
            <FieldArray name="surgeries" component={renderSurgeries} />
        </div>

        <div className="section">
            <h3>Additional medications</h3>
            <FieldArray name="medications" component={renderMedications} />
        </div>

        <div className="section">
            <h3>Habits</h3>

            <Fields label="Are you a smoker?" names={["habits_smoking", "habits_smoking_comment"]} commentWhen="true" component={renderHabitFields} />
            <Fields label="Are you taking drugs?" names={["habits_drugs", "habits_drugs_comment"]} commentWhen="true" component={renderHabitFields} />
        </div>

        <div className="section">
            <h3>Conditions</h3>

            <Fields
                label="Do you have resources for basic hygiene?"
                names={["conditions_basic_hygiene", "conditions_basic_hygiene_comment"]}
                component={renderHabitFields}
            />

            <Fields
                label="Do you have access to clean water?"
                names={["conditions_clean_water", "conditions_clean_water_comment"]}
                component={renderHabitFields}
            />

            <Fields
                label="Do you have sufficient food supply?"
                names={["conditions_food_supply", "conditions_food_supply_comment"]}
                component={renderHabitFields}
            />

            <Fields
                label="Do you have a good appetite?"
                names={["conditions_good_appetite", "conditions_good_appetite_comment"]}
                component={renderHabitFields}
            />

            <Fields label="Does your tent have heating?" names={["conditions_heating", "conditions_heating_comment"]} component={renderHabitFields} />

            <Fields
                label="Does your tent have electricity?"
                names={["conditions_electricity", "conditions_electricity_comment"]}
                component={renderHabitFields}
            />
        </div>
    </div>
)

EditHistory = reduxForm({
    form: "editMedicalHistory"
})(EditHistory)

EditHistory = connect(
    state => ({
        initialState: state.patient.patient
    }),
    {}
)(EditHistory)

export default ({ match }) => (
    <div>
        <Route exact path={match.url} component={History} />
        <Route exact path={match.url + "/edit"} component={EditHistory} />
    </div>
)
