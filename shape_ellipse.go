package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type ellipseOpt func(*ellipse)

func EllipsePen(pen Pen) ellipseOpt {
	return func(e *ellipse) {
		e.pen = pen
	}
}

func EllipseStyle(s Style) ellipseOpt {
	return func(e *ellipse) {
		e.style = s
	}
}

func EllipseFiller(f Filler) ellipseOpt {
	return func(e *ellipse) {
		e.filler = f
	}
}

func EllipseCurveOpt(opt CurveOption) ellipseOpt {
	return func(e *ellipse) {
		e.curveOpt = opt
	}
}

type ellipse struct {
	pen        Pen
	style      Style
	filler     Filler
	curveOpt   CurveOption
	operations []operation
}

func (e ellipse) Name() string {
	return shapeEllipse
}

func (e ellipse) Operations() []operation {
	return e.operations
}

func (e ellipse) Attributes() Attributes {
	return map[string]string{
		"stroke":       e.style.Stroke,
		"stroke-width": tools.FloatToString(e.style.StrokeWidth),
		"fill":         e.style.Fill,
		"fill-weight":  tools.FloatToString(e.style.FillWeight),
	}
}

func (e ellipse) Style() Style {
	return e.style
}

func (e ellipse) Filler() Filler {
	return e.filler
}

func NewEllipse(center Point, width, height float64, opts ...ellipseOpt) (e ellipse) {
	e = ellipse{
		pen:        PenDefault(),
		style:      StyleDefault(),
		curveOpt:   CurveDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&e)
	}

	e.style = e.style.canonicalValues()
	e.operations = ellipseOperations(center, width, height, e.curveOpt, e.style, e.pen, e.filler)
	return
}
