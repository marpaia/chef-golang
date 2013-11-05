package chef

import (
	"encoding/json"
)

// chef.GetUsers returns a map of user names to the users RESTful URL as well
// as an error indicating if the request was successful or not.
//
// Usage:
//
//     users, err := chef.GetUsers()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "user" variable which is a map of
//     // user names to user URLs
//     for user := range users {
//         fmt.Println(user)
//      }
func (chef *Chef) GetUsers() (map[string]string, error) {
	resp, err := chef.Get("users")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	users := map[string]string{}
	json.Unmarshal(body, &users)

	return users, nil
}
