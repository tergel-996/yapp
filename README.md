# Yapp

A macOS identity wrapper for [yazi](https://github.com/sxyazi/yazi).

When you run yazi inside a terminal, it looks like any other window of that
terminal in Cmd+Tab — indistinguishable from a regular shell. Yapp gives yazi
its own labeled entry in Spotlight, Raycast, Dock, and Cmd+Tab so you can
always find the yazi session at a glance. Cmd+Tabbing to "Yapp" (or clicking
its Dock icon) brings the terminal window rendering yazi back to the
foreground.

## Why

Without Yapp, every yazi session shares its host terminal's Cmd+Tab identity
— you can have six Ghostty windows in your switcher and no way to tell which
one is yazi. Yapp adds a stable, clearly-labeled handle for the yazi session
so you can always get back to it.

**Heads up — what Yapp is not.** Yapp is an identity wrapper, not a terminal
emulator. macOS does not allow one app to "adopt" another app's windows, so
while Yapp is running you'll see **both** the "Yapp" entry **and** your
terminal emulator (Ghostty, Kitty, WezTerm, ...) in Cmd+Tab. The terminal is
still the thing actually drawing yazi; Yapp is a clearly-labeled handle
sitting alongside it. If you're looking for something that entirely replaces
your terminal in the switcher, no shell-only wrapper on macOS can do that
— it would need to ship its own terminal renderer, which Yapp intentionally
does not.

## Install

```bash
brew install tergel-996/yapp/yapp
yapp-cli install
```

## Usage

```bash
# Launch yazi as a standalone app
yapp-cli launch

# Launch with a specific path
yapp-cli launch ~/Downloads

# Or just open Yapp.app from Spotlight/Raycast/Dock
```

When launched from the `.app` bundle, Yapp runs as a real Cocoa application:
it shows up in Cmd+Tab with its own icon, lives in the Dock for the duration
of the yazi session, and has a minimal menu bar with `Hide` / `Quit Yapp`
(`⌘Q`). Cmd+Tabbing to Yapp — or clicking Yapp's Dock icon — brings the
terminal window rendering yazi back to the foreground, even if other windows
were covering it.

`⌘Q` in Yapp's menu quits **Yapp only** — the underlying terminal window and
yazi keep running. Quit yazi the usual way (`q`) and Yapp disappears from
Cmd+Tab within ~300 ms.

## Configuration

```bash
# Show current config
yapp-cli config show

# Edit config in $EDITOR
yapp-cli config edit

# Set terminal emulator
yapp-cli set-terminal ghostty    # or: kitty, wezterm, alacritty, iterm, terminal, auto
```

Config lives at `~/.config/yapp/config.toml`:

```toml
[terminal]
name = "auto"

[appearance]
font_size = 14
window_decorations = true
title = "Yapp"

[yazi]
# Optional override; Yapp otherwise finds yazi on PATH, then falls back
# to /opt/homebrew/bin/yazi and /usr/local/bin/yazi.
# path = "/custom/path/to/yazi"

[app]
bundle_id = "com.yapp.filemanager"
install_path = "~/Applications"
```

## Custom Icon

Yapp ships with an embedded default icon, so `yapp-cli install` produces
an `.app` bundle with a real icon in Spotlight, Dock, and Cmd+Tab out of
the box. To override it, pass `--icon`:

```bash
yapp-cli install --icon /path/to/icon.png
```

Provide a 1024x1024 PNG. Yapp converts it to `.icns` using macOS built-in
tools (`sips`, `iconutil`).

To regenerate the embedded default icon from source, run
`go run ./tools/icongen`.

## Supported Terminals

| Terminal | Detection | Notes |
|----------|-----------|-------|
| Ghostty | Auto | Full config support |
| Kitty | Auto | Full config support |
| WezTerm | Auto | `font_size` and `window_decorations` honored; `title` cannot be set via CLI on macOS |
| Alacritty | Auto | Full config support |
| iTerm2 | Auto | Via AppleScript |
| Terminal.app | Fallback | Via AppleScript, always available |

Auto-detection priority is the order above (Ghostty first, Terminal.app last).
Override with `yapp-cli set-terminal <name>`.

## Uninstall

```bash
yapp-cli uninstall
brew uninstall yapp
```

## Upgrading

`brew upgrade yapp` updates the `yapp-cli` binary. The next time you launch
Yapp.app, it automatically refreshes the bundle's copy of the binary from
Homebrew's, so you never run stale code from inside `~/Applications/Yapp.app`.
You don't need to re-run `yapp-cli install` after a brew upgrade.

## Troubleshooting

Launch errors in bundle mode are written to `~/Library/Logs/Yapp/yapp.log`.
If Yapp's Dock icon appears and disappears without opening a window, or
`yapp-cli launch` fails from Spotlight/Raycast with no visible error,
check that file — it records what went wrong (missing yazi, terminal
launch failure, config parse error, etc.).

## Releasing

The Homebrew formula lives in the separate tap repo at
[tergel-996/homebrew-yapp](https://github.com/tergel-996/homebrew-yapp).
Cut a tag here, then bump `url`/`sha256` in the tap's `Formula/yapp.rb`.

## License

MIT
