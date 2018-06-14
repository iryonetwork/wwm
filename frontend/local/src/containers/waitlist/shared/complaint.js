import React, { Component } from "react"
import { connect } from "react-redux"
import { Field, reduxForm } from "redux-form"

import Patient from "shared/containers/patient"
import { renderInput, renderNumericalValuesRadio, renderTextarea, validateRequired } from "shared/forms/renderField"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"
import "./style.css"

const priorityOptions = [{ value: 1, label: "Yes" }, { value: 4, label: "No" }]

class ComplaintFormModalContent extends Component {
    render() {
        const { waitlistItem, patient, onSave, onClose, handleSubmit } = this.props

        return (
            <div className="complaintForm">
                <div className="modal-header">
                    <Patient data={patient} />
                    <h1>
                        {waitlistItem ? (
                            <React.Fragment>
                                <ComplaintIcon />Edit main complaint
                            </React.Fragment>
                        ) : (
                            "Add to Waiting List"
                        )}
                    </h1>
                </div>
                <form onSubmit={handleSubmit(onSave)}>
                    <div className="modal-body">
                        <div className="form-row">
                            <Field name="priority" component={renderNumericalValuesRadio} label="Urgent?" options={priorityOptions} />
                        </div>

                        <div className="form-row">
                            <Field
                                name="mainComplaint"
                                validate={validateRequired}
                                component={renderInput}
                                label="Main Complaint"
                                placeholder="Main Complaint Summary"
                            />
                        </div>

                        <div className="form-row details">
                            <Field name="mainComplaintDetails" component={renderTextarea} optional={true} rows={10} label="Main Complaint Details" />
                        </div>
                    </div>

                    <div className="modal-footer">
                        <div className="form-row">
                            <div className="col">
                                <button
                                    type="button"
                                    tabIndex="-1"
                                    className="btn btn-secondary btn-block"
                                    data-dismiss="has-modal"
                                    onClick={() => {
                                        onClose()
                                    }}
                                >
                                    Cancel
                                </button>
                            </div>
                            <div className="col">
                                <button type="submit" data-dismiss="has-modal" className="btn btn-primary btn-block">
                                    {waitlistItem ? "Save" : "Add"}
                                </button>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        )
    }
}

ComplaintFormModalContent = reduxForm({
    form: "complaintForm"
})(ComplaintFormModalContent)

ComplaintFormModalContent = connect((state, props) => {
    let initialValues
    if (props.waitlistItem) {
        initialValues = {
            priority: props.waitlistItem.priority,
            mainComplaint: props.waitlistItem.mainComplaint.complaint,
            mainComplaintDetails: props.waitlistItem.mainComplaint.comment
        }
    } else {
        initialValues = {
            priority: 4
        }
    }
    return {
        initialValues
    }
})(ComplaintFormModalContent)

class ComplaintSummary extends Component {
    componentDidMount() {
        document.body.classList.add("has-modal")
    }

    componentWillUnmount() {
        document.body.classList.remove("has-modal")
    }

    render() {
        const { waitlistItem, patient, headerMessage, onEnableEdit, onClose } = this.props
        return (
            <div className="summary">
                <div className="modal-header">
                    <Patient data={patient} big={true} />
                    <h2 className="headerMessage">{headerMessage}</h2>
                </div>
                <div className="modal-body">
                    <div className="summaryBox">
                        <div className="row header">
                            <h2>Summary</h2>
                        </div>
                        <div className="row">
                            <label htmlFor="complaint">Main complaint</label>
                            <dt>{waitlistItem.mainComplaint.complaint}</dt>
                        </div>
                        {waitlistItem.mainComplaint.comment && (
                            <div className="row">
                                <label htmlFor="complaint">Details</label>
                                <dd>{waitlistItem.mainComplaint.comment}</dd>
                            </div>
                        )}
                    </div>
                </div>
                <div className="modal-footer">
                    <div className="row">
                        {onClose && (
                            <div className="col">
                                <button
                                    type="button"
                                    tabIndex="-1"
                                    className="btn btn-secondary btn-block"
                                    data-dismiss="has-modal"
                                    onClick={() => {
                                        onClose()
                                    }}
                                >
                                    Close
                                </button>
                            </div>
                        )}
                        {onEnableEdit && (
                            <div className="col">
                                <button
                                    type="button"
                                    onClick={() => {
                                        onEnableEdit()
                                    }}
                                    data-dismiss="has-modal"
                                    className="btn btn-primary btn-block"
                                >
                                    Edit main complaint
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        )
    }
}

export { ComplaintFormModalContent, ComplaintSummary }
