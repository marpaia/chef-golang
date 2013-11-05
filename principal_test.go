package chef

import (
	"testing"
)

func TestGetPrincipal(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	_, ok, err := chef.GetPrincipal(config.RequiredPrincipal.Name)
	if !ok {
		t.Error("Couldn't find required principal")
	}
	if err != nil {
		t.Error(err)
	}
}
