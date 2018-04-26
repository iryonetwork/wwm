import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"

import { loadUsers } from "../../modules/users"
import { loadRules, saveRule, deleteRule } from "../../modules/rules"
import { loadRoles } from "../../modules/roles"
import { ADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"

const Read = 1
const Write = 2
const Delete = 4

class Rules extends React.Component {
    constructor(props) {
        super(props)
        this.state = { loading: true }
    }

    componentDidMount() {
        if (!this.props.users) {
            this.props.loadUsers()
        }
        if (!this.props.allRules) {
            this.props.loadRules()
        }
        if (!this.props.roles) {
            this.props.loadRoles()
        }
        if (this.props.canSee === undefined || this.props.canEdit === undefined) {
            this.props.loadUserRights()
        }
        if (this.props.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(this.props)
    }

    componentWillReceiveProps(nextProps) {
        if (!nextProps.users && this.props.users) {
            this.props.loadUsers()
        }
        if (!nextProps.allRules && this.props.allRules) {
            this.props.loadRules()
        }
        if (!nextProps.roles && this.props.roles) {
            this.props.loadRoles()
        }
        if ((nextProps.canSee === undefined || nextProps.canEdit === undefined) && !nextProps.validationsLoading) {
            this.props.loadUserRights()
        }
        if (nextProps.canSee === false) {
            this.props.history.push(`/`)
        }

        this.determineState(nextProps)
    }

    determineState(props) {
        let loading =
            !props.roles ||
            props.rolesLoading ||
            !props.rules ||
            props.rulesLoading ||
            !props.users ||
            props.usersLoading ||
            props.canEdit === undefined ||
            props.canSee === undefined ||
            props.validationsLoading
        this.setState({ loading: loading })

        if (props.rules) {
            let rules = props.rules
            if (_.isArray(rules)) {
                rules = _.fromPairs(_.map(rules, ruleID => [ruleID, props.allRules[ruleID]]))
            }

            let newRules = []

            if (this.state.rules) {
                let editing = _.fromPairs(_.map(_.filter(this.state.rules, rule => rule.edit && !rule.saving), rule => [rule.id, rule]))
                newRules = _.filter(this.state.rules, rule => !rule.saving && !rule.id)

                rules = _.mapValues(rules, rule => {
                    if (editing[rule.id]) {
                        return editing[rule.id]
                    }
                    return rule
                })
            }

            rules = _.map(rules, rule => _.clone(rule))

            this.setState({ rules: rules.concat(newRules) })
        }
    }

    subjectName(subjectID) {
        if (this.props.users && this.props.users[subjectID]) {
            return this.props.users[subjectID].username
        }
        if (this.props.roles && this.props.roles[subjectID]) {
            return this.props.roles[subjectID].name
        }
        return subjectID
    }

    editRule = index => e => {
        let rules = [...this.state.rules]
        rules[index].edit = !rules[index].edit

        if (!rules[index].edit) {
            if (rules[index].id) {
                rules[index] = _.clone(this.props.allRules[rules[index].id])
            } else {
                rules.splice(index, 1)
            }
        }

        this.setState({ rules })
    }

    editAction = (index, action) => e => {
        let rules = [...this.state.rules]

        if (e.target.checked) {
            rules[index].action |= action
        } else {
            rules[index].action &= ~action
        }

        this.setState({ rules })
    }

    editResource = index => e => {
        let rules = [...this.state.rules]
        rules[index].resource = e.target.value

        this.setState({ rules })
    }

    editSubject = index => e => {
        let rules = [...this.state.rules]
        rules[index].subject = e.target.value

        this.setState({ rules })
    }

    editDeny = index => e => {
        let rules = [...this.state.rules]
        rules[index].deny = e.target.value

        this.setState({ rules })
    }

    saveRule = index => e => {
        let rules = [...this.state.rules]
        rules[index].index = index
        rules[index].saving = true

        this.props.saveRule(this.state.rules[index])

        this.setState({ rules })
    }

    newRule = () => e => {
        if (this.state.rules) {
            let rules = [...this.state.rules, { edit: true, resource: "", subject: this.props.subjectID }]

            this.setState({ rules })
        }
    }

    deleteRule = ruleID => e => {
        this.props.deleteRule(ruleID)
    }

    render() {
        let props = this.props
        if (this.state.loading) {
            return <div>Loading...</div>
        }
        if (!props.canSee || props.forbidden) {
            return null
        }

        return (
            <div>
                {props.embedded ? <h3>ACL</h3> : <h1>ACL</h1>}

                <table className="table table-hover text-center">
                    <thead>
                        <tr>
                            <th scope="col">#</th>
                            {!props.embedded ? <th scope="col">Subject</th> : null}
                            <th scope="col">Resource</th>
                            <th scope="col" />
                            <th scope="col" className="text-center">
                                Read
                            </th>
                            <th scope="col" className="text-center">
                                Write
                            </th>
                            <th scope="col" className="text-center">
                                Delete
                            </th>
                            <th scope="col" />
                        </tr>
                    </thead>
                    <tbody>
                        {this.state.rules
                            ? _.map(this.state.rules, (rule, i) => (
                                  <tr key={rule.id || i}>
                                      <th scope="row">{i + 1}</th>
                                      {!props.embedded ? (
                                          <td>
                                              {props.canEdit && rule.edit ? (
                                                  <select className="form-control form-control-sm" value={rule.subject} onChange={this.editSubject(i)}>
                                                      <option>Select subject</option>
                                                      <optgroup label="Roles">
                                                          {_.map(props.roles, role => (
                                                              <option key={role.id} value={role.id}>
                                                                  {role.name}
                                                              </option>
                                                          ))}
                                                      </optgroup>

                                                      <optgroup label="Users">
                                                          {_.map(props.users, user => (
                                                              <option key={user.id} value={user.id}>
                                                                  {user.username} - {user.email}
                                                              </option>
                                                          ))}
                                                      </optgroup>
                                                  </select>
                                              ) : (
                                                  this.subjectName(rule.subject)
                                              )}
                                          </td>
                                      ) : null}

                                      <td>
                                          {props.canEdit && rule.edit ? (
                                              <input
                                                  type="text"
                                                  className="form-control form-control-sm"
                                                  value={rule.resource}
                                                  onChange={this.editResource(i)}
                                              />
                                          ) : (
                                              rule.resource
                                          )}
                                      </td>

                                      <td>
                                          {props.canEdit && rule.edit ? (
                                              <select className="form-control form-control-sm" value={rule.deny || false} onChange={this.editDeny(i)}>
                                                  <option value={false}>Allow</option>
                                                  <option value={true}>Deny</option>
                                              </select>
                                          ) : rule.deny ? (
                                              "Deny"
                                          ) : (
                                              "Allow"
                                          )}
                                      </td>

                                      <td className="text-center">
                                          <input
                                              type="checkbox"
                                              disabled={!props.canEdit || !rule.edit}
                                              onChange={this.editAction(i, Read)}
                                              checked={(rule.action & Read) === Read}
                                          />
                                      </td>
                                      <td className="text-center">
                                          <input
                                              type="checkbox"
                                              disabled={!props.canEdit || !rule.edit}
                                              onChange={this.editAction(i, Write)}
                                              checked={(rule.action & Write) === Write}
                                          />
                                      </td>
                                      <td className="text-center">
                                          <input
                                              type="checkbox"
                                              disabled={!props.canEdit || !rule.edit}
                                              onChange={this.editAction(i, Delete)}
                                              checked={(rule.action & Delete) === Delete}
                                          />
                                      </td>
                                      <td className="text-right">
                                          {props.canEdit ? (
                                              rule.edit ? (
                                                  <div className="btn-group" role="group">
                                                      <button className="btn btn-sm btn-light" disabled={rule.saving} type="button" onClick={this.editRule(i)}>
                                                          <span className="icon_close" />
                                                      </button>
                                                      <button className="btn btn-sm btn-light" disabled={rule.saving} type="button" onClick={this.saveRule(i)}>
                                                          <span className="icon_floppy" />
                                                      </button>
                                                  </div>
                                              ) : (
                                                  <div className="btn-group" role="group">
                                                      <button className="btn btn-sm btn-light" type="button" onClick={this.editRule(i)}>
                                                          <span className="icon_pencil-edit" />
                                                      </button>
                                                      <button className="btn btn-sm btn-light" type="button" onClick={this.deleteRule(rule.id)}>
                                                          <span className="icon_trash" />
                                                      </button>
                                                  </div>
                                              )
                                          ) : null}
                                      </td>
                                  </tr>
                              ))
                            : null}
                    </tbody>
                </table>
                {props.canEdit ? (
                    <button type="button" className="btn btn-sm btn-outline-secondary float-right" onClick={this.newRule()}>
                        Add new ACL rule
                    </button>
                ) : null}
            </div>
        )
    }
}

const mapStateToProps = (state, ownProps) => {
    return {
        users: state.users.users,
        usersLoading: state.users.loading,
        roles: state.roles.roles,
        rolesLoading: state.roles.loading,
        rules: ownProps.rules ? ownProps.rules : state.rules.rules,
        rulesLoading: state.rules.loading,
        embedded: ownProps.rules ? true : false,
        allRules: state.rules.rules,
        subjectID: ownProps.subject,
        canEdit: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        canSee: state.validations.userRights ? state.validations.userRights[ADMIN_RIGHTS_RESOURCE] : undefined,
        validationsLoading: state.validations.loading,
        forbidden: state.rules.forbidden
    }
}

const mapDispatchToProps = dispatch =>
    bindActionCreators(
        {
            loadUsers,
            loadRoles,
            loadRules,
            saveRule,
            deleteRule,
            loadUserRights
        },
        dispatch
    )

export default connect(mapStateToProps, mapDispatchToProps)(Rules)
