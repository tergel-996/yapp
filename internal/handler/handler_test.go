package handler

import "testing"

func TestBuildRegisterCommands(t *testing.T) {
	cmds := BuildRegisterCommands("com.yapp.filemanager")

	if len(cmds) == 0 {
		t.Fatal("expected at least one command")
	}

	foundHandler := false
	for _, cmd := range cmds {
		if cmd.Binary == "defaults" || cmd.Binary == "duti" {
			foundHandler = true
		}
	}
	if !foundHandler {
		t.Error("expected defaults or duti command")
	}
}

func TestBuildUnregisterCommands(t *testing.T) {
	cmds := BuildUnregisterCommands()

	if len(cmds) == 0 {
		t.Fatal("expected at least one command")
	}
}
