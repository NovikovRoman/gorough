package gorough

import (
	"math"
	"math/rand/v2"
)

type Drawable interface {
	Name() string
	Operations() []operation
	Attributes() Attributes
	Style() Style
	Filler() Filler
}

func doubleLine(p1 Point, p2 Point, pen Pen) (commands []command) {
	commands = make([]command, 0, 4)
	commands = append(commands, oneLine(p1, p2, true, false, pen)...)
	commands = append(commands, oneLine(p1, p2, true, true, pen)...)
	return
}

func oneLine(p1 Point, p2 Point, move bool, overlay bool, pen Pen) (commands []command) {
	commands = make([]command, 0, 2)

	lengthSq := math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2)
	length := math.Sqrt(lengthSq)
	roughnessGain := float64(0)
	if length < 200 {
		roughnessGain = 1
	} else if length > 500 {
		roughnessGain = 0.4
	} else {
		roughnessGain = (-0.0016668)*length + 1.233334
	}

	offset := pen.MaxRandomnessOffset
	if offset*offset*100 > lengthSq {
		offset = length / 10
	}

	halfOffset := offset / 2
	divergePoint := 0.2 + rand.Float64()*0.2
	midDispX := pen.Bowing * pen.MaxRandomnessOffset * (p2.Y - p1.Y) / 200
	midDispY := pen.Bowing * pen.MaxRandomnessOffset * (p1.X - p2.X) / 200
	midDispX = offsetOpt(midDispX, pen.Roughness, roughnessGain)
	midDispY = offsetOpt(midDispY, pen.Roughness, roughnessGain)

	if move && overlay {
		commands = append(commands, commandMove([]float64{
			p1.X + offsetOpt(halfOffset, pen.Roughness, roughnessGain),
			p1.Y + offsetOpt(halfOffset, pen.Roughness, roughnessGain),
		}))

	} else if move && !overlay {
		commands = append(commands, commandMove([]float64{
			p1.X + offsetOpt(offset, pen.Roughness, roughnessGain),
			p1.Y + offsetOpt(offset, pen.Roughness, roughnessGain),
		}))
	}

	if overlay {
		offsetOptValue := offsetOpt(halfOffset, pen.Roughness, roughnessGain)
		commands = append(commands, commandCurveTo([]float64{
			midDispX + p1.X + (p2.X-p1.X)*divergePoint + offsetOptValue,
			midDispY + p1.Y + (p2.Y-p1.Y)*divergePoint + offsetOptValue,
			midDispX + p1.X + 2*(p2.X-p1.X)*divergePoint + offsetOptValue,
			midDispY + p1.Y + 2*(p2.Y-p1.Y)*divergePoint + offsetOptValue,
			p2.X + offsetOptValue,
			p2.Y + offsetOptValue,
		}))
		return
	}

	offsetOptValue := offsetOpt(offset, pen.Roughness, roughnessGain)
	commands = append(commands, commandCurveTo([]float64{
		midDispX + p1.X + (p2.X-p1.X)*divergePoint + offsetOptValue,
		midDispY + p1.Y + (p2.Y-p1.Y)*divergePoint + offsetOptValue,
		midDispX + p1.X + 2*(p2.X-p1.X)*divergePoint + offsetOptValue,
		midDispY + p1.Y + 2*(p2.Y-p1.Y)*divergePoint + offsetOptValue,
		p2.X + offsetOptValue,
		p2.Y + offsetOptValue,
	}))
	return
}
