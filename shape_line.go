package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type lineOpt func(*line)

func LinePen(p Pen) lineOpt {
	return func(l *line) {
		l.pen = p
	}
}

func LineStyle(s Style) lineOpt {
	return func(l *line) {
		l.style = s
	}
}

type line struct {
	pen        Pen
	style      Style
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
		"stroke":       l.style.Stroke,
		"stroke-width": tools.FloatToString(l.style.StrokeWidth),
	}
}

func (l line) Style() Style {
	return l.style
}

func (l line) Filler() Filler {
	return nil
}

func NewLine(p1, p2 Point, opts ...lineOpt) line {
	l := line{
		pen:   PenDefault(),
		style: StyleDefault(),
	}

	for _, opt := range opts {
		opt(&l)
	}

	l.style = l.style.canonicalValues()
	l.operations = lineOperations(p1, p2, l.pen)
	return l
}
