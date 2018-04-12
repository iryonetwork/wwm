import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"
import moment from "moment"

import { loadLocation, saveLocation } from "../../modules/locations"
import { open, close, COLOR_DANGER } from "shared/modules/alert"

class LocationDetail extends React.Component {
    constructor(props) {
        super(props)
        if (props.location) {
            this.state = {
                name: props.location.name,
                capacity: props.location.capacity ? props.location.capacity : "",
                country: props.location.country ? props.location.country : "",
                city: props.location.city ? props.location.city : "",
                electricty: props.location.electricty ? props.location.electricty : false,
                waterSupply: props.location.watterSupply ? props.location.watterSupply : false,
            }
        } else {
            this.state = {
                name: "",
                capacity: "",
                city: "",
                country: "",
                electricty: false,
                waterSupply: false,
            }
        }
    }

    componentDidMount() {
        if (!this.props.location && this.props.locationID !== "new") {
            this.props.loadLocation(this.props.locationID)
        }
    }

    componentWillReceiveProps(props) {
        if (props.location) {
            this.setState({ name: props.location.name })
            this.setState({ capacity: props.location.capacity ? props.location.capacity : "" })
            this.setState({ country: props.location.country ? props.location.country : "" })
            this.setState({ city: props.location.city ? props.location.city : "" })
            this.setState({ electricty: props.location.electricty ? props.location.electricty : false })
            this.setState({ waterSupply: props.location.watterSupply ? props.location.watterSupply : false })
        }
    }

    updateInput = e => {
        const target = e.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const id = target.id;

        this.setState({
          [id]: value
        });
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

        let location = this.props.location

        location.name = this.state.name
        location.capacity = parseInt(this.state.capacity)
        location.country = this.state.country
        location.city = this.state.city
        location.electricty = this.state.electricty
        location.waterSupply = this.state.waterSupply

        this.props.saveLocation(location)
        this.forceUpdate()
    }

    render() {
        let props = this.props
        console.log(props)
        if (!props.location && props.locationID !== "new") {
            return <div>Loading...</div>
        }
        return (
            <div>
                <h1>Locations</h1>
                <h2>{props.location ? this.props.location.name : "Add new location"}</h2>

                <form onSubmit={this.submit}>
                    <div className="form-group">
                        <label htmlFor="name">Name</label>
                        <input className="form-control" id="name" value={this.state.name} onChange={this.updateInput} placeholder="Location name" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="capacity">Capacity</label>
                        <input className="form-control" id="capacity" value={this.state.capacity} onChange={this.updateCapacity} placeholder="e.g. 1000" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="country">Country</label>
                        <input className="form-control" id="country" value={this.state.country} onChange={this.updateInput} placeholder="e.g. Lebanon" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="city">City</label>
                         <input className="form-control" id="city" value={this.state.city} onChange={this.updateInput} placeholder="e.g. Beirut" />
                    </div>
                    <div className="form-group">
                        <label htmlFor="electricty">Electricity</label>
                        <input type="checkbox" className="form-control" id="electricty" checked={this.state.electricty} onChange={this.updateInput} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="waterSupply">Water supply</label>
                        <input type="checkbox" className="form-control" id="waterSupply" checked={this.state.waterSupply} onChange={this.updateInput} />
                    </div>
                    <button type="submit" className="btn btn-sm btn-outline-secondary">
                        Save
                    </button>
                </form>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    let id = ownProps.locationID
    if (!id) {
        id = ownProps.match.params.id
    }

    return {
        location: state.locations.locations ? state.locations.locations[id] : undefined,
        loading: state.locations.loading,
        locationID: id,
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadLocation,
            saveLocation,
            open,
            close
        },
        dispatch
    )

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(LocationDetail))
