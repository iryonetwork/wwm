import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadLocations, deleteLocation } from "../../modules/locations"
import { CATEGORY_COUNTRIES, loadCodes } from "../../modules/codes"

class Locations extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.locations) {
            this.props.loadLocations()
        }
        if (!this.props.countries) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }


        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.locations && !nextProps.locationsLoading) {
            this.props.loadLocations()
        }
        if (!nextProps.countries && !nextProps.codesLoading) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading = !props.locations || props.locationsLoading
        this.setState({loading: loading})
    }

    removeLocation = locationID => e => {
        this.props.deleteLocation(locationID)
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
                            <td>{props.countries[location.country] ? props.countries[location.country].title : location.country}</td>
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
    locations: ownProps.locations ? (state.locations.allLoaded ? _.fromPairs(_.map(ownProps.locations, locationID => [locationID, state.locations.locations[locationID]])) : undefined) : (state.locations.allLoaded ? state.locations.locations : undefined),
    locationsLoading: state.locations.loading,
    countries: state.codes.codes[CATEGORY_COUNTRIES],
    codesLoading: state.codes.loading,
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadLocations,
            deleteLocation,
            loadCodes,
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Locations)
