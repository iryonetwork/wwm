import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import classnames from "classnames"

import { save } from "shared/modules/config"
import {
    LOCALE,
    LENGTH_UNIT,
    WEIGHT_UNIT,
    TEMPERATURE_UNIT,
    BLOOD_PRESSURE_UNIT,
    localeOptions,
    lengthUnitOptions,
    weightUnitOptions,
    temperatureUnitOptions,
    bloodPressureUnitOptions
} from "../../modules/config"

import "./style.css"

const Setting = ({ value, disabled, label, placeholder, options, onChange }) => (
    <label>
        <select disabled={disabled} className={classnames("form-control", { selected: value !== "" })} onChange={onChange} value={value}>
            <option value="" disabled>
                {placeholder ? placeholder : label}
            </option>
            {(options || []).map(option => (
                <option value={option.value} key={option.value}>
                    {option.label}
                </option>
            ))}
        </select>

        <span>{label}</span>
    </label>
)

class SettingsContent extends React.Component {
    save(key) {
        return e => {
            this.props.save(key, e.target.value)
        }
    }

    render() {
        return (
            <div className="settingsBody">
                <div className="section">
                    <h3>General</h3>
                    <div className="sectionBody">
                        <div className="form-row">
                            <div className={classnames("form-group", { "col-sm-6": this.props.wideInput, "col-sm-4": !this.props.wideInput })}>
                                <Setting
                                    name="locale"
                                    options={localeOptions}
                                    label="Language"
                                    placeholder="Language"
                                    value={this.props.localeValue}
                                    onChange={this.save(LOCALE)}
                                />
                            </div>
                        </div>
                    </div>
                </div>
                <div className="section">
                    <h3>Units</h3>
                    <div className="sectionBody">
                        <div className="form-row">
                            <div className={classnames("form-group", { "col-sm-6": this.props.wideInput, "col-sm-4": !this.props.wideInput })}>
                                <Setting
                                    name="length"
                                    options={lengthUnitOptions}
                                    label="Length"
                                    placeholder="Length unit"
                                    value={this.props.lengthUnitValue}
                                    onChange={this.save(LENGTH_UNIT)}
                                />
                            </div>
                            <div className={classnames("form-group", { "col-sm-6": this.props.wideInput, "col-sm-4": !this.props.wideInput })}>
                                <Setting
                                    name="weight"
                                    options={weightUnitOptions}
                                    label="Weight"
                                    placeholder="Weight unit"
                                    value={this.props.weightUnitValue}
                                    onChange={this.save(WEIGHT_UNIT)}
                                />
                            </div>
                        </div>
                        <div className="form-row">
                            <div className={classnames("form-group", { "col-sm-6": this.props.wideInput, "col-sm-4": !this.props.wideInput })}>
                                <Setting
                                    name="temperature"
                                    options={temperatureUnitOptions}
                                    label="Temperature"
                                    placeholder="Temperature unit"
                                    value={this.props.temperatureUnitValue}
                                    onChange={this.save(TEMPERATURE_UNIT)}
                                />
                            </div>
                            <div className={classnames("form-group", { "col-sm-6": this.props.wideInput, "col-sm-4": !this.props.wideInput })}>
                                <Setting
                                    name="bloodPressure"
                                    options={bloodPressureUnitOptions}
                                    label="Blood pressure"
                                    placeholder="Blood pressure unit"
                                    value={this.props.bloodPressureUnitValue}
                                    onChange={this.save(BLOOD_PRESSURE_UNIT)}
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        localeValue: state.config[LOCALE] || "",
        lengthUnitValue: state.config[LENGTH_UNIT] || "",
        weightUnitValue: state.config[WEIGHT_UNIT] || "",
        temperatureUnitValue: state.config[TEMPERATURE_UNIT] || "",
        bloodPressureUnitValue: state.config[BLOOD_PRESSURE_UNIT] || ""
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            save
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(SettingsContent)
