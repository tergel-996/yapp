package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Terminal   TerminalConfig   `toml:"terminal"`
	Appearance AppearanceConfig `toml:"appearance"`
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
