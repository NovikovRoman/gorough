package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func rectangles() (err error) {
	rectangle := gorough.NewRectangle(gorough.Point{X: 20, Y: 20}, 240, 120, nil)

	square := gorough.NewRectangle(gorough.Point{X: 10, Y: 10}, 70, 70, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke: "#ff0080",
		},
	})

	var f *os.File
	if f, err = os.Create("rectangles.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 280
	height := 160
	canvas := svg.New(f)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, rectangle, nil)
	gorough.DrawSVG(canvas, square, nil)
	canvas.End()
	return
}
