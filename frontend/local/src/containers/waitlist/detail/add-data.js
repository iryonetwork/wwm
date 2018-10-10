import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { Field, reduxForm, submit, formValueSelector } from "redux-form"
import classnames from "classnames"
import { goBack, push } from "react-router-redux"
import moment from "moment"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient/card"
import Spinner from "shared/containers/spinner"
import { round } from "shared/utils"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { listAll, update } from "../../../modules/waitlist"
import { cardToObject } from "../../../modules/discovery"
import { SimpleUnitInput, UnitInputWithConversion } from "shared/forms/measurementFields"
import {
    required,
    heightExpectedRange,
    weightExpectedRange,
    temperatureExpectedRange,
    systolicBloodPressureExpectedRange,
    diastolicBloodPressureExpectedRange,
    heartRateExpectedRange,
    oxygenSaturationxpectedRange
} from "shared/forms/validation"
import { POUNDS_OUNCES, POUNDS, KG, MM_HG, CM_HG, CM, FEET_INCHES, CELSIUS, FAHRENHEIT } from "shared/unitConversion/units"
import {
    weightValueToObject,
    weightObjectToValue,
    lengthValueToObject,
    lengthObjectToValue,
    pressureValueToObject,
    pressureObjectToValue,
    temperatureValueToObject,
    temperatureObjectToValue
} from "shared/unitConversion"
import { WEIGHT_UNIT, LENGTH_UNIT, TEMPERATURE_UNIT, BLOOD_PRESSURE_UNIT } from "../../../modules/config"
import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"

const supportedSigns = ["height", "weight", "temperature", "pressure", "heart_rate", "oxygen_saturation"]
const complexSigns = { pressure: ["systolic", "diastolic"] }

const onMedicalDataFormSubmit = (form, dispatch, props) => {
    let vitalSigns = {}

    _.forEach(form, (value, key) => {
        if (key.indexOf("has_") === 0 && value) {
            let sign = key.slice(4)
            vitalSigns[sign] = {}
            vitalSigns[sign].value = form[sign]

            if (form[sign] !== props.initialValues[sign]) {
                vitalSigns[sign].timestamp = moment().format()
            } else {
                vitalSigns[sign].timestamp = props.initialValues["timestmap_" + sign]
            }
        }
    })

    // set BMI if height and weight available
    if (vitalSigns.height && vitalSigns.weight) {
        vitalSigns.bmi = {}
        vitalSigns.bmi.value = round(vitalSigns.weight.value / vitalSigns.height.value / vitalSigns.height.value * 10000, 2)
        vitalSigns.bmi.timestamp = moment.max(moment(vitalSigns.height.timestamp), moment(vitalSigns.weight.timestamp)).format()
    } else if (props.item.vitalSigns && props.item.vitalSigns.bmi) {
        // remove bmi if bmi cannot be set anymore
        delete vitalSigns.bmi
    }

    let newItem = _.clone(props.item)
    newItem.vitalSigns = vitalSigns

    dispatch(update(props.match.params.waitlistID, newItem))
        .then(data => {
            dispatch(listAll(props.match.params.waitlistID))
            dispatch(goBack())
        })
        .catch(ex => {})
}

class MedicalData extends React.Component {
    constructor(props) {
        super(props)
        if (!props.item) {
            props.listAll(props.match.params.waitlistID)
        }
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.goBack()
            this.props.open("Waitlist item was not found", "", COLOR_DANGER, 5)
        }
    }

    render() {
        let { handleSubmit, item, change, dispatch, history } = this.props

        return (
            <Modal>
                <div className="medical-data">
                    {item ? (
                        <form onSubmit={handleSubmit}>
                            <div className="modal-header">
                                <Patient data={item.patient && cardToObject({ connections: item.patient })} />
                                <h1>
                                    <MedicalDataIcon />
                                    Add Medical Data
                                </h1>
                            </div>

                            <div className="modal-body">
                                <h3>Body Measurements</h3>
                                <div>
                                    <Height change={change} dispatch={dispatch} />
                                    <Weight change={change} dispatch={dispatch} />
                                </div>

                                <h3>Vital Signs</h3>

                                <div>
                                    <Temperature change={change} dispatch={dispatch} />
                                    <BloodPressure change={change} dispatch={dispatch} />
                                    <FieldWithUnit
                                        label="Heart rate"
                                        name="heart_rate"
                                        unit="bpm"
                                        change={change}
                                        dispatch={dispatch}
                                        validate={required}
                                        warn={heartRateExpectedRange}
                                    />
                                    <FieldWithUnit
                                        label="Oxygen saturation"
                                        name="oxygen_saturation"
                                        unit="%"
                                        min={0}
                                        max={100}
                                        change={change}
                                        dispatch={dispatch}
                                        validate={required}
                                        warn={oxygenSaturationxpectedRange}
                                    />
                                </div>
                            </div>

                            <div className="modal-footer">
                                <div className="form-row">
                                    <div className="col-sm-4" />
                                    <div className="col-sm-4">
                                        <button
                                            type="button"
                                            tabIndex="-1"
                                            className="btn btn-secondary btn-block"
                                            datadismiss="modal"
                                            onClick={() => history.goBack()}
                                        >
                                            Cancel
                                        </button>
                                    </div>

                                    <div className="col-sm-4">
                                        <button type="submit" className="float-right btn btn-primary btn-block">
                                            Add
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>
                    ) : (
                        <div className="modal-body">
                            <Spinner />
                        </div>
                    )}
                </div>
            </Modal>
        )
    }
}

MedicalData = reduxForm({
    form: "medical-data",
    onSubmit: onMedicalDataFormSubmit
})(MedicalData)

MedicalData = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]

        let initialValues = {}

        // initialize complex signs for correct validation (workaround for bug of handling '<Fields />' in redux-form)
        _.forEach(complexSigns, (obj, sign) => {
            initialValues[sign] = {}
            _.forEach(obj, key => {
                initialValues[sign][key] = undefined
            })
        })

        let selectedSign = props.match.params.sign
        if (selectedSign && _.includes(supportedSigns, selectedSign)) {
            initialValues["has_" + selectedSign] = true
            initialValues["focus"] = selectedSign
        }

        if (item) {
            _.forEach(item.vitalSigns || {}, (obj, key) => {
                initialValues[key] = obj.value
                initialValues["has_" + key] = true
                initialValues["timestmap_" + key] = obj.timestamp
            })
        }

        return {
            listed: state.waitlist.listed,
            item,
            initialValues
        }
    },
    {
        listAll,
        update,
        open,
        goBack,
        push
    }
)(MedicalData)

const selector = formValueSelector("medical-data")

const submitOnEnter = dispatch => e => {
    if (e.key === "Enter") {
        e.preventDefault()
        dispatch(submit("medical-data"))
    }
}

class FieldWithUnit extends React.Component {
    render() {
        return (
            <div
                className={classnames("section", {
                    open: this.props.opened
                })}
            >
                {this.props.opened && (
                    <div className="form-row">
                        <div className="col-sm">
                            <Field
                                name={this.props.name}
                                label={this.props.label}
                                placeholder={this.props.label}
                                component={SimpleUnitInput}
                                unit={this.props.unit}
                                precision={this.props.precision ? this.props.precision : 0}
                                min={this.props.min}
                                max={this.props.max}
                                validate={this.props.validate}
                                warn={this.props.warn}
                                onKeyPress={this.props.dispatch(submitOnEnter)}
                                autoFocus={this.props.focused}
                            />
                        </div>
                        <button className="btn btn-link remove" onClick={() => this.props.change(`has_${this.props.name}`, false)}>
                            <NegativeIcon />
                            Remove
                        </button>
                    </div>
                )}
                {!this.props.opened && (
                    <button
                        className="btn btn-link"
                        onClick={() => {
                            this.props.change(`has_${this.props.name}`, true)
                            this.props.change("focus", this.props.name)
                        }}
                    >
                        Add {this.props.label}
                    </button>
                )}
            </div>
        )
    }
}

FieldWithUnit = connect((state, props) => {
    return {
        opened: selector(state, `has_${props.name}`),
        focused: selector(state, "focus") === props.name
    }
}, {})(FieldWithUnit)

class Height extends React.Component {
    render() {
        return (
            <div
                className={classnames("section", {
                    open: this.props.opened
                })}
            >
                {this.props.opened && (
                    <div className="form-row">
                        <div className="col-sm">
                            {this.props.unit === FEET_INCHES ? (
                                <Field
                                    name="height"
                                    label="Height"
                                    placeholder="Height"
                                    component={UnitInputWithConversion}
                                    inputUnit={FEET_INCHES}
                                    valueUnit={CM}
                                    valuePrecision={0}
                                    valueToObject={lengthValueToObject}
                                    objectToValue={lengthObjectToValue}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={heightExpectedRange}
                                />
                            ) : (
                                <Field
                                    name="height"
                                    label="Height"
                                    placeholder="Height"
                                    component={SimpleUnitInput}
                                    unit={CM}
                                    precision={0}
                                    min={0}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={heightExpectedRange}
                                />
                            )}
                        </div>
                        <button className="btn btn-link remove" onClick={() => this.props.change("has_height", false)}>
                            <NegativeIcon />
                            Remove
                        </button>
                    </div>
                )}
                {!this.props.opened && (
                    <button
                        className="btn btn-link"
                        onClick={() => {
                            this.props.change("has_height", true)
                            this.props.change("focus", "height")
                        }}
                    >
                        Add Height
                    </button>
                )}
            </div>
        )
    }
}

Height = connect((state, props) => {
    return {
        opened: selector(state, "has_height"),
        focused: selector(state, "focus") === "height",
        unit: state.config[LENGTH_UNIT]
    }
}, {})(Height)

class Weight extends React.Component {
    render() {
        return (
            <div
                className={classnames("section", {
                    open: this.props.opened
                })}
            >
                {this.props.opened && (
                    <div className="form-row">
                        <div className="col-sm">
                            {this.props.unit === POUNDS_OUNCES ? (
                                <Field
                                    name="weight"
                                    label="Weight"
                                    placeholder="Weight"
                                    component={UnitInputWithConversion}
                                    inputUnit={POUNDS}
                                    valueUnit={KG}
                                    valuePrecision={8}
                                    inputPrecision={1}
                                    valueToObject={weightValueToObject}
                                    objectToValue={weightObjectToValue}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={weightExpectedRange}
                                />
                            ) : (
                                <Field
                                    name="weight"
                                    label="Weight"
                                    placeholder="Weight"
                                    component={SimpleUnitInput}
                                    unit={KG}
                                    precision={1}
                                    min={0}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={weightExpectedRange}
                                />
                            )}
                        </div>
                        <button className="btn btn-link remove" onClick={() => this.props.change("has_weight", false)}>
                            <NegativeIcon />
                            Remove
                        </button>
                    </div>
                )}
                {!this.props.opened && (
                    <button
                        className="btn btn-link"
                        onClick={() => {
                            this.props.change("has_weight", true)
                            this.props.change("focus", "weight")
                        }}
                    >
                        Add Weight
                    </button>
                )}
            </div>
        )
    }
}

Weight = connect((state, props) => {
    return {
        opened: selector(state, "has_weight"),
        focused: selector(state, "focus") === "weight",
        unit: state.config[WEIGHT_UNIT]
    }
}, {})(Weight)

class Temperature extends React.Component {
    render() {
        return (
            <div
                className={classnames("section", {
                    open: this.props.opened
                })}
            >
                {this.props.opened && (
                    <div className="form-row">
                        <div className="col-sm">
                            {this.props.unit === POUNDS_OUNCES ? (
                                <Field
                                    name="temperature"
                                    label="Body temperature"
                                    placeholder="Body temperature"
                                    component={UnitInputWithConversion}
                                    inputUnit={FAHRENHEIT}
                                    valueUnit={CELSIUS}
                                    valuePrecision={8}
                                    inputPrecision={1}
                                    valueToObject={temperatureValueToObject}
                                    objectToValue={temperatureObjectToValue}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={temperatureExpectedRange}
                                />
                            ) : (
                                <Field
                                    name="temperature"
                                    label="Body temperature"
                                    placeholder="Body temperature"
                                    component={SimpleUnitInput}
                                    unit={CELSIUS}
                                    precision={1}
                                    onKeyPress={this.props.dispatch(submitOnEnter)}
                                    autoFocus={this.props.focused}
                                    validate={required}
                                    warn={temperatureExpectedRange}
                                />
                            )}
                        </div>
                        <button className="btn btn-link remove" onClick={() => this.props.change("has_temperature", false)}>
                            <NegativeIcon />
                            Remove
                        </button>
                    </div>
                )}
                {!this.props.opened && (
                    <button
                        className="btn btn-link"
                        onClick={() => {
                            this.props.change("has_temperature", true)
                            this.props.change("focus", "temperature")
                        }}
                    >
                        Add Body temperature
                    </button>
                )}
            </div>
        )
    }
}

Temperature = connect((state, props) => {
    return {
        opened: selector(state, "has_temperature"),
        focused: selector(state, "focus") === "temperature",
        unit: state.config[TEMPERATURE_UNIT]
    }
}, {})(Temperature)

class BloodPressure extends React.Component {
    render() {
        return (
            <div
                className={classnames("section", {
                    open: this.props.opened
                })}
            >
                {this.props.opened && (
                    <div>
                        <div className="form-row title">
                            <h4>Blood pressure</h4>
                            <div className="col-sm">
                                {this.props.unit === CM_HG ? (
                                    <Field
                                        name="pressure.systolic"
                                        label="Systolic"
                                        placeholder="Systolic"
                                        component={UnitInputWithConversion}
                                        inputUnit={CM_HG}
                                        valueUnit={MM_HG}
                                        valuePrecision={0}
                                        inputPrecision={1}
                                        valueToObject={pressureValueToObject}
                                        objectToValue={pressureObjectToValue}
                                        onKeyPress={this.props.dispatch(submitOnEnter)}
                                        autoFocus={this.props.focused}
                                        validate={required}
                                        warn={systolicBloodPressureExpectedRange}
                                    />
                                ) : (
                                    <Field
                                        name="pressure.systolic"
                                        label="Systolic"
                                        placeholder="Systolic"
                                        component={SimpleUnitInput}
                                        unit={MM_HG}
                                        precision={0}
                                        onKeyPress={this.props.dispatch(submitOnEnter)}
                                        autoFocus={this.props.focused}
                                        validate={required}
                                        warn={systolicBloodPressureExpectedRange}
                                    />
                                )}
                            </div>
                            <div className="col-sm">
                                {this.props.unit === CM_HG ? (
                                    <Field
                                        name="pressure.diatolic"
                                        label="Diastolic"
                                        placeholder="Diastolic"
                                        component={UnitInputWithConversion}
                                        inputUnit={CM_HG}
                                        valueUnit={MM_HG}
                                        valuePrecision={0}
                                        inputPrecision={1}
                                        valueToObject={pressureValueToObject}
                                        objectToValue={pressureObjectToValue}
                                        onKeyPress={this.props.dispatch(submitOnEnter)}
                                        autoFocus={this.props.focused}
                                        validate={required}
                                        warn={diastolicBloodPressureExpectedRange}
                                    />
                                ) : (
                                    <Field
                                        name="pressure.diastolic"
                                        label="Diastolic"
                                        placeholder="Diastolic"
                                        component={SimpleUnitInput}
                                        unit={MM_HG}
                                        precision={0}
                                        onKeyPress={this.props.dispatch(submitOnEnter)}
                                        autoFocus={this.props.focused}
                                        validate={required}
                                        warn={diastolicBloodPressureExpectedRange}
                                    />
                                )}
                            </div>
                            <button className="btn btn-link remove" onClick={() => this.props.change("has_pressure", false)}>
                                <NegativeIcon />
                                Remove
                            </button>
                        </div>
                    </div>
                )}
                {!this.props.opened && (
                    <button
                        className="btn btn-link"
                        onClick={() => {
                            this.props.change("has_pressure", true)
                            this.props.change("focus", "pressure")
                        }}
                    >
                        Add Blood pressure
                    </button>
                )}
            </div>
        )
    }
}

BloodPressure = connect((state, props) => {
    return {
        opened: selector(state, "has_pressure"),
        focused: selector(state, "focus") === "pressure",
        unit: state.config[BLOOD_PRESSURE_UNIT]
    }
}, {})(BloodPressure)

export default MedicalData
