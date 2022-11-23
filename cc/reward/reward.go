package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	ristretto "github.com/bwesterb/go-ristretto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
	"github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

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
type decryptResponse struct {
	cDec []string
	cReq []string
}

const (
	SERVER_TYPE = "tcp"
	SPLIT       = "SPLIT"   //used to split string
	REQUEST     = "request" //
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) ClaimReward(ctx contractapi.TransactionContextInterface, campaignID string, userID string, Tpocs string, timestamp string) (*Reward, error) {

	//## Check if campaign exist and retrieve information ##

	campaignArgs := util.ToChaincodeArgs("ReadCampaign", campaignID)
	response := ctx.GetStub().InvokeChaincode("campaign", campaignArgs, "mychannel")
	if response.Message != "" {
		return nil, fmt.Errorf(response.Message)
	}
	var campaign Campaign
	err := json.Unmarshal([]byte(response.Payload), &campaign)

	if err != nil {
		return nil, err
	}

	//## Check if user already rewarded ##

	addresses := strings.Split(campaign.Verifier, ",")
	nVerifier := len(addresses)

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

	var duration float64
	duration = 0

	startDate, err := time.Parse("2006-01-02T15:04:05", campaign.StartingDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02T15:04:05", campaign.EndingDate)
	if err != nil {
		return nil, err
	}

	tpocSplt := strings.Split(Tpocs, "RWRD")
	//tokens:=make([]TokenCollection,len(tpocSplt))
	var tokens []TokenCollection
	var retrievedTpocU []string
	var retrievedTpocD []string
	//fmt.Println("nTPOC", len(tpocSplt))
	for _, tpoc := range tpocSplt {

		tokenArgs := util.ToChaincodeArgs("ReadToken", tpoc)
		response := ctx.GetStub().InvokeChaincode("token", tokenArgs, "mychannel")
		if response.Message != "" {
			continue
		}

		var token TokenCollection
		err = json.Unmarshal([]byte(response.Payload), &token)
		if err != nil {
			return nil, err
		}
		//test, _ := time.Parse("2006-01-02T15:04:05", string)

		collectDate, err := time.Parse("2006-01-02T15:04:05", token.Timestamp)
		if err != nil {
			continue
		}
		if collectDate.Before(endDate) && collectDate.After(startDate) && token.ID_Campaign == campaignID {
			tokens = append(tokens, token)
			retrievedTpocU = append(retrievedTpocU, token.TPoC_User)
			retrievedTpocD = append(retrievedTpocD, token.TPoC_Device)
		}
		// 	verifiedTPoC_u := verifyTpoc(token.TPoC_User, addresses)
		// 	verifiedTPoC_d := verifyTpoc(token.TPoC_Device, addresses)

		// 	if verifiedTPoC_u == verifiedTPoC_d {
		// 		time, err := strconv.ParseFloat(token.Duration, 64)
		// 		if err != nil {
		// 			continue
		// 		}
		// 		duration = duration + time
		// 	}
		// }

	}
	// fmt.Println("nTokens", len(tokens))
	// fmt.Println("nTPOCU", len(retrievedTpocU))
	// fmt.Println("nTPOCD", len(retrievedTpocD))

	if len(tokens) == 0 {
		return nil, fmt.Errorf("zero tokens")
	}

	requestU := generateRequests(retrievedTpocU, nVerifier, campaignID)
	requestD := generateRequests(retrievedTpocD, nVerifier, campaignID)
	decryptedU := make([]decryptResponse, nVerifier)
	decryptedD := make([]decryptResponse, nVerifier)

	for i := 0; i < nVerifier; i++ {
		decryptedU[i] = decryptRequest(addresses[i], requestU[i])
		decryptedD[i] = decryptRequest(addresses[i], requestD[i])
	}

	for i := 0; i < len(decryptedU[0].cDec); i++ {
		var sumDecU, sumReqU, sumDecD, sumReqD ristretto.Point
		sumDecU.SetZero()
		sumReqU.SetZero()
		sumDecD.SetZero()
		sumReqD.SetZero()
		for j := 0; j < nVerifier; j++ {
			cDecU, _ := utils.ConvertStringToPoint(decryptedU[j].cDec[i])
			cReqU, _ := utils.ConvertStringToPoint(decryptedU[j].cReq[i])
			sumDecU.Add(&sumDecU, cDecU)
			sumReqU.Add(&sumReqU, cReqU)

			cDecD, _ := utils.ConvertStringToPoint(decryptedD[j].cDec[i])
			cReqD, _ := utils.ConvertStringToPoint(decryptedD[j].cReq[i])
			sumDecD.Add(&sumDecD, cDecD)
			sumReqD.Add(&sumReqD, cReqD)
		}
		//fmt.Println(sumDec)
		//fmt.Println(sumReq)
		if sumDecU.Equals(&sumReqU) && sumDecD.Equals(&sumReqD) {
			fmt.Println("verified", i)
			time, err := strconv.ParseFloat(tokens[i].Duration, 64)
			if err != nil {
				continue
			}
			duration = duration + time
		} else {
			fmt.Println("Not verified")
		}
	}

	//## Retrieve Score & Calculate Reward ##
	rewardVal, err := strconv.ParseFloat(campaign.RewardValue, 64)
	if err != nil {
		return nil, err
	}
	score := ConcaveScore(duration)
	rewardAmount := rewardVal * score
	sReward := fmt.Sprintf("%f", rewardAmount)
	// fmt.Println("score-->", score)
	// fmt.Println("reward-->", sReward)

	//## Store transaction to state reward of user ##
	reward := Reward{
		ID_User:       h_ID,
		ID_Campaign:   campaignID,
		Reward_Amount: sReward,
		DocType:       "reward",
		Timestamp:     timestamp,
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
}
func generateRequests(tpocs []string, len int, cID string) []string {
	request := make([]string, len)
	for i := 0; i < len; i++ {
		request[i] = "decrypt" + REQUEST + cID + "CID"
	}
	for _, tpoc := range tpocs {
		key := strings.Split(tpoc, "KEY")
		tpocSplit := strings.Split(key[0], SPLIT)
		j := 1
		for i := 0; i < len; i++ {
			request[i] = request[i] + tpocSplit[i+j] + SPLIT + tpocSplit[i+j+1] + "KEY" + key[1] + "TPOC"
			j += 1
		}

	}
	return request

}
func decryptRequest(address string, req string) decryptResponse {
	connection, err := net.Dial(SERVER_TYPE, address)
	if err != nil {
		panic(err)
	}
	//defer connection.Close()
	//send some data
	//s := "decrypt" + REQUEST + req
	//fmt.Println("request;", s)
	//_, err = connection.Write([]byte("decrypt" + REQUEST + req))
	_, err = connection.Write([]byte(req))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	buffer := make([]byte, 15625)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//fmt.Println("msg received")
	defer connection.Close()
	resp := strings.Split(string(buffer[:mLen]), "RESPONSE")
	//fmt.Println("response:", resp)
	tpocs := strings.Split(resp[1], "TPOC")
	//fmt.Println(tpocs[0])
	var decrypted decryptResponse
	for i := 0; i < (len(tpocs) - 1); i++ {
		tpoc := strings.Split(tpocs[i], SPLIT)
		decrypted.cDec = append(decrypted.cDec, tpoc[0])
		decrypted.cReq = append(decrypted.cReq, tpoc[1])
		//decrypted.cDec[i] = tpoc[0]
		//decrypted.cReq[i] = tpoc[1]
	}
	return decrypted

}
func ConcaveScore(duration float64) float64 {

	//calulate concave score according to BRAVE formula

	var a, b float64
	a = 13000
	b = 11000
	time := duration * 1000
	score := (math.Sqrt((b*b)+(4*a*time)) - b)
	score = score / (2 * a)
	score = math.Round(score)
	if score > 7 {
		score = 7
	}

	return score
}
