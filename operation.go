package gorough


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
		return "M" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1])

	case commandCurveTo:
		return "C" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1]) + ", " +
			tools.FloatToString(c.data[2]) + " " + tools.FloatToString(c.data[3]) + ", " +
			tools.FloatToString(c.data[4]) + " " + tools.FloatToString(c.data[5])

	case commandLineTo:
		return "L" + tools.FloatToString(c.data[0]) + " " + tools.FloatToString(c.data[1])

	default:
		return ""
	}
}
