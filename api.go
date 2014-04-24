package chef

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Chef is the type that contains all of the relevant information about a Chef
// server connection
type Chef struct {
	Host        string
	Url         string
	Port        string
	Version     string
	Key         *rsa.PrivateKey
	UserId      string
	SSLNoVerify bool
}

// Connect looks for knife/chef configuration files and gather connection info
// automagically
func Connect() (*Chef, error) {
	knifeFiles := []string{}
	homedir := os.Getenv("HOME")
	if homedir != "" {
		knifeFiles = append(knifeFiles, filepath.Join(homedir, ".chef/knife.rb"))
	}
	knifeFiles = append(knifeFiles, "/etc/chef/client.rb")
	var knifeFile string
	for _, each := range knifeFiles {
		if _, err := os.Stat(each); err == nil {
			knifeFile = each
			break
		}
	}

	if knifeFile == "" {
		return nil, errors.New("Configuration file not found")
	}

	file, err := os.Open(knifeFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	chef := new(Chef)
	for scanner.Scan() {
		split := splitWhitespace(scanner.Text())
		if len(split) == 2 {
			switch split[0] {
			case "node_name":
				chef.UserId = filterQuotes(split[1])
			case "client_key":
				key, err := keyFromFile(filterQuotes(split[1]))
				if err != nil {
					return nil, err
				}
				chef.Key = key
			case "chef_server_url":
				parsedUrl := filterQuotes(split[1])
				chef.Url = parsedUrl
				chefUrl, err := url.Parse(parsedUrl)
				if err != nil {
					return nil, err
				}
				hostPort := strings.Split(chefUrl.Host, ":")
				if len(hostPort) == 2 {
					chef.Host = hostPort[0]
					chef.Port = hostPort[1]
				} else if len(hostPort) == 1 {
					chef.Host = hostPort[0]
					switch chefUrl.Scheme {
					case "http":
						chef.Port = "80"
					case "https":
						chef.Port = "443"
					default:
						return nil, errors.New("Invalid http scheme")
					}

				} else {
					return nil, errors.New("Invalid host format")
				}
			}
		}
	}

	if chef.Key == nil {
		return nil, errors.New("missing 'client_key' in knife.rb")
	}

	return chef, nil
}

// filterQuotes returns a string with surrounding quotes filtered
func filterQuotes(s string) string {
	re1 := regexp.MustCompile(`^(\'|\")`)
	re2 := regexp.MustCompile(`(\'|\")$`)
	return re2.ReplaceAllString(re1.ReplaceAllString(s, ``), ``)
}

// Given a string with multiple consecutive spaces, splitWhitespace returns a
// slice of strings which represent the given string split by \s characters with
// all duplicates removed
func splitWhitespace(s string) []string {
	re := regexp.MustCompile(`\s+`)
	return strings.Split(re.ReplaceAllString(s, `\s`), `\s`)
}

// Given the appropriate connection parameters, ConnectChef returns a pointer to
// a Chef type so that you can call request methods on it
func ConnectCredentials(host, port, version, userid, key string) (*Chef, error) {
	chef := new(Chef)
	chef.Host = host
	chef.Port = port
	chef.Version = version
	chef.UserId = userid

	var url string
	switch chef.Port {
	case "443":
		url = fmt.Sprintf("https://%s", chef.Host)
	case "80":
		url = fmt.Sprintf("http://%s", chef.Host)
	default:
		url = fmt.Sprintf("%s:%d", chef.Host, chef.Port)
	}

	chef.Url = url

	var rsaKey *rsa.PrivateKey
	var err error

	if strings.Contains(key, "-----BEGIN RSA PRIVATE KEY-----") {
		rsaKey, err = keyFromString([]byte(key))
	} else {
		rsaKey, err = keyFromFile(key)
	}
	if err != nil {
		return nil, err
	}

	chef.Key = rsaKey
	return chef, nil
}

func ConnectUrl(chefServerUrl, version, userid, key string) (*Chef, error) {
	chef := new(Chef)
	chef.Version = version
	chef.UserId = userid
	chef.Url = chefServerUrl

	var rsaKey *rsa.PrivateKey
	var err error

	if strings.Contains(key, "-----BEGIN RSA PRIVATE KEY-----") {
		rsaKey, err = keyFromString([]byte(key))
	} else {
		rsaKey, err = keyFromFile(key)
	}
	if err != nil {
		return nil, err
	}

	chef.Key = rsaKey

	return chef, nil
}

// keyFromFile reads an RSA private key given a filepath
func keyFromFile(filename string) (*rsa.PrivateKey, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return keyFromString(content)
}

// keyFromString parses an RSA private key from a string
func keyFromString(key []byte) (*rsa.PrivateKey, error) {
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

// Get makes an authenticated HTTP request to the Chef server for the supplied
// endpoint
func (chef *Chef) Get(endpoint string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", chef.requestUrl(endpoint), nil)
	return chef.makeRequest(request)
}

// GetWithParams makes an authenticated HTTP request to the Chef server for the
// supplied endpoint and also includes GET query string parameters
func (chef *Chef) GetWithParams(endpoint string, params map[string]string) (*http.Response, error) {
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}
	request, _ := http.NewRequest("GET", chef.requestUrl(endpoint), nil)
	request.Form = values
	return chef.makeRequest(request)
}

// PostForm makes an authenticated POST request to the Chef server With params for the supplied
// endpoint
func (chef *Chef) PostForm(endpoint string, params map[string]string) (*http.Response, error) {
	//TODO: Finish This
	request, _ := http.NewRequest("POST", chef.requestUrl(endpoint), nil)
	return chef.makeRequest(request)
}

// Post  post an io object to the chef server this uses standard http.Post, make sure you close
// your io or send io.Closer
func (chef *Chef) Post(endpoint string, body io.Reader) (*http.Response, error) {
	//TODO: Finish this
	request, _ := http.NewRequest("POST", chef.requestUrl(endpoint), nil)
	return chef.makeRequest(request)
}

// Put makes an authenticated PUT request to the Chef server for the supplied
// endpoint
func (chef *Chef) Put(endpoint string, params map[string]string) (*http.Response, error) {
	//TODO: Finish this
	request, _ := http.NewRequest("PUT", chef.requestUrl(endpoint), nil)
	return chef.makeRequest(request)
}

// Delete makes an authenticated DELETE request to the Chef server for the
// supplied endpoint
func (chef *Chef) Delete(endpoint string, params map[string]string) (*http.Response, error) {
	//TODO: Finish this
	request, _ := http.NewRequest("DELETE", chef.requestUrl(endpoint), nil)
	return chef.makeRequest(request)
}

// requestUrl generate the requestUrl from supplied endpoint
func (chef *Chef) requestUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", chef.Url, endpoint)
}

// sign request takes a pointer to a http.Request and sets up the authentication
// headers so that it can be used against chef-server
func (chef *Chef) SignRequest(request *http.Request) {
	//TODO
}

// makeRequest
// Take a request object, Setup Auth headers and Send it to the server
func (chef *Chef) makeRequest(request *http.Request) (*http.Response, error) {
	chef.SignRequest(request)

	var client *http.Client
	if chef.SSLNoVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	return client.Do(request)
}

// base64BlockEncode takes a byte slice and breaks it up into a slice of strings
// where each string is 60 characters long
func base64BlockEncode(content []byte) []string {
	resultString := base64.StdEncoding.EncodeToString(content)
	var resultSlice []string
	index := 0

	for i := 0; i < len(resultString)/60; i += 1 {
		resultSlice = append(resultSlice, resultString[index:index+60])
		index += 60
	}

	if len(resultString)%60 != 0 {
		resultSlice = append(resultSlice, resultString[index:])
	}

	return resultSlice
}

// hashAndBase64 takes a string a returns a base64 representation of the hash of
// the string in \n seperated 60 character long blocks (don't ask, it's a Chef
// thing apparently)
func hashAndBase64(content io.Reader) string {
	if content == nil {
		content = bytes.NewBufferString("")
	}
	h := sha1.New()
	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	io.WriteString(h, buf.String())
	hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hashed
}

// getTimestamp returns an ISO-8601 formatted timestamp of the current time in
// UTC
func getTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// privateEncrypt implements OpenSSL's RSA_private_encrypt function
func (chef *Chef) privateEncrypt(data []byte) (enc []byte, err error) {
	k := (chef.Key.N.BitLen() + 7) / 8
	tLen := len(data)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	if tLen > k-11 {
		err = errors.New("Data too long")
		return
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], data)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(chef.Key.N) > 0 {
		err = nil
		return
	}
	var m *big.Int
	var ir *big.Int
	if chef.Key.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, chef.Key.D, chef.Key.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, chef.Key.Precomputed.Dp, chef.Key.Primes[0])
		m2 := new(big.Int).Exp(c, chef.Key.Precomputed.Dq, chef.Key.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, chef.Key.Primes[0])
		}
		m.Mul(m, chef.Key.Precomputed.Qinv)
		m.Mod(m, chef.Key.Primes[0])
		m.Mul(m, chef.Key.Primes[1])
		m.Add(m, m2)

		for i, values := range chef.Key.Precomputed.CRTValues {
			prime := chef.Key.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, chef.Key.N)
	}
	enc = m.Bytes()
	return
}

// generateRequestAuthorization returns a string slice of the Chef server
// authorization headers
func (chef *Chef) generateRequestAuthorization(request *http.Request, timestamp string) ([]string, error) {
	path := &request.URL.Path
	var content string
	content += fmt.Sprintf("Method:%s\n", &request.Method)
	content += fmt.Sprintf("Hashed Path:%s\n", hashAndBase64(bytes.NewBufferString(*path)))
	content += fmt.Sprintf("X-Ops-Content-Hash:%s\n", hashAndBase64(request.Body))
	content += fmt.Sprintf("X-Ops-Timestamp:%s\n", timestamp)
	content += fmt.Sprintf("X-Ops-UserId:%s", chef.UserId)
	signature, err := chef.privateEncrypt([]byte(content))
	if err != nil {
		return nil, err
	}
	return base64BlockEncode([]byte(string(signature))), nil
}

// apiRequestHeaders generates a map of all of the request headers that a
// request to the Chef API will need
func (chef *Chef) apiRequestHeaders(request *http.Request) error {
	timestamp := getTimestamp()
	headers := map[string]string{
		"accept":             "application/json",
		"x-chef-version":     chef.Version,
		"x-ops-timestamp":    timestamp,
		"x-ops-userid":       chef.UserId,
		"x-ops-sign":         "version=1.0",
		"x-ops-content-hash": hashAndBase64(request.Body),
	}

	auths, err := chef.generateRequestAuthorization(request, timestamp)
	if err != nil {
		return err
	}
	for index, value := range auths {
		headers[fmt.Sprintf("X-Ops-Authorization-%d", index+1)] = string(value)
	}

	// Write headers to the request
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	//DBG
	spew.Dump(&request)
	return nil
}

// Given an http response object, responseBody returns the response body
func responseBody(resp *http.Response) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}
