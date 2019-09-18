package mode

import (
	"image"
	"image/color"
	"math"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/posener/tiler/internal/clrlib"
	"github.com/posener/tiler/internal/imglib"
)

const quant = clrlib.Quantize(32)

// Mode represents an image with its most common color.
type Mode struct {
	image.Image
	color.Color
	Freq float64
}

// New returns a Mode of an image. if useTransparent is set, the transparent color will be
// participating in the calculation of the most common color.
func New(img image.Image, useTransparent bool) Mode {
	var (
		counter             = make(map[color.Color]float64)
		common  color.Color = color.Transparent
		total   float64
	)
	for i := imglib.Iterate(img.Bounds(), nil); i.Next(); {
		c := quant.Convert(img.At(i.X, i.Y))
		if _, _, _, a := c.RGBA(); a == 0 {
			if useTransparent {
				// All completely transparent color should be counted with the same counter,
				// regradless of their RGB values.
				c = color.Transparent
			} else {
				continue
			}
		}

		total++
		counter[c]++
		if counter[c] > counter[common] {
			common = c
		}
	}
	return Mode{
		Image: img,
		Color: common,
		Freq:  counter[common] / total,
	}
}

// Returns a scaled copy of the mode.
func (m Mode) Scale(scale float64) Mode {
	m.Image = scaleImage(m.Image, scale)
	return m
}

// Returns a rotated copy of the mode.
func (m Mode) Rotate(rotation float64) Mode {
	m.Image = rotateImage(m.Image, rotation)
	return m
}

func (m Mode) Distance(other Mode) float64 {
	return clrlib.Distance(m.Color, other.Color) / m.Freq / other.Freq
}

func scaleImage(img image.Image, scale float64) image.Image {
	dx, dy := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	dx, dy = scale*dx, scale*dy
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(dx)), int(math.Ceil(dy))))
	graphics.Scale(dst, img)
	return dst
}

func rotateImage(img image.Image, rotation float64) image.Image {
	angle := 2 * math.Pi * rotation
	dx, dy := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	cos, sin := math.Cos(angle), math.Sin(angle)
	dx, dy = dx*cos+dy*sin, dx*sin+dy*cos
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(dx)), int(math.Ceil(dy))))
	graphics.Rotate(dst, img, &graphics.RotateOptions{Angle: angle})
	return dst
}
