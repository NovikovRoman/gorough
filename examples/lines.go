package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func lines() (err error) {
	line := gorough.NewLine(gorough.Point{X: 30, Y: 100}, gorough.Point{X: 270, Y: 20},
		gorough.LineStyle(gorough.Style{
			Stroke: "#ff0000",
		}),
	)

	line2 := gorough.NewLine(gorough.Point{X: 60, Y: 10}, gorough.Point{X: 230, Y: 110},
		gorough.LineStyle(gorough.Style{
			Stroke: "#00ff00",
		}),
	)

	line3 := gorough.NewLine(gorough.Point{X: 10, Y: 70}, gorough.Point{X: 250, Y: 90},
		gorough.LineStyle(gorough.Style{
			Stroke: "#0000ff",
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 280
	height := 120
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, line, nil)
	gorough.DrawSVG(canvas, line2, nil)
	gorough.DrawSVG(canvas, line3, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("lines.svg", b.Bytes(), 0644)
	return
}
