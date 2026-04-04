package terminal

import "testing"

func TestGhosttyBuildArgs(t *testing.T) {
	g := &Ghostty{}
	args := g.BuildArgs(LaunchConfig{
		Command:       "/usr/local/bin/yazi",
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"--title=Yapp",
		"--font-size=16",
		"--window-decoration=false",
		"--quit-after-last-window-closed=true",
		"-e", "/usr/local/bin/yazi",
	}

	assertArgs(t, "ghostty", expected, args)
}

func TestKittyBuildArgs(t *testing.T) {
	k := &Kitty{}
	args := k.BuildArgs(LaunchConfig{
		Command:       "/usr/local/bin/yazi",
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"--title", "Yapp",
		"-o", "font_size=16",
		"-o", "hide_window_decorations=yes",
		"/usr/local/bin/yazi",
	}

	assertArgs(t, "kitty", expected, args)
}

func TestWezTermBuildArgs(t *testing.T) {
	w := &WezTerm{}
	args := w.BuildArgs(LaunchConfig{
		Command:       "/usr/local/bin/yazi",
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"start",
		"--class", "Yapp",
		"--",
		"/usr/local/bin/yazi",
	}

	assertArgs(t, "wezterm", expected, args)
}

func TestAlacrittyBuildArgs(t *testing.T) {
	a := &Alacritty{}
	args := a.BuildArgs(LaunchConfig{
		Command:       "/usr/local/bin/yazi",
		Title:         "Yapp",
		FontSize:      16,
		NoDecorations: true,
	})

	expected := []string{
		"--title", "Yapp",
		"-o", "font.size=16",
		"-o", "window.decorations=None",
		"-e", "/usr/local/bin/yazi",
	}

	assertArgs(t, "alacritty", expected, args)
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
