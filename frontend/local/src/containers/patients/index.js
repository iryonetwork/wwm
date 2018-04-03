import React from "react"
import { Link, withRouter } from "react-router-dom"

import Patient from "shared/containers/patient"
import "./style.css"

export default withRouter(({ history }) => (
    <div className="patients">
        <header>
            <h1>Patients</h1>
            <button onClick={() => history.push("/patients/new")} className="btn btn-secondary btn-wide" type="submit">
                Add New Patient
            </button>
        </header>

        <input name="search" placeholder="Search" className="search" />

        <table className="table patients">
            <tbody>
                {Array.from(Array(10), (v, i) => (
                    <tr key={i}>
                        <th scope="row">
                            <Patient />
                        </th>
                        <td>Syrian</td>
                        <td>Syrian ID P349294839</td>
                        <td>Camp 15, Tent 06</td>
                        <td>
                            <Link to="/waitlist/add/abc">Add to Waiting List</Link>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    </div>
))
