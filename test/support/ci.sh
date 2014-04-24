#!/bin/bash
set -e
set -x
test/support/chef_config.sh
go get -v github.com/ctdk/goiardi
test/support/start_server.sh
go get code.google.com/p/go.tools/cmd/vet
go install code.google.com/p/go.tools/cmd/vet
go vet
go test -v -cover
test/support/stop_server.sh
