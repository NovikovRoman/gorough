package gorough

func NewZigZagHatchFiller() Filler {
	return &hatchFiller{
		connectEnds:  true,
		hachureAngle: -41,
		hachureGap:   6,
	}
}
