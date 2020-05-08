package gorough

import (
	"strconv"
)

type circle struct {
	options    *EllipseOptions
	operations []operation
}

func (c circle) Name() string {
	return shapeEllipse
}

func (c circle) Operations() []operation {
	return c.operations
}

func (c circle) Attributes() Attributes {
	return map[string]string{
		"stroke":       c.options.Styles.Stroke,
		"stroke-width": strconv.FormatFloat(c.options.Styles.StrokeWidth, 'f', -1, 64),
		"fill":         c.options.Styles.Fill,
		"fill-weight":  strconv.FormatFloat(c.options.Styles.FillWeight, 'f', -1, 64),
	}
}

func (c circle) Styles() *Styles {
	return c.options.Styles
}

func NewCircle(center Point, diameter float64, opt *EllipseOptions) *circle {
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

	c := &circle{
		options:    opt,
		operations: ellipseOperations(center, diameter, diameter, opt),
	}
	return c
}
