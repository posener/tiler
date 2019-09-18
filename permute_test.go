package tiler

import (
	"image/color"
	"testing"

	"github.com/posener/tiler/internal/clrlib"
	"github.com/stretchr/testify/assert"
)

func TestPermuteColors(t *testing.T) {
	t.Parallel()

	got := permuteColors(0, 1, 2)
	assert.Equal(t, got, []color.Model{clrlib.Scaled{1, 1, 0}, clrlib.Scaled{1, 1, 1}})
}
