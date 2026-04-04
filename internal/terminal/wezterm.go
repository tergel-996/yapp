package terminal

import "os/exec"

type WezTerm struct{}

func (w *WezTerm) Name() string { return "wezterm" }

func (w *WezTerm) Detect() bool {
	_, err := exec.LookPath("wezterm")
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
