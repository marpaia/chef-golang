#!/bin/sh
/opt/chef-server/embedded/bin/runsvdir-start &
go test -v -cover
