package gorough

import (
	"fmt"
	"regexp"
)

type path struct {
	options    *PathOptions
	operations []operation
}

func (p path) Name() string {
	return shapePoligon
}

func (p path) Operations() []operation {
	return p.operations
}

func (p path) Attributes() Attributes {
	return map[string]string{
		"stroke":       p.options.Styles.Stroke,
		"stroke-width": fmt.Sprintf("%g", p.options.Styles.StrokeWidth),
		"fill":         p.options.Styles.Fill,
		"fill-weight":  fmt.Sprintf("%g", p.options.Styles.FillWeight),
	}
}

func (p path) Styles() *Styles {
	return p.options.Styles
}

var (
	reShapePathLt          = regexp.MustCompile(`(!s)\n`)
	reShapePathHyphen      = regexp.MustCompile(`(!s)-\s`)
	reShapePathDoubleSpace = regexp.MustCompile(`(!s)\s\s`)
)

func NewPath(d string, opt *PathOptions) (*path, error) {
	if opt == nil {
		opt = &PathOptions{}
	}

	if opt.PenOptions == nil {
		opt.PenOptions = PenOptionsDefault()
	}

	if opt.Styles == nil {
		opt.Styles = StylesDefault()
	}

	hasFill := opt.Styles.Fill != None && opt.Styles.Fill != "transparent"
	if (!hasFill || opt.Styles.Stroke != "") && opt.Styles.StrokeWidth == 0 {
		opt.Styles.StrokeWidth = 1
	}

	var (
		operations []operation
		combined   []Point
	)
	operations = []operation{}
	p := &path{
		options:    opt,
		operations: operations,
	}

	if d == "" {
		return p, nil
	}

	d = reShapePathLt.ReplaceAllString(d, " ")
	d = reShapePathHyphen.ReplaceAllString(d, "-")
	d = reShapePathDoubleSpace.ReplaceAllString(d, " ")

	simplified := opt.Simplification < 1
	distance := opt.PenOptions.Roughness / 2
	if simplified {
		distance = 4 - 4*opt.Simplification
	}

	points, err := PointsOnPath(d, 1, distance)
	if err != nil {
		return nil, err
	}

	if hasFill {
		if opt.CombineNestedSvgPaths {
			combined = []Point{}
			for _, pp := range points {
				combined = append(combined, pp...)
			}

			if opt.Styles.Filler == nil {
				operations = append(operations, solidFillPolygon(combined, opt.PenOptions))

			} else {
				operations = append(operations, patternFillPolygon(combined, &LineOptions{
					PenOptions: opt.PenOptions,
					Styles:     opt.Styles,
				}))
			}

		} else {
			for _, pp := range points {
				if opt.Styles.Filler == nil {
					operations = append(operations, solidFillPolygon(pp, opt.PenOptions))

				} else {
					operations = append(operations, patternFillPolygon(pp, &LineOptions{
						PenOptions: opt.PenOptions,
						Styles:     opt.Styles,
					}))
				}
			}
		}
	}

	if opt.Styles.Stroke != None {
		if simplified {
			for _, pp := range points {
				operations = append(operations, linearPathOperation(pp, false, opt.PenOptions))
			}

		} else {
			var op operation
			if op, err = svgPath(d, opt.PenOptions); err != nil {
				return nil, err
			}
			operations = append(operations, op)
		}
	}

	p.operations = operations
	return p, nil
}
