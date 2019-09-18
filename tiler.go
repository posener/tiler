// Package tiler tiles a given image from a set of other images. It is a Go port of
// https://github.com/nuno-faria/tiler.
package tiler

import (
	"image"
	"image/draw"
	"log"
	"sort"
	"sync"

	"github.com/posener/tiler/internal/imglib"
	"github.com/posener/tiler/internal/mode"
)

// Config is the configration of the tiling process.
type Config struct {
	// Shift defines the location for which matching the tile will happen. The given (x,y) result in
	// a grid of boxes in size (x,y), on which the tiles will be
	// positioned.
	Shift image.Point
	// Overlap is whether to allow tiles to overlap after being matched.
	Overlap bool
	// TilesPermute is the configuration of the tiles permutations.
	TilesPermute PermuteConfig
}

// UpdateFn is a function for updating on any change to the given image.
type UpdateFn func(img image.Image)

// Tile matches the given tiles with the given configuration over the given image. The tiled image
// is returned in the output.
func Tile(img image.Image, tiles []image.Image, cfg Config, update UpdateFn) image.Image {
	log.Printf("Computing tiles permutations...")
	perms := Permute(tiles, cfg.TilesPermute)
	log.Printf("Using %d tiles permutations!", len(perms))

	log.Printf("Computing tiles matches...")
	matches := computeMatches(img, perms, cfg.Shift)
	log.Printf("Computed tiles matching in %d locations", len(matches))

	log.Print("Composing output...")
	return composeMatches(img.Bounds(), matches, cfg.Overlap, update)
}

// match represetns matching of a tile to a location in the image.
type match struct {
	// which tiled is matched.
	tile image.Image
	// to which area in the original image is it being matched.
	location image.Rectangle
	// how far is it from the original image area.
	distance float64
}

// intersect checks if the match's tile intersects with a corresponding patch in the given image.
// It is used to check if a new tile overlaps existing drawn image.
func (m match) intersect(img *image.RGBA) bool {
	patch := img.SubImage(m.location).(*image.RGBA)
	// Move the rect of the path to the origin.
	patch.Rect = patch.Rect.Sub(m.location.Min)
	return imglib.Intersect(patch, m.tile)
}

// computeMatches computes a 'match' for each tile, according to the distance from boxes
// defined over the image.
func computeMatches(img image.Image, tiles []mode.Mode, shift image.Point) []match {
	// Map tiles according to their size, to improve performance: This result in gridding the image
	// only once, and test all tiles with the same size against the same grid.
	mapped := make(map[image.Point][]mode.Mode)
	for _, tile := range tiles {
		size := tile.Bounds().Size()
		mapped[size] = append(mapped[size], tile)
	}

	var (
		matches []match
		wg      sync.WaitGroup
		lock    sync.Mutex
	)

	// Compute for all the tiles.
	wg.Add(len(mapped))
	for size := range mapped {
		go func(size image.Point) {
			defer wg.Done()

			// Compute for each box (a sub image of the original image) of the
			// current tile size, with the required shift.
			var sizeMatches []match
			for _, box := range grid(img, size, shift) {
				boxMode := mode.New(box, true)
				tile, dist := closestMode(boxMode, mapped[size])
				if tile == nil {
					continue
				}
				sizeMatches = append(sizeMatches,
					match{tile: tile, location: box.Bounds(), distance: dist})
			}

			// Add the matches from the current size to all the matches.
			lock.Lock()
			defer lock.Unlock()
			matches = append(matches, sizeMatches...)
		}(size)
	}
	wg.Wait()
	return matches
}

// composeMatches places the matches over the canvas. It places them in two modes:
// * No overlap: The ones that are closest (smallest distances to image box) and largest are placed
//   first, then other are placed with no overlap.
// * With overlap: All the matches are placed, starting from the most distant and largest.
func composeMatches(rect image.Rectangle, matches []match, overlap bool, update UpdateFn) image.Image {
	log.Printf("Sorting matches...")
	sort.Slice(matches, func(i, j int) bool { return less(matches[i], matches[j], overlap) })

	log.Printf("Placing matches...")
	out := image.NewRGBA(rect)
	for _, match := range matches {
		if !overlap && match.intersect(out) {
			continue
		}
		draw.Draw(out, match.location, match.tile, image.ZP, draw.Over)
		update(out)
	}
	return out
}

func less(left, right match, overlap bool) bool {
	if left.distance == right.distance {
		return imglib.Area(left.tile.Bounds()) > imglib.Area(right.tile.Bounds())
	}
	ret := left.distance < right.distance
	if overlap {
		ret = !ret
	}
	return ret

}

// grid returns a list of subimages of the given image according to the given
// grid size and shift.
func grid(img image.Image, size, shift image.Point) []image.Image {
	if shift.Eq(image.ZP) {
		shift = size
	}
	var boxes []image.Image
	for i := imglib.Iterate(img.Bounds(), &shift); i.Next(); {
		boxes = append(boxes, imglib.SubImage(img, image.Rectangle{Min: i.Point, Max: i.Add(size)}))
	}
	return boxes
}

// closestMode returns the image closes mode and its distances.
func closestMode(m mode.Mode, others []mode.Mode) (image.Image, float64) {
	// if the mode is transparent, return no match.
	if _, _, _, a := m.RGBA(); a == 0 {
		return nil, 0
	}

	minMode := others[0]
	minDist := float64(1)
	for _, other := range others {
		dist := m.Distance(other)
		if dist < minDist {
			minDist = dist
			minMode = other
		}
	}
	return minMode, minDist
}
