package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func poligons() (err error) {
	poligon := gorough.NewPoligon([]gorough.Point{
		{X: 10, Y: 10},
		{X: 200, Y: 10},
		{X: 100, Y: 100},
		{X: 300, Y: 100},
		{X: 60, Y: 200},
	})

	poligon2 := gorough.NewPoligon([]gorough.Point{
		{X: 50, Y: 20},
		{X: 30, Y: 200},
		{X: 300, Y: 180},
		{X: 280, Y: 30},
		{X: 170, Y: 180},
	},
		gorough.PoligonStyle(gorough.Style{
			Stroke: "#ff0000",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 320
	height := 220
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, poligon, nil)
	gorough.DrawSVG(canvas, poligon2, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("poligons.svg", b.Bytes(), 0644)
	return
}
