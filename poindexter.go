// Package poindexter provides functionality for the Poindexter library.
package poindexter

// Version returns the current version of the library.
func Version() string {
	return "0.2.1"
}

// Hello returns a greeting message.
func Hello(name string) string {
	if name == "" {
		return "Hello, World!"
	}
	return "Hello, " + name + "!"
}
