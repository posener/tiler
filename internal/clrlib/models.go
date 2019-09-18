package clrlib

import (
	"image/color"
	"math"
)

// Quantize reduces the number of colors in an image.
// Valid values are between 0 to 255. 0 will make no quantization.
type Quantize uint8

func (q Quantize) Convert(c color.Color) color.Color {
	if q < 1 {
		return c
	}
	rgba := RGBA(c)
	rgba.R = q.quantize(rgba.R)
	rgba.G = q.quantize(rgba.G)
	rgba.B = q.quantize(rgba.B)
	rgba.A = q.quantize(rgba.A)
	return rgba
}

func (q Quantize) quantize(i uint8) uint8 {
	return uint8(math.Round(float64(i)*float64(q)/255) * 255 / float64(q))
}

type Scaled struct {
	R, G, B float64
}

func (s Scaled) Convert(c color.Color) color.Color {
	rgba := RGBA(c)
	rgba.R = scale(rgba.R, s.R)
	rgba.G = scale(rgba.G, s.G)
	rgba.B = scale(rgba.B, s.B)
	return rgba
}

func scale(c uint8, s float64) uint8 {
	return uint8(float64(c) * s)
}
