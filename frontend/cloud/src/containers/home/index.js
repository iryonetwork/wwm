import React from "react"
import { connect } from "react-redux"

import UserDetail from "../users/detail"

const Home = props => (
    <div>
        <UserDetail userID={props.userID} home={true} />
    </div>
)

const mapStateToProps = state => ({
    userID: state.authentication.token.sub
})

export default connect(mapStateToProps)(Home)
