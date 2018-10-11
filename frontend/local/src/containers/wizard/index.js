import React, { Component } from "react"
import { get } from "lodash"
import { connect } from "react-redux"
import { goBack } from "react-router-redux"

import { save } from "shared/modules/config"
import { load as loadUser } from "../../modules/users"
import Spinner from "shared/containers/spinner"
import SettingsContent from "../settings"

import "./style.css"

class Wizard extends Component {
    constructor(props) {
        super(props)

        if (props.user === undefined) {
            props.loadUser("me")
        }
    }

    componentDidMount() {
        document.body.classList.add("has-modal")
    }

    componentWillUnmount() {
        document.body.classList.remove("has-modal")
    }

    close() {
        return e => {
            this.props.save("wizardWasShown", true)
        }
    }

    render() {
        if (this.props.loading) {
            return <Spinner />
        }

        let firstName = get(this.props.user, "personalData.firstName")
        let lastName = get(this.props.user, "personalData.lastName")
        let name = firstName && lastName ? `${firstName} ${lastName}` : this.props.user.username

        return (
            <React.Fragment>
                <div className="wizard modal fade show" tabIndex="-1" role="dialog">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h1>Hello {name}!</h1>
                            </div>

                            <div className="modal-body">
                                <p>It appears you have logged in to Iryo Clinic for the first time on this machine.</p>
                                <p>Please review default interface settings and change them if needed.</p>
                                <div className="settings">
                                    <SettingsContent />
                                </div>
                            </div>
                            <div className="modal-footer">
                                <div className="form-row">
                                    <div className="col-sm">
                                        <button
                                            type="button"
                                            tabIndex="-1"
                                            className="btn btn-secondary btn-block"
                                            data-dismiss="has-modal"
                                            onClick={this.close()}
                                        >
                                            Close
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="modal-backdrop fade show" />
            </React.Fragment>
        )
    }
}

Wizard = connect(
    state => {
        return {
            user: state.users.cache["me"],
            loading: state.users.loading || state.users.cache["me"] === undefined
        }
    },
    {
        loadUser,
        save,
        goBack
    }
)(Wizard)

export default Wizard
