package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	ristretto "github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
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
type Reward struct {
	ID_User       string `json:"ID_User"`
	ID_Campaign   string `json:"ID_Campaign"`
	Reward_Amount string `json:"Reward_Amount"`
	DocType       string `json:"DocType"`
	Timestamp     string `json:"Timestamp"`
}
type TokenCollection struct {
	TPoC_User   string `json:"TPoC_User"`
	TPoC_Device string `json:"TPoC_Device"`
	Timestamp   string `json:"Timestamp"`
	DocType     string `json:"DocType"`
	ID_Campaign string `json:"ID_Campaign"`
	Duration    string `json:"Duration"`
}
type Proof struct {
	SumC string
	SumR string
	Keys []string
}

const (
	SERVER_TYPE = "tcp"
	SPLIT       = "SPLIT"
)

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	campaigns := []Campaign{
		{ID_Campaign: "001", ID_Advertiser: "adv1", ID_Publisher: "pub1", DocType: "campaign", Verifier: "verifier.adv.com:9000,verifier.pub.com:9000", RewardValue: "0.5", StartingDate: "2022-05-01T00:00:01", EndingDate: "2022-09-01T23:59:59"},
		{ID_Campaign: "002", ID_Advertiser: "adv2", ID_Publisher: "pub2", DocType: "campaign", Verifier: "verifier.adv.com:9000,verifier.pub.com:9000", RewardValue: "0.8", StartingDate: "2022-07-01T00:00:01", EndingDate: "2022-12-01T23:59:59"},
	}
	//fmt.Println(assets)
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
func (s *SmartContract) GenerateProof(ctx contractapi.TransactionContextInterface, campaignID string) (*Proof, error) {

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

	addresses := strings.Split(campaign.Verifier, ",")

	var sumC ristretto.Point
	var sumR ristretto.Scalar

	sumC.SetZero()
	sumR.SetZero()
	publicKeys := make([]string, len(addresses))

	for i, addr := range addresses {
		cS, rS, kS := receiveCommit(addr)
		publicKeys[i] = kS
		c, _ := utils.ConvertStringToPoint(cS)
		r, _ := utils.ConvertStringToScalar(rS)

		sumC.Add(&sumC, c)
		sumR.Add(&sumR, r)
	}
	sumCs := utils.ConvertPointToString(&sumC)
	sumRs := utils.ConvertScalarToString(&sumR)

	var proof Proof
	proof.SumC = sumCs
	proof.SumR = sumRs
	proof.Keys = publicKeys

	return &proof, nil
}
func (s *SmartContract) TokenTransaction(ctx contractapi.TransactionContextInterface, TPoC_U string, TPoC_D string, ID_C string, ConnectionTime string) error {
	tokenBytes, err := ctx.GetStub().GetState(TPoC_U)
	if tokenBytes != nil {
		return fmt.Errorf("TPoC already collected")
	}
	if err != nil {
		return err
	}

	currentTime := time.Now()
	//fmt.Println("Current Time in String: ", currentTime.String())
	token := TokenCollection{
		TPoC_User:   TPoC_U,
		TPoC_Device: TPoC_D,
		Timestamp:   currentTime.String(),
		DocType:     "TokenCollection",
		ID_Campaign: ID_C,
		Duration:    ConnectionTime,
	}
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(TPoC_U, tokenJSON)
}
func (s *SmartContract) ClaimReward(ctx contractapi.TransactionContextInterface, campaignID string, userID string, Tpocs []string) (*Reward, error) {

	//## Check if campaign exist ##
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
	//## Check if user already rewarded for campaign c ##
	addresses := strings.Split(campaign.Verifier, ",")
	str := campaignID + userID
	data := []byte(str)
	hash := sha256.Sum256(data)
	h_ID := fmt.Sprintf("%x", hash[:])
	RewardBytes, err := ctx.GetStub().GetState(h_ID)
	if err != nil {
		return nil, fmt.Errorf("error :%s", err)
	}
	if RewardBytes != nil {
		return nil, fmt.Errorf("reward %s already exist", h_ID)
	}
	//## Verify TPoC ##
	//## If verified TPoC time is between starting and ending date of campaign c --> sum duration ##
	var cDSum ristretto.Point
	var cCHKSum ristretto.Point
	var duration float64
	duration = 0
	for _, tpoc := range Tpocs {

		tokenBytes, err := ctx.GetStub().GetState(tpoc)

		if err != nil {
			return nil, fmt.Errorf("failed to get TPoC %s: %v", tpoc, err)
		}
		if campaignBytes == nil {
			return nil, fmt.Errorf("TPoC %s does not exist", tpoc)
		}

		var token TokenCollection
		err = json.Unmarshal(tokenBytes, &token)
		if err != nil {
			return nil, err
		}
		//test, _ := time.Parse("2006-01-02T15:04:05", string)
		startDate, err := time.Parse("2006-01-02T15:04:05", campaign.StartingDate)
		if err != nil {
			continue
		}
		endDate, err := time.Parse("2006-01-02T15:04:05", campaign.EndingDate)
		if err != nil {
			continue
		}
		collectDate, err := time.Parse("2006-01-02T15:04:05", token.Timestamp)
		if err != nil {
			continue
		}
		if collectDate.Before(endDate) && collectDate.After(startDate) && token.ID_Campaign == campaignID {

			keyU := strings.Split(token.TPoC_User, "KEY")
			tpocUSplit := strings.Split(keyU[0], SPLIT)
			// decrypted commit and commit generated from r value decrypted
			ct1D, ct1CHK := verifyTpoc(addresses[0], tpocUSplit[1]+SPLIT+tpocUSplit[2]+"KEY"+keyU[1])
			ct2D, ct2CHK := verifyTpoc(addresses[1], tpocUSplit[3]+SPLIT+tpocUSplit[4]+"KEY"+keyU[1])

			ct1DP, _ := utils.ConvertStringToPoint(ct1D)
			ct2DP, _ := utils.ConvertStringToPoint(ct2D)
			cDSum.SetZero()
			cDSum.Add(&cDSum, ct1DP)
			cDSum.Add(&cDSum, ct2DP)
			cCHKSum.SetZero()
			ct1CHKP, _ := utils.ConvertStringToPoint(ct1CHK)
			ct2CHKP, _ := utils.ConvertStringToPoint(ct2CHK)
			cCHKSum.Add(&cCHKSum, ct1CHKP)
			cCHKSum.Add(&cCHKSum, ct2CHKP)
			if cCHKSum.Equals(&cDSum) {
				keyD := strings.Split(token.TPoC_Device, "KEY")
				tpocDSplit := strings.Split(keyD[0], SPLIT)
				ct1D, ct1CHK := verifyTpoc(addresses[0], tpocDSplit[1]+SPLIT+tpocDSplit[2]+"KEY"+keyD[1])
				ct2D, ct2CHK := verifyTpoc(addresses[1], tpocDSplit[3]+SPLIT+tpocDSplit[4]+"KEY"+keyD[1])

				ct1DP, _ := utils.ConvertStringToPoint(ct1D)
				ct2DP, _ := utils.ConvertStringToPoint(ct2D)
				cDSum.SetZero()
				cDSum.Add(&cDSum, ct1DP)
				cDSum.Add(&cDSum, ct2DP)
				cCHKSum.SetZero()
				ct1CHKP, _ := utils.ConvertStringToPoint(ct1CHK)
				ct2CHKP, _ := utils.ConvertStringToPoint(ct2CHK)
				cCHKSum.Add(&cCHKSum, ct1CHKP)
				cCHKSum.Add(&cCHKSum, ct2CHKP)

				if cCHKSum.Equals(&cDSum) {
					time, err := strconv.ParseFloat(token.Duration, 64)
					if err != nil {
						continue
					}
					duration = duration + time
				}
			}

		}

	}
	//## Retrieve Score & Calculate Reward ##
	rewardVal, err := strconv.ParseFloat(campaign.RewardValue, 64)
	if err != nil {
		return nil, err
	}
	score := ConcaveScore(duration)
	rewardAmount := rewardVal * score
	currentTime := time.Now()
	sReward := fmt.Sprintf("%f", rewardAmount)

	//## Store transaction to state reward of user ##
	reward := Reward{
		ID_User:       h_ID,
		ID_Campaign:   campaignID,
		Reward_Amount: sReward,
		DocType:       "reward",
		Timestamp:     currentTime.String(),
	}
	rewardJSON, err := json.Marshal(reward)
	if err != nil {
		return nil, err
	}
	err = ctx.GetStub().PutState(h_ID, rewardJSON)
	if err != nil {
		return nil, err
	}
	return &reward, nil

	//var reward Reward
	//return &reward, nil
	//return nil, fmt.Errorf("error in claiming reward")

}

func ConcaveScore(duration float64) float64 {

	//calulate concave score according to BRAVE formula

	var a, b float64
	a = 13000
	b = 11000
	time := duration * 1000
	score := (math.Sqrt((b*b)+(4*a*time)) - b)
	score = score / (2 * a)
	if score > 7 {
		score = 7
	}
	return score
}
func receiveCommit(address string) (string, string, string) {

	connection, err := net.Dial(SERVER_TYPE, address)

	if err != nil {
		panic(err)
	}
	defer connection.Close()
	//send some data
	_, err = connection.Write([]byte("proof" + SPLIT))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	resp := strings.Split(string(buffer[:mLen]), SPLIT)
	commit := resp[0]
	random := resp[1]
	key := resp[2]
	//fmt.Println("Commit:", commit)
	//fmt.Println("Random:", random)
	return commit, random, key
}
func verifyTpoc(address string, tpoc string) (string, string) {
	connection, err := net.Dial(SERVER_TYPE, address)
	if err != nil {
		panic(err)
	}
	//defer connection.Close()
	//send some data
	_, err = connection.Write([]byte("decrypt" + SPLIT))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	tpocE := tpoc
	time.Sleep(1 * time.Millisecond)
	//send partial tpoc to verifier
	_, err = connection.Write([]byte(tpocE))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("partial tpoc: ")
	fmt.Println(tpocE)

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//fmt.Println("msg received")
	defer connection.Close()
	resp := strings.Split(string(buffer[:mLen]), SPLIT)
	fmt.Println("response: ", resp)
	commitDec := resp[0]
	commitCHK := resp[1]
	return commitDec, commitCHK

}

//testing function
func (s *SmartContract) TestCommit() string {
	var H ristretto.Point
	var sVal ristretto.Scalar
	var r1Val ristretto.Scalar
	var c1Val ristretto.Point

	sVal.Rand()
	H.Rand()

	c, r := generateCommitment(&H, &sVal)
	r1Val = *r
	c1Val = *c
	fmt.Println("commitment:\n", c1Val)
	fmt.Println("blinding factor:\n", r1Val)

	comm, rand, k := receiveCommit("verifier.pub.com:9000")
	fmt.Println(comm, rand, k)
	return "test"
}
func generateCommitment(H *ristretto.Point, s *ristretto.Scalar) (*ristretto.Point, *ristretto.Scalar) {
	var r ristretto.Scalar
	r.Rand()

	c := pedersen.CommitTo(H, &r, s)
	return c, &r
}
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
