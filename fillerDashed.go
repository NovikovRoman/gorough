package gorough

import "math"

type dashedFiller struct {
	connectEnds  bool
	hachureAngle float64
	hachureGap   float64
	dashOffset   float64
}

func NewDashedFiller(dashOffset ...float64) Filler {
	d := float64(-1)
	if len(dashOffset) > 0 {
		d = dashOffset[0]
	}
	return &dashedFiller{
		connectEnds:  false,
		hachureAngle: 0,
		hachureGap:   4,
		dashOffset:   d,
	}
}

func (f *dashedFiller) SetAngle(angle float64) {
	f.hachureAngle = angle
}

func (f *dashedFiller) SetGap(gap float64) {
	f.hachureGap = gap
}

func (f *dashedFiller) setConnectEnds(b bool) {
	f.connectEnds = b
}

func (f dashedFiller) fillPolygon(points []Point, opt *LineOptions) operation {
	lines := polygonHachureLines(points, f.hachureAngle, f.hachureGap, opt)
	return operation{
		code:     operationFillSketch,
		commands: f.dashedLine(lines, opt),
	}
}

func (f dashedFiller) dashedLine(lines []Line, opt *LineOptions) (commands []command) {
	offset := f.dashOffset
	if offset < 0 {
		offset = initHachureGap(f.hachureGap, opt.Styles.StrokeWidth)
	}
	gap := initHachureGap(f.hachureGap, opt.Styles.StrokeWidth)

	commands = []command{}
	for _, line := range lines {
		length := line.length()
		count := math.Floor(length / (offset + gap))
		startOffset := (length + gap - (count * (offset + gap))) / 2
		p1 := line.P1
		p2 := line.P2
		if p1.X > p2.X {
			p1, p2 = p2, p1
		}

		alpha := math.Atan((p2.Y - p1.Y) / (p2.X - p1.X))
		for i := 0.0; i < count; i++ {
			lstart := i * (offset + gap)
			lend := lstart + offset
			start := Point{
				X: p1.X + (lstart * math.Cos(alpha)) + (startOffset * math.Cos(alpha)),
				Y: p1.Y + lstart*math.Sin(alpha) + (startOffset * math.Sin(alpha)),
			}

			end := Point{
				X: p1.X + (lend * math.Cos(alpha)) + (startOffset * math.Cos(alpha)),
				Y: p1.Y + (lend * math.Sin(alpha)) + (startOffset * math.Sin(alpha)),
			}

			commands = append(commands, doubleLine(start, end, opt.PenOptions)...)
		}
	}
	return
}
