import React from "react"
import { Field, reduxForm } from "redux-form"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { renderInput, renderTextarea } from "shared/forms/renderField"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"

const EditComplaint = props => (
    <Modal>
        <div className="add-to-waitlist">
            <div className="modal-header">
                <Patient />
                <h1>
                    <ComplaintIcon />
                    Edit main complaint
                </h1>
            </div>

            <div className="modal-body">
                <Field name="mainComplaint" component={renderInput} label="Main complaint" />
                <Field name="mainComplaintDetails" component={renderTextarea} optional={true} rows={10} label="Details" />
            </div>

            <div className="modal-footer">
                <div className="form-row">
                    <div className="col-sm-4" />
                    <div className="col-sm-4">
                        <button type="button" tabIndex="-1" className="btn btn-link btn-block" datadismiss="modal" onClick={() => props.history.goBack()}>
                            Cancel
                        </button>
                    </div>

                    <div className="col-sm-4">
                        <button type="submit" className="float-right btn btn-primary btn-block">
                            Save
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Modal>
)

export default reduxForm({
    form: "complaint"
})(EditComplaint)
