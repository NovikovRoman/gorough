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

func (f dotFiller) fillPolygon(points []Point, style Style, pen Pen, filler Filler) operation {
	pen.Roughness = 1
	return f.dotsOnLines(
		polygonHachureLines(points, f.hachureAngle, f.hachureGap, style.StrokeWidth), filler, style, pen)
}

func (f dotFiller) dotsOnLines(lines []Line, filler Filler, style Style, pen Pen) operation {
	var commands []command
	commands = []command{}
	gap := initHachureGap(f.hachureGap, style.StrokeWidth)
	gap = math.Max(gap, 0.1)

	fweight := style.FillWeight
	if fweight < 0 {
		fweight = style.StrokeWidth / 2
	}
	ro := gap / 4

	curveOpt := CurveDefault()
	curveOpt.StepCount = 4

	for _, line := range lines {
		length := line.length()
		dl := length / gap
		count := math.Ceil(dl) - 1
		offset := length - (count * gap)
		x := ((line.P1.X + line.P2.X) / 2) - (gap / 4)
		minY := math.Min(line.P1.Y, line.P2.Y)

		for i := float64(0); i < count; i++ {
			y := minY + offset + (i * gap)
			cx := randOffsetWithRange(x-ro, x+ro, pen.Roughness)
			cy := randOffsetWithRange(y-ro, y+ro, pen.Roughness)
			el := ellipseOperations(Point{
				X: cx,
				Y: cy,
			}, fweight, fweight, curveOpt, style, pen, filler)

			for _, e := range el {
				commands = append(commands, e.commands...)
			}
		}
	}

	return operationFillSketch(commands)
}
