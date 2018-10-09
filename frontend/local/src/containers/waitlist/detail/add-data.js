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
import { SimpleUnitInput, UnitInputWithConversion } from "shared/forms/measurementFields"
import { POUNDS_OUNCES, POUNDS, KG, MM_HG, CM_HG, CM, FEET_INCHES, CELSIUS, FAHRENHEIT } from "shared/unitConversion/units"
import {
    weightValueToObject,
    weightObjectToValue,
    lengthValueToObject,
    lengthObjectToValue,
    pressureValueToObject,
    pressureObjectToValue,
    temperatureValueToObject,
    temperatureObjectToValue
} from "shared/unitConversion"
import { read } from "shared/modules/config"
import { WEIGHT_UNIT, LENGTH_UNIT, TEMPERATURE_UNIT, BLOOD_PRESSURE_UNIT } from "../../../modules/config"
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
                                    <Fields label="Height" names={["has_height", "height", "focus"]} unit="cm" component={renderHeight} change={change} />

                                    <Fields label="Weight" names={["has_weight", "weight", "focus"]} component={renderWeight} change={change} />
                                </div>

                                <h3>Vital Signs</h3>

                                <div>
                                    <Fields
                                        label="Body temperature"
                                        names={["has_temperature", "temperature", "focus"]}
                                        unit="°C"
                                        component={renderTemperature}
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

const submitOnEnter = dispatch => e => {
    if (e.key === "Enter") {
        e.preventDefault()
        dispatch(submit("medical-data"))
    }
}

const renderFieldWithUnit = connect()(fields => {
    return (
        <div
            className={classnames("section", {
                open: fields[fields.names[0]].input.value
            })}
        >
            {fields[fields.names[0]].input.value && (
                <div className="form-row">
                    <div className="col-sm">
                        <SimpleUnitInput
                            input={{
                                ...fields[fields.names[1]].input,
                                onKeyPress: fields.dispatch(submitOnEnter),
                                autoFocus: fields.focus.input.value === fields.names[1] && true
                            }}
                            meta={fields[fields.names[1]].meta}
                            label={fields.label}
                            placeholder={fields.label}
                            unit={fields.unit}
                            precision={fields.precision ? fields.precision : 0}
                        />
                    </div>
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
    )
})

const renderHeight = connect()(fields => {
    return (
        <div
            className={classnames("section", {
                open: fields[fields.names[0]].input.value
            })}
        >
            {fields[fields.names[0]].input.value && (
                <div className="form-row">
                    <div className="col-sm">
                        {fields.dispatch(read(LENGTH_UNIT)) === FEET_INCHES ? (
                            <UnitInputWithConversion
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                inputUnit={FEET_INCHES}
                                valueUnit={CM}
                                valuePrecision={0}
                                valueToObject={lengthValueToObject}
                                objectToValue={lengthObjectToValue}
                            />
                        ) : (
                            <SimpleUnitInput
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                unit={CM}
                                precision={1}
                            />
                        )}
                    </div>
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
    )
})

const renderWeight = connect()(fields => {
    return (
        <div
            className={classnames("section", {
                open: fields[fields.names[0]].input.value
            })}
        >
            {fields[fields.names[0]].input.value && (
                <div className="form-row">
                    <div className="col-sm">
                        {fields.dispatch(read(WEIGHT_UNIT)) === POUNDS_OUNCES ? (
                            <UnitInputWithConversion
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                inputUnit={POUNDS}
                                valueUnit={KG}
                                valuePrecision={8}
                                inputPrecision={1}
                                valueToObject={weightValueToObject}
                                objectToValue={weightObjectToValue}
                            />
                        ) : (
                            <SimpleUnitInput
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                unit={KG}
                                precision={1}
                            />
                        )}
                    </div>
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
    )
})

const renderTemperature = connect()(fields => {
    return (
        <div
            className={classnames("section", {
                open: fields[fields.names[0]].input.value
            })}
        >
            {fields[fields.names[0]].input.value && (
                <div className="form-row">
                    <div className="col-sm">
                        {fields.dispatch(read(TEMPERATURE_UNIT)) === FAHRENHEIT ? (
                            <UnitInputWithConversion
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                inputUnit={FAHRENHEIT}
                                valueUnit={CELSIUS}
                                valuePrecision={8}
                                inputPrecision={1}
                                valueToObject={temperatureValueToObject}
                                objectToValue={temperatureObjectToValue}
                            />
                        ) : (
                            <SimpleUnitInput
                                input={{
                                    ...fields[fields.names[1]].input,
                                    onKeyPress: submitOnEnter(fields.dispatch),
                                    autoFocus: fields.focus.input.value === fields.names[1] && true
                                }}
                                meta={fields[fields.names[1]].meta}
                                label={fields.label}
                                placeholder={fields.label}
                                unit={CELSIUS}
                                precision={1}
                            />
                        )}
                    </div>
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
    )
})

const renderBloodPressure = connect()(fields => {
    return (
        <div
            className={classnames("section", {
                open: fields[fields.names[0]].input.value
            })}
        >
            {fields[fields.names[0]].input.value && (
                <div>
                    <div className="form-row title">
                        <h4>{fields.label}</h4>
                        <div className="col-sm">
                            {fields.dispatch(read(BLOOD_PRESSURE_UNIT)) === CM_HG ? (
                                <UnitInputWithConversion
                                    input={{
                                        ...fields.pressure.systolic.input,
                                        onKeyPress: fields.dispatch(submitOnEnter),
                                        autoFocus: fields.focus.input.value === "pressure" && true
                                    }}
                                    meta={fields.pressure.systolic.meta}
                                    label="Systolic"
                                    placeholder="Systolic"
                                    inputUnit={CM_HG}
                                    valueUnit={MM_HG}
                                    valuePrecision={0}
                                    inputPrecision={1}
                                    valueToObject={pressureValueToObject}
                                    objectToValue={pressureObjectToValue}
                                />
                            ) : (
                                <SimpleUnitInput
                                    input={{
                                        ...fields.pressure.systolic.input,
                                        onKeyPress: fields.dispatch(submitOnEnter),
                                        autoFocus: fields.focus.input.value === "pressure" && true
                                    }}
                                    meta={fields.pressure.systolic.meta}
                                    label="Systolic"
                                    placeholder="Systolic"
                                    unit={MM_HG}
                                    precision={0}
                                />
                            )}
                        </div>
                        <div className="col-sm">
                            {fields.dispatch(read(BLOOD_PRESSURE_UNIT)) === CM_HG ? (
                                <UnitInputWithConversion
                                    input={{
                                        ...fields.pressure.diastolic.input,
                                        onKeyPress: fields.dispatch(submitOnEnter)
                                    }}
                                    meta={fields.pressure.diastolic.meta}
                                    label="Diastolic"
                                    placeholder="Diastolic"
                                    inputUnit={CM_HG}
                                    valueUnit={MM_HG}
                                    valuePrecision={0}
                                    inputPrecision={1}
                                    valueToObject={pressureValueToObject}
                                    objectToValue={pressureObjectToValue}
                                />
                            ) : (
                                <SimpleUnitInput
                                    input={{
                                        ...fields.pressure.diastolic.input,
                                        onKeyPress: fields.dispatch(submitOnEnter)
                                    }}
                                    meta={fields.pressure.diastolic.meta}
                                    label="Diastolic"
                                    placeholder="Diastolic"
                                    unit={MM_HG}
                                    precision={0}
                                />
                            )}
                        </div>
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
})

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
