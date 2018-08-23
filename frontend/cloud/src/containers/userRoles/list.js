import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import classnames from "classnames"

import { loadUsers } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadLocations } from "../../modules/locations"
import { loadOrganizations } from "../../modules/organizations"
import { loadClinics } from "../../modules/clinics"
import { loadAllUserRoles, saveUserRole, deleteUserRole } from "../../modules/userRoles"
import { SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { getName } from "../../utils/user"

class UserRoles extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
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
        if (!this.props.userRoles) {
            this.props.loadAllUserRoles()
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.users && !nextProps.usersLoading) {
            this.props.loadUsers()
        }
        if (!nextProps.roles && !nextProps.rolesLoading) {
            this.props.loadRoles()
        }
        if (!nextProps.clinics && !nextProps.clinicsLoading) {
            this.props.loadClinics()
        }
        if (!nextProps.locations && !nextProps.locationsLoading) {
            this.props.loadLocations()
        }
        if (!nextProps.organizations && !nextProps.organizationsLoading) {
            this.props.loadOrganizations()
        }
        if (!nextProps.userRoles && !nextProps.userRolesLoading) {
            this.props.loadAllUserRoles()
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.userRoles ||
            props.userRolesLoading ||
            !props.users ||
            props.usersLoading ||
            !props.roles ||
            props.rolesLoading ||
            !props.clinics ||
            props.clinicsLoading ||
            !props.locations ||
            props.locationsLoading ||
            !props.organizations ||
            props.organizationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading

        this.setState({
            loading: loading,
            userRoles: _.values(props.userRoles)
        })
    }

    newUserRole() {
        return e => {
            if (this.state.userRoles) {
                let userRoles = [...this.state.userRoles, { edit: true, canSave: false, userID: "", roleID: "", domainType: "", domainID: "" }]
                this.setState({ userRoles: userRoles })
            }
        }
    }

    editUserID(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].userID = e.target.value
            userRoles[index].canSave =
                userRoles[index].userID.length !== 0 &&
                userRoles[index].roleID.length !== 0 &&
                userRoles[index].domainType.length !== 0 &&
                userRoles[index].domainID.length !== 0
            this.setState({ userRoles: userRoles })
        }
    }

    editRoleID(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].roleID = e.target.value
            userRoles[index].canSave =
                userRoles[index].userID.length !== 0 &&
                userRoles[index].roleID.length !== 0 &&
                userRoles[index].domainType.length !== 0 &&
                userRoles[index].domainID.length !== 0
            this.setState({ userRoles: userRoles })
        }
    }

    editDomainType(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].domainType = e.target.value
            if (userRoles[index].domainType === "global") {
                userRoles[index].domainID = "*"
            } else {
                userRoles[index].domainID = ""
            }
            userRoles[index].canSave =
                userRoles[index].userID.length !== 0 &&
                userRoles[index].roleID.length !== 0 &&
                userRoles[index].domainType.length !== 0 &&
                userRoles[index].domainID.length !== 0
            this.setState({ userRoles: userRoles })
        }
    }

    editDomainID(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].domainID = e.target.value
            userRoles[index].canSave =
                userRoles[index].userID.length !== 0 &&
                userRoles[index].roleID.length !== 0 &&
                userRoles[index].domainType.length !== 0 &&
                userRoles[index].domainID.length !== 0
            this.setState({ userRoles: userRoles })
        }
    }

    saveUserRole(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].index = index
            userRoles[index].edit = false
            userRoles[index].saving = true
            this.props.saveUserRole(this.state.userRoles[index])
        }
    }

    cancelNewUserRole(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles.splice(index, 1)
            this.setState({ userRoles: userRoles })
        }
    }

    deleteUserRole(userRoleID) {
        return e => {
            this.props.deleteUserRole(userRoleID)
        }
    }

    getDomainName(domainType, domainID) {
        switch (domainType) {
            case "location":
                if (this.props.locations[domainID]) {
                    return <Link to={`/locations/${domainID}`}>{this.props.locations[domainID].name}</Link>
                }
                return domainID
            case "organization":
                if (this.props.organizations[domainID]) {
                    return <Link to={`/organizations/${domainID}`}>{this.props.organizations[domainID].name}</Link>
                }
                return domainID
            case "clinic":
                if (this.props.clinics[domainID]) {
                    return <Link to={`/clinics/${domainID}`}>{this.props.clinics[domainID].name}</Link>
                }
                return domainID
            case "user":
                if (this.props.users[domainID]) {
                    return (
                        <Link to={`/users/${domainID}`} title={getName(this.props.users[domainID])}>
                            {this.props.users[domainID].username}
                        </Link>
                    )
                }
                return domainID
            default:
                return domainID
        }
    }

    getDomainSelect(index, domainType, domainID) {
        switch (domainType) {
            case "organization":
                return (
                    <select className="form-control" value={domainID || ""} onChange={this.editDomainID(index)}>
                        <option>Select organization</option>
                        {_.map(this.props.organizations, organization => (
                            <option key={organization.id} value={organization.id}>
                                {organization.name}
                            </option>
                        ))}
                        <option value="*">*</option>
                    </select>
                )
            case "clinic":
                return (
                    <select className="form-control" value={domainID || ""} onChange={this.editDomainID(index)}>
                        <option>Select clinic</option>
                        {_.map(this.props.clinics, clinic => (
                            <option key={clinic.id} value={clinic.id}>
                                {clinic.name}
                            </option>
                        ))}
                        <option value="*">*</option>
                    </select>
                )
            case "location":
                return (
                    <select className="form-control" value={domainID || ""} onChange={this.editDomainID(index)}>
                        <option>Select location</option>
                        {_.map(this.props.locations, location => (
                            <option key={location.id} value={location.id}>
                                {location.name}
                            </option>
                        ))}
                        <option value="*">*</option>
                    </select>
                )
            case "user":
                return (
                    <select className="form-control" value={domainID || ""} onChange={this.editDomainID(index)}>
                        <option>Select user</option>
                        {_.map(this.props.users, user => (
                            <option key={user.id} value={user.id}>
                                {user.username} - {getName(user)}
                            </option>
                        ))}
                        <option value="*">*</option>
                    </select>
                )
            case "cloud":
            case "global":
                return (
                    <select className="form-control" value={domainID || ""} onChange={this.editDomainID(index)}>
                        <option key="*" value="*">
                            *
                        </option>
                    </select>
                )
            default:
                return (
                    <select className="form-control" value={domainID || ""} disabled={true}>
                        <option>Select domain type first</option>
                    </select>
                )
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

        return (
            <div>
                <table className="table">
                    <thead>
                        <tr>
                            <th className="w-7" scope="col">
                                #
                            </th>
                            <th className="w-20" scope="col">
                                User
                            </th>
                            <th className="w-15" scope="col">
                                Role
                            </th>
                            <th className="w-17" scope="col">
                                Domain type
                            </th>
                            <th className="w-20" scope="col">
                                Domain
                            </th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {_.map(this.state.userRoles, (userRole, i) => (
                            <tr
                                key={userRole.id || i}
                                className={classnames({
                                    "table-edit": props.canEdit && userRole.edit
                                })}
                            >
                                <th className="w-7" scope="row">
                                    {i + 1}
                                </th>
                                <td className="w-20">
                                    {props.canEdit && userRole.edit ? (
                                        <select className="form-control" value={userRole.userID || ""} onChange={this.editUserID(i)}>
                                            <option>Select user</option>
                                            {_.map(props.users, user => (
                                                <option key={user.id} value={user.id}>
                                                    {user.username} - {getName(user)}
                                                </option>
                                            ))}
                                        </select>
                                    ) : (
                                        <Link to={`/users/${userRole.userID}`} title={getName(props.users[userRole.userID])}>
                                            {props.users[userRole.userID].username}
                                        </Link>
                                    )}
                                </td>
                                <td className="w-15">
                                    {props.canEdit && userRole.edit ? (
                                        <select className="form-control" value={userRole.roleID || ""} onChange={this.editRoleID(i)}>
                                            <option>Select role</option>
                                            {_.map(props.roles, role => (
                                                <option key={role.id} value={role.id}>
                                                    {role.name}
                                                </option>
                                            ))}
                                        </select>
                                    ) : (
                                        <Link to={`/roles/${userRole.roleID}`}>{props.roles[userRole.roleID].name}</Link>
                                    )}
                                </td>
                                <td className="w-17">
                                    {props.canEdit && userRole.edit ? (
                                        <select className="form-control" value={userRole.domainType || ""} onChange={this.editDomainType(i)}>
                                            <option>Select domain type</option>
                                            <option key="global" value="global">
                                                global
                                            </option>
                                            <option key="cloud" value="cloud">
                                                cloud
                                            </option>
                                            <option key="organization" value="organization">
                                                organization
                                            </option>
                                            <option key="clinic" value="clinic">
                                                clinic
                                            </option>
                                            <option key="location" value="location">
                                                location
                                            </option>
                                            <option key="user" value="user">
                                                user
                                            </option>
                                        </select>
                                    ) : (
                                        userRole.domainType
                                    )}
                                </td>
                                <td className="w-20">
                                    {props.canEdit && userRole.edit
                                        ? this.getDomainSelect(i, userRole.domainType, userRole.domainID)
                                        : this.getDomainName(userRole.domainType, userRole.domainID)}
                                </td>
                                <td className="text-right">
                                    {props.canEdit ? (
                                        userRole.edit ? (
                                            <div>
                                                <button
                                                    className="btn btn-secondary"
                                                    disabled={userRole.saving}
                                                    type="button"
                                                    onClick={this.cancelNewUserRole(i)}
                                                >
                                                    Remove
                                                </button>
                                                <button
                                                    className="btn btn-primary"
                                                    disabled={userRole.saving || !userRole.canSave}
                                                    type="button"
                                                    onClick={this.saveUserRole(i)}
                                                >
                                                    Add
                                                </button>
                                            </div>
                                        ) : (
                                            <button onClick={this.deleteUserRole(userRole.id)} className="btn btn-link" type="button">
                                                <span className="remove-link">Remove</span>
                                            </button>
                                        )
                                    ) : null}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {props.canEdit ? (
                    <button type="button" className="btn btn-link" onClick={this.newUserRole()}>
                        Add New User Role
                    </button>
                ) : null}
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    userRoles: ownProps.userRoles
        ? state.userRoles.allLoaded
            ? _.fromPairs(_.map(ownProps.userRoles, userRoleID => [userRoleID, state.userRoles.userRoles[userRoleID]]))
            : undefined
        : state.userRoles.allLoaded
            ? state.userRoles.userRoles
            : undefined,
    userRolesLoading: state.userRoles.loading,
    users: state.users.allLoaded ? state.users.users : undefined,
    usersLoading: state.users.loading,
    roles: state.roles.allLoaded ? state.roles.roles : undefined,
    rolesLoading: state.roles.loading,
    locations: state.locations.allLoaded ? state.locations.locations : undefined,
    locationsLoading: state.locations.loading,
    organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
    organizationsLoading: state.organizations.loading,
    clinics: state.clinics.allLoaded ? state.clinics.clinics : undefined,
    clinicsLoading: state.clinics.loading,
    canEdit: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden:
        state.userRoles.forbidden ||
        state.users.forbidden ||
        state.roles.forbidden ||
        state.locations.forbidden ||
        state.organizations.forbidden ||
        state.clinics.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRoles,
            loadLocations,
            loadOrganizations,
            loadClinics,
            loadAllUserRoles,
            saveUserRole,
            deleteUserRole,
            loadUserRights
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(UserRoles)
