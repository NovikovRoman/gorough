package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func lines() (err error) {
	line := gorough.NewLine(gorough.Point{X: 30, Y: 100}, gorough.Point{X: 270, Y: 20}, &gorough.LineOptions{
		Styles: &gorough.Styles{
			Stroke: "#ff0000",
		},
	})

	line2 := gorough.NewLine(gorough.Point{X: 60, Y: 10}, gorough.Point{X: 230, Y: 110}, &gorough.LineOptions{
		Styles: &gorough.Styles{
			Stroke: "#00ff00",
		},
	})

	line3 := gorough.NewLine(gorough.Point{X: 10, Y: 70}, gorough.Point{X: 250, Y: 90}, &gorough.LineOptions{
		Styles: &gorough.Styles{
			Stroke: "#0000ff",
		},
	})

	var f *os.File
	if f, err = os.Create("lines.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 280
	height := 120
	canvas := svg.New(f)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, line, nil)
	gorough.DrawSVG(canvas, line2, nil)
	gorough.DrawSVG(canvas, line3, nil)
	canvas.End()
	return
}
