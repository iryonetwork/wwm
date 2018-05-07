import React, { Component } from "react"
import { connect } from "react-redux"
import { push } from "react-router-redux"
import { Field, reduxForm } from "redux-form"
//import PropTypes from "prop-types"
//import classnames from "classnames"

import { get, cardToObject } from "../../../modules/discovery"
import { add } from "../../../modules/waitlist"
import { renderInput, renderRadio, renderTextarea } from "shared/forms/renderField"
import { yesNoOptions } from "shared/forms/options"
import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"
import "./style.css"

const priorityOptions = [{ value: 1, label: "Yes" }, { value: 4, label: "No" }]

class AddToWaitlist extends Component {
    constructor(props) {
        super(props)
        props.get(props.match.params.patientID)
        this.onSubmit = this.onSubmit.bind(this)
    }

    componentDidMount() {
        document.body.style.overflow = "hidden"
    }

    componentWillUnmount() {
        document.body.style.overflow = "auto"
    }

    onSubmit(formData) {
        this.props.add(formData, this.props.patient)
    }

    render() {
        const { fetching, patient, push } = this.props
        if (fetching || !patient) {
            return (
                <React.Fragment>
                    <div className="add-to-waitlist modal fade show" tabIndex="-1" role="dialog">
                        <div className="modal-dialog" role="document">
                            <div className="modal-content">
                                <div className="modal-header">
                                    <h1>Add to Waiting List</h1>
                                </div>
                                <form>
                                    <div className="modal-body">
                                        <Spinner />
                                    </div>

                                    <div className="modal-footer">
                                        <div className="form-row">
                                            <div className="col-sm-4">
                                                <button
                                                    type="button"
                                                    tabIndex="-1"
                                                    className="btn btn-link btn-block"
                                                    datadismiss="modal"
                                                    disabled
                                                    onClick={() => {
                                                        push("/")
                                                    }}
                                                >
                                                    Cancel
                                                </button>
                                            </div>
                                            <div className="col-sm-4" />
                                            <div className="col-sm-4">
                                                <button type="submit" className="float-right btn btn-primary btn-block" disabled>
                                                    Add
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>

                    <div className="modal-backdrop fade show" />
                </React.Fragment>
            )
        }

        const p = cardToObject(patient)
        const { handleSubmit } = this.props

        return (
            <React.Fragment>
                <div className="add-to-waitlist modal fade show" tabIndex="-1" role="dialog">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">
                            <div className="modal-header">
                                <Patient data={p} />
                                <h1>Add to Waiting List</h1>
                            </div>
                            <form onSubmit={handleSubmit(this.onSubmit)}>
                                <div className="modal-body">
                                    <div className="form-row">
                                        <Field name="priority" component={renderRadio} label="Urgent?" options={priorityOptions} />
                                    </div>

                                    <div className="form-row">
                                        <Field name="mainComplaint" component={renderInput} label="Main complaint" />
                                    </div>

                                    <div className="form-row details">
                                        <Field name="mainComplaintDetails" component={renderTextarea} optional={true} label="Details" />
                                    </div>

                                    {/* <div className="form-row">
                                        <Field name="doctor" component={renderSelect} options={doctorOptions} label="Doctor" />
                                    </div> */}
                                </div>

                                <div className="modal-footer">
                                    <div className="form-row">
                                        <div className="col-sm-4">
                                            <button
                                                type="button"
                                                tabIndex="-1"
                                                className="btn btn-link btn-block"
                                                datadismiss="modal"
                                                onClick={() => {
                                                    push("/")
                                                }}
                                            >
                                                Cancel
                                            </button>
                                        </div>
                                        <div className="col-sm-4" />
                                        <div className="col-sm-4">
                                            <button type="submit" className="float-right btn btn-primary btn-block">
                                                Add
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>

                <div className="modal-backdrop fade show" />
            </React.Fragment>
        )
    }
}

AddToWaitlist = reduxForm({
    form: "addToWaitlist",
    initialValues: {
        priority: 1
    }
})(AddToWaitlist)

AddToWaitlist = connect(
    state => ({
        patient: state.discovery.patient,
        fetching: state.discovery.fetching
    }),
    {
        push,
        get,
        add
    }
)(AddToWaitlist)

export default AddToWaitlist
