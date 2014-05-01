package chef

import "crypto/rsa"

// chef.KnifeConfig defines all parameters of a knife.rb config file.
// Most attributes can be mapped up with the parameters in knife.rb
// by changing the capitalization to _<lowercase>
// More information can be found at:
// http://docs.opscode.com/config_rb_knife.html
type KnifeConfig struct {
	ChefServerUrl         string
	Host                  string
	Port                  string
	ChefZeroEnabled       bool
	ChefZeroPort          string
	ClientKey             *rsa.PrivateKey
	CookbookCopyright     string
	CookbookEmail         string
	CookbookLicense       string
	CookbookPath          []string
	DataBagEncryptVersion int64
	LocalMode             bool
	NodeName              string
	NoProxy               []string
	SyntaxCheckCachePath  string
	ValidationClientName  string
	ValidationKey         string
	VersionedCookbooks    bool
}

