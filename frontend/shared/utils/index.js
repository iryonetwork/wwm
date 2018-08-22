import React from "react"
import { Route } from "react-router-dom"

const renderMergedProps = (component, ...rest) => {
    const finalProps = Object.assign({}, ...rest)
    return React.createElement(component, finalProps)
}

const PropRoute = ({ component, ...rest }) => {
    return (
        <Route
            {...rest}
            render={routeProps => {
                return renderMergedProps(component, routeProps, rest)
            }}
        />
    )
}

const joinPaths = (...paths) => {
    for (let i = 0; i < paths.length - 1; i++) {
        let current = paths[i]
        let next = paths[i + 1]

        if (current[current.length - 1] === "/" && next[0] === "/") {
            paths[i + 1] = next.substr(1)
        }
        if (current[current.length - 1] !== "/" && next[0] !== "/") {
            paths[i] = current + "/"
        }
    }

    return paths.join("")
}

const round = (number, precision) => {
    var shift = function(number, precision) {
        var numArray = ("" + number).split("e")
        return +(numArray[0] + "e" + (numArray[1] ? +numArray[1] + precision : precision))
    }
    return shift(Math.round(shift(number, +precision)), -precision)
}

const escapeRegex = text => {
    return text.replace(/[-[\]{}()*+?.,\\^$|#\s]/g, "\\$&")
}

const confirmationDialog = (msg, confirmedCallback, rejectedCallback) => {
    if (window.confirm(msg)) {
        confirmedCallback && confirmedCallback()
    } else {
        rejectedCallback && rejectedCallback()
    }
}

export { PropRoute, joinPaths, round, escapeRegex, confirmationDialog }
