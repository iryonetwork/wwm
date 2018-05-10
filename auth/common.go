// package for common auth values and methods
package auth

import (
	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/auth/models"
)

var DomainTypeClinic = "clinic"
var DomainTypeLocation = "location"
var DomainTypeOrganization = "organization"
var DomainTypeUser = "user"
var DomainTypeGlobal = "global"

var DomainIDWildcard = "*"

// predefined role that is given globally to every user
var EveryoneRole = &models.Role{
	ID:   "338fae76-9859-4803-8441-c5c441319cfd",
	Name: swag.String("Everyone"),
}

// predefined role for acl updating/deleting files
var AuthorRole = &models.Role{
	ID:   "d91ff466-125b-4250-9dde-bbc282489430",
	Name: swag.String("Author"),
}

// predefined role to determine basic membership in clinic / organization
var MemberRole = &models.Role{
	ID:   "a422f7f5-291b-4454-ae61-3d98c6091c3e",
	Name: swag.String("Basic member"),
}

// predefined superadmin rol
var SuperadminRole = &models.Role{
	ID:   "3720198b-74ed-40de-a45e-8756f22e67d2",
	Name: swag.String("Superadmin"),
}
