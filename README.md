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

## Configuration

```bash
# Show current config
yapp-cli config show

# Edit config in $EDITOR
yapp-cli config edit

# Set terminal emulator
yapp-cli set-terminal ghostty    # or: kitty, wezterm, alacritty, iterm, terminal, auto

# Register as default folder handler (experimental)
yapp-cli register
```

Config lives at `~/.config/yapp/config.toml`:

```toml
[terminal]
name = "auto"

[appearance]
font_size = 14
window_decorations = false
title = "Yapp"

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
| WezTerm | Auto | Full config support |
| Alacritty | Auto | Full config support |
| iTerm2 | Auto | Via AppleScript |
| Terminal.app | Fallback | Via AppleScript, always available |

## Uninstall

```bash
yapp-cli uninstall
brew uninstall yapp
```

## License

MIT
