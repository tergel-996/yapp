package bundle

import "fmt"

type PlistConfig struct {
	BundleID   string
	BundleName string
	Executable string
	Version    string
	IconFile   string
}

func GeneratePlist(cfg PlistConfig) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>%s</string>
    <key>CFBundleIdentifier</key>
    <string>%s</string>
    <key>CFBundleName</key>
    <string>%s</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleVersion</key>
    <string>%s</string>
    <key>CFBundleShortVersionString</key>
    <string>%s</string>
    <key>CFBundleIconFile</key>
    <string>%s</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSUIElement</key>
    <false/>
</dict>
</plist>
`, cfg.Executable, cfg.BundleID, cfg.BundleName, cfg.Version, cfg.Version, cfg.IconFile)
}
