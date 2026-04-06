package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tergel/yapp/internal/cli"
	"github.com/tergel/yapp/internal/gui"
)

func main() {
	// Cocoa requires its run loop to execute on the process's main OS
	// thread. LockOSThread pins the Go main goroutine to whatever OS
	// thread it started on -- which on macOS is thread 0 (the main
	// thread) -- so NSApplication is happy.
	runtime.LockOSThread()

	// When invoked as the .app bundle launcher (Contents/MacOS/Yapp),
	// run under a Cocoa event loop so macOS recognises us as a real
	// GUI app (Cmd+Tab, Dock icon, Apple Events).
	if filepath.Base(os.Args[0]) == "Yapp" {
		runAsBundle()
		return
	}

	// Plain CLI invocation: no Cocoa, no gui package initialisation
	// beyond the static linker cost.
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

// runAsBundle is the Yapp.app entry path. It:
//  1. Refreshes the bundle's embedded binary from Homebrew if newer.
//  2. Spawns a worker goroutine that runs the normal `launch` flow
//     (create marker + script file, spawn terminal, poll marker).
//  3. Blocks the main thread in the Cocoa runloop.
//  4. When the worker finishes (yazi exited), stops the runloop so the
//     process exits cleanly.
func runAsBundle() {
	selfHealBundle()

	// Rewrite argv so the existing cobra dispatch hits the launch
	// subcommand. Any positional argument (a path) is preserved.
	os.Args = append([]string{os.Args[0], "launch"}, os.Args[1:]...)

	launchErrCh := make(chan error, 1)
	go func() {
		launchErrCh <- cli.Execute()
		gui.Stop()
	}()

	// Blocks until gui.Stop is called from the worker goroutine above
	// or [NSApp terminate:] is invoked from the Quit menu.
	gui.Run()

	// If the runloop was ended by Cmd+Q / terminate: then Cocoa called
	// exit() internally and we don't get here. If it exited via stop:
	// we're here with the launch error in the channel (maybe).
	select {
	case err := <-launchErrCh:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Worker is still going (shouldn't happen in the normal path,
		// but don't block forever if it does).
	}
}

// selfHealBundle checks whether /opt/homebrew/bin/yapp-cli (or the Intel
// Homebrew equivalent) is newer than the binary currently running as
// Yapp.app/Contents/MacOS/Yapp. If so, it atomically replaces the bundle
// binary with the fresh one.
//
// All failures are intentionally silent -- this is best-effort maintenance
// on a hot launch path, not something the user triggered. Worst case, the
// user keeps running stale code until they run `yapp-cli install` manually.
func selfHealBundle() {
	self, err := os.Executable()
	if err != nil {
		return
	}

	var fresh string
	for _, candidate := range []string{
		"/opt/homebrew/bin/yapp-cli",
		"/usr/local/bin/yapp-cli",
	} {
		if _, err := os.Stat(candidate); err == nil {
			fresh = candidate
			break
		}
	}
	if fresh == "" {
		return
	}

	// Paranoia: never try to copy a file over itself. Can't happen under
	// normal circumstances (bundle path vs Homebrew bin path) but guards
	// against an unusual install layout.
	if selfReal, err := filepath.EvalSymlinks(self); err == nil {
		if freshReal, err := filepath.EvalSymlinks(fresh); err == nil {
			if selfReal == freshReal {
				return
			}
		}
	}

	selfInfo, err := os.Stat(self)
	if err != nil {
		return
	}
	freshInfo, err := os.Stat(fresh)
	if err != nil {
		return
	}
	if !freshInfo.ModTime().After(selfInfo.ModTime()) {
		return
	}

	data, err := os.ReadFile(fresh)
	if err != nil {
		return
	}
	tmp := self + ".new"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		_ = os.Remove(tmp)
		return
	}
	// On macOS, atomically replacing the path of a running Mach-O is fine:
	// the kernel keeps the old inode mapped for the current process while
	// future execs from the path see the new file.
	_ = os.Rename(tmp, self)
}
