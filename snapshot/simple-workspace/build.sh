#!/usr/bin/env bash

set -euo pipefail

DIR=$(realpath `dirname $0`/../..)

cd $DIR/examples/simple-workspace

set +e

bash "$NVM_DIR/nvm.sh" use

khayyam deps
khayyam clean
khayyam build
khayyam run

khayyam lint
khayyam test
khayyam run
khayyam run appA
khayyam run appD

echo "exit code expected=0 actual=$?"
