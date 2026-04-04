package terminal

import (
	"fmt"
	"os/exec"
)

type Ghostty struct{}

func (g *Ghostty) Name() string { return "ghostty" }

func (g *Ghostty) Detect() bool {
	_, err := exec.LookPath("ghostty")
	return err == nil
}

func (g *Ghostty) Binary() string {
	path, err := exec.LookPath("ghostty")
	if err != nil {
		return "ghostty"
	}
	return path
}

func (g *Ghostty) BuildArgs(cfg LaunchConfig) []string {
	args := []string{
		fmt.Sprintf("--title=%s", cfg.Title),
	}
	if cfg.FontSize > 0 {
		args = append(args, fmt.Sprintf("--font-size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "--window-decoration=false")
	}
	args = append(args, "--quit-after-last-window-closed=true")
	args = append(args, "-e", cfg.Command)
	args = append(args, cfg.Args...)
	return args
}
