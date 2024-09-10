package gorough

import (
	"math"

	"github.com/NovikovRoman/gorough/internal/tools"
	svg "github.com/NovikovRoman/svgo"
)

func DrawSVG(s *svg.SVG, d Drawable, groupAttrs Attributes) {
	if len(d.Operations()) == 0 {
		return
	}

	s.Group(groupAttrs.String())
	for _, op := range d.Operations() {
		switch op.code {
		case operationPath:
			attr := d.Attributes().Exclude("fill")
			s.Path(operationToPath(op), "fill='none'", attr.String())

		case operationFillPath:
			attr := d.Attributes().Exclude("stroke", "stroke-width", "fill")
			attrFill := ""
			if d.Filler() == nil {
				f := d.Style().Fill
				if f == "" {
					f = "none"
				}
				attrFill = "fill='" + f + "'"
			}
			s.Path(operationToPath(op), attr.String(), "stroke='none'", "stroke-width='0'", attrFill)

		case operationFillSketch:
			fweight := d.Style().FillWeight
			if fweight <= 0 {
				fweight = d.Style().StrokeWidth / 2
			}

			attr := d.Attributes().Exclude("fill", "stroke-width", "stroke")
			s.Path(
				operationToPath(op),
				"stroke='"+d.Style().Fill+"'",
				"stroke-width='"+tools.FloatToString(fweight)+"'",
				"fill='none'",
				attr.String(),
			)
		}
	}
	s.Gend()
}

func svgPath(path string, pen Pen) (op operation, err error) {
	var (
		res      []tools.Segment
		commands []command
	)
	if res, err = tools.ParsePath(path); err != nil {
		return
	}

	segments := tools.Normalize(tools.Absolutize(res))
	first := Point{}
	current := Point{}
	commands = []command{}

	for _, s := range segments {
		switch s.Key {
		case "M":
			dd := make([]float64, len(s.Data))
			for i, v := range s.Data {
				dd[i] = v + offsetOpt(pen.MaxRandomnessOffset, pen.Roughness, 1)
			}
			commands = append(commands, command{
				code: commandMove,
				data: dd,
			})
			current = Point{X: s.Data[0], Y: s.Data[1]}
			first = Point{X: s.Data[0], Y: s.Data[1]}

		case "L":
			commands = append(commands, doubleLine(current, Point{X: s.Data[0], Y: s.Data[1]}, pen)...)
			current = Point{X: s.Data[0], Y: s.Data[1]}

		case "C":
			commands = append(commands, bezierTo(s.Data[0], s.Data[1], s.Data[2], s.Data[3], s.Data[4], s.Data[5], current, pen)...)
			current = Point{X: s.Data[4], Y: s.Data[5]}

		case "Z":
			commands = append(commands, doubleLine(current, Point{X: first.X, Y: first.Y}, pen)...)
			current = Point{X: first.X, Y: first.Y}
		}
	}

	op = operation{
		code:     operationPath,
		commands: commands,
	}
	return
}

func bezierTo(x1, y1, x2, y2, x, y float64, current Point, pen Pen) (commands []command) {
	maxRandomnessOffset := pen.MaxRandomnessOffset
	if maxRandomnessOffset == 0 {
		maxRandomnessOffset = 1
	}
	ros := Point{
		X: maxRandomnessOffset,
		Y: maxRandomnessOffset + 0.3,
	}

	commands = append(commands, command{
		code: commandMove,
		data: []float64{current.X, current.Y},
	})

	f := Point{
		X: x + offsetOpt(ros.X, pen.Roughness, 1),
		Y: y + offsetOpt(ros.X, pen.Roughness, 1),
	}

	commands = append(commands, command{
		code: commandCurveTo,
		data: []float64{
			x1 + offsetOpt(ros.X, pen.Roughness, 1),
			y1 + offsetOpt(ros.X, pen.Roughness, 1),
			x2 + offsetOpt(ros.X, pen.Roughness, 1),
			y2 + offsetOpt(ros.X, pen.Roughness, 1),
			f.X, f.Y,
		},
	})

	commands = append(commands, command{
		code: commandMove,
		data: []float64{
			current.X + offsetOpt(ros.X, pen.Roughness, 1),
			current.Y + offsetOpt(ros.X, pen.Roughness, 1),
		},
	})

	f.X = x + offsetOpt(ros.Y, pen.Roughness, 1)
	f.Y = y + offsetOpt(ros.Y, pen.Roughness, 1)

	commands = append(commands, command{
		code: commandCurveTo,
		data: []float64{
			x1 + offsetOpt(ros.Y, pen.Roughness, 1),
			y1 + offsetOpt(ros.Y, pen.Roughness, 1),
			x2 + offsetOpt(ros.Y, pen.Roughness, 1),
			y2 + offsetOpt(ros.Y, pen.Roughness, 1),
			f.X, f.Y,
		},
	})

	return
}

func solidFillPolygon(points []Point, pen Pen) operation {
	var commands []command

	if len(points) > 2 {
		offset := pen.MaxRandomnessOffset
		commands = make([]command, 0, len(points))
		commands = append(commands, command{
			code: commandMove,
			data: []float64{
				points[0].X + offsetOpt(offset, pen.Roughness, 1),
				points[0].Y + offsetOpt(offset, pen.Roughness, 1),
			},
		})

		for i := 1; i < len(points); i++ {
			commands = append(commands, command{
				code: commandLineTo,
				data: []float64{
					points[i].X + offsetOpt(offset, pen.Roughness, 1),
					points[i].Y + offsetOpt(offset, pen.Roughness, 1),
				},
			})
		}
	}

	return operation{
		code:     operationFillPath,
		commands: commands,
	}
}

func patternFillPolygon(points []Point, style Style, pen Pen, filler Filler) (op operation) {
	return filler.fillPolygon(points, style, pen, filler)
}

func patternFillArc(center Point, width, height, start, stop float64, curveOpt CurveOption, style Style, pen Pen, filler Filler) operation {
	var points []Point
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

	increment := (stop - start) / curveOpt.StepCount
	points = []Point{}
	for angle := start; angle <= stop; angle = angle + increment {
		points = append(points, Point{
			X: center.X + rx*math.Cos(angle),
			Y: center.Y + ry*math.Sin(angle),
		})
	}

	points = append(points, Point{
		X: center.X + rx*math.Cos(stop),
		Y: center.Y + ry*math.Sin(stop),
	}, center)

	return patternFillPolygon(points, style, pen, filler)
}
