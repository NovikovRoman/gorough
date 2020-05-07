package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func ellipses() (err error) {

	circle := gorough.NewCircle(gorough.Point{X: 40, Y: 80}, 50, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:      "#ff0080",
			StrokeWidth: 1,
		},
	})

	ellipseVert := gorough.NewEllipse(gorough.Point{X: 150, Y: 60}, 60, 110, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke: "#00ff80",
		},
	})

	ellipseHoriz := gorough.NewEllipse(gorough.Point{X: 250, Y: 60}, 150, 80, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:      "#660066",
			StrokeWidth: 1,
		},
	})

	var f *os.File
	if f, err = os.Create("ellipses.svg"); err != nil {
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
	gorough.DrawSVG(canvas, circle, nil)
	gorough.DrawSVG(canvas, ellipseVert, nil)
	gorough.DrawSVG(canvas, ellipseHoriz, nil)
	canvas.End()
	return
}
