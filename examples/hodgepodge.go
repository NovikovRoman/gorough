package main

import (
	"bufio"
	"bytes"
	"math"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func hodgepodge() (err error) {
	svgPath, err := gorough.NewPath("M 240 100 C 290 100, 240 225, 290 200 S 290 75, 340 50 S 515 100, 390 150 S 215 200, 90 150 S 90 25, 140 50 S 140 175, 190 200 S 190 100, 240 100",
		gorough.PathStyle(gorough.Style{
			Stroke:      "#660000",
			StrokeWidth: 1,
			Fill:        "#CC6633",
		}),
	)
	if err != nil {
		return err
	}

	line := gorough.NewLine(gorough.Point{X: 200, Y: 350}, gorough.Point{X: 400, Y: 400},
		gorough.LineStyle(gorough.Style{
			Stroke:      "#CC0033",
			StrokeWidth: 1,
		}),
	)

	rectangle := gorough.NewRect(gorough.Point{X: 240, Y: 240}, 50, 50,
		gorough.RectStyle(gorough.Style{
			Stroke:      "#663333",
			StrokeWidth: 1,
			Fill:        "#996666",
			FillWeight:  1,
		}),
		gorough.RectFiller(gorough.NewZigZagFiller()),
	)

	ellipse := gorough.NewEllipse(gorough.Point{X: 400, Y: 300}, 100, 50,
		gorough.EllipseStyle(gorough.Style{
			Stroke:      "#660066",
			StrokeWidth: 1,
			Fill:        "#9900CC",
			FillWeight:  1,
		}),
		gorough.EllipseFiller(gorough.NewZigZagHatchFiller()),
	)

	circle := gorough.NewCircle(gorough.Point{X: 450, Y: 200}, 50,
		gorough.CircleStyle(gorough.Style{
			Stroke:      "#000066",
			StrokeWidth: 1,
			Fill:        "#3333CC",
			FillWeight:  1,
		}),
		gorough.CircleFiller(gorough.NewDotFiller()),
	)

	linearPath := gorough.NewLinearPath([]gorough.Point{
		{X: 10, Y: 230},
		{X: 30, Y: 240},
		{X: 50, Y: 230},
		{X: 70, Y: 240},
		{X: 90, Y: 230},
		{X: 110, Y: 240},
	},
		gorough.LinearPathStyle(gorough.Style{
			Stroke:      "#006633",
			StrokeWidth: 1,
		}),
	)

	poligon := gorough.NewPoligon([]gorough.Point{
		{X: 180, Y: 20},
		{X: 240, Y: 25},
		{X: 250, Y: 60},
		{X: 225, Y: 40},
		{X: 200, Y: 70},
	},
		gorough.PoligonStyle(gorough.Style{
			Stroke:      "#003300",
			StrokeWidth: 1,
			Fill:        "#336633",
			FillWeight:  1,
		}),
		gorough.PoligonFiller(gorough.NewHatchFiller()),
	)

	curve := gorough.NewCurve([]gorough.Point{
		{X: 10, Y: 160},
		{X: 30, Y: 190},
		{X: 50, Y: 160},
		{X: 70, Y: 190},
		{X: 90, Y: 160},
		{X: 110, Y: 190},
	},
		gorough.CurveStyle(gorough.Style{
			Stroke:      "#000000",
			StrokeWidth: 1,
		}),
	)

	arc := gorough.NewArc(gorough.Point{X: 50, Y: 300}, 280, 60, -math.Pi/3, math.Pi/7, true,
		gorough.ArcStyle(gorough.Style{
			Stroke:      "#ff00ff",
			StrokeWidth: 1,
			Fill:        "#000",
			FillWeight:  1,
		}),
		gorough.ArcFiller(gorough.NewDashedFiller()),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 500
	height := 400
	canvas := svg.New(picSvg)

	canvas.Start(width, height)
	gorough.DrawSVG(canvas, svgPath, nil)
	gorough.DrawSVG(canvas, line, nil)
	gorough.DrawSVG(canvas, rectangle, nil)
	gorough.DrawSVG(canvas, ellipse, nil)
	gorough.DrawSVG(canvas, circle, nil)
	gorough.DrawSVG(canvas, linearPath, nil)
	gorough.DrawSVG(canvas, poligon, nil)
	gorough.DrawSVG(canvas, curve, nil)
	gorough.DrawSVG(canvas, arc, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	return os.WriteFile("hodgepodge.svg", b.Bytes(), 0644)
}
