package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
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
		"stroke-width": data_parser.FloatToString(l.options.Styles.StrokeWidth),
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
