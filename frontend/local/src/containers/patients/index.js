import React from "react"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { push } from "react-router-redux"
import classnames from "classnames"

import { read, DEFAULT_WAITLIST_ID } from "shared/modules/config"
import { RESOURCE_PATIENT_IDENTIFICATION, READ, WRITE } from "../../modules/validations"
import { search, cardToObject } from "../../modules/discovery"
import { PatientImage, PatientCard } from "shared/containers/patient"
import Spinner from "shared/containers/spinner"
import { CodeTitle } from "shared/containers/codes"

import { ReactComponent as SearchIcon } from "shared/icons/search.svg"
import { ReactComponent as SearchActiveIcon } from "shared/icons/search-active.svg"
import { ReactComponent as SpinnerIcon } from "shared/icons/spinner.svg"
import { ReactComponent as DeleteIcon } from "shared/icons/delete.svg"

import "./style.css"

const ListRow = ({ patient, canAddToWaitlist, waitlistID }) => {
    const p = cardToObject(patient)
    const id = p["syrian-id"] ? `Syrian ID: ${p["syrian-id"]}` : p["un-id"] ? `UN ID: ${p["un-id"]}` : ""

    return (
        <tr>
            <th scope="row">
                <PatientImage data={p} />
            </th>
            <td className="w-30">
                <Link className="patientLink" to={`/patients/${patient.patientID}`}>
                    <PatientCard data={p} withoutImage={true} />
                </Link>
            </td>
            <td className="w-15">
                <CodeTitle categoryId="countries" codeId={p.nationality} />
            </td>
            <td className="w-25">{id}</td>
            <td className="w-15">{p.region}</td>
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
                        ) : (
                            <table className="table patients">
                                <tbody>
                                    {this.props.patients.map(patient => (
                                        <ListRow
                                            patient={patient}
                                            key={patient.patientID}
                                            canAddToWaitlist={this.props.canAddToWaitlist}
                                            waitlistID={this.state.waitlistID}
                                        />
                                    ))}
                                </tbody>
                            </table>
                        )}
                    </div>
                ) : null}
            </div>
        )
    }
}

PatientList = connect(
    state => ({
        searching: state.discovery.searching || false,
        patients: state.discovery.patients || [],
        canAddPatient: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[WRITE],
        canSeePatients: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canAddToWaitlist: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[WRITE]
    }),
    { search, push, read }
)(PatientList)

export default PatientList

export { ListRow }
