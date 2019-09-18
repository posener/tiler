package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/posener/tiler"
)

var (
	imgPath   = flag.String("img", "", "Image to tile. Required.")
	tilesPath = flag.String("tiles", "", "Path to tiles directory or a tile file. Required.")
	outPath   = flag.String("out", "", "Destination path.")
	shift     = flag.String("shift", "", "Grid shifts in the format: 'x,y'. If omitted, tile size will be used.")
	colors    = flag.String("colors", "", `Scale tiles colors.
Use a number 'n' to define number of scales of each color component.
Use comma separated numbers 'r,g,b' to have different number of scales to each color component.`)
	scale   = flag.String("scale", "", "Scale tiles. Comma separated list of scale factors.")
	rotate  = flag.String("rotate", "", "Rotate tiles. Comma separated list of rotations in range [0..1].")
	overlap = flag.Bool("overlap", false, "Can tiles overlap each other.")
)

func main() {
	flag.Parse()

	if *imgPath == "" {
		log.Fatalf("img flag is required.")
	}
	if *tilesPath == "" {
		log.Fatalf("tiles flag is required.")
	}

	if *outPath == "" {
		*outPath = "tiled.png"
	}

	log.Print("Loading image...")
	img, err := loadImage(*imgPath)
	if err != nil {
		log.Fatalf("Failed loading image %s: %s", *imgPath, err)
	}

	log.Print("Loading tiles...")
	tiles, err := loadTiles(*tilesPath)
	if err != nil {
		log.Fatalf("Failed loading tiles: %s", err)
	}
	log.Printf("Loaded %d tiles", len(tiles))

	if len(tiles) == 0 {
		log.Fatal("No tiles found")
	}

	// TODO: add interactive output here.
	updateFn := func(img image.Image) {}

	cfg := config()
	log.Printf("Tiling with config: %+v", cfg)
	out := tiler.Tile(img, tiles, cfg, updateFn)

	log.Print("Saving result...")
	err = saveImage(*outPath, out)
	if err != nil {
		log.Fatalf("Failed saving output to %q: %s", *outPath, err)
	}
	log.Printf("Done! created %s.", *outPath)
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func loadTiles(path string) ([]image.Image, error) {
	f, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		img, err := loadImage(path)
		return []image.Image{img}, err
	}

	var images []image.Image
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".png") {
			return nil
		}
		img, err := loadImage(path)
		if err != nil {
			return fmt.Errorf("loading tile %q: %w", path, err)
		}
		images = append(images, img)
		return nil
	})
	return images, err
}

func saveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func config() tiler.Config {
	var cfg tiler.Config
	cfg.Overlap = *overlap

	var err error
	if *colors != "" {
		cfg.TilesPermute.NumR, cfg.TilesPermute.NumG, cfg.TilesPermute.NumB, err = parseColors(*colors)
		if err != nil {
			log.Fatalf("Bad value for colors: %s", err)
		}
	}
	if *shift != "" {
		cfg.Shift, err = parsePoint(*shift)
		if err != nil {
			log.Fatalf("Bad value for shift: %s", err)
		}
	}
	if *scale != "" {
		cfg.TilesPermute.Scale, err = parseFloat(*scale)
		if err != nil {
			log.Fatalf("Bad value for scales: %s", err)
		}
	}
	if *rotate != "" {
		cfg.TilesPermute.Rotate, err = parseFloat(*rotate)
		if err != nil {
			log.Fatalf("Bad value for rotations: %s", err)
		}
	}
	return cfg
}

func parseColors(s string) (r uint8, g uint8, b uint8, err error) {
	parts := strings.Split(s, ",")
	if len(parts) == 1 {
		r, err = parseUInt8(parts[0])
		if err != nil {
			return
		}
		g = r
		b = r
		return
	}
	if len(parts) != 3 {
		err = fmt.Errorf("must be of the form 'n' or 'r,g,b'")
		return
	}
	r, err = parseUInt8(parts[0])
	if err != nil {
		err = fmt.Errorf("bad value for r (%s): %s", parts[0], err)
		return
	}
	g, err = parseUInt8(parts[1])
	if err != nil {
		err = fmt.Errorf("bad value for g (%s): %s", parts[1], err)
		return
	}
	b, err = parseUInt8(parts[2])
	if err != nil {
		err = fmt.Errorf("bad value for b (%s): %s", parts[2], err)
	}
	return
}

func parseUInt8(s string) (uint8, error) {
	n, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}

func parsePoint(s string) (image.Point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return image.Point{}, fmt.Errorf("must be of the form x,y")
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return image.Point{}, fmt.Errorf("bad value for x (%s): %s", parts[0], err)
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return image.Point{}, fmt.Errorf("bad value for y (%s): %s", parts[1], err)
	}
	return image.Point{x, y}, nil
}

func parseFloat(s string) ([]float64, error) {
	parts := strings.Split(s, ",")
	var ret []float64
	for _, part := range parts {
		f, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("bad float format for %s: %w", part, err)
		}
		ret = append(ret, f)
	}
	return ret, nil
}
