import React from "react"

import PersonPlaceholder from "shared/public/person.svg"

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

        <table className="table">
            <tbody>
                {Array.from(Array(10), (v, i) => (
                    <tr key={i}>
                        <th scope="row">
                            <img src={PersonPlaceholder} alt="" />
                            <div>
                                <div className="name">Graves, Alma</div>
                                <div className="dob">
                                    3 Jun 1994
                                    <span className="age">24 y</span>
                                    F
                                </div>
                            </div>
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
