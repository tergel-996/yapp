package terminal

import "strings"

// escapeAppleScriptString escapes a Go string so it can be safely embedded
// inside an AppleScript string literal (delimited by double quotes).
//
// AppleScript string literals recognize \\, \", \n, \r, and \t as escape
// sequences. Backslash must be escaped first, then the double quote.
// Newlines and carriage returns are escaped so a multi-line value cannot
// break out of the surrounding quoted expression.
func escapeAppleScriptString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}
