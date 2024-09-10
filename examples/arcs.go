package main

import (
	"bufio"
	"bytes"
	"math"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func arcs() (err error) {
	arc := gorough.NewArc(gorough.Point{X: 40, Y: 80}, 100, 100, -math.Pi/2, 0, true,
		gorough.ArcStyle(gorough.Style{
			Stroke: "#ff0000",
		}),
	)

	arc2 := gorough.NewArc(gorough.Point{X: 150, Y: 60}, 60, 110, -math.Pi/6, math.Pi, true,
		gorough.ArcStyle(gorough.Style{
			Stroke: "#00ff80",
		}),
	)

	arc3 := gorough.NewArc(gorough.Point{X: 250, Y: 60}, 150, 80, -math.Pi/2, math.Pi/2, false,
		gorough.ArcStyle(gorough.Style{
			Stroke: "#660066",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 340
	height := 120
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, arc, nil)
	gorough.DrawSVG(canvas, arc2, nil)
	gorough.DrawSVG(canvas, arc3, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("arcs.svg", b.Bytes(), 0644)
	return
}
