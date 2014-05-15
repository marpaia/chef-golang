package chef

import (
	"testing"
)

func TestGetNodes(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	nodes, err := chef.GetNodes()
	if err != nil {
		t.Error(err)
	}
	if nodes[config.RequiredNode.Name] == "" {
		t.Error("Required node not found")
	}
}

func TestGetNode(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetNode(config.RequiredNode.Name)
	if !ok {
		t.Error("Couldn't find required node")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNode(t *testing.T) {
	chef := testConnectionWrapper(t)
	name := "test-node"
	environment := "_default"
	normal := make(map[string]interface{})
	overrides := make(map[string]interface{})
	defaults := make(map[string]interface{})
	run_list := make([]string, 0)

	normal["hi"] = "hello"
	overrides["hi"] = "hello"
	defaults["hi"] = "hello"

	_, ok, err := chef.CreateNode(name, environment, normal, overrides, defaults, run_list)
	if !ok {
		t.Error("Couldn't create required node")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteNode(t *testing.T) {
	chef := testConnectionWrapper(t)
	_, ok, err := chef.DeleteNode("test-node")
	if !ok {
		t.Error("Couldn't delete required client")
	}
	if err != nil {
		t.Error(err)
	}
}
