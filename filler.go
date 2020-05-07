package gorough

import (
	"math"
	"sort"
)

type edgeEntry struct {
	ymin   float64
	ymax   float64
	x      float64
	islope float64
}

func (e edgeEntry) less(ee edgeEntry) bool {
	if e.ymin < ee.ymin {
		return true
	}

	if e.ymin > ee.ymin {
		return false
	}

	if e.x < ee.x {
		return true
	}

	if e.x > ee.x || e.ymax == ee.ymax {
		return false
	}

	return (e.ymax-ee.ymax)/math.Abs(e.ymax-ee.ymax) < 0
}

type edgeEntries []edgeEntry

func (e edgeEntries) Len() int           { return len(e) }
func (e edgeEntries) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e edgeEntries) Less(i, j int) bool { return e[i].less(e[j]) }

type activeEdgeEntry struct {
	s    float64
	edge edgeEntry
}

type activeEdgeEntries []activeEdgeEntry

func (e activeEdgeEntries) Len() int      { return len(e) }
func (e activeEdgeEntries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e activeEdgeEntries) Less(i, j int) bool {
	if e[i].edge.x == e[j].edge.x {
		return false
	}
	return (e[i].edge.x-e[j].edge.x)/math.Abs(e[i].edge.x-e[j].edge.x) < 0
}

type Filler interface {
	fillPolygon(points []Point, opt *LineOptions) operation
	SetAngle(float64)
	SetGap(float64)
	setConnectEnds(bool)
}

func polygonHachureLines(points []Point, hachureAngle float64, hachureGap float64, opt *LineOptions) []Line {
	rotationCenter := Point{}
	angle := math.Round(hachureAngle + 90)
	if angle != 0 {
		points = RotatePoints(points, rotationCenter, angle)
	}

	lines := straightHachureLines(points, hachureGap, opt)
	if angle != 0 {
		points = RotatePoints(points, rotationCenter, -angle)
		lines = RotateLines(lines, rotationCenter, -angle)
	}
	return lines
}

func straightHachureLines(vertices []Point, hachureGap float64, opt *LineOptions) (lines []Line) {
	if !vertices[0].Eq(vertices[len(vertices)-1]) {
		vertices = append(vertices, vertices[0])
	}

	if len(vertices) <= 2 {
		return
	}

	gap := initHachureGap(hachureGap, opt.Styles.StrokeWidth)
	gap = math.Max(gap, 0.1)

	// Create sorted edges table
	var (
		edges             []edgeEntry
		activeEdges       []activeEdgeEntry
		filterActiveEdges []activeEdgeEntry
	)
	edges = []edgeEntry{}

	for i := 0; i < len(vertices)-1; i++ {
		p1 := vertices[i]
		p2 := vertices[i+1]
		if p1.Y == p2.Y {
			continue
		}

		ymin := math.Min(p1.Y, p2.Y)
		x := p2.X
		if ymin == p1.Y {
			x = p1.X
		}
		edges = append(edges, edgeEntry{
			ymin:   ymin,
			ymax:   math.Max(p1.Y, p2.Y),
			x:      x,
			islope: (p2.X - p1.X) / (p2.Y - p1.Y),
		})
	}

	sort.Sort(edgeEntries(edges))

	if len(edges) == 0 {
		return
	}

	// Start scanning
	activeEdges = []activeEdgeEntry{}
	y := edges[0].ymin

	for len(activeEdges) > 0 || len(edges) > 0 {
		if len(edges) > 0 {
			ix := -1
			for i := range edges {
				if edges[i].ymin > y {
					break
				}
				ix = i
			}

			for _, e := range edges[0 : ix+1] {
				activeEdges = append(activeEdges, activeEdgeEntry{
					s:    y,
					edge: e,
				})
			}
			edges = edges[ix+1:]
		}

		filterActiveEdges = []activeEdgeEntry{}
		for i, a := range activeEdges {
			if activeEdges[i].edge.ymax > y {
				filterActiveEdges = append(filterActiveEdges, a)
			}
		}
		activeEdges = filterActiveEdges
		sort.Sort(activeEdgeEntries(activeEdges))

		// fill between the edges
		if len(activeEdges) > 1 {
			for i := 0; i < len(activeEdges); i = i + 2 {
				nexti := i + 1
				if nexti >= len(activeEdges) {
					break
				}
				ce := activeEdges[i].edge
				ne := activeEdges[nexti].edge

				lines = append(lines, Line{
					P1: Point{X: math.Round(ce.x), Y: y},
					P2: Point{X: math.Round(ne.x), Y: y},
				})
			}
		}

		y += gap
		for i := range activeEdges {
			activeEdges[i].edge.x = activeEdges[i].edge.x + (gap * activeEdges[i].edge.islope)
		}
	}

	return
}

func randOffsetWithRange(min, max float64, opt *PenOptions) float64 {
	return offset(min, max, opt.Roughness, 1)
}

func initHachureGap(hachureGap float64, strokeWidth float64) (gap float64) {
	gap = hachureGap
	if gap < 0 {
		gap = strokeWidth * 4
		if gap == 0 {
			gap = 4
		}
	}
	return
}
