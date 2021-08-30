#!/usr/bin/env bash

set -e

ROOT=`realpath $(dirname $0)/..`
cd $ROOT/examples/simple-workspace

go run $ROOT/main.go $@
