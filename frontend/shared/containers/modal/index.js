import React, { Component } from "react"

class Modal extends Component {
    componentDidMount() {
        document.body.classList.add("has-modal")
    }

    componentWillUnmount() {
        document.body.classList.remove("has-modal")
    }

    render() {
        return (
            <React.Fragment>
                <div className="modal fade show" tabIndex="-1" role="dialog">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">{this.props.children}</div>
                    </div>
                </div>

                <div className="modal-backdrop fade show" />
            </React.Fragment>
        )
    }
}

export default Modal
