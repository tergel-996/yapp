package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yapp-cli",
		Short: "Yazi as a standalone macOS app",
		Long:  "Yapp wraps the yazi terminal file manager in its own macOS application identity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newLaunchCmd())
	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newUninstallCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newSetTerminalCmd())

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print yapp version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("yapp", Version)
		},
	}
}

// Execute runs the root Cobra command and returns its error (if any).
// Callers in cmd/yapp/main.go are responsible for translating the error
// into a process exit code. We do not call os.Exit here because when Yapp
// runs inside the macOS Cocoa runloop, main() needs to unwind cleanly so
// that the runloop can be stopped first.
func Execute() error {
	return NewRootCmd().Execute()
}
