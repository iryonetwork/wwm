import React from "react"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem } from "reactstrap"
import { listAll } from "../../modules/waitlist"

import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"

import "./style.css"

class Waitlist extends React.Component {
    constructor(props) {
        super(props)
        props.listAll(props.match.params.waitlistID)
    }

    render() {
        const { match, listEncounter, listNext, list, listing } = this.props

        if (listing) {
            return (
                <div className="waitlist">
                    <h1>Waiting list</h1>
                    <Spinner />
                </div>
            )
        }

        return (
            <div className="waitlist">
                <h1>Waiting list</h1>

                <Section list={listEncounter} title="Encounter" waitlistID={match.params.waitlistID} />

                <Section list={listNext} title="Up next" waitlistID={match.params.waitlistID} />

                <Section list={list} title="Waiting list" waitlistID={match.params.waitlistID} />
                {/* {listEncounter && <div className="part now">
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
                </div>}

                {listNext && <div className="part next">
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
                </div>}

                {list && <div className="part">
                    <h2>Waiting list</h2>

                    <table className="table patients">
                        <tbody>
                            {(list || []).map(el => (
                                <tr key={el.patient_id}>
                                    <th scope="row">
                                        <Patient />
                                    </th>
                                    <td>
                                        {el.complaint}
                                        {el.priority === 4 && <div>
                                            <span className="badge badge-pill badge-danger">Urgent</span>
                                        </div>}
                                    </td>
                                    <VitalSigns />
                                    <Tools waitlistID={match.params.waitlistID} itemID={"acsd"} />
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>} */}
            </div>
        )
    }
}

const Section = ({list, title, waitlistID}) => {
    if (!list) {
        return null
    }

    return (
        <div className="part">
            <h2>{title}</h2>

            <table className="table patients">
                <tbody>
                    {(list || []).map(el => (
                        <tr key={el.patient_id}>
                            <th scope="row">
                                <Patient />
                            </th>
                            <td>
                                {el.complaint}
                                {el.priority === 4 && <div>
                                    <span className="badge badge-pill badge-danger">Urgent</span>
                                </div>}
                            </td>
                            <VitalSigns />
                            <Tools waitlistID={waitlistID} itemID={el.id} />
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

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

Waitlist = connect(
    state => ({
        listEncounter: state.waitlist.list.length > 0 ? [state.waitlist.list[0]] : [],
        listNext: state.waitlist.list.length > 1 ? [state.waitlist.list[1]] : [],
        list: state.waitlist.list.length > 2 ? state.waitlist.list.slice(2) : [],
        listing: state.waitlist.listing,
    }),
    {
        listAll,
    }
)(Waitlist)

export default Waitlist
