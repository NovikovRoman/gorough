package gorough

import "fmt"

type curve struct {
	options    *EllipseOptions
	operations []operation
}

func (c curve) Name() string {
	return shapeCurve
}

func (c curve) Operations() []operation {
	return c.operations
}

func (c curve) Attributes() Attributes {
	return map[string]string{
		"stroke":       c.options.Styles.Stroke,
		"stroke-width": fmt.Sprintf("%g", c.options.Styles.StrokeWidth),
		"fill":         c.options.Styles.Fill,
		"fill-weight":  fmt.Sprintf("%g", c.options.Styles.FillWeight),
	}
}

func (c curve) Styles() *Styles {
	return c.options.Styles
}

func NewCurve(points []Point, opt *EllipseOptions) *curve {
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
	outline := curveOperation(points, opt)

	if opt.Styles.Fill != None && len(points) >= 3 {
		// Ошибку не проверяем, тк уже проверили количество точек
		bcurve, _ := CurveToBezier(points, 0)
		polyPoints := PointsOnBezierCurves(bcurve, 10, (1+opt.PenOptions.Roughness)/2)

		if opt.Styles.Filler == nil {
			operations = append(operations, solidFillPolygon(polyPoints, opt.PenOptions))

		} else {
			operations = append(operations, patternFillPolygon(polyPoints, &LineOptions{
				PenOptions: opt.PenOptions,
				Styles:     opt.Styles,
			}))
		}
	}

	if opt.Styles.Stroke != None {
		operations = append(operations, outline)
	}

	c := &curve{
		options:    opt,
		operations: operations,
	}

	return c
}
