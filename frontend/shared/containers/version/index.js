import React from "react"
import "./style.css"

const defaultVersion = "dev"

export default () => (<div class="appversion" >{process.env.REACT_APP_VERSION || defaultVersion}</div>)