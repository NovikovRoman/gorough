package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
)

type ellipse struct {
	options    *EllipseOptions
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
		"stroke":       e.options.Styles.Stroke,
		"stroke-width": data_parser.FloatToString(e.options.Styles.StrokeWidth),
		"fill":         e.options.Styles.Fill,
		"fill-weight":  data_parser.FloatToString(e.options.Styles.FillWeight),
	}
}

func (e ellipse) Styles() *Styles {
	return e.options.Styles
}

func NewEllipse(center Point, width, height float64, opt *EllipseOptions) *ellipse {
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

	e := &ellipse{
		options:    opt,
		operations: ellipseOperations(center, width, height, opt),
	}
	return e
}
