import React from "react"
import classnames from "classnames"
import { connect } from "react-redux"
import { Route, Link } from "react-router-dom"
import { UncontrolledDropdown, DropdownToggle, DropdownMenu, DropdownItem } from "reactstrap"
import { listAll, moveToTop } from "../../modules/waitlist"
import { cardToObject } from "../../modules/discovery"
import {
    RESOURCE_WAITLIST,
    RESOURCE_EXAMINATION,
    RESOURCE_PATIENT_IDENTIFICATION,
    RESOURCE_VITAL_SIGNS,
    READ,
    WRITE,
    UPDATE,
    DELETE
} from "../../modules/validations"

import MedicalData from "./detail/add-data"
import EditComplaint from "./detail/edit-complaint"
import RemoveFromWaitlist from "./detail/remove"
import Patient from "shared/containers/patient"
import Spinner from "shared/containers/spinner"

import "./style.css"

const TYPE_ENCOUNTER = "encounter"
const TYPE_NEXT = "next"
const TYPE_WAITLIST = "waitlist"

class Waitlist extends React.Component {
    constructor(props) {
        super(props)

        const waitlistItemID = props.match.params.itemID

        if (!waitlistItemID || !props.listed) {
            props.listAll(props.match.params.waitlistID)
        }
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

        const waitlistID = match.params.waitlistID
        const waitlistItemID = match.params.itemID
        const baseMatchURL = (waitlistID && waitlistItemID) ? "/waitlist/:waitlistID/:itemID" : "/waitlist/:waitlistID"

        return (
            <div className="waitlist">
                <h1>Waiting list</h1>

                <Section
                    type={TYPE_ENCOUNTER}
                    list={listEncounter}
                    title="Encounter"
                    waitlistID={match.params.waitlistID}
                    canAddExamination={this.props.canAddExamination}
                    canSeeExamination={this.props.canSeeExamination}
                    canRemoveFromWaitlist={this.props.canRemoveFromWaitlist}
                    canSeePatients={this.props.canSeePatients}
                    canSeeMainComplaint={this.props.canSeeMainComplaint}
                    canEditMainComplaint={this.props.canEditMainComplaint}
                    canSeeVitalSigns={this.props.canSeeVitalSigns}
                    canAddVitalSigns={this.props.canAddVitalSigns}
                />

                <Section
                    type={TYPE_NEXT}
                    list={listNext}
                    title="Up next"
                    waitlistID={match.params.waitlistID}
                    canSeeExamination={this.props.canSeeExamination}
                    canAddExamination={this.props.canAddExamination}
                    canRemoveFromWaitlist={this.props.canRemoveFromWaitlist}
                    canSeePatients={this.props.canSeePatients}
                    canSeeMainComplaint={this.props.canSeeMainComplaint}
                    canEditMainComplaint={this.props.canEditMainComplaint}
                    canSeeVitalSigns={this.props.canSeeVitalSigns}
                    canAddVitalSigns={this.props.canAddVitalSigns}
                />

                <Section
                    type={TYPE_WAITLIST}
                    list={list}
                    title="Waiting list"
                    waitlistID={match.params.waitlistID}
                    canSeeExamination={this.props.canSeeExamination}
                    canAddExamination={this.props.canAddExamination}
                    canRemoveFromWaitlist={this.props.canRemoveFromWaitlist}
                    canSeePatients={this.props.canSeePatients}
                    canSeeMainComplaint={this.props.canSeeMainComplaint}
                    canEditMainComplaint={this.props.canEditMainComplaint}
                    canSeeVitalSigns={this.props.canSeeVitalSigns}
                    canAddVitalSigns={this.props.canAddVitalSigns}
                />
                {this.props.canEditMainComplaint && <Route path={baseMatchURL + "/edit-complaint"} component={EditComplaint} />}
                {this.props.canAddVitalSigns && <Route path={baseMatchURL + "/add-data"} component={MedicalData} />}
                {this.props.canRemoveFromWaitlist && <Route path={baseMatchURL + "/remove"} component={RemoveFromWaitlist} />}
            </div>
        )
    }
}

const Section = ({
    type,
    list,
    title,
    waitlistID,
    canSeeExamination,
    canAddExamination,
    canRemoveFromWaitlist,
    canSeePatients,
    canSeeMainComplaint,
    canEditMainComplaint,
    canSeeVitalSigns,
    canAddVitalSigns
}) => {
    if (!list) {
        return null
    }

    return canSeePatients ? (
        <div className="part">
            <h2>{title}</h2>

            {(list || []).length === 0 && <p>No one on the list</p>}

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
                            {canSeeVitalSigns && <VitalSigns signs={el.vitalSigns || {}} />}
                            {(canEditMainComplaint || canAddVitalSigns || canRemoveFromWaitlist) && (
                                <Tools
                                    listType={type}
                                    waitlistID={waitlistID}
                                    itemID={el.id}
                                    canSeeExamination={canSeeExamination}
                                    canAddExamination={canAddExamination}
                                    canEditMainComplaint={canEditMainComplaint}
                                    canSeeVitalSigns={canSeeVitalSigns}
                                    canAddVitalSigns={canAddVitalSigns}
                                    canRemoveFromWaitlist={canRemoveFromWaitlist}
                                />
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    ) : null
}

class Tools extends React.Component {
    constructor(props) {
        super(props)
        this.moveToTop = this.moveToTop.bind(this)
    }

    moveToTop(e) {
        this.props.moveToTop(this.props.waitlistID, this.props.itemID)
    }

    render() {
        const { listType, waitlistID, itemID, canSeeExamination, canAddExamination, canEditMainComplaint, canAddVitalSigns, canRemoveFromWaitlist } = this.props
        return (
            <td className="tools">
                <UncontrolledDropdown>
                    <DropdownToggle color="link">
                        <span className="meatballs" />
                    </DropdownToggle>
                    <DropdownMenu right>
                        {canAddExamination &&
                            listType === TYPE_ENCOUNTER && (
                                <DropdownItem>
                                    <Link to={`/waitlist/${waitlistID}/${itemID}/consultation`}>Start consultation</Link>
                                </DropdownItem>
                            )}
                        {canAddExamination &&
                            listType !== TYPE_ENCOUNTER && (
                                <DropdownItem onClick={this.moveToTop}>
                                    <Link onClick={this.moveToTop} to={`/waitlist/${waitlistID}/${itemID}/consultation`}>
                                        Start consultation out of order
                                    </Link>
                                </DropdownItem>
                            )}
                        {canEditMainComplaint && (
                            <DropdownItem>
                                <Link to={`/waitlist/${waitlistID}/${itemID}/edit-complaint`}>Edit main complaint</Link>
                            </DropdownItem>
                        )}
                        {canAddVitalSigns && (
                            <DropdownItem>
                                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data`}>Add vital signs</Link>
                            </DropdownItem>
                        )}
                        {canSeeExamination && (
                            <DropdownItem>
                                <Link to={`/waitlist/${waitlistID}/${itemID}/consultation`}>See consultation data</Link>
                            </DropdownItem>
                        )}
                        {canRemoveFromWaitlist && (
                            <DropdownItem>
                                <Link to={`/waitlist/${waitlistID}/${itemID}/remove`}>Remove from Waiting list</Link>
                            </DropdownItem>
                        )}
                    </DropdownMenu>
                </UncontrolledDropdown>
            </td>
        )
    }
}

Tools = connect(state => ({}), {
    moveToTop
})(Tools)

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
        listed: state.waitlist.listed,
        canSeeExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[READ],
        canAddExamination: ((state.validations.userRights || {})[RESOURCE_EXAMINATION] || {})[WRITE],
        canSeeWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[READ],
        canRemoveFromWaitlist: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[DELETE],
        canSeePatients: ((state.validations.userRights || {})[RESOURCE_PATIENT_IDENTIFICATION] || {})[READ],
        canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
        canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE],
        canSeeMainComplaint: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[READ],
        canEditMainComplaint: ((state.validations.userRights || {})[RESOURCE_WAITLIST] || {})[UPDATE]
    }),
    {
        listAll
    }
)(Waitlist)

export default Waitlist
