import React from "react"
import { connect } from "react-redux"
import { Field, FieldArray, reduxForm } from "redux-form"
import Select from "react-select"

import Footer from "./footer"
import validate from "../shared/validatePersonal"
import { renderInput, renderSelect, renderRadio } from "shared/forms/renderField"
import { livingTogetherOptions, relationOptions } from "shared/forms/options"
import { search, cardToObject } from "../../../modules/discovery"
import { ListRow } from "../index"
import { getCodesAsOptions, loadCategories } from "shared/modules/codes"

import { ReactComponent as RemoveIcon } from "shared/icons/negative.svg"

const Step2 = props => {
    const { handleSubmit, reset, previousPage } = props
    return (
        <form onSubmit={handleSubmit}>
            <div className="modal-body">
                <Form />
            </div>

            <Footer reset={reset} previousPage={previousPage} />
        </form>
    )
}

const Form = () => (
    <div className="patient-form">
        <h3>Summary</h3>

        <div className="form-row">
            <div className="form-group col-sm-4">
                <Field name="peopleInFamily" type="number" min="0" component={renderInput} label="No. of people in the family" />
            </div>
            <div className="form-group col-sm-4">
                <Field name="peopleLivingTogether" type="number" min="0" component={renderInput} label="No. of people living together" />
            </div>
        </div>

        <h3>Family Members</h3>

        <FieldArray name="familyMembers" component={familyMembers} />
    </div>
)

class familyMembers extends React.Component {
    constructor(props) {
        super(props)

        this.getPatients = this.getPatients.bind(this)
        this.props.loadCategories("documentTypes")
    }

    onChange = index => patient => {
        let p = {}
        if (patient) {
            p = cardToObject(patient)
            p.patientID = patient.patientID

            if (p["syrian-id"]) {
                p.documentType = "syrian_id"
                p.documentNumber = p["syrian-id"]
            } else if (p["un-id"]) {
                p.documentType = "un_id"
                p.documentNumber = p["un-id"]
            }
        }

        this.props.fields.remove(index)
        //@TODO: fix this hack
        setTimeout(() => this.props.fields.insert(index, p), 100)
    }

    getPatients(input) {
        if (!input) {
            return Promise.resolve({ options: [] })
        }

        return this.props.search(input).then(data => {
            return { options: data }
        })
    }

    renderSearchLine(patient) {
        return (
            <table className="table patients">
                <tbody>
                    <ListRow patient={patient} key={patient.patientID} />
                </tbody>
            </table>
        )
    }

    renderSearchValue(patient) {
        if (!patient.patientID) {
            return null
        }
        return `${patient.lastName}, ${patient.firstName}`
    }

    render() {
        let {
            fields,
            getCodes,
        } = this.props

        return (
            <div className="familyMembers">
                {fields.map((member, index) => {
                    let editDisabled = fields.get(index).patientID !== undefined
                    return (
                        <div key={fields.get(index).patientID || index}>
                            <div className="form-row">
                                <div className="col-sm-12 search">
                                    <Select.Async
                                        value={editDisabled ? fields.get(index) : undefined}
                                        onChange={this.onChange(index)}
                                        optionRenderer={this.renderSearchLine}
                                        valueRenderer={this.renderSearchValue}
                                        filterOptions={options => options}
                                        valueKey="patientID"
                                        loadOptions={this.getPatients}
                                        multi={false}
                                        backspaceRemoves={true}
                                        placeholder="Search..."
                                    />
                                </div>
                            </div>

                            <div className="form-row">
                                <div className="form-group col-sm-4">
                                    <Field disabled={editDisabled} name={`${member}.firstName`} component={renderInput} label="First name" />
                                </div>
                                <div className="form-group col-sm-4">
                                    <Field disabled={editDisabled} name={`${member}.lastName`} component={renderInput} label="Last name" />
                                </div>
                                <div className="form-group col-sm-4">
                                    <Field disabled={editDisabled} name={`${member}.dateOfBirth`} type="date" component={renderInput} label="Date of birth" />
                                </div>
                            </div>

                            <div className="form-row">
                                <div className="form-group col-sm-4">
                                    <Field
                                        disabled={editDisabled}
                                        name={`${member}.documentType`}
                                        options={getCodes("documentTypes")}
                                        component={renderSelect}
                                        label="ID document type"
                                    />
                                </div>
                                <div className="form-group col-sm-4">
                                    <Field disabled={editDisabled} name={`${member}.documentNumber`} component={renderInput} label="Number" />
                                </div>
                            </div>

                            <div className="form-row">
                                <div className="form-group col-sm-4">
                                    <Field name={`${member}.relation`} component={renderSelect} options={relationOptions} label="Relation" />
                                </div>
                                <div className="form-group col-sm-8">
                                    <Field name={`${member}.livingTogether`} options={livingTogetherOptions} component={renderRadio} label="Living together?" />
                                </div>
                            </div>

                            <a
                                onClick={e => {
                                    e.preventDefault()
                                    fields.remove(index)
                                }}
                                href="/"
                                className="remove"
                            >
                                <RemoveIcon />
                                Remove
                            </a>
                        </div>
                    )
                })}
                <div className="link">
                    <a
                        href="/"
                        onClick={e => {
                            e.preventDefault()
                            fields.push({})
                        }}
                    >
                        Add family member
                    </a>
                </div>
            </div>
        )
    }
}

familyMembers = connect(state => ({}), {
    getCodes: getCodesAsOptions,
    loadCategories,
    search
})(familyMembers)

export { Form }

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step2)
