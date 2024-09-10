package gorough

import (
	"github.com/NovikovRoman/gorough/internal/tools"
)

type arcOpt func(*arc)

func ArcPen(pen Pen) arcOpt {
	return func(a *arc) {
		a.pen = pen
	}
}

func ArcStyle(s Style) arcOpt {
	return func(a *arc) {
		a.style = s
	}
}

func ArcFiller(f Filler) arcOpt {
	return func(a *arc) {
		a.filler = f
	}
}

func ArcCurveOpt(opt CurveOption) arcOpt {
	return func(a *arc) {
		a.curveOpt = opt
	}
}

type arc struct {
	pen        Pen
	style      Style
	filler     Filler
	curveOpt   CurveOption
	operations []operation
}

func (a arc) Name() string {
	return shapeCurve
}

func (a arc) Operations() []operation {
	return a.operations
}

func (a arc) Attributes() Attributes {
	return map[string]string{
		"stroke":       a.style.Stroke,
		"stroke-width": tools.FloatToString(a.style.StrokeWidth),
		"fill":         a.style.Fill,
		"fill-weight":  tools.FloatToString(a.style.FillWeight),
	}
}

func (a arc) Style() Style {
	return a.style
}

func (a arc) Filler() Filler {
	return a.filler
}

func NewArc(center Point, width, height, start, stop float64, closed bool, opts ...arcOpt) (a arc) {
	a = arc{
		pen:        PenDefault(),
		style:      StyleDefault(),
		curveOpt:   CurveDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&a)
	}

	a.style = a.style.canonicalValues()

	outline := arcOperation(center, width, height, start, stop, closed, true, a.curveOpt, a.pen)

	if closed && a.style.Fill != "" {
		if a.filler == nil {
			s := arcOperation(center, width, height, start, stop, true, false, a.curveOpt, a.pen)
			s.code = operationFillPath
			a.operations = append(a.operations, s)

		} else {
			a.operations = append(a.operations, patternFillArc(center, width, height, start, stop, a.curveOpt, a.style, a.pen, a.filler))
		}
	}

	if a.style.Stroke != None {
		a.operations = append(a.operations, outline)
	}
	return
}
