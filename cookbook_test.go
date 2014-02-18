package chef

import (
	"testing"
)

func TestGetCookbooks(t *testing.T) {
	chef := testConnectionWrapper(t)
	cookbooks, err := chef.GetCookbooks()
	if err != nil {
		t.Error(err)
	}
	found := false
	config := testConfig()
	for cookbook := range cookbooks {
		if cookbook == config.RequiredCookbook.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Couldn't find required cookbook")
	}
}

func TestGetCookbook(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetCookbook(config.RequiredCookbook.Name)
	if !ok {
		t.Error("Couldn't find required cookbook")
	}
	if err != nil {
		t.Error("Couldn't find required cookbook")
	}
}

func TestCookbookVersion(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetCookbookVersion(config.RequiredCookbook.Name, config.RequiredCookbook.Version)
	if !ok {
		t.Error(err)
	}
	if !ok {
		t.Error("Cookbook version not found")
	}
}
