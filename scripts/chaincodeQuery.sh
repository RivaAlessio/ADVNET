#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh
. $SCRIPTS_DIR/utils/environment.sh

function queryChaincode() {

    local chaincodeName=$1
    local channelName=$2
    local orgType=$3
    local orgNum=$4
    local peerNum=$5
    local fcnCall=$6

    infoln $orgType

    parsePeerConnectionParameters $orgType $orgNum $peerNum

    infoln "Invoke fcn call:${fcnCall} on peers: $peers"

    set -x
    peer chaincode query --channelID $channelName --name $chaincodeName -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --cafile $ORDERER_CA --tls $peerConnectionParams -c $fcnCall >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution on $peers failed "
    successln "Invoke transaction successful on $peers on channel '$channelName'"
}

CHAINCODE_NAME=$1
CHANNEL_NAME=$2
ORG_TYPE=$3
ORG_NUM=$4
PEER_NUM=$5
fcnCall=$6

queryChaincode $CHAINCODE_NAME $CHANNEL_NAME $ORG_TYPE $ORG_NUM $PEER_NUM $fcnCall