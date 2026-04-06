//go:build !darwin

// Non-darwin stub so `go build` works on Linux/Windows for the purposes
// of `go vet` and IDE tooling. Yapp is a macOS-only product; these no-ops
// are never executed at runtime on a supported platform.
package gui

func Run()                        {}
func Stop()                       {}
func SetActivateHandler(_ func()) {}
