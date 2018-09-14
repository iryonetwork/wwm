import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import moment from "moment"
import _ from "lodash"
import fileDownload from "js-file-download"

import { loadReportsByType, readFile } from "../../modules/reports"

class Reports extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.reports) {
            this.props.loadReportsByType(this.props.reportType)
        }
        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.reports && !nextProps.reportsLoading) {
            this.props.loadReportsByType(nextProps.reportType)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.reports || props.reportsLoading
        this.setState({ loading: loading })
    }

    downloadFile(index) {
        return e => {
            e.preventDefault()
            this.props.readFile(this.props.reportType, this.props.reports[index].name).then(data => {
                let name = this.props.reportType + " (" + moment(this.props.reports[index].created).format("LLL") + ").xlsx"
                fileDownload(data, name)
            })
        }
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        return (
            <table className="table">
                <thead>
                    {_.map(props.reports, (report, index) => (
                        <React.Fragment key={report.name}>
                            <tr>
                                <th className="w-20" scope="row">
                                    Report creation time
                                </th>
                                <th />
                            </tr>
                            <tr>
                                <td className="w-20">
                                    <span>{moment(report.created).format("DD-MM-YYYY HH:mm:ss")}</span>
                                </td>
                                <td>
                                    <Link to="/" onClick={this.downloadFile(index)}>
                                        Download
                                    </Link>
                                </td>
                            </tr>
                        </React.Fragment>
                    ))}
                </thead>
            </table>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        reports: _.get(state.reports.reports, ownProps.reportType, undefined),
        reportsLoading: state.reports.loading,
        forbidden: state.reports.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadReportsByType,
            readFile
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Reports)
