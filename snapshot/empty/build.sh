#!/usr/bin/env bash

set -euo pipefail

khayyam deps
khayyam clean
khayyam build
