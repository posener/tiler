package imglib

import "image"

type Iterator struct {
	image.Point
	rect image.Rectangle
	d    image.Point
}

// d1 is the default dx,dy step.
var d1 = image.Point{X: 1, Y: 1}

func Iterate(rect image.Rectangle, d *image.Point) *Iterator {
	if d == nil {
		d = &d1
	}
	return &Iterator{
		Point: image.Point{X: rect.Min.X - d.X, Y: rect.Min.Y},
		rect:  rect,
		d:     *d,
	}
}

func (i *Iterator) Next() bool {
	i.X += i.d.X
	if i.X > i.rect.Max.X {
		i.X = i.rect.Min.X
		i.Y += i.d.Y
	}
	if i.Y > i.rect.Max.Y {
		return false
	}
	return true
}
