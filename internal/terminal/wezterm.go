package terminal

import (
	"os"
	"os/exec"
)

type WezTerm struct{}

func (w *WezTerm) Name() string { return "wezterm" }

func (w *WezTerm) Detect() bool {
	if _, err := exec.LookPath("wezterm"); err == nil {
		return true
	}
	_, err := os.Stat("/Applications/WezTerm.app")
	return err == nil
}

func (w *WezTerm) Binary() string {
	path, err := exec.LookPath("wezterm")
	if err != nil {
		return "wezterm"
	}
	return path
}

func (w *WezTerm) BuildArgs(cfg LaunchConfig) []string {
	args := []string{
		"start",
		"--class", cfg.Title,
		"--",
		cfg.Command,
	}
	args = append(args, cfg.Args...)
	return args
}
