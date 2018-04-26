import React from "react"
import { Field, reduxForm } from "redux-form"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { renderInput, renderTextarea } from "shared/forms/renderField"

import { ReactComponent as MedicalHistoryIcon } from "shared/icons/medical-history-active.svg"

const AddDiagnosis = props => (
    <Modal>
        <div className="">
            <div className="modal-header">
                <Patient />
                <h1>
                    <MedicalHistoryIcon />
                    Add diagnosis
                </h1>
            </div>

            <div className="modal-body">dia gnosis!</div>

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
    form: "diagnosis"
})(AddDiagnosis)
