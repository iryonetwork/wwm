import React from "react"
import { Link, NavLink, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { joinPaths } from "shared/utils"
import { loadUser } from "../../modules/users"
import { loadRoles } from "../../modules/roles"
import { makeGetWildcardUserUserRoles } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, deleteUserRole } from "../../modules/userRoles"
import { SUPERADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import Spinner from "shared/containers/spinner"

class WildcardUserRolesList extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
    }

    componentDidMount() {
        if (!this.props.user) {
            this.props.loadUser(this.props.userID)
        }
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
        if (!nextProps.user && !nextProps.usersLoading) {
            this.props.loadUser(nextProps.userID)
        }
        if (!nextProps.roles && this.props.roles) {
            this.props.loadRoles()
        }
        if (!nextProps.userRoles && this.props.userRoles) {
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
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading
        this.setState({ loading: loading })
    }

    removeUserRole(userRoleID) {
        return e => {
            this.props.deleteUserRole(userRoleID)
            this.forceUpdate()
        }
    }

    render() {
        let props = this.props
        if (this.props.usersLoading) {
            return <Spinner />
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        let i = 0
        return (
            <div>
                <header>
                    {props.isSelf ? <h1>My Profile</h1> : <h1>Users</h1>}
                    <button onClick={this.submit} className="btn btn-primary btn-wide">
                        {props.usersUpdating ? "Saving..." : "Save"}
                    </button>
                </header>
                <h2>{props.user ? props.user.username : "New user"}</h2>
                {props.user ? (
                    <div className="navigation">
                        {props.canSeePersonal ? (
                            <NavLink exact to={props.basePath}>
                                Personal Info
                            </NavLink>
                        ) : null}
                        {props.canSeeOrganizations ? <NavLink to={joinPaths(props.basePath, "organizations")}>Organizations</NavLink> : null}
                        {props.canSeeClinics ? <NavLink to={joinPaths(props.basePath, "clinics")}>Clinics</NavLink> : null}
                        {props.canSee ? <NavLink to={joinPaths(props.basePath, "userroles")}>Wildcard Roles</NavLink> : null}
                    </div>
                ) : null}
                {this.state.loading ? (
                    <Spinner />
                ) : (
                    <div id="wildcardRoles">
                        <div className="row">
                            <div className="col-12">
                                <table className="table">
                                    <thead>
                                        <tr>
                                            <th className="w-7" scope="col">
                                                #
                                            </th>
                                            <th scope="col">Role</th>
                                            <th scope="col">Domain type</th>
                                            <th />
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {_.map(_.filter(props.wildcardUserRoles, userRole => userRole), userRole => (
                                            <tr key={userRole.id}>
                                                <th className="w-7" scope="row">
                                                    {++i}
                                                </th>
                                                <td>
                                                    {props.canEdit ? (
                                                        <Link to={`/roles/${userRole.roleID}`}>{props.roles[userRole.roleID].name}</Link>
                                                    ) : (
                                                        props.roles[userRole.roleID].name
                                                    )}
                                                </td>
                                                <td>{userRole.domainType}</td>
                                                <td className="text-right">
                                                    {props.canEdit ? (
                                                        <button onClick={this.removeUserRole(userRole.id)} className="btn btn-link" type="button">
                                                            <span className="remove-link">Remove</span>
                                                        </button>
                                                    ) : null}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        )
    }
}

const makeMapStateToProps = () => {
    const getWildcardUserUserRoles = makeGetWildcardUserUserRoles()
    const mapStateToProps = (state, ownProps) => {
        let userID = ownProps.userID
        if (!userID) {
            userID = ownProps.match.params.userID
        }
        let isSelf = state.authentication.token.sub === userID

        return {
            basePath: ownProps.home ? "/me" : `/users/${userID}`,
            isSelf: isSelf,
            userID: userID,
            user: state.users.users ? state.users.users[userID] : undefined,
            usersLoading: state.users.loading,
            roles: state.roles.roles,
            rolesLoading: state.roles.loading,
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            wildcardUserRoles: getWildcardUserUserRoles(state, { userID: userID }),
            canSee: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
            canEdit: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
            canSeePersonal: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeOrganizations: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.userRoles.forbidden || state.users.forbidden || state.roles.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser,
            loadRoles,
            loadUserUserRoles,
            deleteUserRole,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(WildcardUserRolesList))
