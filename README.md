chef-golang [![Build Status](https://drone.io/github.com/marpaia/chef-golang/status.png)](https://drone.io/github.com/marpaia/chef-golang/latest)
===========
This is a Go API client for Opscode's Chef.

## Installation

Use `go get` to install chef-golang:
```
go get github.com/marpaia/chef-golang
```

## Unit tests
There is a helper script in `test/support/start_server.sh` that will setup a goiardi server primed
with unit-testing data. there is also an example knife config and a script named `test/support/chef_config.sh`.

If you don't have a local `~/.chef/knife.rb` you can simply run these 3 commands:
```
./test/support/start_server.sh
./test/support/chef_config.sh
go test -v
```
to shutdown the server run `test/support/stop_server.sh`

The [goiardi](https://github.com/ctdk/goiardi) Server is listening on port `8443` and the keys and config are copied to /tmp/goiardi. If you want to setup your own `knife.rb` and not use the provided one.

## External dependencies

This project has no external dependencies other than the Go standard library.

## Documentation

Like most every other Golang project, this projects documentation can be found
on godoc at [godoc.org/github.com/marpaia/chef-golang](http://godoc.org/github.com/marpaia/chef-golang).

## Examples
```go
package main

import (
    "fmt"
    "os"

    "github.com/marpaia/chef-golang"
)

var findNode = "neo4j.example.org"
var findCookbook = "neo4j"
var findRole = "Neo4j"

func main() {
    c, err := chef.Connect()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    c.SSLNoVerify = true

    // Print detailed information about a specific node
    node, ok, err := c.GetNode(findNode)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    if !ok {
        fmt.Println("\nCouldn't find that node!")
    } else {
        for i := 0; i < 80; i++ {
            fmt.Print("=")
        }

        fmt.Println("\nSystem info:", node.Name, "\n")
        fmt.Println("  [+] IP Address:", node.Info.IPAddress)
        fmt.Println("  [+] MAC Address:", node.Info.MACAddress)
        fmt.Println("  [+] Operating System:", node.Info.Platform)

        fmt.Println("\n  [+] Filesystem Info")
        for partition, info := range node.Info.Filesystem {
            if info.PercentUsed != "" {
                fmt.Println("    - ", partition, "is", info.PercentUsed, "utilized")
            }
        }

        fmt.Println("\n  [+] Roles")
        for _, role := range node.Info.Roles {
            fmt.Println("    - ", role)
        }

        fmt.Println()
        for i := 0; i < 80; i++ {
            fmt.Print("=")
        }
    }

    // Print detailed information about a specific cookbook
    cookbook, ok, err := c.GetCookbook(findCookbook)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    if !ok {
        fmt.Println("\nCouldn't find that cookbook!")
    } else {
        fmt.Println("\n\nCookbook info:", findCookbook)
        for _, version := range cookbook.Versions {
            currentVersion, ok, err := c.GetCookbookVersion(findCookbook, version.Version)
            if err != nil {
                fmt.Println("Error:", err)
                os.Exit(1)
            }
            if ok {
                if len(currentVersion.Files) > 0 {
                    fmt.Println("\n  [+]", findCookbook, currentVersion.Version, "Cookbook Files")
                    for _, cookbookFile := range currentVersion.Files {
                        fmt.Println("    - ", cookbookFile.Name)
                    }
                }
            }
        }

        fmt.Println()
        for i := 0; i < 80; i++ {
            fmt.Print("=")
        }
    }

    // Print detailed information about a specific role
    role, ok, err := c.GetRole(findRole)
    if err != nil {
        fmt.Println("Error:", err)
    }
    if !ok {
        fmt.Println("\nCouldn't find that role!")
    } else {
        fmt.Println("\n\nRole information:", role.Name)
        fmt.Println("\n[+] Runlist")
        for _, recipe := range role.RunList {
            fmt.Println("  - ", recipe)
        }
    }
}
```

Which will output something like this:

```
================================================================================
System info: neo4j.example.org

  [+] IP Address: 10.100.1.2
  [+] MAC Address: AA:BB:CC:DD:EE:FF
  [+] Operating System: centos

  [+] Filesystem Info
    -  /dev/vda is 46% utilized
    -  tmpfs is 1% utilized

  [+] Roles
    -  Base
    -  Neo4j

================================================================================

Cookbook info: neo4j

  [+] neo4j 0.1.6 Cookbook Files
    -  neo4j-server.properties
    -  neo4j-service
    -  neo4j-wrapper.conf

  [+] neo4j 0.1.5 Cookbook Files
    -  neo4j-service
    -  neo4j-wrapper.conf

  [+] neo4j 0.1.4 Cookbook Files
    -  neo4j-service

  [+] neo4j 0.1.3 Cookbook Files
    -  neo4j-service

  [+] neo4j 0.1.2 Cookbook Files
    -  neo4j-service

================================================================================

Role information: Neo4j

[+] Runlist
  -  recipe[yum]
  -  recipe[ldap]
  -  recipe[system]
  -  recipe[neo4j::server]
  -  recipe[neo4j::web-ui]
```

## Contributing

Please contribute and help improve this project!

- Fork the repo
- Make sure the tests pass
- Improve the code
- Make sure your feature has test coverage
- Make sure the tests pass
- Submit a pull request
