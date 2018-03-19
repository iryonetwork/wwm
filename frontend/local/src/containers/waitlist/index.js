import React from "react"

import Patient from "shared/containers/patient"

import "./style.css"

export default () => (
    <div className="waitlist">
        <h1>Waiting list</h1>

        <div className="part now">
            <h2>Encounter</h2>

            <table className="table">
                <tbody>
                    <tr>
                        <th scope="row">
                            <Patient />
                        </th>
                        <td>Knee pain (both knees)</td>
                        <td>
                            <VitalSigns />
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div className="part next">
            <h2>Up Next</h2>

            <table className="table">
                <tbody>
                    <tr>
                        <th scope="row">
                            <Patient />
                        </th>
                        <td>Knee pain (both knees)</td>
                        <td>
                            <VitalSigns />
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div className="part">
            <h2>Waiting list</h2>

            <table className="table">
                <tbody>
                    {Array.from(Array(5), (v, i) => (
                        <tr key={i}>
                            <th scope="row">
                                <Patient />
                            </th>
                            <td>Knee pain (both knees)</td>
                            <td>
                                <VitalSigns />
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
)

const VitalSigns = () => (
    <div className="vital_signs">
        <ul>
            <li className="active">H</li>
            <li className="active">M</li>
            <li className="active">T</li>
            <li>HR</li>
            <li>BP</li>
            <li>OS</li>
        </ul>
    </div>
)
