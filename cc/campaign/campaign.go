package main

import (
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	ristretto "github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

type SmartContract struct {
	contractapi.Contract
}

type Campaign struct {
	ID_Campaign   string `json:"ID_Campaign"`
	ID_Advertiser string `json:"ID_Advertiser"`
	ID_Publisher  string `json:"ID_Publisher"`
	DocType       string `json:"DocType"`
	Verifier      string `json:"Verifier"`
	RewardValue   string `json:"RewardValue"`
	StartingDate  string `json:"StartingDate"`
	EndingDate    string `json:"EndingDate"`
}

const (
	SERVER_TYPE = "tcp"
	SPLIT       = "SPLIT"   //used to split string
	REQUEST     = "request" //
)

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	verifiers := "verifier.adv.com:9000,verifier.pub.com:9000"
	campaigns := []Campaign{
		{ID_Campaign: "001", ID_Advertiser: "adv1", ID_Publisher: "pub1", DocType: "campaign", Verifier: "verifier.adv.com:9000,verifier.pub.com:9000", RewardValue: "0.5", StartingDate: "2022-05-01T00:00:01", EndingDate: "2022-09-01T23:59:59"},
		{ID_Campaign: "002", ID_Advertiser: "adv2", ID_Publisher: "pub2", DocType: "campaign", Verifier: "verifier.adv.com:9000,verifier.pub.com:9000", RewardValue: "0.8", StartingDate: "2022-07-01T00:00:01", EndingDate: "2022-12-01T23:59:59"},
	}
	//fmt.Println(assets)

	addr := strings.Split(verifiers, ",")
	for i := 0; i < len(addr); i++ {
		initParamTest(addr[i], "001")
	}
	for i := 0; i < len(addr); i++ {
		initParamTest(addr[i], "002")
	}
	for _, campaign := range campaigns {
		assetJSON, err := json.Marshal(campaign)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(campaign.ID_Campaign, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CampaignExists(ctx contractapi.TransactionContextInterface, campaignID string) (bool, error) {

	campaignBytes, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return false, fmt.Errorf("failed to read campaign %s from world state. %v", campaignID, err)
	}

	return campaignBytes != nil, nil
}

func (s *SmartContract) ReadCampaign(ctx contractapi.TransactionContextInterface, campaignID string) (*Campaign, error) {

	campaignBytes, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign %s: %v", campaignID, err)
	}
	if campaignBytes == nil {
		return nil, fmt.Errorf("campaign %s does not exist", campaignID)
	}

	var campaign Campaign
	err = json.Unmarshal(campaignBytes, &campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (s *SmartContract) DeleteCampaign(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.CampaignExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the campaign %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}
func (s *SmartContract) CreateCampaign(ctx contractapi.TransactionContextInterface, id string, idadv string, idpub string, doctype string, verifier string, reward string, start string, end string) error {

	exists, err := s.CampaignExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the campaign %s already exists", id)
	}
	addr := strings.Split(verifier, ",")
	for i := 0; i < len(addr); i++ {
		initParam(addr[i], id)
	}

	campaign := Campaign{
		ID_Campaign:   id,
		ID_Advertiser: idadv,
		ID_Publisher:  idpub,
		DocType:       doctype,
		Verifier:      verifier,
		RewardValue:   reward,
		StartingDate:  start,
		EndingDate:    end,
	}
	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, campaignJSON)
}
func (s *SmartContract) CreateTestCampaign(ctx contractapi.TransactionContextInterface, id string, idadv string, idpub string, doctype string, verifier string, reward string, start string, end string) error {

	exists, err := s.CampaignExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the campaign %s already exists", id)
	}
	addr := strings.Split(verifier, ",")
	for i := 0; i < len(addr); i++ {
		initParamTest(addr[i], id)
	}
	campaign := Campaign{
		ID_Campaign:   id,
		ID_Advertiser: idadv,
		ID_Publisher:  idpub,
		DocType:       doctype,
		Verifier:      verifier,
		RewardValue:   reward,
		StartingDate:  start,
		EndingDate:    end,
	}
	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, campaignJSON)
}
func initParam(address string, cID string) {
	connection, err := net.Dial(SERVER_TYPE, address)

	if err != nil {
		panic(err)
	}
	defer connection.Close()
	//send some data
	var H ristretto.Point
	H.Rand()
	Hs := utils.ConvertPointToString(&H)
	_, err = connection.Write([]byte("GenParam" + REQUEST + cID + "CID" + Hs + "HVAL"))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

}
func initParamTest(address string, cID string) {
	connection, err := net.Dial(SERVER_TYPE, address)

	if err != nil {
		panic(err)
	}
	defer connection.Close()
	//send some data

	//yMUdpoU8NCPZiFmpXygUExfNbcEyzvqqr9f8he1f20Q=
	_, err = connection.Write([]byte("GenParam" + REQUEST + cID + "CID" + "yMUdpoU8NCPZiFmpXygUExfNbcEyzvqqr9f8he1f20Q=" + "HVAL"))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

}
