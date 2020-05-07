package data_parser

import "math"

func Normalize(segments []Segment) (out []Segment) {
	lastType := ""
	cx := float64(0)
	cy := float64(0)
	subx := float64(0)
	suby := float64(0)
	lcx := float64(0)
	lcy := float64(0)

	for _, s := range segments {
		switch s.Key {
		case "M":
			out = append(out, Segment{
				Key:  "M",
				Data: s.Data,
			})
			cx = s.Data[0]
			cy = s.Data[1]
			subx = s.Data[0]
			suby = s.Data[1]

		case "C":
			out = append(out, Segment{
				Key:  "C",
				Data: s.Data,
			})
			cx = s.Data[4]
			cy = s.Data[5]
			lcx = s.Data[2]
			lcy = s.Data[3]

		case "L":
			out = append(out, Segment{
				Key:  "L",
				Data: s.Data,
			})
			cx = s.Data[0]
			cy = s.Data[1]

		case "H":
			cx = s.Data[0]
			out = append(out, Segment{
				Key:  "L",
				Data: []float64{cx, cy},
			})

		case "V":
			cy = s.Data[0]
			out = append(out, Segment{
				Key:  "L",
				Data: []float64{cx, cy},
			})

		case "S":
			cx1 := cx
			cy1 := cy
			if lastType == "C" || lastType == "S" {
				cx1 = cx + (cx - lcx)
				cy1 = cy + (cy - lcy)
			}
			newData := []float64{cx1, cy1}
			newData = append(newData, s.Data...)
			out = append(out, Segment{
				Key:  "C",
				Data: newData,
			})
			lcx = s.Data[0]
			lcy = s.Data[1]
			cx = s.Data[2]
			cy = s.Data[3]

		case "T":
			x := s.Data[0]
			y := s.Data[1]
			x1 := cx
			y1 := cy
			if lastType == "Q" || lastType == "T" {
				x1 = cx + (cx - lcx)
				y1 = cy + (cy - lcy)
			}
			cx1 := cx + 2*(x1-cx)/3
			cy1 := cy + 2*(y1-cy)/3
			cx2 := x + 2*(x1-x)/3
			cy2 := y + 2*(y1-y)/3

			out = append(out, Segment{
				Key:  "C",
				Data: []float64{cx1, cy1, cx2, cy2, x, y},
			})
			lcx = x1
			lcy = y1
			cx = x
			cy = y

		case "Q":
			x1 := s.Data[0]
			y1 := s.Data[1]
			x := s.Data[2]
			y := s.Data[3]
			cx1 := cx + 2*(x1-cx)/3
			cy1 := cy + 2*(y1-cy)/3
			cx2 := x + 2*(x1-x)/3
			cy2 := y + 2*(y1-y)/3

			out = append(out, Segment{
				Key:  "C",
				Data: []float64{cx1, cy1, cx2, cy2, x, y},
			})
			lcx = x1
			lcy = y1
			cx = x
			cy = y

		case "A":
			r1 := math.Abs(s.Data[0])
			r2 := math.Abs(s.Data[1])
			angle := s.Data[2]
			largeArcFlag := s.Data[3]
			sweepFlag := s.Data[4]
			x := s.Data[5]
			y := s.Data[6]

			if r1 == 0 || r2 == 0 {
				out = append(out, Segment{
					Key:  "C",
					Data: []float64{cx, cy, x, y, x, y},
				})
				cx = x
				cy = y

			} else if cx != x || cy != y {
				for _, curve := range arcToCubicCurves(cx, cy, x, y, r1, r2, angle, largeArcFlag, sweepFlag) {
					out = append(out, Segment{Key: "C", Data: curve})
				}
				cx = x
				cy = y
			}

		case "Z":
			out = append(out, Segment{
				Key:  "Z",
				Data: []float64{},
			})
			cx = subx
			cy = suby
		}

		lastType = s.Key
	}
	return
}

func degToRad(degrees float64) float64 {
	return math.Pi * degrees / 180
}

func rotate(x float64, y float64, angleRad float64) (resX float64, resY float64) {
	resX = x*math.Cos(angleRad) - y*math.Sin(angleRad)
	resY = x*math.Sin(angleRad) + y*math.Cos(angleRad)
	return
}

func arcToCubicCurves(x1, y1, x2, y2, r1, r2, angle, largeArcFlag, sweepFlag float64, recursive ...float64) [][]float64 {
	angleRad := degToRad(angle)
	var params [][]float64
	params = [][]float64{}
	f1 := float64(0)
	f2 := float64(0)
	cx := float64(0)
	cy := float64(0)
	if len(recursive) > 0 {
		f1 = recursive[0]
		f2 = recursive[1]
		cx = recursive[2]
		cy = recursive[3]

	} else {
		x1, y1 = rotate(x1, y1, -angle)
		x2, y2 = rotate(x2, y2, -angle)
		x := (x1 - x2) / 2
		y := (y1 - y2) / 2
		h := (x*x)/(r1*r1) + (y*y)/(r2*r2)
		if h > 1 {
			h = math.Sqrt(h)
			r1 = h * r1
			r2 = h * r2
		}

		sign := float64(1)
		if largeArcFlag == sweepFlag {
			sign = -1
		}
		r1Pow := r1 * r1
		r2Pow := r2 * r2
		left := r1Pow*r2Pow - r1Pow*y*y - r2Pow*x*x
		right := r1Pow*y*y + r2Pow*x*x
		k := sign * math.Sqrt(math.Abs(left/right))
		cx = k*r1*y/r2 + (x1+x2)/2
		cy = k*-r2*x/r1 + (y1+y2)/2
		f1 = math.Asin((y1 - cy) / r2)
		f2 = math.Asin((y2 - cy) / r2)

		if x1 < cx {
			f1 = math.Pi - f1
		}

		if x2 < cx {
			f2 = math.Pi - f2
		}

		if f1 < 0 {
			f1 = math.Pi*2 - f1
		}
		if f2 < 0 {
			f2 = math.Pi*2 - f2
		}

		if sweepFlag != 0 && f1 > f2 {
			f1 = f1 - math.Pi*2
		}
		if sweepFlag == 0 && f2 > f1 {
			f2 = f2 - math.Pi*2
		}
	}

	df := f2 - f1

	if math.Abs(df) > (math.Pi * 120 / 180) {
		f2old := f2
		x2old := x2
		y2old := y2

		if sweepFlag != 0 && f2 > f1 {
			f2 = f1 + (math.Pi * 120 / 180) //*(1)
		} else {
			f2 = f1 + (-math.Pi * 120 / 180) //*(-1)
		}
		x2 = cx + r1*math.Cos(f2)
		y2 = cy + r2*math.Sin(f2)
		params = arcToCubicCurves(x2, y2, x2old, y2old, r1, r2, angle, 0, sweepFlag, f2, f2old, cx, cy)
	}

	df = f2 - f1

	c1 := math.Cos(f1)
	s1 := math.Sin(f1)
	c2 := math.Cos(f2)
	s2 := math.Sin(f2)
	t := math.Tan(df / 4)
	hx := 4 / 3 * r1 * t
	hy := 4 / 3 * r2 * t

	m1 := []float64{x1, y1}
	m2 := []float64{x1 + hx*s1, y1 - hy*c1}
	m3 := []float64{x2 + hx*s2, y2 - hy*c2}
	m4 := []float64{x2, y2}

	m2[0] = 2*m1[0] - m2[0]
	m2[1] = 2*m1[1] - m2[1]

	res := [][]float64{m2, m3, m4}
	if len(recursive) > 0 {
		return append(res, params...)
	}

	params = append(res, params...)
	var curves [][]float64
	curves = [][]float64{}
	for i := 0; i < len(params); i += 3 {
		r11, r12 := rotate(params[i][0], params[i][1], angleRad)
		r21, r22 := rotate(params[i+1][0], params[i+1][1], angleRad)
		r31, r32 := rotate(params[i+2][0], params[i+2][1], angleRad)
		curves = append(curves, []float64{r11, r12, r21, r22, r31, r32})
	}
	return curves
}
