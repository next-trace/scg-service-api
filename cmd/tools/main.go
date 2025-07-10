package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	protoDir = `proto`
	genDir   = `gen`
	goOutDir = genDir + `/go/` + protoDir + `/scg`
)

func main() {
	// Define command line flags
	protoCmd := flag.Bool("proto", false, "Generate Go code from protobuf definitions")
	cleanCmd := flag.Bool("clean", false, "Clean generated files")
	installCmd := flag.Bool("install-tools", false, "Install required tools")
	buildCmd := flag.Bool("build", false, "Build the SCG tools")
	helpCmd := flag.Bool("help", false, "Show help message")

	// Parse command line flags
	flag.Parse()

	// If no flags are provided, show help
	if !*protoCmd && !*cleanCmd && !*installCmd && !*buildCmd && !*helpCmd {
		*helpCmd = true
	}

	// Execute commands based on flags
	if *helpCmd {
		showHelp()
		return
	}

	if *protoCmd {
		generateProto()
	}

	if *cleanCmd {
		cleanGenerated()
	}

	if *installCmd {
		installTools()
	}

	if *buildCmd {
		BuildTools()
	}
}

// showHelp displays the help message
func showHelp() {
	fmt.Println("SCG Service Base Tools")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  -proto          Generate Go code from protobuf definitions")
	fmt.Println("  -clean          Clean generated files")
	fmt.Println("  -install-tools  Install required tools")
	fmt.Println("  -build          Build the SCG tools")
	fmt.Println("  -help           Show this help message")
}

// generateProto generates Go code from protobuf definitions
func generateProto() {
	fmt.Println("Generating Go code from protobuf definitions...")

	// Ensure output directory exists
	err := os.MkdirAll(goOutDir, 0o750)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", goOutDir, err)
		os.Exit(1)
	}

	// Find all .proto files
	var protoFiles []string
	err = filepath.WalkDir(protoDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error finding proto files: %v\n", err)
		os.Exit(1)
	}

	// Process each proto file
	for _, protoFile := range protoFiles {
		// Validate proto file path for security
		if !isValidProtoFile(protoFile) {
			fmt.Printf("Skipping invalid proto file path: %s\n", protoFile)
			continue
		}

		// Clean and sanitize paths to prevent command injection
		cleanProtoDir := filepath.Clean(protoDir)
		cleanGenDir := filepath.Clean(genDir)
		cleanProtoFile := filepath.Clean(protoFile)

		// Use a fixed set of arguments with a hardcoded binary path to prevent command injection
		protoBinary, err := exec.LookPath("protoc")
		if err != nil {
			fmt.Printf("Error finding protoc binary: %v\n", err)
			os.Exit(1)
		}

		// Create a fixed set of arguments
		args := []string{
			"--proto_path", cleanProtoDir,
			"--go_out", cleanGenDir,
			"--go_opt", "paths=source_relative",
			"--go-grpc_out", cleanGenDir,
			"--go-grpc_opt", "paths=source_relative",
			cleanProtoFile,
		}

		cmd := &exec.Cmd{
			Path: protoBinary,
			Args: append([]string{protoBinary}, args...),
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error generating code for %s: %v\n%s\n", protoFile, err, output)
			os.Exit(1)
		}
	}

	fmt.Println("Done.")
}

// cleanGenerated removes generated files
func cleanGenerated() {
	fmt.Println("Cleaning generated files...")

	err := os.RemoveAll(genDir)
	if err != nil {
		fmt.Printf("Error removing directory %s: %v\n", genDir, err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}

// isValidProtoFile validates that a proto file path is safe to use
func isValidProtoFile(path string) bool {
	// Check if the file exists and is within the proto directory
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Get the absolute path of the proto directory
	absProtoDir, err := filepath.Abs(protoDir)
	if err != nil {
		return false
	}

	// Check if the file is within the proto directory
	if !strings.HasPrefix(absPath, absProtoDir) {
		return false
	}

	// Check if the file has a .proto extension
	if !strings.HasSuffix(path, ".proto") {
		return false
	}

	// Check if the file exists and is a regular file
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}

	return true
}

// installTools installs required tools
func installTools() {
	fmt.Println("Installing required tools...")

	// Find the go binary path
	goBinary, err := exec.LookPath("go")
	if err != nil {
		fmt.Printf("Error finding go binary: %v\n", err)
		os.Exit(1)
	}

	// Install protoc-gen-go with fixed arguments
	protocGenGoCmd := &exec.Cmd{
		Path: goBinary,
		Args: []string{goBinary, "install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"},
	}
	output, err := protocGenGoCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing protoc-gen-go: %v\n%s\n", err, output)
		os.Exit(1)
	}

	// Install protoc-gen-go-grpc with fixed arguments
	protocGenGrpcCmd := &exec.Cmd{
		Path: goBinary,
		Args: []string{goBinary, "install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"},
	}
	output, err = protocGenGrpcCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error installing protoc-gen-go-grpc: %v\n%s\n", err, output)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
