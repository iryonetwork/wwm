import _ from "lodash"
import React from "react"
import classnames from "classnames"
import { connect } from "react-redux"

import { getObjectFromValue } from "../../unitConversion"
import {
    getLabel,
    hasLabelWhitespace,
    canLabelBeSmall,
    convertPrecision,
    FORMATTING_LONG,
    FORMATTING_SHORT,
    FORMATTING_DEFAULT,
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
} from "../../unitConversion/units"

export const LONG = FORMATTING_LONG
export const SHORT = FORMATTING_SHORT
export const DEFAULT = FORMATTING_DEFAULT

export const ValueWithUnit = ({ valueUnit, displayUnit, precision, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT

    if (!hasLabelWhitespace(formatting)(displayUnit)) {
        return (
            <span className="valueWithUnit">
                {_.map(getObjectFromValue(valueUnit, displayUnit, convertPrecision(valueUnit, displayUnit, precision ? precision : 0))(value), (val, unit) => (
                    <SimpleValueWithUnit key={unit} unit={unit} formatting={formatting} value={val} />
                ))}
            </span>
        )
    }
    return _.map(getObjectFromValue(valueUnit, displayUnit, convertPrecision(valueUnit, displayUnit, precision ? precision : 0))(value), (val, unit) => (
        <span className="valueWithUnit" key={unit}>
            <SimpleValueWithUnit unit={unit} formatting={formatting} value={val} />
        </span>
    ))
}

const SimpleValueWithUnit = ({ unit, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT

    return (
        <React.Fragment>
            <React.Fragment>
                <span className={classnames("value", { withMargin: hasLabelWhitespace(formatting)(unit) })}>{value}</span>
                <span className={classnames("unit", { big: !canLabelBeSmall(formatting)(unit) })}>{getLabel(formatting)(unit, value)}</span>
            </React.Fragment>
        </React.Fragment>
    )
}

export const ValueWithUnitString = ({ valueUnit, stringUnit, precision, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT
    let object = getObjectFromValue(valueUnit, stringUnit, convertPrecision(valueUnit, stringUnit, precision ? precision : 0))(value)

    let i = 0
    return _.reduce(
        object,
        (result, value, unit) => {
            if (i !== 0 && hasLabelWhitespace(formatting)(stringUnit)) {
                result = result + " "
            }
            result = result + value.toString()
            if (hasLabelWhitespace(formatting)(unit)) {
                result = result + " "
            }

            i++

            return result + getLabel(formatting)(unit, value)
        },
        ""
    )
}

export const HeightValueWithUnit = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, LENGTH_UNIT, CM)
    let valueUnit = props.unit ? props.unit : CM

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnit)

export const HeightValueWithUnitString = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, LENGTH_UNIT, CM)
    let valueUnit = props.unit ? props.unit : CM

    return {
        stringUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)

export const WeightValueWithUnit = connect((state, props) => {
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
        displayUnit: inputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnit)

export const WeightValueWithUnitString = connect((state, props) => {
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
        stringUnit: inputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)

export const TemperatureValueWithUnitString = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, TEMPERATURE_UNIT, CELSIUS)
    let valueUnit = props.unit ? props.unit : CELSIUS

    return {
        stringUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)

export const TemperatureValueWithUnit = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, TEMPERATURE_UNIT, CELSIUS)
    let valueUnit = props.unit ? props.unit : CELSIUS

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)

export const BloodPressureValueWithUnitString = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, BLOOD_PRESSURE_UNIT, MM_HG)
    let valueUnit = props.unit ? props.unit : MM_HG

    return {
        stringUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)

export const BloodPressureValueWithUnit = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, BLOOD_PRESSURE_UNIT, MM_HG)
    let valueUnit = props.unit ? props.unit : MM_HG

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(ValueWithUnitString)
