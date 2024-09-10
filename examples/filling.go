package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/NovikovRoman/gorough"
	svg "github.com/NovikovRoman/svgo"
)

func filling() (err error) {
	dashedFiller := gorough.NewDashedFiller()
	dotFiller := gorough.NewDotFiller()
	hachureFiller := gorough.NewHachureFiller()
	hatchFiller := gorough.NewHatchFiller()
	zigZagFiller := gorough.NewZigZagFiller()
	zigZagHatchFiller := gorough.NewZigZagHatchFiller()

	// row 1 --------------------------------------------------------------------------
	dashed := gorough.NewRect(gorough.Point{X: 20, Y: 20}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#4e9835",
			Fill:       "#60c247",
			FillWeight: 1,
		}),
		gorough.RectFiller(dashedFiller),
	)

	dot := gorough.NewRect(gorough.Point{X: 140, Y: 20}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#155da5",
			Fill:       "#2171d4",
			FillWeight: 1,
		}),
		gorough.RectFiller(dotFiller),
	)

	hachure := gorough.NewRect(gorough.Point{X: 260, Y: 20}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#c84016",
			Fill:       "#eb551a",
			FillWeight: 1,
		}),
		gorough.RectFiller(hachureFiller),
	)

	hatch := gorough.NewRect(gorough.Point{X: 380, Y: 20}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#492bfd",
			Fill:       "#6239fe",
			FillWeight: 1,
		}),
		gorough.RectFiller(hatchFiller),
	)

	// row 2 --------------------------------------------------------------------------
	zigZag := gorough.NewRect(gorough.Point{X: 20, Y: 160}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#c87800",
			Fill:       "#fc9700",
			FillWeight: 1,
		}),
		gorough.RectFiller(zigZagFiller),
	)

	zigZagHatch := gorough.NewRect(gorough.Point{X: 140, Y: 160}, 100, 100,
		gorough.RectStyle(gorough.Style{
			Stroke:     "#9c001b",
			Fill:       "#c70023",
			FillWeight: 1,
		}),
		gorough.RectFiller(zigZagHatchFiller),
	)

	dashedFiller.SetAngle(-41)
	dashedCircle := gorough.NewCircle(gorough.Point{X: 310, Y: 210}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#4e9835",
			Fill:       "#60c247",
			FillWeight: 1,
		}),
		gorough.CircleFiller(dashedFiller),
	)

	dotFiller.SetGap(8)
	dotCircle := gorough.NewCircle(gorough.Point{X: 430, Y: 210}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#155da5",
			Fill:       "#2171d4",
			FillWeight: 1,
		}),
		gorough.CircleFiller(dotFiller),
	)

	// row 3 --------------------------------------------------------------------------
	hachureFiller.SetAngle(0)
	hachureCircle := gorough.NewCircle(gorough.Point{X: 70, Y: 350}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#c84016",
			Fill:       "#eb551a",
			FillWeight: 1,
		}),
		gorough.CircleFiller(hachureFiller),
	)

	hatchFiller.SetAngle(90)
	hatchFiller.SetGap(8)
	hatchCircle := gorough.NewCircle(gorough.Point{X: 190, Y: 350}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#492bfd",
			Fill:       "#6239fe",
			FillWeight: 1,
		}),
		gorough.CircleFiller(hatchFiller),
	)

	zigZagFiller.SetAngle(49)
	zigZagFiller.SetGap(4)
	zigZagCircle := gorough.NewCircle(gorough.Point{X: 310, Y: 350}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#c87800",
			Fill:       "#fc9700",
			FillWeight: 1,
		}),
		gorough.CircleFiller(zigZagFiller),
	)

	zigZagHatchFiller.SetGap(5)
	zigZagHatchCircle := gorough.NewCircle(gorough.Point{X: 430, Y: 350}, 100,
		gorough.CircleStyle(gorough.Style{
			Stroke:     "#9c001b",
			Fill:       "#c70023",
			FillWeight: 1,
		}),
		gorough.CircleFiller(zigZagHatchFiller),
	)

	var b bytes.Buffer
	picSvg := bufio.NewWriter(&b)

	width := 500
	height := 410
	style := "font-family:'sans-serif';font-size:12px"
	canvas := svg.New(picSvg)
	canvas.Start(width, height)
	canvas.Text(50, 15, "dashed", style)
	gorough.DrawSVG(canvas, dashed, nil)
	canvas.Text(180, 15, "dot", style)
	gorough.DrawSVG(canvas, dot, nil)
	canvas.Text(285, 15, "hachure", style)
	gorough.DrawSVG(canvas, hachure, nil)
	canvas.Text(410, 15, "hatch", style)
	gorough.DrawSVG(canvas, hatch, nil)

	canvas.Text(50, 150, "zigzag", style)
	gorough.DrawSVG(canvas, zigZag, nil)
	canvas.Text(150, 150, "zigzag hatch", style)
	gorough.DrawSVG(canvas, zigZagHatch, nil)
	canvas.Text(270, 150, "tuned dashed", style)
	gorough.DrawSVG(canvas, dashedCircle, nil)
	canvas.Text(400, 150, "tuned dot", style)
	gorough.DrawSVG(canvas, dotCircle, nil)

	canvas.Text(25, 290, "tuned hachure", style)
	gorough.DrawSVG(canvas, hachureCircle, nil)
	canvas.Text(155, 290, "tuned hatch", style)
	gorough.DrawSVG(canvas, hatchCircle, nil)
	canvas.Text(270, 290, "tuned zigzag", style)
	gorough.DrawSVG(canvas, zigZagCircle, nil)
	canvas.Text(375, 290, "tuned zigzag hatch", style)
	gorough.DrawSVG(canvas, zigZagHatchCircle, nil)
	canvas.End()

	if err = picSvg.Flush(); err != nil {
		return
	}
	err = os.WriteFile("filling.svg", b.Bytes(), 0644)
	return
}
