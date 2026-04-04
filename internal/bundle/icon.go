package bundle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertPNGToICNS(pngPath, icnsPath string) error {
	iconsetDir := icnsPath + ".iconset"
	if err := os.MkdirAll(iconsetDir, 0o755); err != nil {
		return fmt.Errorf("creating iconset dir: %w", err)
	}
	defer os.RemoveAll(iconsetDir)

	sizes := []struct {
		name string
		size int
	}{
		{"icon_16x16.png", 16},
		{"icon_16x16@2x.png", 32},
		{"icon_32x32.png", 32},
		{"icon_32x32@2x.png", 64},
		{"icon_128x128.png", 128},
		{"icon_128x128@2x.png", 256},
		{"icon_256x256.png", 256},
		{"icon_256x256@2x.png", 512},
		{"icon_512x512.png", 512},
		{"icon_512x512@2x.png", 1024},
	}

	for _, s := range sizes {
		dest := filepath.Join(iconsetDir, s.name)
		cmd := exec.Command("sips", "-z", fmt.Sprint(s.size), fmt.Sprint(s.size), pngPath, "--out", dest)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("sips resize to %dx%d: %s: %w", s.size, s.size, out, err)
		}
	}

	cmd := exec.Command("iconutil", "-c", "icns", iconsetDir, "-o", icnsPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("iconutil: %s: %w", out, err)
	}

	return nil
}
