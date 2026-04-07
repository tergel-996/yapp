package bundle

import _ "embed"

// DefaultIconPNG is the PNG-encoded default Yapp application icon,
// embedded at build time from defaulticon.png. Regenerate with:
//
//	go run ./tools/icongen
//
// It is used by `yapp-cli install` as the fallback icon when the user
// does not pass --icon. Having a real default is load-bearing: strangers
// installing Yapp for the first time would otherwise see the generic
// macOS binary icon in Spotlight/Dock, which defeats the entire point
// of wrapping yazi in its own .app in the first place.
//
//go:embed defaulticon.png
var DefaultIconPNG []byte
