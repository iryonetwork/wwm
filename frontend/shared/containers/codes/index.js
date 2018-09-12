import React from "react"
import { connect } from "react-redux"

import { fetchCode } from "shared/modules/codes"

class CodeTitle extends React.Component {
    constructor(props) {
        super(props)
        this.loadCode = this.loadCode.bind(this)
        if (props.codeId) {
            this.loadCode(props.categoryId, props.codeId)
        }

        this.componentWillUnmount = this.componentWillUnmount.bind(this)
        this.state = { loading: true, title: "" }
    }

    loadCode(categoryId, codeId) {
        this.props
            .fetchCode(categoryId, codeId)
            .then(code => {
                if (!this.unmounted) {
                    this.setState({
                        loading: false,
                        failed: false,
                        title: code.title
                    })
                }
            })
            .catch(ex => {
                this.setState({ loading: false, failed: true })
            })
    }

    componentWillUnmount() {
        this.unmounted = true
    }

    render() {
        if (this.state.loading) {
            return <span>...</span>
        }

        if (this.state.failed) {
            return <span>Failed to fetch title!</span>
        }

        return <span>{this.state.title}</span>
    }
}

CodeTitle = connect(state => ({}), { fetchCode })(CodeTitle)

export { CodeTitle }
