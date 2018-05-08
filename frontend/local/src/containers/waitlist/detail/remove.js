import React from "react"
import { connect } from "react-redux"

import Modal from "shared/containers/modal"
import Patient from "shared/containers/patient"
import { open, COLOR_DANGER } from "shared/modules/alert"
import { listAll, remove } from "../../../modules/waitlist"

class Remove extends React.Component {
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

    handleSubmit(e) {
        e.preventDefault()

        this.props.remove(this.props.match.params.waitlistID, this.props.match.params.itemID, "canceled")
    }

    render() {
        let { item, history } = this.props
        return (
            <Modal>
                <div className="add-to-waitlist">
                    <form onSubmit={this.handleSubmit}>
                        <div className="modal-header">
                            <Patient />
                            <h1>Remove from Waiting list</h1>
                        </div>

                        {item && item.id ? (
                            <div className="modal-body">Do you really want to remove [person] from waiting list?</div>
                        ) : (
                            <div className="modal-body">Loading...</div>
                        )}

                        <div className="modal-footer">
                            <div className="form-row">
                                <div className="col-sm-4" />
                                <div className="col-sm-4">
                                    <button type="button" tabIndex="-1" className="btn btn-link btn-block" datadismiss="modal" onClick={() => history.goBack()}>
                                        No
                                    </button>
                                </div>

                                <div className="col-sm-4">
                                    <button type="submit" className="float-right btn btn-primary btn-block">
                                        Yes
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

Remove = connect(
    (state, props) => {
        return {
            listed: state.waitlist.listed,
            item: state.waitlist.items[props.match.params.itemID]
        }
    },
    {
        listAll,
        remove,
        open
    }
)(Remove)

export default Remove
