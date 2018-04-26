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

export { joinPaths }
