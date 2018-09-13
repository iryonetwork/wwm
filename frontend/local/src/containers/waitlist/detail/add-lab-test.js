import React from "react"
import { Field, reduxForm } from "redux-form"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient/card"
import { renderRadio } from "shared/forms/renderField"

import { ReactComponent as LaboratoryIcon } from "shared/icons/laboratory.svg"

const options = [
    {
        label: "Positive",
        value: true
    },
    {
        label: "Negative",
        value: false
    }
]

const LaboratoryTest = props => (
    <Modal>
        <div className="lab-test">
            <div className="modal-header">
                <Patient />
                <h1>
                    <LaboratoryIcon />
                    Add Laboratory Tests
                </h1>
            </div>

            <div className="modal-body">
                <Field name="infections" options={options} component={renderRadio} label="Bladder or kidney infections?" />

                <Field name="pregnancy" options={options} component={renderRadio} label="Pregnancy?" />

                <Field name="preeclampsia" options={options} component={renderRadio} label="Preeclampsia?" />

                <Field name="dehydration" options={options} component={renderRadio} label="Dehydration?" />
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
                            Add
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Modal>
)

export default reduxForm({
    form: "lab-test"
})(LaboratoryTest)
