import React from "react"
import { Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadLocations, deleteLocation } from "../../modules/locations"
import { CATEGORY_COUNTRIES, loadCodes } from "../../modules/codes"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"

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
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
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
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.locations ||
            props.locationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading ||
            !props.countries ||
            props.codesLoading
        this.setState({ loading: loading })
    }

    removeLocation(locationID) {
        return e => {
            this.props.deleteLocation(locationID)
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

        let i = 0
        return (
            <table className="table">
                <thead>
                    <tr>
                        <th className="w-7" scope="col">
                            #
                        </th>
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
                            <th className="w-10" scope="row">
                                {++i}
                            </th>
                            <td>
                                <Link to={`/locations/${location.id}`}>{location.name}</Link>
                            </td>
                            <td>{location.city}</td>
                            <td>{props.countries[location.country] ? props.countries[location.country].title : location.country}</td>
                            <td>{location.clinics.length}</td>
                            <td className="text-right">
                                {props.canEdit ? (
                                    <button onClick={this.removeLocation(location.id)} className="btn btn-link" type="button">
                                        <span className="remove-link">Remove</span>
                                    </button>
                                ) : null}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        )
    }
}

const mapStateToProps = (state, ownProps) => ({
    locations: ownProps.locations
        ? state.locations.allLoaded
            ? _.fromPairs(_.map(ownProps.locations, locationID => [locationID, state.locations.locations[locationID]]))
            : undefined
        : state.locations.allLoaded
            ? state.locations.locations
            : undefined,
    locationsLoading: state.locations.loading,
    countries: state.codes.codes[CATEGORY_COUNTRIES],
    codesLoading: state.codes.loading,
    canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
    validationsLoading: state.validations.loading,
    forbidden: state.locations.forbidden
})

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadLocations,
            deleteLocation,
            loadCodes,
            loadUserRights
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Locations)
