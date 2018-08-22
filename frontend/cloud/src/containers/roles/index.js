import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"
import _ from "lodash"
import classnames from "classnames"
import { push } from "react-router-redux"

import { loadRoles, addRole, deleteRole } from "../../modules/roles"
import { SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { open } from "shared/modules/alert"
import RoleDetail from "./detail"
import { confirmationDialog } from "shared/utils"

import "../../styles/style.css"

class Roles extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            roleName: "",
            loading: true
        }
    }

    componentDidMount() {
        if (!this.props.roles) {
            this.props.loadRoles()
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
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }
        if (nextProps.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.roles || props.rolesLoading || props.canEdit === undefined || props.canSee === undefined || props.validationsLoading

        let selectedRoleID = props.roleID
        if (!selectedRoleID) {
            selectedRoleID = props.match.params.roleID
        }

        this.setState({
            loading: loading,
            roles: _.values(props.roles),
            selectedRoleID: selectedRoleID || undefined
        })
    }

    newRole() {
        return e => {
            if (this.state.roles) {
                let roles = [...this.state.roles, { id: "", edit: true, canSave: false, name: "" }]
                this.setState({
                    roles: roles,
                    edit: true
                })
            }
        }
    }

    editRoleName(index) {
        return e => {
            let roles = [...this.state.roles]
            roles[index].name = e.target.value
            roles[index].canSave = roles[index].name.length !== 0
            this.setState({ roles: roles })
        }
    }

    saveRole(index) {
        return e => {
            let roles = [...this.state.roles]

            roles[index].edit = false
            roles[index].saving = true

            this.props.addRole(roles[index].name).then(response => {
                if (response && response.id) {
                    this.props.history.push(`/roles/${response.id}`)
                }
            })
        }
    }

    cancelNewRole(index) {
        return e => {
            let roles = [...this.state.roles]
            roles.splice(index, 1)
            this.setState({
                roles: roles,
                edit: false
            })
        }
    }

    deleteRole(index) {
        return e => {
            confirmationDialog(`Click OK to confirm that you want to remove role ${this.state.roles[index].name}.`, () => {
                this.props.deleteRole(this.state.roles[index].id).then(response => {
                    this.props.history.push(`/roles`)
                })
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
            <div id="roles">
                <div className="row">
                    <div className={props.withDetail ? "col-3" : "col-12"}>
                        <header>
                            <h1>Roles</h1>
                        </header>
                        <table className="table">
                            <thead>
                                <tr>
                                    <th className="w-7" scope="col">
                                        #
                                    </th>
                                    <th scope="col">Name</th>
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {map(this.state.roles, (role, i) => (
                                    <React.Fragment key={role.id || i}>
                                        <tr
                                            className={classnames({
                                                "table-active": role.id === this.state.selectedRoleID,
                                                "table-edit": props.canEdit && role.edit
                                            })}
                                        >
                                            <th className="w-7" sscope="row">
                                                {i + 1}
                                            </th>
                                            {props.canEdit && role.edit ? (
                                                <React.Fragment>
                                                    <td>
                                                        <input
                                                            type="text"
                                                            value={role.name || ""}
                                                            onChange={this.editRoleName(i)}
                                                            className="form-control"
                                                            placeholder="Role Name"
                                                            aria-label="Role Name"
                                                        />
                                                    </td>
                                                    <td className="text-right">
                                                        {props.canEdit ? (
                                                            <div>
                                                                <button
                                                                    className="btn btn-secondary"
                                                                    disabled={role.saving}
                                                                    type="button"
                                                                    onClick={this.cancelNewRole(i)}
                                                                >
                                                                    Remove
                                                                </button>
                                                                <button
                                                                    className="btn btn-primary"
                                                                    disabled={role.saving || !role.canSave}
                                                                    type="button"
                                                                    onClick={this.saveRole(i)}
                                                                >
                                                                    Add
                                                                </button>
                                                            </div>
                                                        ) : null}
                                                    </td>
                                                </React.Fragment>
                                            ) : (
                                                <React.Fragment>
                                                    <td>
                                                        <Link to={`/roles/${role.id}`}>{role.name}</Link>
                                                    </td>
                                                    <td className="text-right">
                                                        {role.id === this.state.selectedRoleID ? (
                                                            <button className="btn btn-link" type="button" onClick={() => this.props.push("/roles")}>
                                                                Hide ACL
                                                                <span className="arrow-up-icon" />
                                                            </button>
                                                        ) : (
                                                            <button className="btn btn-link" type="button" onClick={() => this.props.push(`/roles/${role.id}`)}>
                                                                Show ACL<span className="arrow-down-icon" />
                                                            </button>
                                                        )}
                                                        {props.canEdit ? (
                                                            <button onClick={this.deleteRole(i)} className="btn btn-link" type="button">
                                                                <span className="remove-link">Remove</span>
                                                            </button>
                                                        ) : null}
                                                    </td>
                                                </React.Fragment>
                                            )}
                                        </tr>
                                        {role.id === this.state.selectedRoleID ? (
                                            <React.Fragment>
                                                <tr className="table-active">
                                                    <td colSpan="5" className="row-details-container">
                                                        <RoleDetail roleID={role.id} />
                                                    </td>
                                                </tr>
                                            </React.Fragment>
                                        ) : null}
                                    </React.Fragment>
                                ))}
                            </tbody>
                        </table>
                        {props.canEdit ? (
                            <button type="button" className="btn btn-link" disabled={this.state.edit ? true : null} onClick={this.newRole()}>
                                Add Role
                            </button>
                        ) : null}
                    </div>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        roles: ownProps.roles
            ? state.roles.allLoaded
                ? _.fromPairs(_.map(ownProps.roles, roleID => [roleID, state.roles.roles[roleID]]))
                : undefined
            : state.roles.allLoaded
                ? state.roles.roles
                : undefined,
        rolesLoading: state.roles.loading,
        withDetail: !ownProps.match.isExact,
        path: ownProps.location.pathname,
        canEdit: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
        canSee: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading,
        forbidden: state.roles.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
            addRole,
            deleteRole,
            loadUserRights,
            open,
            push
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Roles)
