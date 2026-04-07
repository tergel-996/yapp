package bundle

import (
	"bytes"
	"image/png"
	"testing"
)

// TestDefaultIconPNG verifies the embedded PNG is present, decodes as a
// real image, and has dimensions appropriate for macOS .icns generation
// (the sips pipeline downsamples to sizes up to 1024x1024).
func TestDefaultIconPNG(t *testing.T) {
	if len(DefaultIconPNG) == 0 {
		t.Fatal("DefaultIconPNG is empty; did tools/icongen run?")
	}

	img, err := png.Decode(bytes.NewReader(DefaultIconPNG))
	if err != nil {
		t.Fatalf("decoding DefaultIconPNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() < 1024 || bounds.Dy() < 1024 {
		t.Errorf("default icon is %dx%d, want at least 1024x1024", bounds.Dx(), bounds.Dy())
	}
}
