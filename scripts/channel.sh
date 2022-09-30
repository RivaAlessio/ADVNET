#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh
. $SCRIPTS_DIR/utils/environment.sh

function createChannelTx() {
    channelName=$1
	local channelTxPath=$(getChannelTxPath $channelName)

    infoln $channelName
    infoln $channelTxPath
	set -x
	configtxgen -profile TwoOrgsChannel -outputCreateChannelTx $channelTxPath -channelID $channelName -configPath $CONFIG_PATH
	res=$?
	{ set +x; } 2>/dev/null

    verifyResult $res "Failed to generate channel configuration transaction..."
}

function createChannel() {
    local channelName=$1
    local orgType=$2
    local orgId=$3

    selectPeer $orgType $orgId 0

    println "Generating channel tx..."
    local channelTxPath=$(getChannelTxPath $channelName)
    local blockPath=$(getBlockPath $channelName)

    println "Creating channel..."
    set -x
    peer channel create -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME -c $channelName -f $channelTxPath --outputBlock $blockPath --tls --cafile $ORDERER_CA -t 40s
    res=$?
    { set +x; } 2>/dev/null

	# cat log.txt
	# verifyResult $res "Channel creation failed"
}

function joinChannel() {
    local channelName=$1
    local orgType=$2
    local orgId=$3
    local peerId=$4

    local orgName="${orgType}${orgId}"

    infoln "Joining Channel $channelName from peer${peerId}.${orgName}"

    selectPeer $orgType $orgId $peerId

    local blockPath=$(getBlockPath $channelName)

    # peer channel getinfo -c $channelName

    set -x
    peer channel join -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA -b $blockPath
    res=$?
    { set +x; } 2>/dev/null

    verifyResult $res "Cannot join $channelName from peer${peerId}.${orgName}"
}

MODE=$1
CHANNEL_NAME=$2

if [ "$MODE" == "create-tx" ]; then
  createChannelTx $CHANNEL_NAME
elif [ "$MODE" == "create" ]; then
  createChannel $CHANNEL_NAME "adv" 0
elif [ "$MODE" == "join" ]; then
  joinChannel $CHANNEL_NAME "adv" 0 0
  sleep 3
  joinChannel $CHANNEL_NAME "pub" 0 0
fi