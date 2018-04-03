import React from "react"
import { Field, FieldArray, reduxForm } from "redux-form"
import classnames from "classnames"

import Footer from "./footer"
import validate from "./validate"
import { renderInput, renderRadio } from "shared/forms/renderField"
import { yesNoOptions } from "./options"

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
                <FieldArray name="medications" component={renderMedications} />

                <h3>Habits and living conditions</h3>
                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="habits_smoking" options={yesNoOptions} component={renderRadio} label="Are you a smoker?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="habits_drugs" options={yesNoOptions} component={renderRadio} label="Are you taking drugs?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field
                            name="conditions_basic_hygiene"
                            options={yesNoOptions}
                            component={renderRadio}
                            label="Do you have resources for basic hygiene?"
                        />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="conditions_clean_water" options={yesNoOptions} component={renderRadio} label="Do you have access to clean water?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="conditions_food_supply" options={yesNoOptions} component={renderRadio} label="Do you have sufficient food supply?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="conditions_good_appetite" options={yesNoOptions} component={renderRadio} label="Do you have a good appetite?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="conditions_heating" options={yesNoOptions} component={renderRadio} label="Does your tent have heating?" />
                    </div>
                </div>

                <div className="form-row habits">
                    <div className="form-group col-sm-12">
                        <Field name="conditions_electricity" options={yesNoOptions} component={renderRadio} label="Does your tent have electricity?" />
                    </div>
                </div>
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
                    <Field name={`${immunization}.date`} component={renderInput} label="Date" />
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
                    <Field name={`${disease}.date`} component={renderInput} label="Date" />
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
                    <Field name={`${injury}.date`} component={renderInput} label="Date" />
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

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step3)
