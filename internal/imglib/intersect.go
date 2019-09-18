package imglib

import "image"

func Intersect(im1, im2 image.Image) bool {
	rect := im1.Bounds().Intersect(im2.Bounds())
	for i := Iterate(rect, nil); i.Next(); {
		_, _, _, a1 := im1.At(i.X, i.Y).RGBA()
		_, _, _, a2 := im2.At(i.X, i.Y).RGBA()
		if a1 > 0 && a2 > 0 {
			return true
		}
	}
	return false
}
