package bundle

import (
	"strings"
	"testing"
)

func TestGeneratePlist(t *testing.T) {
	plist := GeneratePlist(PlistConfig{
		BundleID:   "com.yapp.filemanager",
		BundleName: "Yapp",
		Executable: "Yapp",
		Version:    "1.0.0",
		IconFile:   "AppIcon",
	})

	checks := []struct {
		label string
		want  string
	}{
		{"bundle id", "<string>com.yapp.filemanager</string>"},
		{"bundle name", "<string>Yapp</string>"},
		{"executable", "<string>Yapp</string>"},
		{"version", "<string>1.0.0</string>"},
		{"icon", "<string>AppIcon</string>"},
		{"package type", "<string>APPL</string>"},
		{"high res", "<true/>"},
	}

	for _, c := range checks {
		if !strings.Contains(plist, c.want) {
			t.Errorf("plist missing %s: expected to contain %q", c.label, c.want)
		}
	}
}
