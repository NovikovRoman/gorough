package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
)

type arc struct {
	options    *EllipseOptions
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
		"stroke":       a.options.Styles.Stroke,
		"stroke-width": data_parser.FloatToString(a.options.Styles.StrokeWidth),
		"fill":         a.options.Styles.Fill,
		"fill-weight":  data_parser.FloatToString(a.options.Styles.FillWeight),
	}
}

func (a arc) Styles() *Styles {
	return a.options.Styles
}

func NewArc(center Point, width, height, start, stop float64, closed bool, opt *EllipseOptions) *arc {
	if opt == nil {
		opt = &EllipseOptions{}
	}

	if opt.PenOptions == nil {
		opt.PenOptions = PenOptionsDefault()
	}

	if opt.CurveOptions == nil {
		opt.CurveOptions = CurveOptionsDefault()
	}

	if opt.Styles == nil {
		opt.Styles = StylesDefault()
	}
	opt.Styles.canonicalValues()

	var operations []operation
	operations = []operation{}
	outline := arcOperation(center, width, height, start, stop, closed, true, opt)

	if closed && opt.Styles.Fill != "" {
		if opt.Styles.Filler == nil {
			s := arcOperation(center, width, height, start, stop, true, false, opt)
			s.code = operationFillPath
			operations = append(operations, s)

		} else {
			operations = append(operations, patternFillArc(center, width, height, start, stop, opt))
		}
	}

	if opt.Styles.Stroke != None {
		operations = append(operations, outline)
	}

	a := &arc{
		options:    opt,
		operations: operations,
	}

	return a
}
