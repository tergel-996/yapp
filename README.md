# Yapp

Yazi as a standalone macOS application.

Yapp wraps the [yazi](https://github.com/sxyazi/yazi) terminal file manager in its own macOS app bundle, giving it a separate identity from your terminal emulator. It shows up in Spotlight, Raycast, Dock, and Cmd+Tab as its own app.

## Why

If you use yazi inside your terminal, opening it creates another terminal window. You end up with multiple instances of the same terminal in your app switcher. Yapp solves this by giving yazi its own application identity.

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
(`⌘Q`). Clicking Yapp's Dock icon brings the spawned terminal window back to
the foreground.

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

```bash
yapp-cli install --icon /path/to/icon.png
```

Provide a 1024x1024 PNG. Yapp converts it to icns using macOS built-in tools.

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

## Releasing

The Homebrew formula lives in the separate tap repo at
[tergel-996/homebrew-yapp](https://github.com/tergel-996/homebrew-yapp).
Cut a tag here, then bump `url`/`sha256` in the tap's `Formula/yapp.rb`.

## License

MIT
