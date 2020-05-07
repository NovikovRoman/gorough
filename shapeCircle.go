package gorough

import "fmt"

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
		"stroke-width": fmt.Sprintf("%g", c.options.Styles.StrokeWidth),
		"fill":         c.options.Styles.Fill,
		"fill-weight":  fmt.Sprintf("%g", c.options.Styles.FillWeight),
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
