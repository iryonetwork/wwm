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
    if (precision === undefined) {
        return number
    }

    var shift = function(number, precision) {
        var numArray = ("" + number).split("e")
        return +(numArray[0] + "e" + (numArray[1] ? +numArray[1] + precision : precision))
    }
    return shift(Math.round(shift(number, +precision)), -precision)
}

const isStringNumber = s => {
    s = s.replace(",", ".")
    return !isNaN(parseFloat(s)) && isFinite(s)
}

const getNumberFromString = s => {
    s = s.replace(",", ".")
    return parseFloat(s)
}

const getPrecision = value => {
    if (!isFinite(value)) return 0
    var e = 1,
        p = 0
    while (Math.round(value * e) / e !== value) {
        e *= 10
        p++
    }
    return p
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

const toggleBodyScroll = () => {
    document.body.style.overflow = document.body.style.overflow === "hidden" ? "scroll" : "hidden"
}

export { PropRoute, joinPaths, round, isStringNumber, getNumberFromString, getPrecision, escapeRegex, confirmationDialog, toggleBodyScroll }
