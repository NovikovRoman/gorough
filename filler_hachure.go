package gorough

import (
	"math"
	"sort"
)

type intersectionInfo struct {
	point    Point
	distance float64
}

type intersectionInfos []intersectionInfo

func (f intersectionInfos) Len() int      { return len(f) }
func (f intersectionInfos) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f intersectionInfos) Less(i, j int) bool {
	return f[i].distance-f[j].distance < 0
}

type hachureFiller struct {
	connectEnds  bool
	hachureAngle float64
	hachureGap   float64
}

func NewHachureFiller() Filler {
	return &hachureFiller{
		connectEnds:  false,
		hachureAngle: -41,
		hachureGap:   4,
	}
}

func (f *hachureFiller) SetAngle(angle float64) {
	f.hachureAngle = angle
}

func (f *hachureFiller) SetGap(gap float64) {
	f.hachureGap = gap
}

func (f *hachureFiller) setConnectEnds(b bool) {
	f.connectEnds = b
}

func (f hachureFiller) fillPolygon(points []Point, style Style, pen Pen, filler Filler) operation {
	lines := polygonHachureLines(points, f.hachureAngle, f.hachureGap, style.StrokeWidth)
	if f.connectEnds {
		connectingLines := f.connectingLines(points, lines)
		lines = append(lines, connectingLines...)
	}

	return operationFillSketch(f.renderLines(lines, pen))
}

func (f hachureFiller) renderLines(lines []Line, pen Pen) (commands []command) {
	for _, line := range lines {
		commands = append(commands, doubleLine(line.P1, line.P2, pen)...)
	}
	return commands
}

func (f hachureFiller) connectingLines(polygon []Point, lines []Line) (result []Line) {
	result = []Line{}

	if len(lines) == 0 {
		return
	}

	for i := 1; i < len(lines); i++ {
		prev := lines[i-1]
		if prev.length() < 3 {
			continue
		}

		segment := Line{P1: lines[i].P1, P2: prev.P2}
		if segment.length() > 3 {
			result = append(result, f.splitOnIntersections(polygon, segment)...)
		}
	}
	return
}

func (f hachureFiller) splitOnIntersections(polygon []Point, segment Line) []Line {
	distanceErr := math.Max(5, segment.length()*0.1)
	intersections := []intersectionInfo{}

	for i := 0; i < len(polygon); i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%len(polygon)]

		if DoIntersect(p1, p2, segment.P1, segment.P2) {
			if ip, ok := LineIntersection(p1, p2, segment.P1, segment.P2); ok {
				distance := Line{P1: ip, P2: segment.P1}.length()
				distance2 := Line{P1: ip, P2: segment.P2}.length()

				if distance > distanceErr && distance2 > distanceErr {
					intersections = append(intersections, intersectionInfo{
						point:    ip,
						distance: distance,
					})
				}
			}
		}
	}

	if len(intersections) > 1 {
		ips := make([]Point, len(intersections))
		sort.Sort(intersectionInfos(intersections))
		for i, v := range intersections {
			ips[i] = v.point
		}

		if !IsPointInPolygon(polygon, segment.P1) {
			ips = ips[1:]
		}

		if !IsPointInPolygon(polygon, segment.P2) {
			if len(ips) > 2 {
				ips = ips[:len(ips)-2]
			} else {
				ips = []Point{}
			}
		}

		if len(ips) <= 1 {
			if f.midPointInPolygon(polygon, segment) {
				return []Line{segment}
			}
			return []Line{}
		}

		spoints := []Point{segment.P1}
		spoints = append(spoints, ips...)
		spoints = append(spoints, segment.P2)
		slines := []Line{}
		for i := 0; i < (len(spoints) - 1); i += 2 {
			subSegment := Line{
				P1: spoints[i],
				P2: spoints[i+1],
			}
			if f.midPointInPolygon(polygon, subSegment) {
				slines = append(slines, subSegment)
			}
		}
		return slines
	}

	if f.midPointInPolygon(polygon, segment) {
		return []Line{segment}
	}

	return []Line{}
}

func (f hachureFiller) midPointInPolygon(polygon []Point, segment Line) bool {
	return IsPointInPolygon(polygon, Point{(segment.P1.X + segment.P2.X) / 2, (segment.P1.Y + segment.P2.Y) / 2})
}
