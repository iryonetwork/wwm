import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadLocations, deleteLocation } from "../../modules/locations"

class Locations extends React.Component {
    componentDidMount() {
        this.props.loadLocations()
    }

    removeLocation = locationID => e => {
        this.props.deleteLocation(locationID)
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
                        <th scope="col">City</th>
                        <th scope="col">Country</th>
                        <th scope="col">Clinics</th>
                        <th />
                    </tr>
                </thead>
                <tbody>
                    {_.map(_.filter(props.locations, location => location), location => (
                        <tr key={location.id}>
                            <th scope="row">{++i}</th>
                            <td><Link to={`/locations/${location.id}`}>{location.name}</Link></td>
                            <td>{location.city}</td>
                            <td>{location.country}</td>
                            <td>{location.clinics.length}</td>
                            <td className="text-right">
                                <button onClick={this.removeLocation(location.id)} className="btn btn-sm btn-light" type="button">
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
    locations:
        (ownProps.locations ? (state.locations.locations ? _.fromPairs(_.map(ownProps.locations, locationID => [locationID, state.locations.locations[locationID]])) : {}) : state.locations.locations) ||
        {},
    loading: state.locations.loading
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadLocations,
            deleteLocation
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Locations)
