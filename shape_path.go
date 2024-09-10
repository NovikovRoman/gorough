package gorough

import (
	"strings"

	"github.com/NovikovRoman/gorough/internal/tools"
)

type pathOpt func(*path)

func PathPen(pen Pen) pathOpt {
	return func(p *path) {
		p.pen = pen
	}
}

func PathStyle(s Style) pathOpt {
	return func(p *path) {
		p.style = s
	}
}

func PathSimplification(s float64) pathOpt {
	return func(p *path) {
		p.simplification = s
	}
}

func PathCombineNestedSvgPaths(b bool) pathOpt {
	return func(p *path) {
		p.combineNestedSvgPaths = b
	}
}

func PathFiller(f Filler) pathOpt {
	return func(p *path) {
		p.filler = f
	}
}

type path struct {
	pen                   Pen
	style                 Style
	simplification        float64
	combineNestedSvgPaths bool
	filler                Filler
	operations            []operation
}

func (p path) Name() string {
	return shapePoligon
}

func (p path) Operations() []operation {
	return p.operations
}

func (p path) Attributes() Attributes {
	return map[string]string{
		"stroke":       p.style.Stroke,
		"stroke-width": tools.FloatToString(p.style.StrokeWidth),
		"fill":         p.style.Fill,
		"fill-weight":  tools.FloatToString(p.style.FillWeight),
	}
}

func (p path) Style() Style {
	return p.style
}

func (p path) Filler() Filler {
	return p.filler
}

func NewPath(d string, opts ...pathOpt) (p path, err error) {
	p = path{
		pen:        PenDefault(),
		style:      StyleDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&p)
	}

	if d == "" {
		return
	}

	hasFill := p.style.Fill != None && p.style.Fill != "transparent"
	if (!hasFill || p.style.Stroke != "") && p.style.StrokeWidth == 0 {
		p.style.StrokeWidth = 1
	}

	//d = strings.ReplaceAll(d, "\n", " ")
	d = strings.ReplaceAll(d, "- ", "-")
	//d = reShapePathDoubleSpace.ReplaceAllString(d, " ")

	simplified := p.simplification < 1
	distance := p.pen.Roughness / 2
	if simplified {
		distance = 4 - 4*p.simplification
	}

	var points [][]Point
	if points, err = PointsOnPath(d, 1, distance); err != nil {
		return
	}

	if hasFill {
		p.fill(points)
	}

	if p.style.Stroke != None {
		if simplified {
			for _, pp := range points {
				p.operations = append(p.operations, linearPathOperation(pp, false, p.pen))
			}

		} else {
			var op operation
			if op, err = svgPath(d, p.pen); err != nil {
				return
			}
			p.operations = append(p.operations, op)
		}
	}

	return
}

func (p *path) fill(points [][]Point) {
	var combined []Point

	if p.combineNestedSvgPaths {
		combined = []Point{}
		for _, pp := range points {
			combined = append(combined, pp...)
		}

		if p.filler == nil {
			p.operations = append(p.operations, solidFillPolygon(combined, p.pen))

		} else {
			p.operations = append(p.operations, patternFillPolygon(combined, p.style, p.pen, p.filler))
		}
		return

	}
	for _, pp := range points {
		if p.filler == nil {
			p.operations = append(p.operations, solidFillPolygon(pp, p.pen))

		} else {
			p.operations = append(p.operations, patternFillPolygon(pp, p.style, p.pen, p.filler))
		}
	}
}
