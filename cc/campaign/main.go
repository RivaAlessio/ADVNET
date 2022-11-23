package main

import (
	"log"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincodeCampaign, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincodeCampaign.Start(); err != nil {
		log.Panicf("Error starting campaign chaincode: %v", err)
	}
}
