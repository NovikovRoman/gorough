package data_parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAbsolutize(t *testing.T) {
	segments, err := ParsePath(" M10 10 h 80 v 80 h -80 Z  ")
	require.Nil(t, err)

	absoluteSegments := Absolutize(segments)
	require.Equal(t, absoluteSegments, []Segment{
		{Key: "M", Data: []float64{10, 10}},
		{Key: "H", Data: []float64{90}},
		{Key: "V", Data: []float64{90}},
		{Key: "H", Data: []float64{10}},
		{Key: "Z", Data: []float64{}},
	})
}
