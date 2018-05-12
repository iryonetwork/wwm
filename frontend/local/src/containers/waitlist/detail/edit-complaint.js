import React from "react"
import { Field, reduxForm } from "redux-form"
import { connect } from "react-redux"
import { withRouter } from "react-router"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { renderInput, renderTextarea } from "shared/forms/renderField"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { listAll, update } from "../../../modules/waitlist"
import { cardToObject } from "../../../modules/discovery"

import { ReactComponent as ComplaintIcon } from "shared/icons/complaint.svg"

class EditComplaint extends React.Component {
    constructor(props) {
        super(props)
        if (!props.item) {
            props.listAll(props.match.params.waitlistID)
        }

        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.history.goBack()
            setTimeout(() => this.props.open("Waitlist item was not found", "", COLOR_DANGER, 5), 100)
        }
    }

    handleSubmit(form) {
        let item = this.props.item
        item.mainComplaint.complaint = form.mainComplaint
        item.mainComplaint.comment = form.mainComplaintDetails

        this.props.update(this.props.match.params.waitlistID, item)
    }

    render() {
        let { item, history, handleSubmit } = this.props
        return (
            <Modal>
                <div className="add-to-waitlist">
                    <form onSubmit={handleSubmit(this.handleSubmit)}>
                        <div className="modal-header">
                            <Patient data={item.patient && cardToObject({ connections: item.patient })} />
                            <h1>
                                <ComplaintIcon />
                                Edit main complaint
                            </h1>
                        </div>

                        {item && item.id ? (
                            <div className="modal-body">
                                <Field name="mainComplaint" component={renderInput} label="Main complaint" />
                                <Field name="mainComplaintDetails" component={renderTextarea} optional={true} rows={10} label="Details" />
                            </div>
                        ) : (
                            <div className="modal-body">Loading...</div>
                        )}

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
                </div>
            </Modal>
        )
    }
}

EditComplaint = reduxForm({
    form: "complaint"
})(EditComplaint)

EditComplaint = connect(
    (state, props) => {
        let item = state.waitlist.items[props.match.params.itemID]
        let initialValues
        if (item) {
            initialValues = {
                mainComplaint: item.mainComplaint.complaint,
                mainComplaintDetails: item.mainComplaint.comment
            }
        }

        return {
            listed: state.waitlist.listed,
            item,
            initialValues
        }
    },
    {
        listAll,
        update,
        open
    }
)(EditComplaint)

export default withRouter(EditComplaint)
