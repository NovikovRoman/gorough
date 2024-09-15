package gorough

import "math"

func ellipseOperations(center Point, width, height float64,
	curveOpt CurveOption, style Style, pen Pen, filler Filler) (operations []operation) {

	increment, rx, ry := generateEllipseParams(width, height, curveOpt, pen)
	estimatedPoints, op := ellipseWithParams(center, increment, rx, ry, curveOpt, pen)

	if style.Fill != "" {
		if filler == nil {
			_, s := ellipseWithParams(center, increment, rx, ry, curveOpt, pen)
			s.code = opFillPath
			operations = append(operations, s)

		} else {
			operations = append(operations, patternFillPolygon(estimatedPoints, style, pen, filler))
		}
	}

	if style.Stroke != none {
		operations = append(operations, op)
	}
	return
}

func generateEllipseParams(width, height float64, curveOpt CurveOption, pen Pen) (increment, rx, ry float64) {
	psq := math.Sqrt(math.Pi * 2 * math.Sqrt((math.Pow(width/2, 2)+math.Pow(height/2, 2))/2))
	stepCount := math.Max(curveOpt.StepCount, (curveOpt.StepCount/math.Sqrt(200))*psq)
	increment = (math.Pi * 2) / stepCount
	rx = math.Abs(width / 2)
	ry = math.Abs(height / 2)
	curveFitRandomness := 1 - curveOpt.Fitting
	rx += offsetOpt(rx*curveFitRandomness, pen.Roughness, 1)
	ry += offsetOpt(ry*curveFitRandomness, pen.Roughness, 1)
	return
}

func ellipseWithParams(p Point, increment, rx, ry float64, curveOpt CurveOption, pen Pen) (estimatedPoints []Point, op operation) {
	var ap1 []Point
	ap1, estimatedPoints = computeEllipsePoints(
		p, increment, rx, ry, 1, increment*offset(0.1, offset(0.4, 1, pen.Roughness, 1), pen.Roughness, 1), pen)
	ap2, _ := computeEllipsePoints(p, increment, rx, ry, 1.5, 0, pen)
	commands := curveCommands(ap1, nil, curveOpt, pen)
	commands = append(commands, curveCommands(ap2, nil, curveOpt, pen)...)

	op = operationPath(commands)
	return
}

func computeEllipsePoints(p Point, increment, rx, ry float64, offset float64, overlap float64, pen Pen) (allPoints []Point, corePoints []Point) {
	radOffset := offsetOpt(0.5, pen.Roughness, 1) - (math.Pi / 2)

	allPoints = []Point{
		{
			X: offsetOpt(offset, pen.Roughness, 1) + p.X + 0.9*rx*math.Cos(radOffset-increment),
			Y: offsetOpt(offset, pen.Roughness, 1) + p.Y + 0.9*ry*math.Sin(radOffset-increment),
		},
	}

	for angle := radOffset; angle < (math.Pi*2 + radOffset - 0.01); angle = angle + increment {
		pp := Point{
			X: offsetOpt(offset, pen.Roughness, 1) + p.X + rx*math.Cos(angle),
			Y: offsetOpt(offset, pen.Roughness, 1) + p.Y + ry*math.Sin(angle),
		}
		allPoints = append(allPoints, pp)
		corePoints = append(corePoints, pp)
	}

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, pen.Roughness, 1) + p.X + rx*math.Cos(radOffset+math.Pi*2+overlap*0.5),
		Y: offsetOpt(offset, pen.Roughness, 1) + p.Y + ry*math.Sin(radOffset+math.Pi*2+overlap*0.5),
	})

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, pen.Roughness, 1) + p.X + 0.98*rx*math.Cos(radOffset+overlap),
		Y: offsetOpt(offset, pen.Roughness, 1) + p.Y + 0.98*ry*math.Sin(radOffset+overlap),
	})

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, pen.Roughness, 1) + p.X + 0.9*rx*math.Cos(radOffset+overlap*0.5),
		Y: offsetOpt(offset, pen.Roughness, 1) + p.Y + 0.9*ry*math.Sin(radOffset+overlap*0.5),
	})
	return
}
