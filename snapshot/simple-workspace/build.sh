#!/usr/bin/env bash

set -euo pipefail

DIR=$(realpath `dirname $0`/../..)

cd $DIR/examples/simple-workspace

set +e

bash "$NVM_DIR/nvm.sh" use

monoclean deps
monoclean clean
monoclean build
monoclean run

monoclean lint
monoclean test
monoclean run
monoclean run appA
monoclean run appD

echo "exit code expected=0 actual=$?"
