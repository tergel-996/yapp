package terminal

import (
	"strings"
	"testing"
)

// scriptPath is the canonical ScriptPath value used in terminal unit tests.
// launch.go is responsible for actually creating the file; these tests just
// verify each terminal impl threads the path through unchanged.
const scriptPath = "/var/folders/tmp/yapp-12345.sh"

func TestGhosttyBuildArgs(t *testing.T) {
	g := &Ghostty{}
	args := g.BuildArgs(LaunchConfig{
		ScriptPath:    scriptPath,
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"-na", "Ghostty.app", "--args",
		"--title=Yapp",
		"--font-size=16",
		"--window-decoration=false",
		"--quit-after-last-window-closed=true",
		"-e", scriptPath,
	}

	assertArgs(t, "ghostty", expected, args)
}

func TestKittyBuildArgs(t *testing.T) {
	k := &Kitty{}
	args := k.BuildArgs(LaunchConfig{
		ScriptPath:    scriptPath,
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"-na", "kitty.app", "--args",
		"--title", "Yapp",
		"-o", "font_size=16",
		"-o", "hide_window_decorations=yes",
		scriptPath,
	}

	assertArgs(t, "kitty", expected, args)
}

func TestWezTermBuildArgs(t *testing.T) {
	w := &WezTerm{}
	args := w.BuildArgs(LaunchConfig{
		ScriptPath:    scriptPath,
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"--config", "font_size=16",
		"--config", `window_decorations="NONE"`,
		"start",
		"--class", "Yapp",
		"--",
		scriptPath,
	}

	assertArgs(t, "wezterm", expected, args)
}

func TestWezTermBuildArgsDefaults(t *testing.T) {
	w := &WezTerm{}
	args := w.BuildArgs(LaunchConfig{
		ScriptPath: scriptPath,
		Title:      "Yapp",
	})

	expected := []string{
		"start",
		"--class", "Yapp",
		"--",
		scriptPath,
	}

	assertArgs(t, "wezterm defaults", expected, args)
}

func TestAlacrittyBuildArgs(t *testing.T) {
	a := &Alacritty{}
	args := a.BuildArgs(LaunchConfig{
		ScriptPath:    scriptPath,
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"-na", "Alacritty.app", "--args",
		"--title", "Yapp",
		"-o", "font.size=16",
		"-o", "window.decorations=None",
		"-e", scriptPath,
	}

	assertArgs(t, "alacritty", expected, args)
}

func TestITermName(t *testing.T) {
	i := &ITerm{}
	if i.Name() != "iterm" {
		t.Errorf("expected 'iterm', got %q", i.Name())
	}
}

func TestAppleTerminalName(t *testing.T) {
	a := &AppleTerminal{}
	if a.Name() != "terminal" {
		t.Errorf("expected 'terminal', got %q", a.Name())
	}
}

func TestITermBuildArgs(t *testing.T) {
	i := &ITerm{}
	args := i.BuildArgs(LaunchConfig{
		ScriptPath: scriptPath,
		Title:      "Yapp",
	})

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != "-e" {
		t.Errorf("expected first arg '-e', got %q", args[0])
	}
	if !strings.Contains(args[1], scriptPath) {
		t.Errorf("iTerm script did not embed script path; script: %s", args[1])
	}
	if !strings.Contains(args[1], `write text`) {
		t.Errorf("iTerm script did not use write text; script: %s", args[1])
	}
}

func TestAppleTerminalBuildArgs(t *testing.T) {
	a := &AppleTerminal{}
	args := a.BuildArgs(LaunchConfig{
		ScriptPath: scriptPath,
		Title:      "Yapp",
	})

	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d: %v", len(args), args)
	}
	if args[0] != "-e" {
		t.Errorf("expected first arg '-e', got %q", args[0])
	}
	if !strings.Contains(args[1], scriptPath) {
		t.Errorf("Terminal.app script did not embed script path; script: %s", args[1])
	}
	if !strings.Contains(args[1], `do script`) {
		t.Errorf("Terminal.app script did not use do script; script: %s", args[1])
	}
}

// TestAppleScriptEscapingInBuildArgs verifies that a ScriptPath containing
// a literal double quote is escaped when embedded in the AppleScript, so a
// hostile or unusual path (or title) can't break out of the string literal.
func TestAppleScriptEscapingInBuildArgs(t *testing.T) {
	nasty := `/tmp/"hi"/yapp.sh`
	for _, tc := range []struct {
		name string
		term Terminal
	}{
		{"iterm", &ITerm{}},
		{"terminal", &AppleTerminal{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := tc.term.BuildArgs(LaunchConfig{
				ScriptPath: nasty,
				Title:      `weird"title`,
			})
			if len(args) != 2 {
				t.Fatalf("expected 2 args, got %d", len(args))
			}
			if strings.Contains(args[1], nasty) {
				t.Errorf("unescaped script path leaked into AppleScript: %s", args[1])
			}
			if !strings.Contains(args[1], escapeAppleScriptString(nasty)) {
				t.Errorf("escaped script path not present in AppleScript: %s", args[1])
			}
		})
	}
}

func assertArgs(t *testing.T, name string, expected, got []string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("%s: expected %d args, got %d\nexpected: %v\ngot:      %v", name, len(expected), len(got), expected, got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("%s: arg[%d] = %q, want %q", name, i, got[i], expected[i])
		}
	}
}
