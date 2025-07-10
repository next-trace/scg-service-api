// Package main provides tools for the SCG service base.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildTools builds the SCG tools and moves the binary to the project root.
func BuildTools() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Determine the path to the tools directory
	toolsDir := filepath.Join(currentDir, "cmd", "tools")

	// Change to the tools directory
	err = os.Chdir(toolsDir)
	if err != nil {
		fmt.Printf("Error changing to tools directory: %v\n", err)
		os.Exit(1)
	}

	// Build the tool
	fmt.Println("Building SCG tools...")

	// Find the go binary path
	goBinary, err := exec.LookPath("go")
	if err != nil {
		fmt.Printf("Error finding go binary: %v\n", err)
		os.Exit(1)
	}

	// Create command with fixed arguments
	cmd := &exec.Cmd{
		Path:   goBinary,
		Args:   []string{goBinary, "build", "-o", "scg-tools"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error building tool: %v\n", err)
		os.Exit(1)
	}

	// Move the binary to the project root
	binPath := filepath.Join(toolsDir, "scg-tools")
	destPath := filepath.Join(currentDir, "scg-tools")

	err = os.Rename(binPath, destPath)
	if err != nil {
		fmt.Printf("Error moving binary to project root: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SCG tools built successfully. Binary is at:", destPath)
	fmt.Println("\nUsage:")
	fmt.Println("  ./scg-tools -help           Show help message")
	fmt.Println("  ./scg-tools -proto          Generate Go code from protobuf definitions")
	fmt.Println("  ./scg-tools -clean          Clean generated files")
	fmt.Println("  ./scg-tools -install-tools  Install required tools")
}
