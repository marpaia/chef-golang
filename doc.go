// This is a Go client for Opscode's Chef.
//
// Obviously many of the unit tests require a functioning Chef installation in
// order to verify the results of API requests. Edit the `TEST_CONFIG.json` file
// with the appropriate endpoints and information, and run `go test -v`.
//
// This project has no external dependencies other than the Go standard library.
//
// Like most every other Golang project, this projects documentation can be found
// on godoc at http://godoc.org/github.com/marpaia/chef-golang
//
// The following is an example of a simple Go tool that displays a small amount
// of information about a few different pieces of a Chef infrastructure:
//
//     package main
//
//     import (
//         "fmt"
//         "github.com/marpaia/chef-golang"
//         "os"
//     )
//
//     var findNode = "neo4j.example.org"
//     var findCookbook = "neo4j"
//     var findRole = "Neo4j"
//
//     func main() {
//         c, err := chef.Connect()
//         if err != nil {
//             fmt.Println("Error:", err)
//             os.Exit(1)
//         }
//         c.SSLNoVerify = true
//
//         // Print detailed information about a specific node
//         node, ok, err := c.GetNode(findNode)
//         if err != nil {
//             fmt.Println("Error:", err)
//             os.Exit(1)
//         }
//         if !ok {
//             fmt.Println("\nCouldn't find that node!")
//         } else {
//             for i := 0; i < 80; i++ {
//                 fmt.Print("=")
//             }
//
//             fmt.Println("\nSystem info:", node.Name, "\n")
//             fmt.Println("  [+] IP Address:", node.Info.IPAddress)
//             fmt.Println("  [+] MAC Address:", node.Info.MACAddress)
//             fmt.Println("  [+] Operating System:", node.Info.Platform)
//
//             fmt.Println("\n  [+] Filesystem Info")
//             for partition, info := range node.Info.Filesystem {
//                 if info.PercentUsed != "" {
//                     fmt.Println("    - ", partition, "is", info.PercentUsed, "utilized")
//                 }
//             }
//
//             fmt.Println("\n  [+] Roles")
//             for _, role := range node.Info.Roles {
//                 fmt.Println("    - ", role)
//             }
//
//             fmt.Println()
//             for i := 0; i < 80; i++ {
//                 fmt.Print("=")
//             }
//         }
//
//         // Print detailed information about a specific cookbook
//         cookbook, ok, err := c.GetCookbook(findCookbook)
//         if err != nil {
//             fmt.Println("Error:", err)
//             os.Exit(1)
//         }
//         if !ok {
//             fmt.Println("\nCouldn't find that cookbook!")
//         } else {
//             fmt.Println("\n\nCookbook info:", findCookbook)
//             for _, version := range cookbook.Versions {
//                 currentVersion, ok, err := c.GetCookbookVersion(findCookbook, version.Version)
//                 if err != nil {
//                     fmt.Println("Error:", err)
//                     os.Exit(1)
//                 }
//                 if ok {
//                     if len(currentVersion.Files) > 0 {
//                         fmt.Println("\n  [+]", findCookbook, currentVersion.Version, "Cookbook Files")
//                         for _, cookbookFile := range currentVersion.Files {
//                             fmt.Println("    - ", cookbookFile.Name)
//                         }
//                     }
//                 }
//             }
//
//             fmt.Println()
//             for i := 0; i < 80; i++ {
//                 fmt.Print("=")
//             }
//         }
//
//         // Print detailed information about a specific role
//         role, ok, err := c.GetRole(findRole)
//         if err != nil {
//             fmt.Println("Error:", err)
//         }
//         if !ok {
//             fmt.Println("\nCouldn't find that role!")
//         } else {
//             fmt.Println("\n\nRole information:", role.Name)
//             fmt.Println("\n[+] Runlist")
//             for _, recipe := range role.RunList {
//                 fmt.Println("  - ", recipe)
//             }
//         }
//     }
//
// Which will output something like this:
//
//     ================================================================================
//     System info: neo4j.example.org
//
//       [+] IP Address: 10.100.1.2
//       [+] MAC Address: AA:BB:CC:DD:EE:FF
//       [+] Operating System: centos
//
//       [+] Filesystem Info
//         -  /dev/vda is 46% utilized
//         -  tmpfs is 1% utilized
//
//       [+] Roles
//         -  Base
//         -  Neo4j
//
//     ================================================================================
//
//     Cookbook info: neo4j
//
//       [+] neo4j 0.1.6 Cookbook Files
//         -  neo4j-server.properties
//         -  neo4j-service
//         -  neo4j-wrapper.conf
//
//       [+] neo4j 0.1.5 Cookbook Files
//         -  neo4j-service
//         -  neo4j-wrapper.conf
//
//       [+] neo4j 0.1.4 Cookbook Files
//         -  neo4j-service
//
//       [+] neo4j 0.1.3 Cookbook Files
//         -  neo4j-service
//
//       [+] neo4j 0.1.2 Cookbook Files
//         -  neo4j-service
//
//     ================================================================================
//
//     Role information: Neo4j
//
//     [+] Runlist
//       -  recipe[yum]
//       -  recipe[ldap]
//       -  recipe[system]
//       -  recipe[neo4j::server]
//       -  recipe[neo4j::web-ui]
package chef
