package tiler

import (
	"image"
	"image/color"
	"sync"

	"github.com/posener/tiler/internal/clrlib"
	"github.com/posener/tiler/internal/imglib"
	"github.com/posener/tiler/internal/mode"
)

// PermuteConfig is configuration for the Permute function.
type PermuteConfig struct {
	// NumR, NumG and NumB can configure the specific number of a color component.
	// To use only the original color set it to 0.
	NumR, NumG, NumB uint8
	// Scale creates different scale variants of the given image list.
	Scale []float64
	// Scale creates different rotation variants of the given image list. Values should be in range
	// [0..1]. A rotation of 0 does not rotate the image and a rotation of 1 is 360 degrees.
	Rotate []float64
}

// Permute returns a list of permutations of the provided images, according to the premutation
// configuration. Permute with empty configuration returns the mode of the given images.
func Permute(in []image.Image, cfg PermuteConfig) []mode.Mode {

	if len(cfg.Scale) == 0 {
		cfg.Scale = []float64{1}
	}
	if len(cfg.Rotate) == 0 {
		cfg.Rotate = []float64{0}
	}

	var (
		out    []mode.Mode
		colors = permuteColors(cfg.NumR, cfg.NumG, cfg.NumB)
		wg     sync.WaitGroup
		lock   sync.Mutex
	)

	wg.Add(len(in))
	for _, img := range in {
		go func(img image.Image) {
			defer wg.Done()
			perms := premuteImage(img, colors, cfg.Scale, cfg.Rotate)
			lock.Lock()
			defer lock.Unlock()
			out = append(out, perms...)
		}(img)
	}
	wg.Wait()
	return out
}

func premuteImage(img image.Image, colors []color.Model, scales []float64, rotations []float64) []mode.Mode {
	if img.Bounds().Empty() {
		return nil
	}
	var perms []mode.Mode
	for _, colorModel := range colors {
		// Color the image and calculate mode.
		img := mode.New(imglib.WithModel(img, colorModel), false)

		// Generate tiles in all requested scales.
		for _, scale := range scales {
			img := img.Scale(scale)
			for _, rotation := range rotations {
				img := img.Rotate(rotation)
				perms = append(perms, img)
			}
		}
	}
	return perms
}

// permuteColors returns a list of models that contains all permutations according to the
// number of required permutations of each color component.
func permuteColors(nr, ng, nb uint8) []color.Model {
	var colorModels []color.Model
	for _, r := range iterate(nr) {
		for _, g := range iterate(ng) {
			for _, b := range iterate(nb) {
				colorModels = append(colorModels, clrlib.Scaled{R: r, G: g, B: b})
			}
		}
	}
	return colorModels
}

// iterate returns a slice of float between [0,1] according to the number of given stpes.
func iterate(steps uint8) []float64 {
	if steps <= 1 {
		return []float64{1}
	}
	ret := make([]float64, 0, steps)
	step := 1 / float64(steps-1)
	for i := uint8(0); i < steps; i++ {
		ret = append(ret, float64(i)*step)
	}
	return ret
}
