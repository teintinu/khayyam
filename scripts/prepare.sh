ROOT=`realpath $(dirname $0)/..`
cd $ROOT/npm

node ../scripts/prepare.js $ROOT/npm/package.json
# npm install go-npm
