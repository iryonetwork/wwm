import React from "react"
import { NavLink, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import classnames from "classnames"
import { push } from "react-router-redux"

import { joinPaths } from "shared/utils"
import { ADVANCED_ROLE_IDS } from "shared/modules/config"
import { loadUser } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { loadOrganizations } from "../../modules/organizations"
import { loadLocations } from "../../modules/locations"
import { loadClinics, deleteUserFromClinic } from "../../modules/clinics"
import { makeGetUserClinicIDs, makeGetUserAllowedClinicIDs } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, saveUserRoleCustomMsg, deleteUserRole } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, SUPERADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import Spinner from "shared/containers/spinner"
import ClinicDetail from "./clinicDetail"
import { confirmationDialog } from "shared/utils"

class ClinicsList extends React.Component {
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
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.clinics) {
            this.props.loadClinics()
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
        if (!nextProps.locations && !nextProps.locationsLoading) {
            this.props.loadLocations()
        }
        if (!nextProps.clinics && !nextProps.clinicsLoading) {
            this.props.loadClinics()
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
            !props.clinics ||
            props.clinicsLoading ||
            !props.organizations ||
            props.organizationsLoading ||
            !props.locations ||
            props.locationsLoading ||
            !props.userClinicIDs ||
            !props.allowedClinicIDs ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading

        let selectedClinicID = props.clinicID
        if (!selectedClinicID) {
            selectedClinicID = props.match.params.clinicID
        }
        this.setState({
            loading: loading,
            userClinics: _.map(props.userClinicIDs, clinicID => {
                return { id: clinicID }
            }),
            selectedClinicID: selectedClinicID ? selectedClinicID : undefined
        })
    }

    newUserClinic() {
        return e => {
            if (this.state.userClinics) {
                let userClinics = [
                    ...this.state.userClinics,
                    { id: "", edit: true, canSave: false, userID: this.props.userID, roleID: "", domainType: "clinic" }
                ]
                this.setState({ userClinics: userClinics })
            }
        }
    }

    editClinicID(index) {
        return e => {
            let userClinics = [...this.state.userClinics]
            userClinics[index].id = e.target.value
            userClinics[index].canSave = userClinics[index].id.length !== 0 && userClinics[index].roleID.length !== 0
            this.setState({ userClinics: userClinics })
        }
    }

    editRoleID(index) {
        return e => {
            let userClinics = [...this.state.userClinics]
            userClinics[index].roleID = e.target.value
            userClinics[index].canSave = userClinics[index].id.length !== 0 && userClinics[index].roleID.length !== 0
            this.setState({ userClinics: userClinics })
        }
    }

    saveUserClinic(index) {
        return e => {
            let userClinics = [...this.state.userClinics]
            let userRole = {}
            userRole.userID = userClinics[index].userID
            userRole.roleID = userClinics[index].roleID
            userRole.domainType = userClinics[index].domainType
            userRole.domainID = userClinics[index].id
            userClinics[index].index = index
            userClinics[index].edit = false
            userClinics[index].saving = true

            this.props.saveUserRoleCustomMsg(userRole, "Added User to the Clinic").then(response => {
                if (response && response.domainID) {
                    this.props.history.push(`${this.props.basePath}/clinics/${response.domainID}`)
                }
            })
        }
    }

    cancelNewUserClinic(index) {
        return e => {
            let userClinics = [...this.state.userClinics]
            userClinics.splice(index, 1)
            this.setState({ userClinics: userClinics })
        }
    }

    removeUserClinic(index) {
        return e => {
            confirmationDialog(
                `Click OK to confirm that you want to remove the user from clinic ${this.props.clinics[this.state.userClinics[index].id].name}.`,
                () => {
                    this.props.deleteUserFromClinic(this.state.userClinics[index].id, this.props.userID)
                    if (this.state.selectedClinicID === this.state.userClinics[index].id) {
                        this.props.history.push(`${this.props.basePath}/clinics`)
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
                    {props.canSeeOrganizations ? <NavLink to={joinPaths(props.basePath, "organizations")}>Organizations</NavLink> : null}
                    {props.canSee ? <NavLink to={joinPaths(props.basePath, "clinics")}>Clinics</NavLink> : null}
                    {props.canSeeWildcardUserRoles ? <NavLink to={joinPaths(props.basePath, "userroles")}>Wildcard Roles</NavLink> : null}
                </div>
                {this.state.loading ? (
                    <Spinner />
                ) : (
                    <div id="clinics">
                        <div className="row">
                            <div className="col-12">
                                {this.state.userClinics.length > 0 ? (
                                    <table className="table">
                                        <thead>
                                            <tr>
                                                <th className="w-7" sscope="col">
                                                    #
                                                </th>
                                                <th className="w-20" scope="col">
                                                    Clinic
                                                </th>
                                                <th className="w-20" scope="col">
                                                    Location
                                                </th>
                                                <th scope="col">Organization</th>
                                                <th className="w-25" />
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {_.map(this.state.userClinics, (userClinic, i) => (
                                                <React.Fragment key={userClinic.id || i}>
                                                    {props.canEdit && userClinic.edit ? (
                                                        <tr
                                                            className={classnames({
                                                                "table-active": this.state.selectedClinicID === userClinic.id,
                                                                "table-edit": props.canEdit && userClinic.edit
                                                            })}
                                                        >
                                                            <th className="w-7" scope="row">
                                                                {i + 1}
                                                            </th>
                                                            <td colSpan="2">
                                                                <select className="form-control" value={userClinic.id || ""} onChange={this.editClinicID(i)}>
                                                                    <option value="">Select clinic</option>
                                                                    {_.map(
                                                                        _.difference(
                                                                            props.allowedClinicIDs,
                                                                            _.without(_.map(this.state.userClinics, clinic => clinic.id), userClinic.id)
                                                                        ),
                                                                        clinicID => (
                                                                            <option key={clinicID} value={clinicID}>
                                                                                {props.organizations[props.clinics[clinicID].organization].name} -{" "}
                                                                                {props.clinics[clinicID].name}
                                                                            </option>
                                                                        )
                                                                    )}
                                                                </select>
                                                            </td>
                                                            <td>
                                                                <select className="form-control" value={userClinic.roleID || ""} onChange={this.editRoleID(i)}>
                                                                    <option value="">Select role</option>
                                                                    {_.map(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id)), role => (
                                                                        <option key={role.id} value={role.id}>
                                                                            {role.name}
                                                                        </option>
                                                                    ))}
                                                                </select>
                                                            </td>
                                                            <td className="w-25 text-right">
                                                                <div>
                                                                    <button
                                                                        className="btn btn-secondary"
                                                                        disabled={userClinic.saving}
                                                                        type="button"
                                                                        onClick={this.cancelNewUserClinic(i)}
                                                                    >
                                                                        Remove
                                                                    </button>
                                                                    <button
                                                                        className="btn btn-primary"
                                                                        disabled={userClinic.saving || !userClinic.canSave}
                                                                        type="button"
                                                                        onClick={this.saveUserClinic(i)}
                                                                    >
                                                                        Add
                                                                    </button>
                                                                </div>
                                                            </td>
                                                        </tr>
                                                    ) : (
                                                        <tr className={this.state.selectedClinicID === userClinic.id ? "table-active" : ""}>
                                                            <th className="w-7" scope="row">
                                                                {i + 1}
                                                            </th>
                                                            <td className="w-20">{props.clinics[userClinic.id].name}</td>
                                                            <td className="w-10">
                                                                <Link to={`/locations/${props.clinics[userClinic.id].location}`}>
                                                                    {props.locations[props.clinics[userClinic.id].location].name}
                                                                </Link>
                                                            </td>
                                                            <td>
                                                                <Link to={`/organizations/${props.clinics[userClinic.id].organization}`}>
                                                                    {props.organizations[props.clinics[userClinic.id].organization].name}
                                                                </Link>
                                                            </td>
                                                            <td className="w-25 text-right">
                                                                <div>
                                                                    {this.state.selectedClinicID === userClinic.id ? (
                                                                        <button
                                                                            className="btn btn-link"
                                                                            type="button"
                                                                            onClick={() => this.props.push(`/users/${props.userID}/clinics`)}
                                                                        >
                                                                            Hide Roles
                                                                            <span className="arrow-up-icon" />
                                                                        </button>
                                                                    ) : (
                                                                        <button
                                                                            className="btn btn-link"
                                                                            type="button"
                                                                            onClick={() => this.props.push(`/users/${props.userID}/clinics/${userClinic.id}`)}
                                                                        >
                                                                            Show Roles<span className="arrow-down-icon" />
                                                                        </button>
                                                                    )}
                                                                    {props.canEdit ? (
                                                                        <button className="btn btn-link" type="button" onClick={this.removeUserClinic(i)}>
                                                                            <span className="remove-link">Remove</span>
                                                                        </button>
                                                                    ) : null}
                                                                </div>
                                                            </td>
                                                        </tr>
                                                    )}
                                                    {this.state.selectedClinicID === userClinic.id ? (
                                                        <React.Fragment>
                                                            <tr className="table-active">
                                                                <td colSpan="5" className="row-details-container">
                                                                    <ClinicDetail userID={this.props.userID} clinicID={userClinic.id} />
                                                                </td>
                                                            </tr>
                                                        </React.Fragment>
                                                    ) : null}
                                                </React.Fragment>
                                            ))}
                                        </tbody>
                                    </table>
                                ) : (
                                    <h3>User does not belong to any clinic.</h3>
                                )}
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={
                                            this.state.userClinics.length !== 0 && this.state.userClinics[this.state.userClinics.length - 1].edit ? true : null
                                        }
                                        onClick={this.newUserClinic()}
                                    >
                                        Add Current User to a Clinic
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
    const getUserClinicIDs = makeGetUserClinicIDs()
    const getUserAllowedClinicIDs = makeGetUserAllowedClinicIDs()

    const mapStateToProps = (state, ownProps) => {
        let userID = ownProps.userID
        if (!userID) {
            userID = ownProps.match.params.userID
        }
        let isSelf = state.authentication.token.sub === userID

        let selectedClinicID = ownProps.selectedClinicID
        if (!selectedClinicID) {
            selectedClinicID = ownProps.match.params.clinicID
        }

        return {
            basePath: ownProps.home ? "/me" : `/users/${userID}`,
            isSelf: isSelf,
            userID: userID,
            user: state.users.users ? state.users.users[userID] : undefined,
            usersLoading: state.users.loading,
            selectedClinicID: selectedClinicID,
            clinics: state.clinics.allLoaded ? state.clinics.clinics : undefined,
            clinicsLoading: state.clinics.loading,
            organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
            organizationsLoading: state.organizations.loading,
            locations: state.locations.allLoaded ? state.locations.locations : undefined,
            locationsLoading: state.locations.loading,
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            advancedRoleIDs: state.config[ADVANCED_ROLE_IDS],
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            userClinicIDs: getUserClinicIDs(state, { userID: userID }),
            allowedClinicIDs: getUserAllowedClinicIDs(state, { userID: userID }),
            canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSeePersonal: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeOrganizations: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeWildcardUserRoles: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden:
                state.userRoles.forbidden || state.users.forbidden || state.organizations.forbidden || state.clinics.forbidden | state.locations.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            push,
            loadUser,
            loadRoles,
            loadOrganizations,
            loadLocations,
            loadClinics,
            deleteUserFromClinic,
            loadUserUserRoles,
            saveUserRoleCustomMsg,
            loadUserRights,
            deleteUserRole
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(ClinicsList))
