import store from "../store"

const baseURL = "https://iryo.cloud"

export default (url, method, body) => {
    return fetch(baseURL + url, {
        method: method,
        headers: {
            Authorization: store.getState().authentication.tokenString
        }
    })
        .catch(error => {
            throw error
        })
        .then(response => Promise.all([response.ok, response.json()]))
        .then(([responseOk, body]) => {
            if (responseOk) {
                return body
            } else {
                throw new Error(body)
            }
        })
}
