package gorough

import (
	"fmt"
	"github.com/NovikovRoman/gorough/data_parser"
	"github.com/ajstarks/svgo"
	"math"
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
			if d.Styles().Filler == nil {
				if d.Styles().Fill == "" {
					attrFill = "fill='none'"
				} else {
					attrFill = "fill='" + d.Styles().Fill + "'"
				}
			}
			s.Path(operationToPath(op), attr.String(), "stroke='none'", "stroke-width='0'", attrFill)

		case operationFillSketch:
			fweight := d.Styles().FillWeight
			if fweight <= 0 {
				fweight = d.Styles().StrokeWidth / 2
			}

			attr := d.Attributes().Exclude("fill", "stroke-width", "stroke")
			s.Path(
				operationToPath(op),
				"stroke='"+d.Styles().Fill+"'",
				"stroke-width='"+fmt.Sprintf("%g", fweight)+"'",
				"fill='none'",
				attr.String(),
			)
		}
	}
	s.Gend()
}

func svgPath(path string, opt *PenOptions) (op operation, err error) {
	var (
		res      []data_parser.Segment
		commands []command
	)
	if res, err = data_parser.ParsePath(path); err != nil {
		return
	}

	segments := data_parser.Normalize(data_parser.Absolutize(res))
	first := Point{}
	current := Point{}
	commands = []command{}

	for _, s := range segments {
		switch s.Key {
		case "M":
			dd := make([]float64, len(s.Data))
			for i, v := range s.Data {
				dd[i] = v + offsetOpt(opt.MaxRandomnessOffset, opt.Roughness, 1)
			}
			commands = append(commands, command{
				code: commandMove,
				data: dd,
			})
			current = Point{X: s.Data[0], Y: s.Data[1]}
			first = Point{X: s.Data[0], Y: s.Data[1]}

		case "L":
			commands = append(commands, doubleLine(current, Point{X: s.Data[0], Y: s.Data[1]}, opt)...)
			current = Point{X: s.Data[0], Y: s.Data[1]}

		case "C":
			commands = append(commands, bezierTo(s.Data[0], s.Data[1], s.Data[2], s.Data[3], s.Data[4], s.Data[5], current, opt)...)
			current = Point{X: s.Data[4], Y: s.Data[5]}

		case "Z":
			commands = append(commands, doubleLine(current, Point{X: first.X, Y: first.Y}, opt)...)
			current = Point{X: first.X, Y: first.Y}

		}
	}

	op = operation{
		code:     operationPath,
		commands: commands,
	}
	return
}

func bezierTo(x1, y1, x2, y2, x, y float64, current Point, opt *PenOptions) (commands []command) {
	maxRandomnessOffset := opt.MaxRandomnessOffset
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
		X: x + offsetOpt(ros.X, opt.Roughness, 1),
		Y: y + offsetOpt(ros.X, opt.Roughness, 1),
	}

	commands = append(commands, command{
		code: commandCurveTo,
		data: []float64{
			x1 + offsetOpt(ros.X, opt.Roughness, 1),
			y1 + offsetOpt(ros.X, opt.Roughness, 1),
			x2 + offsetOpt(ros.X, opt.Roughness, 1),
			y2 + offsetOpt(ros.X, opt.Roughness, 1),
			f.X, f.Y,
		},
	})

	commands = append(commands, command{
		code: commandMove,
		data: []float64{
			current.X + offsetOpt(ros.X, opt.Roughness, 1),
			current.Y + offsetOpt(ros.X, opt.Roughness, 1),
		},
	})

	f.X = x + offsetOpt(ros.Y, opt.Roughness, 1)
	f.Y = y + offsetOpt(ros.Y, opt.Roughness, 1)

	commands = append(commands, command{
		code: commandCurveTo,
		data: []float64{
			x1 + offsetOpt(ros.Y, opt.Roughness, 1),
			y1 + offsetOpt(ros.Y, opt.Roughness, 1),
			x2 + offsetOpt(ros.Y, opt.Roughness, 1),
			y2 + offsetOpt(ros.Y, opt.Roughness, 1),
			f.X, f.Y,
		},
	})

	return
}

func solidFillPolygon(points []Point, opt *PenOptions) operation {
	var commands []command
	commands = []command{}
	length := len(points)

	if length > 2 {
		offset := opt.MaxRandomnessOffset
		commands = append(commands, command{
			code: commandMove,
			data: []float64{
				points[0].X + offsetOpt(offset, opt.Roughness, 1),
				points[0].Y + offsetOpt(offset, opt.Roughness, 1),
			},
		})

		for i := 1; i < length; i++ {
			commands = append(commands, command{
				code: commandLineTo,
				data: []float64{
					points[i].X + offsetOpt(offset, opt.Roughness, 1),
					points[i].Y + offsetOpt(offset, opt.Roughness, 1),
				},
			})
		}
	}

	return operation{
		code:     operationFillPath,
		commands: commands,
	}
}

func patternFillPolygon(points []Point, opt *LineOptions) (op operation) {
	op = operation{
		code:     operationPath,
		commands: []command{},
	}
	if opt.Styles == nil || opt.Styles.Filler == nil {
		return
	}
	return opt.Styles.Filler.fillPolygon(points, opt)
}

func patternFillArc(center Point, width, height, start, stop float64, opt *EllipseOptions) operation {
	var points []Point
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

	increment := (stop - start) / opt.CurveOptions.StepCount
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

	return patternFillPolygon(points, &LineOptions{
		PenOptions: opt.PenOptions,
		Styles:     opt.Styles,
	})
}
