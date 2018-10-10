import React, { Component } from "react"
import _ from "lodash"
import classnames from "classnames"
import PropTypes from "prop-types"
import { getPrecisionUnit, getMinimum, getMaximum } from "../unitConversion/units"
import { round, isStringNumber, getNumberFromString, getPrecision } from "../utils"

class SimpleUnitInput extends Component {
    constructor(props) {
        super(props)
        let precision = this.props.precision ? this.props.precision : 0
        let inputValue = isStringNumber(props.input.value) ? round(getNumberFromString(props.input.value), precision).toString() : ""

        this.state = {
            inputValue: inputValue,
            noInput: inputValue === ""
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
    }

    parseInputChange() {
        return event => {
            let { precision, min, max } = this.props
            precision = precision ? precision : 0

            let inputValueAsNumber = getNumberFromString(this.state.inputValue)
            let v = isStringNumber(event.target.value) ? getNumberFromString(event.target.value) : false

            if (
                (event.target.value === "" && this.state.inputValue !== "") ||
                (v !== false && getPrecision(v) <= precision && (min === undefined || v >= min) && (max === undefined || v <= max))
            ) {
                this.setState({
                    inputValue: event.target.value,
                    noInput: event.target.value === ""
                })
                inputValueAsNumber = v ? v : 0
            }

            return event.target.value !== "" ? round(inputValueAsNumber, precision).toString() : ""
        }
    }

    render() {
        const { input, meta, label, placeholder, optional, disabled, unit, autoFocus, onKeyPress } = this.props
        if (label === "Weight") {
            console.log(meta)
        }
        return (
            <div className={classnames("form-group", { "is-invalid": meta.touched && meta.error, "is-warning": !meta.error && meta.warning })}>
                {!this.state.noInput && <span className="label">{label}</span>}
                <div className="inputWithUnit-1">
                    <input
                        {...input}
                        value={this.state.inputValue}
                        disabled={disabled}
                        onBlur={event => input.onBlur(this.parseInputChange(input.value)(event))}
                        onChange={event => input.onChange(this.parseInputChange(input.value)(event))}
                        autoFocus={autoFocus}
                        onKeyPress={onKeyPress && (event => onKeyPress(event))}
                        className={classnames("form-control", {
                            "is-invalid": meta.touched && meta.error,
                            "is-warning": !meta.error && meta.warning
                        })}
                        placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
                        type="text"
                    />
                    <span className="unit">{unit}</span>
                </div>
                {meta.touched && meta.error && <div className="invalid-feedback">{meta.error}</div>}
                {!meta.error && meta.warning && <div className="warning-feedback">{meta.warning}</div>}
            </div>
        )
    }
}

SimpleUnitInput.propTypes = {
    input: PropTypes.object.isRequired,
    meta: PropTypes.object.isRequired,
    unit: PropTypes.string.isRequired,
    label: PropTypes.string.isRequired,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool,
    precision: PropTypes.number,
    min: PropTypes.number,
    max: PropTypes.number
}

export { SimpleUnitInput }

class UnitInputWithConversion extends Component {
    constructor(props) {
        super(props)

        let object = this.convertValueToObject(props.input.value)
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
            noInput: noInput
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
    }

    convertValueToObject(value) {
        let { valueUnit, inputUnit, inputPrecision, valueToObject } = this.props
        inputPrecision = inputPrecision ? inputPrecision : 0

        return valueToObject(valueUnit, inputUnit, inputPrecision)(getNumberFromString(value))
    }

    parseInputChange(inputChangeUnit) {
        return event => {
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
            <div className={classnames("form-group", { "is-invalid": meta.touched && meta.error, "is-warning": meta.touched && !meta.error && meta.warning })}>
                {!this.state.noInput && <span className="label">{label}</span>}
                {_.map(this.state.inputValues, (value, unit) => {
                    return (
                        <div className={`inputWithUnit-${_.size(this.state.inputValues)}`} key={unit}>
                            <input
                                {...input}
                                value={value}
                                autoFocus={autoFocus && 0 === i++}
                                disabled={disabled}
                                onBlur={event => input.onBlur(this.parseInputChange(unit)(event))}
                                onChange={event => input.onChange(this.parseInputChange(unit)(event))}
                                onKeyPress={onKeyPress && (event => onKeyPress(event))}
                                className={classnames("form-control", {
                                    "is-invalid": meta.touched && meta.error,
                                    "is-warning": meta.touched && !meta.error && meta.warning
                                })}
                                placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
                                type="text"
                            />
                            <span className="unit">{unit}</span>
                        </div>
                    )
                })}
                {meta.touched && meta.error && <div className="invalid-feedback">{meta.error}</div>}
                {!meta.error && meta.warning && <div className="warning-feedback">{meta.warning}</div>}
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
    inputPrecision: PropTypes.number,
    valuePrecision: PropTypes.number,
    placeholder: PropTypes.string,
    optional: PropTypes.bool,
    disabled: PropTypes.bool
}

export { UnitInputWithConversion }
