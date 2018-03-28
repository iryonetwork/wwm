import React from "react"
import classnames from "classnames"

const renderInput = ({ input, optional, label, type, meta: { touched, error } }) => (
    <label>
        <input
            {...input}
            className={classnames("form-control", { "is-invalid": touched && error })}
            placeholder={classnames(label, { "(optional)": optional })}
            type={type || "text"}
        />

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderSelect = ({ input, pristine, label, options, meta: { touched, error } }) => (
    <label>
        <select {...input} className={classnames("form-control", { "is-invalid": touched && error, selected: input.value !== "" })}>
            <option value="" disabled>
                {label}
            </option>
            {options.map(option => (
                <option value={option.value} key={option.value}>
                    {option.label}
                </option>
            ))}
        </select>

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderRadio = ({ input, className, label, options, meta: { touched, error } }) => (
    <div className={classnames("form-inline-container", className)}>
        <span className="label">{label}</span>
        {options.map((option, index) => (
            <div key={index} className="form-check form-check-inline">
                <input className="form-check-input" type="radio" name={input.name} id={`${input.name}${index}`} value={option.value} />
                <label className="form-check-label" htmlFor={`${input.name}${index}`}>
                    {option.label}
                </label>
            </div>
        ))}
    </div>
)

export { renderInput, renderSelect, renderRadio }
