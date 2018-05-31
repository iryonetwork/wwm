import React from "react"
import { connect } from "react-redux"
import { Field, FieldArray, reduxForm } from "redux-form"
import { push } from "react-router-redux"
import _ from "lodash"

import { searchCodes } from "shared/modules/codes"
import { update as updateWaitlistItem } from "../../../modules/waitlist"
import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { renderInput, renderTextarea, renderReactSelect } from "shared/forms/renderField"
import { cardToObject } from "../../../modules/discovery"

import { ReactComponent as DiagnosisIcon } from "shared/icons/diagnosis.svg"

import "react-select/dist/react-select.css"

const diagnosisRequired = value => ( !_.isEmpty(value) ? undefined : 'Required')

class AddDiagnosis extends React.Component {
    constructor(props) {
        super(props)
        this.fetchCodes = this.fetchCodes.bind(this)
        this.onSubmit = this.onSubmit.bind(this)
    }

    fetchCodes(input) {
        if (!input) {
            return Promise.resolve({ options: [] })
        }

        return this.props.searchCodes("diagnosis", input).then(results => ({ options: results.map(el => ({ id: el.id, label: el.title })) }))
    }

    onSubmit(formData) {
        // add it to waitlist item
        let diagnosis = {
            diagnosis: formData.diagnosis.id,
            label: formData.diagnosis.label,
            comment: formData.comment,
            therapies: formData.therapies,
        }


        let newItem = Object.assign({}, this.props.waitlistItem)
        newItem.diagnoses = newItem.diagnoses || []
        if (this.props.diagnosisIndex) {
            newItem.diagnoses[this.props.diagnosisIndex] = diagnosis
        } else {
            newItem.diagnoses.push(diagnosis)
        }

        this.props.updateWaitlistItem(this.props.match.params.waitlistID, newItem)

        let {waitlistID, itemID} = this.props.match.params
        this.props.push(`/waitlist/${waitlistID}/${itemID}/consultation`)
    }

    render() {
        const { history, handleSubmit } = this.props
        return !this.props.loading && this.props.waitlistItem ? (
            <Modal>
                <form className="add-diagnosis" onSubmit={handleSubmit(this.onSubmit)}>
                    <div className="modal-header">
                        <Patient data={this.props.waitlistItem.patient && cardToObject({ connections: this.props.waitlistItem.patient })} />
                        <h1>
                            <DiagnosisIcon />
                            {this.props.diagnosisIndex ? "Edit diagnosis" : "Add diagnosis"}
                        </h1>
                    </div>

                    <div className="modal-body">
                        <div className="form-row diagnosisSelect">
                            <div className="form-group col-sm-12">
                                <Field name="diagnosis" validate={diagnosisRequired} component={renderReactSelect} label="Diagnosis" loadOptions={(value) => this.fetchCodes( value ? value : this.props.initialValues.diagnosis.id)} />
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group col-sm-12">
                                <Field name="comment" component={renderTextarea} label="Description" />
                            </div>
                        </div>

                        <h2>Therapy</h2>

                        <FieldArray name="therapies" component={renderTherapies} />
                    </div>

                    <div className="modal-footer">
                        <div className="form-row">
                            <div className="col-sm-4" />
                            <div className="col-sm-4">
                                <button type="button" tabIndex="-1" className="btn btn-link btn-block" datadismiss="modal" onClick={() => history.goBack()}>
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
            </Modal>
        ) : (null)
    }
}

AddDiagnosis = reduxForm({
    form: "diagnosis"
})(AddDiagnosis)

AddDiagnosis = connect(
    (state, props) => {
        let loading = state.waitlist.listing || state.waitlist.fetching || state.waitlist.items[props.match.params.itemID].updating
        let item = state.waitlist.items[props.match.params.itemID]

        let initialValues = {diagnosis: {}}
        if (!loading && props.match.params.diagnosisIndex) {
            initialValues = {
                diagnosis: {
                    id: item.diagnoses[props.match.params.diagnosisIndex].diagnosis,
                    label: item.diagnoses[props.match.params.diagnosisIndex].label,
                },
                comment: item.diagnoses[props.match.params.diagnosisIndex].comment,
                therapies: item.diagnoses[props.match.params.diagnosisIndex].therapies,
            }
        }

        return ({
            loading: loading,
            diagnosisIndex: props.match.params.diagnosisIndex,
            waitlistItem: item,
            initialValues: initialValues,
            searchingCodes: state.codes.searching,
            searchingResults: state.codes.searchResults,
        })
    },
    {
        searchCodes,
        updateWaitlistItem,
        push
    }
)(AddDiagnosis)

export default AddDiagnosis


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
            <div className="section">
                {(fields || []).map((therapy, index) => (
                    <React.Fragment key={index}>
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
                    </React.Fragment>
                ))}
                <div className="form-row">
                    <div className="form-group">
                        <button className="btn btn-link addTherapy" onClick={this.pushNewFields}>
                            Add therapy
                        </button>
                    </div>
                </div>
            </div>
        )
    }
}
