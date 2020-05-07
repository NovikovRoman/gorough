package gorough

type hatchFiller struct {
	connectEnds  bool
	hachureAngle float64
	hachureGap   float64
}

func NewHatchFiller() Filler {
	return &hatchFiller{
		connectEnds:  false,
		hachureAngle: -41,
		hachureGap:   5,
	}
}

func (f *hatchFiller) SetAngle(angle float64) {
	f.hachureAngle = angle
}

func (f *hatchFiller) SetGap(gap float64) {
	f.hachureGap = gap
}

func (f *hatchFiller) setConnectEnds(b bool) {
	f.connectEnds = b
}

func (f hatchFiller) fillPolygon(points []Point, opt *LineOptions) (op operation) {
	fh := NewHachureFiller()
	fh.SetAngle(f.hachureAngle)
	fh.SetGap(f.hachureGap)
	fh.setConnectEnds(f.connectEnds)
	op = fh.fillPolygon(points, opt)
	fh.SetAngle(f.hachureAngle + 90)
	op2 := fh.fillPolygon(points, opt)
	op.commands = append(op.commands, op2.commands...)
	return
}
