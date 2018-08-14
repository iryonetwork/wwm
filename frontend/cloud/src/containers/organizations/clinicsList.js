import React from "react"
import { NavLink, Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadOrganization } from "../../modules/organizations"
import { loadClinics, saveClinic, deleteClinic } from "../../modules/clinics"
import { loadLocations } from "../../modules/locations"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import Spinner from "shared/containers/spinner"

class ClinicsList extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.organization) {
            this.props.loadOrganization(this.props.organizationID)
        }
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.clinics) {
            this.props.loadClinics()
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.organization && !nextProps.organizationsLoading) {
            this.props.loadOrganization(nextProps.organizationID)
        }
        if (!nextProps.locations && !nextProps.locationsLoading) {
            this.props.loadLocations()
        }
        if (!nextProps.clinics && !nextProps.clinicsLoading) {
            this.props.loadClinics()
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.clinics ||
            props.clinicsLoading ||
            !props.organization ||
            props.organizationsLoading ||
            !props.locations ||
            props.locationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading
        this.setState({ loading: loading })

        if (!loading) {
            this.setState({
                organizationClinics: _.map(props.organization.clinics ? props.organization.clinics : [], clinicID => {
                    return props.clinics[clinicID]
                })
            })
        }
    }

    newClinic() {
        return e => {
            if (this.state.organizationClinics) {
                let organizationClinics = [
                    ...this.state.organizationClinics,
                    { id: "", edit: true, canSave: false, name: "", location: "", organization: this.props.organizationID }
                ]
                this.setState({
                    organizationClinics: organizationClinics,
                    edit: true
                })
            }
        }
    }

    editClinicName(index) {
        return e => {
            let organizationClinics = [...this.state.organizationClinics]
            organizationClinics[index].name = e.target.value
            organizationClinics[index].canSave = organizationClinics[index].name.length !== 0 && organizationClinics[index].location.length !== 0
            this.setState({ organizationClinics: organizationClinics })
        }
    }

    editLocationID(index) {
        return e => {
            let organizationClinics = [...this.state.organizationClinics]
            organizationClinics[index].location = e.target.value
            organizationClinics[index].canSave = organizationClinics[index].name.length !== 0 && organizationClinics[index].location.length !== 0
            this.setState({ organizationClinics: organizationClinics })
        }
    }

    saveClinic(index) {
        return e => {
            let organizationClinics = [...this.state.organizationClinics]

            organizationClinics[index].edit = false
            organizationClinics[index].saving = true

            this.props.saveClinic(organizationClinics[index])
            this.setState({ edit: false })
        }
    }

    cancelNewClinic(index) {
        return e => {
            let organizationClinics = [...this.state.organizationClinics]
            organizationClinics.splice(index, 1)
            this.setState({
                organizationClinics: organizationClinics,
                edit: false
            })
        }
    }

    removeClinic(clinicID) {
        return e => {
            this.props.deleteClinic(clinicID)
        }
    }

    render() {
        let props = this.props
        if (!props.organization || props.organizationsLoading) {
            return <Spinner />
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <div>
                <header>
                    <h1>Organizations</h1>
                </header>
                <h2>{props.organization.name}</h2>
                {props.organization ? (
                    <div className="navigation">
                        {props.canSeeOrganization ? (
                            <NavLink exact to={`/organizations/${props.organization.id}`}>
                                Organization's Data
                            </NavLink>
                        ) : null}
                        {props.canSeeUsers ? <NavLink to={`/organizations/${props.organization.id}/users`}>Users</NavLink> : null}
                        {props.canSee ? <NavLink to={`/organizations/${props.organization.id}/clinics`}>Clinics</NavLink> : null}
                    </div>
                ) : null}
                {this.state.loading ? (
                    <Spinner />
                ) : (
                    <div id="clinics">
                        <div className="row">
                            <div className="col-12">
                                {this.state.organizationClinics.length > 0 ? (
                                    <table className="table">
                                        <thead>
                                            <tr>
                                                <th className="w-7" scope="col">
                                                    #
                                                </th>
                                                <th className="w-30" scope="col">
                                                    Name
                                                </th>
                                                <th scope="col">Location</th>
                                                <th />
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {_.map(this.state.organizationClinics, (clinic, i) => (
                                                <tr key={clinic.id || i}>
                                                    <th className="w-7" scope="row">
                                                        {i + 1}
                                                    </th>
                                                    <td className="w-30">
                                                        {props.canEdit ? (
                                                            clinic.edit ? (
                                                                <input
                                                                    type="text"
                                                                    value={clinic.name || ""}
                                                                    onChange={this.editClinicName(i)}
                                                                    className="form-control"
                                                                    placeholder="Clinic name"
                                                                    aria-label="Clinic name"
                                                                />
                                                            ) : (
                                                                <Link to={`/clinics/${clinic.id}`}>{clinic.name}</Link>
                                                            )
                                                        ) : (
                                                            clinic.name
                                                        )}
                                                    </td>
                                                    <td>
                                                        {props.canEdit && clinic.edit ? (
                                                            <select className="form-control" value={clinic.location || ""} onChange={this.editLocationID(i)}>
                                                                <option value="">Select location</option>
                                                                {_.map(props.locations, location => (
                                                                    <option key={location.id} value={location.id || ""}>
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
                                                                <div>
                                                                    <button
                                                                        className="btn btn-secondary"
                                                                        disabled={clinic.saving}
                                                                        type="button"
                                                                        onClick={this.cancelNewClinic(i)}
                                                                    >
                                                                        Remove
                                                                    </button>
                                                                    <button
                                                                        className="btn btn-primary"
                                                                        disabled={clinic.saving || !clinic.canSave}
                                                                        type="button"
                                                                        onClick={this.saveClinic(i)}
                                                                    >
                                                                        Add
                                                                    </button>
                                                                </div>
                                                            ) : (
                                                                <div>
                                                                    <button className="btn btn-link" type="button" onClick={this.removeClinic(clinic.id)}>
                                                                        <span className="remove-link">Remove</span>
                                                                    </button>
                                                                </div>
                                                            )
                                                        ) : null}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                ) : (
                                    <h3>No clinics belong to the organization.</h3>
                                )}
                                {props.canEdit ? (
                                    <button
                                        type="button"
                                        className="btn btn-link"
                                        disabled={
                                            this.state.organizationClinics.length !== 0 &&
                                            this.state.organizationClinics[this.state.organizationClinics.length - 1].edit
                                                ? true
                                                : null
                                        }
                                        onClick={this.newClinic()}
                                    >
                                        Add Clinic
                                    </button>
                                ) : null}
                            </div>
                        </div>
                    </div>
                )}
            </div>
        )
    }
}

const makeMapStateToProps = () => {
    const mapStateToProps = (state, ownProps) => {
        let organizationID = ownProps.organizationID
        if (!organizationID) {
            organizationID = ownProps.match.params.organizationID
        }

        return {
            organizationID: organizationID,
            organization: state.organizations.organizations ? state.organizations.organizations[organizationID] : undefined,
            organizationsLoading: state.organizations.loading,
            clinics: state.clinics.allLoaded ? state.clinics.clinics : undefined,
            clinicsLoading: state.clinics.loading,
            locations: state.locations.allLoaded ? state.locations.locations : undefined,
            locationsLoading: state.locations.loading,
            canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeOrganization: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
            canSeeUsers: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
            validationsLoading: state.validations.loading,
            forbidden: state.organizations.forbidden || state.locations.forbidden || state.clinics.forbidden
        }
    }
    return mapStateToProps
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganization,
            loadClinics,
            saveClinic,
            deleteClinic,
            loadLocations,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(makeMapStateToProps, mapDispatchToProps)(ClinicsList))
