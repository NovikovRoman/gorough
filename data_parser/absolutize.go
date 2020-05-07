package data_parser

func Absolutize(segments []Segment) (out []Segment) {
	cx := float64(0)
	cy := float64(0)
	subx := float64(0)
	suby := float64(0)

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

		case "m":
			cx += s.Data[0]
			cy += s.Data[1]
			out = append(out, Segment{
				Key:  "M",
				Data: []float64{cx, cy},
			})
			subx = cx
			suby = cy

		case "L":
			out = append(out, Segment{
				Key:  "L",
				Data: s.Data,
			})
			cx = s.Data[0]
			cy = s.Data[1]

		case "l":
			cx += s.Data[0]
			cy += s.Data[1]
			out = append(out, Segment{
				Key:  "L",
				Data: []float64{cx, cy},
			})

		case "C":
			out = append(out, Segment{
				Key:  "C",
				Data: s.Data,
			})
			cx = s.Data[4]
			cy = s.Data[5]

		case "c":
			newData := calcNewData(s.Data, cx, cy)
			out = append(out, Segment{
				Key:  "C",
				Data: newData,
			})
			cx = newData[4]
			cy = newData[5]

		case "Q":
			out = append(out, Segment{
				Key:  "Q",
				Data: s.Data,
			})
			cx = s.Data[2]
			cy = s.Data[3]

		case "q":
			newData := calcNewData(s.Data, cx, cy)
			out = append(out, Segment{
				Key:  "Q",
				Data: newData,
			})
			cx = newData[2]
			cy = newData[3]

		case "A":
			out = append(out, Segment{
				Key:  "A",
				Data: s.Data,
			})
			cx = s.Data[5]
			cy = s.Data[6]

		case "a":
			cx += s.Data[5]
			cy += s.Data[6]
			out = append(out, Segment{
				Key:  "A",
				Data: []float64{s.Data[0], s.Data[1], s.Data[2], s.Data[3], s.Data[4], cx, cy},
			})

		case "H":
			out = append(out, Segment{
				Key:  "H",
				Data: s.Data,
			})
			cx = s.Data[0]

		case "h":
			cx += s.Data[0]
			out = append(out, Segment{
				Key:  "H",
				Data: []float64{cx},
			})

		case "V":
			out = append(out, Segment{
				Key:  "V",
				Data: s.Data,
			})
			cy = s.Data[0]

		case "v":
			cy += s.Data[0]
			out = append(out, Segment{
				Key:  "V",
				Data: []float64{cy},
			})

		case "S":
			out = append(out, Segment{
				Key:  "S",
				Data: s.Data,
			})
			cx = s.Data[2]
			cy = s.Data[3]

		case "s":
			newData := calcNewData(s.Data, cx, cy)
			out = append(out, Segment{
				Key:  "S",
				Data: newData,
			})
			cx = newData[2]
			cy = newData[3]

		case "T":
			out = append(out, Segment{
				Key:  "T",
				Data: s.Data,
			})
			cx = s.Data[0]
			cy = s.Data[1]

		case "t":
			cx += s.Data[0]
			cy += s.Data[1]
			out = append(out, Segment{
				Key:  "T",
				Data: []float64{cx, cy},
			})

		case "Z", "z":
			out = append(out, Segment{
				Key:  "Z",
				Data: []float64{},
			})
			cx = subx
			cy = suby
		}
	}
	return
}

func calcNewData(data []float64, cx, cy float64) []float64 {
	var res []float64
	res = []float64{}
	for i, d := range data {
		if i%2 > 0 {
			res = append(res, d+cy)
		} else {
			res = append(res, d+cx)
		}
	}
	return res
}
