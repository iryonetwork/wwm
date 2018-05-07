import React from "react"
import { expect } from "chai"
import { shallow, mount } from "enzyme"
import sinon from "sinon"

import { combineReducers, createStore, applyMiddleware } from "redux"
import { Provider } from "react-redux"
import { ConnectedRouter } from "react-router-redux"
import createHistory from "history/createBrowserHistory"
// import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk"
import rootReducer from "../../../modules"

import Enzyme from "enzyme"
import Adapter from "enzyme-adapter-react-16"

Enzyme.configure({ adapter: new Adapter() })

import NewPatientForm from "./index"
import Step1 from "./step1"
import Step2 from "./step2"
import Step3 from "./step3"

if (!global.window.localStorage) {
    global.window.localStorage = {
        getItem() {
            return "{}"
        },
        setItem() {}
    }
}

const initState = {
    codes: {
        fetching: [],
        cache: {
            countries: [
                { category: "category", id: "SH", locale: "en", title: "Saint Helena" },
                { category: "category", id: "SI", locale: "en", title: "Slovenia" }
            ],
            gender: [
                { category: "category", id: "SH", locale: "en", title: "Saint Helena" },
                { category: "category", id: "SI", locale: "en", title: "Slovenia" }
            ],
            maritalStatus: [
                { category: "category", id: "SH", locale: "en", title: "Saint Helena" },
                { category: "category", id: "SI", locale: "en", title: "Slovenia" }
            ],
            documentTypes: [
                { category: "category", id: "SH", locale: "en", title: "Saint Helena" },
                { category: "category", id: "SI", locale: "en", title: "Slovenia" }
            ]
        }
    },
    patient: {
        newData: {
            documents: [{}]
        }
    }
}

describe("<NewPatientForm />", () => {
    it("renders three header links", () => {
        const store = createStore(rootReducer, initState, applyMiddleware(thunk))
        const history = createHistory()

        const wrapper = mount(
            <Provider store={store}>
                <ConnectedRouter history={history}>
                    <NewPatientForm />
                </ConnectedRouter>
            </Provider>
        )

        expect(wrapper.find("ol li")).to.have.length(3)
    })

    it("renders first page", () => {
        const store = createStore(rootReducer, initState, applyMiddleware(thunk))
        const history = createHistory()

        const wrapper = mount(
            <Provider store={store}>
                <ConnectedRouter history={history}>
                    <NewPatientForm />
                </ConnectedRouter>
            </Provider>
        )

        expect(wrapper.find(Step1)).to.have.length(1)
    })

    it("renders required errors on empty submit", () => {
        const store = createStore(rootReducer, initState, applyMiddleware(thunk))
        const history = createHistory()

        const wrapper = mount(
            <Provider store={store}>
                <ConnectedRouter history={history}>
                    <NewPatientForm />
                </ConnectedRouter>
            </Provider>
        )
        wrapper.find("form").simulate("submit")
        expect(wrapper.find(".is-invalid")).to.have.lengthOf.above(1)
    })

    it("goes to page 2 when entered required values", () => {
        const store = createStore(rootReducer, initState, applyMiddleware(thunk))
        const history = createHistory()

        const wrapper = mount(
            <Provider store={store}>
                <ConnectedRouter history={history}>
                    <NewPatientForm />
                </ConnectedRouter>
            </Provider>
        )

        wrapper.find("form").simulate("submit")
        expect(wrapper.find(".is-invalid")).to.have.lengthOf.above(1)

        wrapper.find(`input[name="firstName"]`).simulate("change", { target: { value: "First" } })
        wrapper.find(`input[name="lastName"]`).simulate("change", { target: { value: "Last" } })
        wrapper.find(`input[name="dateOfBirth"]`).simulate("change", { target: { value: "02/02/2002" } })
        wrapper.find(`select[name="gender"]`).simulate("change", { target: { value: "m" } })
        wrapper.find(`select[name="maritalStatus"]`).simulate("change", { target: { value: "maried" } })
        wrapper.find(`select[name="numberOfKids"]`).simulate("change", { target: { value: "3" } })
        wrapper.find(`select[name="nationality"]`).simulate("change", { target: { value: "syrian" } })
        wrapper.find(`select[name="countryOfOrigin"]`).simulate("change", { target: { value: "syria" } })

        wrapper.find(`select[name="documents[0].type"]`).simulate("change", { target: { value: "un_id" } })
        wrapper.find(`input[name="documents[0].number"]`).simulate("change", { target: { value: "UN12345" } })
        wrapper.find(".addDocument").simulate("click")
        wrapper.find(`select[name="documents[1].type"]`).simulate("change", { target: { value: "syrian_id" } })
        wrapper.find(`input[name="documents[1].number"]`).simulate("change", { target: { value: "54321" } })

        wrapper.find(`select[name="country"]`).simulate("change", { target: { value: "SY" } })
        wrapper.find(`input[name="camp"]`).simulate("change", { target: { value: "2" } })
        wrapper.find(`input[name="tent"]`).simulate("change", { target: { value: "7" } })
        expect(wrapper.find(".is-invalid")).to.have.length(0)

        wrapper.find("form").simulate("submit")

        expect(wrapper.find(Step1)).to.have.length(0)
        expect(wrapper.find(Step2)).to.have.length(1)
        expect(wrapper.find(Step3)).to.have.length(0)

        wrapper.find("form").simulate("submit")

        expect(wrapper.find(Step1)).to.have.length(0)
        expect(wrapper.find(Step2)).to.have.length(0)
        expect(wrapper.find(Step3)).to.have.length(1)
    })
})
