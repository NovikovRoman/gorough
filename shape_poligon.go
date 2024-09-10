package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type poligonOpt func(*poligon)

func PoligonPen(pen Pen) poligonOpt {
	return func(p *poligon) {
		p.pen = pen
	}
}

func PoligonStyle(s Style) poligonOpt {
	return func(p *poligon) {
		p.style = s
	}
}

func PoligonFiller(f Filler) poligonOpt {
	return func(p *poligon) {
		p.filler = f
	}
}

type poligon struct {
	pen        Pen
	style      Style
	filler     Filler
	operations []operation
}

func (p poligon) Name() string {
	return shapePoligon
}

func (p poligon) Operations() []operation {
	return p.operations
}

func (p poligon) Attributes() Attributes {
	return map[string]string{
		"stroke":       p.style.Stroke,
		"stroke-width": tools.FloatToString(p.style.StrokeWidth),
		"fill":         p.style.Fill,
		"fill-weight":  tools.FloatToString(p.style.FillWeight),
	}
}

func (p poligon) Style() Style {
	return p.style
}

func (p poligon) Filler() Filler {
	return p.filler
}

func NewPoligon(points []Point, opts ...poligonOpt) poligon {
	p := poligon{
		pen:        PenDefault(),
		style:      StyleDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&p)
	}

	p.style = p.style.canonicalValues()
	outline := linearPathOperation(points, true, p.pen)

	if p.style.Fill != "" {
		if p.filler == nil {
			p.operations = append(p.operations, solidFillPolygon(points, p.pen))

		} else {
			p.operations = append(p.operations, patternFillPolygon(points, p.style, p.pen, p.filler))
		}
	}

	if p.style.Stroke != None {
		p.operations = append(p.operations, outline)
	}
	return p
}
