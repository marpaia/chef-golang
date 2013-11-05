package chef

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chef.Role represents the relevant attributes of a Chef role
type Role struct {
	Name               string                 `json:"name"`
	ChefType           string                 `json:"chef_type"`
	JSONClass          string                 `json:"json_class"`
	DefaultAttributes  map[string]interface{} `json:"default_attributes"`
	OverrideAttributes map[string]interface{} `json:"override_attributes"`
	RunList            []string               `json:"run_list"`
}

// chef.GetRoles returns a map of role names to a string which represents the
// role's RESTful URL as well as an error indicating if the request was
// successful or not.
//
// Usgae:
//
//     roles, err := chef.GetRoles()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "roles" variable which is a map of
//     // role names to their RESTful URLs
//     for role := range roles {
//         fmt.Println(role)
//      }
func (chef *Chef) GetRoles() (map[string]string, error) {
	resp, err := chef.Get("roles")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	roles := map[string]string{}
	json.Unmarshal(body, &roles)

	return roles, nil
}

// chef.GetRole returns a pointer to the chef.Role type for a given string that
// represents a role name. It also returns a bool indicating whether or not the
// client was found and an error indicating if the request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     role, ok, err := chef.GetRole("neo4j")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that role!")
//     } else {
//         // do what you please with the "role" variable which is of the
//         // *Chef.Role
//         fmt.Printf("%#v\n", role)
//     }
func (chef *Chef) GetRole(name string) (*Role, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("roles/%s", name))
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

	role := new(Role)
	json.Unmarshal(body, role)

	return role, true, nil
}
