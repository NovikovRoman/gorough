package gorough

import (
	"errors"
)

func CurveToBezier(pointsIn []Point, curveTightness float64) (out []Point, err error) {
	if len(pointsIn) < 3 {
		err = errors.New("A curve must have at least three points. ")
		return
	}

	if len(pointsIn) == 3 {
		out = append(out, pointsIn[0])
		out = append(out, pointsIn[1])
		out = append(out, pointsIn[2])
		out = append(out, pointsIn[2])
		return
	}

	points := []Point{pointsIn[0], pointsIn[1]}
	for i := 1; i < len(pointsIn); i++ {
		points = append(points, pointsIn[i])
		if i == len(pointsIn)-1 {
			points = append(points, pointsIn[i])
		}
	}

	b := make([]Point, 4)
	s := 1 - curveTightness
	out = append(out, points[0])
	for i := 1; (i + 2) < len(points); i++ {
		cachedVertArray := points[i]
		b[0] = cachedVertArray
		b[1] = Point{
			X: cachedVertArray.X + (s*points[i+1].X-s*points[i-1].X)/6,
			Y: cachedVertArray.Y + (s*points[i+1].Y-s*points[i-1].Y)/6,
		}
		b[2] = Point{
			X: points[i+1].X + (s*points[i].X-s*points[i+2].X)/6,
			Y: points[i+1].Y + (s*points[i].Y-s*points[i+2].Y)/6,
		}
		b[3] = Point{
			X: points[i+1].X,
			Y: points[i+1].Y,
		}

		out = append(out, b[1], b[2], b[3])
	}

	return
}
