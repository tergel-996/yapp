package terminal

import "testing"

type mockTerminal struct {
	name     string
	detected bool
}

func (m *mockTerminal) Name() string                        { return m.name }
func (m *mockTerminal) DisplayName() string                 { return m.name }
func (m *mockTerminal) Detect() bool                        { return m.detected }
func (m *mockTerminal) Binary() string                      { return "/usr/bin/" + m.name }
func (m *mockTerminal) BuildArgs(cfg LaunchConfig) []string { return nil }

func TestDetectReturnsFirstFound(t *testing.T) {
	terminals := []Terminal{
		&mockTerminal{name: "ghostty", detected: false},
		&mockTerminal{name: "kitty", detected: true},
		&mockTerminal{name: "wezterm", detected: true},
	}

	result, err := Detect(terminals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name() != "kitty" {
		t.Errorf("expected kitty, got %s", result.Name())
	}
}

func TestDetectReturnsErrorWhenNoneFound(t *testing.T) {
	terminals := []Terminal{
		&mockTerminal{name: "ghostty", detected: false},
		&mockTerminal{name: "kitty", detected: false},
	}

	_, err := Detect(terminals)
	if err == nil {
		t.Fatal("expected error when no terminal found")
	}
}

func TestFindByName(t *testing.T) {
	terminals := []Terminal{
		&mockTerminal{name: "ghostty", detected: true},
		&mockTerminal{name: "kitty", detected: true},
	}

	result, err := FindByName("kitty", terminals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name() != "kitty" {
		t.Errorf("expected kitty, got %s", result.Name())
	}
}

func TestFindByNameNotFound(t *testing.T) {
	terminals := []Terminal{
		&mockTerminal{name: "ghostty", detected: true},
	}

	_, err := FindByName("wezterm", terminals)
	if err == nil {
		t.Fatal("expected error for unknown terminal")
	}
}
