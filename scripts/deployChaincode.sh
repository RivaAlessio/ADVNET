#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh
. $SCRIPTS_DIR/utils/environment.sh

function packageChaincode() {
    local chaincode_name=$1
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincode_name}.tar.gz"
    local chaincode_label="${chaincode_name}_1.0"

    infoln "Packaging chaincode $chaincode_name"

    infoln "Vendoring Go dependencies at $CHAINCODE_SRC_PATH"
    pushd $CHAINCODE_SRC_PATH
    GO111MODULE=on go mod vendor
    popd
    successln "Finished vendoring Go dependencies"

    # if [ ! -d $CC_PACKAGE_FOLDER_OUTPUT ]; then
    #     mkdir $CC_PACKAGE_FOLDER_OUTPUT
    # fi

    set -x
    peer lifecycle chaincode package $chaincode_package_path --path $CHAINCODE_PACKAGE_DIR --lang $CHAINCODE_LANGUAGE --label $chaincode_label >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    cat log.txt
    verifyResult $res "Chaincode packaging has failed"
    successln "Chaincode is packaged"
}

function installChaincode() {
    local chaincodeName=$1
    local channelName=$2
    local orgType=$3
    local orgId=$4
    local peerId=$5
    local chaincodePackagePath="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    local peerName="peer${peerId}.${orgType}${orgId}"

    infoln "Installing chaincode ${chaincodeName} in channel ${channelName} of ${peerName}..."

    selectPeer $orgType $orgId $peerId

    packageName="${chaincodeName}_1.0"
    packageId=$(getPackageId $packageName)

    set -x
    packageInfo=$(peer lifecycle chaincode queryinstalled) >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    packageId=$(echo "$packageInfo" | sed -n "s/Package ID: //; s/, Label: ${packageName}$//p")

    if [ -z $packageId ]; then
        # if $packageId is empty, the package is not installed
        infoln "Package ${packageName} is not installed. Installing package ${packageName}..."
        # echo "Empty"
        set -x
        peer lifecycle chaincode install $chaincodePackagePath >&log.txt
        res=$?
        { set +x; } 2>/dev/null
        cat log.txt
        verifyResult $res "Chaincode installation on ${peerName} has failed"
        successln "Chaincode is installed on ${peerName} with id: ${packageId}"
    else
        infoln "Package ${packageName} is installed on ${peerName} with id: ${packageId} and won't be installed again."
    fi
}

function approveForMyOrg() {
    local chaincodeName=$1
    local channelName=$2
    local orgType=$3
    local orgId=$4
    local peerId=$5
    local chaincode_package_path="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    local peer_name="peer${peerId}.${orgType}${orgId}"

    infoln "Approving chaincode ${chaincodeName} in channel ${channelName} of ${peer_name}..."

    selectPeer $orgType $orgId $peerId

    packageName="${chaincodeName}_1.0"
    packageId=$(getPackageId $packageName)

    infoln "My package id:$packageId"

    set -x
    peer lifecycle chaincode approveformyorg -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME --tls --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --version 1.0 --package-id $packageId --sequence 1 >&log.txt
    res=$?
    { set +x; } 2>/dev/null

    cat log.txt
}

function commitChaincode() {
    local chaincodeName=$1
    local channelName=$2
    local orgTypes=$3
    local orgNum=$4
    local peerNum=$5

    local chaincodePackagePath="$CHAINCODE_PACKAGE_DIR/${chaincodeName}.tar.gz"
    local peerName="peer${peerId}.${orgType}${orgId}"

    infoln "Commiting chaincode $chaincodeName in channel '$channelName'..."

    parsePeerConnectionParameters $orgTypes $orgNum $peerNum
    infoln "peerConnectionParams: $peerConnectionParams"
	
	#--signature-policy="OR('adv0MSP.member','pub0MSP.member')"
    set -x
    peer lifecycle chaincode commit -o $ORDERER_ADDRESS --ordererTLSHostnameOverride $ORDERER_HOSTNAME  --cafile $ORDERER_CA --channelID $channelName --name $chaincodeName --tls $peerConnectionParams --version 1.0 --sequence 1
    res=$?
    { set +x; } 2>/dev/null

    # peer lifecycle chaincode querycommitted --channelID $channelName --name $chaincodeName

}



MODE=$1
CHAINCODE_NAME=$2
CHANNEL_NAME=$3
if [ "$MODE" == "package" ]; then
  packageChaincode $CHAINCODE_NAME
elif [ "$MODE" == "install" ]; then
  installChaincode $CHAINCODE_NAME $CHANNEL_NAME "adv" 0 0
  installChaincode $CHAINCODE_NAME $CHANNEL_NAME "pub" 0 0
elif [ "$MODE" == "approve" ]; then
  approveForMyOrg $CHAINCODE_NAME $CHANNEL_NAME "adv" 0 0
  approveForMyOrg $CHAINCODE_NAME $CHANNEL_NAME "pub" 0 0
elif [ "$MODE" == "commit" ]; then
  commitChaincode $CHAINCODE_NAME $CHANNEL_NAME "adv,pub" 1 1
fi