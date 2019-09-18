package clrlib

import (
	"image/color"
)

const maxDist float64 = 255 * 255 * 3

// Distance betwen two colors [0..1]
func Distance(c1, c2 color.Color) float64 {
	rgb1 := RGBA(c1)
	rgb2 := RGBA(c2)

	// Distance according to alpha: if both are transparent, they are the same. If only one of them
	// is transparent, they don't match.
	if rgb1.A == 0 && rgb2.A == 0 {
		return 0
	}
	if rgb1.A == 0 || rgb2.A == 0 {
		return 1
	}

	// Distance according to color components.
	dr := float64(rgb1.R) - float64(rgb2.R)
	dg := float64(rgb1.G) - float64(rgb2.G)
	db := float64(rgb1.B) - float64(rgb2.B)
	return (dr*dr + dg*dg + db*db) / maxDist
}
