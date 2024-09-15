package gorough

import "github.com/NovikovRoman/gorough/internal/tools"

// command types
const (
	cmdMove    = "move"
	cmdCurveTo = "curveTo"
	cmdLineTo  = "lineTo"
)

type command struct {
	code string
	data []float64
}

func commandMove(data []float64) command {
	return command{code: cmdMove, data: data}
}

func commandCurveTo(data []float64) command {
	return command{code: cmdCurveTo, data: data}
}

func commandLineTo(data []float64) command {
	return command{code: cmdLineTo, data: data}
}

func (c command) String() string {
	switch c.code {
	case cmdMove:
		return "M" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1])

	case cmdCurveTo:
		return "C" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1]) + ", " +
			tools.FloatToString(c.data[2]) + " " + tools.FloatToString(c.data[3]) + ", " +
			tools.FloatToString(c.data[4]) + " " + tools.FloatToString(c.data[5])

	case cmdLineTo:
		return "L" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1])

	default:
		return ""
	}
}
