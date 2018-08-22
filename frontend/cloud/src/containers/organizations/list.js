import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadOrganizations, deleteOrganization } from "../../modules/organizations"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { confirmationDialog } from "shared/utils"

import "../../styles/style.css"

class Organizations extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.organizations) {
            this.props.loadOrganizations()
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.organizations && !nextProps.organizationsLoading) {
            this.props.loadOrganizations()
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.organizations || props.organizationsLoading || props.canEdit === undefined || props.canSee === undefined || props.validationsLoading

        this.setState({ loading: loading })
    }

    removeOrganization(organizationID) {
        return e => {
            confirmationDialog(`Click OK to confirm that you want to remove organization ${this.props.organizations[organizationID].name}.`, () => {
                this.props.deleteOrganization(organizationID)
            })
        }
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        let i = 0
        return (
            <table className="table">
                <thead>
                    <tr>
                        <th className="w-7" scope="col">
                            #
                        </th>
                        <th scope="col">Name</th>
                        <th scope="col">Legal status</th>
                        <th scope="col">Service type</th>
                        <th scope="col">Clinics</th>
                        <th />
                    </tr>
                </thead>
                <tbody>
                    {_.map(_.filter(props.organizations, organization => organization), organization => (
                        <tr key={organization.id}>
                            <th className="w-10" scope="row">
                                {++i}
                            </th>
                            <td>
                                <Link to={`/organizations/${organization.id}`}>{organization.name}</Link>
                            </td>
                            <td>{organization.legalStatus}</td>
                            <td>{organization.serviceType}</td>
                            <td>{organization.clinics.length}</td>
                            <td className="text-right">
                                {props.canEdit ? (
                                    <button onClick={this.removeOrganization(organization.id)} className="btn btn-link" type="button">
                                        <span className="remove-link">Remove</span>
                                    </button>
                                ) : null}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    organizations: ownProps.organizations
        ? state.organizations.allLoaded
            ? _.fromPairs(_.map(ownProps.organizations, organizationID => [organizationID, state.organizations.organizations[organizationID]]))
            : undefined
        : state.organizations.allLoaded
            ? state.organizations.organizations
            : undefined,
    organizationsLoading: state.organizations.loading,
    canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden: state.organizations.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganizations,
            deleteOrganization,
            loadUserRights
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Organizations)
