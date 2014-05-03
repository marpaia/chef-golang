#!/bin/bash
#
# This script should only be used if you need ot regen the goiardi data/index and keys
#
#  requires knife. We don't use this on each build beacuse of that.
#  TODO: at some point we can bootstrap with a tool thats not knife ;)
#
#-------------------------------------------------------------------------------
rundir="/tmp/goiardi"

if [ -d $rundir ]; then
  rm -f $rundir/*
else
  mkdir -p $rundir
fi

pushd $(dirname "${0}") > /dev/null
basedir=$(pwd -L)
# Use "pwd -P" for the path without links. man bash for more info.
popd > /dev/null

# shut it if its running
$basedir/stop_server.sh

# remove old data
rm $basedir/keys/*
rm $basedir/seed_data/*

set -x
set -e

go get -u github.com/ctdk/goiardi
go install github.com/ctdk/goiardi

cd $basedir/chef
# start server with our support dirs
goiardi -A -V -H localhost -P 8443 -D $basedir/seed_data/data -i $basedir/seed_data/index -F 30  --conf-root $basedir/keys &

# dumb but we need to wait 1 for the keys to be populated.
while [ ! -f $basedir/keys/admin.pem ] ; do
  echo 'waiting for admin key to be written'
  sleep 1
done

# set environment
knife environment from file environment.json
# setup node
knife node from file node.json
# set up role
knife role from file role.json
# bag
knife data bag create mysql
knife data bag from file mysql mysql_bag.json
# client
knife client create neo4j.example.org -d -f $basedir/keys/neo4j.example.org.pem

# upload cook
# berks3 Berksfile
cd $basedir/chef/test_cook
berks install
berks upload
