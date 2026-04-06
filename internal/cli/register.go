package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/config"
	"github.com/tergel/yapp/internal/handler"
)

func newRegisterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register Yapp as the default folder handler",
		RunE:  runRegister,
	}
}

func runRegister(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Registering %s as folder handler...\n", cfg.App.BundleID)
	if err := handler.Register(cfg.App.BundleID); err != nil {
		return err
	}

	fmt.Println("Done. Yapp is now the default folder handler.")
	fmt.Println("Note: this may not take effect in all contexts. Finder is deeply integrated into macOS.")
	return nil
}
