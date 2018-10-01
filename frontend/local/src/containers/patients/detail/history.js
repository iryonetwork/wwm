import React from "react"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import { reduxForm } from "redux-form"
import moment from "moment"
import classnames from "classnames"

import { RenderForm } from "../new/step3"
import { joinPaths } from "shared/utils"
import { RESOURCE_HEALTH_HISTORY, READ, WRITE } from "../../../modules/validations"
import { getCodesAsOptions, getCodes, loadCategories as loadCategoriesImport } from "shared/modules/codes"
import { updatePatient } from "../../../modules/patient"
import { BABY_MAX_AGE, CHILD_MAX_AGE } from "../../../modules/config"
import Column from "shared/containers/valueColumn"

let History = ({ match, patient, maxBabyAge, maxChildAge, canSee, canEdit }) => {
    const age = moment().diff(moment(patient.dateOfBirth), "years")

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
            {age <= maxBabyAge && <BabyHistory patient={patient} />}
            {maxBabyAge < age && age <= maxChildAge && <ChildHistory patient={patient} />}
            {age > maxChildAge && <AdultHistory patient={patient} />}
        </div>
    ) : null
}

History = connect(
    state => ({
        patient: state.patient.patient,
        maxBabyAge: state.config[BABY_MAX_AGE],
        maxChildAge: state.config[CHILD_MAX_AGE],
        canSee: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[READ],
        canEdit: ((state.validations.userRights || {})[RESOURCE_HEALTH_HISTORY] || {})[WRITE]
    }),
    {}
)(History)

let AdultHistory = ({ patient }) => (
    <React.Fragment>
        {/* @TODO confirm if needed <div className="section">
                <div className="name">Blood type</div>
                <div className="values">A+</div>
            </div> */}
        <Allergies patient={patient} />
        <Immunization patient={patient} />
        <ChronicDiseases patient={patient} />
        <Injuries patient={patient} />
        <Surgeries patient={patient} />
        <Medications patient={patient} />
        <Habits patient={patient} />
        <Conditions patient={patient} />
    </React.Fragment>
)

let BabyHistory = ({ patient }) => (
    <React.Fragment>
        {/* @TODO confirm if needed <div className="section">
                <div className="name">Blood type</div>
                <div className="values">A+</div>
            </div> */}
        <BirthData patient={patient} />
        <Allergies patient={patient} />
        <Immunization patient={patient} />
        <ChronicDiseases patient={patient} />
        <Injuries patient={patient} />
        <Surgeries patient={patient} />
        <Medications patient={patient} />
        <BabyHabits patient={patient} />
        <BabyConditions patient={patient} />
    </React.Fragment>
)

let ChildHistory = ({ patient }) => (
    <React.Fragment>
        {/* @TODO confirm if needed <div className="section">
                <div className="name">Blood type</div>
                <div className="values">A+</div>
            </div> */}
        <ChildVaccination patient={patient} />
        <Allergies patient={patient} />
        <Immunization patient={patient} />
        <ChronicDiseases patient={patient} />
        <Injuries patient={patient} />
        <Surgeries patient={patient} />
        <Medications patient={patient} />
        <Habits patient={patient} />
        <Conditions patient={patient} />
    </React.Fragment>
)

let BirthData = ({ patient, fetchCodes }) => {
    return patient.deliveryType || patient.prematurity || patient.weeksAtBirth || patient.weightAtBirth || patient.heightAtBirth ? (
        <div className="section">
            <div className="name">Birth data</div>
            <div className="values">
                <div className="row">
                    <Column width="4" label="Delivery type" value={patient.deliveryType} key="deliveryType" codes={fetchCodes("deliveryType")} />
                    <Column width="4" label="Prematurity" value={patient.prematurity === "true" ? "Yes" : "No"} key="prematurity" />
                </div>

                <div className="row">
                    <Column width="4" label="Weeks at birth" value={patient.weeksAtBirth} unit="weeks" key="weeksAtBirth" />
                    <Column width="4" label="Weight at birth" value={patient.weightAtBirth} unit="grams" key="weightAtBirth" />
                    <Column width="4" label="Height at birth" value={patient.heightAtBirth} unit="cm" key="heightAtBirth" />
                </div>
            </div>
        </div>
    ) : null
}

BirthData = connect(state => ({}), {
    fetchCodes: getCodes
})(BirthData)

let BabyHabits = ({ patient, fetchCodes }) => {
    return patient.breastfeeding ||
        patient.babyEatsAndDrinks ||
        patient.babyWetDiapers ||
        patient.babyBowelMovements ||
        patient.babyBowelMovementsComment ||
        patient.babySleep ||
        patient.babySleepOnBack ||
        patient.babyVitaminD ||
        patient.babyGetsAround ||
        patient.babyCommunicates ||
        patient.babyAnyoneSmokes ? (
        <div className="section">
            <div className="name">Habits</div>
            <div className="values">
                <dl>
                    {patient.breastfeeding && (
                        <div className="item" key="breastfeeding">
                            <dt
                                className={classnames({
                                    positive: patient.breastfeeding === "true",
                                    negative: patient.breastfeeding === "false"
                                })}
                            >
                                {patient.breastfeeding === "true"
                                    ? "Breastfeeding" + (patient.breastfeedingDuration ? ` for ${patient.breastfeedingDuration} weeks.` : ".")
                                    : "No breastfeeding."}
                            </dt>
                        </div>
                    )}
                    {patient.babyEatsAndDrinks && (
                        <div className="item" key="babyEatsAndDrinks">
                            <div className="row">
                                <Column
                                    width="8"
                                    label="What does your baby eat and drink?"
                                    value={patient.babyEatsAndDrinks}
                                    key="babyEatsAndDrinks"
                                    codes={fetchCodes("babyFood")}
                                />
                            </div>
                        </div>
                    )}
                    {patient.babyWetDiapers && (
                        <div className="item" key="babyWetDiapers">
                            <div className="row">
                                <Column width="8" label="How many diapers does your child wet in 24h?" value={patient.babyWetDiapers} key="babyWetDiapers" />
                            </div>
                        </div>
                    )}
                    {patient.babyBowelMovements && (
                        <div className="item" key="babyBowelMovements">
                            <div className="row">
                                <Column
                                    width="8"
                                    label="How frequent does your baby have bowel movements?"
                                    value={patient.babyBowelMovements}
                                    key="babyBowelMovements"
                                />
                            </div>
                        </div>
                    )}
                    {patient.babyBowelMovementsComment && (
                        <div className="item" key="babyBowelMovementsComment">
                            <div className="row">
                                <Column
                                    width="8"
                                    label="Describe baby's bowel movements"
                                    value={patient.babyBowelMovementsComment}
                                    key="babyBowelMovementsComment"
                                />
                            </div>
                        </div>
                    )}
                    {patient.babySleep && (
                        <div className="item" key="babySleep">
                            <dt
                                className={classnames({
                                    positive: patient.babySleep === "true",
                                    negative: patient.babySleep === "false"
                                })}
                            >
                                {patient.babySleep === "true" ? "Satisfied with child's sleep." : "Not satisfied with child's sleep."}
                            </dt>
                            {patient.babySleepComment && <dd className="withIconComment">{patient.babySleepComment}</dd>}
                        </div>
                    )}
                    {patient.babySleepsOnBack && (
                        <div className="item" key="babySleepsOnBack">
                            <dt
                                className={classnames({
                                    positive: patient.babySleepsOnBack === "true",
                                    negative: patient.babySleepsOnBack === "false"
                                })}
                            >
                                {patient.babySleepsOnBack === "true" ? "Baby sleeps on its back." : "Baby does not sleep on its back."}
                            </dt>
                        </div>
                    )}
                    {patient.babyVitaminD && (
                        <div className="item" key="babyVitaminD">
                            <dt
                                className={classnames({
                                    positive: patient.babyVitaminD === "true",
                                    negative: patient.babyVitaminD === "false"
                                })}
                            >
                                {patient.babyVitaminD === "true" ? "Baby takes vitamin D." : "Baby does not take vitamin D."}
                            </dt>
                        </div>
                    )}
                    {patient.babyGetsAround && (
                        <div className="item" key="babyGetsAround">
                            <div className="row">
                                <Column width="8" label="How does your child get around?" value={patient.babyGetsAround} key="babyGetsAround" />
                            </div>
                        </div>
                    )}
                    {patient.babyCommunicates && (
                        <div className="item" key="babyCommunicates">
                            <div className="row">
                                <Column
                                    width="8"
                                    label="How does your child communicate?"
                                    value={patient.babyCommunicates}
                                    key="babyCommunicates"
                                    codes={fetchCodes("childCommunication")}
                                />
                            </div>
                        </div>
                    )}
                    {patient.babyAnyoneSmokes && (
                        <div className="item" key="babyAnyoneSmokes">
                            <dt
                                className={classnames({
                                    positive: patient.babyAnyoneSmokes === "false",
                                    negative: patient.babyAnyoneSmokes === "true"
                                })}
                            >
                                {patient.babySleep === "true" ? "No smokers in the house." : "Smokers in the house."}
                            </dt>
                            {patient.babyNumberOfSmokers && <dd className="withIconComment">{patient.babyNumberOfSmokers} smokers.</dd>}
                        </div>
                    )}
                </dl>
            </div>
        </div>
    ) : null
}

BabyHabits = connect(state => ({}), {
    fetchCodes: getCodes
})(BabyHabits)

let ChildVaccination = ({ patient }) => {
    return patient.vaccinationUpToDate ||
        patient.vaccinationCertificates ||
        patient.tuberculosisTested ||
        patient.tuberculosisTestResult ||
        patient.vaccinationReaction ? (
        <div className="section">
            <div className="name">Vaccination</div>
            <div className="values">
                <dl>
                    {patient.vaccinationUpToDate && (
                        <div className="item" key="vaccinationUpToDate">
                            <dt
                                className={classnames({
                                    positive: patient.vaccinationUpToDate === "true",
                                    negative: patient.vaccinationUpToDate === "false"
                                })}
                            >
                                {patient.vaccinationUpToDate === "true"
                                    ? "Up to date with the home country vaccination schedule."
                                    : "Not up to date with the home country vaccination schedule."}
                            </dt>
                        </div>
                    )}
                    {patient.vaccinationCertificates && (
                        <div className="item" key="vaccinationCertificates">
                            <dt
                                className={classnames({
                                    positive: patient.vaccinationCertificates === "true",
                                    negative: patient.vaccinationCertificates === "false"
                                })}
                            >
                                {patient.vaccinationCertificates === "true"
                                    ? "Up to date with the home country vaccination schedule."
                                    : "Not up to date with the home country vaccination schedule."}
                            </dt>
                        </div>
                    )}
                    {patient.tuberculosisTested && (
                        <div className="item" key="tuberculosisTested">
                            <dt
                                className={classnames({
                                    positive: patient.tuberculosisTested === "true",
                                    negative: patient.tuberculosisTested === "false"
                                })}
                            >
                                {patient.tuberculosisTested === "true"
                                    ? "Child has been tested for tuberculosis."
                                    : "Child has not been tested for tuberculosis."}
                            </dt>
                        </div>
                    )}
                    {patient.tuberculosisTestResult && (
                        <div className="item" key="tuberculosisTestResult">
                            <dt
                                className={classnames({
                                    danger: patient.tuberculosisTestResult === "true"
                                })}
                            >
                                {patient.tuberculosisTestResult === "true" ? "Positive result for tuberculosis." : "Negative result for tuberculosis."}
                            </dt>
                            {patient.tuberculosisAdditionalInvestigationDetails && (
                                <dd className="withIconComment">{patient.tuberculosisAdditionalInvestigationDetails}</dd>
                            )}
                        </div>
                    )}
                    {patient.vaccinationReaction && (
                        <div className="item" key="vaccinationReaction">
                            <dt
                                className={classnames({
                                    danger: patient.vaccinationReaction === "true"
                                })}
                            >
                                {patient.vaccinationReaction === "true"
                                    ? "Child has experienced vaccination reaction."
                                    : "Child has not experienced vaccination reaction."}
                            </dt>
                            {patient.vaccinationReactionDetails && <dd className="withIconComment">{patient.vaccinationReactionDetails}</dd>}
                        </div>
                    )}
                </dl>
            </div>
        </div>
    ) : null
}

let Allergies = ({ patient }) => {
    return (patient.allergies || []).length === 0 ? null : (
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
                </dl>
            </div>
        </div>
    )
}

let Immunization = ({ patient }) => {
    return (patient.immunizations || []).length === 0 ? null : (
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
                </dl>
            </div>
        </div>
    )
}

let ChronicDiseases = ({ patient }) => {
    return (patient.chronicDiseases || []).length === 0 ? null : (
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
                </dl>
            </div>
        </div>
    )
}

let Injuries = ({ patient }) => {
    return (patient.injuries || []).length === 0 ? null : (
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
                </dl>
            </div>
        </div>
    )
}

let Surgeries = ({ patient }) => {
    return (patient.surgeries || []).length === 0 ? null : (
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
                </dl>
            </div>
        </div>
    )
}

let Medications = ({ patient }) => {
    return (patient.medications || []).length === 0 ? null : (
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
    )
}

let Habits = ({ patient }) => {
    return patient.habits_smoking === "true" || patient.habits_drugs === "true" ? (
        <div className="section">
            <div className="name">Habits</div>
            <div className="values">
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
            </div>
        </div>
    ) : null
}

let BabyConditions = ({ patient }) => {
    return patient.conditions_basic_hygiene ||
        patient.conditions_clean_water ||
        patient.conditions_electricity ||
        patient.conditions_food_supply ||
        patient.conditions_heating ? (
        <div className="section">
            <div className="name">Conditions</div>
            <div className="values">
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
                            {patient.conditions_basic_hygiene_comment && <dd className="withIconComment">{patient.conditions_basic_hygiene_comment}</dd>}
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
                                {patient.conditions_electricity === "true" ? "Accomodation has electricity." : "Accomodation does not have electricity."}
                            </dt>
                            {patient.conditions_electricity_comment && <dd className="withIconComment">{patient.conditions_electricity_comment}</dd>}
                        </div>
                    )}
                </dl>
            </div>
        </div>
    ) : null
}

let Conditions = ({ patient }) => {
    return patient.conditions_basic_hygiene ||
        patient.conditions_clean_water ||
        patient.conditions_electricity ||
        patient.conditions_food_supply ||
        patient.conditions_good_appetite ||
        patient.conditions_heating ? (
        <div className="section">
            <div className="name">Conditions</div>
            <div className="values">
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
                            {patient.conditions_basic_hygiene_comment && <dd className="withIconComment">{patient.conditions_basic_hygiene_comment}</dd>}
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
                            {patient.conditions_good_appetite_comment && <dd className="withIconComment">{patient.conditions_good_appetite_comment}</dd>}
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
                                {patient.conditions_electricity === "true" ? "Accomodation has electricity." : "Accomodation does not have electricity."}
                            </dt>
                            {patient.conditions_electricity_comment && <dd className="withIconComment">{patient.conditions_electricity_comment}</dd>}
                        </div>
                    )}
                </dl>
            </div>
        </div>
    ) : null
}

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
