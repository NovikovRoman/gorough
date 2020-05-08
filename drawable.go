package gorough

import (
	"github.com/NovikovRoman/gorough/data_parser"
)

type Drawable interface {
	Name() string
	Operations() []operation
	Attributes() Attributes
	Styles() *Styles
}

type operation struct {
	code     string
	commands []command
}

type command struct {
	code string
	data []float64
}

func (c command) String() string {
	switch c.code {
	case commandMove:
		return "M" + data_parser.FloatToString(c.data[0]) + " " + data_parser.FloatToString(c.data[1])

	case commandCurveTo:
		return "C" + data_parser.FloatToString(c.data[0]) + " " + data_parser.FloatToString(c.data[1]) + ", " +
			data_parser.FloatToString(c.data[2]) + " " + data_parser.FloatToString(c.data[3]) + ", " +
			data_parser.FloatToString(c.data[4]) + " " + data_parser.FloatToString(c.data[5])

	case commandLineTo:
		return "L" + data_parser.FloatToString(c.data[0]) + " " + data_parser.FloatToString(c.data[1])
	}

	return ""
}
