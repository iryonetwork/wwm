import React from "react"
import { connect } from "react-redux"
import { bindActionCreators } from "redux"
import { Link, withRouter } from "react-router-dom"

import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import Reports from "./list"

class ReportsIndex extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
    }

    componentDidMount() {
        if (this.props.canSee === undefined) {
            this.props.loadUserRights()
        }
        if (this.props.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.canSee === undefined && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }
        if (nextProps.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = props.canSee === undefined || props.validationsLoading
        this.setState({ loading: loading })
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <React.Fragment>
                <div>
                    <h1>Patients report</h1>
                    <Reports reportType="patients" />
                </div>
                <div>
                    <h1>Encounters report</h1>
                    <Reports reportType="encounters" />
                </div>
            </React.Fragment>
        )
    }
}

const mapStateToProps = state => ({
    canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden: false
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ReportsIndex))
