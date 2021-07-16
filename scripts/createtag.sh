#!/usr/bin/env bash

set -e

ROOT=`realpath $(dirname $0)/..`
cd $ROOT

if [[ $1 =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]
then
  git tag v$1
  git push origin v$1
else 
  echo "invalid version"
  exit 1
fi
