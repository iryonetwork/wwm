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

export { PropRoute, joinPaths }
