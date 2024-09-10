package gorough

import "strings"

type Attributes map[string]string

func (a Attributes) Exclude(attrs ...string) Attributes {
	newAttrs := a
	for _, k := range attrs {
		if _, ok := a[k]; ok {
			delete(newAttrs, k)
		}
	}
	return newAttrs
}

func (a Attributes) String() string {
	var attrs []string
	for k, v := range a {
		attrs = append(attrs, k+"='"+v+"'")
	}
	return strings.Join(attrs, " ")
}

type Style struct {
	Stroke      string
	StrokeWidth float64
	Fill        string
	FillWeight  float64
}

func (s Style) canonicalValues() Style {
	if (s.Fill == "" || s.Stroke != "") && s.StrokeWidth == 0 {
		s.StrokeWidth = 1
	}

	if s.FillWeight == 0 {
		s.FillWeight = s.StrokeWidth / 2
	}
	return s
}

func StyleDefault() Style {
	return Style{
		Stroke:      "#000",
		StrokeWidth: 1,
		Fill:        "",
		FillWeight:  0.5,
	}
}

type Pen struct {
	MaxRandomnessOffset float64
	Roughness           float64
	Bowing              float64
}

func PenDefault() Pen {
	return Pen{
		MaxRandomnessOffset: 2,
		Roughness:           1,
		Bowing:              1,
	}
}

type CurveOption struct {
	Tightness float64
	Fitting   float64
	StepCount float64
}

func CurveDefault() CurveOption {
	return CurveOption{
		Tightness: 0,
		Fitting:   0.95,
		StepCount: 9,
	}
}
