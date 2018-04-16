import React from "react"
import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"
import _ from "lodash"

import { loadRoles, addRole, deleteRole } from "../../modules/roles"
import { open, COLOR_DANGER } from "shared/modules/alert"
import RoleDetail from "./detail"

import "./style.css"

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
        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.roles && !nextProps.rolesLoading) {
            this.props.loadRoles()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.roles || props.rolesLoading

        this.setState({loading: loading})
    }

    addRole = () => e => {
        if (this.state.roleName) {
            this.props.addRole(this.state.roleName)
                .then(response => {
                    if (response.id) {
                        this.props.history.push(`/roles/${response.id}`)
                    }
                })
        } else {
            this.props.open("You must enter role name", "", COLOR_DANGER)
        }
    }

    updateRoleName = () => e => {
        this.setState({ roleName: e.target.value })
    }

    deleteRole = id => e => {
        this.props.deleteRole(id)
            .then(response => {
                this.props.history.push(`/roles`)
            })
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        let i = 0
        return (
            <div id="roles">
                <div className="row">
                    <div className={props.withDetail ? "col-3" : "col-12"}>
                        <header>
                            <h1>Roles</h1>
                        </header>
                        <table className="table table-hover">
                            <thead>
                                <tr>
                                    <th scope="col">#</th>
                                    <th scope="col">Name</th>
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {map(props.roles, (role, id) => (
                                    <tr key={role.id} className={props.path.endsWith(role.id) ? "table-active" : ""}>
                                        <th scope="row">{++i}</th>
                                        <td>
                                            <Link to={`/roles/${role.id}`}>{role.name}</Link>
                                        </td>
                                        <td className="text-right">
                                            <button onClick={this.deleteRole(role.id)} className="btn btn-sm btn-light" type="button">
                                                <span className="icon_trash" />
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <div className="input-group mb-3">
                            <input
                                value={this.state.roleName}
                                onChange={this.updateRoleName()}
                                type="text"
                                className="form-control form-control-sm"
                                placeholder="Role name"
                                aria-label="Role name"
                            />
                            <div className="input-group-append">
                                <button onClick={this.addRole()} className="btn btn-sm btn-outline-secondary" type="button">
                                    Add role
                                </button>
                            </div>
                        </div>
                    </div>

                    <div className="col">
                        <Route path="/roles/:roleID" component={RoleDetail} />
                    </div>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        roles: ownProps.roles ? (state.roles.allLoaded ? _.fromPairs(_.map(ownProps.roles, roleID => [roleID, state.roles.roles[roleID]])) : undefined) : (state.roles.allLoaded ? state.roles.roles : undefined),
        rolesLoading: state.roles.loading,
        withDetail: !ownProps.match.isExact,
        path: ownProps.location.pathname,
        forbidden: state.roles.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
            addRole,
            deleteRole,
            open
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Roles)
