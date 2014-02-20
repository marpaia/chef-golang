package chef

import (
	"testing"
)

func TestGetClients(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	clients, err := chef.GetClients()
	if err != nil {
		t.Error(err)
	}
	if clients[config.RequiredClient.Name] == "" {
		t.Error("Required client not found")
	}
}

func TestGetClient(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetClient(config.RequiredClient.Name)
	if !ok {
		t.Error("Couldn't find required client")
	}
	if err != nil {
		t.Error(err)
	}
}
