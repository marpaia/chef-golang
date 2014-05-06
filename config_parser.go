package chef

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ParseConfig(configfile ...string) (*KnifeConfig, error) {
	var knifeFile string

	if len(configfile) == 1 && configfile[0] != "" {
		knifeFile = configfile[0]
	}

	// If knifeFile is not provided during initialization
	// iterate through the same was as the knife command
	if knifeFile == "" {
		knifeFiles := []string{}

		// Check current directory for .chef
		knifeFiles = append(knifeFiles, filepath.Join(".chef", "knife.rb"))
		// Chef ~/.chef
		knifeFiles = append(knifeFiles, filepath.Join(os.Getenv("HOME"), ".chef", "knife.rb"))

		for _, each := range knifeFiles {
			if _, err := os.Stat(each); err == nil {
				knifeFile = each
				break
			}
		}

		if knifeFile == "" {
			return nil, errors.New("knife.rb configuration file not found")
		}
	}

	file, err := os.Open(knifeFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	config := new(KnifeConfig)
	for scanner.Scan() {
		key := regexp.MustCompile(`\s`).Split(scanner.Text(), -1)
		data := regexp.MustCompile(`^[a-z](.+?)\s+`).Split(scanner.Text(), -1)

		if regexp.MustCompile(`^[a-z]+`).MatchString(key[0]) {
			switch {
			case checkLineKey("chef_server_url", key[0]):
				config.ChefServerUrl = filterQuotes(data[1])
				chefUrl, err := url.Parse(config.ChefServerUrl)
				if err != nil {
					return nil, err
				}
				hostPort := strings.Split(chefUrl.Host, ":")
				if len(hostPort) == 2 {
					config.Host = hostPort[0]
					config.Port = hostPort[1]
				} else if len(hostPort) == 1 {
					config.Host = hostPort[0]
					switch chefUrl.Scheme {
					case "http":
						config.Port = "80"
					case "https":
						config.Port = "443"
					default:
						return nil, errors.New("Invalid http scheme")
					}

				} else {
					return nil, errors.New("Invalid host format")
				}

			case checkLineKey("chef_zero[:enabled]", key[0]):
				config.ChefZeroEnabled, _ = strconv.ParseBool(filterQuotes(data[1]))
			case checkLineKey("chef_zero[:port]", key[0]):
				config.ChefZeroPort = filterQuotes(data[1])
			case checkLineKey("client_key", key[0]):
				key, err := KeyFromFile(filterQuotes(data[1]))
				if err != nil {
					return nil, err
				}
				config.ClientKey = key
			case checkLineKey("cookbook_copyright", key[0]):
				config.CookbookCopyright = filterQuotes(data[1])
			case checkLineKey("cookbook_email", key[0]):
				config.CookbookEmail = filterQuotes(data[1])
			case checkLineKey("cookbook_license", key[0]):
				config.CookbookLicense = filterQuotes(data[1])
			case checkLineKey("cookbook_path", key[0]):
				re1 := regexp.MustCompile(`(\n|,|"|\[|\]|^cookbook_path\s)`)
				cookbookPaths := strings.Fields(re1.ReplaceAllString(scanner.Text(), ``))
				for !strings.Contains(scanner.Text(), `]`) {
					for _, item := range strings.Fields(re1.ReplaceAllString(scanner.Text(), ``)) {
						cookbookPaths = append(cookbookPaths, item)
					}
					scanner.Scan()
				}
				config.CookbookPath = cookbookPaths
			case checkLineKey("data_bag_encrypt_version", key[0]):
				config.DataBagEncryptVersion, _ = strconv.ParseInt(filterQuotes(data[1]), 0, 0)
			case checkLineKey("local_mode", key[0]):
				config.LocalMode, _ = strconv.ParseBool(filterQuotes(data[1]))
			case checkLineKey("node_name", key[0]):
				config.NodeName = filterQuotes(data[1])
			case checkLineKey("no_proxy", key[0]):
				config.NoProxy = regexp.MustCompile(`,\s+`).Split(filterQuotes(data[1]), -1)
			case checkLineKey("syntax_check_cache_path", key[0]):
				config.SyntaxCheckCachePath = filterQuotes(data[1])
			case checkLineKey("validation_client_name", key[0]):
				config.ValidationClientName = filterQuotes(data[1])
			case checkLineKey("validation_key", key[0]):
				config.ValidationKey = filterQuotes(data[1])
			case checkLineKey("versioned_cookbooks", key[0]):
				config.VersionedCookbooks, _ = strconv.ParseBool(filterQuotes(data[1]))
			}
		}
	}
	return config, nil
}

// given a string to compare against and an input string, checkLineKey will
// return true if the strings are equal
func checkLineKey(k string, s string) bool {

	if k == s {
		return true
	} else {
		return false
	}
}

// Given a string with multiple consecutive spaces, splitWhitespace returns a
// slice of strings which represent the given string split by \s characters with
// all duplicates removed
func splitWhitespace(s string) []string {
	re := regexp.MustCompile(`\s+`)
	return strings.Split(re.ReplaceAllString(s, `\s`), `\s`)
}

// filterQuotes returns a string with surrounding quotes filtered
func filterQuotes(s string) string {
	re1 := regexp.MustCompile(`^(\'|\")`)
	re2 := regexp.MustCompile(`(\'|\")$`)
	return re2.ReplaceAllString(re1.ReplaceAllString(s, ``), ``)
}

// KeyFromFile reads an RSA private key given a filepath
func KeyFromFile(filename string) (*rsa.PrivateKey, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return KeyFromString(content)
}

// KeyFromString parses an RSA private key from a string
func KeyFromString(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("block size invalid for '%s'", string(key))
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}
