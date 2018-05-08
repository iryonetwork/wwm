import { load as loadCode } from "shared/modules/codes"

export default dispatch => {
    return Promise.all([dispatch(loadCode("countries"))]).then(([countries]) => [
        /* Chronic diseases (array) */
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]",
            formPath: "chronicDiseases",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0018]",
                    formPath: "disease"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0017]",
                    formPath: "date"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0016]",
                    formPath: "comment"
                }
            ]
        },

        /* Immunisations (array) */
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]",
            formPath: "immunizations",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0019]",
                    formPath: "immunization"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0021]",
                    formPath: "date"
                }
            ]
        },

        // /* Allergies (array) */
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]",
            formPath: "allergies",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0010]",
                    formPath: "allergy"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0012]",
                    formPath: "critical"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0013]",
                    formPath: "comment"
                }
            ]
        },

        // // MEDICAL HISTORY
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]",
            formPath: "injuries",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0023]",
                    formPath: "injury"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0024]",
                    formPath: "date"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0025]",
                    formPath: "medication"
                }
            ]
        },

        /* Surgeries (array) */
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]",
            formPath: "surgeries",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0028]",
                    formPath: "injury"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0029]",
                    formPath: "date"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0030]",
                    formPath: "medication"
                }
            ]
        },

        /* Medications (array) */
        {
            type: "array",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]",
            formPath: "medications",
            items: [
                {
                    type: "value",
                    ehrPath: "/items[at0032]",
                    formPath: "medication"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0034]",
                    formPath: "comment"
                },
                {
                    type: "value",
                    ehrPath: "/items[at0016]",
                    formPath: "comment"
                }
            ]
        },

        /* Additional patient info */

        // Number of kids
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0047]",
            formPath: "numberOfKids"
        },
        // Nationality
        {
            type: "code",
            codes: countries,
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0048]",
            formPath: "nationality"
        },
        // Country of origin
        {
            type: "code",

            codes: countries,
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0049]",
            formPath: "countryOfOrigin"
        },
        // Education (@TODO code)
        {
            type: "fixedValue",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0050]",
            value: ">education::secondary|Secondary education|<"
        },
        // Occupation
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0058]",
            formPath: "profession"
        },
        // Date of leaving home country
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0059]",
            formPath: "dateOfLeaving"
        },
        // Date of arriving to camp
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0061]",
            formPath: "dateOfArrival"
        },
        // Transit countries
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0060]",
            formPath: "transitCountries"
        },
        // ])

        // Habits and conditions

        // Are you smoking
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0039]",
            formPath: "habits_smoking"
        },
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0038]",
            formPath: "habits_smoking_comment"
        },
        // Taking drugs
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0051]/items[at0040]",
            formPath: "habits_drugs"
        },
        // Resources for basic hygiene
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0052]/items[at0041]",
            formPath: "conditions_basic_hygiene"
        },
        // Access to clean water
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0053]/items[at0042]",
            formPath: "conditions_clean_water"
        },
        // Sufficient food supply
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0054]/items[at0043]",
            formPath: "conditions_food_supply"
        },
        // Good appetite
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0055]/items[at0044]",
            formPath: "conditions_good_appetite"
        },
        // Accommodations have heating
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0056]/items[at0045]",
            formPath: "conditions_heating"
        },
        // Accommodations have electricity
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0057]/items[at0046]",
            formPath: "conditions_electricity"
        },
        // ])

        /* Vaccine information */

        // On schedule at home (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0065]",
            formPath: "vaccinationUpToDate"
        },
        // Has Immunization documents (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0066]",
            formPath: "vaccinationCertificates"
        },
        // Tested for tuberculosis (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0067]",
            formPath: "tuberculosisTested"
        },
        // Were tuberculosis tests positive (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0068]",
            formPath: "tuberculosisTestResult"
        },
        // Any additional tests done (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0069]",
            formPath: "tuberculosisAdditionalInvestigationDetails"
        },
        // Investigation details
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0070]",
            formPath: "tuberculosisAdditionalInvestigation"
        },
        // Any reaction to vaccines (BOOLEAN)
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0071]",
            formPath: "vaccinationReaction"
        },
        // Details of vaccine reactions
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0072]",
            formPath: "vaccinationReactionDetails"
        },

        // // /* BABY SCREENING */

        // Delivery type ( // @TODO code)
        {
            type: "code",
            codes: [],
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0075]",
            formPath: "deliveryType"
        },
        // Premature
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0076]",
            formPath: "prematurity"
        },
        // Weeks at birth
        {
            type: "integer",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0077]",
            formPath: "weeksAtBirth"
        },
        // Weight at birth
        {
            type: "quantity",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0078]",
            formPath: "weightAtBirth",
            unit: "gm"
        },
        // Height at birth
        {
            type: "quantity",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0079]",
            formPath: "heightAtBirth",
            unit: "cm"
        },
        // Breastfeeding
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0081]",
            formPath: "breastfeeding"
        },
        // Breastfeeding for how long
        {
            type: "integer",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0082]",
            formPath: "breastfeedingDuration"
        },
        // What does baby eat or drink ( // @TODO code)
        {
            type: "code",
            codes: [],
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0083]",
            formPath: "babyEatsAndDrinks"
        },
        // How many diapers does child wet
        {
            type: "integer",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0084]",
            formPath: "babyWetDiapers"
        },
        // How many times does child have bowl movement (CODED) ( // @TODO code)
        {
            type: "code",
            codes: [],
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0085]",
            formPath: "babyBowelMovements"
        },
        // Describe bowl movement
        {
            type: "integer",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0086]",
            formPath: "babyBowelMovementsComment"
        },
        // Satisfied with sleep
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0089]",
            formPath: "babySleep"
        },
        // Comment about the sleep
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0090]",
            formPath: "babySleepComment"
        },
        // Do you take / give vit. D
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0091]",
            formPath: "babyVitaminD"
        },
        // Baby sleeps on her back
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0092]",
            formPath: "babySleepOnBack"
        },
        // Does anyone smoke
        {
            type: "boolean",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0093]",
            formPath: "babyAnyoneSmokes"
        },
        // Number of smokers
        {
            type: "integer",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0094]",
            formPath: "babyNumberOfSmokers"
        },
        // How does child get around
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0095]",
            formPath: "babyGetsAround"
        },
        // How does child communicate
        {
            type: "value",
            ehrPath: "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0096]",
            formPath: "babyCommunicates"
        }
    ])
}
