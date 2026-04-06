package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/config"
	"github.com/tergel/yapp/internal/gui"
	"github.com/tergel/yapp/internal/terminal"
)

// findYazi locates the yazi binary. When launched from a .app bundle, the
// system PATH is limited to /usr/bin:/bin:/usr/sbin:/sbin and does not include
// Homebrew, so we fall back to known installation paths.
//
// Resolution order:
//  1. config.Yazi.Path (if set and the file exists)
//  2. PATH lookup
//  3. /opt/homebrew/bin/yazi  (Apple Silicon Homebrew)
//  4. /usr/local/bin/yazi     (Intel Homebrew)
func findYazi(cfg config.Config) (string, error) {
	if cfg.Yazi.Path != "" {
		expanded := config.ExpandPath(cfg.Yazi.Path)
		if _, err := os.Stat(expanded); err == nil {
			return expanded, nil
		}
		return "", fmt.Errorf("yazi.path %q does not exist", expanded)
	}
	if path, err := exec.LookPath("yazi"); err == nil {
		return path, nil
	}
	for _, candidate := range []string{
		"/opt/homebrew/bin/yazi", // Apple Silicon Homebrew
		"/usr/local/bin/yazi",    // Intel Homebrew
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("yazi not found; install with: brew install yazi, or set yazi.path in config.toml")
}

// shQuote wraps a string in POSIX single quotes, escaping any embedded single
// quotes. Safe for interpolation into /bin/sh script files.
func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func newLaunchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "launch [path]",
		Short: "Launch yazi in the configured terminal emulator",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runLaunch,
	}
}

func runLaunch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	yaziPath, err := findYazi(cfg)
	if err != nil {
		return err
	}

	var term terminal.Terminal
	all := terminal.All()

	if cfg.Terminal.Name == "auto" {
		term, err = terminal.Detect(all)
	} else {
		term, err = terminal.FindByName(cfg.Terminal.Name, all)
	}
	if err != nil {
		return err
	}

	// Register the Dock-icon-click handler so that clicking Yapp's Dock
	// icon (or Cmd+clicking it from Mission Control) brings the spawned
	// terminal window back to the foreground. The handler is a no-op in
	// CLI mode because gui.Run was never called -- SetActivateHandler
	// just stashes the closure in a package-level var either way.
	termDisplayName := term.DisplayName()
	gui.SetActivateHandler(func() {
		_ = exec.Command(
			"/usr/bin/osascript",
			"-e",
			fmt.Sprintf(`tell application "%s" to activate`, termDisplayName),
		).Run()
	})

	// Create a per-PID marker file. The launch script inside the terminal
	// traps EXIT/HUP/INT/TERM and removes the marker when yazi ends. This
	// process then polls for the marker's absence and exits.
	//
	// While we are polling, the Yapp binary is still the running process
	// that owns Yapp.app's bundle identity, so Yapp shows up in Cmd+Tab for
	// the entire lifetime of the yazi session. As soon as the marker
	// disappears, we return and the .app vanishes from the switcher.
	marker, markerCleanup, err := createMarker()
	if err != nil {
		return fmt.Errorf("creating marker: %w", err)
	}
	defer markerCleanup()

	script, scriptCleanup, err := createLaunchScript(marker, yaziPath, args)
	if err != nil {
		return fmt.Errorf("creating launch script: %w", err)
	}
	defer scriptCleanup()

	launchCfg := terminal.LaunchConfig{
		ScriptPath:    script,
		Title:         cfg.Appearance.Title,
		FontSize:      cfg.Appearance.FontSize,
		NoDecorations: !cfg.Appearance.WindowDecorations,
	}

	termArgs := term.BuildArgs(launchCfg)
	binary := term.Binary()

	proc := exec.Command(binary, termArgs...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	if err := proc.Run(); err != nil {
		return fmt.Errorf("launching terminal: %w", err)
	}

	// Terminal-launch commands (open -na, wezterm start, osascript) all
	// return once the terminal window has been requested. Block here until
	// the script wrapper removes the marker.
	waitForMarkerGone(marker)
	return nil
}

// createMarker writes a zero-byte file in the system temp dir keyed to this
// process's PID. It returns the path and a cleanup func that removes the
// file if it still exists (for the case where we exit without the launch
// script having the chance to run).
func createMarker() (string, func(), error) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("yapp-%d.lock", os.Getpid()))
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		return "", func() {}, err
	}
	cleanup := func() { _ = os.Remove(path) }
	return path, cleanup, nil
}

// createLaunchScript writes an executable /bin/sh script that:
//  1. restores a useful PATH (the .app sandbox gives only
//     /usr/bin:/bin:/usr/sbin:/sbin, which hides fzf / ripgrep / zoxide
//     and every other yazi shell-out plugin),
//  2. installs a trap to remove the marker file on every exit path
//     (clean exit, SIGHUP from force-closed window, SIGINT, SIGTERM),
//     and
//  3. runs yazi with the caller's positional args, then waits for it.
//
// Crucially this does NOT use `exec` to replace the shell. `exec` would
// hand the process over to yazi and the shell would vanish, taking the
// EXIT trap with it and leaving the marker behind forever. Instead the
// shell forks yazi, waits for it, and then exits normally so the EXIT
// trap fires. The extra shell process costs ~1 MB RSS for the duration
// of the yazi session.
//
// We use a real script file -- not a `sh -c "expr"` invocation -- because
// terminals like Ghostty concatenate `-e` arguments into a single bash -c
// string before running them, which breaks any inner `/bin/sh -c` argument
// boundary and ends up running the literal `trap` builtin with no effect.
// Passing a one-word script path through `-e` is unambiguous regardless of
// how each terminal tokenises its command line.
//
// PATH restoration order:
//  1. `/usr/libexec/path_helper -s` to pick up /etc/paths and /etc/paths.d/
//     (this is what a normal login shell does on macOS).
//  2. Prepend /opt/homebrew/bin, /usr/local/bin, ~/.cargo/bin, ~/.local/bin
//     as a belt-and-braces fallback for users whose tools aren't covered
//     by path_helper.
//  3. Source ~/.config/yapp/env.sh if present, so users with exotic
//     layouts (Nix, asdf, pyenv shims, ~/.bin, etc.) have a single
//     escape hatch without modifying login shell rc files.
func createLaunchScript(marker, yaziPath string, userArgs []string) (string, func(), error) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("yapp-%d.sh", os.Getpid()))

	// Single-quoted pieces so a path with spaces or special characters
	// cannot break the script.
	parts := []string{shQuote(yaziPath)}
	for _, a := range userArgs {
		parts = append(parts, shQuote(a))
	}
	yaziCmd := strings.Join(parts, " ")

	trapBody := "rm -f " + shQuote(marker)

	const pathSetup = `if [ -x /usr/libexec/path_helper ]; then
    eval "$(/usr/libexec/path_helper -s)"
fi
export PATH="/opt/homebrew/bin:/usr/local/bin:$HOME/.cargo/bin:$HOME/.local/bin:$PATH"
if [ -f "$HOME/.config/yapp/env.sh" ]; then
    . "$HOME/.config/yapp/env.sh"
fi
`
	content := fmt.Sprintf(
		"#!/bin/sh\n%strap %s EXIT HUP INT TERM\n%s\n",
		pathSetup,
		shQuote(trapBody),
		yaziCmd,
	)

	if err := os.WriteFile(path, []byte(content), 0o700); err != nil {
		return "", func() {}, err
	}
	cleanup := func() { _ = os.Remove(path) }
	return path, cleanup, nil
}

// waitForMarkerGone polls the marker file and returns when it is removed.
// 300 ms is short enough to feel snappy when yazi exits and long enough not
// to burn CPU during a long session.
func waitForMarkerGone(marker string) {
	for {
		if _, err := os.Stat(marker); os.IsNotExist(err) {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}
