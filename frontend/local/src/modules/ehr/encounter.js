export default dispatch =>
    Promise.resolve([
        // category
        {
            type: "fixedValue",
            ehrPath: "/content[openEHR-EHR-COMPOSITION.encounter.v1]/category",
            value: "openehr::433"
        },

        // diagnoses
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-EVALUATION.problem_diagnosis.v1]",
            formPath: "diagnoses",
            items: [
                // name
                {
                    type: "value",
                    ehrPath: "/data/items[at0001]/item[at0002]",
                    formPath: "diagnosis.label"
                },
                // comment
                {
                    type: "value",
                    ehrPath: "/data/items[at0001]/item[at0009]",
                    formPath: "comment"
                }
            ]
        },

        // treatments
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-INSTRUCTION.medication_order.v2]",
            formPath: "therapies",
            items: [
                // medication
                {
                    type: "value",
                    ehrPath: "/activities[at0001]/description[at0002]/items[at0070]",
                    formPath: "medication"
                },
                // comment
                {
                    type: "value",
                    ehrPath: "/activities[at0001]/description[at0002]/items[at0044]",
                    formPath: "instructions"
                },
                // link to diagnosis (#num)
                {
                    type: "value",
                    ehrPath: "/activities[at0001]/description[at0002]/items[at0167]",
                    formPath: "diagnosis"
                }
            ]
        },

        // main complaint
        {
            type: "value",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-EVALUATION.complaint.v1]/items[at0001]/item[at0002]",
            formPath: "mainComplaint.complaint"
        },
        {
            type: "value",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-EVALUATION.complaint.v1]/items[at0001]/item[at0003]",
            formPath: "mainComplaint.comment"
        },

        // vital signs

        // weight (QUANTITY)
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.body_weight.v2]/data[at0002]/event[at0003]/items[at0001]",
            unit: "kg",
            formPath: "vitalSigns.weight.value"
        },

        // height
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.height.v2]/data[at0001]/events[at0002]/data[at0003]/items[at0004]",
            unit: "cm",
            formPath: "vitalSigns.height.value"
        },

        // bmi
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.body_mass_index.v2]/data[at0001]/events[at0002]/data[at0003]/items[at0004]",
            unit: "kg/m2",
            formPath: "vitalSigns.bmi.value"
        },

        // temperature
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.body_temperature.v2]/data[at0002]/events[at0003]/data[at0001]/items[at0004]",
            unit: "Cel",
            formPath: "vitalSigns.temperature.value"
        },

        // blood pressure
        // systolic
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.blood_pressure.v1]/data[at0001]/events[at0006]/data[at0003]/items[at0004]",
            unit: "mm[Hg]",
            formPath: "vitalSigns.pressure.value.systolic"
        },
        // diastolic
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.blood_pressure.v1]/data[at0001]/events[at0006]/data[at0003]/items[at0005]",
            unit: "mm[Hg]",
            formPath: "vitalSigns.pressure.value.diastolic"
        },

        // pulse
        {
            type: "fixedValue",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.pulse.v1]/data[at0002]/events[at0003]/data[at0001]/items[at0004]|name",
            value: "local::at1027"
        },
        {
            type: "quantity",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.pulse.v1]/data[at0002]/events[at0003]/data[at0001]/items[at0004]",
            unit: "/min",
            formPath: "vitalSigns.heart_rate.value"
        },

        // oxygen saturation (0 - 100)
        {
            type: "value",
            ehrPath:
                "/content[openEHR-EHR-COMPOSITION.encounter.v1]/context/other_context/items[openEHR-EHR-OBSERVATION.pulse_oximetry.v1]/data[at0001]/events[at0002]/data[at0003]/items[at0006]",
            formPath: "vitalSigns.oxygen_saturation.value"
        }
    ])
