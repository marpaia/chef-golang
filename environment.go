package chef

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chef.Environment dinfes the relevant parameters of a Chef environment. This
// includes the name of the environment, the description strings, etc.
type Environment struct {
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	CookbookVersions   map[string]string      `json:"cookbook_versions"`
	JSONClass          string                 `json:"json_class"`
	ChefType           string                 `json:"chef_type"`
	DefaultAttributes  map[string]interface{} `json:"default_attributes"`
	OverrideAttributes map[string]interface{} `json:"override_attributes"`
}

// chef.GetEnvironments returns a map of environment names to the environment's
// RESTful URL as well as an error indicating if the request was successful or
// not.
//
// Usage:
//
//     environments, err := chef.GetEnvironments()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "environments" variable which is a map of
//     // environment names to environment URLs
//     for environment := range environments {
//         fmt.Println(environment)
//      }
func (chef *Chef) GetEnvironments() (map[string]string, error) {
	resp, err := chef.Get("environments")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	environments := map[string]string{}
	json.Unmarshal(body, &environments)

	return environments, nil
}

// chef.GetEnvironment accepts a string which represents the name of a Chef
// environment and returns a chef.Environment type representing that environment
// as well as a bool indicating whether or not the environment was found and an
// error indicating if the request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     environment, ok, err := chef.GetEnvironment("production")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that environment!")
//     } else {
//         // do what you please with the "environment" variable which is of the
//         // *Chef.Environment type
//         fmt.Printf("%#v\n", environment)
//     }
func (chef *Chef) GetEnvironment(name string) (*Environment, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s", name))
	if err != nil {
		return nil, false, err
	}
	body, err := responseBody(resp)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, false, nil
		}
		return nil, false, err
	}

	environment := new(Environment)
	json.Unmarshal(body, environment)

	return environment, true, nil
}

// chef.GetEnvironmentCookbooks accepts a string which represents the name of a
// Chef environment and returns a map of cookbook names to a *Chef.Cookbook type
// as well as an error indicating whether or not the request failed.
//
// Usage:
//
//     cookbooks, err := chef.GetEnvironmentCookbooks("production")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "cookbooks" variable which is a map of
//     // cookbook names to chef.Cookbook types
//     for name, cookbook := range cookbooks {
//         fmt.Println(name, cookbook.Version[0])
//      }
func (chef *Chef) GetEnvironmentCookbooks(name string) (map[string]*Cookbook, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s/cookbooks", name))
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	cookbooks := map[string]*Cookbook{}
	json.Unmarshal(body, &cookbooks)

	return cookbooks, nil
}

// chef.GetEnvironmentCookbook accepts a string which represents the name of a
// Chef environment as well as a string which represent the name of a cookbook
// and returns a *Chef.Cookbook type, a bool indicating whether or not the
// cookbook was found in the given environment as well as an error indicating
// whether or not the request failed.
//
// Usage:
//
//     cookbook, ok, err := chef.GetEnvironmentCookbook("production", "apache")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that cookbook!")
//     } else {
//         // do what you please with the "cookbook" variable which is of the
//         // *Chef.Cookbook type
//         fmt.Printf("%#v\n", cookbook)
//     }
func (chef *Chef) GetEnvironmentCookbook(env, cb string) (*Cookbook, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s/cookbooks/%s", env, cb))
	if err != nil {
		return nil, false, err
	}
	body, err := responseBody(resp)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, false, nil
		}
		return nil, false, err
	}

	cookbook := map[string]*Cookbook{}
	json.Unmarshal(body, &cookbook)

	return cookbook[cb], true, nil
}

// chef.GetEnvironmentNodes accepts a string which represents the name of a
// Chef environment as well as a string which represent the name of a node
// and returns a map of node names to their corresponding RESTful URL, a bool
// indicating whether or not the cookbook was found in the given environment as
// well as an error indicating whether or not the request failed.
//
// Usage:
//
//     nodes, err := chef.GetEnvironmentNodes("production")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "nodes" variable which is a map of
//     // node names to their corresponding RESTful URL
//     for node := range nodes {
//         fmt.Println(node)
//      }
func (chef *Chef) GetEnvironmentNodes(name string) (map[string]string, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s/nodes", name))
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	nodes := map[string]string{}
	json.Unmarshal(body, &nodes)

	return nodes, nil
}

// chef.GetEnvironmentRecipes accepts a string which represents the name of a
// Chef environment and returns a slice of recipe names as well as an error
// indicating whether or not the request failed.
//
// Usage:
//
//     recipes, err := chef.GetEnvironmentRecipes("production")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "recipes" variable which is a slice of
//     // recipe names
//     for recipe := range recipes {
//         fmt.Println(recipe)
//      }
func (chef *Chef) GetEnvironmentRecipes(name string) ([]string, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s/recipes", name))
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	recipes := []string{}
	json.Unmarshal(body, &recipes)

	return recipes, nil
}

// chef.GetEnvironmentRole accepts a string which represents the name of a
// Chef environment as well as a string which represent the name of a role
// and returns a map of strings (which represent a role attribute like a
// runlist) to a slice of strings which represents the relevant information with
// regards to that attribute, a bool indicating whether or not the role was
// found in the given environment as well as an error indicating whether or not
// the request failed.
//
// Usage:
//
//     role, ok, err := chef.GetEnvironmentRole("production", "base")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that role!")
//     } else {
//         // do what you please with the "role" variable
//         fmt.Println(role)
//     }
func (chef *Chef) GetEnvironmentRole(env, rol string) (map[string][]string, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("environments/%s/roles/%s", env, "WebBase"))
	if err != nil {
		return nil, false, err
	}
	body, err := responseBody(resp)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, false, nil
		}
		return nil, false, err
	}

	role := map[string][]string{}
	json.Unmarshal(body, &role)

	return role, true, nil
}
