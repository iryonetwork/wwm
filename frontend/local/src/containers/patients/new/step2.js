import React from "react"
import { Field, FieldArray, reduxForm } from "redux-form"

import Footer from "./footer"
import validate from "./validate"
import { renderInput, renderSelect, renderRadio } from "shared/forms/renderField"
import { documentTypeOptions, yesNoOptions } from "./options"

import { ReactComponent as RemoveIcon } from "shared/icons/negative.svg"

const Step2 = props => {
    const { handleSubmit, reset, previousPage } = props
    return (
        <form onSubmit={handleSubmit}>
            <div className="modal-body">
                <h3>Summary</h3>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name="people_in_family" type="number" min="0" component={renderInput} label="No. of people in the family" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="people_living_together" type="number" min="0" component={renderInput} label="No. of people living together" />
                    </div>
                </div>

                <h3>Family Members</h3>

                <FieldArray name="familyMembers" component={familyMembers} />
            </div>

            <Footer reset={reset} previousPage={previousPage} />
        </form>
    )
}

const relationOptions = [
    {
        value: "spouse",
        label: "Spouse"
    },
    {
        value: "child",
        label: "Child"
    },
    {
        value: "parent",
        label: "Parent"
    },
    {
        value: "sibling",
        label: "Sibling"
    },
    {
        value: "grandparent",
        label: "Grandparent"
    },
    {
        value: "other",
        label: "Other"
    }
]

const familyMembers = ({ fields, meta: { error, submitFailed } }) => (
    <div className="familyMembers">
        {fields.map((member, index) => (
            <div key={member}>
                <div className="form-row">
                    <div className="col-sm-12 patients">
                        <input name="search" placeholder="Search" className="search" />
                    </div>
                </div>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.firstName`} component={renderInput} label="First name" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.lastName`} component={renderInput} label="Last name" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.dateOfBirth`} type="date" component={renderInput} label="Date of birth" />
                    </div>
                </div>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.relation`} component={renderSelect} options={relationOptions} label="Relation" />
                    </div>
                    <div className="form-group col-sm-8">
                        <Field name={`${member}.livingTogether`} options={yesNoOptions} component={renderRadio} label="Living together?" />
                    </div>
                </div>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.documentType`} options={documentTypeOptions} component={renderSelect} label="ID document type" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name={`${member}.documentNumber`} component={renderInput} label="Number" />
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
        ))}
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

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step2)
