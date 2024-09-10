package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

type linearPathOpt func(*linearPath)

func LinearPathPen(pen Pen) linearPathOpt {
	return func(l *linearPath) {
		l.pen = pen
	}
}

func LinearPathStyle(s Style) linearPathOpt {
	return func(l *linearPath) {
		l.style = s
	}
}

type linearPath struct {
	pen        Pen
	style      Style
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
		"stroke":       l.style.Stroke,
		"stroke-width": tools.FloatToString(l.style.StrokeWidth),
	}
}

func (l linearPath) Style() Style {
	return l.style
}

func (l linearPath) Filler() Filler {
	return nil
}

func NewLinearPath(points []Point, opts ...linearPathOpt) linearPath {
	l := linearPath{
		pen:   PenDefault(),
		style: StyleDefault(),
	}

	for _, opt := range opts {
		opt(&l)
	}

	l.style = l.style.canonicalValues()
	l.operations = []operation{linearPathOperation(points, false, l.pen)}
	return l
}
