import React from "react"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { push } from "react-router-redux"

import { RESOURCE_PATIENT_IDENTIFICATION, READ, WRITE } from "../../modules/validations"
import { search, cardToObject } from "../../modules/discovery"
import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"
import "./style.css"

const ListRow = ({ patient, canAddToWaitlist }) => {
    const p = cardToObject(patient)
    const id = p["syrian-id"] ? `Syrian ID: ${p["syrian-id"]}` : p["un-id"] ? `UN ID: ${p["un-id"]}` : ""

    return (
        <tr>
            <th scope="row">
                <Link className="patientLink" to={`/patients/${patient.patientID}`}>
                    <Patient data={p} />
                </Link>
            </th>
            <td>{p.nationality}</td>
            <td>{id}</td>
            <td>
                Camp {p.camp}, Tent {p.tent}
            </td>
            <td>{canAddToWaitlist && <Link to={`/to-waitlist/${patient.patientID}`}>Add to Waiting List</Link>}</td>
        </tr>
    )
}

class PatientList extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            searchQuery: ""
        }
        this.props.search("")
    }

    updateSearchQuery = e => {
        this.setState({ searchQuery: e.target.value })
    }

    search = e => {
        e.preventDefault()
        this.props.search(this.state.searchQuery)
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
                        <form onSubmit={this.search}>
                            <div className="input-group search">
                                <input
                                    name="search"
                                    placeholder="Search"
                                    className="form-control"
                                    value={this.state.searchQuery}
                                    onChange={this.updateSearchQuery}
                                />
                                <span className="input-group-append">
                                    <button type="submit" className="btn btn-secondary">
                                        Search
                                    </button>
                                </span>
                            </div>
                        </form>
                        {this.props.searching ? (
                            <Spinner />
                        ) : (
                            <table className="table patients">
                                <tbody>
                                    {this.props.patients.map(patient => (
                                        <ListRow patient={patient} key={patient.patientID} canAddToWaitlist={this.props.canAddToWaitlist} />
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
    { search, push }
)(PatientList)

export default PatientList

export { ListRow }
