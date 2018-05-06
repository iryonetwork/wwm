import React from "react"
// import { expect } from "chai"
import { shallow, mount } from "enzyme"
import sinon from "sinon"
import { load as loadClinic } from './clinics'
import { load as loadCodes } from 'shared/modules/codes'
// import locations from './locations'
import { load as loadUser } from './users'
jest.mock('./users')
jest.mock('./clinics')
jest.mock('shared/modules/codes')

import {composePatientData, codeToString} from './ehr'

const fullFormData = {
    documents: [
        {
            type: 'syrian-id',
            number: '1234567'
        },
        {
            type: 'un-id',
            number: '987654'
        }
    ],
    firstName: 'Dominik',
    lastName: 'Znidar',
    dateOfBirth: '1983-05-21',
    gender: 'CODED-at0310',
    maritalStatus: 'SNOMED-125681006',
    numberOfKids: '1',
    nationality: 'SI',
    countryOfOrigin: 'SI',
    education: 'gimnasium,',
    profession: 'developer',
    country: 'SI',
    camp: '19',
    tent: '83',
    clinic: 'ZD Kranj',
    phone: '040123456',
    email: 'email@test.com',
    whatsapp: '987987654',
    dateOfLeaving: '2016-01-01',
    transitCountries: 'Slovenia',
    dateOfArrival: '2016-12-31',
    people_in_family: '1',
    people_living_together: '2',
    familyMembers: [
        {
            firstName: 'Test',
            lastName: 'Test',
            dateOfBirth: '1993-05-04',
            relation: 'child',
            livingTogether: 'true',
            documentType: 'syrian_id',
            documentNumber: '1234567'
        }
    ],
    allergies: [
        {
            allergy: 'allergy 1',
            comment: 'allergy comment 1',
            critical: 'false'
        },
        {
            allergy: 'allergy 2',
            comment: 'allergy comment 2',
            critical: 'true'
        }
    ],
    immunizations: [
        {
            immunization: 'Immunization 1',
            date: '2014-02-01'
        },
        {
            immunization: 'Immunuzation 2',
            date: '2004-03-02'
        }
    ],
    chronicDiseases: [
        {
            disease: 'Chronic 1',
            date: '2004-03-12',
            medication: 'Medication 1'
        },
        {
            disease: 'Chronic 2 ',
            date: '2005-04-23',
            medication: 'Medication 2'
        }
    ],
    injuries: [
        {
            injury: 'Injury 1',
            date: '2007-06-05',
            medication: 'Aids 1'
        },
        {
            injury: 'Injury 2',
            date: '2009-08-07',
            medication: 'Aids 2'
        }
    ],
    surgeries: [
        {
            injury: 'Surgery 1',
            date: '2004-03-12',
            medication: 'Comment 1'
        },
        {
            injury: 'Surgery 2 ',
            date: '2004-03-12',
            medication: 'Comment 2'
        }
    ],
    medications: [
        {
            medication: 'Medication 1',
            comment: 'Comment 1'
        },
        {
            medication: 'Medication 2',
            comment: 'Comment 2'
        }
    ],
    habits_smoking: 'true',
    habits_smoking_comment: '10 boxes a day',
    habits_drugs: 'true',
    habits_drugs_comment: 'All of them',
    conditions_basic_hygiene: 'true',
    conditions_heating: 'true',
    conditions_good_appetite: 'true',
    conditions_food_supply: 'true',
    conditions_clean_water: 'true',
    conditions_electricity: 'true'
}

// const waitlistItem = {
//     added: '2018-05-02T15:30:13.435Z',
// }

if (!global.window.localStorage) {
    global.window.localStorage = {
        getItem() { return '{}'; },
        setItem() {}
    };
}

beforeEach(() => {
    loadClinic.mockClear()
    loadUser.mockClear()
});

describe('ehr', () => {

    describe('composePatientData', () => {
        it("should return object with person and info keys", () => {
            loadClinic.mockResolvedValue({
                id: 'e4ebb41b-7c62-4db7-9e1c-f47058b96dd0',
                name: 'CLINIC1',
            })
            loadUser.mockResolvedValue({
                id: 'c12574e4-acd4-4266-9d53-b614c8a942bc',
                personalData: {
                    firstName: 'Doctor',
                    lastName: 'X',
                },
            })
            loadCodes.mockResolvedValue(() => [
                {category: "category", id: "SH", locale: "en", title: "Saint Helena"},
                {category: "category", id: "SI", locale: "en", title: "Slovenia"},
            ])
            // loadCodes.mockResolvedValue([
            //     { "category": "gender", "id": "CODED-at0310", "locale": "en", "title": "Male" },
            //     { "category": "gender", "id": "CODED-at0311", "locale": "en", "title": "Female" },
            // ])
            // loadCodes.mockResolvedValue([
            //     {category: "maritalStatus", id: "SNOMED-125681006", locale: "en", title: "Single"},
            //     {category: "maritalStatus", id: "SNOMED-20295000", locale: "en", title: "Divorced"},
            // ])
            // loadCodes.mockResolvedValue([
            //     {category: "countries", id: "SH", locale: "en", title: "Saint Helena"},
            //     {category: "countries", id: "SI", locale: "en", title: "Slovenia"},
            // ])

            const getState = () => ({
                locations: {
                    cache: {},
                },
                authentication: {
                    tokenString: 'TOKEN',
                }
            })
            const dispatch = (fn) => { return typeof fn === 'function' ? fn(dispatch, getState) : fn}

            dispatch(composePatientData(fullFormData))
                .then(out => {
                        expect(Object.keys(out)).toEqual(['person', 'info'])
                        expect(loadUser).toHaveBeenLastCalledWith('me')
                        expect(loadClinic).toHaveBeenLastCalledWith('e4ebb41b-7c62-4db7-9e1c-f47058b96dd0')
                        expect(loadCodes).toHaveBeenCalledTimes(4)
                        expect(loadCodes).toHaveBeenCalledWith('gender')
                        expect(loadCodes).toHaveBeenCalledWith('countries')
                        expect(loadCodes).toHaveBeenCalledWith('maritalStatus')
                    })
                .catch(ex => {
                    expect(ex).toBeUndefined()
                })
        })

        //

    })

    describe('codeToString', () => {
        it('should find a match', () => {
            const out = codeToString('key', [
                {category: 'test', id: 'key', title: 'title'},
            ])
            expect(out).toBe('test::key|title|')
        })

        it('should detect SNOMED codes', () => {
            const out = codeToString('SNOMED-key', [
                {category: 'test', id: 'SNOMED-key', title: 'title'},
            ])
            expect(out).toBe('SNOMED::key|title|')
        })

        it('should detect LOCAL codes', () => {
            const out = codeToString('CODED-key', [
                {category: 'test', id: 'CODED-key', title: 'title'},
            ])
            expect(out).toBe('local::key|title|')
        })
    })
})
