import _ from "lodash"
import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { loadStatus } from "../../modules/status"
import { Popover, PopoverBody } from 'reactstrap';

import "./style.css"

class Status extends React.Component {
    constructor(props) {
        super(props)
        props.loadStatus()
        this.state = { loading: true }
    }

    componentDidMount() {
        clearInterval(this.interval)
        this.interval = setInterval(this.props.loadStatus, 5000);
    }

    toggleStatusPopover = () => {
        this.setState({
            statusPopoverOpen: !this.state.statusPopoverOpen
        });
    }

    render() {
        return (
            <div>
            <div id="status" className="status">
                {this.props.status ? (
                    <span id="statusPopover" className="statusInfo clickable" onClick={this.toggleStatusPopover}>
                        <span className={"light " + this.props.status.status} />
                        <span className="link">
                            {this.props.status.status === "ok" ? "Online" : (this.props.status.status === "warning" ? "Warning" : "Offline")}
                        </span>
                    </span>
                ) : (null)}
            </div>
            <Popover modifiers="" placement="bottom-start" className="statusPopover" isOpen={this.state.statusPopoverOpen} target="statusPopover" toggle={this.toggleStatusPopover}>
                    <PopoverBody>
                    <ul>
                    <h4>Status details</h4>
                        {this.props.status.local ? (
                          <li key="local" className="statusInfo" id="localStatus">
                            <span className={"light " + this.props.status.local.status} />
                            Local
                            <ul>
                                {_.map(this.props.status.local.components, (value, key) => (
                                    <li key={key} className="statusInfo">
                                       <span className={"light " + value.status} />
                                       {key}
                                    </li>
                                ))}
                            </ul>
                          </li>
                        ) : (null)}
                        {this.props.status.cloud ? (
                            <li key="cloud" className="statusInfo" id="cloudStatus">
                                <span className={"light " + this.props.status.cloud.status} />
                                Cloud
                                <ul>
                                    {_.map(this.props.status.cloud.components, (value, key) => (
                                    <li key={key} className="statusInfo">
                                           <span className={"light " + value.status} />
                                           {key}
                                        </li>
                                    ))}
                                </ul>
                            </li>
                        ) : (null)}
                        {this.props.status.external ? (
                          <li key="external" className="statusInfo" id="externalStatus">
                            <span className={"light " + this.props.status.external.status} />
                            External
                            <ul>
                                {_.map(this.props.status.external.components, (value, key) => (
                                    <li key={key} className="statusInfo">
                                       <span className={"light " + value.status} />
                                       {key}
                                    </li>
                                ))}
                            </ul>
                          </li>
                        ) : (null)}
                    </ul>
                    </PopoverBody>
            </Popover>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        status: state.status.status
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadStatus
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Status)
