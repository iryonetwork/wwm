import React from "react"
import { NavLink, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { push } from "react-router-redux"
import _ from "lodash"
import classnames from "classnames"

import { ADVANCED_ROLE_IDS } from "shared/modules/config"
import { loadUser } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadOrganizations, clearOrganizations, deleteUserFromOrganization } from "../../modules/organizations"
import { makeGetUserOrganizationIDs } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, saveUserRoleCustomMsg, deleteUserRole } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import OrganizationDetail from "./organizationDetail"
import { joinPaths } from "shared/utils"
import Spinner from "shared/containers/spinner"
import { confirmationDialog } from "shared/utils"

class OrganizationsList extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.user) {
            this.props.loadUser(this.props.userID)
        }
        if (!this.props.roles) {
            this.props.loadRoles()
        }
        if (!this.props.organizations) {
            this.props.loadOrganizations()
        }
        if (!this.props.userRoles) {
            this.props.loadUserUserRoles(this.props.userID)
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.user && !nextProps.usersLoading) {
            this.props.loadUser(nextProps.userID)
        }
        if (!nextProps.roles && !nextProps.rolesLoading) {
            this.props.loadRoles()
        }
        if (!nextProps.organizations && !nextProps.organizationsLoading) {
            this.props.loadOrganizations()
        }
        if (!nextProps.userRoles && !nextProps.userRolesLoading) {
            this.props.loadUserUserRoles(this.props.userID)
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
            !props.roles ||
            props.rolesLoading ||
            !props.organizations ||
            props.organizationsLoading ||
            !props.userOrganizationIDs ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading

        this.setState({
            loading: loading,
            userOrganizations: _.map(props.userOrganizationIDs, organizationID => {
                return { id: organizationID }
            }),
            selectedOrganizationID: props.selectedOrganizationID ? props.selectedOrganizationID : undefined
        })
    }

    newUserOrganization() {
        return e => {
            if (this.state.userOrganizations) {
                let userOrganizations = [
                    ...this.state.userOrganizations,
                    { id: "", edit: true, canSave: false, userID: this.props.userID, roleID: "", domainType: "organization" }
                ]
                this.setState({ userOrganizations: userOrganizations })
            }
        }
    }

    editOrganizationID(index) {
        return e => {
            let userOrganizations = [...this.state.userOrganizations]
            userOrganizations[index].id = e.target.value
            userOrganizations[index].canSave = userOrganizations[index].id.length !== 0 && userOrganizations[index].roleID.length !== 0
            this.setState({ userOrganizations: userOrganizations })
        }
    }

    editRoleID(index) {
        return e => {
            let userOrganizations = [...this.state.userOrganizations]
            userOrganizations[index].roleID = e.target.value
            userOrganizations[index].canSave = userOrganizations[index].id.length !== 0 && userOrganizations[index].roleID.length !== 0
            this.setState({ userOrganizations: userOrganizations })
        }
    }

    saveUserOrganization(index) {
        return e => {
            let userOrganizations = [...this.state.userOrganizations]
            let userRole = {}
            userRole.userID = userOrganizations[index].userID
            userRole.roleID = userOrganizations[index].roleID
            userRole.domainType = userOrganizations[index].domainType
            userRole.domainID = userOrganizations[index].id
            userOrganizations[index].index = index
            userOrganizations[index].edit = false
            userOrganizations[index].saving = true

            this.props.saveUserRoleCustomMsg(userRole, "Added User to the Organization").then(response => {
                if (response && response.domainID) {
                    this.props.history.push(`${this.props.basePath}/organizations/${response.domainID}`)
                }
            })
        }
    }

    cancelNewUserOrganization(index) {
        return e => {
            let userOrganizations = [...this.state.userOrganizations]
            userOrganizations.splice(index, 1)
            this.setState({ userOrganizations: userOrganizations })
        }
    }

    removeUserOrganization(index) {
        return e => {
            confirmationDialog(
                `Click OK to confirm that you want to remove the user from organization ${
                    this.props.organizations[this.state.userOrganizations[index].id].name
                }.`,
                () => {
                    this.props.deleteUserFromOrganization(this.state.userOrganizations[index].id, this.props.userID)
                    if (this.state.selectedOrganizationID === this.state.userOrganizations[index].id) {
                        this.props.history.push(`${this.props.basePath}/organizations`)
                    }
                }
            )
        }
    }

    render() {
        let props = this.props
        if (!props.user || props.usersLoading) {
            return <Spinner />
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <div>
                <header>{props.isSelf ? <h1>My Profile</h1> : <h1>Users</h1>}</header>
                <h2>{props.user.username}</h2>
                <div className="navigation">
                    {props.canSeePersonal ? (
                        <NavLink exact to={props.basePath}>
                            Personal Info
                        </NavLink>
                    ) : null}
                    {props.canSee ? <NavLink to={joinPaths(props.basePath, "organizations")}>Organizations</NavLink> : null}
                    {props.canSeeClinics ? <NavLink to={joinPaths(props.basePath, "clinics")}>Clinics</NavLink> : null}
                    {props.canSeeWildcardUserRoles ? <NavLink to={joinPaths(props.basePath, "userroles")}>Wildcard Roles</NavLink> : null}
                </div>
                {this.state.loading ? (
                    <Spinner />
                ) : (
                    <div id="organizations">
                        <div className="row">
                            <div className="col-12">
                                {this.state.userOrganizations.length > 0 ? (
                                    <table className="table">
                                        <thead>
                                            <tr>
                                                <th className="w-7" scope="col">
                                                    #
                                                </th>
                                                <th className="w-40" scope="col">
                                                    Organization name
                                                </th>
                                                <th />
                                                <th className="w-25" />
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {_.map(this.state.userOrganizations, (userOrganization, i) => (
                                                <React.Fragment key={userOrganization.id || i}>
                                                    <tr
                                                        className={classnames({
                                                            "table-active": this.state.selectedOrganizationID === userOrganization.id,
                                                            "table-edit": props.canEdit && userOrganization.edit
                                                        })}
                                                    >
                                                        <th className="w-7" scope="row">
                                                            {i + 1}
                                                        </th>
                                                        <td className="w-40">
                                                            {props.canEdit && userOrganization.edit ? (
                                                                <select
                                                                    className="form-control"
                                                                    value={userOrganization.id || ""}
                                                                    onChange={this.editOrganizationID(i)}
                                                                >
                                                                    <option value="">Select organization</option>
                                                                    {_.map(
                                                                        _.difference(
                                                                            _.map(_.values(props.organizations), organization => organization.id),
                                                                            _.without(
                                                                                _.map(this.state.userOrganizations, organization => organization.id),
                                                                                userOrganization.id
                                                                            )
                                                                        ),
                                                                        organizationID => (
                                                                            <option key={organizationID} value={organizationID}>
                                                                                {props.organizations[organizationID].name}
                                                                            </option>
                                                                        )
                                                                    )}
                                                                </select>
                                                            ) : (
                                                                <Link to={`/organizations/${userOrganization.id}`}>
                                                                    {props.organizations[userOrganization.id].name}
                                                                </Link>
                                                            )}
                                                        </td>
                                                        <td>
                                                            {props.canEdit && userOrganization.edit ? (
                                                                <select
                                                                    className="form-control"
                                                                    value={userOrganization.roleID || ""}
                                                                    onChange={this.editRoleID(i)}
                                                                >
                                                                    <option value="">Select role</option>
                                                                    {_.map(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id)), role => (
                                                                        <option key={role.id} value={role.id}>
                                                                            {role.name}
                                                                        </option>
                                                                    ))}
                                                                </select>
                                                            ) : null}
                                                        </td>
                                                        <td className="w-25 text-right">
                                                            {props.canEdit ? (
                                                                userOrganization.edit ? (
                                                                    <div>
                                                                        <button
                                                                            className="btn btn-secondary"
                                                                            disabled={userOrganization.saving}
                                                                            type="button"
                                                                            onClick={this.cancelNewUserOrganization(i)}
                                                                        >
                                                                            Remove
                                                                        </button>
                                                                        <button
                                                                            className="btn btn-primary"
                                                                            disabled={userOrganization.saving || !userOrganization.canSave}
                                                                            type="button"
                                                                            onClick={this.saveUserOrganization(i)}
                                                                        >
                                                                            Add
                                                                        </button>
                                                                    </div>
                                                                ) : (
                                                                    <div>
                                                                        {this.state.selectedOrganizationID === userOrganization.id ? (
                                                                            <button
                                                                                className="btn btn-link"
                                                                                type="button"
                                                                                onClick={() => this.props.push(`/users/${props.userID}/organizations`)}
                                                                            >
                                                                                Hide Roles
                                                                                <span className="arrow-up-icon" />
                                                                            </button>
                                                                        ) : (
                                                                            <button
                                                                                className="btn btn-link"
                                                                                type="button"
                                                                                onClick={() =>
                                                                                    this.props.push(
                                                                                        `/users/${props.userID}/organizations/${userOrganization.id}`
                                                                                    )
                                                                                }
                                                                            >
                                                                                Show Roles<span className="arrow-down-icon" />
                                                                            </button>
                                                                        )}
                                                                        {props.canEdit ? (
                                                                            <button
                                                                                className="btn btn-link"
                                                                                type="button"
                                                                                onClick={this.removeUserOrganization(i)}
                                                                            >
                                                                                <span className="remove-link">Remove</span>
                                                                            </button>
                                                                        ) : null}
                                                                    </div>
                                                                )
                                                            ) : null}
                                                        </td>
                                                    </tr>
                                                    {this.state.selectedOrganizationID === userOrganization.id ? (
                                                        <React.Fragment>
                                                            <tr className="table-active">
                                                                <td colSpan="4" className="row-details-container">
                                                                    <OrganizationDetail userID={this.props.userID} organizationID={userOrganization.id} />
                                                                </td>
                                                            </tr>
                                                        </React.Fragment>
                                                    ) : null}
                                                </React.Fragment>
                                            ))}
                                        </tbody>
                                    </table>
                                ) : (
                                    <h3>User does not belong to any organization.</h3>
                                )}
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={
                                            this.state.userOrganizations.length !== 0 &&
                                            this.state.userOrganizations[this.state.userOrganizations.length - 1].edit
                                                ? true
                                                : null
                                        }
                                        onClick={this.newUserOrganization()}
                                    >
                                        Add Current User to an Organization
                                    </button>
                                ) : null}
                            </div>
                        </div>
                    </div>
                )}
            </div>
        )
    }
}

const makeMapStateToProps = () => {
    const getUserOrganizationIDs = makeGetUserOrganizationIDs()

    const mapStateToProps = (state, ownProps) => {
        let userID = ownProps.userID
        if (!userID) {
            userID = ownProps.match.params.userID
        }
        let isSelf = state.authentication.token.sub === userID

        let selectedOrganizationID = ownProps.organizationID
        if (!selectedOrganizationID) {
            selectedOrganizationID = ownProps.match.params.organizationID
        }
        return {
            basePath: ownProps.home ? "/me" : `/users/${userID}`,
            isSelf: isSelf,
            userID: userID,
            user: state.users.users ? state.users.users[userID] : undefined,
            usersLoading: state.users.loading,
            selectedOrganizationID: selectedOrganizationID,
            organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
            organizationsLoading: state.organizations.loading,
            advancedRoleIDs: state.config[ADVANCED_ROLE_IDS],
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            userOrganizationIDs: getUserOrganizationIDs(state, { userID: userID }),
            canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSeePersonal: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeWildcardUserRoles: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.userRoles.forbidden || state.users.forbidden || state.organizations.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            push,
            loadRoles,
            loadUser,
            loadOrganizations,
            clearOrganizations,
            loadUserUserRoles,
            saveUserRoleCustomMsg,
            deleteUserFromOrganization,
            deleteUserRole,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(OrganizationsList))
