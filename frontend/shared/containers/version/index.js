import React from "react"
import "./style.css"

const defaultVersion = "dev"

export default () => (<div className="appversion" >{process.env.REACT_APP_VERSION || defaultVersion}</div>)