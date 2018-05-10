import { get, has, set } from "lodash"
import { read, CLINIC_ID, LOCATION_ID } from "shared/modules/config"
import { load as loadClinic } from "./clinics"
import { load as loadLocation } from "./locations"
import { load as loadUser } from "./users"
import personSpec from "./ehr/person"
import infoSpec from "./ehr/info"

// Converts form data into two separate documents
export const composePatientData = formData => dispatch => {
    return dispatch(buildContextForPatientData(formData)).then(context => {
        return Promise.all([dispatch(buildPersonData(context, formData)), dispatch(buildInfoData(context, formData))]).then(([person, info]) => ({
            person,
            info
        }))
    })
}

export const extractPatientData = (person, info) => dispatch => {
    return Promise.all([dispatch(extractPersonData(person)), dispatch(extractInfoData(info))]).then(([person, info]) => Object.assign(person, info))
}

const buildContextForPatientData = formData => dispatch => {
    return Promise.all([
        dispatch(loadClinic(dispatch(read(CLINIC_ID)))),
        dispatch(loadLocation(dispatch(read(LOCATION_ID)))),
        dispatch(loadUser("me")) // doctor
    ]).then(([clinic, location, doctor]) => {
        return {
            // facility details
            "/context/health_care_facility|name": clinic.name,
            "/context/health_care_facility|identifier": clinic.id,
            "/territory": location.country,
            "/language": "en",

            // time info
            "/context/start_time": new Date().toJSON(),
            "/context/end_time": new Date().toJSON(),

            // // participants
            // // add doctor
            // "/composer|identifier": doctor.id,
            // "/composer|name": `${doctor.personalData.firstName} ${doctor.personalData.lastName}`,

            "/category": "openehr::431|persistent|"
        }
    })
}

const buildPersonData = (ctx, formData) => dispatch => {
    return dispatch(personSpec).then(spec => {
        return specToDocument(spec, Object.assign({}, ctx), formData, "")
    })
}

const buildInfoData = (ctx, formData) => dispatch => {
    return dispatch(infoSpec).then(spec => {
        return specToDocument(spec, Object.assign({}, ctx), formData, "")
    })
}

const extractPersonData = doc => dispatch => {
    return dispatch(personSpec).then(spec => {
        return specToObject(spec, {}, doc, "")
    })
}

const extractInfoData = doc => dispatch => {
    return dispatch(infoSpec).then(spec => {
        return specToObject(spec, {}, doc, "")
    })
}

const specToDocument = (specs, data, formData, ehrPrefix) => {
    const fns = specs.reduce((acc, spec) => {
        switch (spec.type) {
            case "value":
                acc.push(assignValue(ehrPrefix + spec.ehrPath, spec.formPath))
                break

            case "integer":
                acc.push(assignInteger(ehrPrefix + spec.ehrPath, spec.formPath))
                break

            case "boolean":
                acc.push(assignBoolean(ehrPrefix + spec.ehrPath, spec.formPath))
                break

            case "fixedValue":
                acc.push(assignFixedValue(ehrPrefix + spec.ehrPath, spec.value))
                break

            case "quantity":
                acc.push(assignQuantity(ehrPrefix + spec.ehrPath, spec.formPath, spec.unit))
                break

            case "code":
                acc.push(assignCode(spec.codes, ehrPrefix + spec.ehrPath, spec.formPath))
                break

            case "array":
                ;(get(formData, spec.formPath, []) || []).forEach((arrEl, i) => {
                    data = specToDocument(spec.items, data, arrEl, `${spec.ehrPath}:${i}`)
                })
                break

            default:
                throw new Error(`Invalid type "${spec.type}"`)
        }
        return acc
    }, [])

    return compose(data, formData, fns)
}

const codeRe = /^(.+)::(.+)\|(.+)\|$/

const specToObject = (specs, data, doc, ehrPrefix) => {
    specs.forEach((spec, i) => {
        // skip when ehr key is not present in the document (but not for arrays)
        if (!(ehrPrefix + spec.ehrPath in doc) && spec.type !== "array") {
            return
        }

        const value = doc[ehrPrefix + spec.ehrPath]
        let newValue = undefined
        switch (spec.type) {
            case "fixedValue":
                // noop
                return

            case "value":
            case "integer":
            case "boolean":
                newValue = value
                break

            case "quantity":
                const re = new RegExp(`\\,${spec.unit}$`)

                // skip if value is malformed
                if (!re.test(value)) {
                    // noop
                    return
                }

                newValue = value.replace(re, "")
                break

            case "code":
                // skip if value is malformed
                if (!codeRe.test(value)) {
                    // noop
                    return
                }

                // extract values
                const values = value.match(codeRe)
                let key = values[2]

                // handle custom category names
                if (values[1] === "SNOMED") {
                    key = `SNOMED-${key}`
                } else if (values[1] === "local") {
                    key = `CODED-${key}`
                }

                newValue = key
                break

            case "array":
                let arrVal = []
                var j = 0
                while (true) {
                    const out = specToObject(spec.items, {}, doc, `${ehrPrefix}${spec.ehrPath}:${j}`)
                    if (Object.keys(out).length === 0) {
                        break
                    }

                    arrVal.push(out)
                    j++
                }
                newValue = arrVal
                break

            default:
                throw new Error(`Invalid type "${spec.type}"`)
        }

        if (newValue !== undefined) {
            data = set(data, spec.formPath, newValue)
        }
    })

    return data
}

const compose = (data, formData, fns) => {
    return fns.reduce((acc, fn) => fn(acc, formData), data)
}

const assignValue = (ehrPath, formPath) => (data, formData) => {
    if (!has(formData, formPath)) {
        return data
    }
    return Object.assign(data, { [ehrPath]: get(formData, formPath) })
}

const assignInteger = (ehrPath, formPath) => (data, formData) => {
    if (!has(formData, formPath)) {
        return data
    }
    return Object.assign(data, { [ehrPath]: parseInt(get(formData, formPath), 10) })
}

const assignQuantity = (ehrPath, formPath, unit) => (data, formData) => {
    if (!has(formData, formPath)) {
        return data
    }
    return Object.assign(data, { [ehrPath]: `${get(formData, formPath)},${unit}` })
}

const assignBoolean = (ehrPath, formPath) => (data, formData) => {
    if (!has(formData, formPath)) {
        return data
    }
    return Object.assign(data, { [ehrPath]: get(formData, formPath) ? "true" : "false" })
}

const assignFixedValue = (ehrPath, value) => (data, formData) => {
    return Object.assign(data, { [ehrPath]: value })
}

const assignCode = (codes, ehrPath, formPath) => (data, formData) => {
    if (!has(formData, formPath)) {
        return data
    }
    return Object.assign(data, { [ehrPath]: codeToString(get(formData, formPath), codes) })
}

export const codeToString = (key, codes) => {
    for (let i = 0; i < codes.length; i++) {
        const el = codes[i]
        if (el.id === key) {
            let category = el.category
            let id = el.id
            if (id.indexOf("SNOMED-") === 0) {
                category = "SNOMED"
                id = id.substring(7)
            } else if (id.indexOf("CODED-") === 0) {
                category = "local"
                id = id.substring(6)
            }

            return `${category}::${id}|${el.title}|`
        }
    }
}
