import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadLocations } from "../../modules/locations"
import { loadOrganizations } from "../../modules/organizations"
import { loadClinics, deleteClinic } from "../../modules/clinics"

class Clinics extends React.Component {
    componentDidMount() {
        this.props.loadClinics()
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.organizations) {
            this.props.loadOrganizations()
        }
    }

    removeClinic = clinicID => e => {
        this.props.deleteClinic(clinicID)
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
                        <th scope="col">Name</th>
                        <th scope="col">Organization</th>
                        <th scope="col">Location</th>
                        <th />
                    </tr>
                </thead>
                <tbody>
                    {_.map(_.filter(props.clinics, clinic => clinic), clinic => (
                        <tr key={clinic.id}>
                            <th scope="row">{++i}</th>
                            <td><Link to={`/clinics/${clinic.id}`}>{clinic.name}</Link></td>
                            <td>
                                <Link to={`/organizations/${clinic.organization}`}>{props.organizations[clinic.organization].name}</Link>
                            </td>
                            <td>
                                <Link to={`/locations/${clinic.location}`}>{props.locations[clinic.location].name}</Link>
                            </td>
                            <td className="text-right">
                                <button onClick={this.removeClinic(clinic.id)} className="btn btn-sm btn-light" type="button">
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

const mapStateToProps = (state, ownProps) => ({
    clinics:
        (ownProps.clinics ? (state.clinics.clinics ? _.fromPairs(_.map(ownProps.clinics, clinicID => [clinicID, state.clinics.clinics[clinicID]])) : {}) : state.clinics.clinics) ||
        {},
    organizations: state.organizations.organizations,
    locations: state.locations.locations,
    loading: state.locations.loading || state.organizations.loading || state.clinics.loading
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganizations,
            loadLocations,
            loadClinics,
            deleteClinic
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Clinics)
