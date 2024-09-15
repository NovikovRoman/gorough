package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type rectOpt func(*rect)

func RectPen(p Pen) rectOpt {
	return func(r *rect) {
		r.pen = p
	}
}

func RectStyle(s Style) rectOpt {
	return func(r *rect) {
		r.style = s
	}
}
func RectFiller(f Filler) rectOpt {
	return func(r *rect) {
		r.filler = f
	}
}

type rect struct {
	pen        Pen
	style      Style
	filler     Filler
	operations []operation
}

func (r rect) Name() string {
	return shapeRect
}

func (r rect) Operations() []operation {
	return r.operations
}

func (r rect) Attributes() Attributes {
	return map[string]string{
		"stroke":       r.style.Stroke,
		"stroke-width": tools.FloatToString(r.style.StrokeWidth),
		"fill":         r.style.Fill,
		"fill-weight":  tools.FloatToString(r.style.FillWeight),
	}
}

func (r rect) Style() Style {
	return r.style
}

func (r rect) Filler() Filler {
	return r.filler
}

func NewRect(p Point, width, height float64, opts ...rectOpt) rect {
	r := rect{
		pen:        PenDefault(),
		style:      StyleDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&r)
	}

	r.style = r.style.canonicalValues()

	outline := rectangleOperation(p, width, height, r.pen)

	if r.style.Fill != "" {
		points := []Point{p, {X: p.X + width, Y: p.Y}, {X: p.X + width, Y: p.Y + height}, {X: p.X, Y: p.Y + height}}

		if r.filler == nil {
			r.operations = append(r.operations, solidFillPolygon(points, r.pen))

		} else {
			r.operations = append(r.operations, patternFillPolygon(points, r.style, r.pen, r.filler))
		}
	}

	if r.style.Stroke != none {
		r.operations = append(r.operations, outline)
	}
	return r
}
