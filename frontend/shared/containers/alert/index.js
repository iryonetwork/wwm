import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"

import { close } from "../../modules/alert"

const Alert = props => {
    return (
        <div className={`alert alert-${props.color} alert-dismissible fade ${props.open ? "show" : ""}`} role="alert">
            {props.code ? `${props.code}:` : ""} {props.message}
            {props.disableClose ? null : (
                <button type="button" className="close" data-dismiss="alert" aria-label="Close" onClick={props.close}>
                    <span aria-hidden="true">&times;</span>
                </button>
            )}
        </div>
    )
}

const mapStateToProps = state => state.alert

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            close
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Alert)
