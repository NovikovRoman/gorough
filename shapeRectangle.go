package gorough

import (
	"strconv"
)

type rectangle struct {
	options    *RectangleOptions
	operations []operation
}

func (r rectangle) Name() string {
	return shapeRectangle
}

func (r rectangle) Operations() []operation {
	return r.operations
}

func (r rectangle) Attributes() Attributes {
	return map[string]string{
		"stroke":       r.options.Styles.Stroke,
		"stroke-width": strconv.FormatFloat(r.options.Styles.StrokeWidth, 'f', -1, 64),
		"fill":         r.options.Styles.Fill,
		"fill-weight":  strconv.FormatFloat(r.options.Styles.FillWeight, 'f', -1, 64),
	}
}

func (r rectangle) Styles() *Styles {
	return r.options.Styles
}

func NewRectangle(p Point, width, height float64, opt *RectangleOptions) *rectangle {
	if opt == nil {
		opt = &RectangleOptions{}
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
	outline := rectangleOperation(p, width, height, opt)

	if opt.Styles.Fill != "" {
		points := []Point{p, {X: p.X + width, Y: p.Y}, {X: p.X + width, Y: p.Y + height}, {X: p.X, Y: p.Y + height}}

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

	r := &rectangle{
		options:    opt,
		operations: operations,
	}
	return r
}
