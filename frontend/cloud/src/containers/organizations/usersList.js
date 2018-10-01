import React from "react"
import { NavLink, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import classnames from "classnames"
import { push } from "react-router-redux"

import { ADVANCED_ROLE_IDS } from  "../../modules/config"
import { loadUsers } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadOrganization, deleteUserFromOrganization } from "../../modules/organizations"
import { saveUserRoleCustomMsg, loadDomainUserRoles } from "../../modules/userRoles"
import { SELF_RIGHTS_RESOURCE, ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { makeGetOrganizationUserIDs } from "../../selectors/userRolesSelectors"
import { getName } from "../../utils/user"
import UserDetail from "./userDetail"
import Spinner from "shared/containers/spinner"
import { confirmationDialog } from "shared/utils"

class UsersList extends React.Component {
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
        if (!this.props.organization) {
            this.props.loadOrganization(this.props.organizationID)
        }
        if (!this.props.userRoles) {
            this.props.loadDomainUserRoles("organization", this.props.organizationID)
        }
        if (
            this.props.canSee === undefined ||
            this.props.canEdit === undefined ||
            this.props.canSeeOrganization === undefined ||
            this.props.canSeeClinics === undefined
        ) {
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
        if (!nextProps.organization && !nextProps.organizationsLoading) {
            this.props.loadOrganization(nextProps.organizationID)
        }
        if (!nextProps.userRoles && !nextProps.userRolesLoading) {
            this.props.loadDomainUserRoles("organization", this.props.organizationID)
        }
        if (
            (nextProps.canSee === undefined ||
                nextProps.canEdit === undefined ||
                nextProps.canSeeOrganization === undefined ||
                nextProps.canSeeClinics === undefined) &&
            !nextProps.validationsLoading
        ) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.users ||
            props.usersLoading ||
            !props.roles ||
            props.rolesLoading ||
            !props.userRoles ||
            props.userRolesLoading ||
            !props.organization ||
            props.organizationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading
        this.setState({ loading: loading })

        if (!loading) {
            this.setState({
                organizationUsers: _.map(props.organizationUserIDs ? props.organizationUserIDs : [], userID => {
                    return props.users[userID]
                }),
                selectedUserID: props.selectedUserID
            })
        }
    }

    newUser() {
        return e => {
            if (this.state.organizationUsers) {
                let organizationUsers = [
                    ...this.state.organizationUsers,
                    { id: "", edit: true, canSave: false, domainType: "organization", domainID: this.props.organizationID, userID: "", roleID: "" }
                ]
                this.setState({ organizationUsers: organizationUsers })
            }
        }
    }

    editUserID(index) {
        return e => {
            let organizationUsers = [...this.state.organizationUsers]
            organizationUsers[index].userID = e.target.value
            organizationUsers[index].canSave = organizationUsers[index].userID.length !== 0 && organizationUsers[index].roleID.length !== 0
            this.setState({ organizationUsers: organizationUsers })
        }
    }

    editRoleID(index) {
        return e => {
            let organizationUsers = [...this.state.organizationUsers]
            organizationUsers[index].roleID = e.target.value
            organizationUsers[index].canSave = organizationUsers[index].userID.length !== 0 && organizationUsers[index].roleID.length !== 0
            this.setState({ organizationUsers: organizationUsers })
        }
    }

    saveUser(index) {
        return e => {
            let organizationUsers = [...this.state.organizationUsers]

            organizationUsers[index].edit = false
            organizationUsers[index].saving = true

            this.props.saveUserRoleCustomMsg(this.state.organizationUsers[index], "Added User to the Organization")
        }
    }

    cancelNewUser(index) {
        return e => {
            let organizationUsers = [...this.state.organizationUsers]
            organizationUsers.splice(index, 1)
            this.setState({ organizationUsers: organizationUsers })
        }
    }

    removeUser(index) {
        return e => {
            confirmationDialog(`Click OK to confirm that you want to remove user ${this.state.clinicUsers[index].username} from the organization.`, () => {
                this.props.deleteUserFromOrganization(this.props.organizationID, this.state.organizationUsers[index].id)
            })
        }
    }

    render() {
        let props = this.props
        if (!props.organization || props.organizationsLoading) {
            return <Spinner />
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <div>
                <header>
                    <h1>Organizations</h1>
                </header>
                <h2>{props.organization.name}</h2>
                {props.organization ? (
                    <div className="navigation">
                        {props.canSeeOrganization ? (
                            <NavLink exact to={`/organizations/${props.organization.id}`}>
                                Organization's Data
                            </NavLink>
                        ) : null}
                        {props.canSee ? <NavLink to={`/organizations/${props.organization.id}/users`}>Users</NavLink> : null}
                        {props.canSeeClinics ? <NavLink to={`/organizations/${props.organization.id}/clinics`}>Clinics</NavLink> : null}
                    </div>
                ) : null}
                {this.state.loading ? (
                    <Spinner />
                ) : (
                    <div id="users">
                        <div className="row">
                            <div className="col-12">
                                {this.state.organizationUsers.length > 0 ? (
                                    <table className="table">
                                        <thead>
                                            <tr>
                                                <th className="w-7" scope="col">
                                                    #
                                                </th>
                                                <th scope="col">Username</th>
                                                <th scope="col">Name</th>
                                                <th scope="col">Email</th>
                                                <th />
                                                <th />
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {_.map(this.state.organizationUsers, (user, i) => (
                                                <React.Fragment key={user.id || i}>
                                                    {props.canEdit && user.edit ? (
                                                        <tr
                                                            className={classnames({
                                                                "table-active": this.state.selectedUserID === user.id,
                                                                "table-edit": props.canEdit && user.edit
                                                            })}
                                                        >
                                                            <th className="w-7" scope="row">
                                                                {i + 1}
                                                            </th>
                                                            <td colSpan="2">
                                                                <select className="form-control" value={user.userID || ""} onChange={this.editUserID(i)}>
                                                                    <option>Select user</option>
                                                                    {_.map(
                                                                        _.difference(_.map(_.values(props.users), user => user.id), props.organizationUserIDs),
                                                                        userID => (
                                                                            <option key={userID} value={userID}>
                                                                                {props.users[userID].username} - {getName(props.users[userID])} ({
                                                                                    props.users[userID].email
                                                                                })
                                                                            </option>
                                                                        )
                                                                    )}
                                                                </select>
                                                            </td>
                                                            <td colSpan="2">
                                                                <select className="form-control" value={user.roleID || ""} onChange={this.editRoleID(i)}>
                                                                    <option>Select role</option>
                                                                    {_.map(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id)), role => (
                                                                        <option key={role.id} value={role.id}>
                                                                            {role.name}
                                                                        </option>
                                                                    ))}
                                                                </select>
                                                            </td>
                                                            <td className="text-right">
                                                                <div>
                                                                    <button
                                                                        className="btn btn-secondary"
                                                                        disabled={user.saving}
                                                                        type="button"
                                                                        onClick={this.cancelNewUser(i)}
                                                                    >
                                                                        Remove
                                                                    </button>
                                                                    <button
                                                                        className="btn btn-primary"
                                                                        disabled={user.saving || !user.canSave}
                                                                        type="button"
                                                                        onClick={this.saveUser(i)}
                                                                    >
                                                                        Add
                                                                    </button>
                                                                </div>
                                                            </td>
                                                        </tr>
                                                    ) : (
                                                        <tr className={this.state.selectedUserID === user.id ? "table-active" : ""}>
                                                            <th className="w-7" scope="row">
                                                                {i + 1}
                                                            </th>
                                                            <td>
                                                                <Link to={`/users/${user.id}`}>{user.username}</Link>
                                                            </td>
                                                            <td>{getName(user)}</td>
                                                            <td>{user.email}</td>
                                                            <td />
                                                            <td className="text-right">
                                                                <div>
                                                                    {this.state.selectedUserID === user.id ? (
                                                                        <button
                                                                            className="btn btn-link"
                                                                            type="button"
                                                                            onClick={() => this.props.push(`/organizations/${props.organizationID}/users`)}
                                                                        >
                                                                            Hide Roles
                                                                            <span className="arrow-up-icon" />
                                                                        </button>
                                                                    ) : (
                                                                        <button
                                                                            className="btn btn-link"
                                                                            type="button"
                                                                            onClick={() =>
                                                                                this.props.push(`/organizations/${props.organizationID}/users/${user.id}`)
                                                                            }
                                                                        >
                                                                            Show Roles<span className="arrow-down-icon" />
                                                                        </button>
                                                                    )}
                                                                    {props.canEdit ? (
                                                                        <button className="btn btn-link" type="button" onClick={this.removeUser(i)}>
                                                                            <span className="remove-link">Remove</span>
                                                                        </button>
                                                                    ) : null}
                                                                </div>
                                                            </td>
                                                        </tr>
                                                    )}
                                                    {this.state.selectedUserID === user.id ? (
                                                        <tr className="table-active">
                                                            <td colSpan="6" className="row-details-container">
                                                                <UserDetail organizationID={props.organizationID} userID={user.id} />
                                                            </td>
                                                        </tr>
                                                    ) : null}
                                                </React.Fragment>
                                            ))}
                                        </tbody>
                                    </table>
                                ) : (
                                    <h3>No users belong to the organization.</h3>
                                )}
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={
                                            this.state.organizationUsers.length !== 0 &&
                                            this.state.organizationUsers[this.state.organizationUsers.length - 1].edit
                                                ? true
                                                : null
                                        }
                                        onClick={this.newUser()}
                                    >
                                        Add User to the Organization
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
    const getOrganizationUserIDs = makeGetOrganizationUserIDs()

    const mapStateToProps = (state, ownProps) => {
        let organizationID = ownProps.organizationID
        if (!organizationID) {
            organizationID = ownProps.match.params.organizationID
        }
        let selectedUserID = ownProps.userID
        if (!selectedUserID) {
            selectedUserID = ownProps.match.params.userID
        }

        return {
            organizationID: organizationID,
            selectedUserID: selectedUserID ? selectedUserID : undefined,
            organization: state.organizations.organizations ? state.organizations.organizations[organizationID] : undefined,
            organizationsLoading: state.organizations.loading,
            userRoles: state.userRoles.domainUserRoles
                ? state.userRoles.domainUserRoles["organization"] && state.userRoles.domainUserRoles["organization"][organizationID]
                    ? state.userRoles.domainUserRoles["organization"][organizationID]
                    : undefined
                : undefined,
            userRolesLoading: state.userRoles.loading,
            users: state.users.allLoaded ? state.users.users : undefined,
            usersLoading: state.users.loading,
            advancedRoleIDs: state.config[ADVANCED_ROLE_IDS],
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            organizationUserIDs: getOrganizationUserIDs(state, { organizationID: organizationID }),
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSeeOrganization: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.organizations.forbidde || state.users.forbidden || state.userRoles.forbidden || state.roles.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            push,
            loadUsers,
            loadRoles,
            loadOrganization,
            loadDomainUserRoles,
            saveUserRoleCustomMsg,
            deleteUserFromOrganization,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(UsersList))
