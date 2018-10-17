import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import classnames from "classnames"
import moment from "moment"
import { UncontrolledTooltip } from "reactstrap"
import { ReactComponent as WarningIcon } from "shared/icons/warning.svg"

import {
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
import { ValueWithUnit, SHORT } from "../measurement"
import "./style.css"

const card = ({ id, name, value, valueUnit, displayUnit, precision, timestamp, timestampWarning, consultationTooltipOn, onClick, isActive }) => {
    return (
        <div
            key={id}
            className={classnames("card", {
                active: isActive,
                clickable: onClick ? true : false
            })}
            onClick={onClick && onClick()}
        >
            <div className="card-header">{name}</div>
            <div className="card-body">
                <div className="card-text">
                    <p>
                        <ValueWithUnit valueUnit={valueUnit} displayUnit={displayUnit} precision={precision ? precision : 0} value={value} formatting={SHORT} />
                    </p>
                </div>
            </div>
            <div
                className={classnames("card-footer", {
                    timestampWarning: timestampWarning || !timestamp
                })}
            >
                {consultationTooltipOn ? (
                    <React.Fragment>
                        <a href="/" id={`${id}Tooltip`}>
                            {(timestampWarning || !timestamp) && <WarningIcon />}
                            {timestamp ? moment(timestamp).format("Do MMM Y") : "Unknown date"}
                        </a>
                        <UncontrolledTooltip placement="bottom-start" target={`${id}Tooltip`}>
                            {timestampWarning ? "This reading was done in the past encounter." : "This reading was done in the current encounter."}
                        </UncontrolledTooltip>
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        {(timestampWarning || !timestamp) && <WarningIcon />}
                        {timestamp ? moment(timestamp).format("Do MMM Y") : "Unknown date"}
                    </React.Fragment>
                )}
            </div>
        </div>
    )
}

export const VitalSignCard = connect((state, props) => {
    return {
        displayUnit: props.unit,
        valueUnit: props.unit
    }
}, {})(card)

export const HeightVitalSignCard = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, LENGTH_UNIT, CM)
    let valueUnit = props.unit ? props.unit : CM

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(card)

export const WeightVitalSignCard = connect((state, props) => {
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
}, {})(card)

export const TemperatureVitalSignCard = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, TEMPERATURE_UNIT, CELSIUS)
    let valueUnit = props.unit ? props.unit : CELSIUS

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(card)

export const BloodPressureVitalSignCard = connect((state, props) => {
    let configuredInputUnit = _.get(state.config, BLOOD_PRESSURE_UNIT, MM_HG)
    let valueUnit = props.unit ? props.unit : MM_HG

    return {
        displayUnit: configuredInputUnit,
        valueUnit: valueUnit
    }
}, {})(card)
