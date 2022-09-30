package utils

import (
	"github.com/bwesterb/go-ristretto"
)

func Split(targetPoint *ristretto.Point, n int) []*ristretto.Point {
	points := make([]*ristretto.Point, n)

	var sum ristretto.Point
	sum.SetZero()
	// sum
	for i := 0; i < n-1; i++ {
		points[i] = &ristretto.Point{}
		points[i].Rand()

		sum.Add(points[i], &sum)
		// fmt.Println("sum:", sum, "point:", points[i])
	}
	// fmt.Println("sum:", sum)
	// fmt.Println(points)

	// var temp ristretto.Point
	// temp.Set(targetPoint)
	points[n-1] = &ristretto.Point{}
	// fmt.Println(points)

	points[n-1].Set(targetPoint)

	// fmt.Println(points)

	points[n-1].Sub(points[n-1], &sum)

	return points
}
