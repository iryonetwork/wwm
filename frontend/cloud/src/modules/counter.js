import produce from "immer"

export const INCREMENT_REQUESTED = "counter/INCREMENT_REQUESTED"
export const INCREMENT = "counter/INCREMENT"
export const DECREMENT_REQUESTED = "counter/DECREMENT_REQUESTED"
export const DECREMENT = "counter/DECREMENT"

const initialState = {
    count: 0,
    isIncrementing: false,
    isDecrementing: false
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case INCREMENT_REQUESTED:
                draft.isIncrementing = true
                break
            case INCREMENT:
                draft.count = state.count + 1
                draft.isIncrementing = !state.isIncrementing
                break
            case DECREMENT_REQUESTED:
                draft.isDecrementing = true
                break
            case DECREMENT:
                draft.count = state.count - 1
                draft.isDecrementing = !state.isDecrementing
                break
            default:
        }
    })
}

export const increment = () => {
    return dispatch => {
        dispatch({
            type: INCREMENT_REQUESTED
        })

        dispatch({
            type: INCREMENT
        })
    }
}

export const incrementAsync = () => {
    return dispatch => {
        dispatch({
            type: INCREMENT_REQUESTED
        })

        return setTimeout(() => {
            dispatch({
                type: INCREMENT
            })
        }, 3000)
    }
}

export const decrement = () => {
    return dispatch => {
        dispatch({
            type: DECREMENT_REQUESTED
        })

        dispatch({
            type: DECREMENT
        })
    }
}

export const decrementAsync = () => {
    return dispatch => {
        dispatch({
            type: DECREMENT_REQUESTED
        })

        return setTimeout(() => {
            dispatch({
                type: DECREMENT
            })
        }, 3000)
    }
}
