#!/usr/bin/env bash

export PATH="$PWD:$PATH"

go build

cd snapshot/$1
./build.sh
