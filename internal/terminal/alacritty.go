package terminal

import (
	"fmt"
	"os/exec"
)

type Alacritty struct{}

func (a *Alacritty) Name() string { return "alacritty" }

func (a *Alacritty) Detect() bool {
	_, err := exec.LookPath("alacritty")
	return err == nil
}

func (a *Alacritty) Binary() string {
	path, err := exec.LookPath("alacritty")
	if err != nil {
		return "alacritty"
	}
	return path
}

func (a *Alacritty) BuildArgs(cfg LaunchConfig) []string {
	args := []string{
		"--title", cfg.Title,
	}
	if cfg.FontSize > 0 {
		args = append(args, "-o", fmt.Sprintf("font.size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "-o", "window.decorations=None")
	}
	args = append(args, "-e", cfg.Command)
	args = append(args, cfg.Args...)
	return args
}
