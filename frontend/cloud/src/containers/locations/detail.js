import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter, Link } from "react-router-dom"
import _ from "lodash"
import classnames from "classnames"

import { loadLocation, saveLocation } from "../../modules/locations"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { CATEGORY_COUNTRIES, loadCodes } from "../../modules/codes"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import { processStateOnChange, processStateOnBlur } from "../../utils/formFieldsUpdate"

import "../../styles/style.css"

class LocationDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            name: "",
            capacity: "",
            city: "",
            country: "",
            electricity: false,
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
        if (!nextProps.location && nextProps.locationID !== "new" && !this.props.locationsLoading) {
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
            props.locationsLoading ||
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
                electricity: props.location.electricity || false,
                waterSupply: props.location.waterSupply || false,
                manager: manager || {}
            })
        }
    }

    updateInput() {
        return e => {
            this.setState(processStateOnChange(this.state, e))
        }
    }

    onBlurInput() {
        return e => {
            this.setState(processStateOnBlur(this.state, e))
        }
    }

    updateCapacity() {
        return e => {
            var parsed = parseInt(e.target.value)
            if (!isNaN(parsed) && parsed >= 0) {
                this.setState({ capacity: e.target.value })
            }
        }
    }

    submit() {
        return e => {
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
            location.electricity = this.state.electricity
            location.waterSupply = this.state.waterSupply
            location.manager = _.clone(this.state.manager)

            if (!_.isEmpty(validationErrors)) {
                this.props.open("There are errors in the data submitted", "", COLOR_DANGER)
                this.setState({ validationErrors: validationErrors })
                return
            }

            this.props.saveLocation(location).then(response => {
                if (!location.id && response && response.id) {
                    this.props.history.push(`/locations`)
                }
            })
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
            <div>
                <header>
                    <h1>Locations</h1>
                    <Link to="/locations" className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={this.submit()} className="btn btn-primary btn-wide">
                        {props.locationsUpdating ? "Saving..." : "Save"}
                    </button>
                </header>
                <h2>{props.location ? props.location.name : "New Location"}</h2>
                <div className="location-form">
                    <form onSubmit={this.submit()} className="needs-validation" noValidate>
                        <div>
                            <div className="section">
                                <div className="form-row">
                                    <div className="form-group col-sm-3">
                                        <label>
                                            <input
                                                type="text"
                                                className={"form-control" + (this.state.validationErrors["name"] ? " is-invalid" : "")}
                                                id="name"
                                                value={this.state.name || ""}
                                                onChange={this.updateInput()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                                placeholder="Name"
                                                required="true"
                                            />
                                            <span>Name</span>
                                            {props.canEdit ? (
                                                this.state.validationErrors["name"] ? (
                                                    <div className="invalid-feedback">{this.state.validationErrors["name"]}</div>
                                                ) : (
                                                    <small className="form-text text-muted">Required</small>
                                                )
                                            ) : null}
                                        </label>
                                    </div>
                                    <div className="form-group col-sm-3">
                                        <label>
                                            <input
                                                type="text"
                                                className="form-control"
                                                id="capacity"
                                                value={this.state.capacity || ""}
                                                onChange={this.updateCapacity()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                                placeholder="Capacity"
                                            />
                                            <span>Capacity</span>
                                        </label>
                                    </div>
                                    <div className="form-group col-sm-3">
                                        <label>
                                            <select
                                                className={classnames("form-control", { selected: this.state.country })}
                                                id="country"
                                                value={this.state.country || ""}
                                                onChange={this.updateInput()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                            >
                                                <option value="">Select country</option>
                                                {_.map(props.countries, country => (
                                                    <option key={country.id} value={country.id}>
                                                        {country.title}
                                                    </option>
                                                ))}
                                            </select>
                                            <span>Country</span>
                                        </label>
                                    </div>
                                    <div className="form-group col-sm-3">
                                        <label>
                                            <input
                                                type="text"
                                                className="form-control"
                                                id="city"
                                                value={this.state.city || ""}
                                                onChange={this.updateInput()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                                placeholder="City"
                                            />
                                            <span>City</span>
                                        </label>
                                    </div>
                                </div>
                                <div className="form-row">
                                    <div className="form-inline-container">
                                        <span className="label">Electricity</span>
                                        <div className="form-check form-check-inline">
                                            <input
                                                className="form-check-input"
                                                type="radio"
                                                id="electricity"
                                                checked={this.state.electricity === true}
                                                onChange={this.updateInput()}
                                                disabled={!props.canEdit}
                                                value={true}
                                            />
                                            <label className="form-check-label">Yes</label>
                                        </div>
                                        <div className="form-check form-check-inline">
                                            <input
                                                className="form-check-input"
                                                type="radio"
                                                id="electricity"
                                                checked={this.state.electricity === false}
                                                onChange={this.updateInput()}
                                                disabled={!props.canEdit}
                                                value={false}
                                            />
                                            <label className="form-check-label">No</label>
                                        </div>
                                    </div>
                                </div>
                                <div className="form-row">
                                    <div className="form-inline-container">
                                        <span className="label">Water Supply</span>
                                        <div className="form-check form-check-inline">
                                            <input
                                                className="form-check-input"
                                                type="radio"
                                                id="waterSupply"
                                                checked={this.state.waterSupply === true}
                                                onChange={this.updateInput()}
                                                disabled={!props.canEdit}
                                                value={true}
                                            />
                                            <label className="form-check-label">Yes</label>
                                        </div>
                                        <div className="form-check form-check-inline">
                                            <input
                                                className="form-check-input"
                                                type="radio"
                                                id="waterSupply"
                                                checked={this.state.waterSupply === false}
                                                onChange={this.updateInput()}
                                                disabled={!props.canEdit}
                                                value={false}
                                            />
                                            <label className="form-check-label">No</label>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="section">
                                <h3>Manager</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="manager.name"
                                                    value={this.state.manager.name || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Full Name"
                                                />
                                                <span>Full Name</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="email"
                                                    className="form-control"
                                                    id="manager.email"
                                                    value={this.state.manager.email || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="user@email.com"
                                                />
                                                <span>Email Address</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="manager.phoneNumber"
                                                    value={this.state.manager.phoneNumber || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Phone Number"
                                                />
                                                <span>Phone Number</span>
                                            </label>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        {props.canEdit ? (
                            <div className="row buttons">
                                <div className="col-sm-4">
                                    <Link to="locations" className="btn btn-secondary btn-block">
                                        Cancel
                                    </Link>
                                </div>
                                <div className="col-sm-4">
                                    <button type="submit" className="btn btn-primary btn-block">
                                        {props.locationsUpdating ? "Saving..." : "Save"}
                                    </button>
                                </div>
                            </div>
                        ) : null}
                    </form>
                </div>
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
        locationsLoading: state.locations.loading,
        locationsUpdating: state.locations.updating,
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
