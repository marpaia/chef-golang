package chef

import (
	"testing"
)

func TestGetEnvironments(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	environments, err := chef.GetEnvironments()
	if err != nil {
		t.Error(err)
	}
	if environments[config.RequiredEnvironment.Name] == "" {
		t.Error("Required environment not found")
	}
}

func TestGetEnvironment(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetEnvironment(config.RequiredEnvironment.Name)
	if !ok {
		t.Error("Couldn't find required environment")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvironmentCookbooks(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, err := chef.GetEnvironmentCookbooks(config.RequiredEnvironment.Name)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvironmentCookbook(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, ok, err := chef.GetEnvironmentCookbook(config.RequiredEnvironment.Name, config.RequiredCookbook.Name)
	if !ok {
		t.Error("Couldn't find cookbook in environment")
	}
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvironmentNodes(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	_, err := chef.GetEnvironmentNodes(config.RequiredEnvironment.Name)
	if err != nil {
		t.Error(err)
	}
}

func TestGetEnvironmentRecipes(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	recipes, err := chef.GetEnvironmentRecipes(config.RequiredEnvironment.Name)
	if err != nil {
		t.Error(err)
	}
	found := false
	for _, recipe := range recipes {
		if recipe == config.RequiredRecipe.Name {
			found = true
		}
	}
	if !found {
		t.Error("Couldn't find required recipe")
	}
}

// For now I am skiping this test because it is a webui only endpoint
func TestGetEnvironmentRole(t *testing.T) {
	chef := testConnectionWrapper(t)
	config := testConfig()
	println("Looking for role: ", config.RequiredRole.Name)
	println("In", config.RequiredEnvironment.Name)
	test, ok, err := chef.GetEnvironmentRole(config.RequiredEnvironment.Name, config.RequiredRole.Name)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Couldn't find required role in required environment")
	}
	t.Log(test)
}
