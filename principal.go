package chef

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chef.GetPrincipal returns a map of principal item names to their
// corresponding RESTful url. It also returns a bool indicating whether or not
// the client was found and an error indicating if the request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     principal, ok, err := chef.GetCookbookPrincipal("neo4j.example.org")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that principal!")
//     } else {
//         // do what you please with the "principal" variable which is of the
//         // map[string]string type
//         fmt.Println(principal)
//     }
func (chef *Chef) GetPrincipal(name string) (map[string]string, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("principals/%s", name), nil)
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

	principal := map[string]string{}
	json.Unmarshal(body, &principal)

	return principal, true, nil
}
