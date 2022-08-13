#!/usr/bin/env bash

set -e

export

ROOT=`realpath $(dirname $0)/..`
rm -Rf $ROOT/npm
mkdir -p $ROOT/npm
cd $ROOT

cd $ROOT/npm

node ../scripts/prepare.js $ROOT/npm/package.json
cp $ROOT/README.md $ROOT/npm/

npm install go-npm --save
rm -Rf node_modules
