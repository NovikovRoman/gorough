package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func linearPaths() (err error) {
	linearPath := gorough.NewLinearPath([]gorough.Point{
		{X: 10, Y: 10},
		{X: 200, Y: 10},
		{X: 100, Y: 100},
		{X: 300, Y: 100},
		{X: 60, Y: 200},
	}, nil)

	linearPath2 := gorough.NewLinearPath([]gorough.Point{
		{X: 20, Y: 50},
		{X: 50, Y: 200},
		{X: 80, Y: 50},
		{X: 110, Y: 200},
		{X: 140, Y: 50},
		{X: 170, Y: 200},
		{X: 200, Y: 50},
		{X: 230, Y: 200},
		{X: 260, Y: 50},
		{X: 290, Y: 200},
	}, &gorough.LineOptions{
		Styles: &gorough.Styles{
			Stroke: "#00ff00",
		},
	})

	var f *os.File
	if f, err = os.Create("linear_paths.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 320
	height := 220
	canvas := svg.New(f)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, linearPath, nil)
	gorough.DrawSVG(canvas, linearPath2, nil)
	canvas.End()
	return
}
