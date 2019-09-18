package imglib

import (
	"image"
	_ "image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersect(t *testing.T) {
	t.Parallel()
	in := loadImage(t, "testdata/circle_in.png")
	out := loadImage(t, "testdata/circle_out.png")

	assert.True(t, Intersect(in, in))
	assert.True(t, Intersect(out, out))
	assert.False(t, Intersect(in, out))
	assert.False(t, Intersect(out, in))
}

func loadImage(t *testing.T, path string) image.Image {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return img
}
