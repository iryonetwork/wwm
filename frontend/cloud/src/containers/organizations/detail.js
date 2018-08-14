import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter, Route, Link, NavLink, Switch } from "react-router-dom"
import _ from "lodash"
import classnames from "classnames"

import { joinPaths } from "shared/utils"
import { loadOrganization, saveOrganization } from "../../modules/organizations"
import { CATEGORY_COUNTRIES, loadCodes } from "../../modules/codes"
import { ADMIN_RIGHTS_RESOURCE, SELF_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { open, close, COLOR_DANGER } from "shared/modules/alert"
import ClinicsList from "./clinicsList"
import UsersList from "./usersList"
import { processStateOnChange, processStateOnBlur } from "../../utils/formFieldsUpdate"

import "../../styles/style.css"

class OrganizationDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            organization: {},
            name: "",
            legalStatus: "",
            serviceType: "",
            address: {},
            representative: {},
            primaryContact: {},
            loading: true,
            validationErrors: {}
        }
    }

    componentDidMount() {
        if (!this.props.organization && this.props.organizationID !== "new") {
            this.props.loadOrganization(this.props.organizationID)
        }
        if (!this.props.countries) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }
        if (this.props.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.organization && nextProps.organizationID !== "new" && !this.props.organizationsLoading) {
            this.props.loadOrganization(nextProps.organizationID)
        }
        if (!nextProps.countries && !nextProps.codesLoading) {
            this.props.loadCodes(CATEGORY_COUNTRIES)
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }
        if (nextProps.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            (!props.organization && props.organizationID !== "new") ||
            props.organizationsLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading ||
            !props.countries ||
            props.codesLoading
        this.setState({ loading: loading })

        if (props.organization) {
            this.setState({
                organization: props.organization,
                name: props.organization.name,
                legalStatus: props.organization.legalStatus || "",
                serviceType: props.organization.serviceType || "",
                address: _.clone(props.organization.address) || {},
                representative: _.clone(props.organization.representative) || {},
                primaryContact: _.clone(props.organization.primaryContact) || {}
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

    submit() {
        return e => {
            e.preventDefault()
            this.props.close()

            let validationErrors = {}
            if (!this.state.name || this.state.name === "") {
                validationErrors["name"] = "Required"
            }

            let organization = this.state.organization

            organization.name = this.state.name
            organization.legalStatus = this.state.legalStatus
            organization.serviceType = this.state.serviceType
            organization.address = _.clone(this.state.address)
            organization.representative = _.clone(this.state.representative)
            organization.primaryContact = _.clone(this.state.primaryContact)

            if (!_.isEmpty(validationErrors)) {
                this.props.open("There are errors in the data submitted", "", COLOR_DANGER)
                this.setState({ validationErrors: validationErrors })
                return
            }

            this.props.saveOrganization(organization).then(response => {
                if (!organization.id && response && response.id) {
                    this.props.history.push(`/organizations/${response.id}`)
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
                    <h1>Organizations</h1>
                    <Link to="/organizations" className="btn btn-secondary btn-wide">
                        Cancel
                    </Link>
                    <button onClick={this.submit()} className="btn btn-primary btn-wide">
                        {props.organizationsUpdating ? "Saving..." : "Save"}
                    </button>
                </header>
                <h2>{props.organization ? props.organization.name : "New Organization"}</h2>
                {props.organization ? (
                    <div className="navigation">
                        {props.canSee ? (
                            <NavLink exact to={`/organizations/${props.organization.id}`}>
                                Organization's Data
                            </NavLink>
                        ) : null}
                        {props.canSeeUsers ? <NavLink to={`/organizations/${props.organization.id}/users`}>Users</NavLink> : null}
                        {props.canSeeClinics ? <NavLink to={`/organizations/${props.organization.id}/clinics`}>Clinics</NavLink> : null}
                    </div>
                ) : null}
                <div className="organization-form">
                    <form onSubmit={this.submit()} className="needs-validation" noValidate>
                        <div>
                            <div className="section">
                                <div className="form-row">
                                    <div className="form-group col-sm-4">
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
                                    <div className="form-group col-sm-4">
                                        <label>
                                            <input
                                                type="text"
                                                className="form-control"
                                                id="legalStatus"
                                                value={this.state.legalStatus || ""}
                                                onChange={this.updateInput()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                                placeholder="Legal Status"
                                            />
                                            <span>Legal Status</span>
                                        </label>
                                    </div>
                                    <div className="form-group col-sm-4">
                                        <label>
                                            <input
                                                type="text"
                                                className="form-control"
                                                id="serviceType"
                                                value={this.state.serviceType || ""}
                                                onChange={this.updateInput()}
                                                onBlur={this.onBlurInput()}
                                                disabled={!props.canEdit}
                                                placeholder="Service Type"
                                            />
                                            <span>Service Type</span>
                                        </label>
                                    </div>
                                </div>
                            </div>
                            <div className="section">
                                <h3>Address</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="address.addressLine1"
                                                    value={this.state.address.addressLine1 || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Address Line 1"
                                                />
                                                <span>Address Line 1</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="address.addressLine2"
                                                    value={this.state.address.addressLine2 || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Address Line 2"
                                                />
                                                <span>Address Line 1</span>
                                            </label>
                                        </div>
                                    </div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="address.city"
                                                    value={this.state.address.city}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="City"
                                                />
                                                <span>City</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="address.postCode"
                                                    value={this.state.address.postCode || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Post Code"
                                                />
                                                <span>Post Code</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <select
                                                    className={classnames("form-control", { selected: this.state.address.country })}
                                                    id="address.country"
                                                    value={this.state.address.country || ""}
                                                    onChange={this.updateInput()}
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
                                    </div>
                                </div>
                            </div>
                            <div className="section">
                                <h3>Representative</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="representative.name"
                                                    value={this.state.representative.name || ""}
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
                                                    id="representative.email"
                                                    value={this.state.representative.email || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Email Address"
                                                />
                                                <span>Email Address</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="representative.phoneNumber"
                                                    value={this.state.representative.phoneNumber || ""}
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
                            <div className="section">
                                <h3>Primary Contact</h3>
                                <div>
                                    <div className="form-row">
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="text"
                                                    className="form-control"
                                                    id="primaryContact.name"
                                                    value={this.state.primaryContact.name || ""}
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
                                                    id="primaryContact.email"
                                                    value={this.state.primaryContact.email || ""}
                                                    onChange={this.updateInput()}
                                                    onBlur={this.onBlurInput()}
                                                    disabled={!props.canEdit}
                                                    placeholder="Email Address"
                                                />
                                                <span>Email Address</span>
                                            </label>
                                        </div>
                                        <div className="form-group col-sm-4">
                                            <label>
                                                <input
                                                    type="tel"
                                                    className="form-control"
                                                    id="primaryContact.phoneNumber"
                                                    value={this.state.primaryContact.phoneNumber || ""}
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
                                    <Link to="/organizations" className="btn btn-secondary btn-block">
                                        Cancel
                                    </Link>
                                </div>
                                <div className="col-sm-4">
                                    <button type="submit" className="btn btn-primary btn-block">
                                        {props.organizationsUpdating ? "Saving..." : "Save"}
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

const mapOrganizationDetailStateToProps = (state, ownProps) => {
    let id = ownProps.organizationID
    if (!id) {
        id = ownProps.match.params.organizationID
    }

    return {
        organizationID: id,
        organization: state.organizations.organizations ? state.organizations.organizations[id] : undefined,
        organizationsLoading: state.organizations.loading,
        organizationsUpdating: state.organizations.updating,
        countries: state.codes.codes[CATEGORY_COUNTRIES],
        codesLoading: state.codes.loading,
        canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        canSee: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeUsers: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading,
        forbidden: state.organizations.forbidden
    }
}

const mapOrganizationDetailDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadOrganization,
            saveOrganization,
            loadCodes,
            loadUserRights,
            open,
            close
        },
        dispatch
    )

OrganizationDetail = withRouter(connect(mapOrganizationDetailStateToProps, mapOrganizationDetailDispatchToProps)(OrganizationDetail))

class OrganizationRoutes extends React.Component {
    componentDidMount() {
        if (this.props.canSeeOrganization === undefined || this.props.canSeeUsers === undefined || this.props.canSeeClinics === undefined) {
            this.props.loadUserRights()
        }
    }

    componentWillReceiveProps(nextProps) {
        if (
            (nextProps.canSeeOrganization === undefined || nextProps.canSeeUsers === undefined || nextProps.canSeeClinics === undefined) &&
            !nextProps.validationsLoading
        ) {
            this.props.loadUserRights()
        }
    }

    render() {
        let { match, organizationID, canSeeOrganization, canSeeUsers, canSeeClinics } = this.props

        return (
            <Switch>
                {canSeeOrganization && <Route exact path={match.url} component={() => <OrganizationDetail organizationID={organizationID} />} />}
                {canSeeUsers && <Route exact path={joinPaths(match.url, "users")} component={() => <UsersList organizationID={organizationID} />} />}
                {canSeeUsers && <Route exact path={joinPaths(match.url, "users", ":userID")} component={() => <UsersList organizationID={organizationID} />} />}
                {canSeeClinics && <Route exact path={joinPaths(match.url, "clinics")} component={() => <ClinicsList organizationID={organizationID} />} />}
                {canSeeClinics && (
                    <Route exact path={joinPaths(match.url, "clinics", ":clinicID")} component={() => <ClinicsList organizationID={organizationID} />} />
                )}
            </Switch>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let organizationID = ownProps.organizationID
    if (!organizationID) {
        organizationID = ownProps.match.params.organizationID
    }

    return {
        organizationID: organizationID,
        canSeeOrganization: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        canSeeUsers: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        canSeeClinics: state.validations.userRights ? state.validations.userRights[SELF_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUserRights
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(OrganizationRoutes))
