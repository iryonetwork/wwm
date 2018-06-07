import React from "react"
import { connect } from "react-redux"
import { withRouter } from "react-router"
import { push } from "react-router-redux"

import { ComplaintFormModalContent, ComplaintSummary } from "../shared/complaint"
import Modal from "shared/containers/modal"
import Spinner from "shared/containers/spinner"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { update, listAll, resetIndicators } from "../../../modules/waitlist"
import { cardToObject } from "../../../modules/discovery"

class EditComplaint extends React.Component {
    constructor(props) {
        super(props)
        props.resetIndicators()
        if (!props.item) {
            props.listAll(props.match.params.waitlistID)
        }

        this.save = this.save.bind(this)
        this.close = this.close.bind(this)

        this.state = {
            saving: false,
            saved: false
        }
    }

    componentDidMount() {
        document.body.classList.add("has-modal")
    }

    componentWillUnmount() {
        document.body.classList.remove("has-modal")
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.history.goBack()
            setTimeout(() => this.props.open("Waitlist item was not found", "", COLOR_DANGER, 5), 100)
        }
    }

    save(formData) {
        let item = this.props.item
        item.priority = formData.priority
        item.mainComplaint.complaint = formData.mainComplaint
        item.mainComplaint.comment = formData.mainComplaintDetails

        this.props.update(this.props.waitlistID, item)
    }

    close = () => {
        this.props.push(`/waitlist/${this.props.match.params.waitlistID}`)
    }

    render() {
        let { item, waitlistUpdating, waitlistUpdated } = this.props
        let loading = !item || waitlistUpdating
        let patient = item && item.patient && cardToObject({ connections: item.patient })

        return (
            <Modal>
                <div className="add-to-waitlist">
                    {loading ? (
                        <div className="modal-body">
                            <Spinner />
                        </div>
                    ) : !waitlistUpdated ? (
                        <ComplaintFormModalContent waitlistItem={item} patient={patient} onSave={this.save} onClose={this.close} />
                    ) : (
                        <ComplaintSummary
                            waitlistItem={item}
                            patient={patient}
                            onClose={this.close}
                            headerMessage="Main complaint has been succesfully updated"
                        />
                    )}
                </div>
            </Modal>
        )
    }
}

EditComplaint = connect(
    (state, props) => {
        return {
            waitlistID: props.match.params.waitlistID,
            listed: state.waitlist.listed,
            item: state.waitlist.items[props.match.params.itemID],
            waitlistUpdating: state.waitlist.updating,
            waitlistUpdated: state.waitlist.updated
        }
    },
    {
        update,
        listAll,
        open,
        resetIndicators,
        push
    }
)(EditComplaint)

export default withRouter(EditComplaint)
