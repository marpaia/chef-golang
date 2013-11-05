package chef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var testRequiredHeaders []string

var ConfigFilePath = "./TEST_CONFIG.json"

func init() {
	testRequiredHeaders = []string{
		"Accept",
		"X-Ops-Timestamp",
		"X-Ops-Userid",
		"X-Ops-Sign",
		"X-Ops-Content-Hash",
		"X-Ops-Authorization-1",
	}
}

func testConnectionWrapper() *Chef {
	t := new(testing.T)
	chef, err := Connect()
	if err != nil {
		t.Error(err)
	}
	chef.SSLNoVerify = true

	return chef
}

type testConfigFile struct {
	RequiredCookbook struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"required_cookbook"`
	RequiredNode struct {
		Name string `json:"name"`
	} `json:"required_node"`
	RequiredRecipe struct {
		Name string `json:"name"`
	} `json:"required_recipe"`
	RequiredRole struct {
		Name string `json:"name"`
	} `json:"required_role"`
	RequiredClient struct {
		Name string `json:"name"`
	} `json:"required_client"`
	RequiredEnvironment struct {
		Name string `json:"name"`
	} `json:"required_environment"`
	RequiredUser struct {
		Name string `json:"name"`
	} `json:"required_user"`
	RequiredData struct {
		Name string `json:"name"`
	} `json:"required_data"`
	RequiredPrincipal struct {
		Name string `json:"name"`
	} `json:"required_principal"`
	SearchData struct {
		Index string `json:"index"`
		Query string `json:"query"`
	} `json:"search_data"`
	TestCredentials struct {
		Host    string `json:"host"`
		Port    string `json:"port"`
		Version string `json:"version"`
		UserId  string `json:"user_name"`
		Key     string `json:"key"`
	} `json:"test_credentials"`
	KeyPath   string `json:"key_path"`
	KeyString string `json:"key_string"`
}

func testConfig() *testConfigFile {
	t := new(testing.T)
	file, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		t.Error(err)
	}

	var config *testConfigFile
	json.Unmarshal(file, &config)
	if config == nil {
		t.Error("Config is nil")
	}
	return config
}

func TestReadConfig(t *testing.T) {
	_ = testConfig()
}

func TestApiRequest(t *testing.T) {
	chef := testConnectionWrapper()
	method := "GET"
	endpoint := "cookbooks"
	requestURL := fmt.Sprintf("%s/%s", chef.Url, endpoint)
	req, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		t.Error(err)
	}
	chef.apiRequest(req, method, fmt.Sprintf("/%s", endpoint), "")

	for _, value := range testRequiredHeaders {
		if req.Header.Get(value) == "" {
			t.Error("Couldn't find header:", value)
		}
	}
}

func TestResponseBody(t *testing.T) {
	etsy, err := http.Get("https://www.etsy.com/")
	if err != nil {
		t.Error(err)
	}

	bytes, err := responseBody(etsy)
	if err != nil {
		t.Error(err)
	}

	etsyString := "Is code your craft? http://www.etsy.com/careers"
	if !strings.Contains(string(bytes), etsyString) {
		t.Error("Response body didn't return valid string")
	}
}

func TestGet(t *testing.T) {
	c := testConnectionWrapper()
	resp, err := c.Get("cookbooks")
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	cookbooks := map[string]interface{}{}
	json.Unmarshal(body, &cookbooks)
	found := false
	config := testConfig()
	cookbook := config.RequiredCookbook.Name
	for name := range cookbooks {
		if name == cookbook {
			found = true
			break
		}
	}
	if !found {
		t.Error("Required cookbook not found")
	}
}

func TestConnect(t *testing.T) {
	_, err := Connect()
	if err != nil {
		t.Error(err)
	}
}

func TestApiRequestHeaders(t *testing.T) {
	chef := testConnectionWrapper()
	headers := chef.apiRequestHeaders("GET", "/cookbooks", "")
	count := 0
	for _, requiredHeader := range testRequiredHeaders {
		for header := range headers {
			if strings.ToLower(requiredHeader) == strings.ToLower(header) {
				count += 1
				break
			}
		}
	}
	if count != len(testRequiredHeaders) {
		t.Error("apiRequestHeaders didn't return all of testRequiredHeaders")
	}
}

func TestGenerateRequestAuthorization(t *testing.T) {
	chef := testConnectionWrapper()
	auth := chef.generateRequestAuthorization("GET", "/cookbooks", "", "2013-10-27T20:45:25Z")
	if len(auth[0]) != 60 {
		t.Error("Incorrect request authorization string")
	}
}

func TestPrivateEncrypt(t *testing.T) {
	chef := testConnectionWrapper()
	enc, err := chef.privateEncrypt([]byte("encrypt_this"))
	if err != nil {
		t.Error(err)
	}
	if len(enc) != 256 {
		t.Error("Wrong size of encrypted data")
	}
}

func TestBase64BlockEncode(t *testing.T) {
	toEncode := []byte("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz")
	results := base64BlockEncode(toEncode)
	expected := []string{"YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXphYmNkZWZnaGlqa2xtbm9wcXJz", "dHV2d3h5emFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6"}
	if !reflect.DeepEqual(results, expected) {
		t.Error("Results not matching")
	}
}

func TestHashAndBase64(t *testing.T) {
	if len(hashAndBase64("hash_this")) != 28 {
		t.Error("Wrong length for hashAndBase64")
	}
}

func TestDo(t *testing.T) {
	chef := testConnectionWrapper()
	req, err := http.NewRequest("GET", "https://www.etsy.com/", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = chef.Do(req)
	if err != nil {
		t.Error(err)
	}

}

func TestGenerateRequest(t *testing.T) {
	chef := testConnectionWrapper()
	_, err := chef.generateRequest("GET", "cookbooks", nil)
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromString(t *testing.T) {
	config := testConfig()
	_, err := keyFromString([]byte(config.KeyString))
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromFile(t *testing.T) {
	config := testConfig()
	_, err := keyFromFile(config.KeyPath)
	if err != nil {
		t.Error(err)
	}
}

func TestConnectCredentials(t *testing.T) {
	config := testConfig()
	host := config.TestCredentials.Host
	port := config.TestCredentials.Port
	version := config.TestCredentials.Version
	userid := config.TestCredentials.UserId
	key := config.TestCredentials.Key
	_, err := ConnectCredentials(host, port, version, userid, key)
	if err != nil {
		t.Error(err)
	}
}

func TestSplitWhitespace(t *testing.T) {
	str := "c   h   e   f"
	if !reflect.DeepEqual(splitWhitespace(str), []string{"c", "h", "e", "f"}) {
		t.Error("splitWhitespace slices not equal")
	}
}

func TestFilterQuotes(t *testing.T) {
	known := map[string]string{
		`'this`: "this",
		`this'`: "this",
		`"this`: "this",
		`this"`: "this",
	}

	for bad, good := range known {
		if filterQuotes(bad) != good {
			t.Error("filterQuotes didn't produce an expected string")
		}
	}
}
