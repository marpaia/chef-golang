package chef

import (
	"testing"
)

func TestGetSearchIndexes(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	results, err := chef.GetSearchIndexes()
	if err != nil {
		t.Error(err)
	}
	found := false
	for index := range results {
		if index == config.SearchData.Index {
			found = true
		}
	}
	if !found {
		t.Error("Couldn't find required index")
	}
}

func TestSearch(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	_, err := chef.Search(config.SearchData.Index, config.SearchData.Query)
	if err != nil {
		t.Error(err)
	}
}

func TestSearchWithParams(t *testing.T) {}

func TestNewSearchQuery(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	query := chef.NewSearchQuery(config.SearchData.Index, config.SearchData.Query)
	if query.Index != config.SearchData.Index {
		t.Error("Search index isn't correctly set")
	}
	if query.Query != config.SearchData.Query {
		t.Error("Search query isn't correctly set")
	}
}
