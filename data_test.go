package chef

import (
	"testing"
)

func TestGetData(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	data, err := chef.GetData()
	if err != nil {
		t.Error(err)
	}
	found := false
	for point := range data {
		if point == config.RequiredData.Name {
			found = true
			break
		}
	}
	if !found {
		t.Error("Couldn't find required data in data results")
	}
}

func TestGetDataByName(t *testing.T) {
	chef := testConnectionWrapper()
	config := testConfig()
	_, ok, err := chef.GetDataByName(config.RequiredData.Name)
	if !ok {
		t.Error("Couldn't find required data")
	}
	if err != nil {
		t.Error(err)
	}
}
