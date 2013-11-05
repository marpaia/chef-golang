package chef

import (
	"testing"
)

func TestGetUsers(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	users, err := chef.GetUsers()
	if err != nil {
		t.Error(err)
	}
	found := false
	for user := range users {
		if user == config.RequiredUser.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Couldn't find required user")
	}
}
