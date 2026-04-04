package config

import (
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
	if cfg.Appearance.WindowDecorations != false {
		t.Error("expected window decorations false")
	}
	if cfg.Appearance.Title != "Yapp" {
		t.Errorf("expected title 'Yapp', got %q", cfg.Appearance.Title)
	}
	if cfg.App.BundleID != "com.yapp.filemanager" {
		t.Errorf("expected bundle ID 'com.yapp.filemanager', got %q", cfg.App.BundleID)
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

