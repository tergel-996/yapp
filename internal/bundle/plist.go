package bundle

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type PlistConfig struct {
	BundleID   string
	BundleName string
	Executable string
	Version    string
	IconFile   string
}

// GeneratePlist builds the Info.plist body for a Yapp-style .app bundle.
// All interpolated values are run through encoding/xml's escaper so a
// BundleName or BundleID containing `<`, `>`, `&`, or quotes cannot break
// out of the enclosing <string> element.
func GeneratePlist(cfg PlistConfig) string {
	const header = `<?xml version="1.0" encoding="UTF-8"?>
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
`
	return fmt.Sprintf(header,
		xmlEscape(cfg.Executable),
		xmlEscape(cfg.BundleID),
		xmlEscape(cfg.BundleName),
		xmlEscape(cfg.Version),
		xmlEscape(cfg.Version),
		xmlEscape(cfg.IconFile),
	)
}

func xmlEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}
