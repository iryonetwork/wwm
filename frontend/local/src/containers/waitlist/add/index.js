import React, { Component } from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { push } from "react-router-redux"

import { ComplaintFormModalContent, ComplaintSummary } from "../shared/complaint"
import { get, cardToObject } from "../../../modules/discovery"
import { add, update, listAll, resetIndicators } from "../../../modules/waitlist"
import Spinner from "shared/containers/spinner"
import Modal from "shared/containers/modal"

class AddToWaitlist extends Component {
    constructor(props) {
        super(props)
        props.resetIndicators()
        props.get(props.match.params.patientID)
        props.listAll(props.match.params.destinationWaitlistID)

        this.edit = this.edit.bind(this)
        this.save = this.save.bind(this)
        this.close = this.close.bind(this)

        this.state = {
            edit: false
        }
    }

    componentDidMount() {
        document.body.classList.add("has-modal")
    }

    componentWillUnmount() {
        document.body.classList.remove("has-modal")
    }

    save(formData) {
        if (!this.props.waitlistItem) {
            this.setState({
                edit: false
            })
            this.props.add(this.props.waitlistID, formData, this.props.patient)
        } else {
            this.setState({
                edit: false
            })
            let item = this.props.waitlistItem
            item.priority = formData.priority
            item.mainComplaint.complaint = formData.mainComplaint
            item.mainComplaint.comment = formData.mainComplaintDetails

            this.props.update(this.props.waitlistID, item).then(data => {
                this.props.listAll(this.props.waitlistID)
            })
        }
    }

    edit = () => {
        this.setState({ edit: true })
    }

    close = () => {
        this.props.push("/")
    }

    render() {
        const { waitlistFetching, patientFetching, waitlistItem, waitlistAdding, waitlistAdded, waitlistUpdating, waitlistUpdated } = this.props
        let patient = this.props.patient && cardToObject(this.props.patient)
        let loading = waitlistFetching || patientFetching || !patient || waitlistAdding || waitlistUpdating

        return (
            <Modal>
                <div className="add-to-waitlist">
                    {loading ? (
                        <div className="modal-body">
                            <Spinner />
                        </div>
                    ) : !waitlistItem || this.state.edit ? (
                        <ComplaintFormModalContent waitlistItem={waitlistItem} patient={patient} onSave={this.save} onClose={this.close} />
                    ) : (
                        <ComplaintSummary
                            waitlistItem={waitlistItem}
                            patient={patient}
                            onEnableEdit={!waitlistUpdated && !waitlistAdded && this.edit}
                            onClose={this.close}
                            headerMessage={
                                (waitlistAdded && "Patient has been succesfully added to Waiting List") ||
                                (waitlistUpdated && "Main complaint has been succesfully updated") ||
                                "Patient is already in the Waiting List"
                            }
                        />
                    )}
                </div>
            </Modal>
        )
    }
}

AddToWaitlist = connect(
    (state, props) => ({
        waitlistID: props.match.params.destinationWaitlistID,
        patient: state.discovery.patient,
        patientFetching: state.discovery.fetching,
        waitlistItem: _.find(state.waitlist.items, ["patientID", props.match.params.patientID]),
        waitlistFetching: state.waitlist.listing,
        waitlistAdding: state.waitlist.adding,
        waitlistAdded: state.waitlist.added,
        waitlistUpdating: state.waitlist.updating,
        waitlistUpdated: state.waitlist.updated
    }),
    {
        get,
        add,
        update,
        listAll,
        resetIndicators,
        push
    }
)(AddToWaitlist)

export default AddToWaitlist
