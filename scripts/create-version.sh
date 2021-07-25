#!/usr/bin/env bash

set -e

ROOT=`realpath $(dirname $0)/..`
cd $ROOT

next=`git tag -l | node $ROOT/scripts/create-version.js`

echo "Publish next version: $next"
read -p "Are you sure? " -n 1 -r
echo 
if [[ $REPLY =~ ^[Yy]$ ]]
then
  git tag $next
  git push origin $next
fi
