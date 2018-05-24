import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { Fields, reduxForm } from "redux-form"
import classnames from "classnames"
import { goBack } from "react-router-redux"
import moment from "moment"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { listAll, update } from "../../../modules/waitlist"
import { cardToObject } from "../../../modules/discovery"

import { ReactComponent as MedicalDataIcon } from "shared/icons/vitalsigns.svg"
import { ReactComponent as NegativeIcon } from "shared/icons/negative.svg"


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

class MedicalData extends React.Component {
    constructor(props) {
        super(props)
        if (!props.item) {
            props.listAll(props.match.params.waitlistID)
        }

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.history.goBack()
            setTimeout(() => this.props.open("Waitlist item was not found", "", COLOR_DANGER, 5), 100)
        }
    }

    handleSubmit(form) {
        let vitalSigns = {}
        _.forEach(form, (value, key) => {
            if (key.indexOf("has_") === 0 && value) {
                let sign = key.slice(4)
                vitalSigns[sign] = {}
                vitalSigns[sign].value = form[sign]

                if (form[sign] !== this.props.initialValues[sign]) {
                    vitalSigns[sign].timestamp = moment().format("X")
                } else {
                    vitalSigns[sign].timestamp = this.props.initialValues["timestmap_" + sign]
                }
            }
        })

        // set BMI if height and weight available
        if (vitalSigns.height && vitalSigns.weight) {
            vitalSigns.bmi = {}
            vitalSigns.bmi.value = round(
                vitalSigns.weight.value /
                    vitalSigns.height.value /
                    vitalSigns.height.value *
                    10000,
                2
            )
            vitalSigns.bmi.timestamp = (moment(vitalSigns.height.timestamp, "X").isAfter(moment(vitalSigns.weight.timestamp, "X")) ? moment(vitalSigns.height.timestamp, "X") : moment(vitalSigns.weight.timestamp, "X")).format("X")
        }

        this.props.item.vitalSigns = vitalSigns
        this.props
            .update(this.props.match.params.waitlistID, this.props.item)
            .then(data => {
                this.props.listAll(this.props.match.params.waitlistID)
                this.props.goBack()
            })
            .catch(ex => {})
    }

    render() {
        let props = this.props
        return (
            <Modal>
                <div className="medical-data">
                    <form onSubmit={props.handleSubmit(this.handleSubmit)}>
                        <div className="modal-header">
                            <Patient data={props.item.patient && cardToObject({ connections: props.item.patient })} />
                            <h1>
                                <MedicalDataIcon />
                                Add medical data
                            </h1>
                        </div>

                        <div className="modal-body">
                            <h3>Body measurements</h3>
                            <div>
                                <Fields
                                    label="Height"
                                    names={["has_height", "height"]}
                                    unit="cm"
                                    component={renderFieldWithUnit}
                                    change={props.change}
                                />

                                <Fields
                                    label="Weight"
                                    names={["has_weight", "weight"]}
                                    unit="kg"
                                    component={renderFieldWithUnit}
                                    change={props.change}
                                />
                            </div>

                            <h3>Vital signs</h3>

                            <div>
                                <Fields
                                    label="Body temperature"
                                    names={["has_temperature", "temperature"]}
                                    unit="°C"
                                    component={renderFieldWithUnit}
                                    change={props.change}
                                />

                                <Fields
                                    label="Blood pressure"
                                    names={["has_pressure", "pressure.systolic", "pressure.diastolic"]}
                                    component={renderBloodPressure}
                                    change={props.change}
                                />

                                <Fields
                                    label="Heart rate"
                                    names={["has_heart_rate", "heart_rate"]}
                                    unit="bpm"
                                    component={renderFieldWithUnit}
                                    change={props.change}
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
                                    change={props.change}
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
                                    change={props.change}
                                />*/}

                                <Fields
                                    label="Oxygen saturation"
                                    names={["has_oxygen_saturation", "oxygen_saturation"]}
                                    unit="%"
                                    component={renderFieldWithUnit}
                                    change={props.change}
                                />
                                {/*
                    <Fields
                        label="Lung function test"
                        names={["has_lung_function_test", "heart_lung_function_test"]}
                        unit="%"
                        component={renderFieldWithUnit}
                        change={props.change}
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
                                        className="btn btn-link btn-block"
                                        datadismiss="modal"
                                        onClick={() => props.history.goBack()}
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

const renderFieldWithUnit = fields => (
    <div className={classnames("section", { open: fields[fields.names[0]].input.value })}>
        {fields[fields.names[0]].input.value && (
            <div className="form-row">
                <div className="col-sm-4">
                    <label>
                        <input {...fields[fields.names[1]].input} type="number" className={classnames("form-control", { "is-invalid": fields[fields.names[1]].meta.touched && fields[fields.names[1]].meta.error })} placeholder={fields.label} />

                        <span>{fields.label}</span>
                        {fields[fields.names[1]].meta.touched && fields[fields.names[1]].meta.error && <div className="invalid-feedback">{fields[fields.names[1]].meta.error}</div>}
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
            <button className="btn btn-link" onClick={() => fields.change(fields.names[0], true)}>
                Add {fields.label}
            </button>
        )}
    </div>
)

const renderBloodPressure = fields => (
    <div className={classnames("section", { open: fields[fields.names[0]].input.value })}>
        {fields[fields.names[0]].input.value && (
            <div>
                <div className="form-row title">
                    <h4>{fields.label}</h4>
                    <div className="col-sm-4">
                        <label>
                            <input {...fields.pressure.systolic.input} type="number" className={classnames("form-control", { "is-invalid": fields.pressure.systolic.meta.touched && fields.pressure.systolic.meta.error })} placeholder="Systolic" />
                            <span>Systolic</span>
                            {fields.pressure.systolic.meta.touched && fields.pressure.systolic.meta.error && <div className="invalid-feedback">{fields.pressure.systolic.meta.error}</div>}
                        </label>
                    </div>

                    <div className="col-sm-2 unit">mmHg</div>

                    <div className="col-sm-4">
                        <label>
                            <input {...fields.pressure.diastolic.input} type="number" className={classnames("form-control", { "is-invalid": fields.pressure.diastolic.meta.touched && fields.pressure.diastolic.meta.error })} placeholder="Diastolic" />
                            <span>Diastolic</span>
                            {fields.pressure.diastolic.meta.touched && fields.pressure.diastolic.meta.error && <div className="invalid-feedback">{fields.pressure.diastolic.meta.error}</div>}
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
            <button className="btn btn-link" onClick={() => fields.change(fields.names[0], true)}>
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
    validate
})(MedicalData)

MedicalData = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]
        let initialValues = {}
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
        goBack
    }
)(MedicalData)

export default MedicalData

const round = (number, precision) => {
    var shift = function(number, precision) {
        var numArray = ("" + number).split("e")
        return +(numArray[0] + "e" + (numArray[1] ? +numArray[1] + precision : precision))
    }
    return shift(Math.round(shift(number, +precision)), -precision)
}
