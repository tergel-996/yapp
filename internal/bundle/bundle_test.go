package bundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	appPath := filepath.Join(dir, "Yapp.app")

	opts := CreateOptions{
		AppPath:    appPath,
		BundleID:   "com.yapp.filemanager",
		BundleName: "Yapp",
		Version:    "1.0.0",
		YappBinary: "/usr/local/bin/yapp",
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

	// Check launcher script exists and is executable
	launcherPath := filepath.Join(appPath, "Contents", "MacOS", "Yapp")
	info, err := os.Stat(launcherPath)
	if err != nil {
		t.Fatalf("missing launcher: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Error("launcher is not executable")
	}

	// Check launcher content
	launcher, err := os.ReadFile(launcherPath)
	if err != nil {
		t.Fatalf("reading launcher: %v", err)
	}
	if !strings.Contains(string(launcher), "/usr/local/bin/yapp") {
		t.Error("launcher doesn't reference yapp binary")
	}
	if !strings.Contains(string(launcher), "launch") {
		t.Error("launcher doesn't call 'launch' subcommand")
	}
}
