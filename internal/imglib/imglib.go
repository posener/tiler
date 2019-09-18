package imglib

import (
	"image"
	"image/color"
)

// RGBA copies the given image to an RBGA image.
func RGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	for i := Iterate(src.Bounds(), nil); i.Next(); {
		dst.Set(i.X, i.Y, src.At(i.X, i.Y))
	}
	return dst
}

// WithModel returns an image with different color model. The image is not copied.
func WithModel(parent image.Image, model color.Model) image.Image {
	return img{
		parent: parent,
		rect:   parent.Bounds(),
		model:  model,
	}
}

// SubImage returns the image in the intersection of the given image and rectangle.
// The image is not copied.
func SubImage(parent image.Image, rect image.Rectangle) image.Image {
	rect = rect.Intersect(parent.Bounds())
	if rect.Empty() {
		return &image.RGBA{}
	}
	return img{
		parent: parent,
		rect:   rect,
		model:  parent.ColorModel(),
	}
}

type img struct {
	parent image.Image
	rect   image.Rectangle
	model  color.Model
}

func (i img) Bounds() image.Rectangle {
	return i.rect
}

func (i img) ColorModel() color.Model {
	return i.model
}

func (i img) At(x, y int) color.Color {
	c := i.parent.At(x, y)
	if i.model != i.parent.ColorModel() {
		c = i.model.Convert(c)
	}
	return c
}

func Area(r image.Rectangle) int {
	return r.Dx() * r.Dy()
}
