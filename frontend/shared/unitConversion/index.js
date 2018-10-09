import convert from "convert-units"
import _ from "lodash"

import { round } from "../utils"
import * as units from "./units"

const CONVERT_UNITS_SYMBOLS = {
    [units.KG]: "kg",
    [units.GRAMS]: "g",
    [units.POUNDS]: "lb",
    [units.OUNCES]: "oz",
    [units.CM]: "cm",
    [units.METRES]: "m",
    [units.FEET]: "ft",
    [units.INCHES]: "in",
    [units.MM_HG]: "mm",
    [units.CM_HG]: "cm",
    [units.CELSIUS]: "C",
    [units.FAHRENHEIT]: "F"
}

export const getObjectFromValue = (valueUnit, objectUnit, precision) => value => {
    // special case for pressure with two values (diastolic and systolioc together) passed needs to be processed first
    if (_.includes(units.PRESSURE_VALUE_UNITS, valueUnit) && value instanceof Array) {
        let result = ""
        _.map(value, (val, key) => {
            if (key === 0) {
                result = pressureValueToObject(valueUnit, objectUnit, precision)(val)[objectUnit]
            } else {
                result = result + "/" + pressureValueToObject(valueUnit, objectUnit, precision)(val)[objectUnit]
            }
        })

        return { [objectUnit]: result }
    }

    if (valueUnit === objectUnit) {
        return { [objectUnit]: round(parseFloat(value), precision) }
    }

    switch (objectUnit) {
        case units.KG:
        case units.POUNDS:
        case units.OUNCES:
        case units.GRAMS:
        case units.POUNDS_OUNCES:
        case units.KG_GRAMS:
            return weightValueToObject(valueUnit, objectUnit, precision)(value)
        case units.CM:
        case units.METRES:
        case units.FEET:
        case units.INCHES:
        case units.FEET_INCHES:
        case units.METRES_CM:
            return lengthValueToObject(valueUnit, objectUnit, precision)(value)
        case units.MM_HG:
        case units.CM_HG:
            return pressureValueToObject(valueUnit, objectUnit, precision)(value)
        case units.CELSIUS:
        case units.FAHRENHEIT:
            return temperatureValueToObject(valueUnit, objectUnit, precision)(value)
        default:
            console.log("Unsupported object unit. Returns 0.")
            return ""
    }
}

export const getValueFromObject = (valueUnit, objectUnit, precision) => object => {
    if (valueUnit === objectUnit) {
        return round(object[objectUnit], precision)
    }

    switch (objectUnit) {
        case units.KG:
        case units.POUNDS:
        case units.OUNCES:
        case units.GRAMS:
        case units.POUNDS_OUNCES:
        case units.KG_GRAMS:
            return weightObjectToValue(valueUnit, precision)(object)
        case units.CM:
        case units.METRES:
        case units.FEET:
        case units.INCHES:
        case units.FEET_INCHES:
        case units.METRES_CM:
            return lengthObjectToValue(valueUnit, precision)(object)
        case units.MM_HG:
        case units.CM_HG:
            return pressureObjectToValue(valueUnit, precision)(object)
        case units.CELSIUS:
        case units.FAHRENHEIT:
            return temperatureObjectToValue(valueUnit, precision)(object)
        default:
            console.log("Unsupported object unit. Returns 0.")
            return ""
    }
}

const objectToValue = (valueUnits, simpleUnits) => (valueUnit, precision) => object => {
    if (!_.includes(valueUnits, valueUnit)) {
        console.log("Unsupported value unit `" + valueUnit + "`. Returns 0.")
        return 0
    }

    try {
        return round(
            _.reduce(
                object,
                (result, value, unit) => {
                    if (!_.includes(simpleUnits, unit)) {
                        throw new Error("Invalid unit in the object. Returns 0.")
                    }
                    result =
                        result +
                        convert(value)
                            .from(CONVERT_UNITS_SYMBOLS[unit])
                            .to(CONVERT_UNITS_SYMBOLS[valueUnit])
                    return result
                },
                0
            ),
            precision
        )
    } catch (err) {
        console.log(err.message)
        return 0
    }
}

const valueToObject = (valueUnits, simpleUnits, complexUnits, complexUnitsLimits) => (valueUnit, objectUnit, precision) => value => {
    if (!_.includes(valueUnits, valueUnit)) {
        console.log("Unsupported value unit `" + valueUnit + "`. Returns empty object.")
        return {}
    }

    if (_.includes(simpleUnits, objectUnit)) {
        if (isNaN(value)) {
            return { [objectUnit]: 0 }
        }
        return {
            [objectUnit]: round(
                convert(value)
                    .from(CONVERT_UNITS_SYMBOLS[valueUnit])
                    .to(CONVERT_UNITS_SYMBOLS[objectUnit]),
                precision
            )
        }
    } else if (_.includes(complexUnits, objectUnit)) {
        try {
            let units = objectUnit.split("_")
            let prevUnit
            return _.reduce(
                units,
                (result, unit, index) => {
                    if (!_.includes(simpleUnits, unit)) {
                        throw new Error("Invalid part unit `" + unit + "` in object unit`" + objectUnit + "`. Returns empty object.")
                    }
                    if (isNaN(value)) {
                        result[unit] = 0
                        return result
                    }

                    if (index < units.length - 1) {
                        result[unit] = Math.floor(
                            convert(value)
                                .from(CONVERT_UNITS_SYMBOLS[valueUnit])
                                .to(CONVERT_UNITS_SYMBOLS[unit])
                        )
                    } else {
                        let v = round(
                            convert(value)
                                .from(CONVERT_UNITS_SYMBOLS[valueUnit])
                                .to(CONVERT_UNITS_SYMBOLS[unit]) - objectToValue(valueUnits, simpleUnits)(unit)(result),
                            precision
                        )

                        if (v >= _.get(complexUnitsLimits, `${objectUnit}.${unit}`, v + 1)) {
                            result[prevUnit]++

                            v = round(
                                convert(value)
                                    .from(CONVERT_UNITS_SYMBOLS[valueUnit])
                                    .to(CONVERT_UNITS_SYMBOLS[unit]) - objectToValue(valueUnits, simpleUnits)(unit)(result),
                                precision
                            )
                        }
                        result[unit] = v
                    }
                    prevUnit = unit

                    return result
                },
                {}
            )
        } catch (err) {
            console.log(err.message)
            return {}
        }
    }

    console.log("Unsupported object unit `" + objectUnit + "`. Returns empty object.")
    return {}
}

export const weightObjectToValue = (valueUnit, precision) => object => {
    return objectToValue(units.WEIGHT_VALUE_UNITS, units.WEIGHT_SIMPLE_UNITS)(valueUnit, precision)(object)
}

export const weightValueToObject = (valueUnit, objectUnit, precision) => value => {
    return valueToObject(units.WEIGHT_VALUE_UNITS, units.WEIGHT_SIMPLE_UNITS, units.WEIGHT_COMPLEX_UNITS, units.WEIGHT_COMPLEX_UNITS_LIMITS)(
        valueUnit,
        objectUnit,
        precision
    )(value)
}

export const lengthObjectToValue = (valueUnit, precision) => object => {
    return objectToValue(units.LENGTH_VALUE_UNITS, units.LENGTH_SIMPLE_UNITS)(valueUnit, precision)(object)
}

export const lengthValueToObject = (valueUnit, objectUnit, precision) => value => {
    return valueToObject(units.LENGTH_VALUE_UNITS, units.LENGTH_SIMPLE_UNITS, units.LENGTH_COMPLEX_UNITS, units.LENGTH_COMPLEX_UNITS_LIMITS)(
        valueUnit,
        objectUnit,
        precision
    )(value)
}

export const pressureObjectToValue = (valueUnit, precision) => object => {
    return objectToValue(units.PRESSURE_VALUE_UNITS, units.PRESSURE_SIMPLE_UNITS)(valueUnit, precision)(object)
}

export const pressureValueToObject = (valueUnit, objectUnit, precision) => value => {
    return valueToObject(units.PRESSURE_VALUE_UNITS, units.PRESSURE_SIMPLE_UNITS, [], {})(valueUnit, objectUnit, precision)(value)
}

export const temperatureObjectToValue = (valueUnit, precision) => object => {
    return objectToValue(units.TEMPERATURE_VALUE_UNITS, units.TEMPERATURE_SIMPLE_UNITS)(valueUnit, precision)(object)
}

export const temperatureValueToObject = (valueUnit, objectUnit, precision) => value => {
    return valueToObject(units.TEMPERATURE_VALUE_UNITS, units.TEMPERATURE_SIMPLE_UNITS, [], {})(valueUnit, objectUnit, precision)(value)
}
