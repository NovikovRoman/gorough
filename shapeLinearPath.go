package gorough

import (
	"strconv"
)

type linearPath struct {
	options    *LineOptions
	operations []operation
}

func (l linearPath) Name() string {
	return shapeLinearPath
}

func (l linearPath) Operations() []operation {
	return l.operations
}

func (l linearPath) Attributes() Attributes {
	return map[string]string{
		"stroke":       l.options.Styles.Stroke,
		"stroke-width": strconv.FormatFloat(l.options.Styles.StrokeWidth, 'f', -1, 64),
	}
}

func (l linearPath) Styles() *Styles {
	return l.options.Styles
}

func NewLinearPath(points []Point, opt *LineOptions) *linearPath {
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

	l := &linearPath{
		options:    opt,
		operations: []operation{linearPathOperation(points, false, opt.PenOptions)},
	}

	return l
}
