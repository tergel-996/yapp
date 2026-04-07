package bundle

import (
	"fmt"
	"os"
	"os/exec"
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

	// Copy the yapp-cli binary directly -- the .app executable must be a
	// native Mach-O binary. A shell script causes Rosetta crashes on Apple Silicon.
	// The binary detects invocation as "Yapp" (via os.Args[0]) and runs launch.
	launcherPath := filepath.Join(macosDir, opts.BundleName)
	data, err := os.ReadFile(opts.YappBinary)
	if err != nil {
		return fmt.Errorf("reading yapp binary: %w", err)
	}
	if err := os.WriteFile(launcherPath, data, 0o755); err != nil {
		return fmt.Errorf("writing launcher binary: %w", err)
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

	// Re-sign the bundle with an ad-hoc signature whose identifier
	// matches CFBundleIdentifier. Without this step the Mach-O binary
	// inherits Go's linker-generated signature (`Identifier=a.out`,
	// Info.plist not bound, sealed resources=none) which disagrees with
	// the enclosing bundle's declared identity. On recent macOS (tested
	// on Tahoe / Darwin 25) that mismatch silently prevents the running
	// process from appearing in Cmd+Tab even when LaunchServices has
	// already registered it as type=Foreground. Re-signing makes the
	// binary and the Info.plist agree cryptographically.
	if err := signBundle(opts.AppPath, opts.BundleName, opts.BundleID); err != nil {
		return fmt.Errorf("signing bundle: %w", err)
	}

	return nil
}

func Remove(appPath string) error {
	return os.RemoveAll(appPath)
}

// signBundle ad-hoc signs the bundle inside-out: the Mach-O binary at
// Contents/MacOS/<BundleName> first, then the bundle wrapper itself.
// Both get the same --identifier (the CFBundleIdentifier) so that the
// binary's signed identity matches the enclosing bundle's declared id,
// which is a prerequisite for modern macOS Cmd+Tab visibility.
//
// `--sign -` is an ad-hoc signature (no Developer ID); that is fine for
// a tool installed from source via a Homebrew tap -- the goal is just a
// properly-bound signature, not Gatekeeper notarization.
func signBundle(appPath, bundleName, bundleID string) error {
	binaryPath := filepath.Join(appPath, "Contents", "MacOS", bundleName)

	// Sign the inner Mach-O first. Apple recommends inside-out signing
	// order for bundles; --deep is deprecated and does not always bind
	// the parent plist correctly.
	if out, err := exec.Command(
		"/usr/bin/codesign",
		"--force",
		"--sign", "-",
		"--identifier", bundleID,
		binaryPath,
	).CombinedOutput(); err != nil {
		return fmt.Errorf("signing %s: %s: %w", binaryPath, out, err)
	}

	// Then sign the bundle wrapper. This binds Info.plist into the
	// signature and produces the sealed-resources manifest.
	if out, err := exec.Command(
		"/usr/bin/codesign",
		"--force",
		"--sign", "-",
		"--identifier", bundleID,
		appPath,
	).CombinedOutput(); err != nil {
		return fmt.Errorf("signing %s: %s: %w", appPath, out, err)
	}

	return nil
}
