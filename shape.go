package gorough

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	None            = "none"
	shapeLine       = "line"
	shapeRectangle  = "rectangle"
	shapeEllipse    = "ellipse"
	shapeLinearPath = "linear path"
	shapePoligon    = "poligon"
	shapeCurve      = "curve"

	// operation types
	operationPath       = "path"
	operationFillPath   = "fillPath"
	operationFillSketch = "fillSketch"

	// command types
	commandMove    = "move"
	commandCurveTo = "curveTo"
	commandLineTo  = "lineTo"
)

func lineOperations(p1 Point, p2 Point, opt *PenOptions) []operation {
	return []operation{{code: operationPath, commands: doubleLine(p1, p2, opt)}}
}

func doubleLine(p1 Point, p2 Point, opt *PenOptions) (commands []command) {
	commands = oneLine(p1, p2, true, false, opt)
	commands = append(commands, oneLine(p1, p2, true, true, opt)...)
	return
}

func oneLine(p1 Point, p2 Point, move bool, overlay bool, opt *PenOptions) (commands []command) {
	lengthSq := math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
	length := math.Sqrt(lengthSq)
	roughnessGain := float64(0)
	if length < 200 {
		roughnessGain = 1
	} else if length > 500 {
		roughnessGain = 0.4
	} else {
		roughnessGain = (-0.0016668)*length + 1.233334
	}

	offset := opt.MaxRandomnessOffset
	if offset*offset*100 > lengthSq {
		offset = length / 10
	}

	halfOffset := offset / 2
	divergePoint := 0.2 + rand.Float64()*0.2
	midDispX := opt.Bowing * opt.MaxRandomnessOffset * (p2.Y - p1.Y) / 200
	midDispY := opt.Bowing * opt.MaxRandomnessOffset * (p1.X - p2.X) / 200
	midDispX = offsetOpt(midDispX, opt.Roughness, roughnessGain)
	midDispY = offsetOpt(midDispY, opt.Roughness, roughnessGain)

	if move && overlay {
		commands = append(commands, command{
			code: commandMove,
			data: []float64{
				p1.X + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				p1.Y + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
			},
		})

	} else if move && !overlay {
		commands = append(commands, command{
			code: commandMove,
			data: []float64{
				p1.X + offsetOpt(offset, opt.Roughness, roughnessGain),
				p1.Y + offsetOpt(offset, opt.Roughness, roughnessGain),
			},
		})
	}

	if overlay {
		commands = append(commands, command{
			code: commandCurveTo,
			data: []float64{
				midDispX + p1.X + (p2.X-p1.X)*divergePoint + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				midDispY + p1.Y + (p2.Y-p1.Y)*divergePoint + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				midDispX + p1.X + 2*(p2.X-p1.X)*divergePoint + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				midDispY + p1.Y + 2*(p2.Y-p1.Y)*divergePoint + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				p2.X + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
				p2.Y + offsetOpt(halfOffset, opt.Roughness, roughnessGain),
			},
		})
		return
	}

	commands = append(commands, command{
		code: commandCurveTo,
		data: []float64{
			midDispX + p1.X + (p2.X-p1.X)*divergePoint + offsetOpt(offset, opt.Roughness, roughnessGain),
			midDispY + p1.Y + (p2.Y-p1.Y)*divergePoint + offsetOpt(offset, opt.Roughness, roughnessGain),
			midDispX + p1.X + 2*(p2.X-p1.X)*divergePoint + offsetOpt(offset, opt.Roughness, roughnessGain),
			midDispY + p1.Y + 2*(p2.Y-p1.Y)*divergePoint + offsetOpt(offset, opt.Roughness, roughnessGain),
			p2.X + offsetOpt(offset, opt.Roughness, roughnessGain),
			p2.Y + offsetOpt(offset, opt.Roughness, roughnessGain),
		},
	})
	return
}

func offsetOpt(x float64, roughness float64, roughnessGain float64) float64 {
	return offset(-x, x, roughness, roughnessGain)
}

func offset(min float64, max float64, roughness float64, roughnessGain float64) float64 {
	return roughness * roughnessGain * ((rand.Float64() * (max - min)) + min)
}

func operationToPath(op operation) string {
	path := make([]string, len(op.commands))
	for i, c := range op.commands {
		switch c.code {
		case commandMove:
			path[i] = fmt.Sprintf("M%g %g", c.data[0], c.data[1])

		case commandCurveTo:
			path[i] = fmt.Sprintf(
				"C%g %g, %g %g, %g %g", c.data[0], c.data[1], c.data[2], c.data[3], c.data[4], c.data[5])

		case commandLineTo:
			path[i] = fmt.Sprintf("L%g %g", c.data[0], c.data[1])
		}
	}
	return strings.Join(path, "")
}

func rectangleOperation(p Point, width, height float64, opt *RectangleOptions) operation {
	return polygonOperation([]Point{
		p,
		{
			X: p.X + width,
			Y: p.Y,
		},
		{
			X: p.X + width,
			Y: p.Y + height,
		},
		{
			X: p.X,
			Y: p.Y + height,
		},
	}, opt.PenOptions)
}

func polygonOperation(points []Point, opt *PenOptions) operation {
	return linearPathOperation(points, true, opt)
}

func linearPathOperation(points []Point, close bool, opt *PenOptions) operation {
	var (
		length   = len(points)
		commands []command
	)
	commands = []command{}

	if length == 2 {
		commands = doubleLine(points[0], points[1], opt)

	} else if length > 2 {
		for i := 0; i < length-1; i++ {
			commands = append(commands, doubleLine(points[i], points[i+1], opt)...)
		}
		if close {
			commands = append(commands, doubleLine(points[length-1], points[0], opt)...)
		}
	}

	return operation{
		code:     operationPath,
		commands: commands,
	}
}

func generateEllipseParams(width, height float64, opt *EllipseOptions) (increment, rx, ry float64) {
	psq := math.Sqrt(math.Pi * 2 * math.Sqrt((math.Pow(width/2, 2)+math.Pow(height/2, 2))/2))
	stepCount := math.Max(opt.CurveOptions.StepCount, (opt.CurveOptions.StepCount/math.Sqrt(200))*psq)
	increment = (math.Pi * 2) / stepCount
	rx = math.Abs(width / 2)
	ry = math.Abs(height / 2)
	curveFitRandomness := 1 - opt.CurveOptions.Fitting
	rx += offsetOpt(rx*curveFitRandomness, opt.PenOptions.Roughness, 1)
	ry += offsetOpt(ry*curveFitRandomness, opt.PenOptions.Roughness, 1)
	return
}

func ellipseWithParams(p Point, increment, rx, ry float64, opt *EllipseOptions) (estimatedPoints []Point, op operation) {
	var ap1 []Point
	ap1, estimatedPoints = computeEllipsePoints(p, increment, rx, ry, 1, increment*offset(0.1, offset(0.4, 1, opt.PenOptions.Roughness, 1), opt.PenOptions.Roughness, 1), opt)
	ap2, _ := computeEllipsePoints(p, increment, rx, ry, 1.5, 0, opt)
	commands := curveCommands(ap1, nil, opt)
	commands = append(commands, curveCommands(ap2, nil, opt)...)

	op = operation{
		code:     operationPath,
		commands: commands,
	}
	return
}

func computeEllipsePoints(p Point, increment, rx, ry float64, offset float64, overlap float64, opt *EllipseOptions) (allPoints []Point, corePoints []Point) {
	radOffset := offsetOpt(0.5, opt.PenOptions.Roughness, 1) - (math.Pi / 2)

	allPoints = []Point{
		{
			X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.X + 0.9*rx*math.Cos(radOffset-increment),
			Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.Y + 0.9*ry*math.Sin(radOffset-increment),
		},
	}

	for angle := radOffset; angle < (math.Pi*2 + radOffset - 0.01); angle = angle + increment {
		pp := Point{
			X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.X + rx*math.Cos(angle),
			Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.Y + ry*math.Sin(angle),
		}
		allPoints = append(allPoints, pp)
		corePoints = append(corePoints, pp)
	}

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.X + rx*math.Cos(radOffset+math.Pi*2+overlap*0.5),
		Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.Y + ry*math.Sin(radOffset+math.Pi*2+overlap*0.5),
	})

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.X + 0.98*rx*math.Cos(radOffset+overlap),
		Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.Y + 0.98*ry*math.Sin(radOffset+overlap),
	})

	allPoints = append(allPoints, Point{
		X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.X + 0.9*rx*math.Cos(radOffset+overlap*0.5),
		Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + p.Y + 0.9*ry*math.Sin(radOffset+overlap*0.5),
	})
	return
}

func curveCommands(points []Point, closePoint *Point, opt *EllipseOptions) (commands []command) {
	length := len(points)
	commands = []command{}
	if length < 2 {
		return
	}

	if length == 2 {
		return doubleLine(points[0], points[1], opt.PenOptions)

	}

	if length == 3 {
		commands = append(commands, command{
			code: commandMove,
			data: []float64{points[1].X, points[1].Y},
		})

		commands = append(commands, command{
			code: commandCurveTo,
			data: []float64{
				points[1].X, points[1].Y,
				points[2].X, points[2].Y,
				points[2].X, points[2].Y,
			},
		})

	} else {
		b := make([]Point, 4)
		s := 1 - opt.CurveOptions.Tightness
		commands = append(commands, command{
			code: commandMove,
			data: []float64{points[1].X, points[1].Y},
		})

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

			commands = append(commands, command{
				code: commandCurveTo,
				data: []float64{b[1].X, b[1].Y, b[2].X, b[2].Y, b[3].X, b[3].Y},
			})
		}

		if closePoint != nil {
			ro := opt.PenOptions.MaxRandomnessOffset
			commands = append(commands, command{
				code: commandLineTo,
				data: []float64{
					closePoint.X + offsetOpt(ro, opt.PenOptions.Roughness, 1),
					closePoint.Y + offsetOpt(ro, opt.PenOptions.Roughness, 1),
				},
			})
		}
	}

	return
}

func curveOperation(points []Point, opt *EllipseOptions) operation {
	commands := curveWithOffset(points, 1+opt.PenOptions.Roughness*0.2, opt)
	commands = append(commands, curveWithOffset(points, 1.5+opt.PenOptions.Roughness*0.22, opt)...)
	return operation{
		code:     operationPath,
		commands: commands,
	}
}

func curveWithOffset(points []Point, offset float64, opt *EllipseOptions) []command {
	var ps []Point
	ps = []Point{
		{
			X: points[0].X + offsetOpt(offset, opt.PenOptions.Roughness, 1),
			Y: points[0].Y + offsetOpt(offset, opt.PenOptions.Roughness, 1),
		},
		{
			X: points[0].X + offsetOpt(offset, opt.PenOptions.Roughness, 1),
			Y: points[0].Y + offsetOpt(offset, opt.PenOptions.Roughness, 1),
		},
	}

	for i := 1; i < len(points); i++ {
		ps = append(ps, Point{
			X: points[i].X + offsetOpt(offset, opt.PenOptions.Roughness, 1),
			Y: points[i].Y + offsetOpt(offset, opt.PenOptions.Roughness, 1),
		})

		if i == len(points)-1 {
			ps = append(ps, Point{
				X: points[i].X + offsetOpt(offset, opt.PenOptions.Roughness, 1),
				Y: points[i].Y + offsetOpt(offset, opt.PenOptions.Roughness, 1),
			})
		}
	}

	return curveCommands(ps, nil, opt)
}

func arcOperation(center Point, width, height, start, stop float64, closed bool, roughClosure bool, opt *EllipseOptions) operation {
	rx := math.Abs(width / 2)
	ry := math.Abs(height / 2)
	rx += offsetOpt(rx*0.01, opt.PenOptions.Roughness, 1)
	ry += offsetOpt(ry*0.01, opt.PenOptions.Roughness, 1)

	for start < 0 {
		start += math.Pi * 2
		stop += math.Pi * 2
	}

	if stop-start > math.Pi*2 {
		start = 0
		stop = math.Pi * 2
	}
	ellipseInc := (math.Pi * 2) * opt.CurveOptions.StepCount
	arcInc := math.Min(ellipseInc/2, (stop-start)/2)
	commands := arcCommands(center, arcInc, rx, ry, start, stop, 1, opt)
	commands = append(commands, arcCommands(center, arcInc, rx, ry, start, stop, 1.5, opt)...)

	if closed {
		if roughClosure {
			commands = append(commands, doubleLine(center, Point{
				X: center.X + rx*math.Cos(start),
				Y: center.Y + ry*math.Sin(start),
			}, opt.PenOptions)...)

			commands = append(commands, doubleLine(center, Point{
				X: center.X + rx*math.Cos(stop),
				Y: center.Y + ry*math.Sin(stop),
			}, opt.PenOptions)...)

		} else {
			commands = append(commands,
				command{
					code: commandLineTo,
					data: []float64{center.X, center.Y},
				},
				command{
					code: commandLineTo,
					data: []float64{center.X + rx*math.Cos(start), center.Y + ry*math.Sin(start)},
				},
			)
		}
	}

	return operation{
		code:     operationPath,
		commands: commands,
	}
}

func arcCommands(center Point, increment, rx, ry, start, stop, offset float64, opt *EllipseOptions) []command {
	var points []Point
	radOffset := start + offsetOpt(0.1, opt.PenOptions.Roughness, 1)
	points = []Point{
		{
			X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + center.X + 0.9*rx*math.Cos(radOffset-increment),
			Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + center.Y + 0.9*ry*math.Cos(radOffset-increment),
		},
	}

	for angle := radOffset; angle <= stop; angle = angle + increment {
		points = append(points, Point{
			X: offsetOpt(offset, opt.PenOptions.Roughness, 1) + center.X + rx*math.Cos(angle),
			Y: offsetOpt(offset, opt.PenOptions.Roughness, 1) + center.Y + ry*math.Sin(angle),
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
	return curveCommands(points, nil, opt)
}

func ellipseOperations(center Point, width, height float64, opt *EllipseOptions) []operation {
	var operations []operation
	operations = []operation{}

	increment, rx, ry := generateEllipseParams(width, height, opt)
	estimatedPoints, op := ellipseWithParams(center, increment, rx, ry, opt)

	if opt.Styles.Fill != "" {
		if opt.Styles.Filler == nil {
			_, s := ellipseWithParams(center, increment, rx, ry, opt)
			s.code = operationFillPath
			operations = append(operations, s)

		} else {
			operations = append(operations, patternFillPolygon(estimatedPoints, &LineOptions{
				PenOptions: opt.PenOptions,
				Styles:     opt.Styles,
			}))
		}
	}

	if opt.Styles.Stroke != None {
		operations = append(operations, op)
	}
	return operations
}
