import React, { Component } from "react"
import _ from "lodash"
import { connect } from "react-redux"
import classnames from "classnames"
import PropTypes from "prop-types"

import {
    getPrecisionUnit,
    convertPrecision,
    getMinimum,
    getMaximum,
    KG_GRAMS,
    KG,
    GRAMS,
    POUNDS,
    POUNDS_OUNCES,
    CM,
    CELSIUS,
    MM_HG,
    WEIGHT_UNIT,
    LENGTH_UNIT,
    TEMPERATURE_UNIT,
    BLOOD_PRESSURE_UNIT
} from "../unitConversion/units"
import { round, isStringNumber, getNumberFromString, getPrecision } from "../utils"
import {
    weightValueToObject,
    weightObjectToValue,
    lengthValueToObject,
    lengthObjectToValue,
    pressureValueToObject,
    pressureObjectToValue,
    temperatureValueToObject,
    temperatureObjectToValue
} from "../unitConversion"

class SimpleUnitInput extends Component {
    constructor(props) {
        super(props)
        let precision = this.props.precision ? this.props.precision : 0
        let inputValue = isStringNumber(String(props.input.value)) ? round(getNumberFromString(String(props.input.value)), precision).toString() : ""

        this.state = {
            inputValue: inputValue,
            noInput: inputValue === "",
            changed: props.meta.touched,
            warning: props.meta.warning
        }
    }

    componentDidUpdate(prevProps) {
        let precision = this.props.precision ? this.props.precision : 0

        if (this.props.precision !== prevProps.precision || this.props.input.value !== prevProps.input.value) {
            let currentValue = this.state.inputValue !== "" ? round(getNumberFromString(this.state.inputValue), precision) : 0
            let value = isStringNumber(this.props.input.value) ? round(getNumberFromString(this.props.input.value), precision) : 0

            if (currentValue !== value) {
                this.setState({
                    inputValue: value.toString()
                })
            }
        }

        if (this.props.meta.warning !== prevProps.meta.warning) {
            this.state.warningTimeout && clearTimeout(this.state.warningTimeout)
            if (this.props.meta.warning === undefined || !this.state.changed) {
                this.setState({
                    warning: this.props.meta.warning
                })
            } else {
                this.setState({
                    warningTimeout: window.setTimeout(this.setWarning(this.props.meta.warning), 600)
                })
            }
        }
    }

    setWarning(warning) {
        return () => this.setState({ warning: warning })
    }

    setCurrentWarning() {
        this.setState({ warning: this.props.meta.warning })
    }

    parseInputChange() {
        return event => {
            this.setState({ changed: true })

            let { precision, min, max } = this.props
            precision = precision ? precision : 0

            let inputValueAsNumber = getNumberFromString(this.state.inputValue)
            let v = isStringNumber(event.target.value) ? getNumberFromString(event.target.value) : false

            if (
                (event.target.value === "" && this.state.inputValue !== "") ||
                (v !== false && getPrecision(v) <= precision && (min === undefined || v >= min) && (max === undefined || v <= max))
            ) {
                this.setState({
                    inputValue: precision === 0 ? event.target.value.replace(",", "").replace(".", "") : event.target.value,
                    noInput: event.target.value === ""
                })
                inputValueAsNumber = v ? v : 0
            }

            return event.target.value !== "" ? round(inputValueAsNumber, precision).toString() : ""
        }
    }

    render() {
        const { input, meta, label, placeholder, optional, disabled, unit, autoFocus, onKeyPress } = this.props

        return (
            <div className={classnames("form-group", { "is-invalid": this.state.changed && meta.error, "is-warning": !meta.error && this.state.warning })}>
                {!this.state.noInput && <span className="label">{label}</span>}
                <div className="inputWithUnit-1">
                    <input
                        {...input}
                        value={this.state.inputValue}
                        disabled={disabled}
                        onBlur={event => {
                            this.setCurrentWarning()
                            input.onBlur(this.parseInputChange()(event))
                        }}
                        onChange={event => input.onChange(this.parseInputChange()(event))}
                        autoFocus={autoFocus}
                        onKeyPress={onKeyPress && (event => onKeyPress(event))}
                        className={classnames("form-control", {
                            "is-invalid": this.state.changed && meta.error,
                            "is-warning": !meta.error && this.state.warning
                        })}
                        placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
                        type="text"
                    />
                    <span className="unit">{unit}</span>
                </div>
                {this.state.changed && meta.error && <div className="invalid-feedback">{meta.error}</div>}
                {!meta.error && this.state.warning && <div className="warning-feedback">{this.state.warning}</div>}
            </div>
        )
    }
}

SimpleUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    unit: PropTypes.string.isRequired,
    label: PropTypes.string.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool,
    precision: PropTypes.number,
    min: PropTypes.number,
    max: PropTypes.number
}

class UnitInputWithConversion extends Component {
    constructor(props) {
        super(props)
        let valuePrecision = props.valuePrecision ? props.valuePrecision : 0

        let inputValue = isStringNumber(String(props.input.value)) ? round(getNumberFromString(String(props.input.value)), valuePrecision).toString() : ""
        let object = this.convertValueToObject(inputValue)
        let noInput = true
        let inputValues = _.reduce(
            object,
            (result, value, unit) => {
                result[unit] = value > 0 ? value.toString() : ""
                noInput = result[unit] !== "" ? false : noInput
                return result
            },
            {}
        )

        this.state = {
            inputValues: inputValues,
            noInput: noInput,
            changed: props.meta.touched,
            warning: props.meta.warning
        }
    }

    componentDidUpdate(prevProps) {
        if (
            this.props.input.value !== prevProps.input.value ||
            this.props.inputUnit !== prevProps.inputUnit ||
            this.props.valuePrecision !== prevProps.valuePrecision ||
            this.props.inputPrecision !== prevProps.inputPrecision
        ) {
            let inputPrecision = this.props.precision ? this.props.precision : 0
            let object = this.convertValueToObject(this.props.input.value)
            let noInput = this.state.noInput

            let inputValues = _.reduce(
                object,
                (result, value, unit) => {
                    let precision = unit === getPrecisionUnit(this.props.inputUnit) ? inputPrecision : 0
                    let currentValue = this.state.inputValues[unit] !== "" ? round(getNumberFromString(this.state.inputValues[unit]), precision) : 0
                    result[unit] = currentValue !== value ? value.toString() : this.state.inputValues[unit]
                    noInput = result[unit] !== "" ? false : noInput
                    return result
                },
                {}
            )

            this.setState({
                inputValues: inputValues,
                noInput: noInput
            })
        }

        if (this.props.meta.warning !== prevProps.meta.warning) {
            this.state.warningTimeout && clearTimeout(this.state.warningTimeout)
            if (this.props.meta.warning === undefined || !this.state.changed) {
                this.setState({
                    warning: this.props.meta.warning
                })
            } else {
                this.setState({
                    warningTimeout: window.setTimeout(this.setWarning(this.props.meta.warning), 600)
                })
            }
        }
    }

    convertValueToObject(value) {
        let { valueUnit, inputUnit, inputPrecision, valueToObject } = this.props
        inputPrecision = inputPrecision ? inputPrecision : 0

        return valueToObject(valueUnit, inputUnit, inputPrecision)(getNumberFromString(value))
    }

    setWarning(warning) {
        return () => this.setState({ warning: warning })
    }

    setCurrentWarning() {
        this.setState({ warning: this.props.meta.warning })
    }

    parseInputChange(inputChangeUnit) {
        return event => {
            this.setState({ changed: true })

            let { inputUnit, valueUnit, inputPrecision, valuePrecision, objectToValue } = this.props
            valuePrecision = valuePrecision ? valuePrecision : 0
            inputPrecision = getPrecisionUnit(inputUnit) === inputChangeUnit && inputPrecision ? inputPrecision : 0
            let min = getMinimum(inputUnit)
            let max = getMaximum(inputUnit)

            let inputValues = this.state.inputValues
            let noInput = this.state.noInput
            let v = isStringNumber(event.target.value) ? getNumberFromString(event.target.value) : false

            if (
                (event.target.value === "" && inputValues[inputChangeUnit] !== "") ||
                (v !== false &&
                    getPrecision(v) <= inputPrecision &&
                    (min[inputChangeUnit] === undefined || v >= min[inputChangeUnit]) &&
                    (max[inputChangeUnit] === undefined || v <= max[inputChangeUnit]))
            ) {
                inputValues[inputChangeUnit] = event.target.value
                noInput = _.reduce(
                    inputValues,
                    (result, value) => {
                        if (value !== "") {
                            result = false
                        }
                        return result
                    },
                    true
                )

                this.setState({
                    inputValues: inputValues,
                    noInput: noInput
                })
            }

            return noInput
                ? ""
                : objectToValue(valueUnit, valuePrecision)(
                      _.reduce(
                          inputValues,
                          (object, value, unit) => {
                              object[unit] = isStringNumber(value) ? getNumberFromString(value) : 0
                              return object
                          },
                          {}
                      )
                  ).toString()
        }
    }

    render() {
        const { input, meta, label, placeholder, optional, disabled, autoFocus, onKeyPress } = this.props
        let i = 0
        return (
            <div
                className={classnames("form-group", {
                    "is-invalid": this.state.changed && meta.error,
                    "is-warning": !meta.error && this.state.warning
                })}
            >
                {!this.state.noInput && <span className="label">{label}</span>}
                {_.map(this.state.inputValues, (value, unit) => {
                    return (
                        <div className={`inputWithUnit-${_.size(this.state.inputValues)}`} key={unit}>
                            <input
                                {...input}
                                value={value}
                                autoFocus={autoFocus && 0 === i++}
                                disabled={disabled}
                                onBlur={event => {
                                    this.setCurrentWarning()
                                    input.onBlur(this.parseInputChange(unit)(event))
                                }}
                                onChange={event => input.onChange(this.parseInputChange(unit)(event))}
                                onKeyPress={onKeyPress && (event => onKeyPress(event))}
                                className={classnames("form-control", {
                                    "is-invalid": this.state.changed && meta.error,
                                    "is-warning": !meta.error && this.state.warning
                                })}
                                placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
                                type="text"
                            />
                            <span className="unit">{unit}</span>
                        </div>
                    )
                })}
                {this.state.changed && meta.error && <div className="invalid-feedback">{meta.error}</div>}
                {!meta.error && this.state.warning && <div className="warning-feedback">{this.state.warning}</div>}
            </div>
        )
    }
}

UnitInputWithConversion.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    label: PropTypes.string.isRequired,
    inputUnit: PropTypes.string.isRequired,
    valueUnit: PropTypes.string.isRequired,
    objectToValue: PropTypes.func.isRequired,
    valueToObject: PropTypes.func.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    inputPrecision: PropTypes.number,
    valuePrecision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

class HeightUnitInput extends Component {
    render() {
        if (this.props.inputUnit === this.props.valueUnit) {
            return <SimpleUnitInput {...this.props} unit={this.props.valueUnit} min={0} />
        }

        return (
            <UnitInputWithConversion
                {...this.props}
                valuePrecision={8}
                inputPrecision={convertPrecision(this.props.valueUnit, this.props.inputUnit, this.props.precision)}
                valueToObject={lengthValueToObject}
                objectToValue={lengthObjectToValue}
            />
        )
    }
}

HeightUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    label: PropTypes.string.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    unit: PropTypes.string,
    precision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

HeightUnitInput = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, LENGTH_UNIT, CM)
    let valueUnit = props.unit ? props.unit : CM

    return {
        inputUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(HeightUnitInput)

class WeightUnitInput extends Component {
    render() {
        if (this.props.inputUnit === this.props.valueUnit) {
            return <SimpleUnitInput {...this.props} unit={this.props.valueUnit} min={0} />
        }

        return (
            <UnitInputWithConversion
                {...this.props}
                valuePrecision={8}
                inputPrecision={convertPrecision(this.props.valueUnit, this.props.inputUnit, this.props.precision)}
                valueToObject={weightValueToObject}
                objectToValue={weightObjectToValue}
            />
        )
    }
}

WeightUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    label: PropTypes.string.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    unit: PropTypes.string,
    precision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

WeightUnitInput = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, WEIGHT_UNIT, KG_GRAMS)
    let valueUnit = props.unit ? props.unit : KG

    let inputUnit
    switch (valueUnit) {
        case GRAMS:
            switch (configuredInputUnit) {
                case POUNDS_OUNCES:
                    inputUnit = POUNDS_OUNCES
                    break
                default:
                    inputUnit = GRAMS
                    break
            }
            break
        default:
            switch (configuredInputUnit) {
                case POUNDS_OUNCES:
                    inputUnit = POUNDS
                    break
                default:
                    inputUnit = KG
                    break
            }
    }

    return {
        inputUnit: inputUnit,
        valueUnit: valueUnit
    }
}, {})(WeightUnitInput)

class TemperatureUnitInput extends Component {
    render() {
        if (this.props.inputUnit === this.props.valueUnit) {
            return <SimpleUnitInput {...this.props} unit={this.props.valueUnit} min={0} />
        }

        return (
            <UnitInputWithConversion
                {...this.props}
                valuePrecision={8}
                inputPrecision={convertPrecision(this.props.valueUnit, this.props.inputUnit, this.props.precision)}
                valueToObject={temperatureValueToObject}
                objectToValue={temperatureObjectToValue}
            />
        )
    }
}

TemperatureUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    label: PropTypes.string.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    unit: PropTypes.string,
    precision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

TemperatureUnitInput = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, TEMPERATURE_UNIT, CELSIUS)
    let valueUnit = props.unit ? props.unit : CELSIUS

    return {
        inputUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(TemperatureUnitInput)

class BloodPressureUnitInput extends Component {
    render() {
        if (this.props.inputUnit === this.props.valueUnit) {
            return <SimpleUnitInput {...this.props} unit={this.props.valueUnit} min={0} />
        }

        return (
            <UnitInputWithConversion
                {...this.props}
                valuePrecision={8}
                inputPrecision={convertPrecision(this.props.valueUnit, this.props.inputUnit, this.props.precision)}
                valueToObject={pressureValueToObject}
                objectToValue={pressureObjectToValue}
            />
        )
    }
}

BloodPressureUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    label: PropTypes.string.isRequired,
    onKeyPress: PropTypes.func,
    autoFocus: PropTypes.bool,
    unit: PropTypes.string,
    precision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

BloodPressureUnitInput = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, BLOOD_PRESSURE_UNIT, MM_HG)
    let valueUnit = props.unit ? props.unit : MM_HG

    return {
        inputUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(BloodPressureUnitInput)

export { SimpleUnitInput, UnitInputWithConversion, HeightUnitInput, WeightUnitInput, TemperatureUnitInput, BloodPressureUnitInput }
