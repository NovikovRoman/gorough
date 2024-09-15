package gorough

import "math"

func arcOperation(center Point, width, height, start, stop float64, closed bool, roughClosure bool, curveOpt CurveOption, pen Pen) operation {
	rx := math.Abs(width / 2)
	ry := math.Abs(height / 2)
	rx += offsetOpt(rx*0.01, pen.Roughness, 1)
	ry += offsetOpt(ry*0.01, pen.Roughness, 1)

	for start < 0 {
		start += math.Pi * 2
		stop += math.Pi * 2
	}

	if stop-start > math.Pi*2 {
		start = 0
		stop = math.Pi * 2
	}
	ellipseInc := (math.Pi * 2) * curveOpt.StepCount
	arcInc := math.Min(ellipseInc/2, (stop-start)/2)
	commands := arcCommands(center, arcInc, rx, ry, start, stop, 1, curveOpt, pen)
	commands = append(commands, arcCommands(center, arcInc, rx, ry, start, stop, 1.5, curveOpt, pen)...)

	if closed {
		if roughClosure {
			commands = append(commands, doubleLine(center, Point{
				X: center.X + rx*math.Cos(start),
				Y: center.Y + ry*math.Sin(start),
			}, pen)...)

			commands = append(commands, doubleLine(center, Point{
				X: center.X + rx*math.Cos(stop),
				Y: center.Y + ry*math.Sin(stop),
			}, pen)...)

		} else {
			commands = append(commands,
				commandLineTo([]float64{center.X, center.Y}),
				commandLineTo([]float64{center.X + rx*math.Cos(start), center.Y + ry*math.Sin(start)}),
			)
		}
	}

	return operationPath(commands)
}

func arcCommands(center Point, increment, rx, ry, start, stop, offset float64, curveOpt CurveOption, pen Pen) []command {
	var points []Point
	radOffset := start + offsetOpt(0.1, pen.Roughness, 1)
	points = []Point{
		{
			X: offsetOpt(offset, pen.Roughness, 1) + center.X + 0.9*rx*math.Cos(radOffset-increment),
			Y: offsetOpt(offset, pen.Roughness, 1) + center.Y + 0.9*ry*math.Cos(radOffset-increment),
		},
	}

	for angle := radOffset; angle <= stop; angle = angle + increment {
		points = append(points, Point{
			X: offsetOpt(offset, pen.Roughness, 1) + center.X + rx*math.Cos(angle),
			Y: offsetOpt(offset, pen.Roughness, 1) + center.Y + ry*math.Sin(angle),
		})
	}
	points = append(points,
		Point{
			X: center.X + rx*math.Cos(stop),
			Y: center.Y + ry*math.Sin(stop),
		},
		Point{
			X: center.X + rx*math.Cos(stop),
			Y: center.Y + ry*math.Sin(stop),
		},
	)
	return curveCommands(points, nil, curveOpt, pen)
}
