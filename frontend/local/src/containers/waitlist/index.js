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
import PatientImage from "shared/containers/patient/image"
import PatientCard from "shared/containers/patient/card"
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
        const baseMatchURL = waitlistID && waitlistItemID ? "/waitlist/:waitlistID/:itemID" : "/waitlist/:waitlistID"

        return (
            <div className="waitlist">
                <h1>Waiting list</h1>

                <Section
                    type={TYPE_ENCOUNTER}
                    list={listEncounter}
                    title="Consultation"
                    waitlistID={match.params.waitlistID}
                    canAddExamination={this.props.canAddExamination}
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
                    title="Waiting List"
                    waitlistID={match.params.waitlistID}
                    canAddExamination={this.props.canAddExamination}
                    canRemoveFromWaitlist={this.props.canRemoveFromWaitlist}
                    canSeePatients={this.props.canSeePatients}
                    canSeeMainComplaint={this.props.canSeeMainComplaint}
                    canEditMainComplaint={this.props.canEditMainComplaint}
                    canSeeVitalSigns={this.props.canSeeVitalSigns}
                    canAddVitalSigns={this.props.canAddVitalSigns}
                />
                {this.props.canEditMainComplaint && <Route path={baseMatchURL + "/edit-complaint"} component={EditComplaint} />}
                {this.props.canAddVitalSigns && <Route exact path={baseMatchURL + "/add-data/:sign"} component={MedicalData} />}
                {this.props.canAddVitalSigns && <Route exact path={baseMatchURL + "/add-data"} component={MedicalData} />}
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
                    {(list || []).map(el => {
                        let patientObject =
                            el.patient &&
                            cardToObject({
                                patientID: el.patientID,
                                connections: el.patient
                            })
                        return (
                            <tr key={el.id}>
                                <th scope="row">
                                    <PatientImage data={patientObject} />
                                </th>
                                <td className="w-30">
                                    <Link className="patientLink" to={`/patients/${patientObject.patientID}`}>
                                        <PatientCard data={patientObject} withoutImage={true} />
                                    </Link>
                                </td>
                                <td>
                                    {canSeeMainComplaint && (
                                        <div>
                                            {el.mainComplaint.complaint}
                                            {el.priority === 1 && (
                                                <div>
                                                    <span className="badge badge-pill badge-danger">Urgent</span>
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </td>
                                {canSeeVitalSigns && <VitalSigns waitlistID={waitlistID} itemID={el.id} signs={el.vitalSigns || {}} />}
                                {(canEditMainComplaint || canAddVitalSigns || canRemoveFromWaitlist) && (
                                    <Tools
                                        listType={type}
                                        waitlistID={waitlistID}
                                        itemID={el.id}
                                        canAddExamination={canAddExamination}
                                        canEditMainComplaint={canEditMainComplaint}
                                        canSeeVitalSigns={canSeeVitalSigns}
                                        canAddVitalSigns={canAddVitalSigns}
                                        canRemoveFromWaitlist={canRemoveFromWaitlist}
                                    />
                                )}
                            </tr>
                        )
                    })}
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
        const { listType, waitlistID, itemID, canAddExamination, canEditMainComplaint, canAddVitalSigns, canRemoveFromWaitlist } = this.props
        return (
            <td className="tools">
                <UncontrolledDropdown direction="down">
                    <DropdownToggle color="link">
                        <span className="meatballs" />
                    </DropdownToggle>
                    <DropdownMenu right flip={false}>
                        {canAddExamination &&
                            listType === TYPE_ENCOUNTER && (
                                <Link to={`/waitlist/${waitlistID}/${itemID}/consultation`}>
                                    <DropdownItem>Start consultation</DropdownItem>
                                </Link>
                            )}
                        {canAddExamination &&
                            listType !== TYPE_ENCOUNTER && (
                                <Link onClick={this.moveToTop} to={`/waitlist/${waitlistID}/${itemID}/consultation`}>
                                    <DropdownItem onClick={this.moveToTop}>Start consultation out of order</DropdownItem>{" "}
                                </Link>
                            )}
                        {canEditMainComplaint && (
                            <Link to={`/waitlist/${waitlistID}/${itemID}/edit-complaint`}>
                                <DropdownItem>Edit main complaint</DropdownItem>
                            </Link>
                        )}
                        {canAddVitalSigns && (
                            <Link to={`/waitlist/${waitlistID}/${itemID}/add-data`}>
                                <DropdownItem>Add vital signs</DropdownItem>
                            </Link>
                        )}
                        {canRemoveFromWaitlist && (
                            <Link to={`/waitlist/${waitlistID}/${itemID}/remove`}>
                                <DropdownItem>Remove from Waiting list</DropdownItem>
                            </Link>
                        )}
                        <span className="arrow" />
                    </DropdownMenu>
                </UncontrolledDropdown>
            </td>
        )
    }
}

Tools = connect(state => ({}), {
    moveToTop
})(Tools)

const VitalSigns = ({ waitlistID, itemID, signs }) => (
    <td className="vital_signs">
        <div>
            <ul>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/height`}>
                    <li className={classnames({ active: signs.height })}>H</li>
                </Link>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/weight`}>
                    <li className={classnames({ active: signs.weight })}>W</li>
                </Link>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/temperature`}>
                    <li className={classnames({ active: signs.temperature })}>T</li>
                </Link>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/heart_rate`}>
                    <li className={classnames({ active: signs.heart_rate })}>HR</li>
                </Link>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/pressure`}>
                    <li className={classnames({ active: signs.pressure })}>BP</li>
                </Link>
                <Link to={`/waitlist/${waitlistID}/${itemID}/add-data/oxygen_saturation`}>
                    <li className={classnames({ active: signs.oxygen_saturation })}>OS</li>
                </Link>
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
