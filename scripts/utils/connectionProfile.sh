#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${PEER0_PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        $SCRIPTS_DIR/utils/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG="adv0"
PEER0_PORT=1050
CAPORT=7054
PEERPEM=organizations/peerOrganizations/adv0.advnet.com/tlsca/tlsca.adv0.advnet.com-cert.pem
CAPEM=organizations/peerOrganizations/adv0.advnet.com/ca/ca.adv0.advnet.com-cert.pem

echo "$(yaml_ccp $ORG $PEER0_PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/adv0.advnet.com/connection-adv0.yaml

ORG="pub0"
PEER0_PORT=2050
CAPORT=8054
PEERPEM=organizations/peerOrganizations/pub0.advnet.com/tlsca/tlsca.pub0.advnet.com-cert.pem
CAPEM=organizations/peerOrganizations/pub0.advnet.com/ca/ca.pub0.advnet.com-cert.pem

echo "$(yaml_ccp $ORG $PEER0_PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/pub0.advnet.com/connection-pub0.yaml
