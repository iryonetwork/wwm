import React from "react"
import { Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadRoles } from "../../modules/roles"
import { makeGetWildcardUserUserRoles } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, deleteUserRole } from "../../modules/userRoles"

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
        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.roles && this.props.roles) {
            this.props.loadRoles()
        }
        if (!nextProps.userRoles && this.props.userRoles) {
            this.props.loadUserUserRoles(this.props.userID)
        }
        this.determineState(nextProps)
    }

    determineState(props) {
        this.setState({})
    }

    removeUserRole = userRoleID => e => {
        this.props.deleteUserRole(userRoleID)
        this.forceUpdate()
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (props.loading) {
            return <div>Loading...</div>
        }
        let i = 0
        return (
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
                                <Link to={`/roles/${userRole.roleID}`}>{props.roles[userRole.roleID].name}</Link>
                            </td>
                            <td>{userRole.domainType}</td>
                            <td className="text-right">
                                <button onClick={this.removeUserRole(userRole.id)} className="btn btn-sm btn-light" type="button">
                                    <span className="icon_trash" />
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
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
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            wildcardUserRoles: getWildcardUserUserRoles(state, {userID: userID}),
            loading: state.userRoles.loading || state.roles.loading,
            forbidden: state.userRoles.forbidden || state.users.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
            loadUserUserRoles,
            deleteUserRole
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(WildcardUserRolesList))
