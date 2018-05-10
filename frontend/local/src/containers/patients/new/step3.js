import React, { Component } from "react"
import { Field, Fields, FieldArray, formValueSelector, reduxForm } from "redux-form"
import { connect } from "react-redux"
import classnames from "classnames"

import Footer from "./footer"
import validate from "./validate"
import {
    renderInput,
    renderHabitFields,
    renderRadio,
    renderSelect,
    renderHorizontalInput,
    renderHorizontalSelect,
    renderHorizontalRadio
} from "shared/forms/renderField"
import { yesNoOptions, positiveNegativeOptions } from "shared/forms/options"
import { read, BABY_MAX_AGE, CHILD_MAX_AGE } from "shared/modules/config"
import { getCodesAsOptions, loadCategories as loadCategoriesImport } from "shared/modules/codes"

import { ReactComponent as RemoveIcon } from "shared/icons/negative.svg"

const numberOptions = Array.from(Array(9), (x, i) => ({
    label: i,
    value: i
}))

class Step3 extends Component {
    componentWillMount() {
        this.props.loadCategories("babyFood", "childCommunication", "deliveryType")
    }

    render() {
        const { handleSubmit, reset, previousPage, dateOfBirth, codesLoading, getCodes } = this.props
        return (
            <form onSubmit={handleSubmit}>
                <RenderForm
                    dateOfBirth={dateOfBirth}
                    babyFoods={getCodes("babyFood")}
                    communicationTypes={getCodes("childCommunication")}
                    deliveryTypes={getCodes("deliveryType")}
                    codesLoading={codesLoading}
                />

                <Footer reset={reset} previousPage={previousPage} />
            </form>
        )
    }
}

let RenderForm = ({ dateOfBirth, babyFoods, communicationTypes, deliveryTypes, codesLoading, maxBabyAge, maxChildAge }) => {
    if (codesLoading) {
        return null
    }

    const age = (Date.now() - new Date(dateOfBirth).getTime()) / (1000 * 60 * 60 * 24 * 365)

    if (age <= maxBabyAge) {
        return renderBabyForm({ babyFoods, deliveryTypes, communicationTypes })
    } else if (age <= maxChildAge) {
        return renderChildForm()
    } else {
        return renderAdultForm()
    }
}

RenderForm = connect(
    state => ({
        maxBabyAge: state.config[BABY_MAX_AGE],
        maxChildAge: state.config[CHILD_MAX_AGE]
    }),
    {}
)(RenderForm)

const renderAdultForm = props => (
    <div className="modal-body">
        <HealthAttributes />
        <HabitsAndLivingConditions />
    </div>
)

const renderBabyForm = ({ babyFoods, deliveryTypes, communicationTypes }) => (
    <div className="modal-body">
        <h3>Birth data</h3>
        <div className="baby-form">
            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="deliveryType" component={renderSelect} options={deliveryTypes} label="Delivery type" />
                </div>
                <div className="form-group col-sm-8">
                    <Field name="prematurity" options={yesNoOptions} component={renderRadio} label="Prematurity?" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-2">
                    <Field name="weeksAtBirth" component={renderInput} label="Weeks at birth" />
                </div>
                <div className="col-sm-2 unit">weeks</div>
                <div className="form-group col-sm-2">
                    <Field name="weightAtBirth" component={renderInput} label="Weight at birth" />
                </div>
                <div className="col-sm-2 unit">grams</div>
                <div className="form-group col-sm-2">
                    <Field name="heightAtBirth" component={renderInput} label="Height at birth" />
                </div>
                <div className="col-sm-2 unit">cm</div>
            </div>
        </div>
        <HealthAttributes />
        <h3>Habits and living conditions</h3>
        <Field name="breastfeeding" component={renderHorizontalRadio} options={yesNoOptions} label="Breastfeeding?" />
        <Field name="breastfeedingDuration" component={renderHorizontalInput} label="For how long?" />
        <Field name="babyEatsAndDrinks" component={renderHorizontalSelect} options={babyFoods} label="What does your baby eat and drink?" /> {/*@TODO codes */}
        <Field name="babyWetDiapers" component={renderHorizontalSelect} options={numberOptions} label="How many diapers does your child wet in 24h?" />
        <Field
            name="babyBowelMovements"
            component={renderHorizontalSelect}
            options={numberOptions}
            label="How frequent does your baby have bowel movements?"
        />{" "}
        {/*@ TODO codes */}
        <Field name="babyBowelMovementsComment" component={renderHorizontalInput} label="Describe baby's bowel movements" />
        <Field name="babySleep" component={renderHorizontalRadio} options={yesNoOptions} label="Are you satisfied with child's sleep?" />
        <Field name="babySleepComment" component={renderHorizontalInput} label="Comment" hideLabel={true} />
        <Field name="babyVitaminD" component={renderHorizontalRadio} options={yesNoOptions} label="Do you or your baby take vitamin D?" />
        <Field name="babySleepOnBack" component={renderHorizontalRadio} options={yesNoOptions} label="Does your baby sleep on her back?" />
        <Field name="babyAnyoneSmokes" component={renderHorizontalRadio} options={yesNoOptions} label="Does anyone at your house smoke?" />
        <Field name="babyNumberOfSmokers" component={renderHorizontalSelect} options={numberOptions} label="How many smokers?" />
        <Field name="babyGetsAround" component={renderHorizontalInput} label="How does your child get around?" />
        <Field name="babyCommunicates" component={renderHorizontalSelect} options={communicationTypes} label="How does your child communicate?" />
        <Field name="conditions_clean_water" component={renderHorizontalRadio} options={yesNoOptions} label="Do you have access to clean water?" />
        <Field name="conditions_food_supply" component={renderHorizontalRadio} options={yesNoOptions} label="Do you have sufficient food supply?" />
        <Field name="conditions_heating" component={renderHorizontalRadio} options={yesNoOptions} label="Does your tent have heating?" />
        <Field name="conditions_electricity" component={renderHorizontalRadio} options={yesNoOptions} label="Does your tent have electricity?" />
    </div>
)

const renderChildForm = () => (
    <div className="modal-body">
        <h3>Vaccine information</h3>

        <Field
            name="vaccinationUpToDate"
            component={renderHorizontalRadio}
            options={yesNoOptions}
            label="Was this child up to date with the home country vaccination schedule?"
        />
        <Field
            name="vaccinationCertificates"
            component={renderHorizontalRadio}
            options={yesNoOptions}
            label="Do you have this child's immunization certificates?"
        />
        <Field name="tuberculosisTested" component={renderHorizontalRadio} options={yesNoOptions} label="Has this child been tested for tuberculosis?" />
        <Field name="tuberculosisTestResult" component={renderHorizontalRadio} options={positiveNegativeOptions} label="• What was the result?" />
        <Field name="tuberculosisAdditionalInvestigationDetails" component={renderHorizontalInput} label="• Investigation details" />
        <Field
            name="tuberculosisAdditionalInvestigation"
            component={renderHorizontalRadio}
            options={yesNoOptions}
            label="• Any additional investigation done?"
        />
        <Field name="vaccinationReaction" component={renderHorizontalRadio} options={yesNoOptions} label="Has the child ever experienced vaccine reaction?" />
        <Field name="vaccinationReactionDetails" component={renderHorizontalRadio} options={yesNoOptions} label="• Any additional investigation done?" />

        <HealthAttributes />

        <HabitsAndLivingConditions />
    </div>
)

const HealthAttributes = () => (
    <div>
        <h3>Permanent health attributes</h3>
        <FieldArray name="allergies" component={renderAllergies} />
        <FieldArray name="immunizations" component={renderImmunizations} />
        <FieldArray name="chronicDiseases" component={renderChronicDiseases} />
        <FieldArray name="injuries" component={renderInjuries} />
        <FieldArray name="surgeries" component={renderSurgeries} />
        <FieldArray name="medications" component={renderMedications} />
    </div>
)

const HabitsAndLivingConditions = () => (
    <div>
        <h3>Habits and living conditions</h3>

        <Fields label="Are you a smoker?" names={["habits_smoking", "habits_smoking_comment"]} commentWhen="true" component={renderHabitFields} />
        <Fields label="Are you taking drugs?" names={["habits_drugs", "habits_drugs_comment"]} commentWhen="true" component={renderHabitFields} />

        <Fields
            label="Do you have resources for basic hygiene?"
            names={["conditions_basic_hygiene", "conditions_basic_hygiene_comment"]}
            component={renderHabitFields}
        />

        <Fields label="Do you have access to clean water?" names={["conditions_clean_water", "conditions_clean_water_comment"]} component={renderHabitFields} />

        <Fields
            label="Do you have sufficient food supply?"
            names={["conditions_food_supply", "conditions_food_supply_comment"]}
            component={renderHabitFields}
        />

        <Fields label="Do you have a good appetite?" names={["conditions_good_appetite", "conditions_good_appetite_comment"]} component={renderHabitFields} />

        <Fields label="Does your tent have heating?" names={["conditions_heating", "conditions_heating_comment"]} component={renderHabitFields} />

        <Fields label="Does your tent have electricity?" names={["conditions_electricity", "conditions_electricity_comment"]} component={renderHabitFields} />
    </div>
)

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

Step3 = reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step3)

const selector = formValueSelector("newPatient")
Step3 = connect(
    state => ({
        dateOfBirth: selector(state, "dateOfBirth"),
        codesLoading: state.codes.loading
    }),
    {
        getCodes: getCodesAsOptions,
        loadCategories: loadCategoriesImport
    }
)(Step3)

export default Step3
