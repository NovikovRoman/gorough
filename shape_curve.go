package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type curveOpt func(*curve)

func CurvePen(pen Pen) curveOpt {
	return func(c *curve) {
		c.pen = pen
	}
}

func CurveStyle(s Style) curveOpt {
	return func(c *curve) {
		c.style = s
	}
}

func CurveFiller(f Filler) curveOpt {
	return func(c *curve) {
		c.filler = f
	}
}

func CurveCurveOpt(opt CurveOption) curveOpt {
	return func(c *curve) {
		c.curveOpt = opt
	}
}

type curve struct {
	pen        Pen
	style      Style
	filler     Filler
	curveOpt   CurveOption
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
		"stroke":       c.style.Stroke,
		"stroke-width": tools.FloatToString(c.style.StrokeWidth),
		"fill":         c.style.Fill,
		"fill-weight":  tools.FloatToString(c.style.FillWeight),
	}
}

func (c curve) Style() Style {
	return c.style
}

func (c curve) Filler() Filler {
	return c.filler
}

func NewCurve(points []Point, opts ...curveOpt) (c curve) {
	c = curve{
		pen:        PenDefault(),
		style:      StyleDefault(),
		curveOpt:   CurveDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&c)
	}

	c.style = c.style.canonicalValues()

	outline := curveOperation(points, c.curveOpt, c.pen)

	if c.style.Fill != none && len(points) >= 3 {
		// It does not check the error, because the number of points has already been checked
		bcurve, _ := CurveToBezier(points, 0)
		polyPoints := PointsOnBezierCurves(bcurve, 10, (1+c.pen.Roughness)/2)

		if c.filler == nil {
			c.operations = append(c.operations, solidFillPolygon(polyPoints, c.pen))

		} else {
			c.operations = append(c.operations, patternFillPolygon(polyPoints, c.style, c.pen, c.filler))
		}
	}

	if c.style.Stroke != none {
		c.operations = append(c.operations, outline)
	}
	return
}
