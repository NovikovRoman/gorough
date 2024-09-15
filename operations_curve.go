package gorough

func curveCommands(points []Point, closePoint *Point, curveOpt CurveOption, pen Pen) (commands []command) {
	length := len(points)
	commands = []command{}
	if length < 2 {
		return
	}

	if length == 2 {
		return doubleLine(points[0], points[1], pen)
	}

	if length == 3 {
		commands = append(commands, commandMove([]float64{points[1].X, points[1].Y}))

		commands = append(commands, commandCurveTo([]float64{
			points[1].X, points[1].Y,
			points[2].X, points[2].Y,
			points[2].X, points[2].Y,
		}))
		return
	}

	b := make([]Point, 4)
	s := 1 - curveOpt.Tightness
	commands = append(commands, commandMove([]float64{points[1].X, points[1].Y}))

	for i := 1; (i + 2) < length; i++ {
		b[0] = points[i]
		b[1] = Point{
			X: points[i].X + (s*points[i+1].X-s*points[i-1].X)/6,
			Y: points[i].Y + (s*points[i+1].Y-s*points[i-1].Y)/6,
		}
		b[2] = Point{
			X: points[i+1].X + (s*points[i].X-s*points[i+2].X)/6,
			Y: points[i+1].Y + (s*points[i].Y-s*points[i+2].Y)/6,
		}
		b[3] = Point{
			X: points[i+1].X,
			Y: points[i+1].Y,
		}

		commands = append(commands, commandCurveTo([]float64{b[1].X, b[1].Y, b[2].X, b[2].Y, b[3].X, b[3].Y}))
	}

	if closePoint != nil {
		ro := pen.MaxRandomnessOffset
		commands = append(commands, commandLineTo([]float64{
			closePoint.X + offsetOpt(ro, pen.Roughness, 1),
			closePoint.Y + offsetOpt(ro, pen.Roughness, 1),
		}))
	}
	return
}

func curveOperation(points []Point, curveOpt CurveOption, pen Pen) operation {
	commands := curveWithOffset(points, 1+pen.Roughness*0.2, curveOpt, pen)
	commands = append(commands, curveWithOffset(points, 1.5+pen.Roughness*0.22, curveOpt, pen)...)
	return operationPath(commands)
}

func curveWithOffset(points []Point, offset float64, curveOpt CurveOption, pen Pen) []command {
	var ps []Point
	ps = []Point{
		{
			X: points[0].X + offsetOpt(offset, pen.Roughness, 1),
			Y: points[0].Y + offsetOpt(offset, pen.Roughness, 1),
		},
		{
			X: points[0].X + offsetOpt(offset, pen.Roughness, 1),
			Y: points[0].Y + offsetOpt(offset, pen.Roughness, 1),
		},
	}

	for i := 1; i < len(points); i++ {
		ps = append(ps, Point{
			X: points[i].X + offsetOpt(offset, pen.Roughness, 1),
			Y: points[i].Y + offsetOpt(offset, pen.Roughness, 1),
		})

		if i == len(points)-1 {
			ps = append(ps, Point{
				X: points[i].X + offsetOpt(offset, pen.Roughness, 1),
				Y: points[i].Y + offsetOpt(offset, pen.Roughness, 1),
			})
		}
	}

	return curveCommands(ps, nil, curveOpt, pen)
}
