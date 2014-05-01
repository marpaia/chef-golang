package chef

import (
	"testing"
	"strconv"
)

func TestParseConfig(t *testing.T) {
	c, err := ParseConfig("test/support/knife.rb")
	//c, err := ParseConfig("/Users/bb016494/.chef/knife.rb")
	if err != nil {
		t.Error(err)
	}

	// Verify we've read in the config file correctly
	if c.ChefServerUrl != "http://127.0.0.1:8443" {
		t.Error("Incorrect chef_server_url: " + c.ChefServerUrl)
	}

	if c.Host != "127.0.0.1" {
		t.Error("Incorrect chef server host: " + c.Host)
	}

	if c.Port != "8443" {
		t.Error("Incorrect chef server port: " + c.Port)
	}

	if c.ChefZeroEnabled != true {
		t.Error("Incorrect chef_zero[:enabled]: " + strconv.FormatBool(c.ChefZeroEnabled))
	}

	if c.ChefZeroPort != "8889" {
		t.Error("Incorrect chef_zero[:port]: " + c.ChefZeroPort)
	}

	if c.CookbookCopyright != "chef-golang-copyright" {
		t.Error("Incorrect cookbook_copyright: " + c.CookbookCopyright)
	}

	if c.CookbookEmail != "chef-golang@chef-golang.github.com" {
		t.Error("Incorrect cookbook_email: " + c.CookbookEmail)
	}

	if c.CookbookLicense != "chef-golang-license" {
		t.Error("Incorrect cookbook_license: " + c.CookbookLicense)
	}

	if len(c.CookbookPath) == 2 {
		if c.CookbookPath[0] != `/var/chef/cookbooks` {
			t.Error("Incorrect cookbook_path: " + c.CookbookPath[0])
		}
		if c.CookbookPath[1] != `/var/chef/site-cookbooks` {
			t.Error("Incorrect cookbook_path: " + c.CookbookPath[1])
		}
	} else {
		t.Error("Incorrect cookbook_path")
	}

	if c.DataBagEncryptVersion != 2 {
		t.Error("Incorrect data_bag_encrypt_version: " + strconv.FormatInt(c.DataBagEncryptVersion, 10))
	}

	if c.LocalMode != true {
		t.Error("Incorrect local_mode: " + strconv.FormatBool(c.LocalMode))
	}

	if c.NodeName != "admin" {
		t.Error("Incorrect node_name: " + c.NodeName)
	}

	//t.Error(c.NoProxy)
	if len(c.NoProxy) == 4 {
		if c.NoProxy[0] != "localhost" {
			t.Error("Incorrect no_proxy: " + c.NoProxy[0] )
		}
		if c.NoProxy[1] != "10.*" {
			t.Error("Incorrect no_proxy: " + c.NoProxy[1] )
		}
		if c.NoProxy[2] != "*.example.com" {
			t.Error("Incorrect no_proxy: " + c.NoProxy[2] )
		}
		if c.NoProxy[3] != "*.dev.example.com" {
			t.Error("Incorrect no_proxy: " + c.NoProxy[3] )
		}
	} else {
		t.Error("Incorrect no_proxy")
	}

	if c.SyntaxCheckCachePath != "" {
		t.Error("Incorrect syntax_check_cache_path: " + c.SyntaxCheckCachePath)
	}

	if c.ValidationClientName != "chef-validator" {
		t.Error("Incorrect chef-validator: " + c.ValidationClientName)
	}

	if c.ValidationKey != "/tmp/goiardi/chef-validator.pem" {
		t.Error("Incorrect validation_key: " + c.ValidationKey)
	}

	if c.VersionedCookbooks != true {
		t.Error("Incorrect versioned_cookbooks: " + strconv.FormatBool(c.VersionedCookbooks))
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

