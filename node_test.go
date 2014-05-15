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
	node := new(Node)
	node.Name = "test-node"
	_, ok, err := chef.CreateNode(node)
	if !ok {
		t.Error("Couldn't create required node")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteNode(t *testing.T) {
	chef := testConnectionWrapper(t)
	node := new(Node)
	node.Name = "test-node"
	_, ok, err := chef.DeleteNode(node)
	if !ok {
		t.Error("Couldn't delete required client")
	}
	if err != nil {
		t.Error(err)
	}
}
