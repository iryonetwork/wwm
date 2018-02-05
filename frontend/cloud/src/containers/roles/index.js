import React from "react"
import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import map from "lodash/map"

import { loadRoles } from "../../modules/roles"
import RoleDetail from "./detail"

import "./style.css"

class Roles extends React.Component {
    componentDidMount() {
        this.props.loadRoles()
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
                                </tr>
                            </thead>
                            <tbody>
                                {map(props.roles, (role, id) => (
                                    <tr
                                        key={role.id}
                                        className={
                                            props.path.endsWith(role.id)
                                                ? "table-active"
                                                : ""
                                        }
                                    >
                                        <th scope="row">{++i}</th>
                                        <td>
                                            <Link to={`/roles/${role.id}`}>
                                                {role.name}
                                            </Link>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
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
            loadRoles
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Roles)
