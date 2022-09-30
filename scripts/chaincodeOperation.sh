#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh
. $SCRIPTS_DIR/utils/environment.sh

function chaincodeInovkeInit() {
    local chaincodeName=$1
    local channelName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5
    local fcnCall=$6

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum

    set -x
    peer chaincode invoke --channelID $channelName --name $chaincodeName -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --cafile $ORDERER_CA --tls $peerConnectionParams -c $fcnCall >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    cat log.txt
    verifyResult $res "Invoke execution on $peers failed "
    successln "Invoke transaction successful on $peers on channel '$channelName'"
}

CHAINCODE_NAME=$1
CHANNEL_NAME=$2
ORG_TYPES=$3
ORG_NUM=$4
PEER_NUM=$5
fcnCall=$6

chaincodeInovkeInit $CHAINCODE_NAME $CHANNEL_NAME $ORG_TYPES $ORG_NUM $PEER_NUM $fcnCall
