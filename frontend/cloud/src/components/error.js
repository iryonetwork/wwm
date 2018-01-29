import React from "react"

const Error = props => {
    if (props.error) {
        return (
            <div className="alert alert-danger" role="alert">
                {props.code ? `${props.code}:` : ""} {props.error}
            </div>
        )
    }

    return null
}

export default Error
