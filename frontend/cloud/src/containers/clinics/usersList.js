import React from "react"
import { Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import { push } from "react-router-redux"

import { ADVANCED_ROLE_IDS } from "../../modules/config"
import { loadUsers } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { deleteUserFromClinic } from "../../modules/clinics"
import { saveUserRoleCustomMsg, loadDomainUserRoles } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { makeGetClinicUserIDs, makeGetOrganizationUserIDs } from "../../selectors/userRolesSelectors"
import { getName } from "../../utils/user"
import UserDetail from "./userDetail"
import { confirmationDialog } from "shared/utils"

class UsersList extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.clinic) {
            this.props.loadClinic(this.props.clinicID)
        }
        if (!this.props.users) {
            this.props.loadUsers()
        }
        if (!this.props.roles) {
            this.props.loadRoles()
        }
        if (!this.props.userRoles) {
            this.props.loadDomainUserRoles("clinic", this.props.clinicID)
        }
        if (this.props.organizationID && !this.props.clinicsOrganizationUserRoles) {
            this.props.loadDomainUserRoles("organization", this.props.organizationID)
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.clinic && !nextProps.clinicsLoading) {
            this.props.loadClinic(this.props.clinicID)
        }
        if (!nextProps.users && !nextProps.usersLoading) {
            this.props.loadUsers()
        }
        if (!nextProps.roles && !nextProps.rolesLoading) {
            this.props.loadRoles()
        }
        if (!nextProps.userRolesLoading) {
            if (!nextProps.userRoles) {
                this.props.loadDomainUserRoles("clinic", this.props.clinicID)
            }
            if (!nextProps.clinicsOrganizationUserRoles && !nextProps.userRolesLoading) {
                this.props.loadDomainUserRoles("organization", this.props.organizationID)
            }
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
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading
        this.setState({ loading: loading })

        if (!loading) {
            this.setState({
                clinicUsers: _.map(props.clinicUserIDs ? props.clinicUserIDs : [], userID => {
                    return props.users[userID]
                }),
                selectedUserID: props.selectedUserID
            })
        }
    }

    newUser() {
        return e => {
            if (this.state.clinicUsers) {
                let clinicUsers = [
                    ...this.state.clinicUsers,
                    { id: "", edit: true, canSave: false, domainType: "clinic", domainID: this.props.clinicID, userID: "", roleID: "" }
                ]
                this.setState({ clinicUsers: clinicUsers })
            }
        }
    }

    editUserID(index) {
        return e => {
            let clinicUsers = [...this.state.clinicUsers]
            clinicUsers[index].userID = e.target.value
            clinicUsers[index].canSave = clinicUsers[index].userID.length !== 0 && clinicUsers[index].roleID.length !== 0
            this.setState({ clinicUsers: clinicUsers })
        }
    }

    editRoleID(index) {
        return e => {
            let clinicUsers = [...this.state.clinicUsers]
            clinicUsers[index].roleID = e.target.value
            clinicUsers[index].canSave = clinicUsers[index].userID.length !== 0 && clinicUsers[index].roleID.length !== 0
            this.setState({ clinicUsers: clinicUsers })
        }
    }

    saveUser(index) {
        return e => {
            let clinicUsers = [...this.state.clinicUsers]

            clinicUsers[index].edit = false
            clinicUsers[index].saving = true

            this.props.saveUserRoleCustomMsg(this.state.clinicUsers[index], "Added User to the Clinic")
        }
    }

    cancelNewUser(index) {
        return e => {
            let clinicUsers = [...this.state.clinicUsers]
            clinicUsers.splice(index, 1)
            this.setState({ clinicUsers: clinicUsers })
        }
    }

    removeUser(index) {
        return e => {
            confirmationDialog(`Click OK to confirm that you want to remove user ${this.state.clinicUsers[index].username} from the clinic.`, () => {
                this.props.deleteUserFromClinic(this.props.clinicID, this.state.clinicUsers[index].id)
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

        return (
            <div>
                <table className="table">
                    <thead>
                        <tr>
                            <td className="w-7 row-details-header-column">
                                <span className="row-details-icon" />
                            </td>
                            <th className="w-7" scope="col">
                                #
                            </th>
                            <th className="w-15" scope="col">
                                Username
                            </th>
                            <th className="w-20" scope="col">
                                Name
                            </th>
                            <th scope="col">Email</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {_.map(this.state.clinicUsers, (user, i) => (
                            <React.Fragment key={user.id || i}>
                                {props.canEdit && user.edit ? (
                                    <tr>
                                        <td className="w-7 row-details-header-column" />
                                        <th className="w-7" scope="row">
                                            {i + 1}
                                        </th>
                                        <td className="w-35" colSpan="2">
                                            <select className="form-control" value={user.userID || ""} onChange={this.editUserID(i)}>
                                                <option>Select user</option>
                                                {_.map(_.difference(props.allowedClinicUserIDs, props.clinicUserIDs), userID => (
                                                    <option key={userID} value={userID || ""}>
                                                        {props.users[userID].username} - {getName(props.users[userID])} ({props.users[userID].email})
                                                    </option>
                                                ))}
                                            </select>
                                        </td>
                                        <td>
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
                                                <button className="btn btn-secondary" disabled={user.saving} type="button" onClick={this.cancelNewUser(i)}>
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
                                        <td className="w-7 row-details-header-column" />
                                        <th className="w-7" scope="row">
                                            {i + 1}
                                        </th>
                                        <td className="w-15">
                                            {this.state.selectedUserID === user.id ? (
                                                <Link to={`/users/${user.id}`}>{user.username}</Link>
                                            ) : (
                                                <Link to={`/clinics/${props.clinicID}/users/${user.id}`}>{user.username}</Link>
                                            )}
                                        </td>
                                        <td className="w-20">{getName(user)}</td>
                                        <td>{user.email}</td>
                                        <td className="text-right">
                                            <div>
                                                {this.state.selectedUserID === user.id ? (
                                                    <button
                                                        className="btn btn-link"
                                                        type="button"
                                                        onClick={() => this.props.push(`/clinics/${props.clinicID}`)}
                                                    >
                                                        Hide Roles
                                                        <span className="arrow-up-icon" />
                                                    </button>
                                                ) : (
                                                    <button
                                                        className="btn btn-link"
                                                        type="button"
                                                        onClick={() => this.props.push(`/clinics/${props.clinicID}/users/${user.id}`)}
                                                    >
                                                        Show Roles
                                                        <span className="arrow-down-icon" />
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
                                        <td colSpan="7" className="row-details-container">
                                            <UserDetail clinicID={props.clinicID} userID={user.id} />
                                        </td>
                                    </tr>
                                ) : null}
                            </React.Fragment>
                        ))}
                        <tr className="table-edit">
                            <td className="w-7 row-details-header-column" />
                            <td colSpan="7">
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={
                                            this.state.clinicUsers.length !== 0 && this.state.clinicUsers[this.state.clinicUsers.length - 1].edit ? true : null
                                        }
                                        onClick={this.newUser()}
                                    >
                                        Add User to the Clinic
                                    </button>
                                ) : null}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        )
    }
}

const makeMapStateToProps = () => {
    const getClinicUserIDs = makeGetClinicUserIDs()
    const getOrganizationUserIDs = makeGetOrganizationUserIDs()

    const mapStateToProps = (state, ownProps) => {
        let clinicID = ownProps.clinicID
        if (!clinicID) {
            clinicID = ownProps.match.params.clinicID
        }
        let selectedUserID = ownProps.userID
        if (!selectedUserID) {
            selectedUserID = ownProps.match.params.userID
        }

        let clinic = state.clinics.clinics ? state.clinics.clinics[clinicID] : undefined
        let organizationID = clinic ? clinic.organization : undefined

        return {
            clinicID: clinicID,
            organizationID: organizationID,
            selectedUserID: selectedUserID ? selectedUserID : undefined,
            clinic: clinic,
            clinicsLoading: state.clinics.loading,
            userRoles:
                state.userRoles.domainUserRoles && state.userRoles.domainUserRoles["clinic"] && state.userRoles.domainUserRoles["clinic"][clinicID]
                    ? state.userRoles.domainUserRoles["clinic"][clinicID]
                    : undefined,
            clinicsOrganizationUserRoles:
                organizationID &&
                state.userRoles.domainUserRoles &&
                state.userRoles.domainUserRoles["organization"] &&
                state.userRoles.domainUserRoles["organization"][organizationID]
                    ? state.userRoles.domainUserRoles["organization"][organizationID]
                    : undefined,
            userRolesLoading: state.userRoles.loading,
            users: state.users.allLoaded ? state.users.users : undefined,
            usersLoading: state.users.loading,
            advancedRoleIDs: state.config[ADVANCED_ROLE_IDS],
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            clinicUserIDs: getClinicUserIDs(state, { clinicID: clinicID }),
            allowedClinicUserIDs: getOrganizationUserIDs(state, { organizationID: organizationID }),
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.clinics.forbidden || state.users.forbidden || state.userRoles.forbidden || state.roles.forbidden
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
            loadDomainUserRoles,
            saveUserRoleCustomMsg,
            deleteUserFromClinic,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(UsersList))
