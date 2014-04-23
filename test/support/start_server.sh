#!/bin/bash
#
# script to setup and pull goiardi into local system, create users, and run tests against.
#
#-------------------------------------------------------------------------------
set -e
set -x

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

go get github.com/ctdk/goiardi


cp $basedir/keys/* $rundir/
cp $basedir/seed_data/* $rundir/

goiardi -V -H localhost -P 8443 -D $rundir/data -i $rundir/index -F 30  -A  --conf-root $rundir &
