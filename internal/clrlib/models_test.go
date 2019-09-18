package clrlib

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuantize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		q       Quantize
		c, want color.Color
	}{
		{
			q:    Quantize(4),
			c:    color.RGBA{1, 50, 70, 255},
			want: color.RGBA{0, 63, 63, 255},
		},
		{
			q:    Quantize(1),
			c:    color.RGBA{1, 50, 127, 128},
			want: color.RGBA{0, 0, 0, 255},
		},
	}

	for _, tt := range tests {
		got := tt.q.Convert(tt.c)
		assert.Equal(t, tt.want, got)
	}
}
