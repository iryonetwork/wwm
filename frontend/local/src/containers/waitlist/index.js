import React from "react"
import classnames from "classnames"
import { connect } from "react-redux"
import { Link } from "react-router-dom"
import { UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem } from "reactstrap"
import { listAll } from "../../modules/waitlist"
import { cardToObject } from "../../modules/discovery"
import { RESOURCE_WAITLIST, RESOURCE_PATIENT_IDENTIFICATION, RESOURCE_VITAL_SIGNS, READ, WRITE, UPDATE, DELETE } from "../../modules/validations"

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
        if (!this.props.canSeeWaitlist) {
            return null
        }

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

                <Section list={listEncounter} title="Encounter" waitlistID={match.params.waitlistID} canRemoveFromWaitlist={this.props.canRemoveFromWaitlist} canSeePatients={this.props.canSeePatients} canSeeMainComplaint={this.props.canSeeMainComplaint} canEditMainComplaint={this.props.canEditMainComplaint} canSeeVitalSigns={this.props.canSeeVitalSigns} canAddVitalSigns={this.props.canAddVitalSigns}/>

                <Section list={listNext} title="Up next" waitlistID={match.params.waitlistID} canRemoveFromWaitlist={this.props.canRemoveFromWaitlist} canSeePatients={this.props.canSeePatients} canSeeMainComplaint={this.props.canSeeMainComplaint} canEditMainComplaint={this.props.canEditMainComplaint} canSeeVitalSigns={this.props.canSeeVitalSigns} canAddVitalSigns={this.props.canAddVitalSigns}/>

                <Section list={list} title="Waiting list" waitlistID={match.params.waitlistID} canRemoveFromWaitlist={this.props.canRemoveFromWaitlist} canSeePatients={this.props.canSeePatients} canSeeMainComplaint={this.props.canSeeMainComplaint} canEditMainComplaint={this.props.canEditMainComplaint} canSeeVitalSigns={this.props.canSeeVitalSigns} canAddVitalSigns={this.props.canAddVitalSigns}/>
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

const Section = ({ list, title, waitlistID, canRemoveFromWaitlist, canSeePatients, canSeeMainComplaint, canEditMainComplaint, canSeeVitalSigns, canAddVitalSigns }) => {
    if (!list) {
        return null
    }

    return canSeePatients ? (
        <div className="part">
            <h2>{title}</h2>

            <table className="table patients">
                <tbody>
                    {(list || []).map(el => (
                        <tr key={el.id}>
                            <th scope="row">
                                <Patient data={el.patient && cardToObject({ connections: el.patient })} />
                            </th>
                            <td>
                                {canSeeMainComplaint && (
                                    <div>
                                        {el.complaint}
                                        {el.priority === 1 && (
                                            <div>
                                                <span className="badge badge-pill badge-danger">Urgent</span>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </td>
                            {canSeeVitalSigns && (<VitalSigns signs={el.vital_signs || {}} />)}
                            {(canEditMainComplaint || canAddVitalSigns || canRemoveFromWaitlist) && (
                                <Tools waitlistID={waitlistID} itemID={el.id} canEditMainComplaint={canEditMainComplaint} canSeeVitalSigns={canSeeVitalSigns} canAddVitalSigns={canAddVitalSigns} canRemoveFromWaitlist={canRemoveFromWaitlist} />
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    ) : (null)
}

const Tools = ({ waitlistID, itemID, canEditMainComplaint, canSeeVitalSigns, canAddVitalSigns, canRemoveFromWaitlist }) => (
    <td className="tools">
        <UncontrolledDropdown>
            <DropdownToggle color="link">
                <span className="meatballs" />
            </DropdownToggle>
            <DropdownMenu right>
                {canSeeVitalSigns && (
                    <DropdownItem>
                        <Link to={`/waitlist/${waitlistID}/${itemID}/consultation`}>Consultation</Link>
                    </DropdownItem>
                )}
                {canEditMainComplaint && (
                    <DropdownItem>
                        <Link to={`/waitlist/${waitlistID}/${itemID}/consultation/edit-complaint`}>Edit main complaint</Link>
                    </DropdownItem>
                )}
                {canAddVitalSigns && (
                    <DropdownItem>
                        <Link to={`/waitlist/${waitlistID}/${itemID}/consultation/add-data`}>Add vital signs</Link>
                    </DropdownItem>
                )}
                {canRemoveFromWaitlist && (
                    <DropdownItem>
                        <Link to={`/waitlist/${waitlistID}/${itemID}/consultation/remove`}>Remove from Waiting list</Link>
                    </DropdownItem>
                )}
            </DropdownMenu>
        </UncontrolledDropdown>
    </td>
)

const VitalSigns = ({ signs }) => (
    <td className="vital_signs">
        <div>
            <ul>
                <li className={classnames({ active: signs.height })}>H</li>
                <li className={classnames({ active: signs.weight })}>W</li>
                <li className={classnames({ active: signs.temperature })}>T</li>
                <li className={classnames({ active: signs.heart_rate })}>HR</li>
                <li className={classnames({ active: signs.pressure })}>BP</li>
                <li className={classnames({ active: signs.oxygen_saturation })}>OS</li>
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
        canSeeWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[READ],
        canRemoveFromWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[DELETE],
        canSeePatients: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
        canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE],
        canSeeMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canEditMainComplaint: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[UPDATE],
    }),
    {
        listAll,
    }
)(Waitlist)

export default Waitlist
