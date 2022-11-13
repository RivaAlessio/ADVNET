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
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	elgamal "github.com/tuhoag/elliptic-curve-cryptography-go/elgamal"
	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
	"github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

type SmartContract struct {
	contractapi.Contract
}

// LEDGER ASSETS
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

// NOT LEDGER ASSETS
type Proof struct {
	SumC string   `json:"SumC"`
	SumR string   `json:"SumR"`
	Keys []string `json:"Keys"`
}
type PoCTPoC struct {
	Proof Proof    `json:"PoC"`
	Tpocs []string `json:"TPoCs"`
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

//##################################################################################################################################################################################################################################
//			CAMPAIGN FUNCTIONS
//##################################################################################################################################################################################################################################

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

//##################################################################################################################################################################################################################################
//##################################################################################################################################################################################################################################
//			PROOF
//##################################################################################################################################################################################################################################

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
func (s *SmartContract) GeneratePoCandTPoC(ctx contractapi.TransactionContextInterface, campaignID string, ntpoc int) (*PoCTPoC, error) {
	var secretK ristretto.Scalar

	proof, err := s.GenerateProof(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	fmt.Println(ntpoc)
	fmt.Println(proof)
	tpocs := make([]string, ntpoc)
	commit, err := utils.ConvertStringToPoint(proof.SumC)
	if err != nil {
		return nil, err
	}
	rands, err := utils.ConvertStringToScalar(proof.SumR)
	if err != nil {
		return nil, err
	}
	for i := 0; i < ntpoc; i++ {
		secretK.Rand()
		tpocs[i] = generateTpoc(commit, rands, len(proof.Keys), &secretK, proof.Keys)
		fmt.Println("TPOC:-->", tpocs[i])
	}
	var poctpoc PoCTPoC
	poctpoc.Proof = *proof
	poctpoc.Tpocs = tpocs

	return &poctpoc, nil
}

//##################################################################################################################################################################################################################################
//##################################################################################################################################################################################################################################
//			TOKEN TRANSACTION
//##################################################################################################################################################################################################################################

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

//##################################################################################################################################################################################################################################
//			REWARD
//##################################################################################################################################################################################################################################

func (s *SmartContract) ClaimReward(ctx contractapi.TransactionContextInterface, campaignID string, userID string, Tpocs string, timestamp string) (*Reward, error) {

	//## Check if campaign exist and retrieve information ##

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

		tokenBytes, err := ctx.GetStub().GetState(tpoc)

		if err != nil {
			return nil, fmt.Errorf("failed to get TPoC %s: %v", tpoc, err)
		}
		if tokenBytes == nil {
			//return nil, fmt.Errorf("TPoC %s does not exist", tpoc)
			continue
		}

		var token TokenCollection
		err = json.Unmarshal(tokenBytes, &token)
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

func (s *SmartContract) TestingProtocol(ctx contractapi.TransactionContextInterface, campaignID string, userID string, tpoc string, timestamp string) (*Reward, error) {

	//## Check if campaign exist and retrieve information ##
	fmt.Println("inside claim reward")
	campaignBytes, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign %s: %v", campaignID, err)
	}
	if campaignBytes == nil {
		return nil, fmt.Errorf("campaign %s does not exist", campaignID)
	}
	fmt.Println("before unmarshal")
	var campaign Campaign
	err = json.Unmarshal(campaignBytes, &campaign)
	if err != nil {
		return nil, err
	}
	fmt.Println("after unmarshal")
	//## Check if user already rewarded ##
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

	tokenBytes, err := ctx.GetStub().GetState(tpoc)

	if err != nil {
		return nil, fmt.Errorf("failed to get TPoC %s: %v", tpoc, err)
	}
	if tokenBytes == nil {
		return nil, fmt.Errorf("TPoC %s does not exist", tpoc)
	}

	var token TokenCollection
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}
	//test, _ := time.Parse("2006-01-02T15:04:05", string)

	collectDate, err := time.Parse("2006-01-02T15:04:05", token.Timestamp)
	if err != nil {
		return nil, err
	}
	if collectDate.Before(endDate) && collectDate.After(startDate) && token.ID_Campaign == campaignID {

		verifiedTPoC_u := verifyTpoc(token.TPoC_User, addresses)
		verifiedTPoC_d := verifyTpoc(token.TPoC_Device, addresses)

		if verifiedTPoC_u == verifiedTPoC_d {
			time, err := strconv.ParseFloat(token.Duration, 64)
			if err != nil {
				return nil, err
			}
			duration = duration + time
			fmt.Println("duration-->", duration)
		}
	}
	//## Retrieve Score & Calculate Reward ##
	rewardVal, err := strconv.ParseFloat(campaign.RewardValue, 64)
	if err != nil {
		return nil, err
	}
	score := ConcaveScore(duration)
	fmt.Println("score-->", score)
	rewardAmount := rewardVal * score
	fmt.Println("reward-->", rewardAmount)
	sReward := fmt.Sprintf("%f", rewardAmount)

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
func (s *SmartContract) DeleteReward(ctx contractapi.TransactionContextInterface, campaignID string, userID string) error {

	str := campaignID + userID
	data := []byte(str)
	hash := sha256.Sum256(data)
	h_ID := fmt.Sprintf("%x", hash[:])
	RewardBytes, err := ctx.GetStub().GetState(h_ID)
	if err != nil {
		return err
	}
	if RewardBytes == nil {
		return nil
	}

	return ctx.GetStub().DelState(h_ID)
}

//##################################################################################################################################################################################################################################
//##################################################################################################################################################################################################################################
//			UTILS
//##################################################################################################################################################################################################################################

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
func decryptTpoc(address, tpoc string) (string, string) {
	connection, err := net.Dial(SERVER_TYPE, address)
	if err != nil {
		panic(err)
	}
	//defer connection.Close()
	//send some data
	_, err = connection.Write([]byte("decrypt" + REQUEST + tpoc))
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//fmt.Println("msg received")
	defer connection.Close()
	resp := strings.Split(string(buffer[:mLen]), SPLIT)
	//fmt.Println("response: ", resp)
	commitDec := resp[0]
	commitCHK := resp[1]
	return commitDec, commitCHK
}
func ConcaveScore(duration float64) float64 {

	//calulate concave score according to BRAVE formula

	var a, b float64
	a = 13000
	b = 11000
	time := duration * 1000
	score := (math.Sqrt((b*b)+(4*a*time)) - b)
	score = score / (2 * a)
	if score >= 1 {
		score = math.Round(score)
	}
	score = math.Round(score)
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
	_, err = connection.Write([]byte("proof" + REQUEST))
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
func verifyTpoc(tpoc string, addr []string) bool {

	//TPoC has the form of : partial commit + "SPLIT" + partial r values + ... + "KEY" + key value
	var cDecryptedSum, cToCheckSum ristretto.Point
	cDecryptedSum.SetZero()
	cToCheckSum.SetZero()
	key := strings.Split(tpoc, "KEY")
	tpocSplit := strings.Split(key[0], SPLIT)
	j := 1
	for i := 0; i < len(addr); i++ {
		cDecStr, cToCheckStr := decryptTpoc(addr[i], tpocSplit[i+j]+SPLIT+tpocSplit[i+j+1]+"KEY"+key[1])
		cDec, _ := utils.ConvertStringToPoint(cDecStr)
		cToCheck, _ := utils.ConvertStringToPoint(cToCheckStr)
		cDecryptedSum.Add(&cDecryptedSum, cDec)
		cToCheckSum.Add(&cToCheckSum, cToCheck)
		j += 1
	}
	if cDecryptedSum.Equals(&cToCheckSum) {
		return true
	} else {
		return false
	}
}
func generateTpoc(sumC *ristretto.Point, sumR *ristretto.Scalar, n int, secretK *ristretto.Scalar, publicK []string) string {
	//split commit into n parts
	cSplitted := SplitPoint(sumC, n)
	//var sumRP ristretto.Point
	rSplitted := SplitScalar(sumR, n)
	//var rString string
	var rPoint ristretto.Point
	var tpoc string = ""
	var key ristretto.Point
	fmt.Print("splitted commit: ")
	fmt.Println(cSplitted)
	fmt.Print("splitted randoms: ")
	fmt.Println(rSplitted)

	//encrypt each partial commit
	for i := 0; i < n; i++ {
		//set point val to scalar
		rPoint.SetZero()
		rPoint.ScalarMultBase(rSplitted[i])

		public, _ := utils.ConvertStringToPoint(publicK[i])

		k, encCommit := elgamal.Encrypt(secretK, cSplitted[i], public)

		encCS := utils.ConvertPointToString(encCommit)
		k, encPoint := elgamal.Encrypt(secretK, &rPoint, public)
		encPS := utils.ConvertPointToString(encPoint)
		tpoc = tpoc + SPLIT + encCS + SPLIT + encPS
		key = *k
	}

	tpoc = tpoc + "KEY" + utils.ConvertPointToString(&key)
	return tpoc
}
func SplitScalar(target *ristretto.Scalar, n int) []*ristretto.Scalar {
	scalars := make([]*ristretto.Scalar, n)

	var sum ristretto.Scalar
	sum.SetZero()
	// sum
	for i := 0; i < n-1; i++ {
		scalars[i] = &ristretto.Scalar{}
		scalars[i].Rand()
		sum.Add(scalars[i], &sum)
	}
	scalars[n-1] = &ristretto.Scalar{}
	scalars[n-1].Set(target)
	scalars[n-1].Sub(scalars[n-1], &sum)

	return scalars
}
func SplitPoint(targetPoint *ristretto.Point, n int) []*ristretto.Point {
	points := make([]*ristretto.Point, n)

	var sum ristretto.Point
	sum.SetZero()
	// sum
	for i := 0; i < n-1; i++ {
		points[i] = &ristretto.Point{}
		points[i].Rand()
		sum.Add(points[i], &sum)
	}

	points[n-1] = &ristretto.Point{}
	points[n-1].Set(targetPoint)
	points[n-1].Sub(points[n-1], &sum)

	return points
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

//##################################################################################################################################################################################################################################
//##################################################################################################################################################################################################################################
//testing function  	---> DELETE THEM AT THE END <---
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
func (s *SmartContract) PartialClaim(ctx contractapi.TransactionContextInterface, campaignID string, userID string) error {
	//## Check if campaign exist ##
	fmt.Println("inside claim reward")
	campaignBytes, err := ctx.GetStub().GetState(campaignID)
	if err != nil {
		return fmt.Errorf("failed to get campaign %s: %v", campaignID, err)
	}
	if campaignBytes == nil {
		return fmt.Errorf("campaign %s does not exist", campaignID)
	}
	fmt.Println("before unmarshal")
	var campaign Campaign
	err = json.Unmarshal(campaignBytes, &campaign)
	if err != nil {
		return err
	}
	fmt.Println("after unmarshal")
	//## Check if user already rewarded for campaign c ##

	str := campaignID + userID
	data := []byte(str)
	hash := sha256.Sum256(data)
	h_ID := fmt.Sprintf("%x", hash[:])
	RewardBytes, err := ctx.GetStub().GetState(h_ID)
	if err != nil {
		return fmt.Errorf("error :%s", err)
	}
	if RewardBytes != nil {
		return fmt.Errorf("reward %s already exist", h_ID)
	}
	fmt.Println("no reward obtained", h_ID)
	return nil
}

//##################################################################################################################################################################################################################################
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
