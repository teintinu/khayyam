#!/usr/bin/env bash

set -euo pipefail

DIR=$(realpath `dirname $0`/../..)

cd $DIR/examples/new

bash "$NVM_DIR/nvm.sh" use
khayyam deps
khayyam clean

set +e

khayyam lint
khayyam test
echo "exit code expected=0 actual=$?"
