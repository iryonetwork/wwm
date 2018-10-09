import _ from "lodash"
import React from "react"
import classnames from "classnames"

import { getObjectFromValue } from "../../unitConversion"
import { getLabel, hasLabelWhitespace, canLabelBeSmall, FORMATTING_LONG, FORMATTING_SHORT, FORMATTING_DEFAULT } from "../../unitConversion/units"

export const LONG = FORMATTING_LONG
export const SHORT = FORMATTING_SHORT
export const DEFAULT = FORMATTING_DEFAULT

export const ValueWithUnit = ({ valueUnit, displayUnit, precision, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT

    if (!hasLabelWhitespace(formatting)(displayUnit)) {
        return (
            <span className="valueWithUnit" key={displayUnit}>
                {_.map(getObjectFromValue(valueUnit, displayUnit, precision ? precision : 0)(value), (val, unit) => (
                    <SimpleValueWithUnit unit={unit} formatting={formatting} value={val} />
                ))}
            </span>
        )
    }
    return _.map(getObjectFromValue(valueUnit, displayUnit, precision ? precision : 0)(value), (val, unit) => (
        <span className="valueWithUnit" key={unit}>
            <SimpleValueWithUnit unit={unit} formatting={formatting} value={val} />
        </span>
    ))
}

const SimpleValueWithUnit = ({ unit, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT

    return (
        <React.Fragment key={unit}>
            <React.Fragment>
                <span className={classnames("value", { withMargin: hasLabelWhitespace(formatting)(unit) })}>{value}</span>
                <span className={classnames("unit", { big: !canLabelBeSmall(formatting)(unit) })}>{getLabel(formatting)(unit, value)}</span>
            </React.Fragment>
        </React.Fragment>
    )
}

export const ValueWithUnitString = ({ valueUnit, stringUnit, precision, formatting, value }) => {
    formatting = formatting ? formatting : DEFAULT
    let object = getObjectFromValue(valueUnit, stringUnit, precision)(value)

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
