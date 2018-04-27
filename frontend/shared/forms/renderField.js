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
                <input {...input} className="form-check-input" type="radio" id={`${input.name}${index}`} value={option.value} />
                <label className="form-check-label" htmlFor={`${input.name}${index}`}>
                    {option.label}
                </label>
            </div>
        ))}
    </div>
)

const renderTextarea = ({ input, rows, optional, label, meta: { touched, error } }) => (
    <label>
        <textarea
            {...input}
            rows={rows || 5}
            className={classnames("form-control", { "is-invalid": touched && error })}
            placeholder={classnames(label, { "(optional)": optional })}
        />

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderHabitFields = fields => (
    <div className="habits">
        <div className="form-row">
            <div className="label">{fields.label}</div>
            <div className="form-group col-sm-6">
                <div className="form-check form-check-inline">
                    <input {...fields[fields.names[0]].input} value={true} className="form-check-input" type="radio" id={fields.names[0] + "yes"} />
                    <label className="form-check-label" htmlFor={fields.names[0] + "yes"}>
                        Yes
                    </label>
                </div>
                <div className="form-check form-check-inline">
                    <input {...fields[fields.names[0]].input} value={false} className="form-check-input" type="radio" id={fields.names[0] + "no"} />
                    <label className="form-check-label" htmlFor={fields.names[0] + "no"}>
                        No
                    </label>
                </div>
            </div>
        </div>
        {fields[fields.names[0]].input.value === (fields.commentWhen || "false") && (
            <div className="row comment">
                <div className="label" />
                <div className="col-sm-4">
                    <label>
                        <input {...fields[fields.names[1]].input} className="form-control" placeholder="Comment (optional)" />
                        <span>Comment</span>
                    </label>
                </div>
            </div>
        )}
    </div>
)

export { renderInput, renderSelect, renderRadio, renderTextarea, renderHabitFields }
