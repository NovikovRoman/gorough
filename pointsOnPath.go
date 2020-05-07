package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
)

func PointsOnPath(path string, tolerance float64, distance float64) (out [][]Point, err error) {
	var segments []data_parser.Segment
	if segments, err = data_parser.ParsePath(path); err != nil {
		return
	}

	var (
		sets          [][]Point
		currentPoints []Point
		pendingCurve  []Point
	)
	start := Point{}

	normalize := data_parser.Normalize(data_parser.Absolutize(segments))
	for _, s := range normalize {
		switch s.Key {
		case "M":
			sets, currentPoints, pendingCurve = appendPendingPoints(sets, currentPoints, pendingCurve, tolerance)
			start = Point{
				X: s.Data[0],
				Y: s.Data[1],
			}
			currentPoints = append(currentPoints, start)

		case "L":
			currentPoints, pendingCurve = appendPendingCurve(currentPoints, pendingCurve, tolerance)
			currentPoints = append(currentPoints, Point{
				X: s.Data[0],
				Y: s.Data[1],
			})

		case "C":
			if len(pendingCurve) == 0 {
				lastPoint := start
				if len(currentPoints) > 0 {
					lastPoint = currentPoints[len(currentPoints)-1]
				}
				pendingCurve = append(pendingCurve, lastPoint)
			}
			pendingCurve = append(pendingCurve, Point{
				X: s.Data[0],
				Y: s.Data[1],
			})
			pendingCurve = append(pendingCurve, Point{
				X: s.Data[2],
				Y: s.Data[3],
			})
			pendingCurve = append(pendingCurve, Point{
				X: s.Data[4],
				Y: s.Data[5],
			})

		case "Z":
			currentPoints, pendingCurve = appendPendingCurve(currentPoints, pendingCurve, tolerance)
			currentPoints = append(currentPoints, start)
		}
	}

	sets, currentPoints, pendingCurve = appendPendingPoints(sets, currentPoints, pendingCurve, tolerance)
	if distance == 0 {
		return sets, nil
	}



	for _, set := range sets {
		simplifiedSet := Simplify(set, distance)
		if len(simplifiedSet) > 0 {
			out = append(out, simplifiedSet)
		}
	}
	return
}

func appendPendingCurve(currentPoints []Point, pendingCurve []Point, tolerance float64) ([]Point, []Point) {
	if len(pendingCurve) >= 4 {
		currentPoints = append(currentPoints, PointsOnBezierCurves(pendingCurve, tolerance, 0)...)
	}
	return currentPoints, []Point{}
}

func appendPendingPoints(sets [][]Point, currentPoints []Point, pendingCurve []Point, tolerance float64) ([][]Point, []Point, []Point) {
	currentPoints, pendingCurve = appendPendingCurve(currentPoints, pendingCurve, tolerance)
	if len(currentPoints) > 0 {
		sets = append(sets, currentPoints)
		currentPoints = []Point{}
	}
	return sets, currentPoints, pendingCurve
}
