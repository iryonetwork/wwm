import React from "react"
import { connect } from "react-redux"

import User from "../users/detail"

const Home = props => (
    <div>
        <User userID={props.userID} home={true} />
    </div>
)

const mapStateToProps = state => ({
    userID: state.authentication.token.sub
})

export default connect(mapStateToProps)(Home)
