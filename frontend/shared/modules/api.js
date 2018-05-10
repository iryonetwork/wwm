import { read, BASE_URL } from "./config"
import { getToken } from "./authentication"

function APIError(error) {
    this.message = error.message
    this.code = error.code
    this.name = "API Error"
}

export default (url, method, body) => {

    let fullUrl = read(BASE_URL) + url
    return fetch(fullUrl, {
        method: method,
        headers: {
            Authorization: getToken(),
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
