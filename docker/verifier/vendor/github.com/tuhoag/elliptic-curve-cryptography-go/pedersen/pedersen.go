package pedersen

import (
	"github.com/bwesterb/go-ristretto"
)

func CommitTo(H *ristretto.Point, r *ristretto.Scalar, x *ristretto.Scalar) *ristretto.Point {
	var result, rPoint, transferPoint ristretto.Point
	rPoint.ScalarMultBase(r)
	transferPoint.ScalarMult(H, x)
	result.Add(&rPoint, &transferPoint)
	return &result
}

func GenerateH() *ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)

	return &H
}
