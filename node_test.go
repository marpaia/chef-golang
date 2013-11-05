package chef

import (
	"testing"
)

func TestGetNodes(t *testing.T) {
	chef := testConnectionWrapper()
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
	chef := testConnectionWrapper()
	config := testConfig()
	_, ok, err := chef.GetNode(config.RequiredNode.Name)
	if !ok {
		t.Error("Couldn't find required node")
	}
	if err != nil {
		t.Error(err)
	}
}
