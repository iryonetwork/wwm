import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadUsers } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadLocations } from "../../modules/locations"
import { loadOrganizations } from "../../modules/organizations"
import { loadClinics } from "../../modules/clinics"
import { loadUserRoles, deleteUserRole } from "../../modules/userRoles"

class UserRoles extends React.Component {
    componentDidMount() {
        this.props.loadUserRoles()
        if (!this.props.users) {
            this.props.loadUsers()
        }
        if (!this.props.roles) {
            this.props.loadRoles()
        }
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.organizations) {
            this.props.loadOrganizations()
        }
        if (!this.props.clinics) {
            this.props.loadClinics()
        }
    }

    removeUserRole = userRoleID => e => {
        this.props.deleteUserRole(userRoleID)
    }

    getDomainName(domainType, domainID) {
        switch(domainType) {
            case "location":
                if (this.props.locations[domainID]) {
                    return this.props.locations[domainID].name
                }
                return domainID
            case "organization":
                if (this.props.organizations[domainID]) {
                    return this.props.organizations[domainID].name
                }
                return domainID
            case "clinic":
                if (this.props.clinics[domainID]) {
                    return this.props.clinics[domainID].name
                }
                return domainID
            case "user":
                if (this.props.users[domainID]) {
                    return this.props.users[domainID].username
                }
                return domainID
            default:
                return domainID
        }
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
                        <th scope="col">Username</th>
                        <th scope="col">Role</th>
                        <th scope="col">Domain type</th>
                        <th scope="col">Domain</th>
                        <th />
                    </tr>
                </thead>
                <tbody>
                    {_.map(_.filter(props.userRoles, userRole => userRole), userRole => (
                        <tr key={userRole.id}>
                            <th scope="row">{++i}</th>
                            <td>
                                <Link to={`/users/${userRole.userID}`}>{props.users[userRole.userID].username}</Link>
                            </td>
                            <td>
                                <Link to={`/roles/${userRole.roleID}`}>{props.roles[userRole.roleID].name}</Link>
                            </td>
                            <td>{userRole.domainType}</td>
                            <td>{this.getDomainName(userRole.domainType, userRole.domainID)}</td>
                            <td className="text-right">
                                <button onClick={this.removeUserRole(userRole.id)} className="btn btn-sm btn-light" type="button">
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
    userRoles:
        (ownProps.userRoles ? (state.userRoles.userRoles ? _.fromPairs(_.map(ownProps.userRoles, userRoleID => [userRoleID, state.userRoles.userRoles[userRoleID]])) : {}) : state.userRoles.userRoles) ||
        {},
    users: state.users.users,
    roles: state.roles.roles,
    locations: state.locations.locations,
    organizations: state.organizations.organizations,
    clinics: state.clinics.clinics,
    loading: state.userRoles.loading || state.users.loading || state.roles.loading || state.locations.loading || state.organizations.loading || state.clinics.loading,
    forbidden: state.userRoles.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRoles,
            loadLocations,
            loadOrganizations,
            loadClinics,
            loadUserRoles,
            deleteUserRole
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(UserRoles)
