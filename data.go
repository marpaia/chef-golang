package chef

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chef.GetData returns a map of databag names to their related REST URL
// endpoint as well as an error indicating if the request was successful or not
//
// Usage:
//
//     data, err := chef.GetData()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     for d := range data {
//         fmt.Println(d)
//     }
func (chef *Chef) GetData() (map[string]string, error) {
	resp, err := chef.Get("data")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	data := map[string]string{}
	json.Unmarshal(body, &data)

	return data, nil
}

// chef.GetDataByName accept a string which represents the name of a specific
// databag and returns a map of information about that databag, a bool
// indicating whether or not the databag was found and an error indicating if
// the request failed or not.
//
// Note that if the request is successful but no such data item existed, the
// error return value will be nil but the bool will be false
//
// Usage:
//
//     data, ok, err := chef.GetDataByName("apache")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that databag!")
//     } else {
//         // do what you please with the "data" variable which is of the
//         // map[string]string type
//         fmt.Println(data)
//     }
func (chef *Chef) GetDataByName(name string) (map[string]string, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("data/%s", name))
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

	data := map[string]string{}
	json.Unmarshal(body, &data)

	return data, true, nil
}
