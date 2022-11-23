package main

import(
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric/common/util"
	"fmt"
	"encoding/json"
	ristretto "github.com/bwesterb/go-ristretto"
	"github.com/tuhoag/elliptic-curve-cryptography-go/utils"
	"net"
	"strings"
	elgamal "github.com/tuhoag/elliptic-curve-cryptography-go/elgamal"
	
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

type Proof struct {
	SumC string   `json:"SumC"`
	SumR string   `json:"SumR"`
	Keys []string `json:"Keys"`
}
type PoCTPoC struct {
	Proof Proof    `json:"PoC"`
	Tpocs []string `json:"TPoCs"`
}

const (
	SERVER_TYPE = "tcp"
	SPLIT       = "SPLIT"   //used to split string
	REQUEST     = "request" //
)

func (s *SmartContract) GenerateProof(ctx contractapi.TransactionContextInterface, campaignID string) (*Proof, error) {

	campaignArgs := util.ToChaincodeArgs("ReadCampaign", campaignID)
	response:=ctx.GetStub().InvokeChaincode("campaign",campaignArgs, "mychannel")
	if response.Message !=""{
		return nil, fmt.Errorf(response.Message)
	}
	var campaign Campaign
	err := json.Unmarshal([]byte(response.Payload), &campaign)
	
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