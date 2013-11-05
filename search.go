package chef

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// chef.SearchParams represents the necessary parameters of a Chef search query
type SearchParams struct {
	Index string
	Query string
	Sort  string
	Rows  int
	Start int
	chef  *Chef
}

// chef.SearchResults represents the results of a Chef search query
type SearchResults struct {
	Total int           `json:"total"`
	Start int           `json:"start"`
	Rows  []interface{} `json:"rows"`
}

// chef.GetSearchIndexes returns a map of search indexes to the indexes RESTful
// URL as well as an error indicating if the request was successful or not.
//
// Usage:
//
//     indexes, err := chef.GetSearchIndexes()
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "indexes" variable which is a map of
//     // index names to index URLs
//     for index := range indexes {
//         fmt.Println(index)
//      }
func (chef *Chef) GetSearchIndexes() (map[string]string, error) {
	resp, err := chef.Get("search")
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	results := map[string]string{}
	json.Unmarshal(body, &results)

	return results, nil
}

// chef.NewSearchQuery accepts an index and a query and returns a struct which
// represents the appropriate search parameters. This is mostly just used by the
// chef.Search method, but you can call it yourself if you'd like some more
// control over your parameters
func (chef *Chef) NewSearchQuery(index, query string) *SearchParams {
	params := new(SearchParams)
	params.Index = index
	params.Query = query
	params.Rows = -1
	params.Start = -1
	params.chef = chef
	return params
}

// chef.Search accepts an index and a query and returns a *Chef.Search results
// type as well as an error indicating if the request was successful or not
//
// Usage:
//
//     results, err := chef.Search("nodes", "hostname:memcached*")
//     if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//     }
//     // do what you please with the "results" variable which is of the type
//     // *Chef.SearchResults
//     fmt.Println(results)
func (chef *Chef) Search(index, query string) (*SearchResults, error) {
	return chef.NewSearchQuery(index, query).Execute()
}

// chef.SearchWithParams is similar to chef.Search, but you can define
// additional Chef search parameters
func (chef *Chef) SearchWithParams(index, query string, params map[string]interface{}) (*SearchResults, error) {
	searchParams := chef.NewSearchQuery(index, query)
	if params["rows"] != nil {
		searchParams.Rows = params["rows"].(int)
	}
	if params["start"] != nil {
		searchParams.Start = params["start"].(int)
	}
	if params["sort"] != nil {
		searchParams.Sort = params["sort"].(string)
	}
	return searchParams.Execute()
}

// chef.Execute is a method on the chef.SearchParams type that executes a given
// search that has a given set of paramters. This is mostly used by the
// chef.Search method, but you can call it yourself if you'd like some more
// control over your parameters
func (search *SearchParams) Execute() (*SearchResults, error) {
	params := map[string]string{
		"q": search.Query,
	}
	if search.Rows != -1 {
		params["rows"] = strconv.Itoa(search.Rows)
	}
	if search.Start != -1 {
		params["start"] = strconv.Itoa(search.Start)
	}
	if search.Sort != "" {
		params["sort"] = search.Sort
	}
	resp, err := search.chef.GetWithParams(fmt.Sprintf("search/%s", search.Index), params)
	if err != nil {
		return nil, err
	}
	body, err := responseBody(resp)
	if err != nil {
		return nil, err
	}

	results := new(SearchResults)
	json.Unmarshal(body, results)

	return results, nil
}
