package chef

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	_, err := ParseConfig("test/support/knife.rb")
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromString(t *testing.T) {
	config := testConfig()
	_, err := KeyFromString([]byte(config.KeyString))
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromFile(t *testing.T) {
	config := testConfig()
	_, err := KeyFromFile(config.KeyPath)
	if err != nil {
		t.Error(err)
	}
}
