import React, { Component } from "react"
import { withRouter } from "react-router-dom"
import { Field, reduxForm } from "redux-form"
//import PropTypes from "prop-types"
//import classnames from "classnames"

import { renderInput, renderRadio, renderTextarea, renderSelect } from "shared/forms/renderField"
import { yesNoOptions } from "shared/forms/options"
import Patient from "shared/containers/patient"
import "./style.css"

const doctorOptions = [
    {
        label: "Dr. Doctor",
        value: "uuid-of-the-doctor"
    }
]

class AddToWaitlist extends Component {
    constructor(props) {
        super(props)

        this.state = {}
    }

    componentDidMount() {
        document.body.style.overflow = "hidden"
    }

    componentWillUnmount() {
        document.body.style.overflow = "auto"
    }

    render() {
        let { history } = this.props
        return (
            <React.Fragment>
                <div className="add-to-waitlist modal fade show" tabIndex="-1" role="dialog">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">
                            <div className="modal-header">
                                <Patient />
                                <h1>Add to Waiting List</h1>
                            </div>
                            <form>
                                <div className="modal-body">
                                    <div className="form-row">
                                        <Field name="urgent" component={renderRadio} label="Urgent?" options={yesNoOptions} />
                                    </div>

                                    <div className="form-row">
                                        <Field name="mainComplaint" component={renderInput} label="Main complaint" />
                                    </div>

                                    <div className="form-row details">
                                        <Field name="mainComplaintDetails" component={renderTextarea} optional={true} label="Details" />
                                    </div>

                                    <div className="form-row">
                                        <Field name="doctor" component={renderSelect} options={doctorOptions} label="Doctor" />
                                    </div>
                                </div>

                                <div className="modal-footer">
                                    <div className="form-row">
                                        <div className="col-sm-4">
                                            <button
                                                type="button"
                                                tabIndex="-1"
                                                className="btn btn-link btn-block"
                                                datadismiss="modal"
                                                onClick={() => {
                                                    history.push("/")
                                                }}
                                            >
                                                Cancel
                                            </button>
                                        </div>
                                        <div className="col-sm-4" />
                                        <div className="col-sm-4">
                                            <button type="submit" className="float-right btn btn-primary btn-block">
                                                Add
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>

                <div className="modal-backdrop fade show" />
            </React.Fragment>
        )
    }
}

export default withRouter(
    reduxForm({
        form: "addToWaitlist"
    })(AddToWaitlist)
)
