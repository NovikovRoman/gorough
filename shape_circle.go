package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type circleOpt func(*circle)

func CirclePen(pen Pen) circleOpt {
	return func(c *circle) {
		c.pen = pen
	}
}

func CircleStyle(s Style) circleOpt {
	return func(c *circle) {
		c.style = s
	}
}

func CircleFiller(f Filler) circleOpt {
	return func(c *circle) {
		c.filler = f
	}
}

func CircleCurveOpt(opt CurveOption) circleOpt {
	return func(c *circle) {
		c.curveOpt = opt
	}
}

type circle struct {
	pen        Pen
	style      Style
	filler     Filler
	curveOpt   CurveOption
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
		"stroke":       c.style.Stroke,
		"stroke-width": tools.FloatToString(c.style.StrokeWidth),
		"fill":         c.style.Fill,
		"fill-weight":  tools.FloatToString(c.style.FillWeight),
	}
}

func (c circle) Style() Style {
	return c.style
}

func (c circle) Filler() Filler {
	return c.filler
}

func NewCircle(center Point, diameter float64, opts ...circleOpt) (c circle) {
	c = circle{
		pen:        PenDefault(),
		style:      StyleDefault(),
		curveOpt:   CurveDefault(),
		operations: []operation{},
	}

	for _, opt := range opts {
		opt(&c)
	}

	c.style = c.style.canonicalValues()
	c.operations = ellipseOperations(center, diameter, diameter, c.curveOpt, c.style, c.pen, c.filler)
	return
}
