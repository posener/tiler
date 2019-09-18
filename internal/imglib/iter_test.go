package imglib

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	t.Parallel()

	r := image.Rect(-1, -2, 1, 0)
	var got []image.Point
	for i := Iterate(r, &image.Point{1, 2}); i.Next(); {
		got = append(got, i.Point)
	}

	want := []image.Point{
		{-1, -2}, {0, -2}, {1, -2},
		{-1, 0}, {0, 0}, {1, 0},
	}
	assert.Equal(t, want, got)
}
