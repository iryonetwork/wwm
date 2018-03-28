import React from "react"
import { Field, FieldArray, reduxForm } from "redux-form"

import validate from "./validate"
import Footer from "./footer"
import { renderInput, renderSelect } from "./renderField"
import { documentTypeOptions } from "./options"

const genderOptions = [
    {
        label: "Male",
        value: "m"
    },
    {
        label: "Female",
        value: "f"
    }
]

const maritalStatusOptions = [
    {
        label: "Single",
        value: "single"
    },
    {
        label: "Maried",
        value: "maried"
    },
    {
        label: "Divorced",
        value: "divorced"
    },
    {
        label: "Widowed",
        value: "widowed"
    }
]

const numberOfKidsOptions = Array.from(Array(9), (x, i) => ({
    label: i,
    value: i
}))

const nationalityOptions = [
    {
        label: "Syrian",
        value: "syrian"
    }
]

const countryOptions = [
    {
        label: "Syria",
        value: "syria"
    }
]

const Step1 = props => {
    const { handleSubmit, reset } = props
    return (
        <form onSubmit={handleSubmit}>
            <div className="modal-body">
                <h3>Identification</h3>

                <Field name="image" type="file" component={renderImageField} />

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
                        <Field name="gender" component={renderSelect} options={genderOptions} label="Gender" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="maritalStatus" component={renderSelect} options={maritalStatusOptions} label="Marital status" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="numberOfKids" component={renderSelect} options={numberOfKidsOptions} label="Number of kids" />
                    </div>
                </div>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name="nationality" component={renderSelect} options={nationalityOptions} label="Nationality" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="countryOfOrigin" component={renderSelect} options={countryOptions} label="Country of origin" />
                    </div>
                </div>

                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name="education" optional={true} component={renderInput} label="Education" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="profession" optional={true} component={renderInput} label="Profession" />
                    </div>
                </div>

                <FieldArray name="documents" component={renderDocuments} />

                <h3>Contact and Demographics</h3>
                <div className="form-row">
                    <div className="form-group col-sm-4">
                        <Field name="country" component={renderSelect} options={countryOptions} label="Country" />
                    </div>
                    <div className="form-group col-sm-2">
                        <Field name="camp" component={renderInput} label="Camp" />
                    </div>
                    <div className="form-group col-sm-2">
                        <Field name="tent" component={renderInput} label="Tent" />
                    </div>
                    <div className="form-group col-sm-4">
                        <Field name="clinic" component={renderInput} label="Clinic" />
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

            <Footer reset={reset} />
        </form>
    )
}

const renderDocuments = ({ fields, meta: { error, submitFailed } }) =>
    fields.map((doc, index) => (
        <div className="form-row" key={index}>
            <div className="form-group col-sm-4">
                <Field name={`${doc}.type`} options={documentTypeOptions} component={renderSelect} label="ID document type" />
            </div>
            <div className="form-group col-sm-4">
                <Field name={`${doc}.number`} component={renderInput} label="Number" />
            </div>
            {index === fields.length - 1 && (
                <div className="form-group col-sm-4">
                    <button className="btn btn-link addDocument" onClick={() => fields.push({})}>
                        Add addidional document
                    </button>
                </div>
            )}
        </div>
    ))

const renderImageField = field => {
    let value = field.input.value
    delete field.input.value

    let image = null
    if (value && value.length) {
        let reader = new FileReader()
        reader.onload = e => {
            if (image) {
                image.style.backgroundImage = "url('" + e.target.result.replace(/(\r\n|\n|\r)/gm, "") + "')"
            }
        }
        reader.readAsDataURL(value[0])
    }

    return (
        <div className="form-row">
            <div className="form-group col-sm-2">
                <div
                    className="image"
                    ref={div => {
                        image = div
                    }}
                />
            </div>
            <div className="form-group col-sm-10">
                <input type="file" className="custom-file-input" accept="image/*" {...field.input} />
                <button type="button" className="btn btn-image btn-secondary btn-wide">
                    Add profile picture
                </button>
            </div>
        </div>
    )
}

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    initialValues: {
        documents: [{}]
    },
    validate
})(Step1)
