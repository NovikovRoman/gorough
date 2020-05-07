package data_parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNormalize(t *testing.T) {
	segments, err := ParsePath(" M 10 80 Q 52.5 10, 95 80 T 180 80")
	require.Nil(t, err)
	absSegments := Absolutize(segments)
	require.Equal(t, Normalize(absSegments), []Segment{
		{Key: "M", Data: []float64{10, 80}},
		{Key: "C", Data: []float64{38.33333333333333, 33.333333333333336, 66.66666666666667, 33.333333333333336, 95, 80}},
		{Key: "C", Data: []float64{123.33333333333333, 126.66666666666666, 151.66666666666666, 126.66666666666666, 180, 80}},
	})
}

