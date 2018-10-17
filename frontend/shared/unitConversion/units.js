// weight units
export const LENGTH_UNIT = "lengthUnit"
export const WEIGHT_UNIT = "weightUnit"
export const TEMPERATURE_UNIT = "temperatureUnit"
export const BLOOD_PRESSURE_UNIT = "bloodPressureUnit"

export const KG = "kg"
export const GRAMS = "g"
export const POUNDS = "lb"
export const OUNCES = "oz"
export const KG_GRAMS = "kg_g"
export const POUNDS_OUNCES = "lb_oz"
export const WEIGHT_VALUE_UNITS = [KG, GRAMS, POUNDS, OUNCES]
export const WEIGHT_SIMPLE_UNITS = [KG, GRAMS, POUNDS, OUNCES]
export const WEIGHT_COMPLEX_UNITS = [KG_GRAMS, POUNDS_OUNCES]
export const WEIGHT_COMPLEX_UNITS_LIMITS = {
    [POUNDS_OUNCES]: {
        [OUNCES]: 16
    },
    [KG_GRAMS]: {
        [GRAMS]: 1000
    }
}

// length units
export const CM = "cm"
export const METERS = "m"
export const FEET = "ft"
export const INCHES = "in"
export const METERS_CM = "m_cm"
export const FEET_INCHES = "ft_in"
export const LENGTH_VALUE_UNITS = [CM, METERS, FEET, INCHES]
export const LENGTH_SIMPLE_UNITS = [CM, METERS, FEET, INCHES]
export const LENGTH_COMPLEX_UNITS = [METERS_CM, FEET_INCHES]
export const LENGTH_COMPLEX_UNITS_LIMITS = {
    [METERS_CM]: {
        [CM]: 100
    },
    [FEET_INCHES]: {
        [INCHES]: 12
    }
}

// pressure units
export const MM_HG = "mm[Hg]"
export const CM_HG = "cm[Hg]"
export const PRESSURE_VALUE_UNITS = [MM_HG]
export const PRESSURE_SIMPLE_UNITS = [MM_HG, CM_HG]

// temperature units
export const CELSIUS = "°C"
export const FAHRENHEIT = "°F"
export const TEMPERATURE_VALUE_UNITS = [CELSIUS]
export const TEMPERATURE_SIMPLE_UNITS = [CELSIUS, FAHRENHEIT]

export const FORMATTING_LONG = "LONG"
export const FORMATTING_SHORT = "SHORT"
export const FORMATTING_DEFAULT = "DEFAULT"

export const FORMATTING = {
    [KG]: {
        SHORT: {
            LABEL: KG
        },
        LONG: {
            SINGULAR_LABEL: "kilogram",
            LABEL: "kilograms"
        },
        LABEL: KG,
        WHITESPACE: true
    },
    [GRAMS]: {
        SHORT: {
            LABEL: GRAMS
        },
        LONG: {
            SINGULAR_LABEL: "gram",
            LABEL: "grams"
        },
        LABEL: GRAMS
    },
    [KG_GRAMS]: { SHORT: {}, LONG: {} },
    [POUNDS]: {
        SHORT: {
            LABEL: POUNDS
        },
        LONG: {
            SINGULAR_LABEL: "pound",
            LABEL: "pounds"
        },
        LABEL: POUNDS
    },
    [OUNCES]: {
        SHORT: {
            LABEL: OUNCES
        },
        LONG: {
            SINGULAR_LABEL: "ounce",
            LABEL: "ounces"
        },
        LABEL: OUNCES
    },
    [POUNDS_OUNCES]: { SHORT: {}, LONG: {} },
    [CM]: {
        SHORT: {
            LABEL: CM
        },
        LONG: {
            SINGULAR_LABEL: "centimeter",
            LABEL: "centimeters"
        },
        LABEL: CM
    },
    [METERS]: {
        SHORT: {
            LABEL: METERS
        },
        LONG: {
            SINGULAR_LABEL: "meters",
            LABEL: "meters"
        },
        LABEL: METERS
    },
    [FEET]: {
        SHORT: {
            LABEL: "′",
            NO_WHITESPACE: true,
            CAN_BE_SMALL: false
        },
        LONG: {
            LABEL: "feet",
            SINGULAR_LABEL: "foot"
        },
        LABEL: FEET
    },
    [INCHES]: {
        SHORT: {
            LABEL: "″",
            NO_WHITESPACE: true,
            CAN_BE_SMALL: false
        },
        LONG: {
            LABEL: "inches",
            SINGULAR_LABEL: "inch"
        },
        LABEL: INCHES
    },
    [FEET_INCHES]: { SHORT: { NO_WHITESPACE: true }, LONG: {} },
    [MM_HG]: {
        SHORT: {
            LABEL: MM_HG
        },
        LONG: {
            LABEL: MM_HG
        },
        LABEL: MM_HG
    },
    [CM_HG]: {
        SHORT: {
            LABEL: CM_HG
        },
        LONG: {
            LABEL: CM_HG
        },
        LABEL: CM_HG
    },
    [CELSIUS]: {
        SHORT: {
            LABEL: CELSIUS
        },
        LONG: {
            LABEL: CELSIUS
        },
        LABEL: CELSIUS
    },
    [FAHRENHEIT]: {
        SHORT: {
            LABEL: FAHRENHEIT
        },
        LONG: {
            LABEL: FAHRENHEIT
        },
        LABEL: FAHRENHEIT
    }
}

export const MININUM = {
    [KG]: {
        [KG]: 0
    },
    [GRAMS]: {
        [GRAMS]: 0
    },
    [KG_GRAMS]: {
        [KG]: 0,
        [GRAMS]: 0
    },
    [POUNDS]: {
        [POUNDS]: 0
    },
    [OUNCES]: {
        [OUNCES]: 0
    },
    [POUNDS_OUNCES]: {
        [POUNDS]: 0,
        [OUNCES]: 0
    },
    [CM]: {
        [CM]: 0
    },
    [METERS]: {
        [METERS]: 0
    },
    [METERS_CM]: {
        [METERS]: 0,
        [CM]: 0
    },
    [FEET]: {
        [FEET]: 0
    },
    [INCHES]: {
        [INCHES]: 0
    },
    [FEET_INCHES]: {
        [FEET]: 0,
        [INCHES]: 0
    },
    [MM_HG]: {
        [MM_HG]: 0
    },
    [CM_HG]: {
        [CM_HG]: 0
    },
    [CELSIUS]: {
        [CELSIUS]: -273.15
    },
    [FAHRENHEIT]: {
        [FAHRENHEIT]: -459.67
    }
}

export const MAXIMUM = {
    [KG]: {},
    [GRAMS]: {},
    [KG_GRAMS]: {},
    [POUNDS]: {},
    [OUNCES]: {},
    [POUNDS_OUNCES]: {
        [OUNCES]: 15.99999999
    },
    [CM]: {},
    [METERS]: {},
    [METERS_CM]: {
        [CM]: 99.99999999
    },
    [FEET]: {},
    [INCHES]: {},
    [FEET_INCHES]: {
        [INCHES]: 11.99999999
    },
    [MM_HG]: {},
    [CM_HG]: {},
    [CELSIUS]: {},
    [FAHRENHEIT]: {}
}

export const getLabel = formattingType => (unit, value) => {
    if (FORMATTING[unit]) {
        switch (formattingType) {
            case FORMATTING_SHORT:
            case FORMATTING_LONG:
                if (FORMATTING[unit][formattingType] && FORMATTING[unit][formattingType].LABEL) {
                    if (FORMATTING[unit][formattingType].SINGULAR_LABEL && value <= 1 && value > 0) {
                        return FORMATTING[unit][formattingType].SINGULAR_LABEL
                    }
                    return FORMATTING[unit][formattingType].LABEL
                }
                return getLabel()(unit, value)
            default:
                if (FORMATTING[unit].LABEL) {
                    if (FORMATTING[unit].SINGULAR_LABEL && value <= 1 && value > 0) {
                        return FORMATTING[unit].SINGULAR_LABEL
                    }
                    return FORMATTING[unit].LABEL
                }
        }
    }

    return unit
}

export const hasLabelWhitespace = formattingType => unit => {
    if (FORMATTING[unit]) {
        switch (formattingType) {
            case FORMATTING_SHORT:
            case FORMATTING_LONG:
                return FORMATTING[unit][formattingType] ? (FORMATTING[unit][formattingType].NO_WHITESPACE === true ? false : true) : hasLabelWhitespace()(unit)
            default:
                return FORMATTING[unit].NO_WHITESPACE === true ? false : true
        }
    }

    return true
}

export const canLabelBeSmall = formattingType => unit => {
    if (FORMATTING[unit]) {
        switch (formattingType) {
            case FORMATTING_SHORT:
            case FORMATTING_LONG:
                return FORMATTING[unit][formattingType] ? (FORMATTING[unit][formattingType].CAN_BE_SMALL === false ? false : true) : canLabelBeSmall()(unit)
            default:
                return FORMATTING[unit].CAN_BE_SMALL === false ? false : true
        }
    }

    return true
}

export const getPrecisionUnit = unit => {
    switch (unit) {
        case POUNDS_OUNCES:
            return OUNCES
        case KG_GRAMS:
            return GRAMS
        case FEET_INCHES:
            return INCHES
        case METERS_CM:
            return CM
        default:
            return unit
    }
}

export const convertPrecision = (valueUnit, inputUnit, precision) => {
    if (valueUnit === MM_HG && inputUnit === CM_HG) {
        return precision + 1
    }
    return precision
}

export const getMinimum = unit => {
    return MININUM[unit] ? MININUM[unit] : {}
}

export const getMaximum = unit => {
    return MAXIMUM[unit] ? MAXIMUM[unit] : {}
}
