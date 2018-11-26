import React from "react"
import _ from "lodash"

export default ({ value, unit, label, codes, width }) => {
    // don't render if empty
    if (value === undefined) {
        return null
    }

    // convert for a code
    if (codes && _.size(codes) > 0) {
        value = _.reduce(codes, (acc, code) => {
            if (code.id === value) {
                return code.title
            }
            return acc
        }, undefined)
        // don't render if code is not found
        if (value === undefined) {
            return null
        }
    }

    return (
        <div className={`col-sm-${width}`}>
            <div className="label" key="label">
                {label}
            </div>
            <div className="value" key="value">
                {value}
                {unit && ` ${unit}`}
            </div>
        </div>
    )
}
