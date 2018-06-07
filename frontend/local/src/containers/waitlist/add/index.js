import React, { Component } from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { push } from "react-router-redux"

import { ComplaintFormModalContent, ComplaintSummary } from "../shared/complaint"
import { get, cardToObject } from "../../../modules/discovery"
import { add, update, listAll } from "../../../modules/waitlist"
import Spinner from "shared/containers/spinner"
import Modal from "shared/containers/modal"

class AddToWaitlist extends Component {
    constructor(props) {
        super(props)
        props.get(props.match.params.patientID)
        props.listAll(props.match.params.destinationWaitlistID)

        this.edit = this.edit.bind(this)
        this.save = this.save.bind(this)
        this.close = this.close.bind(this)

        this.state = {
            edit: false,
            saving: false,
            saved: false,
            adding: false,
            added: false
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
                edit: false,
                adding: true,
                added: false
            })
            this.props
                .add(this.props.waitlistID, formData, this.props.patient)
                .then(data => {
                    this.setState({
                        adding: false,
                        added: true
                    })
                })
                .catch(ex => {
                    console.log(ex)
                })
        } else {
            let item = this.props.waitlistItem
            item.priority = formData.priority
            item.mainComplaint.complaint = formData.mainComplaint
            item.mainComplaint.comment = formData.mainComplaintDetails
            this.setState({
                edit: false,
                saving: true,
                saved: false
            })
            this.props
                .update(this.props.waitlistID, item)
                .then(data => {
                    this.props.listAll(this.props.waitlistID)
                    this.setState({
                        saving: false,
                        saved: true
                    })
                })
                .catch(ex => {
                    console.log(ex)
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
        const { waitlistFetching, patientFetching, waitlistItem } = this.props
        let patient = this.props.patient && cardToObject(this.props.patient)
        let loading = waitlistFetching || patientFetching || !patient || this.state.adding || this.state.saving

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
                            onEnableEdit={!this.state.saved && !this.state.added && this.edit}
                            onClose={this.close}
                            headerMessage={
                                (this.state.added && "Patient has been succesfully added to Waiting List") ||
                                (this.state.saved && "Main complaint has been succesfully updated") ||
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
        waitlistFetching: state.waitlist.listing
    }),
    {
        get,
        add,
        update,
        listAll,
        push
    }
)(AddToWaitlist)

export default AddToWaitlist
