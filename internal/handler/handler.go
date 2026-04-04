package handler

import (
	"fmt"
	"os/exec"
)

type ShellCommand struct {
	Binary string
	Args   []string
}

func BuildRegisterCommands(bundleID string) []ShellCommand {
	return []ShellCommand{
		{
			Binary: "defaults",
			Args: []string{
				"write", "com.apple.LaunchServices/com.apple.launchservices.secure",
				"LSHandlers", "-array-add",
				fmt.Sprintf(`{LSHandlerContentType="public.folder";LSHandlerRoleAll="%s";}`, bundleID),
			},
		},
	}
}

func BuildUnregisterCommands() []ShellCommand {
	return []ShellCommand{
		{
			Binary: "/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister",
			Args:   []string{"-kill", "-r", "-domain", "local", "-domain", "system", "-domain", "user"},
		},
	}
}

func Register(bundleID string) error {
	for _, cmd := range BuildRegisterCommands(bundleID) {
		proc := exec.Command(cmd.Binary, cmd.Args...)
		if out, err := proc.CombinedOutput(); err != nil {
			return fmt.Errorf("%s failed: %s: %w", cmd.Binary, out, err)
		}
	}
	return nil
}

func Unregister() error {
	for _, cmd := range BuildUnregisterCommands() {
		proc := exec.Command(cmd.Binary, cmd.Args...)
		if out, err := proc.CombinedOutput(); err != nil {
			return fmt.Errorf("%s failed: %s: %w", cmd.Binary, out, err)
		}
	}
	return nil
}
