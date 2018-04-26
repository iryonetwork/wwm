import React from "react"
import { Link } from "react-router-dom"
import { UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem } from "reactstrap"

import Patient from "shared/containers/patient"

import "./style.css"

export default ({ match }) => (
    <div className="waitlist">
        <h1>Waiting list</h1>

        <div className="part now">
            <h2>Encounter</h2>

            <table className="table patients">
                <tbody>
                    <tr>
                        <th scope="row">
                            <Patient />
                        </th>
                        <td>Knee pain (both knees)</td>
                        <VitalSigns />
                        <Tools waitlistID={match.params.waitlistID} itemID={"acsd"} />
                    </tr>
                </tbody>
            </table>
        </div>

        <div className="part next">
            <h2>Up Next</h2>

            <table className="table patients">
                <tbody>
                    <tr>
                        <th scope="row">
                            <Patient />
                        </th>
                        <td>
                            Knee pain (both knees)
                            <div>
                                <span className="badge badge-pill badge-danger">Urgent</span>
                            </div>
                        </td>
                        <VitalSigns />
                        <Tools waitlistID={match.params.waitlistID} itemID={"acsd"} />
                    </tr>
                </tbody>
            </table>
        </div>

        <div className="part">
            <h2>Waiting list</h2>

            <table className="table patients">
                <tbody>
                    {Array.from(Array(5), (v, i) => (
                        <tr key={i}>
                            <th scope="row">
                                <Patient />
                            </th>
                            <td>
                                Knee pain (both knees)
                                <div>
                                    <span className="badge badge-pill badge-danger">Urgent</span>
                                </div>
                            </td>
                            <VitalSigns />
                            <Tools waitlistID={match.params.waitlistID} itemID={"acsd"} />
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    </div>
)

const Tools = ({ waitlistID, itemID }) => (
    <td className="tools">
        <UncontrolledDropdown>
            <DropdownToggle color="link">
                <span className="meatballs" />
            </DropdownToggle>
            <DropdownMenu right>
                <DropdownItem>Edit main complaint</DropdownItem>
                <DropdownItem>
                    <Link to={`/waitlist/${waitlistID}/${itemID}/consultation/add-data`}>Add vital signs</Link>
                </DropdownItem>
                <DropdownItem>Remove from Waiting list</DropdownItem>
            </DropdownMenu>
        </UncontrolledDropdown>
    </td>
)

const VitalSigns = () => (
    <td className="vital_signs">
        <div>
            <ul>
                <li className="active">H</li>
                <li className="active">M</li>
                <li className="active">T</li>
                <li>HR</li>
                <li>BP</li>
                <li>OS</li>
            </ul>
        </div>
    </td>
)
