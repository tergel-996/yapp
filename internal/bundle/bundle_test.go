package bundle

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	appPath := filepath.Join(dir, "Yapp.app")

	// Create a fake binary to copy into the bundle
	fakeBinary := filepath.Join(dir, "yapp-cli")
	fakeContent := []byte{0x7f, 0x45, 0x4c, 0x46} // ELF magic bytes (fake binary content)
	if err := os.WriteFile(fakeBinary, fakeContent, 0o755); err != nil {
		t.Fatal(err)
	}

	opts := CreateOptions{
		AppPath:    appPath,
		BundleID:   "com.yapp.filemanager",
		BundleName: "Yapp",
		Version:    "1.0.0",
		YappBinary: fakeBinary,
	}

	if err := Create(opts); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check directory structure
	for _, path := range []string{
		"Contents/MacOS",
		"Contents/Resources",
	} {
		full := filepath.Join(appPath, path)
		info, err := os.Stat(full)
		if err != nil {
			t.Errorf("missing directory %s: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", path)
		}
	}

	// Check Info.plist exists and has correct content
	plistData, err := os.ReadFile(filepath.Join(appPath, "Contents", "Info.plist"))
	if err != nil {
		t.Fatalf("missing Info.plist: %v", err)
	}
	if !strings.Contains(string(plistData), "com.yapp.filemanager") {
		t.Error("Info.plist missing bundle ID")
	}

	// Check launcher binary exists, is executable, and matches the source binary
	launcherPath := filepath.Join(appPath, "Contents", "MacOS", "Yapp")
	info, err := os.Stat(launcherPath)
	if err != nil {
		t.Fatalf("missing launcher: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Error("launcher is not executable")
	}
	launcher, err := os.ReadFile(launcherPath)
	if err != nil {
		t.Fatalf("reading launcher: %v", err)
	}
	if string(launcher) != string(fakeContent) {
		t.Error("launcher binary content doesn't match source binary")
	}
}

func TestConvertPNGToICNS(t *testing.T) {
	pngPath := filepath.Join(t.TempDir(), "icon.png")
	icnsPath := filepath.Join(t.TempDir(), "icon.icns")

	// Create a valid 256x256 PNG
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			img.Set(x, y, color.RGBA{R: 0x4A, G: 0x9E, B: 0xF5, A: 0xFF})
		}
	}
	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	err = ConvertPNGToICNS(pngPath, icnsPath)
	if err != nil {
		t.Fatalf("ConvertPNGToICNS failed: %v", err)
	}

	info, err := os.Stat(icnsPath)
	if err != nil {
		t.Fatalf("icns file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("icns file is empty")
	}
}
