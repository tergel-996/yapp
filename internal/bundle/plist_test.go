package bundle

import (
	"encoding/xml"
	"strings"
	"testing"
)

// plistDict matches the handful of fields GeneratePlist writes. It lets
// tests parse the output and assert both structural validity and field
// values without string matching.
type plistDict struct {
	XMLName xml.Name `xml:"plist"`
	Dict    struct {
		Keys    []string `xml:"key"`
		Strings []string `xml:"string"`
	} `xml:"dict"`
}

// dictLookup walks the parallel Keys/Strings slices from plistDict and
// returns the string value for a given key, or "" if absent.
func (p *plistDict) dictLookup(key string) string {
	// The <key>/<string> elements alternate inside <dict>, but Go's xml
	// decoder surfaces them as two separate slices in document order.
	// Since Keys and Strings are each in order, and there are also
	// <true/>/<false/> elements between them, we can't just zip by index.
	// Instead, walk the raw key list and match by position among string
	// values, skipping the non-string keys (CFBundlePackageType is a
	// string so it's fine; NSHighResolutionCapable / LSUIElement are
	// booleans and fall outside the Strings slice).
	stringKeys := []string{}
	for _, k := range p.Dict.Keys {
		if k == "NSHighResolutionCapable" || k == "LSUIElement" {
			continue
		}
		stringKeys = append(stringKeys, k)
	}
	for i, k := range stringKeys {
		if k == key && i < len(p.Dict.Strings) {
			return p.Dict.Strings[i]
		}
	}
	return ""
}

func TestGeneratePlist(t *testing.T) {
	out := GeneratePlist(PlistConfig{
		BundleID:   "com.yapp.filemanager",
		BundleName: "Yapp",
		Executable: "Yapp",
		Version:    "1.0.0",
		IconFile:   "AppIcon",
	})

	var parsed plistDict
	if err := xml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("plist is not valid XML: %v\n%s", err, out)
	}

	cases := map[string]string{
		"CFBundleExecutable":        "Yapp",
		"CFBundleIdentifier":        "com.yapp.filemanager",
		"CFBundleName":              "Yapp",
		"CFBundlePackageType":       "APPL",
		"CFBundleVersion":           "1.0.0",
		"CFBundleShortVersionString": "1.0.0",
		"CFBundleIconFile":          "AppIcon",
	}
	for key, want := range cases {
		if got := parsed.dictLookup(key); got != want {
			t.Errorf("plist %s = %q, want %q", key, got, want)
		}
	}
}

// TestGeneratePlistEscapesHostileValues is a regression guard against
// XML injection via user-controlled config fields. If someone sets
// `bundle_id = "foo</string><string>bar"` in config.toml, the output must
// still parse as a single <string> element containing the literal text.
func TestGeneratePlistEscapesHostileValues(t *testing.T) {
	hostile := `foo</string><string>bar`
	out := GeneratePlist(PlistConfig{
		BundleID:   hostile,
		BundleName: `Yapp & "friends"`,
		Executable: "Yapp",
		Version:    "1.0.0",
		IconFile:   "AppIcon",
	})

	var parsed plistDict
	if err := xml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("hostile plist did not parse: %v\n%s", err, out)
	}

	if got := parsed.dictLookup("CFBundleIdentifier"); got != hostile {
		t.Errorf("CFBundleIdentifier did not round-trip hostile input.\n got:  %q\nwant: %q", got, hostile)
	}
	if got := parsed.dictLookup("CFBundleName"); got != `Yapp & "friends"` {
		t.Errorf(`CFBundleName did not round-trip. got: %q`, got)
	}

	// Sanity: raw output must not contain the unescaped closing tag from
	// the hostile bundle id, or we'd have successfully injected a second
	// <string> element.
	if strings.Contains(out, hostile) {
		t.Errorf("hostile value leaked into plist unescaped")
	}
}
