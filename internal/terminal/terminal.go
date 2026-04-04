package terminal

type LaunchConfig struct {
	Command       string
	Args          []string
	Title         string
	FontSize      int
	NoDecorations bool
	ConfigFile    string
}

type Terminal interface {
	Name() string
	Detect() bool
	Binary() string
	BuildArgs(cfg LaunchConfig) []string
}
