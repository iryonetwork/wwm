import React from "react"
import { Link, withRouter } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import classnames from "classnames"
import { push } from "react-router-redux"

import { loadLocations } from "../../modules/locations"
import { loadOrganizations } from "../../modules/organizations"
import { saveClinic, loadClinics, deleteClinic } from "../../modules/clinics"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { confirmationDialog } from "shared/utils"

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
        let loading =
            !props.clinics ||
            props.clinicsLoading ||
            !props.locations ||
            props.locationsLoading ||
            !props.organizations ||
            props.organizationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading

        let selectedClinicID = props.clinicID
        if (!selectedClinicID) {
            selectedClinicID = props.match.params.clinicID
        }
        this.setState({
            loading: loading,
            clinics: _.values(props.clinics),
            selectedClinicID: selectedClinicID || undefined
        })
    }

    newClinic() {
        return e => {
            if (this.state.clinics) {
                let clinics = [...this.state.clinics, { id: "", edit: true, canSave: false, name: "", organization: "", location: "" }]
                this.setState({
                    clinics: clinics,
                    edit: true
                })
            }
        }
    }

    editClinicName(index) {
        return e => {
            let clinics = [...this.state.clinics]
            clinics[index].name = e.target.value
            clinics[index].canSave = clinics[index].name.length !== 0 && clinics[index].organization.length !== 0 && clinics[index].location.length !== 0
            this.setState({ clinics: clinics })
        }
    }

    editOrganizationID(index) {
        return e => {
            let clinics = [...this.state.clinics]
            clinics[index].organization = e.target.value
            clinics[index].canSave = clinics[index].name.length !== 0 && clinics[index].organization.length !== 0 && clinics[index].location.length !== 0
            this.setState({ clinics: clinics })
        }
    }

    editLocationID(index) {
        return e => {
            let clinics = [...this.state.clinics]
            clinics[index].location = e.target.value
            clinics[index].canSave = clinics[index].name.length !== 0 && clinics[index].organization.length !== 0 && clinics[index].location.length !== 0
            this.setState({ clinics: clinics })
        }
    }

    saveClinic(index) {
        return e => {
            let clinics = [...this.state.clinics]

            clinics[index].edit = false
            clinics[index].saving = true

            this.props.saveClinic(clinics[index]).then(response => {
                if (response && response.id) {
                    this.setState({ edit: false })
                    this.props.history.push(`/clinics/${response.id}`)
                }
            })
        }
    }

    cancelNewClinic(index) {
        return e => {
            let clinics = [...this.state.clinics]
            clinics.splice(index, 1)
            this.setState({
                clinics: clinics,
                edit: false
            })
        }
    }

    removeClinic(index) {
        return e => {
            confirmationDialog(`Click OK to confirm that you want to remove clinic ${this.state.clinics[index].name}.`, () =>
                this.props.deleteClinic(this.state.clinics[index].id)
            )
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
            <div id="clinics">
                <div className="row">
                    <table className="table">
                        <thead>
                            <tr>
                                <th className="w-7" scope="col">
                                    #
                                </th>
                                <th className="w-20" scope="col">
                                    Name
                                </th>
                                <th className="w-20" scope="col">
                                    Organization
                                </th>
                                <th className="w-20" scope="col">
                                    Location
                                </th>
                                <th />
                            </tr>
                        </thead>
                        <tbody>
                            {_.map(this.state.clinics, (clinic, i) => (
                                <React.Fragment key={clinic.id || i}>
                                    <tr
                                        className={classnames({
                                            "table-active": this.state.selectedClinicID === clinic.id,
                                            "table-edit": props.canEdit && clinic.edit
                                        })}
                                    >
                                        <th className="w-7" scope="row">
                                            {i + 1}
                                        </th>
                                        <td className="w-20">
                                            {props.canEdit && clinic.edit ? (
                                                <input
                                                    value={clinic.name || ""}
                                                    onChange={this.editClinicName(i)}
                                                    type="text"
                                                    className="form-control"
                                                    placeholder="Clinic Name"
                                                    aria-label="Clinic Name"
                                                />
                                            ) : (
                                                clinic.name
                                            )}
                                        </td>
                                        <td className="w-20">
                                            {props.canEdit && clinic.edit ? (
                                                <select className="form-control" value={clinic.organization || ""} onChange={this.editOrganizationID(i)}>
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
                                        <td className="w-20">
                                            {props.canEdit && clinic.edit ? (
                                                <select className="form-control" value={clinic.location || ""} onChange={this.editLocationID(i)}>
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
                                                        {this.state.selectedClinicID === clinic.id ? (
                                                            <button className="btn btn-link" type="button" onClick={() => this.props.push(`/clinics`)}>
                                                                Hide Users
                                                                <span className="arrow-up-icon" />
                                                            </button>
                                                        ) : (
                                                            <button
                                                                className="btn btn-link"
                                                                type="button"
                                                                onClick={() => this.props.push(`/clinics/${clinic.id}`)}
                                                            >
                                                                Show Users
                                                                <span className="arrow-down-icon" />
                                                            </button>
                                                        )}
                                                        <button className="btn btn-link" type="button" onClick={this.removeClinic(i)}>
                                                            <span className="remove-link">Remove</span>
                                                        </button>
                                                    </div>
                                                )
                                            ) : null}
                                        </td>
                                    </tr>
                                    {this.state.selectedClinicID === clinic.id ? (
                                        <tr className="table-active">
                                            <td colSpan="5" className="row-details-container">
                                                <UsersList clinicID={clinic.id} />
                                            </td>
                                        </tr>
                                    ) : null}
                                </React.Fragment>
                            ))}
                        </tbody>
                    </table>
                    {props.canEdit ? (
                        <button type="button" className="btn btn-link" disabled={this.state.edit ? true : null} onClick={this.newClinic()}>
                            Add Clinic
                        </button>
                    ) : null}
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    clinics: ownProps.clinics
        ? state.clinics.allLoaded
            ? _.fromPairs(_.map(ownProps.clinics, clinicID => [clinicID, state.clinics.clinics[clinicID]]))
            : undefined
        : state.clinics.allLoaded
            ? state.clinics.clinics
            : undefined,
    clinicsLoading: state.clinics.loading,
    organizations: state.organizations.allLoaded ? state.organizations.organizations : undefined,
    organizationsLoading: state.organizations.loading,
    locations: state.locations.allLoaded ? state.locations.locations : undefined,
    locationsLoading: state.locations.loading,
    canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden: state.locations.forbidden || state.clinics.forbidden || state.organizations.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            push,
            loadOrganizations,
            loadLocations,
            loadClinics,
            saveClinic,
            deleteClinic,
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Clinics))
