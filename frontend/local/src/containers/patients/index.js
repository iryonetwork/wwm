import React from "react"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { push } from "react-router-redux"
import classnames from "classnames"
import _ from "lodash"

import { read } from "shared/modules/config"
import { DEFAULT_WAITLIST_ID } from  "../../modules/config"
import { RESOURCE_PATIENT_IDENTIFICATION, READ, WRITE } from "../../modules/validations"
import { search, cardToObject } from "../../modules/discovery"
import PatientImage from "shared/containers/patient/image"
import PatientCard from "shared/containers/patient/card"
import Spinner from "shared/containers/spinner"
import CodeTitle from "shared/containers/codes/title"

import { ReactComponent as SearchIcon } from "shared/icons/search.svg"
import { ReactComponent as SearchActiveIcon } from "shared/icons/search-active.svg"
import { ReactComponent as SpinnerIcon } from "shared/icons/spinner.svg"
import { ReactComponent as DeleteIcon } from "shared/icons/delete.svg"

import "./style.css"

const ListRow = ({ patient, canAddToWaitlist, waitlistID }) => {
    const id = patient["syrian-id"] ? `Syrian ID: ${patient["syrian-id"]}` : patient["un-id"] ? `UN ID: ${patient["un-id"]}` : ""

    return (
        <tr>
            <th scope="row">
                <PatientImage data={patient} />
            </th>
            <td className="w-30">
                <Link className="patientLink" to={`/patients/${patient.patientID}`}>
                    <PatientCard data={patient} withoutImage={true} />
                </Link>
            </td>
            <td className="w-15">
                <CodeTitle categoryId="countries" codeId={patient.nationality} />
            </td>
            <td className="w-25">{id}</td>
            <td className="w-15">{patient.region}</td>
            <td>
                {canAddToWaitlist && (
                    <Link className="btn btn-link" to={`/to-waitlist/${waitlistID}/${patient.patientID}`}>
                        Add to Waiting List
                    </Link>
                )}
            </td>
        </tr>
    )
}

class PatientList extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            waitlistID: this.props.read(DEFAULT_WAITLIST_ID),
            searching: this.props.searching,
            searchQuery: ""
        }
        this.props.search("")
    }

    componentDidUpdate(prevProps) {
        if (!this.props.searching && prevProps.searching && this.state.searching) {
            // trigger delayed searching state change to false to prevent spinner flickering too much
            window.setTimeout(() => {
                this.setState({ searching: false })
            }, 250)
        } else if (this.props.searching !== this.state.searching) {
            this.setState({ searching: this.props.searching })
        }
    }

    updateSearchQuery = e => {
        this.state.searchTimeout && clearTimeout(this.state.searchTimeout)

        let timeout = 800
        if (e.target.value === "") {
            timeout = 0
        }

        this.setState({
            searchQuery: e.target.value,
            searchTimeout: window.setTimeout(this.search(e.target.value), timeout)
        })
    }

    searchBoxKeyPress = e => {
        if (e.key === "Enter") {
            this.state.searchTimeout && clearTimeout(this.state.searchTimeout)
            this.setState({
                searchTimeout: window.setTimeout(this.search(this.state.searchQuery), 0)
            })
        }
    }

    clearSearchQuery = e => {
        this.state.searchTimeout && window.clearTimeout(this.state.searchTimeout)
        this.setState({
            searchQuery: ""
        })
        this.props.search("")
    }

    search = value => {
        return () => this.props.search(value)
    }

    render() {
        const { push } = this.props

        return (
            <div className="patients">
                <header>
                    <h1>Patients</h1>
                    {this.props.canAddPatient ? (
                        <button onClick={() => push("/new-patient")} className="btn btn-secondary btn-wide" type="submit">
                            Add New Patient
                        </button>
                    ) : null}
                </header>

                {this.props.canSeePatients ? (
                    <div>
                        <div className="input-group search">
                            <span className="input-group-prepend">
                                {this.state.searching ? <SpinnerIcon /> : this.state.searchFocused ? <SearchActiveIcon /> : <SearchIcon />}
                            </span>
                            <input
                                name="search"
                                placeholder="Search"
                                className={classnames("form-control", {
                                    searching: this.state.searching
                                })}
                                value={this.state.searchQuery}
                                onChange={this.updateSearchQuery}
                                onKeyPress={this.searchBoxKeyPress}
                                onFocus={() => this.setState({ searchFocused: true })}
                                onBlur={() => this.setState({ searchFocused: false })}
                            />
                            {this.state.searchQuery.length !== 0 && (
                                <span className="input-group-append">
                                    <button className="btn" onClick={this.clearSearchQuery}>
                                        <DeleteIcon />
                                    </button>
                                </span>
                            )}
                        </div>
                        {this.state.searching ? (
                            <Spinner />
                        ) : this.props.numberOfPatients === 0 ? (
                            <div className="noResults">
                                <h3>No patients found</h3>
                                {this.state.searchQuery && (
                                    <button onClick={this.clearSearchQuery} className="btn btn-primary btn-wide">
                                        Clear Search
                                    </button>
                                )}
                            </div>
                        ) : this.state.searchQuery ? (
                            <div className="section">
                                <h3>
                                    Showing {this.props.numberOfPatients} result{this.props.numberOfPatients > 1 && "s"}
                                </h3>
                                <table className="table patients">
                                    <tbody>
                                        {_.map(this.props.patientsByInitial, (patients, initial) =>
                                            patients.map(patient => (
                                                <ListRow
                                                    patient={patient}
                                                    key={patient.patientID}
                                                    canAddToWaitlist={this.props.canAddToWaitlist}
                                                    waitlistID={this.state.waitlistID}
                                                />
                                            ))
                                        )}
                                    </tbody>
                                </table>
                            </div>
                        ) : (
                            _.map(this.props.patientsByInitial, (patients, initial) => (
                                <div key={`initial-${initial}`} className="section">
                                    <h3>{initial}</h3>
                                    <table className="table patients">
                                        <tbody>
                                            {patients.map(patient => (
                                                <ListRow
                                                    patient={patient}
                                                    key={patient.patientID}
                                                    canAddToWaitlist={this.props.canAddToWaitlist}
                                                    waitlistID={this.state.waitlistID}
                                                />
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            ))
                        )}
                    </div>
                ) : null}
            </div>
        )
    }
}

PatientList = connect(
    state => ({
        searching: state.discovery.searching || state.codes.isFetching || state.codes.loading || false,
        numberOfPatients: state.discovery.patients ? state.discovery.patients.length : 0,
        patientsByInitial: state.discovery.patients ? sortPatientsByLastNameInitial(state.discovery.patients) : {},
        canAddPatient: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[WRITE],
        canSeePatients: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canAddToWaitlist: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[WRITE]
    }),
    { search, push, read }
)(PatientList)

const sortPatientsByLastNameInitial = patients => {
    return _.groupBy(
        _.orderBy(
            _.flatMap(patients, cardToObject),
            [
                p => {
                    let name = "" + p.lastName + p.firstName
                    return name.toLowerCase()
                }
            ],
            ["asc"]
        ),
        p => (p.lastName ? p.lastName.charAt(0) : " ")
    )
}

export default PatientList

export { ListRow }
