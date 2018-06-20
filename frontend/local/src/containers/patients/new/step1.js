import React, { Component } from "react"
import { connect } from "react-redux"
import { Field, FieldArray, reduxForm } from "redux-form"

import validate from "../shared/validatePersonal"
import Footer from "./footer"
import Spinner from "shared/containers/spinner"
import { renderInput, renderSelect } from "shared/forms/renderField"
import { getCodesAsOptions, loadCategories as loadCategoriesImport } from "shared/modules/codes"

const numberOfKidsOptions = Array.from(Array(9), (x, i) => ({
    label: i,
    value: i
}))

const educationOptions = [
    {
        label: "Primary",
        value: "primary"
    },
    {
        label: "Secondary",
        value: "secondary"
    },
    {
        label: "Tertiary",
        value: "tertiary"
    }
]

const Form = props => (
    <div className="patient-form">
        <h3>Identification</h3>
        <div className="section">
            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="firstName" component={renderInput} label="First name" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="lastName" component={renderInput} label="Last name" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="dateOfBirth" component={renderInput} type="date" label="Date of birth" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="gender" component={renderSelect} options={props.genders} label="Gender" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="maritalStatus" component={renderSelect} options={props.maritalStatus} label="Marital status" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="numberOfKids" component={renderSelect} options={numberOfKidsOptions} label="Number of kids" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="nationality" component={renderSelect} options={props.countries} label="Nationality" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="countryOfOrigin" component={renderSelect} options={props.countries} label="Country of origin" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="education" optional={true} component={renderSelect} options={educationOptions} label="Education" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="profession" optional={true} component={renderInput} label="Profession" />
                </div>
            </div>

            <FieldArray name="documents" component={renderDocuments} documentTypes={props.documentTypes} />
        </div>

        <h3>Contact and Demographics</h3>
        <div className="section">
            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="country" component={renderSelect} options={props.countries} label="Country" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="region" component={renderInput} label="Region" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="address" component={renderInput} label="Address" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="phone" optional={true} component={renderInput} label="Phone number" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="email" optional={true} component={renderInput} label="Email address" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="whatsapp" optional={true} component={renderInput} label="WhatsApp" />
                </div>
            </div>

            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="dateOfLeaving" component={renderInput} type="date" label="Date of leaving home country" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="transitCountries" component={renderInput} label="Transit countries" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="dateOfArrival" component={renderInput} type="date" label="Date of arrival" />
                </div>
            </div>
        </div>
    </div>
)

class Step1 extends Component {
    constructor(props) {
        super(props)
        props.loadCategories("gender", "maritalStatus", "countries", "documentTypes")
    }

    render() {
        const { handleSubmit, reset, codesLoading, getCodes } = this.props

        if (codesLoading) {
            return <Spinner />
        }

        return (
            <form onSubmit={handleSubmit}>
                <div className="modal-body">
                    <Form
                        countries={getCodes("countries")}
                        maritalStatus={getCodes("maritalStatus")}
                        genders={getCodes("gender")}
                        documentTypes={getCodes("documentTypes")}
                    />
                </div>

                <Footer reset={reset} />
            </form>
        )
    }
}

const renderDocuments = props => {
    const { fields, documentTypes } = props
    return fields.map((doc, index) => (
        <div className="form-row" key={index}>
            <div className="form-group col-sm-4">
                <Field name={`${doc}.type`} options={documentTypes} component={renderSelect} label="ID document type" />
            </div>
            <div className="form-group col-sm-4">
                <Field name={`${doc}.number`} component={renderInput} label="Number" />
            </div>
            {index === fields.length - 1 && (
                <div className="form-group col-sm-4">
                    <button className="btn btn-link addDocument" onClick={() => fields.push({})}>
                        Add additional document
                    </button>
                </div>
            )}
        </div>
    ))
}

Step1 = reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    // initialValues: {
    //     documents: [{}]
    // },
    validate
})(Step1)

Step1 = connect(
    state => ({
        codesLoading: state.codes.loading,
        initialValues: state.patient.newData
    }),
    {
        getCodes: getCodesAsOptions,
        loadCategories: loadCategoriesImport
    }
)(Step1)

export default Step1

export { Form }
