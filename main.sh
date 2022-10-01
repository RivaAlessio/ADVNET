#!/bin/bash

. $PWD/settings.sh

export CHANNEL_NAME="mychannel"
export LOG_LEVEL=INFO
export FABRIC_LOGGING_SPEC=INFO
export CHAINCODE_NAME="main"

function initialize() {
    $SCRIPTS_DIR/init.sh "orgs"
    sleep 1
    $SCRIPTS_DIR/init.sh "system-genesis-block"
}

function networkUp() {
    $SCRIPTS_DIR/network.sh "start" $LOG_LEVEL
}

function networkDown() {
    $SCRIPTS_DIR/network.sh "stop" $LOG_LEVEL
}

function clear() {
    $SCRIPTS_DIR/network.sh "clear"
}

function createChannel() {
    $SCRIPTS_DIR/channel.sh "create-tx" $CHANNEL_NAME
    sleep 3
    $SCRIPTS_DIR/channel.sh "create" $CHANNEL_NAME
}

function joinChannel() {
    $SCRIPTS_DIR/channel.sh "join" $CHANNEL_NAME
}

function packageChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "package" $CHAINCODE_NAME
}

function installChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "install" $CHAINCODE_NAME $CHANNEL_NAME
    $SCRIPTS_DIR/deployChaincode.sh "install" $CHAINCODE_NAME $CHANNEL_NAME
}

function approveChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "approve" $CHAINCODE_NAME $CHANNEL_NAME
    $SCRIPTS_DIR/deployChaincode.sh "approve" $CHAINCODE_NAME $CHANNEL_NAME
}

function commitChaincode() {
    $SCRIPTS_DIR/deployChaincode.sh "commit" $CHAINCODE_NAME $CHANNEL_NAME
}


function invokeChaincodeInit() {
    fcnCall='{"function":"'initLedger'","Args":[]}'
    $SCRIPTS_DIR/chaincodeOperation.sh $CHAINCODE_NAME $CHANNEL_NAME "adv,pub" 1 1 $fcnCall
}

function queryChaincode() {
    #fcnCall='{"Args":["org.hyperledger.fabric:GetMetadata"]}'
	fcnCall='{"function":"'ReadCampaign'","Args":["'001'"]}'
    $SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
	
	
	#sleep 2
	#fcnCall='{"function":"'TestCommit'","Args":[]}'
    #$SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
}
function proof(){
	fcnCall='{"function":"'GenerateProof'","Args":["'001'"]}'
    $SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
	sleep 1
	fcnCall='{"function":"'GenerateProof'","Args":["'002'"]}'
    $SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
}
function pocandtpoc(){
	fcnCall='{"function":"'GeneratePoCandTPoC'","Args":["'001'","2"]}'
	$SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
}
function tokencollection(){
	fcnCall='{"function":"'TokenTransaction'","Args":["'SPLITyqLgA1RXL3S+npfSVXOSYMF5pSfAz7dhZqCYyNhnAVM=SPLITHlGMv1ecFIx7GMAn2vXG9xFUGbBdopa7s8sYOHHKTEY=SPLITBFj2xt4L2+9yVf2VPlSPsmqG+aAJEuC0mNsiV8tq0R0=SPLITRqIYlcqmL/O80Xn4XvNb+fbHLuE/WjELmMr6+M/G1A0=KEYQAMga5l4/JX/CSEVnKAB52LL4erdkUP78iHVqQ0q3zQ='","'SPLITcJPDG9nehhVvTZgXYZeWpyLL4bQI1ZNkjN9nnARoV1c=SPLITlC5CJVplpAFbPxsaT4q7u7B/LHPnEU7bjJhZufqE720=SPLIT+Or332HtiSHNXCmuMQ7lhFuPTxBy7dzeL8iWdRsRcG4=SPLITdLOPcveNJW9ltUiy5BZpIqujK0ZFtPG401MFwFD8wQo=KEYolN5IPZASlN4tsdeG7HjLmHDEj/mtrUQvqQh1hdfdic='","'001'","'25'","'2022-08-11T00:00:01'"]}'
    $SCRIPTS_DIR/chaincodeOperation.sh $CHAINCODE_NAME $CHANNEL_NAME "adv,pub" 1 1 $fcnCall
	sleep 2
	fcnCall='{"function":"'TokenTransaction'","Args":["'SPLITPuxpWUVdwrJIMd7Tf/tNwZbp9PHpUJ0ycXx8sYWGrUg=SPLITkovHacagnAAXLwfWxTQDeBw1eJ56N0Tu5bPf9PCjPRc=SPLITiGZWSXH5UlAqaHJPhr5OVQnJUxLvca3KsBtroB5/uGA=SPLIT3q62s8F+cztIW33SQAtotvjURtozzu2qC+TXk3lUJV0=KEYMvfCBn23SX3raSd4yeXMjYIafy8n4CJiUgJL02Pfxm4='","'SPLITRHCApqI8BpHPkNf+ZlCYpoYQtix9c1Afno3UWuy8oT4=SPLITeIN8O7BrY5UisLcCqbLw/Ywb6Nuu4wBWtUwTIfLHdSY=SPLIT4lklmvhSYyo+GcjCA5bZ1CWyVA1xXE1zGt1TBHbOBSQ=SPLIT9BIqQ2Hixc5bA5okTZSZ0IxXYJz4hltJZWAyxRMf3ho=KEY+tc1b2quPlIVGm7Fwkkg+j/JDRbWWaS2M1CklJvqrwo='","'001'","'76'","'2022-08-20T00:00:01'"]}'
	$SCRIPTS_DIR/chaincodeOperation.sh $CHAINCODE_NAME $CHANNEL_NAME "adv,pub" 1 1 $fcnCall
}
function queryAllToken(){
	fcnCall='{"function":"'QueryAllToken'","Args":[]}'
    $SCRIPTS_DIR/chaincodeQuery.sh $CHAINCODE_NAME $CHANNEL_NAME "adv" 1 1 $fcnCall
}
function initCaliper() {
    $SCRIPTS_DIR/caliper.sh "init" $CALIPER_VERSION $FABRIC_VERSION
}

function caliperLaunch() {
    $SCRIPTS_DIR/caliper.sh "launch" $CALIPER_VERSION $FABRIC_VERSION $CALIPER_WORKSPACE $CALIPER_NETWORK_CONFIG $CALIPER_BENCH_CONFIG
}
function testingprotocol(){
	
	fcnCall='{"function":"'ClaimReward'","Args":["'001'","'TestId'","'SPLITyqLgA1RXL3S+npfSVXOSYMF5pSfAz7dhZqCYyNhnAVM=SPLITHlGMv1ecFIx7GMAn2vXG9xFUGbBdopa7s8sYOHHKTEY=SPLITBFj2xt4L2+9yVf2VPlSPsmqG+aAJEuC0mNsiV8tq0R0=SPLITRqIYlcqmL/O80Xn4XvNb+fbHLuE/WjELmMr6+M/G1A0=KEYQAMga5l4/JX/CSEVnKAB52LL4erdkUP78iHVqQ0q3zQ=RWRDSPLITPuxpWUVdwrJIMd7Tf/tNwZbp9PHpUJ0ycXx8sYWGrUg=SPLITkovHacagnAAXLwfWxTQDeBw1eJ56N0Tu5bPf9PCjPRc=SPLITiGZWSXH5UlAqaHJPhr5OVQnJUxLvca3KsBtroB5/uGA=SPLIT3q62s8F+cztIW33SQAtotvjURtozzu2qC+TXk3lUJV0=KEYMvfCBn23SX3raSd4yeXMjYIafy8n4CJiUgJL02Pfxm4='","today"]}'
	$SCRIPTS_DIR/chaincodeOperation.sh $CHAINCODE_NAME $CHANNEL_NAME "adv,pub" 1 1 $fcnCall
	
}

function clearCaliper() {
    $SCRIPTS_DIR/caliper.sh "clear"
}

MODE=$1

if [ $MODE = "network" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "up" ]; then
        initialize
        networkUp
    elif [ $SUB_MODE = "down" ]; then
        networkDown
        clear
    elif [ $SUB_MODE = "restart" ]; then
        networkDown
        clear
        initialize
        networkUp
		sleep 5
        createChannel
		sleep 2
        joinChannel
        packageChaincode
        installChaincode
        approveChaincode
        commitChaincode
		sleep 2
        invokeChaincodeInit
		#sleep 2
        #queryChaincode
		echo "---> Network Ready <---"
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "channel" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "create" ]; then
        createChannel
    elif [ $SUB_MODE = "join" ]; then
        joinChannel
    else
        echo "Unsupported $MODE $SUB_MODE command."
    fi

elif [ $MODE = "chaincode" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "package" ]; then
        packageChaincode
    elif [ $SUB_MODE = "install" ]; then
        installChaincode
    elif [ $SUB_MODE = "approve" ]; then
        approveChaincode
    elif [ $SUB_MODE = "commit" ]; then
        commitChaincode
	elif [ $SUB_MODE = "reinstall" ]; then
		packageChaincode
		installChaincode
		approveChaincode
        commitChaincode
    elif [ $SUB_MODE = "invoke-init" ]; then
        invokeChaincodeInit
    elif [ $SUB_MODE = "query" ]; then
        queryChaincode
	elif [ $SUB_MODE = "proof" ]; then
        proof
	elif [ $SUB_MODE = "poctpoc" ]; then
        pocandtpoc
	elif [ $SUB_MODE = "tokencollection" ]; then
        tokencollection
	elif [ $SUB_MODE = "queryalltoken" ]; then
        queryAllToken
	elif [ $SUB_MODE = "testing-protocol" ]; then
		tokencollection
		sleep 1
		testingprotocol
    else
        echo "Unsupported '$MODE $SUB_MODE' command."
    fi

elif [ $MODE = "caliper" ]; then
    SUB_MODE=$2
    if [ $SUB_MODE = "init" ]; then
        initCaliper
    elif [ $SUB_MODE = "launch" ]; then
        caliperLaunch
    elif [ $SUB_MODE = "clear" ]; then
        clearCaliper
    else
        echo "Unsupported '$MODE $SUB_MODE' command."
    fi
else
    echo "Unsupported $MODE command."
fi