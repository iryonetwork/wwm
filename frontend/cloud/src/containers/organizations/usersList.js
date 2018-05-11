import React from "react"
import { Route, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { ADVANCED_ROLE_IDS } from "shared/modules/config"
import { loadUsers } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadOrganization, deleteUserFromOrganization } from "../../modules/organizations"
import { saveUserRoleCustomMsg, loadDomainUserRoles } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { makeGetOrganizationUserIDs } from "../../selectors/userRolesSelectors"
import { getName } from "../../utils/user"
import UserDetail from "./userDetail"

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
        if (!nextProps.organization && !nextProps.organizationsLoading) {
            this.props.loadOrganization(nextProps.organizationID)
        }
        if (!nextProps.userRoles && !nextProps.userRolesLoading) {
            this.props.loadDomainUserRoles("organization", this.props.organizationID)
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
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

    newUser = () => e => {
        if (this.state.organizationUsers) {
            let organizationUsers = [
                ...this.state.organizationUsers,
                { id: "", edit: true, canSave: false, domainType: "organization", domainID: this.props.organizationID, userID: "", roleID: "" }
            ]
            this.setState({ organizationUsers: organizationUsers })
        }
    }

    editUserID = index => e => {
        let organizationUsers = [...this.state.organizationUsers]
        organizationUsers[index].userID = e.target.value
        organizationUsers[index].canSave = organizationUsers[index].userID.length !== 0 && organizationUsers[index].roleID.length !== 0
        this.setState({ organizationUsers: organizationUsers })
    }

    editRoleID = index => e => {
        let organizationUsers = [...this.state.organizationUsers]
        organizationUsers[index].roleID = e.target.value
        organizationUsers[index].canSave = organizationUsers[index].userID.length !== 0 && organizationUsers[index].roleID.length !== 0
        this.setState({ organizationUsers: organizationUsers })
    }

    saveUser = index => e => {
        let organizationUsers = [...this.state.organizationUsers]

        organizationUsers[index].edit = false
        organizationUsers[index].saving = true

        this.props.saveUserRoleCustomMsg(this.state.organizationUsers[index], "Added user to organization")
    }

    cancelNewUser = index => e => {
        let organizationUsers = [...this.state.organizationUsers]
        organizationUsers.splice(index, 1)
        this.setState({ organizationUsers: organizationUsers })
    }

    removeUser = userID => e => {
        this.props.deleteUserFromOrganization(this.props.organizationID, userID)
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
            <div id="users">
                <h2>Users</h2>
                <div className="row">
                    <div className={this.state.selectedUserID ? "col-8" : "col-12"}>
                        <table className="table table-hover text-center">
                            <thead>
                                <tr>
                                    <th scope="col">#</th>
                                    <th scope="col">Username</th>
                                    <th scope="col">Name</th>
                                    <th scope="col">Email</th>
                                    <th />
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {_.map(this.state.organizationUsers, (user, i) => {
                                    return props.canEdit && user.edit ? (
                                        <tr key={i}>
                                            <th scope="row">{i + 1}</th>
                                            <td colSpan="2">
                                                <select className="form-control form-control-sm" value={user.userID} onChange={this.editUserID(i)}>
                                                    <option>Select user</option>
                                                    {_.map(_.difference(_.map(_.values(props.users), user => user.id), props.organizationUserIDs), userID => (
                                                        <option key={userID} value={userID}>
                                                            {props.users[userID].username} - {getName(props.users[userID])} ({props.users[userID].email})
                                                        </option>
                                                    ))}
                                                </select>
                                            </td>
                                            <td colSpan="2">
                                                <select className="form-control form-control-sm" value={user.roleID} onChange={this.editRoleID(i)}>
                                                    <option>Select role</option>
                                                    {_.map(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id)), role => (
                                                        <option key={role.id} value={role.id}>
                                                            {role.name}
                                                        </option>
                                                    ))}
                                                </select>
                                            </td>
                                            <td className="text-right">
                                                <div className="btn-group" role="group">
                                                    <button
                                                        className="btn btn-sm btn-light"
                                                        disabled={user.saving}
                                                        type="button"
                                                        onClick={this.cancelNewUser(i)}
                                                    >
                                                        <span className="icon_close" />
                                                    </button>
                                                    <button
                                                        className="btn btn-sm btn-light"
                                                        disabled={user.saving || !user.canSave}
                                                        type="button"
                                                        onClick={this.saveUser(i)}
                                                    >
                                                        <span className="icon_floppy" />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    ) : (
                                        <tr key={user.id} className={this.state.selectedUserID === user.id ? "table-active" : ""}>
                                            <th scope="row">{i + 1}</th>
                                            <td>
                                                {this.state.selectedUserID === user.id ? (
                                                    <Link to={`/users/${user.id}`}>{user.username}</Link>
                                                ) : (
                                                    <Link to={`/organizations/${props.organizationID}/users/${user.id}`}>{user.username}</Link>
                                                )}
                                            </td>
                                            <td>{getName(user)}</td>
                                            <td>{user.email}</td>
                                            <td />
                                            <td className="text-right">
                                                <div className="btn-group" role="group">
                                                    <button onClick={this.removeUser(user.id)} className="btn btn-sm btn-light" type="button">
                                                        <span className="icon_trash" />
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    )
                                })}
                            </tbody>
                        </table>
                        {props.canEdit ? (
                            <button
                                type="button"
                                className="btn btn-sm btn-outline-primary col"
                                disabled={
                                    this.state.organizationUsers.length !== 0 && this.state.organizationUsers[this.state.organizationUsers.length - 1].edit
                                        ? true
                                        : null
                                }
                                onClick={this.newUser()}
                            >
                                Add user
                            </button>
                        ) : null}
                    </div>
                    <div className="col">
                        <Route path="/organizations/:organizationID/users/:userID" component={UserDetail} />
                    </div>
                </div>
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
            validationsLoading: state.validations.loading,
            forbidden: state.organizations.forbidde || state.users.forbidden || state.userRoles.forbidden || state.roles.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
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
