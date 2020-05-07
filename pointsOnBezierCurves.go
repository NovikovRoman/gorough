package gorough

func PointsOnBezierCurves(points []Point, tolerance, distance float64) (newPoints []Point) {
	if tolerance == 0 {
		tolerance = 0.15
	}
	numSegments := (len(points) - 1) / 3
	for i := 0; i < numSegments; i++ {
		offset := i * 3
		newPoints = getPointsOnBezierCurveWithSplitting(points, offset, tolerance, newPoints...)
	}

	if distance > 0 {
		return simplifyPoints(newPoints, 0, len(newPoints), distance)
	}
	return
}

func getPointsOnBezierCurveWithSplitting(points []Point, offset int, tolerance float64, newPoints ...Point) (outPoints []Point) {
	outPoints = append(outPoints, newPoints...)

	if flatness(points, offset) < tolerance {
		p0 := points[offset]
		if len(outPoints) == 0 || Distance(outPoints[len(outPoints)-1], p0) > 1 {
			outPoints = append(outPoints, p0)
		}
		outPoints = append(outPoints, points[offset+3])
		return outPoints
	}

	// subdivide
	t := 0.5
	p1 := points[offset+0]
	p2 := points[offset+1]
	p3 := points[offset+2]
	p4 := points[offset+3]

	q1 := lerp(p1, p2, t)
	q2 := lerp(p2, p3, t)
	q3 := lerp(p3, p4, t)

	r1 := lerp(q1, q2, t)
	r2 := lerp(q2, q3, t)

	red := lerp(r1, r2, t)

	outPoints = getPointsOnBezierCurveWithSplitting(
		[]Point{p1, q1, r1, red}, 0, tolerance, outPoints...)
	outPoints = getPointsOnBezierCurveWithSplitting(
		[]Point{red, r2, q3, p4}, 0, tolerance, outPoints...)
	return outPoints
}

func lerp(a Point, b Point, t float64) Point {
	return Point{
		X: a.X + (b.X-a.X)*t,
		Y: a.Y + (b.Y-a.Y)*t,
	}
}
