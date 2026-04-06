package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/bundle"
	"github.com/tergel/yapp/internal/config"
)

func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the Yapp.app bundle",
		RunE:  runUninstall,
	}
}

func runUninstall(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	appPath := filepath.Join(cfg.App.InstallPath, "Yapp.app")

	fmt.Printf("Removing %s...\n", appPath)
	if err := bundle.Remove(appPath); err != nil {
		return fmt.Errorf("removing bundle: %w", err)
	}

	fmt.Printf("Yapp.app uninstalled. Config at %s was kept.\n", config.Dir())
	return nil
}
