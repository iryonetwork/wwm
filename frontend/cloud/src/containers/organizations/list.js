import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadOrganizations, deleteOrganization } from "../../modules/organizations"

class Organizations extends React.Component {
    componentDidMount() {
        this.props.loadOrganizations()
    }

    removeOrganization = organizationID => e => {
        this.props.deleteOrganization(organizationID)
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (props.loading) {
            return <div>Loading...</div>
        }
        let i = 0
        return (
            <table className="table table-hover">
                <thead>
                    <tr>
                        <th scope="col">#</th>
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
                            <th scope="row">{++i}</th>
                            <td><Link to={`/organizations/${organization.id}`}>{organization.name}</Link></td>
                            <td>{organization.legalStatus}</td>
                            <td>{organization.serviceType}</td>
                            <td>{organization.clinics.length}</td>
                            <td className="text-right">
                                <button onClick={this.removeOrganization(organization.id)} className="btn btn-sm btn-light" type="button">
                                    <span className="icon_trash" />
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    organizations:
        (ownProps.organizations ? (state.organizations.organizations ? _.fromPairs(_.map(ownProps.organizations, organizationID => [organizationID, state.organizations.organizations[organizationID]])) : {}) : state.organizations.organizations) ||
        {},
    loading: state.organizations.loading
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganizations,
            deleteOrganization
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Organizations)
