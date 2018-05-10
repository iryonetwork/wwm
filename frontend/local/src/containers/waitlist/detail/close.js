import React from "react"
import { connect } from "react-redux"

import { saveConsultation } from "../../../modules/patient"
import Spinner from "shared/containers/spinner"
import Modal from "shared/containers/modal"
import { open, COLOR_DANGER } from "shared/modules/alert"

class CloseConsultation extends React.Component {
    constructor(props) {
        super(props)
        props.saveConsultation(props.match.params.waitlistID, props.match.params.itemID)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.item && nextProps.listed) {
            this.props.history.goBack()
            setTimeout(() => this.props.open("Waiting list item was not found", "", COLOR_DANGER, 5), 100)
        }
    }

    // handleSubmit(e) {
    //     e.preventDefault()

    //     this.props.remove(this.props.match.params.waitlistID, this.props.match.params.itemID, "canceled")
    // }

    render() {
        if (this.props.saving) {
            return <Spinner />
        }

        return "done..."
        // let { item, history } = this.props
        // return (
        //     <Modal>
        //         <div className="add-to-waitlist">
        //             <form onSubmit={this.handleSubmit}>
        //                 <div className="modal-header">
        //                     <Patient data={item.patient && cardToObject({ connections: item.patient })} />
        //                     <h1>Remove from Waiting list</h1>
        //                 </div>

        //                 {item && item.id ? (
        //                     <div className="modal-body">Do you really want to remove patient from waiting list?</div>
        //                 ) : (
        //                     <div className="modal-body">Loading...</div>
        //                 )}

        //                 <div className="modal-footer">
        //                     <div className="form-row">
        //                         <div className="col-sm-4" />
        //                         <div className="col-sm-4">
        //                             <button type="button" tabIndex="-1" className="btn btn-link btn-block" datadismiss="modal" onClick={() => history.goBack()}>
        //                                 No
        //                             </button>
        //                         </div>

        //                         <div className="col-sm-4">
        //                             <button type="submit" className="float-right btn btn-primary btn-block">
        //                                 Yes
        //                             </button>
        //                         </div>
        //                     </div>
        //                 </div>
        //             </form>
        //         </div>
        //     </Modal>
        // )
    }
}

CloseConsultation = connect(
    (state, props) => {
        return {
            saving: state.patient.saving,
            saved: state.patient.saved
            // listed: state.waitlist.listed,
            // item: state.waitlist.items[props.match.params.itemID]
        }
    },
    {
        saveConsultation
        // listAll
        // remove,
        // open
    }
)(CloseConsultation)

export default CloseConsultation
