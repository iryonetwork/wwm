import store from "../store"

const baseURL = "https://iryo.cloud"

function APIError(error) {
    this.message = error.message
    this.code = error.code
    this.name = "API Error"
}

export default (url, method, body) => {
    return fetch(baseURL + url, {
        method: method,
        headers: {
            Authorization: store.getState().authentication.tokenString,
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    })
        .catch(error => {
            throw new Error("Failed to connect to server")
        })
        .then(response => {
            if (response.status === 204) {
                return Promise.all([response.ok, {}])
            }
            return Promise.all([response.ok, response.json()])
        })
        .then(([responseOk, body]) => {
            if (responseOk) {
                return body
            } else {
                throw new APIError(body)
            }
        })
}
