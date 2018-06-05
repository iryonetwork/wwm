import React from "react"
// import { expect } from "chai"
import { shallow, mount } from "enzyme"
import sinon from "sinon"
import { load as loadClinic } from "./clinics"
import { load as loadCodes } from "shared/modules/codes"
import { load as loadLocation } from "./locations"
import { load as loadUser } from "./users"
jest.mock("./users")
jest.mock("./clinics")
jest.mock("./locations")
jest.mock("shared/modules/codes")

import { composePatientData, codeToString, extractPatientData } from "./ehr"

const fullFormData = {
    documents: [
        {
            type: "syrian-id",
            number: "1234567"
        },
        {
            type: "un-id",
            number: "987654"
        }
    ],
    firstName: "Dominik",
    lastName: "Znidar",
    dateOfBirth: "1983-05-21",
    gender: "CODED-at0310",
    maritalStatus: "SNOMED-125681006",
    numberOfKids: "1",
    nationality: "SI",
    countryOfOrigin: "SI",
    education: "gimnasium,",
    profession: "developer",
    country: "SI",
    region: "Kranj",
    phone: "040123456",
    email: "email@test.com",
    whatsapp: "987987654",
    dateOfLeaving: "2016-01-01",
    transitCountries: "Slovenia",
    dateOfArrival: "2016-12-31",
    people_in_family: "1",
    people_living_together: "2",
    familyMembers: [
        {
            firstName: "Test",
            lastName: "Test",
            dateOfBirth: "1993-05-04",
            relation: "child",
            livingTogether: "true",
            documentType: "syrian_id",
            documentNumber: "1234567"
        }
    ],
    allergies: [
        {
            allergy: "allergy 1",
            comment: "allergy comment 1",
            critical: "false"
        },
        {
            allergy: "allergy 2",
            comment: "allergy comment 2",
            critical: "true"
        }
    ],
    immunizations: [
        {
            immunization: "Immunization 1",
            date: "2014-02-01"
        },
        {
            immunization: "Immunization 2",
            date: "2004-03-02"
        }
    ],
    chronicDiseases: [
        {
            disease: "Chronic 1",
            date: "2004-03-12",
            medication: "Medication 1"
        },
        {
            disease: "Chronic 2 ",
            date: "2005-04-23",
            medication: "Medication 2"
        }
    ],
    injuries: [
        {
            injury: "Injury 1",
            date: "2007-06-05",
            medication: "Aids 1"
        },
        {
            injury: "Injury 2",
            date: "2009-08-07",
            medication: "Aids 2"
        }
    ],
    surgeries: [
        {
            injury: "Surgery 1",
            date: "2004-03-12",
            medication: "Comment 1"
        },
        {
            injury: "Surgery 2 ",
            date: "2004-03-12",
            medication: "Comment 2"
        }
    ],
    medications: [
        {
            medication: "Medication 1",
            comment: "Comment 1"
        },
        {
            medication: "Medication 2",
            comment: "Comment 2"
        }
    ],
    habits_smoking: "true",
    habits_smoking_comment: "10 boxes a day",
    habits_drugs: "true",
    habits_drugs_comment: "All of them",
    conditions_basic_hygiene: "true",
    conditions_heating: "true",
    conditions_good_appetite: "true",
    conditions_food_supply: "true",
    conditions_clean_water: "true",
    conditions_electricity: "true"
}

const personDocument = {
    "/context/health_care_facility|name": "CLINIC1",
    "/context/health_care_facility|identifier": "e4ebb41b-7c62-4db7-9e1c-f47058b96dd0",
    "/territory": "CNT",
    "/language": "en",
    "/context/start_time": "2018-05-07T21:53:45.635Z",
    "/context/end_time": "2018-05-07T21:53:45.635Z",
    "/category": "openehr::431|persistent|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:0|id":
        "1234567",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:0|type":
        "syrian-id",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:1|id":
        "987654",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0005]/items[openEHR-DEMOGRAPHIC-CLUSTER.person_identifier.v1]/item[at0001]:1|type":
        "un-id",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0002]": "Dominik",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0003]": "Znidar",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/identities[openEHR-DEMOGRAPHIC-PARTY_IDENTITY.person_name.v1]/details[at0001]/items[at0008]": "true",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0010]": "1983-05-21",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0012]": "category::SI|Slovenia|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0017]": "local::at0310|Male|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/details[openEHR-DEMOGRAPHIC-ITEM_TREE.person_details.v1.0.0]/items[at0033]": "SNOMED::125681006|Married|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0033]":
        "local::at0463|Temporary Accommodation|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0002]/items[at0009]":
        "category::SI|Slovenia|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.address.v1]:0/details[at0001]/items[at0003]": "Kranj",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/name[at0014]":
        "local::at0022|Mobile|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:1/details[at0001]/items[at0007]":
        "040123456",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/name[at0014]":
        "local::at0024|Email|",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:2/details[at0001]/items[at0007]":
        "email@test.com",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/name[at0013]": "whatsapp",
    "/content[openEHR-DEMOGRAPHIC-PERSON.person.v1]/contacts[openEHR-DEMOGRAPHIC-ADDRESS.electronic_communication.v1.0.0]:3/details[at0001]/items[at0007]":
        "987987654"
}

const infoDocument = {
    "/context/health_care_facility|name": "CLINIC1",
    "/context/health_care_facility|identifier": "e4ebb41b-7c62-4db7-9e1c-f47058b96dd0",
    "/territory": "CNT",
    "/language": "en",
    "/context/start_time": "2018-05-07T21:53:45.635Z",
    "/context/end_time": "2018-05-07T21:53:45.635Z",
    "/category": "openehr::431|persistent|",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:0/items[at0018]": "Chronic 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:0/items[at0017]": "2004-03-12",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:1/items[at0018]": "Chronic 2 ",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0015]:1/items[at0017]": "2005-04-23",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:0/items[at0019]": "Immunization 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:0/items[at0021]": "2014-02-01",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:1/items[at0019]": "Immunization 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0014]:1/items[at0021]": "2004-03-02",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:0/items[at0010]": "allergy 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:0/items[at0012]": "false",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:0/items[at0013]": "allergy comment 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:1/items[at0010]": "allergy 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:1/items[at0012]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0009]:1/items[at0013]": "allergy comment 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0023]": "Injury 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0024]": "2007-06-05",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:0/items[at0025]": "Aids 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:1/items[at0023]": "Injury 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:1/items[at0024]": "2009-08-07",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0022]:1/items[at0025]": "Aids 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0028]": "Surgery 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0029]": "2004-03-12",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:0/items[at0030]": "Comment 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:1/items[at0028]": "Surgery 2 ",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:1/items[at0029]": "2004-03-12",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0026]:1/items[at0030]": "Comment 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0032]": "Medication 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0034]": "Comment 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:0/items[at0016]": "Comment 1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:1/items[at0032]": "Medication 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:1/items[at0034]": "Comment 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0031]:1/items[at0016]": "Comment 2",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0047]": "1",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0048]": "category::SI|Slovenia|",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0049]": "category::SI|Slovenia|",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0050]": ">education::secondary|Secondary education|<",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0058]": "developer",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0059]": "2016-01-01",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0061]": "2016-12-31",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0062]/items[at0060]": "Slovenia",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0039]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0036]/items[at0038]": "10 boxes a day",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0051]/items[at0040]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0052]/items[at0041]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0053]/items[at0042]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0054]/items[at0043]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0055]/items[at0044]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0056]/items[at0045]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0036]/items[at0057]/items[at0046]": "true",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0078]": "weightAtBirth",
    "/content[openEHR-EHR-ITEM_TREE.patient_info.v0]/items[at0073]/items[at0074]/items[at0079]": "heightAtBirth"
}

// const waitlistItem = {
//     added: '2018-05-02T15:30:13.435Z',
// }

if (!global.window.localStorage) {
    global.window.localStorage = {
        getItem() {
            return "{}"
        },
        setItem() {}
    }
}

beforeEach(() => {
    loadClinic.mockClear()
    loadUser.mockClear()
    loadLocation.mockClear()
    loadCodes.mockClear()
})

describe("ehr", () => {
    describe("composePatientData", () => {
        it("should return object with person and info keys", done => {
            loadClinic.mockResolvedValue({
                id: "e4ebb41b-7c62-4db7-9e1c-f47058b96dd0",
                name: "CLINIC1"
            })
            loadUser.mockResolvedValue({
                id: "c12574e4-acd4-4266-9d53-b614c8a942bc",
                personalData: {
                    firstName: "Doctor",
                    lastName: "X"
                }
            })
            loadCodes.mockResolvedValue([
                { category: "category", id: "SH", locale: "en", title: "Saint Helena" },
                { category: "category", id: "SI", locale: "en", title: "Slovenia" },
                { category: "gender", id: "CODED-at0310", locale: "en", title: "Male" },
                { category: "maritalStatus", id: "SNOMED-125681006", locale: "en", title: "Married" }
            ])
            loadLocation.mockResolvedValue({
                country: "CNT"
            })

            const getState = () => ({
                locations: {
                    cache: {}
                },
                authentication: {
                    tokenString: "TOKEN"
                }
            })

            const dispatch = fn => {
                // console.log(fn, new Error().stack)
                return typeof fn === "function" ? fn(dispatch, getState) : fn
            }

            dispatch(composePatientData(fullFormData))
                .then(out => {
                    expect(Object.keys(out)).toEqual(["person", "info"])

                    // copy over timestamps
                    out.person["/context/start_time"] = personDocument["/context/start_time"]
                    out.person["/context/end_time"] = personDocument["/context/end_time"]
                    out.info["/context/start_time"] = infoDocument["/context/start_time"]
                    out.info["/context/end_time"] = infoDocument["/context/end_time"]
                    // expect(out.person).toEqual(personDocument) // @TODO should match
                    // expect(out.info).toEqual(infoDocument) // @TODO should match
                    expect(loadUser).toHaveBeenLastCalledWith("me")
                    expect(loadClinic).toHaveBeenLastCalledWith("e4ebb41b-7c62-4db7-9e1c-f47058b96dd0")
                    expect(loadCodes).toHaveBeenCalledTimes(4)
                    expect(loadCodes).toHaveBeenCalledWith("gender")
                    expect(loadCodes).toHaveBeenCalledWith("countries")
                    expect(loadCodes).toHaveBeenCalledWith("maritalStatus")
                    // console.log(out)
                    done()
                })
                .catch(ex => {
                    // console.log(ex)
                    expect(ex).toBeUndefined()
                })
        })

        //
    })

    describe("extractPatientData", () => {
        it("should return an object", done => {
            loadClinic.mockResolvedValue({})
            loadUser.mockResolvedValue({})
            loadCodes.mockResolvedValue([])
            loadLocation.mockResolvedValue({})

            const getState = () => ({})

            const dispatch = fn => {
                return typeof fn === "function" ? fn(dispatch, getState) : fn
            }

            dispatch(extractPatientData(personDocument, infoDocument))
                .then(out => {
                    // expect(out).toEqual(fullFormData) // @TODO should match!
                    expect(loadCodes).toHaveBeenCalledTimes(4)
                    expect(loadCodes).toHaveBeenCalledWith("gender")
                    expect(loadCodes).toHaveBeenCalledWith("countries")
                    expect(loadCodes).toHaveBeenCalledWith("maritalStatus")
                    done()
                })
                .catch(ex => {
                    console.log(ex)
                    expect(ex).toBeUndefined()
                })
        })

        //
    })

    describe("codeToString", () => {
        it("should find a match", () => {
            const out = codeToString("key", [{ category: "test", id: "key", title: "title" }])
            expect(out).toBe("test::key|title|")
        })

        it("should detect SNOMED codes", () => {
            const out = codeToString("SNOMED-key", [{ category: "test", id: "SNOMED-key", title: "title" }])
            expect(out).toBe("SNOMED::key|title|")
        })

        it("should detect LOCAL codes", () => {
            const out = codeToString("CODED-key", [{ category: "test", id: "CODED-key", title: "title" }])
            expect(out).toBe("local::key|title|")
        })
    })
})
