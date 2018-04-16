import React from "react"
import { Route, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadRoles } from "../../modules/roles"
import { loadOrganizations } from "../../modules/organizations"
import { loadClinics, deleteUserFromClinic } from "../../modules/clinics"
import { makeGetUserClinicIDs, makeGetUserAllowedClinicIDs } from "../../selectors/userRolesSelectors"
import { loadUserUserRoles, saveUserRoleCustomMsg, deleteUserRole } from "../../modules/userRoles"
import ClinicDetail from "./clinicDetail"

class ClinicsList extends React.Component {
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
        if (!this.props.clinics) {
            this.props.loadClinics()
        }
        if (!this.props.userRoles) {
            this.props.loadUserUserRoles(this.props.userID)
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
        if (!nextProps.clinics && !nextProps.clinicsLoading) {
            this.props.loadClinics()
        }
        if (!nextProps.userRoles && !nextProps.userRolesLoading) {
            this.props.loadUserUserRoles(this.props.userID)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.userRoles || props.userRolesLoading || !props.roles || props.rolesLoading || !props.clinics || props.clinicsLoading || !props.organizations || props.organizationsLoading || !props.userClinicIDs || !props.allowedClinicIDs
        let selectedClinicID = props.clinicID
        if (!selectedClinicID) {
            selectedClinicID = props.match.params.clinicID
        }
        this.setState({
            loading: loading,
            userClinics: _.map(props.userClinicIDs, clinicID => {return {"id": clinicID}}),
            selectedClinicID: selectedClinicID ? selectedClinicID : undefined
        })
    }

   newUserClinic = () => e => {
        if (this.state.userClinics) {
            let userClinics = [...this.state.userClinics, { id: "", edit: true, canSave: false, userID: this.props.userID, roleID: "", domainType: "clinic" }]
            this.setState({ userClinics: userClinics })
        }
    }

    editClinicID = index => e => {
        let userClinics = [...this.state.userClinics]
        userClinics[index].id = e.target.value
        userClinics[index].canSave = (userClinics[index].id.length !== 0) && (userClinics[index].roleID.length !== 0)
        this.setState({ userClinics: userClinics })
    }

    editRoleID = index => e => {
        let userClinics = [...this.state.userClinics]
        userClinics[index].roleID = e.target.value
        userClinics[index].canSave = (userClinics[index].id.length !== 0) && (userClinics[index].roleID.length !== 0)
        this.setState({ userClinics: userClinics })
    }

    saveUserClinic = index => e => {
        let userClinics = [...this.state.userClinics]
        let userRole = {}
        userRole.userID = userClinics[index].userID
        userRole.roleID = userClinics[index].roleID
        userRole.domainType = userClinics[index].domainType
        userRole.domainID = userClinics[index].id
        userClinics[index].index = index
        userClinics[index].edit = false
        userClinics[index].saving = true

        this.props.saveUserRoleCustomMsg(userRole, "Added user to clinic")
            .then(response => {
                if (response && response.domainID) {
                    this.props.history.push(`/users/${this.props.userID}/clinics/${response.domainID}`)
                }
            })
    }

    cancelNewUserClinic = index => e => {
        let userClinics = [...this.state.userClinics]
        userClinics.splice(index, 1)
        this.setState({ userClinics: userClinics })
    }

    removeUserClinic = clinicID => e => {
        this.props.deleteUserFromClinic(clinicID, this.props.userID)
        if (this.state.selectedClinicID === clinicID) {
            this.props.history.push(`/users/${this.props.userID}`)
        }
    }

    render() {
        let props = this.props
        if (props.forbidden) {
            return null
        }
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        return (
            <div id="clinics">
                <div className="row">
                    <div className={this.state.selectedClinicID ? "col-4" : "col-12"}>
                        <table className="table table-hover">
                            <thead>
                                <tr>
                                    <th scope="col">#</th>
                                    <th scope="col">Clinic name</th>
                                    <th />
                                    <th />
                                </tr>
                            </thead>
                            <tbody>
                                {_.map(this.state.userClinics, (userClinic, i) => (
                                    <tr key={userClinic.id || i} className={(this.state.selectedClinicID === userClinic.id) ? "table-active" : ""}>
                                        <th scope="row">{i+1}</th>
                                        <td>
                                        {userClinic.edit ? (
                                            <select className="form-control form-control-sm" value={userClinic.id} onChange={this.editClinicID(i)}>
                                                <option value="">Select clinic</option>
                                                {_.map(_.difference(props.allowedClinicIDs, _.without(_.map(this.state.userClinics, clinic => clinic.id), userClinic.id)), clinicID => (
                                                    <option key={clinicID} value={clinicID}>
                                                        {props.organizations[props.clinics[clinicID].organization].name} - {props.clinics[clinicID].name}
                                                    </option>
                                                ))}
                                            </select>
                                          ) : (
                                            (this.state.selectedClinicID === userClinic.id) ? (
                                                <Link to={`/clinics/${userClinic.id}`}>{props.clinics[userClinic.id].name}</Link>
                                            ) : (
                                                <Link to={`/users/${props.userID}/clinics/${userClinic.id}`}>{props.clinics[userClinic.id].name}</Link>
                                            )
                                          )}
                                        </td>
                                        <td>
                                        {userClinic.edit ? (
                                            <select className="form-control form-control-sm" value={userClinic.roleID} onChange={this.editRoleID(i)}>
                                                <option value="">Select role</option>
                                                {_.map(props.roles, role => (
                                                    <option key={role.id} value={role.id}>
                                                        {role.name}
                                                    </option>
                                                ))}
                                            </select>
                                        ) : ("")}
                                        </td>
                                        <td className="text-right">
                                          {userClinic.edit ? (
                                              <div className="btn-group" role="group">
                                                  <button className="btn btn-sm btn-light" disabled={userClinic.saving} type="button" onClick={this.cancelNewUserClinic(i)}>
                                                      <span className="icon_close" />
                                                  </button>
                                                  <button className="btn btn-sm btn-light" disabled={userClinic.saving || !userClinic.canSave} type="button" onClick={this.saveUserClinic(i)}>
                                                      <span className="icon_floppy" />
                                                  </button>
                                              </div>
                                          ) : (
                                              <div className="btn-group" role="group">
                                                  <button className="btn btn-sm btn-light" type="button" onClick={this.removeUserClinic(userClinic.id)}>
                                                      <span className="icon_trash" />
                                                  </button>
                                              </div>
                                          )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <button type="button" className="btn btn-sm btn-outline-primary col" disabled={(this.state.userClinics.length !== 0 && this.state.userClinics[this.state.userClinics.length - 1].edit) ? true : null} onClick={this.newUserClinic()}>
                            Add user to clinic
                        </button>
                    </div>
                    <div className="col">
                        <Route path="/users/:userID/clinics/:clinicID" component={ClinicDetail} />
                    </div>
                </div>
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

        return {
            userID: userID,
            clinics: state.clinics.allLoaded ? state.clinics.clinics : undefined,
            clinicsLoading: state.clinics.loading,
            organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
            organizationsLoading: state.organizations.loading,
            roles: state.roles.allLoaded ? state.roles.roles : undefined,
            rolesLoading: state.roles.loading,
            userRoles: state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[userID] ? state.userRoles.userUserRoles[userID] : undefined) : undefined,
            userRolesLoading: state.userRoles.loading,
            userClinicIDs: getUserClinicIDs(state, {userID: userID}),
            allowedClinicIDs: getUserAllowedClinicIDs(state, {userID: userID}),
            forbidden: state.userRoles.forbidden || state.users.forbidden || state.organizations.forbidden || state.clinics.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadRoles,
            loadOrganizations,
            loadClinics,
            deleteUserFromClinic,
            loadUserUserRoles,
            saveUserRoleCustomMsg,
            deleteUserRole
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(ClinicsList))
