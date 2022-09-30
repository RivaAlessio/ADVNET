function selectPeer() {
    local orgType=$1
    local orgId=$2
    local peerId=$3

    # calculate port
    if [ $orgType = "adv" ]; then
        local basePort=$ADV_BASE_PORT
    elif [ $orgType = "pub" ]; then
        local basePort=$PUB_BASE_PORT
    else
        errorln "Org type $orgType is unsupported."
    fi
    local port=$(($basePort + $orgId * 100 + $peerId))
    local orgName="${orgType}${orgId}"
    local orgDomain="${orgName}.${PROJECT_NAME}.com"
    local peerDomain="peer${peerId}.${orgDomain}"

    infoln "Selecting organization ${peerDomain} with port $port"

    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="${orgName}MSP"
    export CORE_PEER_ADDRESS="0.0.0.0:${port}"
    export PEER_ORG_CA="${ORGANIZATION_OUTPUTS}/peerOrganizations/$orgDomain/peers/$peerDomain/tls/ca.crt"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PEER_ORG_CA
    export CORE_PEER_MSPCONFIGPATH="${ORGANIZATION_OUTPUTS}/peerOrganizations/${orgDomain}/users/Admin@${orgDomain}/msp"

    infoln "Selected peer ${peerDomain}: ${CORE_PEER_MSPCONFIGPATH}"
}

function getChannelTxPath() {
    echo $CHANNEL_PATH/$1.tx
}

function getBlockPath() {
    echo "${CHANNEL_PATH}/$1.block"
}

function getPackageId() {
    # packageInfo=$1
    packageName=$1
    set -x
    packageInfo=$(peer lifecycle chaincode queryinstalled) >&log.txt
    res=$?
    packageId=$(echo "$packageInfo" | sed -n "s/Package ID: ${packageName}:*//; s/, Label: ${packageName}$//p")
    { set +x; } 2>/dev/null


    echo "${packageName}:${packageId}"
}


function parsePeerConnectionParameters() {
    IFS=',' read -r -a orgTypes <<< $1
    local maxOrdId=$(($2 - 1))
    local maxPeerId=$(($3 - 1))

    peerConnectionParams=""
    peers=""
    for orgType in ${orgTypes[@]}; do
        for orgId in $(seq 0 $maxOrdId); do
            for peerId in $(seq 0 $maxOrdId); do
                selectPeer $orgType $orgId $peerId

                peers="$peers $CORE_PEER_ADDRESS"
                peerConnectionParams="$peerConnectionParams --peerAddresses $CORE_PEER_ADDRESS"
                peerConnectionParams="$peerConnectionParams --tlsRootCertFiles $CORE_PEER_TLS_ROOTCERT_FILE"
            done
        done
    done
}