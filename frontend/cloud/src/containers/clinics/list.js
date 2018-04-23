import React from "react"
import { Route, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadLocations } from "../../modules/locations"
import { loadOrganizations } from "../../modules/organizations"
import { saveClinic, loadClinics, deleteClinic } from "../../modules/clinics"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"

import UsersList from "./usersList"

class Clinics extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.clinics) {
            this.props.loadClinics()
        }
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.organizations) {
            this.props.loadOrganizations()
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.clinics && !nextProps.clinicsLoading) {
            this.props.loadClinics()
        }
        if (!nextProps.locations && !nextProps.locationsLoading) {
            this.props.loadLocations()
        }
        if (!nextProps.organizations && !nextProps.organizationsLoading) {
            this.props.loadOrganizations()
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.clinics || props.clinicsLoading || !props.locations || props.locationsLoading || !props.organizations || props.organizationsLoading || props.canEdit === undefined || props.canSee === undefined || props.validationsLoading

        let selectedClinicID = props.clinicID
        if (!selectedClinicID) {
            selectedClinicID = props.match.params.clinicID
        }
        this.setState({
            loading: loading,
            clinics: _.values(props.clinics),
            selectedClinicID: selectedClinicID ? selectedClinicID : undefined
        })
    }

    newClinic = () => e => {
        if (this.state.clinics) {
            let clinics = [...this.state.clinics, { id: "", edit: true, canSave: false, name: "", organization: "", location: "" }]
            this.setState({
                clinics: clinics,
                edit: true
            })
        }
    }

    editClinicName = index => e => {
        let clinics = [...this.state.clinics]
        clinics[index].name = e.target.value
        clinics[index].canSave = (clinics[index].name.length !== 0) && (clinics[index].organization.length !== 0) && (clinics[index].location.length !== 0)
        this.setState({ clinics: clinics })
    }

    editOrganizationID = index => e => {
        let clinics = [...this.state.clinics]
        clinics[index].organization = e.target.value
        clinics[index].canSave = (clinics[index].name.length !== 0) && (clinics[index].organization.length !== 0) && (clinics[index].location.length !== 0)
        this.setState({ clinics: clinics })
    }

    editLocationID = index => e => {
        let clinics = [...this.state.clinics]
        clinics[index].location = e.target.value
        clinics[index].canSave = (clinics[index].name.length !== 0) && (clinics[index].organization.length !== 0) && (clinics[index].location.length !== 0)
        this.setState({ clinics: clinics })
    }

    saveClinic = index => e => {
        let clinics = [...this.state.clinics]

        clinics[index].edit = false
        clinics[index].saving = true

        this.props.saveClinic(clinics[index])
            .then(response => {
                if (response.id) {
                    this.setState({ edit: false })
                    this.props.history.push(`/clinics/${response.id}`)
                }
            })
    }

    cancelNewClinic = index => e => {
        let clinics = [...this.state.clinics]
        clinics.splice(index, 1)
        this.setState({
            clinics: clinics,
            edit: false
        })
    }

    removeClinic = clinicID => e => {
        this.props.deleteClinic(clinicID)
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
            <div id="clinics">
                <div className="row">
                    <table className="table table-hover">
                        <thead>
                            <tr>
                                <th scope="col">#</th>
                                <th scope="col">Name</th>
                                <th scope="col">Organization</th>
                                <th scope="col">Location</th>
                                <th />
                            </tr>
                        </thead>
                        <tbody>
                            {_.map(this.state.clinics, (clinic, i) => (
                                <tr key={clinic.id || i} className={((this.state.edit && clinic.edit) || (!this.state.edit && this.state.selectedClinicID === clinic.id)) ? "table-active" : ""}>
                                    <th scope="row">{i+1}</th>
                                    <td>
                                    {(props.canEdit && clinic.edit) ? (
                                        <input
                                            value={clinic.name}
                                            onChange={this.editClinicName(i)}
                                            type="text"
                                            className="form-control form-control-sm"
                                            placeholder="Clinic name"
                                            aria-label="Clinic name"
                                        />
                                    ) : (
                                        <Link to={`/clinics/${clinic.id}`}>{clinic.name}</Link>
                                    )}
                                    </td>
                                    <td>
                                    {(props.canEdit && clinic.edit) ? (
                                        <select className="form-control form-control-sm" value={clinic.organization} onChange={this.editOrganizationID(i)}>
                                            <option value="">Select organization</option>
                                            {_.map(props.organizations, organization => (
                                                <option key={organization.id} value={organization.id}>
                                                    {organization.name}
                                                </option>
                                            ))}
                                        </select>
                                    ) : (
                                        <Link to={`/organizations/${clinic.organization}`}>{props.organizations[clinic.organization].name}</Link>
                                    )}
                                    </td>
                                    <td>
                                    {(props.canEdit && clinic.edit) ? (
                                        <select className="form-control form-control-sm" value={clinic.location} onChange={this.editLocationID(i)}>
                                            <option value="">Select location</option>
                                            {_.map(props.locations, location => (
                                                <option key={location.id} value={location.id}>
                                                    {location.name}
                                                </option>
                                            ))}
                                        </select>
                                    ) : (
                                        <Link to={`/locations/${clinic.location}`}>{props.locations[clinic.location].name}</Link>
                                    )}
                                    </td>
                                    <td className="text-right">
                                    {props.canEdit ? (
                                        clinic.edit ? (
                                            <div className="btn-group" role="group">
                                                <button className="btn btn-sm btn-light" disabled={clinic.saving} type="button" onClick={this.cancelNewClinic(i)}>
                                                    <span className="icon_close" />
                                                </button>
                                                <button className="btn btn-sm btn-light" disabled={clinic.saving || !clinic.canSave} type="button" onClick={this.saveClinic(i)}>
                                                    <span className="icon_floppy" />
                                                </button>
                                            </div>
                                        ) : (
                                            <div className="btn-group" role="group">
                                                <button onClick={this.removeClinic(clinic.id)} className="btn btn-sm btn-light" type="button">
                                                    <span className="icon_trash" />
                                                </button>
                                            </div>
                                        )
                                    ) : (null)}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    { props.canEdit ? (
                        <button type="button" className="btn btn-sm btn-outline-primary col" disabled={this.state.edit ? true : null} onClick={this.newClinic()}>
                            Add clinic
                        </button>
                    ) : (null)}
                </div>
                <div>
                    {this.state.edit ? (null) : (
                        <div className="m-4">
                            <Route exact path="/clinics/:clinicID" component={UsersList} />
                            <Route exact path="/clinics/:clinicID/users/:userID" component={UsersList} />
                        </div>
                    )}
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    clinics: ownProps.clinics ? (state.clinics.allLoaded ? _.fromPairs(_.map(ownProps.clinics, clinicID => [clinicID, state.clinics.clinics[clinicID]])) : undefined) : (state.clinics.allLoaded ? state.clinics.clinics : undefined),
    clinicsLoading: state.clinics.loading,
    organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
    organizationsLoading: state.organizations.loading,
    locations: state.locations.allLoaded ? state.locations.locations: undefined,
    locationsLoading: state.locations.loading,
    canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden: state.locations.forbidden || state.clinics.forbidden || state.organizations.forbidden,
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganizations,
            loadLocations,
            loadClinics,
            saveClinic,
            deleteClinic,
            loadUserRights,
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Clinics))
