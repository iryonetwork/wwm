import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import classnames from "classnames"
import { Popover, PopoverBody } from "reactstrap"

import { toggleBodyScroll } from "shared/utils"
import SettingsContent from "./index"
import { ReactComponent as SettingsIcon } from "shared/icons/settings-active.svg"

import "./style.css"

class SettingsPopover extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
        this.toggleSettingsPopover = this.toggleSettingsPopover.bind(this)
    }

    toggleSettingsPopover(e) {
        e.preventDefault()
        toggleBodyScroll()
        this.setState({
            settingsPopoverOpen: !this.state.settingsPopoverOpen
        })
    }

    render() {
        return (
            <React.Fragment>
                <a
                    id="settingsPopover"
                    className={classnames("navigation", { active: this.state.settingsPopoverOpen })}
                    href="/"
                    onClick={this.toggleSettingsPopover}
                >
                    <SettingsIcon />
                    Settings
                </a>
                <Popover
                    placement="right"
                    modifiers={{
                        preventOverflow: { enabled: true, boundariesElement: "html" },
                        flip: { enabled: false },
                        offset: { offset: "-140px,-95px", order: 800, enabled: true }
                    }}
                    className="settingsPopover"
                    isOpen={this.state.settingsPopoverOpen}
                    target="settingsPopover"
                    toggle={this.toggleSettingsPopover}
                >
                    <PopoverBody>
                        <SettingsContent wideInput={true} />
                    </PopoverBody>
                </Popover>
            </React.Fragment>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        status: state.status.status
    }
}

const mapDispatchToProps = dispatch => bindActionCreators({}, dispatch)

export default connect(mapStateToProps, mapDispatchToProps)(SettingsPopover)
