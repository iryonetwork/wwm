import React from "react"
import _ from "lodash"
import { connect } from "react-redux"
import { Collapse } from "reactstrap"
import { Route, Link } from "react-router-dom"

import { joinPaths } from "shared/utils"
import Spinner from "shared/containers/spinner"
import { VitalSignCard, HeightVitalSignCard, WeightVitalSignCard, TemperatureVitalSignCard, BloodPressureVitalSignCard } from "shared/containers/vitalSign"
import { RESOURCE_VITAL_SIGNS, READ, WRITE } from "../../../modules/validations"
import { fetchHealthRecords } from "../../../modules/patient"
import AddMedicalData from "../../waitlist/detail/add-data"

class MedicalData extends React.Component {
    constructor(props) {
        super(props)
        if ((props.patientID && !props.patientRecords) || props.patientID !== props.loadedPatientID) {
            props.fetchHealthRecords(props.patientID)
        }

        this.toggleMeasurementsHistory = this.toggleMeasurementsHistory.bind(this)
        this.toggleVitalSignsHistory = this.toggleVitalSignsHistory.bind(this)

        this.state = {
            bodyMeasurementsHistory: false,
            vitalSignsHistory: false
        }
    }

    componentWillReceiveProps(nextProps) {
        if ((nextProps.medicalData === undefined || nextProps.patientID !== nextProps.loadedPatientID) && !nextProps.medicalDataLoading) {
            this.props.fetchHealthRecords(nextProps.patientID)
        }
    }

    toggleMeasurementsHistory = sign => () => {
        this.setState({ bodyMeasurementsHistory: this.state.bodyMeasurementsHistory !== sign ? sign : false })
    }

    toggleVitalSignsHistory = sign => () => {
        this.setState({ vitalSignsHistory: this.state.vitalSignsHistory !== sign ? sign : false })
    }

    render() {
        let { match, medicalData, medicalDataLoading, canSeeVitalSigns, canAddVitalSigns, inConsultation } = this.props

        if (medicalDataLoading) {
            return <Spinner />
        }

        return canSeeVitalSigns ? (
            <div className="medicalData">
                <header>
                    <h1>Medical Data</h1>
                    {inConsultation &&
                        canAddVitalSigns && (
                            <React.Fragment>
                                <Route exact path={match.path + "/add-data"} component={AddMedicalData} />
                                <Link to={joinPaths(match.url, "add-data")} className="btn btn-secondary btn-wide">
                                    Add Medical Data
                                </Link>
                            </React.Fragment>
                        )}
                </header>
                <div>
                    {medicalData && !_.isEmpty(medicalData) ? (
                        <React.Fragment>
                            {medicalData.height || medicalData.weight || medicalData.bmi ? (
                                <React.Fragment>
                                    <h2>Body Measurements</h2>
                                    <div className="section" key="bodyMeasurements">
                                        <div className="card-group">
                                            {medicalData.height && (
                                                <div className="cardContainer">
                                                    <HeightVitalSignCard
                                                        id="height0"
                                                        name="Height"
                                                        precision={0}
                                                        value={medicalData.height[0].value}
                                                        unit="cm"
                                                        timestamp={medicalData.height[0].timestamp}
                                                        timestampWarning={medicalData.height[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleMeasurementsHistory("height")}
                                                        isActive={this.state["bodyMeasurementsHistory"] === "height"}
                                                    />
                                                </div>
                                            )}

                                            {medicalData.weight && (
                                                <div className="cardContainer">
                                                    <WeightVitalSignCard
                                                        id="weight0"
                                                        name="Body mass"
                                                        precision={1}
                                                        value={medicalData.weight[0].value}
                                                        unit="kg"
                                                        timestamp={medicalData.weight[0].timestamp}
                                                        timestampWarning={medicalData.weight[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleMeasurementsHistory("weight")}
                                                        isActive={this.state["bodyMeasurementsHistory"] === "weight"}
                                                    />
                                                </div>
                                            )}

                                            {medicalData.bmi && (
                                                <div className="cardContainer">
                                                    <VitalSignCard
                                                        id="bmi0"
                                                        name="BMI"
                                                        value={medicalData.bmi[0].value}
                                                        unit=""
                                                        precision={2}
                                                        timestamp={medicalData.bmi[0].timestamp}
                                                        timestampWarning={medicalData.bmi[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleMeasurementsHistory("bmi")}
                                                        isActive={this.state["bodyMeasurementsHistory"] === "bmi"}
                                                    />
                                                </div>
                                            )}
                                        </div>
                                        <Collapse isOpen={this.state["bodyMeasurementsHistory"] !== false}>
                                            <div className="part" key="bodyMeasurementsHistory">
                                                <h3>History</h3>
                                            </div>
                                        </Collapse>
                                        <Collapse isOpen={this.state["bodyMeasurementsHistory"] === "height"}>
                                            {medicalData.height && medicalData.height.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.height.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"height" + i} className="cardContainer">
                                                                    <HeightVitalSignCard
                                                                        id={"height" + i}
                                                                        name="Height"
                                                                        precision={0}
                                                                        value={reading.value}
                                                                        unit="cm"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for height.</div>
                                            )}
                                        </Collapse>
                                        <Collapse isOpen={this.state["bodyMeasurementsHistory"] === "weight"}>
                                            {medicalData.weight && medicalData.weight.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.weight.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"weight" + i} className="cardContainer">
                                                                    <WeightVitalSignCard
                                                                        id={"weight" + i}
                                                                        name="Body mass"
                                                                        precision={1}
                                                                        value={reading.value}
                                                                        unit="kg"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for body mass.</div>
                                            )}
                                        </Collapse>
                                        <Collapse isOpen={this.state["bodyMeasurementsHistory"] === "bmi"}>
                                            {medicalData.bmi && medicalData.bmi.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.bmi.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"bmi" + i} className="cardContainer">
                                                                    <VitalSignCard
                                                                        id={"bmi" + i}
                                                                        name="BMI"
                                                                        value={reading.value}
                                                                        unit=""
                                                                        precision={2}
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for body mass index.</div>
                                            )}
                                        </Collapse>
                                    </div>
                                </React.Fragment>
                            ) : null}

                            {medicalData.temperature || medicalData.heart_rate || medicalData.pressure || medicalData.oxygen_saturation ? (
                                <React.Fragment>
                                    <h2>Vital Signs</h2>
                                    <div className="section" key="vitalSigns">
                                        <div className="card-group">
                                            {medicalData.temperature && (
                                                <div className="cardContainer">
                                                    <TemperatureVitalSignCard
                                                        id="temperature0"
                                                        name="Temperature"
                                                        precision={1}
                                                        value={medicalData.temperature[0].value}
                                                        unit="°C"
                                                        timestamp={medicalData.temperature[0].timestamp}
                                                        timestampWarning={medicalData.temperature[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleVitalSignsHistory("temperature")}
                                                        isActive={this.state["vitalSignsHistory"] === "temperature"}
                                                    />
                                                </div>
                                            )}
                                            {medicalData.heart_rate && (
                                                <div className="cardContainer">
                                                    <VitalSignCard
                                                        id="heart_rate0"
                                                        name="Heart rate"
                                                        value={medicalData.heart_rate[0].value}
                                                        unit="bpm"
                                                        timestamp={medicalData.heart_rate[0].timestamp}
                                                        timestampWarning={medicalData.heart_rate[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleVitalSignsHistory("heart_rate")}
                                                        isActive={this.state["vitalSignsHistory"] === "heart_rate"}
                                                    />
                                                </div>
                                            )}
                                            {medicalData.pressure &&
                                                medicalData.pressure[0].value &&
                                                medicalData.pressure[0].value.systolic &&
                                                medicalData.pressure[0].value.diastolic && (
                                                    <div className="cardContainer">
                                                        <BloodPressureVitalSignCard
                                                            id="pressure0"
                                                            name="Blood pressure"
                                                            precision={this.props.bloodPressureUnit === "mm[Hg]" ? 0 : 1}
                                                            value={[medicalData.pressure[0].value.systolic, medicalData.pressure[0].value.diastolic]}
                                                            unit="mm[Hg]"
                                                            timestamp={medicalData.pressure[0].timestamp}
                                                            timestampWarning={medicalData.pressure[0].timestampWarning}
                                                            consultationTooltipOn={inConsultation}
                                                            onClick={() => this.toggleVitalSignsHistory("pressure")}
                                                            isActive={this.state["vitalSignsHistory"] === "pressure"}
                                                        />
                                                    </div>
                                                )}
                                            {medicalData.oxygen_saturation && (
                                                <div className="cardContainer">
                                                    <VitalSignCard
                                                        id="oxygen_saturation0"
                                                        name="Oxygen saturation"
                                                        value={medicalData.oxygen_saturation[0].value}
                                                        unit="%"
                                                        timestamp={medicalData.oxygen_saturation[0].timestamp}
                                                        timestampWarning={medicalData.oxygen_saturation[0].timestampWarning}
                                                        consultationTooltipOn={inConsultation}
                                                        onClick={() => this.toggleVitalSignsHistory("oxygen_saturation")}
                                                        isActive={this.state["vitalSignsHistory"] === "oxygen_saturation"}
                                                    />
                                                </div>
                                            )}
                                        </div>
                                        <Collapse isOpen={this.state["vitalSignsHistory"] !== false}>
                                            <div className="part" key="vitalSignsHistory">
                                                <h3>History</h3>
                                            </div>
                                        </Collapse>
                                        <Collapse isOpen={this.state["vitalSignsHistory"] === "temperature"}>
                                            {medicalData.temperature && medicalData.temperature.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.temperature.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"temperature" + i} className="cardContainer">
                                                                    <TemperatureVitalSignCard
                                                                        id={"temperature" + i}
                                                                        name="Temperature"
                                                                        precision={1}
                                                                        value={reading.value}
                                                                        unit="°C"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for temperature.</div>
                                            )}
                                        </Collapse>
                                        <Collapse isOpen={this.state["vitalSignsHistory"] === "heart_rate"}>
                                            {medicalData.heart_rate && medicalData.heart_rate.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.heart_rate.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"heart_rate" + i} className="cardContainer">
                                                                    <VitalSignCard
                                                                        id={"heart_rate" + i}
                                                                        name="Heart rate"
                                                                        value={reading.value}
                                                                        unit="bpm"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for heart rate.</div>
                                            )}
                                        </Collapse>
                                        <Collapse isOpen={this.state["vitalSignsHistory"] === "pressure"}>
                                            {medicalData.pressure && medicalData.pressure.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.pressure.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"reading" + i} className="cardContainer">
                                                                    <BloodPressureVitalSignCard
                                                                        id={"pressure" + i}
                                                                        name="Blood pressure"
                                                                        precision={this.props.bloodPressureUnit === "mm[Hg]" ? 0 : 1}
                                                                        value={[reading.value.systolic, reading.value.diastolic]}
                                                                        unit="mm[Hg]"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for blood pressure.</div>
                                            )}
                                        </Collapse>
                                        <Collapse isOpen={this.state["vitalSignsHistory"] === "oxygen_saturation"}>
                                            {medicalData.oxygen_saturation && medicalData.oxygen_saturation.length > 1 ? (
                                                <div className="card-group">
                                                    {medicalData.oxygen_saturation.map((reading, i) => {
                                                        return (
                                                            i !== 0 && (
                                                                <div key={"oxygen_saturation" + i} className="cardContainer">
                                                                    <VitalSignCard
                                                                        id={"oxygen_saturation" + i}
                                                                        name="Oxygen saturation"
                                                                        value={reading.value}
                                                                        unit="%"
                                                                        timestamp={reading.timestamp}
                                                                        timestampWarning={reading.timestampWarning}
                                                                        consultationTooltipOn={inConsultation}
                                                                    />
                                                                </div>
                                                            )
                                                        )
                                                    })}
                                                </div>
                                            ) : (
                                                <div className="missingHistory">No historical measurements are available for oxygen saturation.</div>
                                            )}
                                        </Collapse>
                                    </div>
                                </React.Fragment>
                            ) : null}
                        </React.Fragment>
                    ) : (
                        <h3>No medical data found</h3>
                    )}
                </div>
            </div>
        ) : null
    }
}

MedicalData = connect(
    (state, props) => {
        let medicalData = undefined
        let inConsultation = props.match.params.waitlistID && props.match.params.itemID
        if (state.patient.patientRecords.data) {
            let records = state.patient.patientRecords.data
            // sort records by creation time and reverse to have latest record as first
            records = _.reverse(
                _.sortBy(records, [
                    function(obj) {
                        return obj.meta.created
                    }
                ])
            )

            medicalData = {}
            // if in consultation, fetch data from current consultation as well
            if (inConsultation && state.waitlist.item) {
                _.forEach(state.waitlist.item.vitalSigns, (obj, key) => {
                    // set for warning only if timestamp is missing
                    let vitalSign = _.clone(obj)
                    vitalSign.timestampWarning = vitalSign.timestamp ? false : true
                    medicalData[key] = medicalData[key] || []
                    medicalData[key].push(vitalSign)
                })
            }

            // collect latest data for each category
            _.forEach(records, ({ data, meta }) => {
                _.forEach(data.vitalSigns, (obj, key) => {
                    // if in consultation, mark historical data with warning
                    let vitalSign = _.clone(obj)
                    vitalSign.timestamp = vitalSign.timestamp ? vitalSign.timestamp : meta.created
                    vitalSign.timestampWarning = vitalSign.timestamp ? inConsultation : true
                    medicalData[key] = medicalData[key] || []
                    medicalData[key].push(vitalSign)
                })
            })
        }

        return {
            inConsultation: inConsultation,
            patientID: props.match.params.patientID || state.patient.patient.ID,
            loadedPatientID: state.patient.patient.ID,
            medicalData: medicalData,
            medicalDataLoading: state.patient.patientRecords.loading,
            patientRecords: state.patient.patientRecords.data,
            canSeeVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[READ],
            canAddVitalSigns: ((state.validations.userRights || {})[RESOURCE_VITAL_SIGNS] || {})[WRITE]
        }
    },
    { fetchHealthRecords }
)(MedicalData)

export default MedicalData
