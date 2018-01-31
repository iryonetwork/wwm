import React from "react"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import { withRouter } from "react-router-dom"

import { loadUser } from "../../modules/users"

class UserDetail extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            email: "",
            password: "",
            password2: ""
        }

        this.updateEmail = this.updateEmail.bind(this)
        this.updatePassword = this.updatePassword.bind(this)
        this.updatePassword2 = this.updatePassword2.bind(this)
    }

    componentDidMount() {
        this.props.loadUser(this.props.userID)
    }

    componentWillReceiveProps(props) {
        if (props.user) {
            this.setState({ email: props.user.email })
        }
    }

    updateEmail(e) {
        this.setState({ email: e.target.value })
    }

    updatePassword(e) {
        this.setState({ password: e.target.value })
    }

    updatePassword2(e) {
        this.setState({ password2: e.target.value })
    }

    render() {
        let props = this.props
        if (props.loading) {
            return <div>Loading...</div>
        }
        return (
            <div>
                <h1>Users</h1>

                <h2>{props.user.username}</h2>

                <form>
                    <div className="form-group">
                        <label htmlFor="email">Email address</label>
                        <input
                            type="email"
                            className="form-control"
                            id="email"
                            value={this.state.email}
                            onChange={this.updateEmail}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password">Change password</label>
                        <input
                            type="password"
                            className="form-control"
                            id="paswword"
                            value={this.state.password}
                            onChange={this.updatePassword}
                        />
                    </div>

                    <div className="form-group">
                        <label htmlFor="password2">New password again</label>
                        <input
                            type="password"
                            className="form-control"
                            id="paswword2"
                            value={this.state.password2}
                            onChange={this.updatePassword2}
                        />
                    </div>
                </form>
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        user: state.users.user,
        loading: state.users.loading,
        userID: ownProps.match.params.id
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUser
        },
        dispatch
    )

export default withRouter(
    connect(mapStateToProps, mapDispatchToProps)(UserDetail)
)
