package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func linearPaths() (err error) {
	linearPath := gorough.NewLinearPath([]gorough.Point{
		{X: 10, Y: 10},
		{X: 200, Y: 10},
		{X: 100, Y: 100},
		{X: 300, Y: 100},
		{X: 60, Y: 200},
	})

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
	},
		gorough.LinearPathStyle(gorough.Style{
			Stroke: "#00ff00",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 320
	height := 220
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, linearPath, nil)
	gorough.DrawSVG(canvas, linearPath2, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("linear_paths.svg", b.Bytes(), 0644)
	return
}
