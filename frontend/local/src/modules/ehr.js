import { get, has } from 'lodash'
import { read, CLINIC_ID, LOCATION_ID } from 'shared/modules/config'
import { load as loadClinic } from './clinics'
import { load as loadCode } from 'shared/modules/codes'
import { load as loadLocation } from './locations'
import { load as loadUser } from './users'

// Converts form data into two separate documents
export const composePatientData = (formData) => (
    (dispatch) => {
        return dispatch(buildContextForPatientData(formData))
            .then(context => {
                return Promise.all([
                    buildPersonData(context, formData),
                    buildInfoData(context, formData),
                ])
                    .then(([person, info]) => ({person,info}))
            })
    }
)

const buildContextForPatientData = (formData) => (dispatch) => {
    return Promise.all([
        dispatch(loadClinic(read(CLINIC_ID))),
        dispatch(loadLocation(read(LOCATION_ID))),
        dispatch(loadUser('me')), // doctor
    ])
        .then(([clinic, location, doctor]) => {

            return {
                // facility details
                "/context/health_care_facility|name": clinic.name,
                "/context/health_care_facility|identifier": clinic.id,
                "/territory": location.country,
                "/language":"en",

                // time info
                "/context/start_time":  (new Date()).toJSON(),
                "/context/end_time": (new Date()).toJSON(),

                // // participants
                // // add doctor
                // "/composer|identifier": doctor.id,
                // "/composer|name": `${doctor.personalData.firstName} ${doctor.personalData.lastName}`,

                "/category":"openehr::431|persistent|",
            }
        })
}

const buildPersonData = (ctx, formData) => (dispatch) => {
    return Promise.all([
        dispatch(loadCode('countries')),
        dispatch(loadCode('gender')),
        dispatch(loadCode('maritalStatus')),
    ])
        .then(([countries, genders, maritalStatuses]) => {
            let data = compose({}, formData, [
                // First Name
                assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0002]",
                    "firstName"
                ),
                // Last Name
                assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0003]",
                    "lastName"
                ),
                // Preferred Name
                assignFixedValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0008]",
                    "true"
                ),
                // Date of Birth
                assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0010]",
                    "dateOfBirth"
                ),
                // Country of Birth
                assignCode(
                    countries,
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]",
                    "country"
                ),
                // Gender
                assignCode(
                    genders,
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]",
                    "gender"
                ),
                // Marital Status
                assignCode(
                    maritalStatuses,
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]",
                    codeToString(formData.maritalStatus, maritalStatuses)
                ),
            	/* ADDRESS */
                // Type of address
                assignFixedValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0033]",
                    "local::at0463|Temporary Accommodation|"
                ),
                // Country
                assignCode(
                    countries,
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at0009]",
                    "country"
                ),
                // Camp (stored under "Address site name")
                assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00014]",
                    "camp"
                ),
                // Tent (stored under "Building/Complex sub-unit number")
                assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at00013]",
                    "tent"
                ),

            	/* Phone number */
                // Type of address (phone)
                assignFixedValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/name[at0014]",
                    "local::at0022|Mobile|"
                ),
                // Phone number
            	assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/details[at0001]/items[at0007]",
                    "phone"
                ),

            	/* Email address */
                // Type of address (email)
                assignFixedValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/name[at0014]",
                    "local::at0024|Email|"
                ),
                // Email address
            	assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]",
                    "email"
                ),

            	/* Whatsapp */
                // Type of address (whatsapp)
                assignFixedValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/name[at0013]",
                    "whatsapp"
                ),
                // Address
            	assignValue(
                    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]",
                    "whatsapp"
                ),
            ]);

            // add identities
            (formData.documents || []).forEach((doc, i) => {
                data = compose(data, doc, [
                    // document id
                    assignValue(
                        `/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:${i}|id`,
                        "id"
                    ),
                    assignValue(
                        `/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:${i}|type`,
                        "type"
                    ),
                ])
            });

            // 	// Relationships
            // 	// type of relationship
            //     "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0/details[at0041]": ">relationships::stranger|Thing|<",
            //     "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|namespace": "local-together", (local-together, local)
            //     "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|type": "PERSON",
            //     "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/relationships[at0004]:0|id": ">PERSON-UUID<"
            // }

            return data
        })
        // .catch(ex => {
        //     console.log('why?!?', ex)
        // })
}

const buildInfoData = (ctx, formData) => (dispatch) => {
    return Promise.all([
        dispatch(loadCode('countries')),
    ])
        .then(([countries]) => {
            let data = {};

            /* Chronic diseases (array) */
            (formData.chronicDiseases || []).forEach((el, i) => {
                data = compose(data, el, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:${i}/items[at0018]`,
                        "disease"
                    ),
                    // Date
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:${i}/items[at0017]`,
                        "date"
                    ),
                    // Comment
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:${i}/items[at0016]`,
                        "comment"
                    ),
                ])
            });

            /* Immunisations (array) */
            (formData.immunizations || []).forEach((el, i) => {
                data = Object.assign(data, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:${i}/items[at0019]`,
                        "immunization"
                    ),
                    // Date
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:${i}/items[at0021]`,
                        "date"
                    ),
                ])
            });

            // /* Allergies (array) */
            (formData.allergies || []).forEach((el, i) => {
                data = Object.assign(data, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:${i}/items[at0010]`,
                        "allergy"
                    ),
                    // Critical (BOOL)
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at00009]:${i}/items[at0012]`,
                        "critical"
                    ),
                    // Comment
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:${i}/items[at0013]`,
                        "comment"
                    ),
                ])
            });


            // // MEDICAL HISTORY

            /* Injuries or handicaps (array) */
            (formData.injuries || []).forEach((el, i) => {
                data = Object.assign(data, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:${i}/items[at0023]`,
                        "injury"
                    ),
                    // Date
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:${i}/items[at0024]`,
                        "date"
                    ),
                    // Comment
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:${i}/items[at0025]`,
                        "medication"
                    ),
                ])
            });

            /* Surgeries (array) */
            (formData.surgeries || []).forEach((el, i) => {
                data = Object.assign(data, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:${i}/items[at0028]`,
                        "injury"
                    ),
                    // Date
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:${i}/items[at0029]`,
                        "date"
                    ),
                    // Comment
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:${i}/items[at0030]`,
                        "medication"
                    ),
                ])
            });

            /* Medications (array) */
            (formData.medications || []).forEach((el, i) => {
                data = Object.assign(data, [
                    // Name
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:${i}/items[at0032]`,
                        "medication"
                    ),
                    // Comment
                    assignValue(
                        `/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:${i}/items[at0034]`,
                        "comment"
                    ),
                ])
            });

            /* Additional patient info */

            data = compose(data, formData, [
                // Number of kids
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0047]",
                    "numberOfKids"
                ),
                // Nationality
                assignCode(
                    countries,
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0048]",
                    "nationality"
                ),
                // Country of origin
                assignCode(
                    countries,
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0049]",
                    "countryOfOrigin"
                ),
                // Education
                assignValue( // @TODO
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0050]",
                    ">education::secondary|Secondary education|<"
                ),
                // Occupation
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0058]",
                    "profession"
                ),
                // Date of leaving home country
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0059]",
                    "dateOfLeaving"
                ),
                // Date of arriving to camp
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0061]",
                    "dateOfArrival"
                ),
                // Transit countries
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0060]",
                    "transitCountries"
                ),
            ])

            // Habits and conditions

            data = compose(data, formData, [
                // Are you smoking
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0039]",
                    "habits_smoking"
                ),
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0038]",
                    "habits_smoking_comment"
                ),
                // Taking drugs
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0051]/items[at0040]",
                    "habits_drugs"
                ),
                // Resources for basic hygiene
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0052]/items[at0041]",
                    "conditions_basic_hygiene"
                ),
                // Access to clean water
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0053]/items[at0042]",
                    "conditions_clean_water"
                ),
                // Sufficient food supply
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0054]/items[at0043]",
                    "conditions_food_supply"
                ),
                // Good appetite
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0055]/items[at0044]",
                    "conditions_good_appetite"
                ),
                // Accommodations have heating
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0056]/items[at0045]",
                    "conditions_heating"
                ),
                // Accommodations have electricity
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0057]/items[at0046]",
                    "conditions_electricity"
                ),
            ]);

            /* Vaccine information */

            data = compose(data, formData, [
                // On schedule at home (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0065]",
                    "vaccinationUpToDate"
                ),
                // Has Immunization documents (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0066]",
                    "vaccinationCertificates"
                ),
                // Tested for tuberculosis (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0067]",
                    "tuberculosisTested"
                ),
                // Were tuberculosis tests positive (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0068]",
                    "tuberculosisTestResult"
                ),
                // Any additional tests done (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0069]",
                    "tuberculosisAdditionalInvestigationDetails"
                ),
                // Investigation details
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0070]",
                    "tuberculosisAdditionalInvestigation"
                ),
                // Any reaction to vaccines (BOOLEAN)
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0071]",
                    "vaccinationReaction"
                ),
                // Details of vaccine reactions
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0064]/items[at0072]",
                    "vaccinationReactionDetails"
                ),
            ])

            // /* BABY SCREENING */
            data = compose(data, formData, [
                // Delivery type
                assignCode(
                    [], // @TODO code
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0075]",
                    "deliveryType"
                ),
                // Premature
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0076]",
                    "prematurity"
                ),
                // Weeks at birth
                assignInteger(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0077]",
                    "weeksAtBirth"
                ),
                // Weight at birth
                assignQuantity(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0078]",
                    "weightAtBirth",
                    "gm"
                ),
                // Height at birth
                assignQuantity(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0079]",
                    "heightAtBirth",
                    "cm"
                ),
                // Breastfeeding
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0081]",
                    "breastfeeding"
                ),
                // Breastfeeding for how long
                assignInteger(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0087]/items[at0082]",
                    "breastfeedingDuration"
                ),
                // What does baby eat or drink
                assignCode(
                    [], // @TODO code
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0083]",
                    "babyEatsAndDrinks"
                ),
                // How many diapers does child wet
                assignInteger(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0084]",
                    "babyWetDiapers"
                ),
                // How many times does child have bowl movement (CODED)
                assignCode(
                    [], // @TODO code
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0085]",
                    "babyBowelMovements"
                ),
                // Describe bowl movement
                assignInteger(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0086]",
                    "babyBowelMovementsComment"
                ),
                // Satisfied with sleep
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0089]",
                    "babySleep"
                ),
                // Comment about the sleep
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0088]/items[at0090]",
                    "babySleepComment"
                ),
                // Do you take / give vit. D
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0091]",
                    "babyVitaminD"
                ),
                // Baby sleeps on her back
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0092]",
                    "babySleepOnBack"
                ),
                // Does anyone smoke
                assignBoolean(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0093]",
                    "babyAnyoneSmokes"
                ),
                // Number of smokers
                assignInteger(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0094]",
                    "babyNumberOfSmokers"
                ),
                // How does child get around
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0095]",
                    "babyGetsAround"
                ),
                // How does child communicate
                assignValue(
                    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0080]/items[at0096]",
                    "babyCommunicates"
                )
            ])

            return data
        })
}

const compose = (data, formData, fns) => {
    return fns.reduce((acc, fn) => fn(acc, formData), data)
}

const assignValue = (ehrPath, formPath, staticValue) => (
    (data, formData) => {
        if (!has(formData, formPath)) {
            return data
        }
        return Object.assign(data, {[ehrPath]: staticValue || get(formData, formPath)})
    }
)

const assignInteger = (ehrPath, formPath) => (
    (data, formData) => {
        if (!has(formData, formPath)) {
            return data
        }
        return Object.assign(data, {[ehrPath]: parseInt(get(formData, formPath), 10)})
    }
)

const assignQuantity = (ehrPath, formPath, unit) => (
    (data, formData) => {
        if (!has(formData, formPath)) {
            return data
        }
        return Object.assign(data, {[ehrPath]: `${get(formData, formPath)},${unit}`})
    }
)

const assignBoolean = (ehrPath, formPath) => (
    (data, formData) => {
        if (!has(formData, formPath)) {
            return data
        }
        return Object.assign(data, {[ehrPath]: get(formData, formPath) ? "true" : "false"})
    }
)

const assignFixedValue = (ehrPath, value) => (
    (data, formData) => {
        return Object.assign(data, {[ehrPath]: value})
    }
)

const assignCode = (codes, ehrPath, formPath) => (
    (data, formData) => {
        if (!has(formData, formPath)) {
            return data
        }
        return Object.assign(data, {[ehrPath]: codeToString(get(formData, formPath), codes)})
    }
)

export const codeToString = (key, codes) => {
    for (let i=0; i < codes.length; i++) {
        const el = codes[i]
        if (el.id === key) {
            let category = el.category
            let id = el.id
            if (id.indexOf('SNOMED-') === 0) {
                category = 'SNOMED'
                id = id.substring(7)
            } else if (id.indexOf('CODED-') === 0) {
                category = 'local'
                id = id.substring(6)
            }

            return `${category}::${id}|${el.title}|`
        }
    }
}
