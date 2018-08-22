import React from "react"
//import { Route, Link } from "react-router-dom"
import { bindActionCreators } from "redux"
import { connect } from "react-redux"
import _ from "lodash"
import classnames from "classnames"

import { loadUsers } from "../../modules/users"
import { loadRules, saveRule, deleteRule } from "../../modules/rules"
import { loadRoles } from "../../modules/roles"
import { SUPERADMIN_RIGHTS_RESOURCE, loadUserRights } from "../../modules/validations"
import { confirmationDialog } from "shared/utils"

import "../../styles/style.css"

const Read = 1
const Write = 2
const Delete = 4
const Update = 16

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

  editRule(index) {
    return e => {
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
  }

  editAction(index, action) {
    return e => {
      let rules = [...this.state.rules]

      if (e.target.checked) {
        rules[index].action |= action
      } else {
        rules[index].action &= ~action
      }

      this.setState({ rules })
    }
  }

  editResource(index) {
    return e => {
      let rules = [...this.state.rules]
      rules[index].resource = e.target.value

      this.setState({ rules })
    }
  }

  editSubject(index) {
    return e => {
      let rules = [...this.state.rules]
      rules[index].subject = e.target.value

      this.setState({ rules })
    }
  }

  editDeny(index) {
    return e => {
      let rules = [...this.state.rules]
      rules[index].deny = e.target.value

      this.setState({ rules })
    }
  }

  saveRule(index) {
    return e => {
      let rules = [...this.state.rules]
      rules[index].index = index
      rules[index].saving = true

      this.props.saveRule(this.state.rules[index])

      this.setState({ rules })
    }
  }

  newRule() {
    return e => {
      if (this.state.rules) {
        let rules = [...this.state.rules, { edit: true, resource: "", subject: this.props.subjectID }]

        this.setState({ rules })
      }
    }
  }

  deleteRule(index) {
    return e => {
      confirmationDialog(
        `Click OK to confirm that you want to remove ACL rule for resource ${this.state.rules[index].resource} for subject ${this.subjectName(
          this.state.rules[index].subject
        )}.`,
        () => {
          this.props.deleteRule(this.state.rules[index].id)
        }
      )
    }
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
        {!props.embedded ? (
          <header>
            <h1>Access Control List</h1>
          </header>
        ) : null}
        <table className="table">
          <thead>
            <tr>
              {props.embedded ? (
                <td className="w-7 row-details-header-column">
                  <span className="row-details-icon" />
                </td>
              ) : null}
              <th className="w-5" scope="col">
                #
              </th>
              {!props.embedded ? (
                <th className="w-17" scope="col">
                  Subject
                </th>
              ) : null}
              <th className="w-17" scope="col">
                Resource
              </th>
              <th className="w-10" scope="col" />
              <th className="w-5" scope="col">
                Read
              </th>
              <th className="w-5" scope="col">
                Write
              </th>
              <th className="w-5" scope="col">
                Delete
              </th>
              <th className="w-5" scope="col">
                Update
              </th>
              <th scope="col" />
            </tr>
          </thead>
          <tbody>
            {this.state.rules
              ? _.map(this.state.rules, (rule, i) => (
                  <tr
                    className={classnames({
                      "table-edit": props.canEdit && rule.edit
                    })}
                    key={rule.id || i}
                  >
                    {props.embedded ? <td className="w-7 row-details-header-column" /> : null}
                    <th className="w-5" scope="row">
                      {i + 1}
                    </th>
                    {!props.embedded ? (
                      <td className="w-17">
                        {props.canEdit && rule.edit ? (
                          <select className="form-control" value={rule.subject || ""} onChange={this.editSubject(i)}>
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

                    <td className={classnames({ "w-17": !props.embedded, "w-20": props.embedded })}>
                      {props.canEdit && rule.edit ? (
                        <input type="text" className="form-control" value={rule.resource || ""} onChange={this.editResource(i)} />
                      ) : (
                        rule.resource
                      )}
                    </td>

                    <td className={classnames({ "w-17": !props.embedded, "w-20": props.embedded })}>
                      {props.canEdit && rule.edit ? (
                        <select className="form-control" value={rule.deny || false} onChange={this.editDeny(i)}>
                          <option value={false}>Allow</option>
                          <option value={true}>Deny</option>
                        </select>
                      ) : rule.deny ? (
                        "Deny"
                      ) : (
                        "Allow"
                      )}
                    </td>

                    <td className="w-5">
                      <input
                        className="form-check-input"
                        type="checkbox"
                        disabled={!props.canEdit || !rule.edit}
                        onChange={this.editAction(i, Read)}
                        checked={(rule.action & Read) === Read}
                      />
                    </td>
                    <td className="w-5">
                      <input
                        className="form-check-input"
                        type="checkbox"
                        disabled={!props.canEdit || !rule.edit}
                        onChange={this.editAction(i, Write)}
                        checked={(rule.action & Write) === Write}
                      />
                    </td>
                    <td className="w-5">
                      <input
                        className="form-check-input"
                        type="checkbox"
                        disabled={!props.canEdit || !rule.edit}
                        onChange={this.editAction(i, Delete)}
                        checked={(rule.action & Delete) === Delete}
                      />
                    </td>
                    <td className="w-5">
                      <input
                        className="form-check-input"
                        type="checkbox"
                        disabled={!props.canEdit || !rule.edit}
                        onChange={this.editAction(i, Update)}
                        checked={(rule.action & Update) === Update}
                      />
                    </td>
                    <td className="text-right">
                      {props.canEdit ? (
                        rule.edit ? (
                          <div>
                            <button className="btn btn-secondary" disabled={rule.saving} type="button" onClick={this.editRule(i)}>
                              Cancel
                            </button>
                            <button className="btn btn-primary" disabled={rule.saving} type="button" onClick={this.saveRule(i)}>
                              Save
                            </button>
                          </div>
                        ) : (
                          <div>
                            <button onClick={this.editRule(i)} className="btn btn-link" type="button">
                              Edit
                            </button>
                            <button onClick={this.deleteRule(i)} className="btn btn-link" type="button">
                              <span className="remove-link">Remove</span>
                            </button>
                          </div>
                        )
                      ) : null}
                    </td>
                  </tr>
                ))
              : null}
            {props.canEdit && props.embedded ? (
              <tr className="table-edit">
                <td className="w-7 row-details-header-column" />
                <td colSpan="8">
                  {props.canEdit ? (
                    <button
                      type="button"
                      className="btn btn-link"
                      disabled={this.state.rules.length !== 0 && this.state.rules[this.state.rules.length - 1].edit ? true : null}
                      onClick={this.newRule()}
                    >
                      Add New Access Control List Rule
                    </button>
                  ) : null}
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
        {props.canEdit && !props.embedded ? (
          <button
            type="button"
            disabled={this.state.rules.length !== 0 && this.state.rules[this.state.rules.length - 1].edit ? true : null}
            className="btn btn-link"
            onClick={this.newRule()}
          >
            Add New Access Control List Rule
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
    canEdit: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
    canSee: state.validations.userRights ? state.validations.userRights[SUPERADMIN_RIGHTS_RESOURCE] : undefined,
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
