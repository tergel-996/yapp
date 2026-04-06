package terminal

import "testing"

func TestEscapeAppleScriptString(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", `yazi /tmp`, `yazi /tmp`},
		{"double quote", `say "hi"`, `say \"hi\"`},
		{"backslash", `a\b`, `a\\b`},
		{"backslash before quote", `a\"b`, `a\\\"b`},
		{"newline", "a\nb", `a\nb`},
		{"carriage return", "a\rb", `a\rb`},
		{"injection attempt", `"; do shell script "rm -rf /"; --`,
			`\"; do shell script \"rm -rf /\"; --`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := escapeAppleScriptString(c.in)
			if got != c.want {
				t.Errorf("escape(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
