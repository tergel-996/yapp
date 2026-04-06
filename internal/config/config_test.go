package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()

	if cfg.Terminal.Name != "auto" {
		t.Errorf("expected terminal name 'auto', got %q", cfg.Terminal.Name)
	}
	if cfg.Appearance.FontSize != 14 {
		t.Errorf("expected font size 14, got %d", cfg.Appearance.FontSize)
	}
	if cfg.Appearance.WindowDecorations != true {
		t.Error("expected window decorations true")
	}
	if cfg.Appearance.Title != "Yapp" {
		t.Errorf("expected title 'Yapp', got %q", cfg.Appearance.Title)
	}
	if cfg.Yazi.Path != "" {
		t.Errorf("expected empty yazi path by default, got %q", cfg.Yazi.Path)
	}
	if cfg.App.BundleID != "com.yapp.filemanager" {
		t.Errorf("expected bundle ID 'com.yapp.filemanager', got %q", cfg.App.BundleID)
	}
}

func TestLoadYaziPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfgDir := filepath.Join(dir, "yapp")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(cfgDir, "config.toml"),
		[]byte("[yazi]\npath = \"/custom/yazi\"\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Yazi.Path != "/custom/yazi" {
		t.Errorf("yazi.path = %q, want %q", cfg.Yazi.Path, "/custom/yazi")
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"/abs/path", "/abs/path"},
		{"~", home},
		{"~/foo", filepath.Join(home, "foo")},
		{"./rel", "./rel"},
	}
	for _, c := range cases {
		if got := ExpandPath(c.in); got != c.want {
			t.Errorf("ExpandPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLoadReturnsDefaultsWhenNoFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := Default()
	expected.App.InstallPath = Default().App.InstallPath

	if cfg.Terminal.Name != expected.Terminal.Name {
		t.Errorf("expected terminal %q, got %q", expected.Terminal.Name, cfg.Terminal.Name)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg := Default()
	cfg.Terminal.Name = "ghostty"
	cfg.Appearance.FontSize = 18

	if err := Save(cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Terminal.Name != "ghostty" {
		t.Errorf("expected terminal 'ghostty', got %q", loaded.Terminal.Name)
	}
	if loaded.Appearance.FontSize != 18 {
		t.Errorf("expected font size 18, got %d", loaded.Appearance.FontSize)
	}
}

