// Package testdata provides Go fixture files used by integration tests.
// The functions here are intentionally simple and contain no issues.
package testdata

// Add returns the sum of a and b.
func Add(a, b int) int { return a + b }

// Greet returns a greeting string.
func Greet(name string) string { return "Hello, " + name }
