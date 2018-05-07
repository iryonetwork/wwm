import React from "react"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { push } from "react-router-redux"

import { search, cardToObject } from "../../modules/discovery"
import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"
import "./style.css"

const ListRow = ({ patient }) => {
    const p = cardToObject(patient)
    const id = p["syrian-id"] ? `Syrian ID: ${p["syrian-id"]}` : p["un-id"] ? `UN ID: ${p["un-id"]}` : ""

    return (
        <tr>
            <th scope="row">
                <Patient data={p} />
            </th>
            <td>{p.nationality}</td>
            <td>{id}</td>
            <td>
                Camp {p.camp}, Tent {p.tent}
            </td>
            <td>
                <Link to={`/to-waitlist/${patient.patientID}`}>Add to Waiting List</Link>
            </td>
        </tr>
    )
}

class PatientList extends React.Component {
    constructor(props) {
        super(props)
        this.props.search("")
    }

    render() {
        const { push } = this.props

        return (
            <div className="patients">
                <header>
                    <h1>Patients</h1>
                    <button onClick={() => push("/new-patient")} className="btn btn-secondary btn-wide" type="submit">
                        Add New Patient
                    </button>
                </header>

                <input name="search" placeholder="Search" className="search" />

                {this.props.searching ? (
                    <Spinner />
                ) : (
                    <table className="table patients">
                        <tbody>{this.props.patients.map(patient => <ListRow patient={patient} key={patient.patientID} />)}</tbody>
                    </table>
                )}
            </div>
        )
    }
}

PatientList = connect(
    state => ({
        searching: state.discovery.searching || false,
        patients: state.discovery.patients || []
    }),
    { search, push }
)(PatientList)

export default PatientList
