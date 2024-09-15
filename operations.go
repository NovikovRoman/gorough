package gorough

import (
	"strings"
)

// operation types
const (
	opPath       = "path"
	opFillPath   = "fillPath"
	opFillSketch = "fillSketch"
)

type operation struct {
	code     string
	commands []command
}

func operationPath(c []command) operation {
	return operation{code: opPath, commands: c}
}

func operationFillPath(c []command) operation {
	return operation{code: opFillPath, commands: c}
}

func operationFillSketch(c []command) operation {
	return operation{code: opFillSketch, commands: c}
}

func lineOperations(p1 Point, p2 Point, pen Pen) []operation {
	return []operation{{code: opPath, commands: doubleLine(p1, p2, pen)}}
}

func operationToPath(op operation) string {
	path := make([]string, len(op.commands))
	for i, c := range op.commands {
		path[i] = c.String()
	}
	return strings.Join(path, "")
}

func rectangleOperation(p Point, width, height float64, pen Pen) operation {
	return polygonOperation([]Point{
		p,
		{
			X: p.X + width,
			Y: p.Y,
		},
		{
			X: p.X + width,
			Y: p.Y + height,
		},
		{
			X: p.X,
			Y: p.Y + height,
		},
	}, pen)
}

func polygonOperation(points []Point, pen Pen) operation {
	return linearPathOperation(points, true, pen)
}

func linearPathOperation(points []Point, close bool, pen Pen) operation {
	var commands []command
	if len(points) == 2 {
		commands = doubleLine(points[0], points[1], pen)

	} else if len(points) > 2 {
		commands = make([]command, 0, len(points)-1*4)

		for i := 0; i < len(points)-1; i++ {
			commands = append(commands, doubleLine(points[i], points[i+1], pen)...)
		}
		if close {
			commands = append(commands, doubleLine(points[len(points)-1], points[0], pen)...)
		}
	}

	return operationPath(commands)
}
