package terminal

import "fmt"

func Detect(terminals []Terminal) (Terminal, error) {
	for _, t := range terminals {
		if t.Detect() {
			return t, nil
		}
	}
	return nil, fmt.Errorf("no supported terminal emulator found; install one of: ghostty, kitty, wezterm, alacritty, iterm2")
}

func FindByName(name string, terminals []Terminal) (Terminal, error) {
	for _, t := range terminals {
		if t.Name() == name {
			return t, nil
		}
	}
	return nil, fmt.Errorf("unknown terminal %q", name)
}

func All() []Terminal {
	return []Terminal{
		&Ghostty{},
		&Kitty{},
		&WezTerm{},
		&Alacritty{},
		&ITerm{},
		&AppleTerminal{},
	}
}
