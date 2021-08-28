#!/usr/bin/env bash

set -euo pipefail

(
  set +e
  khayyam env
  echo "exit code expected=1 actual=$?"
)
