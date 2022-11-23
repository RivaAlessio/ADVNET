package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}
type TokenCollection struct {
	TPoC_User   string `json:"TPoC_User"`
	TPoC_Device string `json:"TPoC_Device"`
	Timestamp   string `json:"Timestamp"`
	DocType     string `json:"DocType"`
	ID_Campaign string `json:"ID_Campaign"`
	Duration    string `json:"Duration"`
}

func (s *SmartContract) TokenTransaction(ctx contractapi.TransactionContextInterface, TPoC_U string, TPoC_D string, ID_C string, ConnectionTime string, timestamp string) error {
	tokenBytes, err := ctx.GetStub().GetState(TPoC_U)
	if tokenBytes != nil {
		return fmt.Errorf("TPoC already collected")
	}
	if err != nil {
		return err
	}

	//currentTime := time.Now()
	fmt.Println("Collection time: ", timestamp)
	token := TokenCollection{
		TPoC_User:   TPoC_U,
		TPoC_Device: TPoC_D,
		Timestamp:   timestamp,
		DocType:     "TokenCollection",
		ID_Campaign: ID_C,
		Duration:    ConnectionTime,
	}
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}
	fmt.Println("inside token collection")
	return ctx.GetStub().PutState(TPoC_U, tokenJSON)
}
func (s *SmartContract) ReadToken(ctx contractapi.TransactionContextInterface, ID string) (*TokenCollection, error) {
	tokenBytes, err := ctx.GetStub().GetState(ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get TPoC %s: %v", ID, err)
	}
	if tokenBytes == nil {
		return nil, fmt.Errorf("TPoC %s does not exist", ID)
	}

	var token TokenCollection
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
func (s *SmartContract) QueryAllToken(ctx contractapi.TransactionContextInterface) ([]*TokenCollection, error) {
	queryString := fmt.Sprintf(`{"selector":{"DocType":"TokenCollection"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}
func (s *SmartContract) QueryAllTokenOfCampaign(ctx contractapi.TransactionContextInterface, campaignID string) ([]*TokenCollection, error) {
	queryString := fmt.Sprintf(`{"selector":{"DocType":"TokenCollection","ID_Campaign":"%s"}}`, campaignID)
	return getQueryResultForQueryString(ctx, queryString)
}
func (s *SmartContract) DeleteToken(ctx contractapi.TransactionContextInterface, tokenID string) error {
	exists, err := s.TokenExists(ctx, tokenID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the token %s does not exist", tokenID)
	}

	return ctx.GetStub().DelState(tokenID)
}
func (s *SmartContract) TokenExists(ctx contractapi.TransactionContextInterface, tokenID string) (bool, error) {

	tokenBytes, err := ctx.GetStub().GetState(tokenID)
	if err != nil {
		return false, fmt.Errorf("failed to read token %s from world state. %v", tokenID, err)
	}

	return tokenBytes != nil, nil
}
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*TokenCollection, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructTokenQueryResponseFromIterator(resultsIterator)
}

func constructTokenQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*TokenCollection, error) {
	var tokens []*TokenCollection
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var token TokenCollection
		err = json.Unmarshal(queryResult.Value, &token)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}
