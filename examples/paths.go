package main

import (
	"github.com/NovikovRoman/gorough"
	"github.com/ajstarks/svgo"
	"os"
)

func paths() (err error) {
	svgPath, err := gorough.NewPath("M148.14 67.19L136.6 40.44L88.53 26.93L48.7 40.44L26.93 72.92L26.93 117.19L39.59 158.59L58.35 184.48L92.11 196.53L130.7 188.15L132.31 200.2L152.67 200.2L150.53 128.44L104.98 126.87L104.98 141.01C113.41 142.41 118.68 143.28 120.78 143.63C128.69 144.94 133.14 153.47 129.68 160.7C127.95 164.31 128.35 163.48 126.29 167.79C123.09 174.47 116.34 178.72 108.93 178.72C102.57 178.72 92.59 178.72 85.47 178.72C80.25 178.72 75.38 176.1 72.5 171.75C67.74 164.55 62.24 156.24 58.53 150.63C52.11 140.92 48.56 129.6 48.3 117.97C48.17 112.1 48.12 110.09 47.97 103.67C47.75 93.76 49.69 83.92 53.65 74.84C55.55 70.48 54.79 72.21 56.59 68.11C60.59 58.94 69.64 53.01 79.65 53.01C84.71 53.01 84.9 53.01 90.33 53.01C99.7 53.01 108.68 56.77 115.25 63.45C117.52 65.76 123.21 71.53 132.31 80.77L148.14 67.19Z", &gorough.PathOptions{
		Styles: &gorough.Styles{
			Stroke: "#ff0000",
		},
	})
	if err != nil {
		return err
	}

	svgPath2, err := gorough.NewPath("M234.19 86.4L234.19 106.27L234.19 133.75L230.58 150.97L225.17 167.19L210.75 167.19L193.74 161.56L183.7 141.37L183.7 110.9L183.7 82.09L193.74 67.19L219.25 70.5L234.19 86.4ZM164.56 105.38L163.54 154.11L189.35 187.87L230.05 195.08L247.36 155.88L248.57 98.67L249.4 59.11L192.44 40.96L169.98 52.49L164.56 105.38Z", &gorough.PathOptions{
		Styles: &gorough.Styles{
			Stroke: "#0000ff",
		},
	})
	if err != nil {
		return err
	}

	var f *os.File
	if f, err = os.Create("paths.svg"); err != nil {
		return
	}
	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	width := 280
	height := 240
	canvas := svg.New(f)
	canvas.Start(width, height)
	gorough.DrawSVG(canvas, svgPath, nil)
	gorough.DrawSVG(canvas, svgPath2, nil)
	canvas.End()
	return
}
