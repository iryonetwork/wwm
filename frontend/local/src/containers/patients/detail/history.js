import React from "react"

//import "./style.css"

export default () => (
    <div className="history">
        <header>
            <h1>Medical History</h1>
        </header>

        <div className="section">
            <div className="name">Blood type</div>
            <div className="values">A+</div>
        </div>

        <div className="section">
            <div className="name">Allergies</div>
            <div className="values">
                <dl>
                    <dt>Peanuts</dt>
                    <dd>high risk</dd>
                    <dt>Pollen</dt>
                    <dd />
                </dl>
            </div>
        </div>
    </div>
)
