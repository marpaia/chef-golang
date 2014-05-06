log_level                :info
log_location             STDOUT
node_name                'admin'
client_key               '/tmp/goiardi/admin.pem'
validation_client_name   'chef-validator'
validation_key           '/tmp/goiardi/chef-validator.pem'
chef_server_url          'http://127.0.0.1:8443'
cookbook_path [
  "/var/chef/cookbooks",
  "/var/chef/site-cookbooks"
]
chef_zero[:enabled] true
chef_zero[:port] 8889
data_bag_encrypt_version 2
local_mode true
no_proxy "localhost, 10.*, *.example.com, *.dev.example.com"
versioned_cookbooks true
cookbook_license "chef-golang-license"
cookbook_email "chef-golang@chef-golang.github.com"
cookbook_copyright "chef-golang-copyright"

