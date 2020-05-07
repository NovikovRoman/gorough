package gorough

import "strings"

type Attributes map[string]string

func (a Attributes) Exclude(attrs ...string) Attributes {
	newAttrs := a
	for _, k := range attrs {
		if a.HasAttr(k) {
			delete(newAttrs, k)
		}
	}
	return newAttrs
}

func (a Attributes) HasAttr(name string) bool {
	for k := range a {
		if k == name {
			return true
		}
	}
	return false
}

func (a Attributes) String() string {
	var attrs []string
	for k, v := range a {
		attrs = append(attrs, k+"='"+v+"'")
	}
	return strings.Join(attrs, " ")
}

type LineOptions struct {
	PenOptions *PenOptions
	Styles     *Styles
}

type EllipseOptions struct {
	PenOptions   *PenOptions
	CurveOptions *CurveOptions
	Styles       *Styles
}

type RectangleOptions struct {
	PenOptions *PenOptions
	Styles     *Styles
}

type PathOptions struct {
	PenOptions            *PenOptions
	Styles                *Styles
	Simplification        float64
	CombineNestedSvgPaths bool
}

type Styles struct {
	Stroke      string
	StrokeWidth float64
	Fill        string
	FillWeight  float64
	Filler      Filler
}

func (s *Styles) canonicalValues() {
	if (s.Fill == "" || s.Stroke != "") && s.StrokeWidth == 0 {
		s.StrokeWidth = 1
	}

	if s.FillWeight == 0 {
		s.FillWeight = s.StrokeWidth / 2
	}
}

func StylesDefault() *Styles {
	return &Styles{
		Stroke:      "#000",
		StrokeWidth: 1,
		Fill:        "",
		FillWeight:  0.5,
	}
}

type PenOptions struct {
	MaxRandomnessOffset float64
	Roughness           float64
	Bowing              float64
}

func PenOptionsDefault() *PenOptions {
	return &PenOptions{
		MaxRandomnessOffset: 2,
		Roughness:           1,
		Bowing:              1,
	}
}

type CurveOptions struct {
	Tightness float64
	Fitting   float64
	StepCount float64
}

func CurveOptionsDefault() *CurveOptions {
	return &CurveOptions{
		Tightness: 0,
		Fitting:   0.95,
		StepCount: 9,
	}
}
