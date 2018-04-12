import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import moment from "moment"

import { loadOrganization, saveOrganization } from "../../modules/organizations"
import { open, close, COLOR_DANGER } from "shared/modules/alert"

class OrganizationDetail extends React.Component {
    constructor(props) {
        super(props)
        if (props.organization) {
            this.state = {
                name: props.organization.name,
                legalStatus: props.organization.legalStatus ? props.organization.legalStatus : "",
                serviceType: props.organization.serviceType ? props.organization.serviceType : "",
            }
        } else {
            this.state = {
                name: "",
                legalStatus: "",
                serviceType: "",
            }
        }
    }

    componentDidMount() {
        if (!this.props.organization && this.props.organizationID !== "new") {
            this.props.loadOrganization(this.props.organizationID)
        }
    }

    componentWillReceiveProps(props) {
        if (props.organization) {
            this.setState({ name: props.organization.name })
            this.setState({ legalStatus: props.organization.legalStatus ? props.organization.legalStatus : "" })
            this.setState({ serviceType: props.organization.serviceType ? props.organization.serviceType : "" })
        }
    }

    updateInput = e => {
        const target = e.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const id = target.id;

        this.setState({
          [id]: value
        });
    }

    submit = e => {
        e.preventDefault()
        this.props.close()

        let organization = this.props.organization

        organization.name = this.state.name
        organization.legalStatus = this.state.legalStatus
        organization.serviceType = this.state.serviceType

        this.props.saveOrganization(organization)
        this.forceUpdate()
    }

    render() {
        let props = this.props
        console.log(props)
        if (!props.organization && props.organizationID !== "new") {
            return <div>Loading...</div>
        }
        return (
            <div>
                <h1>Organizations</h1>
                <h2>{props.organization ? this.props.organization.name : "Add new organization"}</h2>

                <form onSubmit={this.submit}>
                    <div className="form-group">
                        <label htmlFor="name">Name</label>
                        <input className="form-control" id="name" value={this.state.name} onChange={this.updateInput} placeholder="Organization name" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="legalStatus">Legal status</label>
                        <input className="form-control" id="legalStatus" value={this.state.legalStatus} onChange={this.updateInput} placeholder="e.g. NGO" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="country">Service type</label>
                        <input className="form-control" id="serviceType" value={this.state.serviceType} onChange={this.updateInput} placeholder="e.g. Basic care" />
                    </div>
                    <button type="submit" className="btn btn-sm btn-outline-secondary">
                        Save
                    </button>
                </form>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.organizationID
    if (!id) {
        id = ownProps.match.params.id
    }

    return {
        organization: state.organizations.organizations ? state.organizations.organizations[id] : undefined,
        loading: state.organizations.loading,
        organizationID: id,
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganization,
            saveOrganization,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(OrganizationDetail))
