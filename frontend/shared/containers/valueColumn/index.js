import React from "react"

export default ({ value, unit, label, codes, width }) => {
    // don't render if empty
    if (value === undefined) {
        return null
    }

    // convert for a code
    if (codes && codes.length > 0) {
        value = codes.reduce((acc, code) => {
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
