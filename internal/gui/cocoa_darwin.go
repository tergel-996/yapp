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

// activateHandler is the Go callback invoked when the user clicks Yapp's
// Dock icon while Yapp is running. launch.go installs a closure that
// runs `osascript -e 'tell application "<term>" to activate'` to bring
// the correct terminal window back to the foreground.
//
// It is set from a goroutine and read from the Cocoa main thread, so
// access is guarded by a mutex.
var (
	activateMu      sync.Mutex
	activateHandler func()
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
	f := activateHandler
	activateMu.Unlock()
	if f != nil {
		// Run on a goroutine so we don't block the Cocoa main thread
		// on exec.Command waiting for osascript. Cocoa delegate methods
		// should return quickly.
		go f()
	}
}
