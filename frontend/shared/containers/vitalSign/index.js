import React from "react"
import classnames from "classnames"
import moment from "moment"
import { UncontrolledTooltip } from "reactstrap"
import { ReactComponent as WarningIcon } from "shared/icons/warning.svg"

import { ValueWithUnit, SHORT } from "../measurement"
import "./style.css"

export default ({ id, name, value, valueUnit, displayUnit, precision, timestamp, timestampWarning, consultationTooltipOn, onClick, isActive }) => {
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
