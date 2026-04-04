package bundle

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateOptions struct {
	AppPath    string
	BundleID   string
	BundleName string
	Version    string
	YappBinary string
	IconPath   string
}

func Create(opts CreateOptions) error {
	macosDir := filepath.Join(opts.AppPath, "Contents", "MacOS")
	resourcesDir := filepath.Join(opts.AppPath, "Contents", "Resources")

	if err := os.MkdirAll(macosDir, 0o755); err != nil {
		return fmt.Errorf("creating MacOS dir: %w", err)
	}
	if err := os.MkdirAll(resourcesDir, 0o755); err != nil {
		return fmt.Errorf("creating Resources dir: %w", err)
	}

	plist := GeneratePlist(PlistConfig{
		BundleID:   opts.BundleID,
		BundleName: opts.BundleName,
		Executable: opts.BundleName,
		Version:    opts.Version,
		IconFile:   "AppIcon",
	})

	plistPath := filepath.Join(opts.AppPath, "Contents", "Info.plist")
	if err := os.WriteFile(plistPath, []byte(plist), 0o644); err != nil {
		return fmt.Errorf("writing Info.plist: %w", err)
	}

	launcher := fmt.Sprintf("#!/bin/bash\nexec \"%s\" launch \"$@\"\n", opts.YappBinary)

	launcherPath := filepath.Join(macosDir, opts.BundleName)
	if err := os.WriteFile(launcherPath, []byte(launcher), 0o755); err != nil {
		return fmt.Errorf("writing launcher: %w", err)
	}

	if opts.IconPath != "" {
		iconDest := filepath.Join(resourcesDir, "AppIcon.icns")
		data, err := os.ReadFile(opts.IconPath)
		if err != nil {
			return fmt.Errorf("reading icon: %w", err)
		}
		if err := os.WriteFile(iconDest, data, 0o644); err != nil {
			return fmt.Errorf("writing icon: %w", err)
		}
	}

	return nil
}

func Remove(appPath string) error {
	return os.RemoveAll(appPath)
}
