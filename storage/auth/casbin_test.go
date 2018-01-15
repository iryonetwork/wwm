package auth

import (
	"strconv"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/models"
)

func TestRules(t *testing.T) {
	storage := newTestStorage()
	defer storage.Close()

	// add user
	u1, _ := storage.AddUser(&models.User{Username: swag.String("user1")})
	u2, _ := storage.AddUser(&models.User{Username: swag.String("user2")})
	u3, _ := storage.AddUser(&models.User{Username: swag.String("user3")})

	g1, _ := storage.AddRole(&models.Role{Name: swag.String("group1"), Users: []string{u1.ID, u2.ID}})
	g2, _ := storage.AddRole(&models.Role{Name: swag.String("group2"), Users: []string{u2.ID, u3.ID}})

	//storage.enforcer.Enforce

	tests := []struct {
		rules       []*models.Rule
		userID      string
		validations []*models.ValidationPair
		results     []bool
	}{
		{
			rules: []*models.Rule{
				{
					Action:   swag.Int64(Read),
					Subject:  &g1.ID,
					Resource: swag.String("/storage*"),
				},
				{
					Action:   swag.Int64(Read | Write),
					Subject:  &g2.ID,
					Resource: swag.String("/other*"),
				},
				{
					Action:   swag.Int64(Read | Write),
					Deny:     true,
					Subject:  &u2.ID,
					Resource: swag.String("/storage/location*"),
				},
			},
			userID: u2.ID,
			validations: []*models.ValidationPair{
				{
					Actions:  swag.Int64(Write),
					Resource: swag.String("/storage/test"),
				},
				{
					Actions:  swag.Int64(Read),
					Resource: swag.String("/storage/test"),
				},
				{
					Actions:  swag.Int64(Write),
					Resource: swag.String("/other"),
				},
				{
					Actions:  swag.Int64(Read),
					Resource: swag.String("/storage/location/other"),
				},
			},
			results: []bool{
				false,
				true,
				true,
				false,
			},
		},
		{
			rules: []*models.Rule{
				{
					Action:   swag.Int64(Read),
					Subject:  &g1.ID,
					Resource: swag.String("/*"),
				},
				{
					Action:   swag.Int64(Write),
					Subject:  &u1.ID,
					Resource: swag.String("/auth/user*"),
				},
				{
					Action:   swag.Int64(Write),
					Subject:  &g1.ID,
					Deny:     true,
					Resource: swag.String("/auth/user/5"),
				},
			},
			userID: u1.ID,
			validations: []*models.ValidationPair{
				{
					Actions:  swag.Int64(Read),
					Resource: swag.String("/storage"),
				},
				{
					Actions:  swag.Int64(Read),
					Resource: swag.String("/auth/user"),
				},
				{
					Actions:  swag.Int64(Write),
					Resource: swag.String("/auth/user/1"),
				},
				{
					Actions:  swag.Int64(Delete),
					Resource: swag.String("/auth/user/1"),
				},
				{
					Actions:  swag.Int64(Delete | Write),
					Resource: swag.String("/auth/user/1"),
				},
				{
					Actions:  swag.Int64(Write),
					Resource: swag.String("/auth/user/5"),
				},
			},
			results: []bool{
				true,
				true,
				true,
				false,
				false,
				false,
			},
		},
		{
			rules:  []*models.Rule{},
			userID: u3.ID,
			validations: []*models.ValidationPair{
				{
					Actions:  swag.Int64(Read),
					Resource: swag.String("/storage"),
				},
			},
			results: []bool{
				false,
			},
		},
	}

	for testIndex, test := range tests {
		for _, rule := range test.rules {
			storage.AddRule(rule)
		}

		storage.enforcer.LoadPolicy()

		for i, validation := range test.validations {
			res := storage.enforcer.Enforce(test.userID, *validation.Resource, strconv.FormatInt(*validation.Actions, 10))
			if res != test.results[i] {
				t.Fatalf("Test %d: Expected validation '%v' to be %t", testIndex, validation, test.results[i])
			}
		}

		for _, rule := range test.rules {
			storage.RemoveRule(rule.ID)
		}
	}
}
