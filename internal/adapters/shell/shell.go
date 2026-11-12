package shell

import (
"fmt"
"os"
"path/filepath"
"strings"
)

// Adapter handles shell-specific syntax for aliases and functions.
type Adapter struct {
	shellType string
}

// New creates a new shell adapter for the given shell type.
func New(shellType string) *Adapter {
	return &Adapter{shellType: shellType}
}

// FormatAlias formats a bookmark alias for the configured shell.
// A non-empty description is written as a trailing comment on the alias line
// (on the function line for fish). The description must already be escaped so
// it fits on a single line.
func (a *Adapter) FormatAlias(alias, command, description string) string {
	comment := ""
	if description != "" {
		comment = " # " + description
	}
	switch a.shellType {
	case "fish":
		return a.formatFishFunction(alias, command, comment)
	case "nu", "nushell":
		return a.formatNushellAlias(alias, command, comment)
	default: // bash, zsh, sh
		return a.formatPosixAlias(alias, command, comment)
	}
}

// formatPosixAlias formats an alias for POSIX-compatible shells (bash, zsh, sh).
func (a *Adapter) formatPosixAlias(alias, command, comment string) string {
	escaped := strings.ReplaceAll(command, "'", "'\\''")
	return fmt.Sprintf("alias %s='%s'%s\n\n", alias, escaped, comment)
}

// formatFishFunction formats a function for fish shell.
func (a *Adapter) formatFishFunction(alias, command, comment string) string {
	return fmt.Sprintf("function %s%s\n\t%s\nend\n\n", alias, comment, command)
}

// formatNushellAlias formats an alias for nushell.
func (a *Adapter) formatNushellAlias(alias, command, comment string) string {
	// Nushell requires spaces around = in alias
	return fmt.Sprintf("alias %s = %s%s\n\n", alias, command, comment)
}

// GetFileExtension returns the appropriate file extension for the shell.
func GetFileExtension(shellType string) string {
	switch shellType {
	case "fish":
		return ".fish"
	case "nu", "nushell":
		return ".nu"
	default: // bash, zsh, sh
		return ".sh"
	}
}

// DetectShell attempts to detect the current shell from environment variables.
// Returns the detected shell type or "bash" as a fallback.
func DetectShell() string {
	// Check SHELL environment variable (most common)
	if shellPath := os.Getenv("SHELL"); shellPath != "" {
		shellName := filepath.Base(shellPath)
		return normalizeShellName(shellName)
	}

	// Check for shell-specific environment variables
	if os.Getenv("BASH_VERSION") != "" {
		return "bash"
	}
	if os.Getenv("ZSH_VERSION") != "" {
		return "zsh"
	}
	if os.Getenv("FISH_VERSION") != "" {
		return "fish"
	}
	if os.Getenv("NU_VERSION") != "" {
		return "nu"
	}

	// Fallback to bash
	return "bash"
}

// normalizeShellName normalizes shell names to standard identifiers.
func normalizeShellName(name string) string {
	name = strings.ToLower(name)
	
	// Handle common variations
	switch {
	case strings.Contains(name, "bash"):
		return "bash"
	case strings.Contains(name, "zsh"):
		return "zsh"
	case strings.Contains(name, "fish"):
		return "fish"
	case strings.Contains(name, "nu"):
		return "nu"
	case name == "sh":
		return "sh"
	default:
		// If unknown, assume POSIX-compatible
		return "bash"
	}
}
