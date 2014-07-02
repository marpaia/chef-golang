package chef

import (
	"fmt"
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

func TestCreateClient(t *testing.T) {
	chef := testConnectionWrapper(t)
	client := new(Client)
	client.Name = "test-client"
	client.Admin = true
	client, ok, err := chef.CreateClient(client)
	if !ok {
		t.Error("Couldn't create required client")
	}
	if err != nil {
		t.Error(err)
	}
	if client.URI != fmt.Sprintf("http://localhost:8443/clients/%s", "test-client") {
		t.Error("Client URI doesn't match", client.URI)
	}
	if ok && client.PrivateKey == "" {
		t.Error("New client private key was empty")
	}
}

func TestDeleteClient(t *testing.T) {
	chef := testConnectionWrapper(t)
	client := new(Client)
	client.Name = "test-client"
	ok, err := chef.DeleteClient(client)
	if !ok {
		t.Error("Couldn't delete required client")
	}
	if err != nil {
		t.Error(err)
	}
}
