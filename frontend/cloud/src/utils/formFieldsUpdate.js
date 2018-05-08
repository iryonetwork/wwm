export const processStateOnChange = (state, e) => {
    const target = e.target
    const value = target.type === "checkbox" ? target.checked : target.value

    let id
    let toAssign
    let splitID = target.id.split(".")

    id = splitID[0]

    switch (splitID.length) {
        case 3:
            toAssign = state[id]
            toAssign[splitID[1]][splitID[2]] = value
            break
        case 2:
            toAssign = state[id]
            toAssign[splitID[1]] = value
            break
        default:
            id = target.id
            toAssign = value
    }
    state[id] = toAssign

    if (target.required && !value) {
        state.validationErrors[target.id] = "Required"
        return state
    }

    if (state.validationErrors[target.id]) {
        delete state.validationErrors[target.id]
    }

    return state
}

export const processStateOnBlur = (state, e) => {
    // trims input on blur
    const target = e.target
    const value = target.type === "checkbox" ? target.checked : target.value

    let id
    let toAssign
    let splitID = target.id.split(".")

    id = splitID[0]

    switch (splitID.length) {
        case 3:
            toAssign = state[id]
            toAssign[splitID[1]][splitID[2]] = value.trim()
            break
        case 2:
            toAssign = state[id]
            toAssign[splitID[1]] = value.trim()
            break
        default:
            toAssign = value.trim()
    }

    state[id] = toAssign

    if (target.required && !value) {
        state.validationErrors[target.id] = "Required"
        return state
    }

    if (state.validationErrors[target.id]) {
        delete state.validationErrors[target.id]
    }

    return state
}
