import React from "react"

import Patient from "shared/containers/patient"
import "./style.css"

export default () => (
    <div className="patients">
        <header>
            <h1>Patients</h1>
            <button className="btn btn-secondary btn-wide" type="submit">
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
                            <button className="btn btn-link">Add to Waiting List</button>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    </div>
)
