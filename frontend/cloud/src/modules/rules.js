import produce from "immer"
import keyBy from "lodash/keyBy"
import forEach from "lodash/forEach"

import api from "./api"
import { open, close, COLOR_DANGER, COLOR_SUCCESS } from "./alert"

const LOAD_RULES = "roles/LOAD_RULES"
const LOAD_RULES_SUCCESS = "roles/LOAD_RULES_SUCCESS"
const LOAD_RULES_FAIL = "roles/LOAD_RULES_FAIL"

const SAVE_RULE = "roles/SAVE_RULE"
const SAVE_RULE_SUCCESS = "roles/SAVE_RULE_SUCCESS"
const SAVE_RULE_FAIL = "roles/SAVE_RULE_FAIL"

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
                forEach(action.rules, rule => {
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
                if (!draft.subjects[action.rule.subject]) {
                    draft.subjects[action.rule.subject] = []
                }
                draft.subjects[action.rule.subject].push(action.rule.id)
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
                    rule: response
                })
                dispatch(open("Saved!", "", COLOR_SUCCESS, 5))
            })
            .catch(error => {
                dispatch({
                    type: SAVE_RULE_FAIL
                })
                console.log(error)
                dispatch(open(error.message, error.code, COLOR_DANGER))
            })
    }
}
