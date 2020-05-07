package gorough

import "fmt"

type line struct {
	options    *LineOptions
	operations []operation
}

func (l line) Name() string {
	return shapeLine
}

func (l line) Operations() []operation {
	return l.operations
}

func (l line) Attributes() Attributes {
	return map[string]string{
		"stroke":       l.options.Styles.Stroke,
		"stroke-width": fmt.Sprintf("%g", l.options.Styles.StrokeWidth),
	}
}

func (l line) Styles() *Styles {
	return l.options.Styles
}

func NewLine(p1, p2 Point, opt *LineOptions) *line {
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

	l := &line{
		options:    opt,
		operations: lineOperations(p1, p2, opt.PenOptions),
	}

	return l
}
