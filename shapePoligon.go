package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
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
		"stroke-width": data_parser.FloatToString(p.options.Styles.StrokeWidth),
		"fill":         p.options.Styles.Fill,
		"fill-weight":  data_parser.FloatToString(p.options.Styles.FillWeight),
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
