#!/bin/bash

function initCaliper() {
    cd $CALIPER_DIR 

    set -x
    npm install --only=prod @hyperledger/caliper-cli@$CALIPER_VERSION
    { set +x; } 2>/dev/null

    set -x
    #npx caliper bind --caliper-bind-sut fabric:$FABRIC_VERSION
	npx caliper bind --caliper-bind-sut fabric:2.1
    { set +x; } 2>/dev/null

    cd $PROJECT_DIR
}

function caliperLaunchManager() {
    echo $PWD
    cd $CALIPER_DIR 

    set -x
    npx caliper launch manager --caliper-workspace $CALIPER_WORKSPACE --caliper-networkconfig $CALIPER_NETWORK_CONFIG --caliper-benchconfig $CALIPER_BENCH_CONFIG --caliper-flow-only-test --caliper-fabric-gateway-enabled
    { set +x; } 2>/dev/null

    cd $PROJECT_DIR
}

function clearCaliper() {
    cd $CALIPER_DIR 

    rm -rf node_modules
    rm -rf package-lock.json

    cd $PROJECT_DIR
}



MODE=$1
CALIPER_VERSION=$2
FABRIC_VERSION=$3
CALIPER_WORKSPACE=$4
CALIPER_NETWORK_CONFIG=$5
CALIPER_BENCH_CONFIG=$6

if [ $MODE = "init" ]; then
    initCaliper
elif [ $MODE = "launch" ]; then
    caliperLaunchManager
elif [ $MODE = "clear" ]; then
    clearCaliper
else
    echo "Unsupported '$MODE' command."
fi