package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tergel/yapp/internal/config"
	"github.com/tergel/yapp/internal/terminal"
)

func newLaunchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "launch [path]",
		Short: "Launch yazi in the configured terminal emulator",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runLaunch,
	}
}

func runLaunch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	yaziPath, err := exec.LookPath("yazi")
	if err != nil {
		return fmt.Errorf("yazi not found in PATH; install it with: brew install yazi")
	}

	var term terminal.Terminal
	all := terminal.All()

	if cfg.Terminal.Name == "auto" {
		term, err = terminal.Detect(all)
	} else {
		term, err = terminal.FindByName(cfg.Terminal.Name, all)
	}
	if err != nil {
		return err
	}

	yaziArgs := []string{}
	if len(args) > 0 {
		yaziArgs = append(yaziArgs, args[0])
	}

	launchCfg := terminal.LaunchConfig{
		Command:       yaziPath,
		Args:          yaziArgs,
		Title:         cfg.Appearance.Title,
		FontSize:      cfg.Appearance.FontSize,
		NoDecorations: !cfg.Appearance.WindowDecorations,
	}

	termArgs := term.BuildArgs(launchCfg)
	binary := term.Binary()

	proc := exec.Command(binary, termArgs...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}
