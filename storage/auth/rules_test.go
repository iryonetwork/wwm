package auth

import (
	"reflect"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/iryonetwork/wwm/gen/auth/models"
	"github.com/iryonetwork/wwm/log/errorChecker"
	"github.com/iryonetwork/wwm/utils"
)

// method to ensure that rules used for tests are always fresh
func getTestRules() (*models.Rule, *models.Rule) {
	testRule := &models.Rule{
		Action:   swag.Int64(1),
		Resource: swag.String("/"),
	}
	testRule2 := &models.Rule{
		Action:   swag.Int64(1),
		Resource: swag.String("/auth/*"),
	}

	return testRule, testRule2
}

func TestAddRule(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testRule, _ := getTestRules()
	_, testUser2 := getTestUsers()

	// add user
	_, err := storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	testRule.Subject = &testUser2.ID

	// add rule
	rule, err := storage.AddRule(testRule)
	if rule.ID == "" {
		t.Fatalf("Expected ID to be set, got an empty string")
	}
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// load policy manually because it's launched in goroutine
	errorChecker.FatalTesting(t, storage.enforcer.LoadPolicy())

	// can't add rule that already exist
	_, err = storage.AddRule(testRule)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	// add rule with invalid user id
	testRule.Subject = swag.String("wrong")
	_, err = storage.AddRule(testRule)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}

	rules, err := storage.GetRules()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(rules) != 1 {
		t.Fatalf("Expected 1 rules; got %d", len(rules))
	}
}

func TestGetRule(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testRule, _ := getTestRules()
	_, testUser2 := getTestUsers()

	// add user and rule
	_, err := storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)

	testRule.Subject = &testUser2.ID
	_, err = storage.AddRule(testRule)
	errorChecker.FatalTesting(t, err)

	// get rule
	rule, err := storage.GetRule(testRule.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if !reflect.DeepEqual(*testRule, *rule) {
		t.Fatalf("Expected returned rule to be '%v', got '%v'", *testRule, *rule)
	}

	// get rule with wrong uuid
	_, err = storage.GetRule("something")
	uErr, ok := err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrBadRequest {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrBadRequest, uErr.Code())
	}

	// get non existing rule
	_, err = storage.GetRule("E4363A8D-4041-4B17-A43E-17705C96C1CD")
	uErr, ok = err.(utils.Error)
	if !ok {
		t.Fatalf("Expected error to be of type 'utils.Error'; got '%T'", err)
	}
	if uErr.Code() != utils.ErrNotFound {
		t.Fatalf("Expected error code to be '%s'; got '%s'", utils.ErrNotFound, uErr.Code())
	}
}

func TestGetRules(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testRule, testRule2 := getTestRules()
	_, testRole2 := getTestRoles()
	_, testUser2 := getTestUsers()

	// add user, role and rules
	_, err := storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRole(testRole2)
	errorChecker.FatalTesting(t, err)
	testRule.Subject = &testUser2.ID
	testRule2.Subject = &testRole2.ID

	_, err = storage.AddRule(testRule)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRule(testRule2)
	errorChecker.FatalTesting(t, err)

	// get rules
	rules, err := storage.GetRules()
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}
	if len(rules) != 2 {
		t.Fatalf("Expected 2 rules; got %d", len(rules))
	}

	rulesMap := map[string]*models.Rule{}
	for _, rule := range rules {
		rulesMap[rule.ID] = rule
	}

	if !reflect.DeepEqual(*testRule, *rulesMap[testRule.ID]) {
		t.Fatalf("Expected rule one to be '%v', got '%v'", *testRule, *rulesMap[testRule.ID])
	}

	if !reflect.DeepEqual(*testRule2, *rulesMap[testRule2.ID]) {
		t.Fatalf("Expected rule one to be '%v', got '%v'", *testRule2, *rulesMap[testRule2.ID])
	}
}

func TestUpdateRule(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testRule, _ := getTestRules()
	_, testRole2 := getTestRoles()
	_, testUser2 := getTestUsers()

	// add user, role and rule
	_, err := storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	_, err = storage.AddRole(testRole2)
	errorChecker.FatalTesting(t, err)
	testRule.Subject = &testUser2.ID
	_, err = storage.AddRule(testRule)
	errorChecker.FatalTesting(t, err)

	// update rule
	updateRule := &models.Rule{
		ID:      testRule.ID,
		Action:  swag.Int64(3),
		Subject: &testRole2.ID,
	}
	rule, err := storage.UpdateRule(updateRule)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	if !reflect.DeepEqual(*rule, *updateRule) {
		t.Fatalf("Expected rule one to be '%v', got '%v'", *rule, *updateRule)
	}

	// update rule with invalid subject
	updateRule.Subject = swag.String("wrong")
	_, err = storage.UpdateRule(updateRule)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}

func TestRemoveRule(t *testing.T) {
	storage := newTestStorage(nil)
	defer storage.Close()

	testRule, _ := getTestRules()
	_, testUser2 := getTestUsers()

	// add user and rule
	_, err := storage.AddUser(testUser2)
	errorChecker.FatalTesting(t, err)
	testRule.Subject = &testUser2.ID
	_, err = storage.AddRule(testRule)
	errorChecker.FatalTesting(t, err)

	// remove rule
	err = storage.RemoveRule(testRule.ID)
	if err != nil {
		t.Fatalf("Expected error to be nil; got '%v'", err)
	}

	// remove rule again
	err = storage.RemoveRule(testRule.ID)
	if err == nil {
		t.Fatalf("Expected error; got nil")
	}
}
