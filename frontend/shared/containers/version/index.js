import React from "react"
import "./style.css"

const defaultVersion = "dev"

export default () => (<div id="version" class="version" >{process.env.REACT_APP_VERSION || defaultVersion}</div>)