package gorough

import (
	"math"
)

// Simplify function
func Simplify(points []Point, distance float64) []Point {
	return simplifyPoints(points, 0, len(points), distance)
}

// simplifyPoints
//
// Ramer–Douglas–Peucker algorithm
// https://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm
func simplifyPoints(points []Point, start int, end int, epsilon float64, newPoints ...Point) (outPoints []Point) {
	outPoints = append(outPoints, newPoints...)
	s := points[start]
	e := points[end-1]

	maxDistSq := 0.0
	maxNdx := 1

	for i := start + 1; i < end-1; i++ {
		distSq := DistanceToSegmentSq(points[i], s, e)
		if distSq > maxDistSq {
			maxDistSq = distSq
			maxNdx = i
		}
	}

	// if that point is too far, split
	if math.Sqrt(maxDistSq) > epsilon {
		outPoints = simplifyPoints(points, start, maxNdx+1, epsilon, outPoints...)
		outPoints = simplifyPoints(points, maxNdx, end, epsilon, outPoints...)

	} else if len(outPoints) > 0 {
		outPoints = append(outPoints, e)

	} else {
		outPoints = append(outPoints, s)
	}

	return outPoints
}
