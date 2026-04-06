package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/config"
	"github.com/tergel/yapp/internal/terminal"
)

func newSetTerminalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set-terminal <name>",
		Short: "Set the terminal emulator (auto, ghostty, kitty, wezterm, alacritty, iterm, terminal)",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetTerminal,
	}
}

func runSetTerminal(cmd *cobra.Command, args []string) error {
	name := args[0]

	if name != "auto" {
		all := terminal.All()
		t, err := terminal.FindByName(name, all)
		if err != nil {
			valid := make([]string, len(all))
			for i, t := range all {
				valid[i] = t.Name()
			}
			return fmt.Errorf("unknown terminal %q; valid options: auto, %v", name, valid)
		}
		if !t.Detect() {
			fmt.Printf("Warning: %s is not installed on this system\n", name)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfg.Terminal.Name = name
	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("Terminal set to %q\n", name)
	return nil
}
