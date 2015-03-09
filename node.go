package chef

import (
	"encoding/json"
	"fmt"
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
		Languages map[string]interface{} `json:"languages"`
		Kernel    struct {
			Name    string                            `json:"name"`
			Release string                            `json:"release"`
			Version string                            `json:"version"`
			Machine string                            `json:"machine"`
			Modules map[string]map[string]interface{} `json:"modules"`
		} `json:"kernel"`
		OS        string `json:"os"`
		OSVersion string `json:"os_version"`
		Hostname  string `json:"hostname"`
		FQDN      string `json:"fqdn"`
		Domain    string `json:"domain"`
		Network   struct {
			Interfaces map[string]struct {
				Type          string `json:"type"`
				Encapsulation string `json:"encapsulation"`
				Addresses     map[string]struct {
					Family    string `json:"family"`
					Broadcast string `json:"broadcast"`
					Netmast   string `json:"netmast"`
				} `json:"addresses"`
				Routes []struct {
					Destination string `json:"destination"`
					Family      string `json:"family"`
					Metric      string `json:"metric"`
				} `json:"routes"`
				State string            `json:"state"`
				Flags []string          `json:"flags"`
				MTU   string            `json:"mtu"`
				Arp   map[string]string `json:"arp"`
			} `json:"interfaces"`
			DefaultInterface string `json:"default_interface"`
			DefaultGateway   string `json:"default_gateway"`
		} `json:"network"`
		IPAddress       string                       `json:"ipaddress"`
		MACAddress      string                       `json:"macaddress"`
		ChefPackages    map[string]map[string]string `json:"chef_packages"`
		Keys            map[string]map[string]string `json:"keys"`
		Platform        string                       `json:"platform"`
		PlatformVersion string                       `json:"platform_version"`
		PlatformFamily  string                       `json:"platform_family"`
		CPU             map[string]interface{}       `json:"cpu"`
		Filesystem      map[string]struct {
			KBSize       interface{} `json:"ks_size"`
			KBUsed       interface{} `json:"ks_used"`
			KBavailable  interface{} `json:"ks_available"`
			PercentUsed  interface{} `json:"percent_used"`
			Mount        string      `json:"mount"`
			FSType       string      `json:"fs_type"`
			MountOptions []string    `json:"mount_options"`
		} `json:"filesystem"`
		Memory          map[string]interface{} `json:"memory"`
		UptimeSeconds   int                    `json:"uptime_seconds"`
		Uptime          string                 `json:"uptime"`
		IdletimeSeconds int                    `json:"idletime_seconds"`
		Idletime        string                 `json:"idletime"`
		BlockDevice     map[string]interface{} `json:"block_device"`
		Recipes         []string               `json:"recipes"`
		Roles           []string               `json:"roles"`
	} `json:"automatic"`
	Normal  map[string]interface{} `json:"normal"`
	Default map[string]interface{} `json:"default"`
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
// returns a chef.Environment type representing that role as well as a bool
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
