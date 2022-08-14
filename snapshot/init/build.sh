#!/usr/bin/env bash

set -euo pipefail

khayyam init

expected='a'
actual=`cat khayyam.yml``

if [ expected != actual ]
then
  echo "init failed"
  exit 1
fi