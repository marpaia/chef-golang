package chef

import (
	"testing"
)

func TestGetRoles(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	roles, err := chef.GetRoles()
	if err != nil {
		t.Error(err)
	}
	if roles[config.RequiredRole.Name] == "" {
		t.Error("Required role not found")
	}
}

func TestGetRole(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	_, ok, err := chef.GetRole(config.RequiredRole.Name)
	if !ok {
		t.Error("Couldn't find required role")
	}
	if err != nil {
		t.Error(err)
	}
}
