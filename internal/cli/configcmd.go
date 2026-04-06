package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show or edit yapp configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.Path())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Print the current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			fmt.Printf("terminal.name = %q\n", cfg.Terminal.Name)
			fmt.Printf("appearance.font_size = %d\n", cfg.Appearance.FontSize)
			fmt.Printf("appearance.window_decorations = %v\n", cfg.Appearance.WindowDecorations)
			fmt.Printf("appearance.title = %q\n", cfg.Appearance.Title)
			fmt.Printf("yazi.path = %q\n", cfg.Yazi.Path)
			fmt.Printf("app.bundle_id = %q\n", cfg.App.BundleID)
			fmt.Printf("app.install_path = %q\n", cfg.App.InstallPath)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "edit",
		Short: "Open config in $EDITOR",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.Path()
			if _, err := os.Stat(path); os.IsNotExist(err) {
				cfg := config.Default()
				if err := config.Save(cfg); err != nil {
					return err
				}
			}
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			proc := exec.Command(editor, path)
			proc.Stdin = os.Stdin
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr
			return proc.Run()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Create default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.Path()
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config already exists at %s", path)
			}
			cfg := config.Default()
			if err := config.Save(cfg); err != nil {
				return err
			}
			fmt.Printf("Config created at %s\n", path)
			return nil
		},
	})

	return cmd
}
