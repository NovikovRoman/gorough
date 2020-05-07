package gorough

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPointsOnPath(t *testing.T) {
	points, err := PointsOnPath("M240,100c50,0,0,125,50,100s0,-125,50,-150s175,50,50,100s-175,50,-300,0s0,-125,50,-100s0,125,50,150s0,-100,50,-100", 0, 0)
	require.Nil(t, err)

	require.Len(t, points, 1)
	require.Len(t, points[0], 252)
	require.Equal(t, points[0][4], Point{X: 246.39266967773438, Y: 100.79565048217773})
}
