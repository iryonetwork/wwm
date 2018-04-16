import React from "react"
import { Field, Fields, FieldArray, reduxForm } from "redux-form"
import classnames from "classnames"

import Footer from "./footer"
import validate from "./validate"
import { renderInput, renderHabitFields, renderRadio } from "shared/forms/renderField"
import { yesNoOptions } from "shared/forms/options"

import { ReactComponent as RemoveIcon } from "shared/icons/negative.svg"

const Step3 = props => {
    const { handleSubmit, reset, previousPage } = props
    return (
        <form onSubmit={handleSubmit}>
            <div className="modal-body">
                <h3>Permanent health attributes</h3>
                <FieldArray name="allergies" component={renderAllergies} />
                <FieldArray name="immunizations" component={renderImmunizations} />
                <FieldArray name="chronicDiseases" component={renderChronicDiseases} />
                <FieldArray name="injuries" component={renderInjuries} />
                <FieldArray name="surgeries" component={renderSurgeries} />
                <FieldArray name="medications" component={renderMedications} />

                <h3>Habits and living conditions</h3>

                <Fields label="Are you a smoker?" names={["habits_smoking", "habits_smoking_comment"]} commentWhen="true" component={renderHabitFields} />
                <Fields label="Are you taking drugs?" names={["habits_drugs", "habits_drugs_comment"]} commentWhen="true" component={renderHabitFields} />

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

            <Footer reset={reset} previousPage={previousPage} />
        </form>
    )
}

const renderAllergies = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((allergy, index) => (
            <div key={allergy} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${allergy}.allergy`} component={renderInput} label="Allergy" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${allergy}.comment`} optional={true} component={renderInput} label="Comment" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${allergy}.critical`} options={yesNoOptions} component={renderRadio} label="Critical?" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add allergy
            </a>
        </div>
    </div>
)

const renderImmunizations = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((immunization, index) => (
            <div key={immunization} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${immunization}.immunization`} component={renderInput} label="Immunization" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${immunization}.date`} type="date" component={renderInput} label="Date" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add immunization
            </a>
        </div>
    </div>
)

const renderChronicDiseases = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((disease, index) => (
            <div key={disease} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${disease}.disease`} component={renderInput} label="Disease" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${disease}.date`} type="date" component={renderInput} label="Date" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${disease}.medication`} optional={true} component={renderInput} label="Medication" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add chronic disease
            </a>
        </div>
    </div>
)

const renderInjuries = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((injury, index) => (
            <div key={injury} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.injury`} component={renderInput} label="Injury or handicap" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.date`} type="date" component={renderInput} label="Date" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.medication`} optional={true} component={renderInput} label="Prosthetics &amp; aids" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add injury or handicap
            </a>
        </div>
    </div>
)

const renderSurgeries = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((injury, index) => (
            <div key={injury} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.injury`} component={renderInput} label="Surgery" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.date`} type="date" component={renderInput} label="Date" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${injury}.medication`} optional={true} component={renderInput} label="Comment" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add surgery
            </a>
        </div>
    </div>
)

const renderMedications = ({ fields, meta: { error, submitFailed } }) => (
    <div className={classnames("attributes", { open: fields.length })}>
        {fields.map((medication, index) => (
            <div key={medication} className="form-row">
                <div className="form-group col-sm-4">
                    <Field name={`${medication}.medication`} component={renderInput} label="Additional medication" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name={`${medication}.comment`} optional={true} component={renderInput} label="Comment" />
                </div>

                <a
                    onClick={e => {
                        e.preventDefault()
                        fields.remove(index)
                    }}
                    href="/"
                    className="remove"
                >
                    <RemoveIcon />
                    Remove
                </a>
            </div>
        ))}
        <div className="link">
            <a
                href="/"
                onClick={e => {
                    e.preventDefault()
                    fields.push({})
                }}
            >
                Add additional medication
            </a>
        </div>
    </div>
)

export { renderMedications, renderSurgeries, renderInjuries, renderChronicDiseases, renderImmunizations, renderAllergies }

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step3)
