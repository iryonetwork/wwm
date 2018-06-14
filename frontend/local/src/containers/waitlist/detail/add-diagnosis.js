import React from "react"
import { connect } from "react-redux"
import _ from "lodash"
import { Field, FieldArray, reduxForm } from "redux-form"
import { push, goBack } from "react-router-redux"

import { searchCodes } from "shared/modules/codes"
import { update as updateWaitlistItem, resetIndicators } from "../../../modules/waitlist"
import { saveConsultation } from "../../../modules/patient"
import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { renderInput, renderTextarea, renderReactSelect, validateRequired } from "shared/forms/renderField"
import { cardToObject } from "../../../modules/discovery"
import Spinner from "shared/containers/spinner"
import { open, COLOR_SUCCESS } from "shared/modules/alert"

import { ReactComponent as DiagnosisIcon } from "shared/icons/diagnosis.svg"

import "react-select/dist/react-select.css"

class AddDiagnosis extends React.Component {
    constructor(props) {
        super(props)
        props.resetIndicators()
        if (!props.waitlistItem) {
            props.listAll(props.match.params.waitlistID)
        } else if (!props.match.params.diagnosisIndex) {
            props.push(
                `/waitlist/${this.props.match.params.waitlistID}/${this.props.match.params.itemID}/consultation/diagnoses/${
                    props.waitlistItem.diagnoses ? props.waitlistItem.diagnoses.length : 0
                }/edit`
            )
        }

        this.onSave = this.onSave.bind(this)
        this.onClose = this.onClose.bind(this)
        this.onCloseConsultation = this.onCloseConsultation.bind(this)
    }

    onSave(formData) {
        // add it to waitlist item
        let diagnosis = {
            diagnosis: formData.diagnosis.id,
            label: formData.diagnosis.label,
            comment: formData.comment,
            therapies: formData.therapies
        }

        let newItem = Object.assign({}, this.props.waitlistItem)
        newItem.diagnoses = newItem.diagnoses || []
        if (this.props.diagnosisIndex) {
            newItem.diagnoses[this.props.diagnosisIndex] = diagnosis
        } else {
            newItem.diagnoses.push(diagnosis)
        }

        this.props.updateWaitlistItem(this.props.match.params.waitlistID, newItem)
    }

    onClose() {
        this.props.push(`/waitlist/${this.props.match.params.waitlistID}/${this.props.match.params.itemID}/consultation`)
    }

    onCloseConsultation() {
        this.props.saveConsultation(this.props.match.params.waitlistID, this.props.match.params.itemID).then(() => {
            this.props.push(`/waitlist/${this.props.match.params.waitlistID}`)
            this.props.open("Consultation was closed", "", COLOR_SUCCESS, 5)
        })
    }

    render() {
        const { waitlistLoading, waitlistUpdated, waitlistItem, diagnosis } = this.props
        let loading = !waitlistItem || waitlistLoading
        return !loading ? (
            !waitlistUpdated ? (
                <Modal>
                    <DiagnosisFormModalContent
                        diagnosis={diagnosis}
                        patient={this.props.waitlistItem.patient && cardToObject({ connections: this.props.waitlistItem.patient })}
                        onSave={this.onSave}
                        onCancel={this.onClose}
                    />
                </Modal>
            ) : (
                <Modal>
                    <DiagnosisSummary
                        diagnosis={diagnosis}
                        patient={this.props.waitlistItem.patient && cardToObject({ connections: this.props.waitlistItem.patient })}
                        onCloseConsultation={this.onCloseConsultation}
                        onContinueConsultation={this.onClose}
                    />
                </Modal>
            )
        ) : (
            <Modal>
                <Spinner />
            </Modal>
        )
    }
}

AddDiagnosis = connect(
    (state, props) => {
        let loading = state.waitlist.listing || state.waitlist.fetching || state.waitlist.items[props.match.params.itemID].updating
        let item = state.waitlist.items[props.match.params.itemID]

        let diagnosis =
            !loading && props.match.params.diagnosisIndex && item.diagnoses && item.diagnoses[props.match.params.diagnosisIndex]
                ? {
                      diagnosis: {
                          id: item.diagnoses[props.match.params.diagnosisIndex].diagnosis,
                          label: item.diagnoses[props.match.params.diagnosisIndex].label
                      },
                      comment: item.diagnoses[props.match.params.diagnosisIndex].comment,
                      therapies: item.diagnoses[props.match.params.diagnosisIndex].therapies
                  }
                : { diagnosis: {} }

        return {
            waitlistLoading: loading,
            waitlistUpdated: state.waitlist.updated,
            diagnosisIndex: props.match.params.diagnosisIndex,
            waitlistItem: item,
            diagnosis: diagnosis
        }
    },
    {
        updateWaitlistItem,
        push,
        resetIndicators,
        goBack,
        open,
        saveConsultation
    }
)(AddDiagnosis)

export default AddDiagnosis

class DiagnosisFormModalContent extends React.Component {
    constructor(props) {
        super(props)
        this.fetchCodes = this.fetchCodes.bind(this)
    }

    fetchCodes(input) {
        if (!input) {
            return Promise.resolve({ options: [] })
        }

        return this.props.searchCodes("diagnosis", input).then(results => ({ options: results.map(el => ({ id: el.id, label: el.title })) }))
    }

    render() {
        const { patient, onSave, onCancel, handleSubmit } = this.props
        return (
            <form className="add-diagnosis" onSubmit={handleSubmit(onSave)}>
                <div className="modal-header">
                    <Patient data={patient} />
                    <h1>
                        <DiagnosisIcon />
                        {!_.isEmpty(this.props.diagnosis.diagnosis) ? "Edit diagnosis" : "Add diagnosis"}
                    </h1>
                </div>

                <div className="modal-body">
                    <div className="form-row diagnosisSelect">
                        <div className="form-group col-sm-12">
                            <Field
                                name="diagnosis"
                                validate={validateRequired}
                                component={renderReactSelect}
                                label="Diagnosis"
                                loadOptions={value => this.fetchCodes(value ? value : this.props.initialValues.diagnosis.id)}
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-group col-sm-12">
                            <Field name="comment" component={renderTextarea} label="Description" />
                        </div>
                    </div>
                    <div className="therapies">
                        <h2>Therapies</h2>
                        <FieldArray name="therapies" component={renderTherapies} />
                    </div>
                </div>

                <div className="modal-footer">
                    <div className="form-row">
                        <div className="col-sm-4" />
                        <div className="col-sm-4">
                            <button type="button" tabIndex="-1" className="btn btn-secondary btn-block" datadismiss="modal" onClick={() => onCancel()}>
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
            </form>
        )
    }
}

DiagnosisFormModalContent = reduxForm({
    form: "diagnosis"
})(DiagnosisFormModalContent)

DiagnosisFormModalContent = connect(
    (state, props) => ({
        initialValues: props.diagnosis
    }),
    {
        searchCodes
    }
)(DiagnosisFormModalContent)

class DiagnosisSummary extends React.Component {
    constructor(props) {
        super(props)
        this.fetchCodes = this.fetchCodes.bind(this)
    }

    fetchCodes(input) {
        if (!input) {
            return Promise.resolve({ options: [] })
        }

        return this.props.searchCodes("diagnosis", input).then(results => ({ options: results.map(el => ({ id: el.id, label: el.title })) }))
    }

    render() {
        const { patient, diagnosis, onContinueConsultation, onCloseConsultation } = this.props

        return (
            <div className="summary">
                <div className="modal-header">
                    <Patient data={patient} big={true} />
                    <h2 className="headerMessage">Diagnosis has been succesfully saved</h2>
                </div>
                <div className="modal-body">
                    <div className="summaryBox">
                        <div className="row header">
                            <h2>Summary</h2>
                        </div>
                        <div className="row">
                            <label htmlFor="diagnosis">Diagnosis</label>
                            <dl id="diagnosis">
                                <dt>{diagnosis.diagnosis.label}</dt>
                                {diagnosis.comment && <dd>{diagnosis.comment}</dd>}
                            </dl>
                        </div>
                        {diagnosis.therapies &&
                            diagnosis.therapies.length > 0 && (
                                <div className="row">
                                    <label htmlFor="therapies">Therapy</label>
                                    <dl id="therapies">
                                        {diagnosis.therapies.map(el => (
                                            <React.Fragment>
                                                <dt>{el.medicine}</dt>
                                                {el.instructions && <dd>{el.instructions}</dd>}
                                            </React.Fragment>
                                        ))}
                                    </dl>
                                </div>
                            )}
                    </div>
                </div>
                <div className="modal-footer">
                    <div className="row">
                        {onContinueConsultation && (
                            <div className="col">
                                <button
                                    type="button"
                                    tabIndex="-1"
                                    className="btn btn-primary btn-block"
                                    data-dismiss="has-modal"
                                    onClick={() => {
                                        onContinueConsultation()
                                    }}
                                >
                                    Continue consultation
                                </button>
                            </div>
                        )}
                        {onCloseConsultation && (
                            <div className="col">
                                <button
                                    type="button"
                                    onClick={() => {
                                        onCloseConsultation()
                                    }}
                                    data-dismiss="has-modal"
                                    className="btn btn-success btn-block"
                                >
                                    Close consultation
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        )
    }
}

class renderTherapies extends React.Component {
    constructor(props) {
        super(props)
        this.pushNewFields = this.pushNewFields.bind(this)
    }

    pushNewFields(e) {
        e.preventDefault()
        this.props.fields.push({})
    }

    render() {
        const { fields } = this.props
        return (
            <React.Fragment>
                {(fields || []).map((therapy, index) => (
                    <div className="section" key={index}>
                        <div className="form-row">
                            <div className="form-group col-sm-12">
                                <Field name={`${therapy}.medicine`} component={renderInput} label="Medicine" />
                            </div>
                        </div>
                        <div className="form-row">
                            <div className="form-group col-sm-12">
                                <Field name={`${therapy}.instructions`} component={renderTextarea} label="Instructions" />
                            </div>
                        </div>
                    </div>
                ))}
                <div className="section">
                    <button className="btn btn-link addTherapy" onClick={this.pushNewFields}>
                        Add therapy
                    </button>
                </div>
            </React.Fragment>
        )
    }
}
