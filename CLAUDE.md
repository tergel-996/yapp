# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Yapp wraps [yazi](https://github.com/sxyazi/yazi) in a macOS `.app` bundle so it gets its own identity in Spotlight, Dock, Cmd+Tab, and Raycast — separate from whatever terminal emulator actually renders it.

Distribution is via a Homebrew tap at [`tergel-996/homebrew-yapp`](https://github.com/tergel-996/homebrew-yapp). The formula lives **only** in the tap repo; there is no in-repo copy to keep in sync.

## Commands

```bash
make build          # builds ./bin/yapp-cli with version ldflag from `git describe`
make vet            # go vet ./...
make test           # go test ./... -v
make clean
make install        # copies bin/yapp-cli to /usr/local/bin/yapp-cli
make install-app    # build + run `yapp-cli install` (regenerates ~/Applications/Yapp.app)

go test ./internal/terminal/...                    # run one package
go test ./internal/bundle -run TestGeneratePlist   # run one test
```

The module path is `github.com/tergel/yapp` but the **binary is `yapp-cli`**, not `yapp`. This is intentional: `yapp` collides with Perl's `Parse::Yapp` at `/usr/bin/yapp`. Don't "fix" this. The version string is injected via `-ldflags "-X github.com/tergel/yapp/internal/cli.Version=..."`.

## Architecture

### Dual-mode binary, Cocoa runloop, self-heal (`cmd/yapp/main.go` + `internal/gui`)

The same binary serves two roles and dispatches on `os.Args[0]`:

- Invoked as `yapp-cli` → normal Cobra CLI. No Cocoa, no gui initialisation beyond static linker cost.
- Invoked as `Yapp` (i.e. from `Yapp.app/Contents/MacOS/Yapp`) → runs under a Cocoa `NSApplication` event loop. This is what makes macOS treat Yapp as a real GUI app: **Cmd+Tab visibility, Dock icon, Apple Event handling, no "not responding" dialog from LaunchServices.**

**Why the Cocoa runloop is load-bearing.** LaunchServices dispatches `kAEOpenApplication` to every `APPL`-type bundle it launches and expects an `NSApplication`-level response within ~8–15 seconds. A plain Go CLI binary that merely stays alive in a `for { stat(marker) }` loop is *not* a real Cocoa app from LaunchServices' perspective — it gets evicted from Cmd+Tab and eventually triggers the "application is not responding" dialog. The fix isn't to run longer; it's to actually be an NSApplication. We do that via cgo in `internal/gui/`.

The Yapp-mode launch sequence in `main.go`:

1. `runtime.LockOSThread()` pins the Go main goroutine to the process's main OS thread. Cocoa requires its run loop to execute on this specific thread.
2. `selfHealBundle()` stat-compares the currently running bundle binary against `/opt/homebrew/bin/yapp-cli` (or the Intel equivalent) and atomically replaces the bundle binary if Homebrew has a newer one. Silent, best-effort; closes the "stale bundle after `brew upgrade`" gap. The current run still executes the mmap'd old binary; the replacement takes effect next launch.
3. A worker goroutine runs `cli.Execute()` (which dispatches to the `launch` subcommand → writes marker + script, spawns terminal, polls marker). When the poll returns, the worker calls `gui.Stop()`.
4. The main goroutine calls `gui.Run()`, which blocks in `[NSApp run]` until the worker calls `gui.Stop()` (normal case) or the user hits Cmd+Q (hard `[NSApp terminate:]`).

`cli.Execute()` returns `error` rather than calling `os.Exit` directly, specifically so main can unwind the runloop cleanly.

### Cocoa package (`internal/gui`)

`internal/gui/cocoa_darwin.go` + `cocoa_darwin.m` is a cgo wrapper around a minimal Cocoa app. **Keep the `@implementation` block in the `.m` file, not in the .go cgo preamble** — cgo will otherwise compile the preamble into more than one translation unit and the Objective-C class symbols become duplicates at link time (I learned this the hard way).

What it gives us:

- **`Run()`** — initialises `NSApplication`, sets `NSApplicationActivationPolicyRegular` (so the app has a Dock icon and Cmd+Tab entry), installs a `YappAppDelegate`, builds a minimal main menu (Hide / Hide Others / Show All / Quit), activates the app, and enters `[NSApp run]`. Must be called from the thread that started `main()` (hence `runtime.LockOSThread`).
- **`Stop()`** — dispatches `[NSApp stop:]` + a synthetic event to the main queue so the run loop wakes up and exits. Safe to call from any goroutine.
- **`SetActivateHandler(func())`** — registers a callback invoked from two distinct reactivation paths in the delegate: `applicationDidBecomeActive:` (fires on Cmd+Tab switches and Dock-icon clicks that transition Yapp from inactive to active) and `applicationShouldHandleReopen:hasVisibleWindows:` (fires on Dock-icon clicks even when Yapp is *already* active — without this, clicking the Dock icon while you're already in Yapp's menu bar would be a no-op). `launch.go` registers a closure that runs `osascript -e 'tell application "<terminal DisplayName>" to activate'` to refocus the spawned terminal window. **Both handler paths are necessary** — dropping either one breaks a real user action (removing `didBecomeActive` breaks Cmd+Tab reactivation; removing `shouldHandleReopen` breaks Dock click while already active).
- **`//export yappCocoaHandleReopen`** — the Go end of the reactivation callback. Gates out the *first* call via the `firstActivationSeen` flag because `[NSApp activateIgnoringOtherApps:YES]` during `yapp_run()` fires `applicationDidBecomeActive:` before the worker goroutine has spawned the terminal — running the handler at that moment would either no-op (terminal not yet running) or activate a pre-existing instance of the same terminal app (not the one that will host yazi). All subsequent calls run the registered handler on a new goroutine so the Cocoa main thread isn't blocked by `exec.Command`.

The menu's Quit item is bound to `[NSApp terminate:]`, which hard-exits Yapp and **does not** run Go deferred cleanups on the worker goroutine. That means Cmd+Q while yazi is still running leaks the marker + script files in `$TMPDIR` until the shell's EXIT trap fires (on yazi exit) or the OS sweeps `/tmp`. Acceptable for now; if it becomes a problem, the fix is channel plumbing from the Cocoa delegate back to the worker goroutine to signal a graceful shutdown.

`internal/gui/stub_other.go` (`//go:build !darwin`) provides no-op stubs so `go vet` / IDE tooling works on Linux. Yapp itself is macOS-only.

When modifying `main.go` or `gui/`, preserve: the `LockOSThread` call, the worker-goroutine + `gui.Run()` structure, the `@implementation` being in the `.m` file, and the self-heal call before the argv rewrite.

### `.app` bundle creation (`internal/bundle`)

`bundle.Create` builds `Yapp.app/` by writing `Info.plist` and **copying the current `yapp-cli` binary** to `Contents/MacOS/Yapp`. It must be a native Mach-O binary, not a shell script — shell scripts as the bundle executable cause Rosetta crashes on Apple Silicon. The install flow in `internal/cli/install.go` resolves its own executable path via `os.Executable` + `filepath.EvalSymlinks` (important because Homebrew installs via a symlink) and passes it to `bundle.Create`.

All user-controlled fields that land in `Info.plist` (`BundleID`, `BundleName`, `Version`, `IconFile`) are escaped through `xml.EscapeText` in `plist.go`. Don't switch this back to raw `fmt.Sprintf` interpolation — a config with `bundle_id = "foo</string><string>bar"` would otherwise inject a second `<string>` element into the plist. There's a regression test (`TestGeneratePlistEscapesHostileValues`).

**Ad-hoc code signing.** After copying the binary and writing the plist, `bundle.Create` calls `signBundle` which shells out to `/usr/bin/codesign --force --sign - --identifier <BundleID>` inside-out (the Mach-O at `Contents/MacOS/<BundleName>` first, then the bundle wrapper). This step is **load-bearing on macOS Tahoe (Darwin 25) and later**: without it, the binary ships with Go's linker-generated ad-hoc signature whose identifier is `a.out`, the `Info.plist` is not bound to the signature, and `Sealed Resources` is empty. That mismatch between the binary's signed identity and the enclosing bundle's declared identity **silently prevents the running process from appearing in Cmd+Tab** even though `lsappinfo` correctly reports it as `type="Foreground"` and `NSApplicationActivationPolicyRegular` has been set. The symptom is brutal: the app is technically running and registered, but just doesn't show in the switcher, with no error anywhere. The fix is to re-sign with an identifier that matches `CFBundleIdentifier`, which produces a properly bound signature (matching ad-hoc-signed apps like Ghostty's Developer-ID-signed one in structure if not trust level). Use `codesign -dvvv ~/Applications/Yapp.app` to diagnose: a broken install shows `Identifier=a.out`, `Info.plist=not bound`, `Sealed Resources=none`; a correct one shows `Identifier=com.yapp.filemanager`, `Info.plist entries=N`, `Sealed Resources version=2 rules=13 files=1`. `--deep` is deprecated and does not reliably bind parent plists; use the explicit inside-out sequence.

Icons: `ConvertPNGToICNS` shells out to macOS `sips` + `iconutil` to produce a multi-resolution `.icns` from a single PNG. `internal/bundle/defaulticon.png` is embedded into the binary via `//go:embed` in `default_icon.go` and used by `internal/cli/install.go` as the fallback when the user does not pass `--icon`, so a fresh install always ends up with a real icon in Spotlight/Dock/Cmd+Tab rather than the generic Unix-executable icon. The PNG is reproducible from source via `go run ./tools/icongen` (stdlib-only SDF renderer; no image-processing deps).

### Terminal strategy (`internal/terminal`)

`Terminal` is an interface implemented once per supported emulator (`ghostty.go`, `kitty.go`, `wezterm.go`, `alacritty.go`, `iterm.go`, `apple_terminal.go`). `All()` returns them in detection priority order; `Detect()` picks the first one present on disk; `FindByName` resolves the explicit `terminal.name` config value. Each impl also exposes a `DisplayName()` string — the macOS `.app` display name used by the Dock-click handler in `gui` (`tell application "Ghostty" to activate`, etc.).

`LaunchConfig` carries a `ScriptPath` string — the absolute path to an executable `/bin/sh` script that `launch.go` writes to `$TMPDIR`. Each terminal impl passes that path to its native runner:

- **Ghostty / Alacritty:** `open -na App.app --args ... -e <scriptPath>`
- **Kitty:** `open -na kitty.app --args ... <scriptPath>` (command is positional at end)
- **WezTerm:** `wezterm [--config K=V]... start --class Yapp -- <scriptPath>`. WezTerm does not expose a CLI flag for the window title on macOS, so `cfg.Title` is intentionally not honored; the `--class` flag is still set for the OS-level window class.
- **iTerm2:** AppleScript `write text "<scriptPath>"` on a new session. iTerm's `create window ... command "X"` tokenizes X via execvp-style splitting and mangles paths with special characters, so `write text` is the safer path.
- **Terminal.app:** AppleScript `do script "<scriptPath>"`, which runs the script in a new Terminal tab.

**Why a script file and not `/bin/sh -c "..."`?** Ghostty's `-e` flag joins all following argv elements into a single bash `-c` string, which breaks any inner `/bin/sh -c` argument boundary (bash re-tokenises and `sh` ends up with just the literal word `trap` as its `-c` arg). Passing a single-word script path is unambiguous regardless of how each terminal tokenises its command line. Don't try to "simplify" this back to `sh -c` inline — you'll reintroduce the original Ghostty bug, and the fix is not obvious from symptoms (Ghostty reports "failed to launch" after ~20ms).

Two non-obvious rules govern `Binary()` / `BuildArgs()` on macOS:

1. **GUI terminals are launched via `/usr/bin/open -na App.app --args ...`**, not by calling their CLI binary directly. The `ghostty`/`kitty`/etc. binaries in `PATH` are CLI helpers that can't open new windows when invoked from a non-terminal context (like a `.app` bundle).
2. **Ghostty-specific:** always pass `--window-decoration=true|false` explicitly. Ghostty reads its own config file otherwise and will override the user's `window_decorations` setting. See commit `07ac075`.

Any string interpolated into an AppleScript literal must go through `escapeAppleScriptString` (`applescript.go`) — both `cfg.ShellCommand` and `cfg.Title` are passed through user config and would otherwise be a command-injection vector for the iTerm and Terminal.app paths.

When adding a new terminal, implement `Terminal`, wrap `cfg.ShellCommand` via `/bin/sh -c` (or your native equivalent), add the impl to `All()` in `detect.go` in the desired priority position, and add it to the switch in `internal/cli/set_terminal.go`.

### Launch path + marker-file wait (`internal/cli/launch.go`)

`runLaunch` → `config.Load` → `findYazi(cfg)` → pick terminal → register Dock-click activate handler via `gui.SetActivateHandler` → create per-PID marker file → write shell script to `$TMPDIR` → `term.BuildArgs(LaunchConfig{ScriptPath: ...})` → `exec.Command(term.Binary(), args...).Run()` → `waitForMarkerGone(marker)`.

**The marker-file wait is a secondary keep-alive** (the primary is now the Cocoa runloop in `cmd/yapp/main.go`). Even with Cocoa keeping Yapp alive from LaunchServices' perspective, the launch goroutine still needs to know when yazi exits so it can call `gui.Stop()` and let the process terminate. Terminal-launch commands (`open -na`, `wezterm start`, `osascript`) all return in ~50ms while yazi is still running, so we can't rely on the terminal-launch process exit alone.

The script file at `$TMPDIR/yapp-<pid>.sh` looks like:

```sh
#!/bin/sh
trap 'rm -f '\''$TMPDIR/yapp-<pid>.lock'\''' EXIT HUP INT TERM
'/opt/homebrew/bin/yazi' '/path/arg'
```

**Important: no `exec` before the yazi line.** `exec` would replace the shell with yazi and the `EXIT` trap would never fire, leaking the marker forever. The shell forks yazi, waits for it, then exits normally so the trap runs. Costs ~1 MB RSS of extra shell process for the session.

The `trap` covers both clean yazi exit **and** a force-closed terminal window (shell receives `SIGHUP` → trap fires → `rm` runs). Both paths were verified end-to-end (normal exit ~1s, SIGHUP ~50ms). A per-PID marker means two concurrent Yapp invocations don't interfere with each other.

All shell interpolation goes through `shQuote` (POSIX single-quote escaping with `'\''` for embedded quotes). The trap body itself is single-quoted and escaped via the same helper — do not rewrite this to use double quotes without thinking about `$`/backtick injection via the marker path.

`findYazi` resolution order: `config.Yazi.Path` (with `~/` expansion) → `exec.LookPath("yazi")` → `/opt/homebrew/bin/yazi` → `/usr/local/bin/yazi`. The hardcoded fallbacks exist because when `Yapp.app` is launched from Finder/Spotlight, `PATH` is only `/usr/bin:/bin:/usr/sbin:/sbin`. Any new dependency on an external binary needs the same treatment.

### Config (`internal/config`)

TOML at `~/.config/yapp/config.toml` (or `$XDG_CONFIG_HOME/yapp/config.toml`). Sections:

- `[terminal] name` — `auto` or one of the supported terminal IDs.
- `[appearance] font_size`, `window_decorations`, `title` — honored per-terminal; WezTerm ignores `title`.
- `[yazi] path` — optional absolute path; overrides auto-detection. Leading `~/` is expanded via `ExpandPath`.
- `[app] bundle_id`, `install_path` — where to install `Yapp.app`.

`Load()` returns defaults when the file doesn't exist — there's no "not configured" error state. Yapp does **not** manage yazi's own config; yazi reads its own `~/.config/yazi/` independently.

## CI

`.github/workflows/ci.yml` runs `go vet`, `go build`, and `go test` on `macos-latest` for pushes and PRs to `main`. The pipeline runs on macOS (not Linux) because:

1. A few tests shell out to `sips` / `iconutil` / `osascript`, which only exist on macOS.
2. The `internal/gui` package is cgo + AppKit and only builds on darwin. The `stub_other.go` build tag makes `go vet` pass on Linux for local tooling, but the real binary cannot be produced off-macOS.
