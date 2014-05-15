package chef

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// chef.Node represents the relevant parameters of a Chef node
type Node struct {
	Name        string   `json:"name"`
	Environment string   `json:"chef_environment"`
	JSONClass   string   `json:"json_class"`
	RunList     []string `json:"run_list"`
	ChefType    string   `json:"chef_type"`
	Info        struct {
		Languages map[string]map[string]string `json:"languages"`
		Kernel    struct {
			Name    string                       `json:"name"`
			Release string                       `json:"release"`
			Version string                       `json:"version"`
			Machine string                       `json:"machine"`
			Modules map[string]map[string]string `json:"modules"`
		} `json:"kernel"`
		OS        string `json:"os"`
		OSVersion string `json:"os_version"`
		Hostname  string `json:"hostname"`
		FQDN      string `json:"fqdn"`
		Domain    string `json:"domain"`
		Network   struct {
			Interfaces map[string]struct {
				Type          string `json:"type"`
				Number        string `json:"number"`
				Encapsulation string `json:"encapsulation"`
				Addresses     map[string]struct {
					Family    string `json:"family"`
					Broadcast string `json:"broadcast"`
					Netmast   string `json:"netmast"`
				} `json:"addresses"`
				Flags []string          `json:"flags"`
				MTU   string            `json:"mtu"`
				Arp   map[string]string `json:"arp"`
			} `json:"interfaces"`
		} `json:"network"`
		IPAddress       string                       `json:"ipaddress"`
		MACAddress      string                       `json:"macaddress"`
		ChefPackages    map[string]map[string]string `json:"chef_packages"`
		Keys            map[string]map[string]string `json:"keys"`
		Platform        string                       `json:"platform"`
		PlatformVersion string                       `json:"platform_version"`
		CPU             map[string]interface{}       `json:"cpu"`
		Filesystem      map[string]struct {
			KBSize       string   `json:"ks_size"`
			KBUsed       string   `json:"ks_used"`
			KBavailable  string   `json:"ks_available"`
			PercentUsed  string   `json:"percent_used"`
			Mount        string   `json:"mount"`
			FSType       string   `json:"fs_type"`
			MountOptions []string `json:"mount_options"`
		} `json:"filesystem"`
		Memory          map[string]interface{} `json:"memory"`
		UptimeSeconds   int                    `json:"uptime_seconds"`
		Uptime          string                 `json:"uptime"`
		IdletimeSeconds int                    `json:"idletime_seconds"`
		Idletime        string                 `json:"idletime"`
		BlockDevice     map[string]interface{} `json:"block_device"`
		Datacenter      map[string]interface{} `json:"datacenter"`
		Ipmi            struct {
			Address    string `json:"address"`
			MacAddress string `json:"mac_address"`
		} `json:"ipmi"`
		Recipes []string `json:"recipes"`
		Roles   []string `json:"roles"`
	} `json:"automatic"`
	Normal  map[string]interface{} `json:"normal"`
	Default map[string]interface{} `json:"default"`
}

// chef.PutNode represents the node structure for PUTs/POSTs
type PutNode struct {
	Name        string                 `json:"name"`
	ChefType    string                 `json:"chef_type"`
	JSONClass   string                 `json:"json_class"`
	Attributes  map[string]interface{} `json:"normal"`
	Overrides   map[string]interface{} `json:"override"`
	Defaults    map[string]interface{} `json:"default"`
	RunList     []string               `json:"run_list"`
	Environment string                 `json:"chef_environment"`
}

// chef.GetNodes returns a map of nodes names to the nodes's RESTful URL as well
// as an error indicating if the request was successful or not.
//
// Usage:
//
//     nodes, err := chef.GetNodes()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "node" variable which is a map of
//     // node names to node URLs
//     for node := range nodes {
//         fmt.Println(node)
//      }
func (chef *Chef) GetNodes() (map[string]string, error) {
	resp, err := chef.Get("nodes")
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

// chef.GetNode accepts a string which represents the name of a Chef role and
// returns a chef.Node type representing that role as well as a bool
// indicating whether or not the role was found and an error indicating if the
// request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     node, ok, err := chef.GetNode("neo4j.example.com")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that node!")
//     } else {
//         // do what you please with the "node" variable which is of the
//         // *Chef.Node type
//         fmt.Printf("%#v\n", node)
//     }
func (chef *Chef) GetNode(name string) (*Node, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("nodes/%s", name))
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

	node := new(Node)
	json.Unmarshal(body, node)

	return node, true, nil
}

// CreateNode accepts a node name string and normal/override/default maps for
// attributes, and a []string run list. It will create a new node and error if a
// node already exists. It returns the node name, a success bool, and a potential error.
func (chef *Chef) CreateNode(name, environment string, normal, overrides, defaults map[string]interface{}, run_list []string) (string, bool, error) {
	node := new(PutNode)
	node.Name = name
	node.ChefType = "node"
	node.JSONClass = "Chef::Node"
	node.Attributes = normal
	node.Overrides = overrides
	node.Defaults = defaults
	node.RunList = run_list
	node.Environment = environment

	json_body, err := json.Marshal(node)
	if err != nil {
		return name, false, err
	}

	post_body := bytes.NewReader(json_body)
	resp, err := chef.Post("nodes", nil, post_body)
	if err != nil {
		return name, false, err
	}
	if resp.StatusCode != 201 {
		err = errors.New(fmt.Sprintf("Server returned %s", resp.Status))
		return name, false, err
	}

	// The body contains the URI to the node, not sure if anyone will care.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return name, false, err
	}

	requestError := new(Error)
	json.Unmarshal(body, requestError)

	if len(requestError.Error) != 0 {
		err = errors.New(requestError.Error[0])
		return name, false, err
	}

	return name, true, nil
}

// DeleteNode accepts a node name string and returns the node before it was deleted,
// a true bool when the delete is successful and a false bool with error when its not
func (chef *Chef) DeleteNode(name string) (*Node, bool, error) {
	params := make(map[string]string)

	resp, err := chef.Delete(fmt.Sprintf("nodes/%s", name), params)
	if err != nil {
		return nil, false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	requestError := new(Error)
	json.Unmarshal(body, requestError)

	if len(requestError.Error) != 0 {
		err = errors.New(requestError.Error[0])
		return nil, false, err
	}

	node := new(Node)
	json.Unmarshal(body, &node)
	return node, true, nil
}
