package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	ristretto "github.com/bwesterb/go-ristretto"
	"github.com/go-redis/redis"
	elgamal "github.com/tuhoag/elliptic-curve-cryptography-go/elgamal"
	pedersen "github.com/tuhoag/elliptic-curve-cryptography-go/pedersen"
	utils "github.com/tuhoag/elliptic-curve-cryptography-go/utils"
)

const (
	SERVER_HOST = "0.0.0.0"
	SERVER_PORT = "9000"
	SERVER_TYPE = "tcp"
	SPLIT       = "SPLIT"
	REQUEST     = "request"
)

type CParam struct {
	Hval      string `json:"Hval"`
	Secretval string `json:"Secretval"`
	Sk        string `json:"Sk"`
	Pk        string `json:"Pk"`
}

var secretVal ristretto.Scalar
var H ristretto.Point
var sk ristretto.Scalar
var pk ristretto.Point

func main() {

	// //secretVal.Rand()
	// //H = *pedersen.GenerateH()
	// string := "yMUdpoU8NCPZiFmpXygUExfNbcEyzvqqr9f8he1f20Q="
	// Hgen, _ := utils.ConvertStringToPoint(string)
	// H = *Hgen
	// fmt.Println(H)
	// //stringSecret := "R9oNT1Lg5e3ntdAYvEl0BfHaM4ys0Qwuy/IoEfG05Ag="
	// stringSecret := os.Getenv("SECRET")
	// Sgen, _ := utils.ConvertStringToScalar(stringSecret)
	// secretVal = *Sgen

	// //stringSecretK := "DYMm0G2G3zc25FM6Xxuk07jtd9V3TeV0DE8rYHlgdQU="
	// stringSecretK := os.Getenv("SECRETKEY")
	// SgenK, _ := utils.ConvertStringToScalar(stringSecretK)
	// sk = *SgenK
	// //generate secret key and public key
	// //sk.Rand()
	// pk.ScalarMultBase(&sk)

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		//go processClient(connection)
		go processRequest(connection)
	}
}
func processRequest(connection net.Conn) {
	buffer := make([]byte, 15625)
	var s string
	var response string
	// var secretVal ristretto.Scalar
	// var H ristretto.Point
	// var sk ristretto.Scalar
	// var pk ristretto.Point

	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//defer connection.Close()
	s = string(buffer[:mLen])
	req := strings.Split(s, REQUEST)
	if req[0] == "proof" {

		cID := strings.Split(req[1], "CID")
		cp := RetrieveParam(cID[0])

		H, _ := utils.ConvertStringToPoint(cp.Hval)
		fmt.Println(H)
		secretVal, _ := utils.ConvertStringToScalar(cp.Secretval)
		pk, _ := utils.ConvertStringToPoint(cp.Pk)

		c, r := generateCommitment(H, secretVal)

		response = utils.ConvertPointToString(c) + SPLIT + utils.ConvertScalarToString(r) + SPLIT + utils.ConvertPointToString(pk)
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if req[0] == "Req_commit" {
		//r_received, _ := utils.ConvertStringToScalar(req[1])
		//c_req := pedersen.CommitTo(&H, r_received, &secretVal)

		cID := strings.Split(req[1], "CID")
		cp := RetrieveParam(cID[0])

		H, _ := utils.ConvertStringToPoint(cp.Hval)
		fmt.Println(H)
		secretVal, _ := utils.ConvertStringToScalar(cp.Secretval)

		r_received, _ := utils.ConvertStringToPoint(cID[1])
		c_req := CommitToTest(H, r_received, secretVal)

		response = utils.ConvertPointToString(c_req) + SPLIT
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if req[0] == "decrypt" {

		fmt.Println("decrypting")
		fmt.Println(req[1])
		// bufferz := make([]byte, 2024)
		// mLeng, err := connection.Read(bufferz)
		// if err != nil {
		// 	fmt.Println("Error reading:", err.Error())
		// }
		// fmt.Println("read buffer")
		// str := string(bufferz[:mLeng])
		// fmt.Println(str)
		cID := strings.Split(req[1], "CID")
		cp := RetrieveParam(cID[0])

		H, _ := utils.ConvertStringToPoint(cp.Hval)
		secretVal, _ := utils.ConvertStringToScalar(cp.Secretval)
		sk, _ := utils.ConvertStringToScalar(cp.Pk)

		response := "RESPONSE"
		reqs := strings.Split(cID[1], "TPOC")
		for i := 0; i < len(reqs)-1; i++ {
			reqDec := strings.Split(reqs[i], "KEY")
			key, _ := utils.ConvertStringToPoint(reqDec[1])
			tpoc := strings.Split(reqDec[0], SPLIT)
			c_Enc, _ := utils.ConvertStringToPoint(tpoc[0])
			c_Dec := elgamal.Decrypt(sk, key, c_Enc)
			rand, _ := utils.ConvertStringToPoint(tpoc[1])
			randDec := elgamal.Decrypt(sk, key, rand)
			c_Resp := CommitToTest(H, randDec, secretVal)
			c_DecS := utils.ConvertPointToString(c_Dec)
			c_RespS := utils.ConvertPointToString(c_Resp)
			response = response + c_DecS + SPLIT + c_RespS + "TPOC"

		}

		//fmt.Println(c_Dec)
		//fmt.Println(c_Resp)
		fmt.Println(response)
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if req[0] == "GenParam" {
		cID := strings.Split(req[1], "CID")
		HVal := strings.Split(cID[1], "HVAL")
		GenParam(HVal[0], cID[0])
	}
	defer connection.Close()
}

func generateCommitment(H *ristretto.Point, s *ristretto.Scalar) (*ristretto.Point, *ristretto.Scalar) {
	var r ristretto.Scalar
	r.Rand()

	c := pedersen.CommitTo(H, &r, s)
	return c, &r
}
func CommitToTest(H *ristretto.Point, r *ristretto.Point, x *ristretto.Scalar) *ristretto.Point {
	var result, transferPoint ristretto.Point
	transferPoint.ScalarMult(H, x)
	result.Add(r, &transferPoint)
	return &result
}
func RetrieveParam(CID string) *CParam {
	redisAddr := os.Getenv("REDIS")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	campaign, err := client.Get(CID).Result()
	if err == redis.Nil {
		fmt.Println("is null")
	} else if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println(campaign)
	var c CParam
	err = json.Unmarshal([]byte(campaign), &c)
	client.Close()
	return &c
}
func GenParam(Hval string, CID string) {
	redisAddr := os.Getenv("REDIS")
	var cParam CParam
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	if os.Getenv("TESTING") == "TRUE" {
		var pKey ristretto.Point
		cParam.Hval = Hval
		cParam.Secretval = os.Getenv("SECRET")
		cParam.Sk = os.Getenv("SECRETKEY")
		sKey, _ := utils.ConvertStringToScalar(cParam.Sk)
		pKey.ScalarMultBase(sKey)
		pKeyS := utils.ConvertPointToString(&pKey)
		cParam.Pk = pKeyS
	} else {
		var secretVal, pKey ristretto.Point
		var sKey ristretto.Scalar
		secretVal.Rand()
		sKey.Rand()
		pKey.ScalarMultBase(&sKey)
		secretS := utils.ConvertPointToString(&secretVal)
		skeyS := utils.ConvertScalarToString(&sKey)
		pkeyS := utils.ConvertPointToString(&pKey)

		cParam.Hval = Hval
		cParam.Secretval = secretS
		cParam.Sk = skeyS
		cParam.Pk = pkeyS
	}
	paramBytes, err := json.Marshal(cParam)
	if err != nil {
		fmt.Println(err)
	}
	err = client.Set(CID, paramBytes, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	client.Close()
}
