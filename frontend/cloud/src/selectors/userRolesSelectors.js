import _ from "lodash"
import { createSelector } from 'reselect'

const getUserUserRoles = (state, props) =>
  state.userRoles.userUserRoles ? (state.userRoles.userUserRoles[props.userID] ? state.userRoles.userUserRoles[props.userID] : {}) : {}

const getOrganizationUserRoles = (state, props) =>
  state.userRoles.domainUserRoles ? ((state.userRoles.domainUserRoles["organization"] && state.userRoles.domainUserRoles["organization"][props.organizationID]) ? state.userRoles.domainUserRoles["organization"][props.organizationID] : {}) : {}

const getClinicUserRoles = (state, props) =>
  state.userRoles.domainUserRoles ? ((state.userRoles.domainUserRoles["clinic"] && state.userRoles.domainUserRoles["clinic"][props.clinicID]) ? state.userRoles.domainUserRoles["clinic"][props.clinicID] : {}) : {}

const getUserID = (state, props) => props.userID

const getOrganizationID = (state, props) => props.organizationID

const getClinicID = (state, props) => props.clinicID

const getOrganizations = (state, props) => state.organizations.organizations ? state.organizations.organizations : {}

export const makeGetOrganizationUserUserRoles = () => {
  return createSelector(
    [ getUserUserRoles, getOrganizationID ],
    (userRoles, organizationID) => {
      return _.pickBy(userRoles, userRole => (userRole.domainType === "organization" && userRole.domainID === organizationID))
    }
  )
}

export const makeGetClinicUserUserRoles = () => {
  return createSelector(
    [ getUserUserRoles, getClinicID ],
    (userRoles, clinicID) => {
        return _.pickBy(userRoles, userRole => (userRole.domainType === "clinic" && userRole.domainID === clinicID))
    }
  )
}

export const makeGetWildcardUserUserRoles = () => {
  return createSelector(
    [ getUserUserRoles ],
    (userRoles) => {
        return _.pickBy(userRoles, userRole => userRole.domainID === "*")
    }
  )
}

export const makeGetUserOrganizationIDs = () => {
  return createSelector(
    [ getUserUserRoles ],
    (userRoles) => {
        return _.uniq(_.map(_.pickBy(userRoles, userRole => (userRole.domainType === "organization" && userRole.domainID !== "*")), 'domainID'))
    }
  )
}

export const makeGetUserClinicIDs = () => {
  return createSelector(
    [ getUserUserRoles ],
    (userRoles) => {
      return _.uniq(_.map(_.pickBy(userRoles, userRole => (userRole.domainType === "clinic" && userRole.domainID !== "*")), 'domainID'))
    }
  )
}

export const makeGetUserAllowedClinicIDs = () => {
  return createSelector(
    [ makeGetUserOrganizationIDs(), getOrganizations ],
    (userOrganizationIDs, organizations) => {
        let clinicIDs = []
        _.forEach(userOrganizationIDs, organizationID => {
          let organizationClinicIDs = organizations[organizationID] ? (organizations[organizationID].clinics ? organizations[organizationID].clinics : []) : []
          clinicIDs.push(...organizationClinicIDs)
        })

        return clinicIDs
    }
  )
}

export const makeGetOrganizationUserIDs = () => {
  return createSelector(
    [ getOrganizationUserRoles ],
    (userRoles) => {
        return _.uniq(_.map(userRoles, 'userID'))
    }
  )
}

export const makeGetClinicUserIDs = () => {
  return createSelector(
    [ getClinicUserRoles ],
    (userRoles) => {
        return _.uniq(_.map(userRoles, 'userID'))
    }
  )
}

export const makeGetUserOrganizationUserRoles = () => {
  return createSelector(
    [ getOrganizationUserRoles, getUserID ],
    (userRoles, userID) => {
      return _.pickBy(userRoles, userRole => (userRole.userID === userID))
    }
  )
}

export const makeGetUserClinicUserRoles = () => {
  return createSelector(
    [ getClinicUserRoles, getUserID ],
    (userRoles, userID) => {
        return _.pickBy(userRoles, userRole => (userRole.userID === userID))
    }
  )
}
