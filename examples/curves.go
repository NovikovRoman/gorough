package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func curves() (err error) {
	curve := gorough.NewCurve([]gorough.Point{
		{X: 10, Y: 10},
		{X: 200, Y: 10},
		{X: 100, Y: 100},
		{X: 300, Y: 100},
		{X: 60, Y: 200},
	})

	curve2 := gorough.NewCurve([]gorough.Point{
		{X: 50, Y: 20},
		{X: 30, Y: 200},
		{X: 300, Y: 180},
		{X: 280, Y: 30},
		{X: 170, Y: 180},
	},
		gorough.CurveStyle(gorough.Style{
			Stroke: "#ff0000",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 320
	height := 220
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, curve, nil)
	gorough.DrawSVG(canvas, curve2, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("curves.svg", b.Bytes(), 0644)
	return
}
