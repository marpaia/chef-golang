package chef

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// chef.Client defines the relevant parameters of a Chef client. This includes
// it's name, whether or not it's an admin, etc.
type Client struct {
	Name        string `json:"name"`
	Admin       bool   `json:"admin"`
	JSONClass   string `json:"json_class,omitempty"`
	ChefType    string `json:"chef_type,omitempty"`
	ClientName  string `json:"clientname,omitempty"`
	Org         string `json:"orgname,omitempty"`
	Validator   bool   `json:"validator,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
	URI         string `json:"uri,omitempty"`
}

// chef.GetClients returns a map of client name's to client REST urls as well as
// an error indicating if the request was successful or not.
//
// Usage:
//
//     clients, err := chef.GetClients()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "clients" variable which is map of client
//     // names to client REST urls
//     for client := range clients {
//         fmt.Println("Client:", client)
//     }
func (chef *Chef) GetClients() (map[string]string, error) {
	resp, err := chef.Get("clients")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	clients := map[string]string{}
	json.Unmarshal(body, &clients)

	return clients, nil
}

// GetClient accept a string representing the client name and returns a Client
// type which illustrates various information about the client. It also returns
// a bool indicating whether or not the client was found and an error indicating
// if the request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     client, ok, err := chef.GetClient("clientname")
//     if err != nil {
//         fmt.Println(err)
//		   os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that client!")
//     } else {
//         // do what you please with the "client" variable which is of the
//         // *Chef.Client type
//         fmt.Printf("%#v\n", client)
//     }
func (chef *Chef) GetClient(name string) (*Client, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("clients/%s", name))
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

	client := new(Client)
	json.Unmarshal(body, client)

	return client, true, nil
}

// CreateClient accepts a Client struct and admin bool to create a new client.
// It returns a Client struct that includes a PrivateKey, PublicKey, and URI as
// described in http://docs.opscode.com/api_chef_server.html#post
func (chef *Chef) CreateClient(client *Client) (*Client, bool, error) {
	jsonBody, err := json.Marshal(client)
	if err != nil {
		return client, false, err
	}

	postBody := bytes.NewReader(jsonBody)
	resp, err := chef.Post("clients", nil, postBody)
	if err != nil {
		return client, false, err
	}
	if resp.StatusCode != 201 {
		err = errors.New(fmt.Sprintf("Server returned %s", resp.Status))
		return client, false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, false, err
	}

	requestError := new(Error)
	json.Unmarshal(body, requestError)

	if len(requestError.Error) != 0 {
		err = errors.New(requestError.Error[0])
		return client, false, err
	}

	json.Unmarshal(body, &client)
	return client, true, nil
}

// DeleteClient accepts a Client struct and returns a success bool and error
func (chef *Chef) DeleteClient(client *Client) (bool, error) {
	params := make(map[string]string)
	_, err := chef.Delete(fmt.Sprintf("clients/%s", client.Name), params)
	if err != nil {
		return false, err
	}
	return true, nil
}
