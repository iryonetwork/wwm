import _ from "lodash"

import api from "shared/modules/api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "shared/modules/alert"

const LOAD_RULES = "rules/LOAD_RULES"
const LOAD_RULES_SUCCESS = "rules/LOAD_RULES_SUCCESS"
const LOAD_RULES_FAIL = "rules/LOAD_RULES_FAIL"

const SAVE_RULE_SUCCESS = "rules/SAVE_RULE_SUCCESS"
const SAVE_RULE_FAIL = "rules/SAVE_RULE_FAIL"

const DELETE_RULE_SUCCESS = "rules/DELETE_RULE_SUCCESS"

const initialState = {
    loading: false,
    allLoaded: false,
    forbidden: false
}

export default (state = initialState, action) => {
    let subjects, rules
    switch (action.type) {
        case LOAD_RULES:
            return {
                ...state,
                loading: true
            }

        case LOAD_RULES_SUCCESS:
            subjects = {}
            _.forEach(action.rules, rule => {
                if (!subjects[rule.subject]) {
                    subjects[rule.subject] = []
                }
                subjects[rule.subject].push(rule.id)
            })

            return {
                loading: false,
                allLoaded: true,
                rules: _.keyBy(action.rules, "id"),
                subjects
            }

        case LOAD_RULES_FAIL:
            let forbidden = false
            if (action.code === 403) {
                forbidden = true
            }
            return {
                ...state,
                forbidden,
                loading: false
            }

        case SAVE_RULE_SUCCESS:
            subjects = { ...state.subjects }
            rules = { ...state.rules }
            rules[action.rule.id] = action.rule
            rules[action.rule.id].edit = false
            rules[action.rule.id].saving = false
            rules[action.rule.id].index = action.index
            if (subjects[action.rule.subject]) {
                if (_.indexOf(subjects[action.rule.subject], action.rule.id) === -1) {
                    subjects[action.rule.subject].push(action.rule.id)
                }
            } else {
                subjects[action.rule.subject] = [action.rule.id]
            }
            return {
                loading: false,
                rules,
                subjects
            }

        case SAVE_RULE_FAIL:
            rules = { ...state.rules }
            if (action.rule.id) {
                rules[action.rule.id] = action.rule
                rules[action.rule.id].saving = false
            } else {
                action.rule.saving = false
            }

            return {
                ...state,
                rules
            }

        case DELETE_RULE_SUCCESS:
            subjects = { ...state.subjects }
            subjects[state.rules[action.id].subject] = _.without(subjects[state.rules[action.id].subject], action.id)

            return {
                loading: false,
                rules: _.pickBy(state.rules, rule => rule.id !== action.id),
                subjects
            }

        default:
            return state
    }
}

export const loadRules = () => {
    return dispatch => {
        dispatch({
            type: LOAD_RULES
        })
        dispatch(close())

        return dispatch(api(`/auth/rules`, "GET"))
            .then(response => {
                dispatch({
                    type: LOAD_RULES_SUCCESS,
                    rules: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_RULES_FAIL,
                    code: error.code
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveRule = rule => {
    rule.deny = rule.deny === "true" || rule.deny === true
    rule.action |= 0

    return dispatch => {
        dispatch(close())

        let url = "/auth/rules"
        let method = "POST"
        if (rule.id) {
            url += "/" + rule.id
            method = "PUT"
        }

        return dispatch(api(url, method, rule))
            .then(response => {
                if (rule.id) {
                    response = rule
                }
                dispatch({
                    type: SAVE_RULE_SUCCESS,
                    index: rule.index,
                    rule: response
                })
                setTimeout(() => dispatch(open("Saved ACL rule", "", COLOR_SUCCESS, 5)), 100)

                return response
            })
            .catch(error => {
                dispatch({
                    type: SAVE_RULE_FAIL,
                    rule: rule
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteRule = ruleID => {
    return dispatch => {
        dispatch(close())

        return dispatch(api(`/auth/rules/${ruleID}`, "DELETE"))
            .then(response => {
                dispatch({
                    type: DELETE_RULE_SUCCESS,
                    id: ruleID
                })
                setTimeout(() => dispatch(open("Deleted ACL rule", "", COLOR_SUCCESS, 5)), 100)
            })
            .catch(error => {
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
