package gorough

import "math"

type dotFiller struct {
	connectEnds  bool
	hachureAngle float64
	hachureGap   float64
}

func NewDotFiller() Filler {
	return &dotFiller{
		connectEnds:  false,
		hachureAngle: 0,
		hachureGap:   4,
	}
}

func (f *dotFiller) SetAngle(_ float64) {
	// Do not change the angle
}

func (f *dotFiller) SetGap(gap float64) {
	f.hachureGap = gap
}

func (f *dotFiller) setConnectEnds(b bool) {
	f.connectEnds = b
}

func (f dotFiller) fillPolygon(points []Point, opt *LineOptions) operation {
	o := &EllipseOptions{
		PenOptions:   opt.PenOptions,
		CurveOptions: CurveOptionsDefault(),
		Styles:       opt.Styles,
	}
	o.CurveOptions.StepCount = 4
	o.PenOptions.Roughness = 1

	return f.dotsOnLines(polygonHachureLines(points, f.hachureAngle, f.hachureGap, opt), o)
}

func (f dotFiller) dotsOnLines(lines []Line, opt *EllipseOptions) operation {
	var commands []command
	commands = []command{}
	gap := initHachureGap(f.hachureGap, opt.Styles.StrokeWidth)
	gap = math.Max(gap, 0.1)

	fweight := opt.Styles.FillWeight
	if fweight < 0 {
		fweight = opt.Styles.StrokeWidth / 2
	}
	ro := gap / 4

	for _, line := range lines {
		length := line.length()
		dl := length / gap
		count := math.Ceil(dl) - 1
		offset := length - (count * gap)
		x := ((line.P1.X + line.P2.X) / 2) - (gap / 4)
		minY := math.Min(line.P1.Y, line.P2.Y)

		for i := float64(0); i < count; i++ {
			y := minY + offset + (i * gap)
			cx := randOffsetWithRange(x-ro, x+ro, opt.PenOptions)
			cy := randOffsetWithRange(y-ro, y+ro, opt.PenOptions)
			el := ellipseOperations(Point{
				X: cx,
				Y: cy,
			}, fweight, fweight, opt)

			for _, e := range el {
				commands = append(commands, e.commands...)
			}
		}
	}

	return operation{
		code:     operationFillSketch,
		commands: commands,
	}
}
