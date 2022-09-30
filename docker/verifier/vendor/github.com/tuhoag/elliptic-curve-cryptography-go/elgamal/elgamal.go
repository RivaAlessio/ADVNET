package elgamal

import (
	"github.com/bwesterb/go-ristretto"
	// "strings"
	// "bytes"
)

func Encrypt(privateKey *ristretto.Scalar, secretData *ristretto.Point, receiverPublicKey *ristretto.Point) (*ristretto.Point, *ristretto.Point) {
	var c1 ristretto.Point
	var c2 ristretto.Point

	c1.ScalarMultBase(privateKey)
	//
	// c1.ScalarMultBase(privateKey)

	c2.ScalarMult(receiverPublicKey, privateKey)
	c2.Add(&c2, secretData)

	return &c1, &c2
}

func Decrypt(privateKey *ristretto.Scalar, c1 *ristretto.Point, c2 *ristretto.Point) *ristretto.Point {
	var sharedKey ristretto.Point
	var rawData ristretto.Point

	sharedKey.ScalarMult(c1, privateKey)

	rawData.Sub(c2, &sharedKey)
	return &rawData
}
