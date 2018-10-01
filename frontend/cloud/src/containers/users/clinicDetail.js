import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { Link, withRouter } from "react-router-dom"
import _ from "lodash"

import { ADVANCED_ROLE_IDS } from  "../../modules/config"
import { loadRoles } from "../../modules/roles"
import { makeGetClinicUserUserRoles } from "../../selectors/userRolesSelectors"
import { deleteUserFromClinic } from "../../modules/clinics"
import { loadUserUserRoles, saveUserRole, deleteUserRole } from "../../modules/userRoles"
import { SUPERADMIN_RIGHTS_RESOURCE, ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { open } from "shared/modules/alert"
import Spinner from "shared/containers/spinner"
import { confirmationDialog } from "shared/utils"

class ClinicDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.roles) {
            this.props.loadRoles()
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
        if (!nextProps.userRoles & !nextProps.userRolesLoading) {
            this.props.loadUserUserRoles(this.props.userID)
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.roles ||
            props.rolesLoading ||
            !props.userRoles ||
            props.userRolesLoading ||
            !props.clinicUserRoles ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading

        this.setState({
            loading: loading,
            userRoles: _.values(props.clinicUserRoles)
        })
    }

    newUserRole() {
        return e => {
            if (this.state.userRoles) {
                let userRoles = [
                    ...this.state.userRoles,
                    { edit: true, canSave: false, userID: this.props.userID, roleID: "", domainType: "clinic", domainID: this.props.clinicID }
                ]
                this.setState({ userRoles: userRoles })
            }
        }
    }

    editRoleID(index) {
        return e => {
            let userRoles = [...this.state.userRoles]
            userRoles[index].roleID = e.target.value
            userRoles[index].canSave = userRoles[index].roleID.length !== 0
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

    deleteUserRole(index) {
        return e => {
            confirmationDialog(
                `Click OK to confirm that you want to remove role ${this.props.roles[this.state.userRoles[index].roleID].name} from the user.`,
                () => {
                    // if there's no more clinicUserRoles after removal, remove user from clinic
                    if (_.values(this.props.clinicUserRoles).length === 1) {
                        this.props.deleteUserFromClinic(this.props.clinicID, this.props.userID)
                        this.props.history.push(`/users/${this.props.userID}/clinics`)
                    } else {
                        this.props.deleteUserRole(this.state.userRoles[index].id)
                    }
                }
            )
        }
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <Spinner />
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
                            <th scope="col">Role</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>
                        {_.map(this.state.userRoles, (userRole, i) => (
                            <tr key={userRole.id || i + 1}>
                                <td className="w-7 row-details-header-column" />
                                <th className="w-7" scope="row">
                                    {i + 1}
                                </th>
                                <td>
                                    {props.canEdit && userRole.edit ? (
                                        <select className="form-control" value={userRole.roleID || ""} onChange={this.editRoleID(i)}>
                                            <option value="">Select role</option>
                                            {_.map(
                                                _.difference(
                                                    _.map(
                                                        _.values(_.pickBy(props.roles, role => !_.includes(props.advancedRoleIDs, role.id))),
                                                        role => role.id
                                                    ),
                                                    _.map(_.values(props.clinicUserRoles), userRole => userRole.roleID)
                                                ),
                                                roleID => (
                                                    <option key={roleID} value={roleID}>
                                                        {props.roles[roleID].name}
                                                    </option>
                                                )
                                            )}
                                        </select>
                                    ) : props.canAccessRoles ? (
                                        <Link to={`/roles/${userRole.roleID}`}>{props.roles[userRole.roleID].name}</Link>
                                    ) : (
                                        props.roles[userRole.roleID].name
                                    )}
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
                                            <div className="btn-group" role="group">
                                                <button className="btn btn-link" type="button" onClick={this.deleteUserRole(i)}>
                                                    <span className="remove-link">Remove</span>
                                                </button>
                                            </div>
                                        )
                                    ) : null}
                                </td>
                            </tr>
                        ))}
                        <tr className="table-edit">
                            <td className="w-7 row-details-header-column" />
                            <td colSpan="3">
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={this.state.userRoles.length !== 0 && this.state.userRoles[this.state.userRoles.length - 1].edit ? true : null}
                                        onClick={this.newUserRole()}
                                    >
                                        Add New Role at the Clinic
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
    const getClinicUserUserRoles = makeGetClinicUserUserRoles()

    const mapStateToProps = (state, ownProps) => {
        let userID = ownProps.userID
        if (!userID) {
            userID = ownProps.match.params.userID
        }
        let clinicID = ownProps.clinicID
        if (!clinicID) {
            clinicID = ownProps.match.params.clinicID
        }

        return {
            userID: userID,
            clinicID: clinicID,
            advancedRoleIDs: state.config[ADVANCED_ROLE_IDS],
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            clinicUserRoles: getClinicUserUserRoles(state, { userID: userID, clinicID: clinicID }),
            canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canAccessRoles: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.userRoles.forbidden || state.users.forbidden || state.roles.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
            loadUserUserRoles,
            deleteUserFromClinic,
            saveUserRole,
            deleteUserRole,
            loadUserRights,
            open
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(ClinicDetail))
