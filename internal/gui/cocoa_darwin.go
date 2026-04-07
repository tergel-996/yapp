//go:build darwin

// Package gui hosts the Cocoa main-thread runloop that makes the Yapp
// binary appear to macOS as a real GUI application (Cmd+Tab visible,
// responds to Apple Events, no "not responding" dialog from
// LaunchServices). It is only meaningful in Yapp-bundle mode; the plain
// yapp-cli CLI path never calls Run.
package gui

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework AppKit -framework Foundation

// Declarations of the functions implemented in cocoa_darwin.m. Keeping
// the @implementation out of this preamble is important: cgo processes
// the preamble into a generated C file that can end up in the final
// link more than once, which produces duplicate Objective-C class
// symbols. A standalone .m file is compiled exactly once.
void yapp_run(void);
void yapp_stop(void);
*/
import "C"

import "sync"

// activateHandler is the Go callback invoked when Yapp is reactivated
// (Cmd+Tab switch, Dock icon click, or Finder re-launch of the already-
// running app). launch.go installs a closure that runs
// `osascript -e 'tell application "<term>" to activate'` to bring the
// terminal window rendering yazi back to the foreground.
//
// firstActivationSeen gates out the *initial* activation that Cocoa
// fires during `[NSApp activateIgnoringOtherApps:YES]` in yapp_run().
// That one happens before the terminal has been spawned, so running
// the handler at that moment would either no-op (terminal not yet
// running) or activate a pre-existing instance of the same terminal
// application (not the one that will shortly be hosting yazi). Skipping
// the first reactivation call avoids both of those wrong answers.
//
// activateHandler / firstActivationSeen are set from a worker goroutine
// and read from the Cocoa main thread, so access is guarded by a mutex.
var (
	activateMu          sync.Mutex
	activateHandler     func()
	firstActivationSeen bool
)

// Run blocks the calling goroutine in the Cocoa NSApplication event
// loop. It must be called from the main OS thread; callers should use
// runtime.LockOSThread before entering Run.
func Run() { C.yapp_run() }

// Stop asks the Cocoa event loop to exit. Safe to call from any
// goroutine; the work is dispatched to the main queue.
func Stop() { C.yapp_stop() }

// SetActivateHandler registers a function to be called when the user
// clicks Yapp's Dock icon while Yapp is already running. Passing nil
// clears the handler.
func SetActivateHandler(f func()) {
	activateMu.Lock()
	activateHandler = f
	activateMu.Unlock()
}

//export yappCocoaHandleReopen
func yappCocoaHandleReopen() {
	activateMu.Lock()
	if !firstActivationSeen {
		// Swallow the startup activation (see the comment on
		// firstActivationSeen). Subsequent reactivations — Cmd+Tab
		// switches, Dock-icon clicks, Finder re-launches — are real
		// user actions and must fire the handler.
		firstActivationSeen = true
		activateMu.Unlock()
		return
	}
	f := activateHandler
	activateMu.Unlock()
	if f != nil {
		// Run on a goroutine so we don't block the Cocoa main thread
		// on exec.Command waiting for osascript. Cocoa delegate methods
		// should return quickly.
		go f()
	}
}
