package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func rectangles() (err error) {
	rectangle := gorough.NewRect(gorough.Point{X: 20, Y: 20}, 240, 120)

	square := gorough.NewRect(gorough.Point{X: 10, Y: 10}, 70, 70,
		gorough.RectStyle(gorough.Style{
			Stroke: "#ff0080",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 280
	height := 160
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, rectangle, nil)
	gorough.DrawSVG(canvas, square, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("rectangles.svg", b.Bytes(), 0644)
	return
}
