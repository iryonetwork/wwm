import React from "react"
import { Route, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { ADVANCED_ROLE_IDS } from "shared/modules/config"
import { loadRoles } from "../../modules/roles"
import { loadOrganizations, clearOrganizations, deleteUserFromOrganization } from "../../modules/organizations"
import { makeGetUserOrganizationIDs } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, saveUserRoleCustomMsg, deleteUserRole } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import OrganizationDetail from "./organizationDetail"

class OrganizationsList extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
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

    newUserOrganization = () => e => {
        if (this.state.userOrganizations) {
            let userOrganizations = [
                ...this.state.userOrganizations,
                { id: "", edit: true, canSave: false, userID: this.props.userID, roleID: "", domainType: "organization" }
            ]
            this.setState({ userOrganizations: userOrganizations })
        }
    }

    editOrganizationID = index => e => {
        let userOrganizations = [...this.state.userOrganizations]
        userOrganizations[index].id = e.target.value
        userOrganizations[index].canSave = userOrganizations[index].id.length !== 0 && userOrganizations[index].roleID.length !== 0
        this.setState({ userOrganizations: userOrganizations })
    }

    editRoleID = index => e => {
        let userOrganizations = [...this.state.userOrganizations]
        userOrganizations[index].roleID = e.target.value
        userOrganizations[index].canSave = userOrganizations[index].id.length !== 0 && userOrganizations[index].roleID.length !== 0
        this.setState({ userOrganizations: userOrganizations })
    }

    saveUserOrganization = index => e => {
        let userOrganizations = [...this.state.userOrganizations]
        let userRole = {}
        userRole.userID = userOrganizations[index].userID
        userRole.roleID = userOrganizations[index].roleID
        userRole.domainType = userOrganizations[index].domainType
        userRole.domainID = userOrganizations[index].id
        userOrganizations[index].index = index
        userOrganizations[index].edit = false
        userOrganizations[index].saving = true

        this.props.saveUserRoleCustomMsg(userRole, "Added user to organization").then(response => {
            if (response && response.domainID) {
                this.props.history.push(`/users/${this.props.userID}/organizations/${response.domainID}`)
            }
        })
    }

    cancelNewUserOrganization = index => e => {
        let userOrganizations = [...this.state.userOrganizations]
        userOrganizations.splice(index, 1)
        this.setState({ userOrganizations: userOrganizations })
    }

    removeUserOrganization = organizationID => e => {
        this.props.deleteUserFromOrganization(organizationID, this.props.userID)
        if (this.state.selectedOrganizationID === organizationID) {
            this.props.history.push(`/users/${this.props.userID}`)
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
            <div id="organizations">
                <h2>Organizations</h2>
                <div className="row">
                    <div className={this.state.selectedOrganizationID ? "col-4" : "col-12"}>
                        <table className="table table-hover text-center">
                            <thead>
                                <tr>
                                    <th scope="col">#</th>
                                    <th scope="col">Organization name</th>
                                    <th />
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {_.map(this.state.userOrganizations, (userOrganization, i) => (
                                    <tr
                                        key={userOrganization.id || i}
                                        className={this.state.selectedOrganizationID === userOrganization.id ? "table-active" : ""}
                                    >
                                        <th scope="row">{i + 1}</th>
                                        <td>
                                            {props.canEdit && userOrganization.edit ? (
                                                <select
                                                    className="form-control form-control-sm"
                                                    value={userOrganization.id || ""}
                                                    onChange={this.editOrganizationID(i)}
                                                >
                                                    <option value="">Select organization</option>
                                                    {_.map(
                                                        _.difference(
                                                            _.map(_.values(props.organizations), organization => organization.id),
                                                            _.without(_.map(this.state.userOrganizations, organization => organization.id), userOrganization.id)
                                                        ),
                                                        organizationID => (
                                                            <option key={organizationID} value={organizationID}>
                                                                {props.organizations[organizationID].name}
                                                            </option>
                                                        )
                                                    )}
                                                </select>
                                            ) : this.state.selectedOrganizationID === userOrganization.id ? (
                                                <Link to={`/organizations/${userOrganization.id}`}>{props.organizations[userOrganization.id].name}</Link>
                                            ) : (
                                                <Link to={`/users/${props.userID}/organizations/${userOrganization.id}`}>
                                                    {props.organizations[userOrganization.id].name}
                                                </Link>
                                            )}
                                        </td>
                                        <td>
                                            {props.canEdit && userOrganization.edit ? (
                                                <select className="form-control form-control-sm" value={userOrganization.roleID || ""} onChange={this.editRoleID(i)}>
                                                    <option value="">Select role</option>
                                                    {_.map(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id)), role => (
                                                        <option key={role.id} value={role.id}>
                                                            {role.name}
                                                        </option>
                                                    ))}
                                                </select>
                                            ) : null}
                                        </td>
                                        <td className="text-right">
                                            {props.canEdit ? (
                                                userOrganization.edit ? (
                                                    <div className="btn-group" role="group">
                                                        <button
                                                            className="btn btn-sm btn-light"
                                                            disabled={userOrganization.saving}
                                                            type="button"
                                                            onClick={this.cancelNewUserOrganization(i)}
                                                        >
                                                            <span className="icon_close" />
                                                        </button>
                                                        <button
                                                            className="btn btn-sm btn-light"
                                                            disabled={userOrganization.saving || !userOrganization.canSave}
                                                            type="button"
                                                            onClick={this.saveUserOrganization(i)}
                                                        >
                                                            <span className="icon_floppy" />
                                                        </button>
                                                    </div>
                                                ) : (
                                                    <div className="btn-group" role="group">
                                                        <button
                                                            className="btn btn-sm btn-light"
                                                            type="button"
                                                            onClick={this.removeUserOrganization(userOrganization.id)}
                                                        >
                                                            <span className="icon_trash" />
                                                        </button>
                                                    </div>
                                                )
                                            ) : null}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        {props.canEdit ? (
                            <button
                                type="button"
                                className="btn btn-sm btn-outline-primary col"
                                disabled={
                                    this.state.userOrganizations.length !== 0 && this.state.userOrganizations[this.state.userOrganizations.length - 1].edit
                                        ? true
                                        : null
                                }
                                onClick={this.newUserOrganization()}
                            >
                                Add current user to organization
                            </button>
                        ) : null}
                    </div>
                    <div className="col">
                        <Route path="/users/:userID/organizations/:organizationID" component={OrganizationDetail} />
                    </div>
                </div>
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
        let selectedOrganizationID = ownProps.organizationID
        if (!selectedOrganizationID) {
            selectedOrganizationID = ownProps.match.params.organizationID
        }

        return {
            userID: userID,
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
            validationsLoading: state.validations.loading,
            forbidden: state.userRoles.forbidden || state.users.forbidden || state.organizations.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
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
