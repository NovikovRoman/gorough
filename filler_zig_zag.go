package gorough

func NewZigZagFiller() Filler {
	return &hachureFiller{
		connectEnds:  true,
		hachureAngle: -41,
		hachureGap:   6,
	}
}
