import React from "react"
import classnames from "classnames"
import Select from "react-select"
import _ from "lodash"

const validateRequired = value => (!_.isEmpty(value) ? undefined : "Required")

const renderInput = ({ input, optional, disabled, label, placeholder, type, meta: { touched, error } }) => (
    <label>
        <input
            {...input}
            disabled={disabled}
            className={classnames("form-control", { "is-invalid": touched && error })}
            placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
            type={type || "text"}
        />

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderHorizontalInput = ({ input, optional, label, placeholder, unit, hideLabel, type, meta: { touched, error } }) => (
    <div className="form-row">
        <div className="label">{!hideLabel && label}</div>
        <div className={"form-group " + (unit ? "col-sm-4" : "col-sm-6")}>
            <input
                {...input}
                className={classnames("form-control", { "is-invalid": touched && error })}
                placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
                type={type || "text"}
            />
            {touched && error && <div className="invalid-feedback">{error}</div>}
        </div>
        {unit && <div className="col-sm-2">{unit}</div>}
    </div>
)

const renderSelect = ({ input, disabled, pristine, label, placeholder, options, meta: { touched, error } }) => (
    <label>
        <select {...input} disabled={disabled} className={classnames("form-control", { "is-invalid": touched && error, selected: input.value !== "" })}>
            <option value="" disabled>
                {placeholder ? placeholder : label}
            </option>
            {(options || []).map(option => (
                <option value={option.value} key={option.value}>
                    {option.label}
                </option>
            ))}
        </select>

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderReactSelect = ({ input, label, loadOptions, meta: { touched, error } }) => {
    return (
        <label>
            <Select.Async
                value={input.value}
                className={classnames("form-control", { "is-invalid": touched && error })}
                multi={false}
                loadOptions={loadOptions}
                onChange={input.onChange}
            />
            <span>{label}</span>
            {touched && error && <div className="invalid-feedback">{error}</div>}
        </label>
    )
}

const renderHorizontalSelect = ({ input, pristine, label, placeholder, unit, options, meta: { touched, error } }) => (
    <div className="form-row">
        <div className="label">{label}</div>
        <div className={"form-group " + (unit ? "col-sm-4" : "col-sm-6")}>
            <select {...input} className={classnames("form-control", { "is-invalid": touched && error, selected: input.value !== "" })}>
                <option value="" disabled>
                    {placeholder ? placeholder : label}
                </option>
                {options.map(option => (
                    <option value={option.value} key={option.value}>
                        {option.label}
                    </option>
                ))}
            </select>
            {touched && error && <div className="invalid-feedback">{error}</div>}
        </div>
        {unit && <div className="col-sm-2">{unit}</div>}
    </div>
)

const renderNumericalValuesRadio = ({ input, className, label, options, hideLabel, meta: { touched, error } }) => (
    <div className={classnames("form-inline-container", className)}>
        {!hideLabel && <span className="label">{label}</span>}
        {options.map((option, index) => (
            <div key={index} className="form-check form-check-inline">
                <input
                    {...input}
                    className="form-check-input"
                    type="radio"
                    id={`${input.name}${index}`}
                    checked={Number(input.value) === option.value}
                    value={option.value}
                />
                <label className="form-check-label" htmlFor={`${input.name}${index}`}>
                    {option.label}
                </label>
            </div>
        ))}
    </div>
)

const renderRadio = ({ input, className, label, options, hideLabel, meta: { touched, error } }) => (
    <div className={classnames("form-inline-container", className)}>
        {!hideLabel && <span className="label">{label}</span>}
        {options.map((option, index) => (
            <div key={index} className="form-check form-check-inline">
                <input
                    {...input}
                    className="form-check-input"
                    type="radio"
                    id={`${input.name}${index}`}
                    checked={input.value === option.value}
                    value={option.value}
                />
                <label className="form-check-label" htmlFor={`${input.name}${index}`}>
                    {option.label}
                </label>
            </div>
        ))}
    </div>
)

const renderHorizontalRadio = ({ input, className, label, options, hideLabel, meta: { touched, error } }) => (
    <div className="form-row">
        <div className="label">{label}</div>
        <div className="form-inline-container">
            {options.map((option, index) => (
                <div key={index} className="form-check form-check-inline">
                    <input
                        {...input}
                        className="form-check-input"
                        type="radio"
                        id={`${input.name}${index}`}
                        checked={input.value === option.value}
                        value={option.value}
                    />
                    <label className="form-check-label" htmlFor={`${input.name}${index}`}>
                        {option.label}
                    </label>
                </div>
            ))}
        </div>
    </div>
)

const renderTextarea = ({ input, rows, optional, label, placeholder, meta: { touched, error } }) => (
    <label>
        <textarea
            {...input}
            rows={rows || 5}
            className={classnames("form-control", { "is-invalid": touched && error })}
            placeholder={classnames(placeholder ? placeholder : label, { "(optional)": optional })}
        />

        <span>{label}</span>
        {touched && error && <div className="invalid-feedback">{error}</div>}
    </label>
)

const renderHabitFields = fields => (
    <div className="habits">
        <div className={classnames("complexHorizontalRadio", { subField: fields.subField })}>
            <div className="form-row">
                <div className="label">{fields.label}</div>
                <div className="form-group col-sm-6">
                    <div className="form-check form-check-inline">
                        <input
                            {...fields[fields.names[0]].input}
                            value="true"
                            checked={fields[fields.names[0]].input.value === "true"}
                            className="form-check-input"
                            type="radio"
                            id={fields.names[0] + "yes"}
                        />
                        <label className="form-check-label" htmlFor={fields.names[0] + "yes"}>
                            Yes
                        </label>
                    </div>
                    <div className="form-check form-check-inline">
                        <input
                            {...fields[fields.names[0]].input}
                            value="false"
                            checked={fields[fields.names[0]].input.value === "false"}
                            className="form-check-input"
                            type="radio"
                            id={fields.names[0] + "no"}
                        />
                        <label className="form-check-label" htmlFor={fields.names[0] + "no"}>
                            No
                        </label>
                    </div>
                </div>
            </div>
            {fields[fields.names[0]].input.value === (fields.commentWhen || "false") && (
                <div className="form-row comment">
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
    </div>
)

const renderComplexHorizontalRadio = fields => {
    return (
        <div className={classnames("complexHorizontalRadio", { subField: fields.subField })}>
            <div className="form-row">
                <div className="label">
                    {fields.subField && "• "}
                    {fields.label}
                </div>
                <div className="form-group col-sm-6">
                    {fields.options.map((option, index) => (
                        <div key={index} className="form-check form-check-inline">
                            <input
                                {...fields[fields.names[0]].input}
                                className="form-check-input"
                                type="radio"
                                id={`${fields.names[0]}${index}`}
                                checked={fields[fields.names[0]].input.value === option.value}
                                value={option.value}
                            />
                            <label className="form-check-label" htmlFor={`${fields.names[0]}${index}`}>
                                {option.label}
                            </label>
                        </div>
                    ))}
                </div>
            </div>
            {fields.names[1] &&
                (fields.commentAlways || fields[fields.names[0]].input.value === (fields.commentWhen || "false")) && (
                    <div className="form-row comment">
                        <div className="label" />
                        <div className={classnames({ "col-sm-4": fields.commentShort }, { "col-sm-6": !fields.commentShort })}>
                            <label>
                                {fields.commentMultiline ? (
                                    <textarea
                                        {...fields[fields.names[1]].input}
                                        className="form-control"
                                        placeholder={`${
                                            fields.commentPlaceholder ? fields.commentPlaceholder : fields.commentLabel ? fields.commentLabel : "Comment"
                                        } (optional)`}
                                    />
                                ) : (
                                    <input
                                        {...fields[fields.names[1]].input}
                                        className="form-control"
                                        placeholder={`${
                                            fields.commentPlaceholder ? fields.commentPlaceholder : fields.commentLabel ? fields.commentLabel : "Comment"
                                        } (optional)`}
                                    />
                                )}
                                <span className={classnames({ alwaysVisible: fields.commentLabelAlwaysVisible })}>{`${
                                    fields.commentLabel ? fields.commentLabel : "Comment"
                                } (optional)`}</span>
                            </label>
                        </div>
                    </div>
                )}
        </div>
    )
}

export {
    validateRequired,
    renderInput,
    renderHorizontalInput,
    renderSelect,
    renderReactSelect,
    renderHorizontalSelect,
    renderRadio,
    renderNumericalValuesRadio,
    renderHorizontalRadio,
    renderTextarea,
    renderHabitFields,
    renderComplexHorizontalRadio
}
