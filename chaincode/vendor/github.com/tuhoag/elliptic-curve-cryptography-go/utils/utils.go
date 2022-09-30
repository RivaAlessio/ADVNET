package utils

import (
	b64 "encoding/base64"
	"fmt"

	ristretto "github.com/bwesterb/go-ristretto"
)

// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
// var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

func ConvertStringToPoint(s string) (*ristretto.Point, error) {
	bytes, err := b64.StdEncoding.DecodeString(s)

	if err != nil {
		return nil, err
	}

	point, err := ConvertBytesToPoint(bytes)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert %s to point", s)
	}

	return point, nil
}

func ConvertStringToScalar(s string) (*ristretto.Scalar, error) {
	bytes, err := b64.StdEncoding.DecodeString(s)

	if err != nil {
		return nil, err
	}

	scalar, err := ConvertBytesToScalar(bytes)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert %s to scalar", s)
	}

	return scalar, nil
}

func ConvertBytesToPoint(b []byte) (*ristretto.Point, error) {
	var H ristretto.Point
	var hBytes [32]byte

	copy(hBytes[:32], b[:])

	result := H.SetBytes(&hBytes)

	if !result {
		return nil, fmt.Errorf("Cannot convert point")
	}
	// fmt.Println("in convertBytesToPoint result:", result)

	return &H, nil
}

func ConvertBytesToScalar(b []byte) (*ristretto.Scalar, error) {
	var r ristretto.Scalar
	var rBytes [32]byte

	copy(rBytes[:32], b[:])

	r.SetBytes(&rBytes)

	return &r, nil
}

func ConvertScalarToString(scalar *ristretto.Scalar) string {
	s := b64.StdEncoding.EncodeToString(scalar.Bytes())
	return s
}

func ConvertPointToString(point *ristretto.Point) string {
	s := b64.StdEncoding.EncodeToString(point.Bytes())

	return s
}
