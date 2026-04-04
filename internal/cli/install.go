package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/bundle"
	"github.com/tergel/yapp/internal/config"
)

func newInstallCmd() *cobra.Command {
	var iconPath string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Generate and install the Yapp.app bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(iconPath)
		},
	}

	cmd.Flags().StringVar(&iconPath, "icon", "", "path to PNG icon (optional, will be converted to icns)")

	return cmd
}

func runInstall(iconPath string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	yappBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding yapp binary: %w", err)
	}
	yappBinary, err = filepath.EvalSymlinks(yappBinary)
	if err != nil {
		return fmt.Errorf("resolving yapp binary path: %w", err)
	}

	appPath := filepath.Join(cfg.App.InstallPath, "Yapp.app")

	// Convert icon if provided
	icnsPath := ""
	if iconPath != "" {
		icnsPath = filepath.Join(os.TempDir(), "yapp-icon.icns")
		fmt.Printf("Converting icon %s to icns...\n", iconPath)
		if err := bundle.ConvertPNGToICNS(iconPath, icnsPath); err != nil {
			return fmt.Errorf("converting icon: %w", err)
		}
		defer os.Remove(icnsPath)
	}

	fmt.Printf("Installing Yapp.app to %s...\n", appPath)

	opts := bundle.CreateOptions{
		AppPath:    appPath,
		BundleID:   cfg.App.BundleID,
		BundleName: "Yapp",
		Version:    Version,
		YappBinary: yappBinary,
		IconPath:   icnsPath,
	}

	if err := bundle.Create(opts); err != nil {
		return fmt.Errorf("creating bundle: %w", err)
	}

	// Initialize default config if it doesn't exist
	if _, err := os.Stat(config.Path()); os.IsNotExist(err) {
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving default config: %w", err)
		}
		fmt.Printf("Created config at %s\n", config.Path())
	}

	fmt.Println("Yapp.app installed successfully.")
	fmt.Println("You can find it in Spotlight, Raycast, or your Applications folder.")
	return nil
}
