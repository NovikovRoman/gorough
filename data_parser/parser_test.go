package data_parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParsePath(t *testing.T) {
	segments, err := ParsePath("M10 10 h 80 v 80 h -80 C Z")
	require.NotNil(t, err)

	segments, err = ParsePath(" M10 10 h 80 v 80 h -80 Z  ")
	require.Nil(t, err)
	require.Equal(t, segments, []Segment{
		{Key: "M", Data: []float64{10, 10}},
		{Key: "h", Data: []float64{80}},
		{Key: "v", Data: []float64{80}},
		{Key: "h", Data: []float64{-80}},
		{Key: "Z", Data: []float64{}},
	})

	segments, err = ParsePath(" M 10 80 Q 52.5 10, 95 80 T 180 80")
	require.Nil(t, err)
	require.Equal(t, segments, []Segment{
		{Key: "M", Data: []float64{10, 80}},
		{Key: "Q", Data: []float64{52.5, 10, 95, 80}},
		{Key: "T", Data: []float64{180, 80}},
	})
}

func TestSerialize(t *testing.T) {
	segments, err := ParsePath(",,, M 10 80 Q 52.5 10, 95 80 T 180 80   ")
	require.Nil(t, err)
	require.Equal(t, Serialize(segments), "M 10 80 Q 52.5 10, 95 80 T 180 80")

	segments, err = ParsePath(" M240,100c50,0,0,125,50,100s0,-125,50,-150s175,50,50,100s-175,50,-300,0s0,-125,50,-100s0,125,50,150s0,-100,50,-100")
	require.Nil(t, err)
	require.Equal(t, Serialize(Absolutize(segments)), "M 240 100 C 290 100, 240 225, 290 200 S 290 75, 340 50 S 515 100, 390 150 S 215 200, 90 150 S 90 25, 140 50 S 140 175, 190 200 S 190 100, 240 100")
}
