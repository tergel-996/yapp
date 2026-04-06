package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Terminal   TerminalConfig   `toml:"terminal"`
	Appearance AppearanceConfig `toml:"appearance"`
	Yazi       YaziConfig       `toml:"yazi"`
	App        AppConfig        `toml:"app"`
}

type TerminalConfig struct {
	Name string `toml:"name"`
}

type AppearanceConfig struct {
	FontSize          int    `toml:"font_size"`
	WindowDecorations bool   `toml:"window_decorations"`
	Title             string `toml:"title"`
}

type YaziConfig struct {
	// Path is an optional absolute path to the yazi binary. When set, it
	// bypasses the PATH lookup and hardcoded Homebrew fallbacks in
	// launch.findYazi. A leading "~/" is expanded to the user's home dir.
	Path string `toml:"path"`
}

type AppConfig struct {
	BundleID    string `toml:"bundle_id"`
	InstallPath string `toml:"install_path"`
}

func Default() Config {
	return Config{
		Terminal: TerminalConfig{
			Name: "auto",
		},
		Appearance: AppearanceConfig{
			FontSize:          14,
			WindowDecorations: true,
			Title:             "Yapp",
		},
		App: AppConfig{
			BundleID:    "com.yapp.filemanager",
			InstallPath: filepath.Join(homeDir(), "Applications"),
		},
	}
}

// ExpandPath resolves a config-supplied path, turning a leading "~/" into an
// absolute home-relative path. Empty strings are returned unchanged.
func ExpandPath(p string) string {
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(homeDir(), p[2:])
	}
	if p == "~" {
		return homeDir()
	}
	return p
}

func Dir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "yapp")
	}
	return filepath.Join(homeDir(), ".config", "yapp")
}

func Path() string {
	return filepath.Join(Dir(), "config.toml")
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("HOME")
	}
	return home
}

func Load() (Config, error) {
	cfg := Default()
	path := Path()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path := Path()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(cfg)
}
