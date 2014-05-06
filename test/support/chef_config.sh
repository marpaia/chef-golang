#!/bin/bash
# setup the chef/knife config for use with local testing
#
# TODO: these are shares stuffs
set -e
set -x

rundir="/tmp/goiardi"

pushd $(dirname "${0}") > /dev/null
basedir=$(pwd -L)
# Use "pwd -P" for the path without links. man bash for more info.
popd > /dev/null

if [ ! -f .chef/knife.rb ] ; then
    echo "no local knife config using our stub"
    mkdir -p .chef
    cp $basedir/knife.rb .chef/
fi

