package gorough

import (
	"strconv"
)

type poligon struct {
	options    *LineOptions
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
		"stroke":       p.options.Styles.Stroke,
		"stroke-width": strconv.FormatFloat(p.options.Styles.StrokeWidth, 'f', -1, 64),
		"fill":         p.options.Styles.Fill,
		"fill-weight":  strconv.FormatFloat(p.options.Styles.FillWeight, 'f', -1, 64),
	}
}

func (p poligon) Styles() *Styles {
	return p.options.Styles
}

func NewPoligon(points []Point, opt *LineOptions) *poligon {
	if opt == nil {
		opt = &LineOptions{}
	}

	if opt.PenOptions == nil {
		opt.PenOptions = PenOptionsDefault()
	}

	if opt.Styles == nil {
		opt.Styles = StylesDefault()
	}
	opt.Styles.canonicalValues()

	var operations []operation
	operations = []operation{}
	outline := linearPathOperation(points, true, opt.PenOptions)

	if opt.Styles.Fill != "" {
		if opt.Styles.Filler == nil {
			operations = append(operations, solidFillPolygon(points, opt.PenOptions))

		} else {
			operations = append(operations, patternFillPolygon(points, &LineOptions{
				PenOptions: opt.PenOptions,
				Styles:     opt.Styles,
			}))
		}
	}

	if opt.Styles.Stroke != None {
		operations = append(operations, outline)
	}

	p := &poligon{
		options:    opt,
		operations: operations,
	}

	return p
}
