import React, { Component } from "react"
//import PropTypes from "prop-types"
import classnames from "classnames"

import Step1 from "./step1"
import Step2 from "./step2"
import Step3 from "./step3"

import "./style.css"

class NewPatientForm extends Component {
    constructor(props) {
        super(props)
        this.nextPage = this.nextPage.bind(this)
        this.previousPage = this.previousPage.bind(this)
        this.state = {
            page: 1,
            maxPage: 1
        }
    }
    nextPage() {
        this.setState({
            page: this.state.page + 1,
            maxPage: this.state.maxPage === this.state.page ? this.state.page + 1 : this.state.maxPage
        })
    }

    previousPage() {
        this.setState({ page: this.state.page - 1 })
    }

    setPage = page => () => {
        if (page <= this.state.maxPage) {
            this.setState({ page })
        }
    }

    render() {
        const { onSubmit } = this.props
        const { page } = this.state
        return (
            <React.Fragment>
                <div className="new-patient modal fade show" style={{ display: "block", overflowY: "scroll" }} tabIndex="-1" role="dialog">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h1>Add Patient</h1>

                                <ol>
                                    <li onClick={this.setPage(1)} className={classnames({ active: page === 1 })}>
                                        Patient
                                    </li>
                                    <li onClick={this.setPage(2)} className={classnames({ active: page === 2 })}>
                                        Family Details
                                    </li>
                                    <li onClick={this.setPage(3)} className={classnames({ active: page === 3 })}>
                                        Medical History
                                    </li>
                                </ol>
                            </div>

                            <div>
                                {page === 1 && <Step1 onSubmit={this.nextPage} />}
                                {page === 2 && <Step2 previousPage={this.previousPage} onSubmit={this.nextPage} />}
                                {page === 3 && <Step3 previousPage={this.previousPage} onSubmit={onSubmit} />}
                            </div>
                        </div>
                    </div>
                </div>

                <div className="modal-backdrop fade show" />
            </React.Fragment>
        )
    }
}

NewPatientForm.propTypes = {
    //onSubmit: PropTypes.func.isRequired
}

export default NewPatientForm
