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

        this.state = {}
    }
    componentDidMount() {
        this.props.loadRules()
    }

    render() {
        let props = this.props
        return (
            <div>
                <Rules rules={props.rules} subject={props.roleID} />
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        role: state.roles.roles[ownProps.match.params.id],
        rules: state.rules.subjects ? state.rules.subjects[ownProps.match.params.id] || [] : [],
        roleID: ownProps.match.params.id
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
