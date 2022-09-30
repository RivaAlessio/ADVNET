#!/bin/bash

. $SCRIPTS_DIR/utils/output.sh

function generateOrgs() {
    configPath=$ORG_CONFIG_PATH/crypto-config.yaml
    outputPath=$ORGANIZATION_OUTPUTS

    infoln "Config path: $configPath"
    infoln "Output path: $outputPath"

    set -x
    cryptogen generate --config=$configPath --output=$outputPath
    res=$?
    { set +x; } 2>/dev/null

    if [ $res -ne 0 ]; then
        fatalln "Failed to generate certificates..."
    fi

    . $SCRIPTS_DIR/utils/connectionProfile.sh

}

function createSystemGenesisBlock() {
    which configtxgen
    if [ "$?" -ne 0 ]; then
        fatalln "configtxgen tool not found."
    fi

    if [ ! -d $CHANNEL_PATH ]; then
        mkdir $CHANNEL_PATH
    fi

    infoln "Generating Orderer Genesis block"

    set -x
    configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock $CHANNEL_PATH/genesis.block -configPath $CONFIG_PATH
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        fatalln "Failed to generate orderer genesis block..."
    fi
}


MODE=$1
if [ "$MODE" == "orgs" ]; then
  generateOrgs
elif [ "$MODE" == "system-genesis-block" ]; then
  createSystemGenesisBlock
fi
