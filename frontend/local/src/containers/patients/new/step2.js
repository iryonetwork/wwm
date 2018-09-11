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
import { ReactComponent as SearchIcon } from "shared/icons/search.svg"
import { ReactComponent as SearchActiveIcon } from "shared/icons/search-active.svg"
import { ReactComponent as SpinnerIcon } from "shared/icons/spinner.svg"

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

        <div className="section">
            <div className="form-row">
                <div className="form-group col-sm-4">
                    <Field name="peopleInFamily" type="number" min="0" component={renderInput} label="No. of people in the family" />
                </div>
                <div className="form-group col-sm-4">
                    <Field name="peopleLivingTogether" type="number" min="0" component={renderInput} label="No. of people living together" />
                </div>
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
        this.state = {}
    }

    componentDidUpdate(prevProps) {
        if (!this.props.searching && prevProps.searching && this.state.searchingFamily) {
            // trigger delayed searching state change to false to prevent spinner flickering too much
            window.setTimeout(() => {
                this.setState({ searchingFamily: false })
            }, 250)
        } else if (this.props.searching !== this.state.searchingFamily) {
            this.setState({ searchingFamily: this.props.searching })
        }
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
        let { fields, getCodes } = this.props

        return (
            <div className="familyMembers">
                {fields.map((member, index) => {
                    let editDisabled = fields.get(index).patientID !== undefined
                    return (
                        <div className="section" key={fields.get(index).patientID || index}>
                            <div className="form-row searchBar">
                                <div className="col-sm-12 search">
                                    <span className="search-prepend">
                                        {this.state.searchingFamily ? <SpinnerIcon /> : this.state.searchFocused ? <SearchActiveIcon /> : <SearchIcon />}
                                    </span>
                                    <Select.Async
                                        value={editDisabled ? fields.get(index) : undefined}
                                        onChange={this.onChange(index)}
                                        onFocus={() =>
                                            this.setState({
                                                searchFocused: true
                                            })
                                        }
                                        onBlur={() =>
                                            this.setState({
                                                searchFocused: false
                                            })
                                        }
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
                                <button
                                    className="btn btn-link remove"
                                    onClick={e => {
                                        e.preventDefault()
                                        fields.remove(index)
                                    }}
                                >
                                    <RemoveIcon />
                                    Remove
                                </button>
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
                        </div>
                    )
                })}
                <div className="section">
                    <button
                        className="btn btn-link addFamilyMember"
                        onClick={e => {
                            e.preventDefault()
                            fields.push({})
                        }}
                    >
                        Add family member
                    </button>
                </div>
            </div>
        )
    }
}

familyMembers = connect(
    state => ({
        searching: state.discovery.searching || false
    }),
    {
        getCodes: getCodesAsOptions,
        loadCategories,
        search
    }
)(familyMembers)

export { Form }

export default reduxForm({
    form: "newPatient",
    destroyOnUnmount: false,
    forceUnregisterOnUnmount: true,
    validate
})(Step2)
