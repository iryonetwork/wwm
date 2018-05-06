import produce from "immer"

export const SHOW_ALERT = "alert/SHOW_ALERT"
export const HIDE_ALERT = "alert/HIDE_ALERT"
export const COLOR_DANGER = "danger"
export const COLOR_PRIMARY = "primary"
export const COLOR_SUCCESS = "success"

const initialState = {
    open: false
}

export default (state = initialState, action) => {
    return produce(state, draft => {
        switch (action.type) {
            case SHOW_ALERT:
                draft.open = true
                draft.message = action.message
                draft.code = action.code
                draft.color = action.color ? action.color : COLOR_PRIMARY
                break
            case HIDE_ALERT:
                draft.message = ''
                break
            default:
        }
    })
}

export const open = (message, code, color, closeIn) => {
    return dispatch => {
        dispatch({
            type: SHOW_ALERT,
            message,
            code,
            color
        })

        if (closeIn) {
            setTimeout(() => {
                dispatch({ type: HIDE_ALERT })
            }, closeIn * 1000)
        }
    }
}

export const close = () => {
    return dispatch => {
        dispatch({ type: HIDE_ALERT })
    }
}
