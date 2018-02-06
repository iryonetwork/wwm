import React from "react"
import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"

import { loadRoles, addRole, deleteRole } from "../../modules/roles"
import { open, COLOR_DANGER } from "../../modules/alert"
import RoleDetail from "./detail"

import "./style.css"

class Roles extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            roleName: ""
        }
    }
    componentDidMount() {
        this.props.loadRoles()
    }

    addRole = () => e => {
        if (this.state.roleName) {
            this.props.addRole(this.state.roleName)
        } else {
            this.props.open("You must enter role name", "", COLOR_DANGER)
        }
    }

    updateRoleName = () => e => {
        this.setState({ roleName: e.target.value })
    }

    deleteRole = id => e => {
        this.props.deleteRole(id)
    }

    render() {
        let props = this.props
        if (props.loading) {
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
                        <Route path="/roles/:id" component={RoleDetail} />
                    </div>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        roles: state.roles.roles || {},
        loading: state.roles.loading,
        withDetail: !ownProps.match.isExact,
        path: ownProps.location.pathname
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
