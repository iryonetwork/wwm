import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import _ from "lodash"

import { loadLocation, saveLocation } from "../../modules/locations"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { CATEGORY_COUNTRIES, loadCodes } from "../../modules/codes"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import { processStateOnChange, processStateOnBlur } from "../../utils/formFieldsUpdate"

class LocationDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            name: "",
            capacity: "",
            city: "",
            country: "",
            electricty: false,
            waterSupply: false,
            manager: {},
            loading: true,
            validationErrors: {}
        }
    }

    componentDidMount() {
        if (!this.props.location && this.props.locationID !== "new") {
            this.props.loadLocation(this.props.locationID)
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
        if (!nextProps.location && nextProps.locationID !== "new" && !this.props.locationLoading) {
            this.props.loadLocation(this.props.locationID)
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
            (!props.location && props.locationID !== "new") ||
            props.locationLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading ||
            !props.countries ||
            props.codesLoading
        this.setState({ loading: loading })

        if (props.location) {
            let manager = _.clone(props.location.manager)

            this.setState({
                name: props.location.name,
                capacity: props.location.capacity || "",
                country: props.location.country || "",
                city: props.location.city || "",
                electricty: props.location.electricty || false,
                waterSupply: props.location.waterSupply || false,
                manager: manager || {}
            })
        }
    }

    updateInput = e => {
        this.setState(processStateOnChange(this.state, e))
    }

    onBlurInput = e => {
        this.setState(processStateOnBlur(this.state, e))
    }

    updateCapacity = e => {
        var parsed = parseInt(e.target.value)
        if (!isNaN(parsed) && parsed >= 0) {
            this.setState({ capacity: e.target.value })
        }
    }

    submit = e => {
        e.preventDefault()
        this.props.close()

        let validationErrors = {}
        if (!this.state.name || this.state.name === "") {
            validationErrors["name"] = "Required"
        }

        let location = this.props.location ? this.props.location : {}

        location.name = this.state.name
        location.capacity = parseInt(this.state.capacity)
        location.country = this.state.country
        location.city = this.state.city
        location.electricty = this.state.electricty
        location.waterSupply = this.state.waterSupply
        location.manager = _.clone(this.state.manager)

        if (!_.isEmpty(validationErrors)) {
            this.props.open("There are errors in the data submitted", "", COLOR_DANGER)
            this.setState({ validationErrors: validationErrors })
            return
        }

        this.props.saveLocation(location).then(response => {
            if (!location.id && response && response.id) {
                this.props.history.push(`/locations/${response.id}`)
            }
        })
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
            <div>
                <h1>Locations</h1>
                <h2>{props.location ? this.props.location.name : "Add new location"}</h2>

                <form onSubmit={this.submit} className="needs-validation" noValidate>
                    <div className="form-group">
                        <label htmlFor="name">Name</label>
                        <input
                            type="text"
                            className={"form-control" + (this.state.validationErrors["name"] ? " is-invalid" : "")}
                            id="name"
                            value={this.state.name || ""}
                            onChange={this.updateInput}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                            placeholder="Location name"
                            required="true"
                        />
                        {props.canEdit ? (this.state.validationErrors["name"] ? (
                            <div className="invalid-feedback">{this.state.validationErrors["name"]}</div>
                        ) : (
                            <small className="form-text text-muted">Required</small>
                        )) : (null)}
                    </div>
                    <div className="form-group">
                        <label htmlFor="capacity">Capacity</label>
                        <input
                            type="text"
                            className="form-control"
                            id="capacity"
                            value={this.state.capacity || ""}
                            onChange={this.updateCapacity}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                            placeholder="e.g. 1000"
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="country">Country</label>
                        <select
                            className="form-control form-control-sm"
                            id="country"
                            value={this.state.country || ""}
                            onChange={this.updateInput}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                        >
                            <option value="">Select country</option>
                            {_.map(props.countries, country => (
                                <option key={country.id} value={country.id}>
                                    {country.title}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div className="form-group">
                        <label htmlFor="city">City</label>
                        <input
                            type="text"
                            className="form-control"
                            id="city"
                            value={this.state.city || ""}
                            onChange={this.updateInput}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                            placeholder="e.g. Beirut"
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="electricty">Electricity</label>
                        <input
                            type="checkbox"
                            className="form-control"
                            id="electricty"
                            checked={this.state.electricty}
                            onChange={this.updateInput}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="waterSupply">Water supply</label>
                        <input
                            type="checkbox"
                            className="form-control"
                            id="waterSupply"
                            checked={this.state.waterSupply}
                            onChange={this.updateInput}
                            onBlur={this.onBlurInput}
                            disabled={!props.canEdit}
                        />
                    </div>
                    <div className="form-group">
                        <h3>Manager</h3>
                        <div className="form-group">
                            <label htmlFor="firstName">Name</label>
                            <input
                                type="text"
                                className="form-control"
                                id="manager.name"
                                value={this.state.manager.name || ""}
                                onChange={this.updateInput}
                                onBlur={this.onBlurInput}
                                disabled={!props.canEdit}
                                placeholder="Full name"
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="email">Email address</label>
                            <input
                                type="email"
                                className="form-control"
                                id="manager.email"
                                value={this.state.manager.email || ""}
                                onChange={this.updateInput}
                                onBlur={this.onBlurInput}
                                disabled={!props.canEdit}
                                placeholder="user@email.com"
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="specialisation">Phone number</label>
                            <input
                                type="tel"
                                className="form-control"
                                id="manager.phoneNumber"
                                value={this.state.manager.phoneNumber || ""}
                                onChange={this.updateInput}
                                onBlur={this.onBlurInput}
                                disabled={!props.canEdit}
                                placeholder="+38640..."
                            />
                        </div>
                    </div>
                    <div className="form-group">
                        {props.canEdit ? (
                            <button type="submit" className="btn btn-outline-primary col">
                                Save
                            </button>
                        ) : null}
                    </div>
                </form>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.locationID
    if (!id) {
        id = ownProps.match.params.locationID
    }

    return {
        locationID: id,
        location: state.locations.locations ? state.locations.locations[id] : undefined,
        locationLoading: state.locations.loading,
        countries: state.codes.codes[CATEGORY_COUNTRIES],
        codesLoading: state.codes.loading,
        canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading,
        forbidden: state.locations.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadLocation,
            saveLocation,
            loadCodes,
            loadUserRights,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(LocationDetail))
