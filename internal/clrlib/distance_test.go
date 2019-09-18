package clrlib

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		c1, c2 color.Color
		want   float64
	}{
		{c1: color.White, c2: color.Black, want: 1},
		{c1: color.Black, c2: color.White, want: 1},
		{c1: color.White, c2: color.White, want: 0},
		{c1: color.Black, c2: color.Black, want: 0},
		{c1: color.RGBA{0, 0, 0, 255}, c2: color.RGBA{255, 255, 255, 0}, want: 1},
		{c1: color.RGBA{255, 255, 255, 0}, c2: color.RGBA{0, 0, 0, 255}, want: 1},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, Distance(tt.c1, tt.c2), "Distance(%+v, %+v)", tt.c1, tt.c2)
	}
}
