import React from "react"
import { Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadRoles } from "../../modules/roles"
import { makeGetWildcardUserUserRoles } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, deleteUserRole } from "../../modules/userRoles"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"

class WildcardUserRolesList extends React.Component {
    constructor(props) {
        super(props)
        this.state = {}
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

    removeUserRole = userRoleID => e => {
        this.props.deleteUserRole(userRoleID)
        this.forceUpdate()
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        let i = 0
        return (
            <div id="wildcardRoles">
                <h2>Wildcard roles</h2>
                <div className="row">
                    <div className="col-12">
                        <table className="table table-hover">
                            <thead>
                                <tr>
                                    <th scope="col">#</th>
                                    <th scope="col">Role</th>
                                    <th scope="col">Domain type</th>
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {_.map(_.filter(props.wildcardUserRoles, userRole => userRole), userRole => (
                                    <tr key={userRole.id}>
                                        <th scope="row">{++i}</th>
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
                                                <button onClick={this.removeUserRole(userRole.id)} className="btn btn-sm btn-light" type="button">
                                                    <span className="icon_trash" />
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

        return {
            userID: userID,
            roles: state.roles.roles,
            rolesLoading: state.roles.loading,
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            wildcardUserRoles: getWildcardUserUserRoles(state, { userID: userID }),
            canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
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
            deleteUserRole,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(WildcardUserRolesList))
