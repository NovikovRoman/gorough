package gorough

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

func (p Point) Eq(point Point) bool {
	return p.X == point.X && p.Y == point.Y
}

type Line struct {
	P1 Point
	P2 Point
}

func (l Line) length() float64 {
	return Distance(l.P1, l.P2)
}

// Distance returns the distance between 2 points
func Distance(p1 Point, p2 Point) float64 {
	return math.Sqrt(DistanceSq(p1, p2))
}

// DistanceSq returns the distance between 2 points square
func DistanceSq(p1 Point, p2 Point) float64 {
	return math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
}

// DistanceToSegmentSq returns the distance squared from a point p to the line segment vw
func DistanceToSegmentSq(p Point, v Point, w Point) float64 {
	l2 := DistanceSq(v, w)
	if l2 == 0 {
		return DistanceSq(p, v)
	}
	t := ((p.X-v.X)*(w.X-v.X) + (p.Y-v.Y)*(w.Y-v.Y)) / l2
	t = math.Max(0, math.Min(1, t))
	return DistanceSq(p, lerp(v, w, t))
}

func RotatePoints(points *[]Point, center Point, degrees float64) {
	if len(*points) == 0 {
		return
	}

	angle := (math.Pi / 180) * degrees
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	for i, p := range *points {
		(*points)[i].X = ((p.X - center.X) * cos) - ((p.Y - center.Y) * sin) + center.X
		(*points)[i].Y = ((p.X - center.X) * sin) + ((p.Y - center.Y) * cos) + center.Y
	}
	return
}

func RotateLines(lines *[]Line, center Point, degrees float64) {
	for i := range *lines {
		points := []Point{(*lines)[i].P1, (*lines)[i].P2}
		RotatePoints(&points, center, degrees)
		(*lines)[i].P1 = points[0]
		(*lines)[i].P2 = points[1]
	}
	return
}

func LineIntersection(a, b, c, d Point) (Point, bool) {
	a1 := b.Y - a.Y
	b1 := a.X - b.X
	c1 := a1*a.X + b1*a.Y
	a2 := d.Y - c.Y
	b2 := c.X - d.X
	c2 := a2*c.X + b2*c.Y
	determinant := a1*b2 - a2*b1
	if determinant == 0 {
		return Point{}, false
	}

	return Point{
		X: (b2*c1 - b1*c2) / determinant,
		Y: (a1*c2 - a2*c1) / determinant,
	}, true
}

func IsPointInPolygon(points []Point, point Point) bool {
	vertices := len(points)
	if vertices < 3 {
		return false
	}

	extreme := Point{X: float64(int(^uint(0) >> 1)), Y: point.Y}
	count := 0
	for i := 0; i < vertices; i++ {
		current := points[i]
		next := points[(i+1)%vertices]
		if DoIntersect(current, next, point, extreme) {
			if orientation(current, point, next) == 0 {
				return onSegment(current, point, next)
			}
			count++
		}
	}

	// true if count is off
	return count%2 == 1
}

func onSegment(p, q, r Point) bool {
	return q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) &&
		q.X <= math.Max(p.Y, r.Y) && q.X >= math.Min(p.Y, r.Y)
}

// For the ordered points p, q, r, return
// 0 if p, q, r are collinear
// 1 if Clockwise
// 2 if counterclickwise
func orientation(p, q, r Point) int {
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)
	if val == 0 {
		return 0
	}
	if val > 0 {
		return 1
	}
	return 2
}

func DoIntersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}

	// p1, q1 and p2 are colinear and p2 lies on segment p1q1
	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}

	// p1, q1 and p2 are colinear and q2 lies on segment p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2 and p1 are colinear and p1 lies on segment p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	// p2, q2 and q1 are colinear and q1 lies on segment p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}

// flatness (Adapted from https://seant23.wordpress.com/2010/11/12/offset-bezier-curves/)
func flatness(points []Point, offset int) float64 {
	p1 := points[offset+0]
	p2 := points[offset+1]
	p3 := points[offset+2]
	p4 := points[offset+3]

	ux := 3*p2.X - 2*p1.X - p4.X
	ux *= ux
	uy := 3*p2.Y - 2*p1.Y - p4.Y
	uy *= uy

	vx := 3*p3.X - 2*p4.X - p1.X
	vx *= vx
	vy := 3*p3.Y - 2*p4.Y - p1.Y
	vy *= vy

	if ux < vx {
		ux = vx
	}

	if uy < vy {
		uy = vy
	}

	return ux + uy
}
