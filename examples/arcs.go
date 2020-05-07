package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"math"
	"os"
)

func arcs() (err error) {
	arc := gorough.NewArc(gorough.Point{X: 40, Y: 80}, 100, 100, -math.Pi/2, 0, true,
		&gorough.EllipseOptions{
			Styles: &gorough.Styles{
				Stroke: "#ff0080",
			},
		})

	arc2 := gorough.NewArc(gorough.Point{X: 150, Y: 60}, 60, 110, -math.Pi/6, math.Pi, true,
		&gorough.EllipseOptions{
			Styles: &gorough.Styles{
				Stroke: "#00ff80",
			},
		})

	arc3 := gorough.NewArc(gorough.Point{X: 250, Y: 60}, 150, 80, -math.Pi/2, math.Pi/2, false,
		&gorough.EllipseOptions{
			Styles: &gorough.Styles{
				Stroke: "#660066",
			},
		})

	var f *os.File
	if f, err = os.Create("arcs.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 340
	height := 120
	canvas := svg.New(f)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, arc, nil)
	gorough.DrawSVG(canvas, arc2, nil)
	gorough.DrawSVG(canvas, arc3, nil)
	canvas.End()
	return
}
