import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { Fields, reduxForm, submit } from "redux-form"
import classnames from "classnames"
import { goBack, push } from "react-router-redux"
import moment from "moment"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient/card"
import Spinner from "shared/containers/spinner"
import { round } from "shared/utils"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { listAll, update } from "../../../modules/waitlist"
import { cardToObject } from "../../../modules/discovery"

import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"

const supportedSigns = ["height", "weight", "temperature", "pressure", "heart_rate", "oxygen_saturation"]
const complexSigns = { pressure: ["systolic", "diastolic"] }

const validate = form => {
    const errors = {}

    _.forEach(form, (value, key) => {
        if (key.indexOf("has_") === 0 && value) {
            let sign = key.slice(4)
            if (_.isObject(form[sign])) {
                errors[sign] = {}
                _.forEach(form[sign], (value, key) => {
                    if (!value) {
                        errors[sign][key] = "Required"
                    }
                })
            } else if (!form[sign]) {
                errors[sign] = "Required"
            }
        }
    })
    return errors
}

const onMedicalDataFormSubmit = (form, dispatch, props) => {
    let vitalSigns = {}

    _.forEach(form, (value, key) => {
        if (key.indexOf("has_") === 0 && value) {
            let sign = key.slice(4)
            vitalSigns[sign] = {}
            vitalSigns[sign].value = form[sign]

            if (form[sign] !== props.initialValues[sign]) {
                vitalSigns[sign].timestamp = moment().format()
            } else {
                vitalSigns[sign].timestamp = props.initialValues["timestmap_" + sign]
            }
        }
    })

    // set BMI if height and weight available
    if (vitalSigns.height && vitalSigns.weight) {
        vitalSigns.bmi = {}
        vitalSigns.bmi.value = round(vitalSigns.weight.value / vitalSigns.height.value / vitalSigns.height.value * 10000, 2)
        vitalSigns.bmi.timestamp = moment.max(moment(vitalSigns.height.timestamp), moment(vitalSigns.weight.timestamp)).format()
    } else if (props.item.vitalSigns && props.item.vitalSigns.bmi) {
        // remove bmi if bmi cannot be set anymore
        delete vitalSigns.bmi
    }

    let newItem = _.clone(props.item)
    newItem.vitalSigns = vitalSigns

    dispatch(update(props.match.params.waitlistID, newItem))
        .then(data => {
            dispatch(listAll(props.match.params.waitlistID))
            dispatch(goBack())
        })
        .catch(ex => {})
}

class MedicalData extends React.Component {
    constructor(props) {
        super(props)
        if (!props.item) {
            props.listAll(props.match.params.waitlistID)
        }
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.goBack()
            this.props.open("Waitlist item was not found", "", COLOR_DANGER, 5)
        }
    }

    render() {
        let { handleSubmit, item, change, history } = this.props
        return (
            <Modal>
                <div className="medical-data">
                    {item ? (
                        <form onSubmit={handleSubmit}>
                            <div className="modal-header">
                                <Patient data={item.patient && cardToObject({ connections: item.patient })} />
                                <h1>
                                    <MedicalDataIcon />
                                    Add Medical Data
                                </h1>
                            </div>

                            <div className="modal-body">
                                <h3>Body Measurements</h3>
                                <div>
                                    <Fields
                                        label="Height"
                                        names={["has_height", "height", "focus"]}
                                        unit="cm"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />

                                    <Fields
                                        label="Weight"
                                        names={["has_weight", "weight", "focus"]}
                                        unit="kg"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />
                                </div>

                                <h3>Vital Signs</h3>

                                <div>
                                    <Fields
                                        label="Body temperature"
                                        names={["has_temperature", "temperature", "focus"]}
                                        unit="°C"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />

                                    <Fields
                                        label="Blood pressure"
                                        names={["has_pressure", "pressure.systolic", "pressure.diastolic", "focus"]}
                                        component={renderBloodPressure}
                                        change={change}
                                    />

                                    <Fields
                                        label="Heart rate"
                                        names={["has_heart_rate", "heart_rate", "focus"]}
                                        unit="bpm"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />

                                    {/*<Fields
                                        label="Hearing screaning"
                                        names={[
                                            "has_hearing",
                                            "hearing.left.500",
                                            "hearing.left.1000",
                                            "hearing.left.2000",
                                            "hearing.left.3000",
                                            "hearing.left.4000",
                                            "hearing.left.6000",
                                            "hearing.left.8000",
                                            "hearing.right.500",
                                            "hearing.right.1000",
                                            "hearing.right.2000",
                                            "hearing.right.3000",
                                            "hearing.right.4000",
                                            "hearing.right.6000",
                                            "hearing.right.8000"
                                        ]}
                                        component={renderHearingScreening}
                                        change={change}
                                    />*/}

                                    {/*<Fields
                                        label="Visual screening"
                                        names={[
                                            "has_visual_screening",
                                            "visual_screening.left.distance",
                                            "visual_screening.left.value",
                                            "visual_screening.right.distance",
                                            "visual_screening.right.value"
                                        ]}
                                        unit="%"
                                        component={renderVisualScreening}
                                        change={change}
                                    />*/}

                                    <Fields
                                        label="Oxygen saturation"
                                        names={["has_oxygen_saturation", "oxygen_saturation", "focus"]}
                                        unit="%"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />
                                    {/*<Fields
                                        label="Lung function test"
                                        names={["has_lung_function_test", "heart_lung_function_test"]}
                                        unit="%"
                                        component={renderFieldWithUnit}
                                        change={change}
                                    />
                                    */}
                                </div>
                            </div>

                            <div className="modal-footer">
                                <div className="form-row">
                                    <div className="col-sm-4" />
                                    <div className="col-sm-4">
                                        <button
                                            type="button"
                                            tabIndex="-1"
                                            className="btn btn-secondary btn-block"
                                            datadismiss="modal"
                                            onClick={() => history.goBack()}
                                        >
                                            Cancel
                                        </button>
                                    </div>

                                    <div className="col-sm-4">
                                        <button type="submit" className="float-right btn btn-primary btn-block">
                                            Add
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>
                    ) : (
                        <div className="modal-body">
                            <Spinner />
                        </div>
                    )}
                </div>
            </Modal>
        )
    }
}

// const renderFieldWithUnits = fields => (
//     <div className={classnames("section", { open: fields[fields.names[0]].input.value })}>
//         {fields[fields.names[0]].input.value && (
//             <div className="form-row">
//                 <div className="col-sm-4">
//                     <label>
//                         <input {...fields[fields.names[1]].input} type="number" className="form-control" placeholder={fields.label} />

//                         <span>{fields.label}</span>
//                     </label>
//                 </div>

//                 <div className="col-sm-2">
//                     <select {...fields[fields.names[2]].input} className="form-control" validate={required}>
//                         {fields.units.map(unit => (
//                             <option key={unit.value} value={unit.value}>
//                                 {unit.label}
//                             </option>
//                         ))}
//                     </select>
//                 </div>

//                 <button className="btn btn-link remove" onClick={() => fields.change(fields.names[0], false)}>
//                     <NegativeIcon />
//                     Remove
//                 </button>
//             </div>
//         )}
//         {!fields[fields.names[0]].input.value && (
//             <button className="btn btn-link" onClick={() => fields.change(fields.names[0], true)}>
//                 Add {fields.label}
//             </button>
//         )}
//     </div>
// )

const SimpleInput = connect()(({ input, dispatch }) => (
    <input
        {...input}
        onKeyPress={ev => {
            if (ev.key === "Enter") {
                ev.preventDefault()
                dispatch(submit("medical-data"))
            }
        }}
    />
))

const renderFieldWithUnit = connect()(fields => (
    <div
        className={classnames("section", {
            open: fields[fields.names[0]].input.value
        })}
    >
        {fields[fields.names[0]].input.value && (
            <div className="form-row">
                <div className="col-sm-4">
                    <label>
                        <SimpleInput
                            input={{
                                ...fields[fields.names[1]].input,
                                type: "number",
                                className: classnames("form-control", {
                                    "is-invalid": fields[fields.names[1]].meta.touched && fields[fields.names[1]].meta.error
                                }),
                                placeholder: fields.label,
                                autoFocus: fields.focus.input.value === fields.names[1] && true
                            }}
                        />
                        <span>{fields.label}</span>
                        {fields[fields.names[1]].meta.touched &&
                            fields[fields.names[1]].meta.error && <div className="invalid-feedback">{fields[fields.names[1]].meta.error}</div>}
                    </label>
                </div>

                <div className="col-sm-2 unit">{fields.unit}</div>

                <button className="btn btn-link remove" onClick={() => fields.change(fields.names[0], false)}>
                    <NegativeIcon />
                    Remove
                </button>
            </div>
        )}
        {!fields[fields.names[0]].input.value && (
            <button
                className="btn btn-link"
                onClick={() => {
                    fields.change(fields.names[0], true)
                    fields.change("focus", fields.names[1])
                }}
            >
                Add {fields.label}
            </button>
        )}
    </div>
))

const renderBloodPressure = fields => (
    <div
        className={classnames("section", {
            open: fields[fields.names[0]].input.value
        })}
    >
        {fields[fields.names[0]].input.value && (
            <div>
                <div className="form-row title">
                    <h4>{fields.label}</h4>
                    <div className="col-sm-4">
                        <label>
                            <SimpleInput
                                input={{
                                    ...fields.pressure.systolic.input,
                                    type: "number",
                                    className: classnames("form-control", {
                                        "is-invalid": fields.pressure.systolic.meta.touched && fields.pressure.systolic.meta.error
                                    }),
                                    placeholder: "Systolic",
                                    autoFocus: fields.focus.input.value === "pressure" && true
                                }}
                            />
                            <span>Systolic</span>
                            {fields.pressure.systolic.meta.touched &&
                                fields.pressure.systolic.meta.error && <div className="invalid-feedback">{fields.pressure.systolic.meta.error}</div>}
                        </label>
                    </div>

                    <div className="col-sm-2 unit">mmHg</div>

                    <div className="col-sm-4">
                        <label>
                            <SimpleInput
                                input={{
                                    ...fields.pressure.diastolic.input,
                                    type: "number",
                                    className: classnames("form-control", {
                                        "is-invalid": fields.pressure.diastolic.meta.touched && fields.pressure.diastolic.meta.error
                                    }),
                                    placeholder: "Diastolic"
                                }}
                            />
                            <span>Diastolic</span>
                            {fields.pressure.diastolic.meta.touched &&
                                fields.pressure.diastolic.meta.error && <div className="invalid-feedback">{fields.pressure.diastolic.meta.error}</div>}
                        </label>
                    </div>

                    <div className="col-sm-2 unit">mmHg</div>

                    <button className="btn btn-link remove" onClick={() => fields.change(fields.names[0], false)}>
                        <NegativeIcon />
                        Remove
                    </button>
                </div>
            </div>
        )}
        {!fields[fields.names[0]].input.value && (
            <button
                className="btn btn-link"
                onClick={() => {
                    fields.change(fields.names[0], true)
                    fields.change("focus", "pressure")
                }}
            >
                Add {fields.label}
            </button>
        )}
    </div>
)

// const renderHearingScreening = fields => {
//     let frequencies = [500, 1000, 2000, 3000, 4000, 6000, 8000]
//     return (
//         <div className={classnames("section", { open: fields[fields.names[0]].input.value })}>
//             {fields[fields.names[0]].input.value && (
//                 <div>
//                     <h4>{fields.label}</h4>
//                     <div className="hearing-row">
//                         <div className="label">Left ear</div>
//                         {frequencies.map(f => (
//                             <div className="col-sm-1" key={`left-${f}`}>
//                                 <label>
//                                     <input {...fields.hearing.left[f].input} type="number" className="form-control" placeholder={f} />
//                                     <span>{`${f} Hz`}</span>
//                                 </label>
//                             </div>
//                         ))}
//                         <div className="col-sm-1 unit">dB</div>

//                         <button className="btn btn-link remove" onClick={() => fields.change(fields.names[0], false)}>
//                             <NegativeIcon />
//                             Remove
//                         </button>
//                     </div>

//                     <div className="hearing-row">
//                         <div className="label">Right ear</div>
//                         {frequencies.map(f => (
//                             <div className="col-sm-1" key={`right-${f}`}>
//                                 <label>
//                                     <input {...fields.hearing.right[f].input} type="number" className="form-control" placeholder={f} />
//                                     <span>{`${f} Hz`}</span>
//                                 </label>
//                             </div>
//                         ))}
//                         <div className="col-sm-1 unit">dB</div>
//                     </div>
//                 </div>
//             )}
//             {!fields[fields.names[0]].input.value && (
//                 <button className="btn btn-link" onClick={() => fields.change(fields.names[0], true)}>
//                     Add {fields.label}
//                 </button>
//             )}
//         </div>
//     )
// }

// const renderVisualScreening = fields => (
//     <div className={classnames("section", { open: fields[fields.names[0]].input.value })}>
//         {fields[fields.names[0]].input.value && (
//             <div className="form-row title">
//                 <h4>{fields.label}</h4>
//                 <div className="col-sm-2 visual">
//                     <label>
//                         <input {...fields.visual_screening.left.distance} type="number" className="form-control" placeholder="Left eye" />
//                         <span>OS - left eye</span>
//                     </label>
//                 </div>
//                 <div className="col-sm-2">
//                     <label>
//                         <input {...fields.visual_screening.left.value} type="number" className="form-control" />
//                     </label>
//                 </div>

//                 <div className="col-sm-1" />

//                 <div className="col-sm-2 visual">
//                     <label>
//                         <input {...fields.visual_screening.right.distance} type="number" className="form-control" placeholder="Right eye" />
//                         <span>OS - right eye</span>
//                     </label>
//                 </div>
//                 <div className="col-sm-2">
//                     <label>
//                         <input {...fields.visual_screening.right.value} type="number" className="form-control" />
//                     </label>
//                 </div>

//                 <button className="btn btn-link remove" onClick={() => fields.change(fields.names[0], false)}>
//                     <NegativeIcon />
//                     Remove
//                 </button>
//             </div>
//         )}
//         {!fields[fields.names[0]].input.value && (
//             <button className="btn btn-link" onClick={() => fields.change(fields.names[0], true)}>
//                 Add {fields.label}
//             </button>
//         )}
//     </div>
// )

MedicalData = reduxForm({
    form: "medical-data",
    validate,
    onSubmit: onMedicalDataFormSubmit
})(MedicalData)

MedicalData = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]

        let initialValues = {}

        // initialize complex signs for correct validation (workaround for bug of handling '<Fields />' in redux-form)
        _.forEach(complexSigns, (obj, sign) => {
            initialValues[sign] = {}
            _.forEach(obj, key => {
                initialValues[sign][key] = undefined
            })
        })

        let selectedSign = props.match.params.sign
        if (selectedSign && _.includes(supportedSigns, selectedSign)) {
            initialValues["has_" + selectedSign] = true
            initialValues["focus"] = selectedSign
        }

        if (item) {
            _.forEach(item.vitalSigns || {}, (obj, key) => {
                initialValues[key] = obj.value
                initialValues["has_" + key] = true
                initialValues["timestmap_" + key] = obj.timestamp
            })
        }

        return {
            listed: state.waitlist.listed,
            item,
            initialValues
        }
    },
    {
        listAll,
        update,
        open,
        goBack,
        push
    }
)(MedicalData)

export default MedicalData
