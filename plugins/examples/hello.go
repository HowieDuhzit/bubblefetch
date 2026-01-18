// Example plugin for bubblefetch
// Build with: go build -buildmode=plugin -o hello.so hello.go

package main

import (
	"fmt"

	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

// ModuleName is the name of this plugin module
// This must be a package-level variable named "ModuleName"
var ModuleName = "hello"

// Render is the function that renders this module's output
// It must have this exact signature
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	// Simple example: display a greeting with the hostname
	greeting := fmt.Sprintf("Hello from %s!", info.Hostname)

	// Use theme styles to render with proper colors
	label := styles.Label.Render("Greeting")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render(greeting)

	return label + separator + value
}
