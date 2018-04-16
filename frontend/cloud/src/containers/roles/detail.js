import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"

import { loadRules } from "../../modules/rules"
import { open } from "shared/modules/alert"
import Rules from "../rules"

class DetailRole extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }

    }

    componentDidMount() {
        if (!this.props.rules) {
            this.props.loadRules()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.rules && !nextProps.rulesLoading) {
            this.props.loadRules()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.rules || props.rulesLoading
        this.setState({loading: loading})
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        return (
            <div>
                <Rules rules={props.rules} subject={props.roleID} />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.roleID
    if (!id) {
        id = ownProps.match.params.roleID
    }
    return {
        roleID: ownProps.match.params.id,
        rules: state.rules.subjects ? (state.rules.subjects[id] ? state.rules.subjects[id] : {}) : undefined,
        rulesLoading: state.rules.loading,
        forbidden: state.roles.forbidden || state.rules.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRules,
            open
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(DetailRole)
