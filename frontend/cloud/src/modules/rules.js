import produce from "immer"
import keyBy from "lodash/keyBy"
import _ from "lodash"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "./alert"

const LOAD_RULES = "rules/LOAD_RULES"
const LOAD_RULES_SUCCESS = "rules/LOAD_RULES_SUCCESS"
const LOAD_RULES_FAIL = "rules/LOAD_RULES_FAIL"

const SAVE_RULE = "rules/SAVE_RULE"
const SAVE_RULE_SUCCESS = "rules/SAVE_RULE_SUCCESS"
const SAVE_RULE_FAIL = "rules/SAVE_RULE_FAIL"

const DELETE_RULE = "rules/DELETE_RULE"
const DELETE_RULE_SUCCESS = "rules/DELETE_RULE_SUCCESS"
const DELETE_RULE_FAIL = "rules/DELETE_RULE_FAIL"

const initialState = {
    loading: true
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case LOAD_RULES:
                draft.loading = true
                break
            case LOAD_RULES_SUCCESS:
                draft.loading = false
                draft.rules = keyBy(action.rules, "id")
                draft.subjects = {}
                _.forEach(action.rules, rule => {
                    if (!draft.subjects[rule.subject]) {
                        draft.subjects[rule.subject] = []
                    }
                    draft.subjects[rule.subject].push(rule.id)
                })
                break
            case LOAD_RULES_FAIL:
                draft.loading = false
                break

            case SAVE_RULE:
                draft.loading = true
                break
            case SAVE_RULE_FAIL:
                draft.loading = false
                break
            case SAVE_RULE_SUCCESS:
                draft.loading = false
                draft.rules[action.rule.id] = action.rule
                draft.rules[action.rule.id].edit = false
                draft.rules[action.rule.id].saving = false
                draft.rules[action.rule.id].index = action.index
                if (!_.get(state, `subjects['${action.rule.subject}'][${action.index}]`)) {
                    draft.subjects[action.rule.subject] = state.subjects[action.rule.subject] || []
                    draft.subjects[action.rule.subject].push(action.rule.id)
                }
                break

            case DELETE_RULE:
                draft.loading = true
                break
            case DELETE_RULE_SUCCESS:
                draft.loading = false
                draft.subjects[state.rules[action.id].subject] = _.without(state.subjects[state.rules[action.id].subject], action.id)
                draft.rules = state.rules
                delete draft.rules[action.id]
                break
            case DELETE_RULE_FAIL:
                draft.loading = false
                break
            default:
        }
    })
}

export const loadRules = () => {
    return dispatch => {
        dispatch({
            type: LOAD_RULES
        })
        dispatch(close())

        return api(`/auth/rules`, "GET")
            .then(response => {
                dispatch({
                    type: LOAD_RULES_SUCCESS,
                    rules: response
                })
            })
            .catch(error => {
                dispatch({
                    type: LOAD_RULES_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const saveRule = rule => {
    rule.deny = rule.deny === "true" || rule.deny === true
    rule.action |= 0

    return dispatch => {
        dispatch({
            type: SAVE_RULE
        })
        dispatch(close())

        let url = "/auth/rules"
        let method = "POST"
        if (rule.id) {
            url += "/" + rule.id
            method = "PUT"
        }

        return api(url, method, rule)
            .then(response => {
                if (rule.id) {
                    response = rule
                }
                dispatch({
                    type: SAVE_RULE_SUCCESS,
                    index: rule.index,
                    rule: response
                })
                dispatch(open("Saved!", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: SAVE_RULE_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}

export const deleteRule = ruleID => {
    return dispatch => {
        dispatch({
            type: DELETE_RULE
        })
        dispatch(close())

        return api(`/auth/rules/${ruleID}`, "DELETE")
            .then(response => {
                dispatch({
                    type: DELETE_RULE_SUCCESS,
                    id: ruleID
                })
                dispatch(open("Deleted!", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: DELETE_RULE_FAIL
                })
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
