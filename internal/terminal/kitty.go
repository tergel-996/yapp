package terminal

import (
	"fmt"
	"os/exec"
)

type Kitty struct{}

func (k *Kitty) Name() string { return "kitty" }

func (k *Kitty) Detect() bool {
	_, err := exec.LookPath("kitty")
	return err == nil
}

func (k *Kitty) Binary() string {
	path, err := exec.LookPath("kitty")
	if err != nil {
		return "kitty"
	}
	return path
}

func (k *Kitty) BuildArgs(cfg LaunchConfig) []string {
	args := []string{
		"--title", cfg.Title,
	}
	if cfg.FontSize > 0 {
		args = append(args, "-o", fmt.Sprintf("font_size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "-o", "hide_window_decorations=yes")
	}
	args = append(args, cfg.Command)
	args = append(args, cfg.Args...)
	return args
}
