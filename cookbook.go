package chef

import (
	"encoding/json"
	"fmt"
	"strings"
)

// chef.Cookbook defines the relevant parameters of a Chef cookbook. This
// includes the RESTful URL of a cookbook and a slice of all of the cookbooks
// versions. The versions each have two attributes: Url, which represents the
// RESTful URL of the cookbook version and Version, which represents the version
// number (identifier) of the cookbook version
type Cookbook struct {
	Url      string `json:"url"`
	Versions []struct {
		Url     string `json:"url"`
		Version string `json:"version"`
	} `json:"versions"`
}

// chef.CookbookVersion defines the relevant parameters of a specific Chef
// cookbook version. This includes, but is not limited to, information about
// recipes, files, etc, various pieces of metadata about the cookbook at that
// point in time, such as the name of the cookbook, the description, the
// license, etc.
type CookbookVersion struct {
	Recipes []struct {
		CookbookItem
	} `json:"recipes"`
	Files []struct {
		CookbookItem
	} `json:"files"`
	RootFiles []struct {
		CookbookItem
	} `json:"root_file"`
	Metadata struct {
		Name            string            `json:"name"`
		Description     string            `json:"description"`
		LongDescription string            `json:"long_description"`
		Maintainer      string            `json:"maintainer"`
		MaintainerEmail string            `json:"maintainer_email"`
		License         string            `json:"license"`
		Providing       map[string]string `json:"providing"`
		Dependencies    map[string]string `json:dependencies`
	} `json:"metadata"`
	Name      string `json:"cookbook_name"`
	Version   string `json:"version"`
	FullName  string `json:"name"`
	Frozen    bool   `json:"frozen?"`
	ChefType  string `json:"chef_type"`
	JSONClass string `json:"json_class"`
}

// chef.CookbookItem defines the relevant parameters of various items that are
// found in a chef Cookbook such as the name, checksum, etc. This type is
// embedded in the chef.CookVersion type to reduce code repetition.
type CookbookItem struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Checksum    string `json:"checksum"`
	Specificity string `json:"specificity"`
	Url         string `json:"url"`
}

// chef.GetCookbooks returns a map of cookbook names to a pointer to the
// chef.Cookbook type as well as an error indicating if the request was
// successful or not.
//
// Usgae:
//
//     cookbooks, err := chef.GetCookbooks()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "cookbooks" variable which is a map of
//     // cookbook names to chef.Cookbook types
//     for name, cookbook := range cookbooks {
//         fmt.Println(name, cookbook.Version[0])
//      }
func (chef *Chef) GetCookbooks() (map[string]*Cookbook, error) {
	resp, err := chef.Get("cookbooks")
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

// chef.GetCookbook returns a pointer to the chef.Cookbook type for a given
// string that represents a cookbook name. It also returns a bool indicating
// whether or not the client was found and an error indicating if the request
// failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     cookbook, ok, err := chef.GetCookbook("apache")
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
func (chef *Chef) GetCookbook(name string) (*Cookbook, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("cookbooks/%s", name))
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

	return cookbook[name], true, nil
}

// chef.GetCookbookVersion returns a pointer to the chef.CookbookVersion type
// for a given string that represents a cookbook version. It also returns a bool
// indicating whether or not the client was found and an error indicating if
// the request failed or not.
//
// Note that if the request is successful but no such client existed, the error
// return value will be nil but the bool will be false.
//
// Usage:
//
//     cookbook, ok, err := chef.GetCookbookVersion("apache", "1.0.0")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     if !ok {
//         fmt.Println("Couldn't find that cookbook version!")
//     } else {
//         // do what you please with the "cookbook" variable which is of the
//         // *Chef.CookbookVersion type
//         fmt.Printf("%#v\n", cookbook)
//     }
func (chef *Chef) GetCookbookVersion(name, version string) (*CookbookVersion, bool, error) {
	resp, err := chef.Get(fmt.Sprintf("cookbooks/%s/%s", name, version))
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
	cookbook := new(CookbookVersion)
	json.Unmarshal(body, &cookbook)
	return cookbook, true, nil
}
