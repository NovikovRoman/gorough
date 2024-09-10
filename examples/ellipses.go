package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func ellipses() (err error) {

	circle := gorough.NewCircle(gorough.Point{X: 40, Y: 80}, 50,
		gorough.CircleStyle(gorough.Style{
			Stroke:      "#ff0080",
			StrokeWidth: 1,
		}),
	)

	ellipseVert := gorough.NewEllipse(gorough.Point{X: 150, Y: 60}, 60, 110,
		gorough.EllipseStyle(gorough.Style{
			Stroke:      "#00ff80",
			StrokeWidth: 1,
		}),
	)

	ellipseHoriz := gorough.NewEllipse(gorough.Point{X: 250, Y: 60}, 150, 80,
		gorough.EllipseStyle(gorough.Style{
			Stroke:      "#660066",
			StrokeWidth: 1,
		}),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 340
	height := 120
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, circle, nil)
	gorough.DrawSVG(canvas, ellipseVert, nil)
	gorough.DrawSVG(canvas, ellipseHoriz, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("ellipses.svg", b.Bytes(), 0644)
	return
}
