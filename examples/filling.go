package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func filling() (err error) {
	dashedFiller := gorough.NewDashedFiller()
	dotFiller := gorough.NewDotFiller()
	hachureFiller := gorough.NewHachureFiller()
	hatchFiller := gorough.NewHatchFiller()
	zigZagFiller := gorough.NewZigZagFiller()
	zigZagHatchFiller := gorough.NewZigZagHatchFiller()

	// row 1 --------------------------------------------------------------------------
	dashed := gorough.NewRectangle(gorough.Point{X: 20, Y: 20}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#4e9835",
			Fill:       "#60c247",
			FillWeight: 1,
			Filler:     dashedFiller,
		},
	})

	dot := gorough.NewRectangle(gorough.Point{X: 140, Y: 20}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#155da5",
			Fill:       "#2171d4",
			FillWeight: 1,
			Filler:     dotFiller,
		},
	})

	hachure := gorough.NewRectangle(gorough.Point{X: 260, Y: 20}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#c84016",
			Fill:       "#eb551a",
			FillWeight: 1,
			Filler:     hachureFiller,
		},
	})

	hatch := gorough.NewRectangle(gorough.Point{X: 380, Y: 20}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#492bfd",
			Fill:       "#6239fe",
			FillWeight: 1,
			Filler:     hatchFiller,
		},
	})

	// row 2 --------------------------------------------------------------------------
	zigZag := gorough.NewRectangle(gorough.Point{X: 20, Y: 160}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#c87800",
			Fill:       "#fc9700",
			FillWeight: 1,
			Filler:     zigZagFiller,
		},
	})

	zigZagHatch := gorough.NewRectangle(gorough.Point{X: 140, Y: 160}, 100, 100, &gorough.RectangleOptions{
		Styles: &gorough.Styles{
			Stroke:     "#9c001b",
			Fill:       "#c70023",
			FillWeight: 1,
			Filler:     zigZagHatchFiller,
		},
	})

	dashedFiller.SetAngle(-41)
	dashedCircle := gorough.NewCircle(gorough.Point{X: 310, Y: 210}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#4e9835",
			Fill:       "#60c247",
			FillWeight: 1,
			Filler:     dashedFiller,
		},
	})

	dotFiller.SetGap(8)
	dotCircle := gorough.NewCircle(gorough.Point{X: 430, Y: 210}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#155da5",
			Fill:       "#2171d4",
			FillWeight: 1,
			Filler:     dotFiller,
		},
	})

	// row 3 --------------------------------------------------------------------------
	hachureFiller.SetAngle(0)
	hachureCircle := gorough.NewCircle(gorough.Point{X: 70, Y: 350}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#c84016",
			Fill:       "#eb551a",
			FillWeight: 1,
			Filler:     hachureFiller,
		},
	})

	hatchFiller.SetAngle(90)
	hatchFiller.SetGap(8)
	hatchCircle := gorough.NewCircle(gorough.Point{X: 190, Y: 350}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#492bfd",
			Fill:       "#6239fe",
			FillWeight: 1,
			Filler:     hatchFiller,
		},
	})

	zigZagFiller.SetAngle(49)
	zigZagFiller.SetGap(4)
	zigZagCircle := gorough.NewCircle(gorough.Point{X: 310, Y: 350}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#c87800",
			Fill:       "#fc9700",
			FillWeight: 1,
			Filler:     zigZagFiller,
		},
	})

	zigZagHatchFiller.SetGap(5)
	zigZagHatchCircle := gorough.NewCircle(gorough.Point{X: 430, Y: 350}, 100, &gorough.EllipseOptions{
		Styles: &gorough.Styles{
			Stroke:     "#9c001b",
			Fill:       "#c70023",
			FillWeight: 1,
			Filler:     zigZagHatchFiller,
		},
	})

	var f *os.File
	if f, err = os.Create("filling.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 500
	height := 410
	style := "font-family:'sans-serif';font-size:12px"
	canvas := svg.New(f)
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
	return
}
