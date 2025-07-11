#!/bin/bash

# SCG Service Base Build Script

# Variables
PROTO_DIR="proto"
GEN_DIR="gen"
GO_OUT_DIR="${GEN_DIR}/go"

# Tools
PROTOC="protoc"
PROTOC_GEN_GO="protoc-gen-go"
PROTOC_GEN_GO_GRPC="protoc-gen-go-grpc"
GOLANGCI_LINT="golangci-lint"

# Function to show help message
show_help() {
  echo "Available commands:"
  echo "  all          - Generate all code (default)"
  echo "  proto        - Generate Go code from protobuf definitions"
  echo "  clean        - Clean generated files"
  echo "  lint         - Run linter on the codebase"
  echo "  lint-fix     - Run linter and fix issues automatically when possible"
  echo "  install-tools - Install required tools"
  echo "  help         - Show this help message"
}

# Function to generate Go code from protobuf definitions
generate_proto() {
  echo "Generating Go code from protobuf definitions..."
  # Ensure directories exist
  mkdir -p "${GO_OUT_DIR}"

  find "${PROTO_DIR}" -name "*.proto" -exec \
    ${PROTOC} \
    --proto_path=${PROTO_DIR} \
    --go_out=${GEN_DIR} \
    --go_opt=paths=source_relative \
    --go-grpc_out=${GEN_DIR} \
    --go-grpc_opt=paths=source_relative \
    {} \;
  echo "Done."
}

# Function to clean generated files
clean_generated() {
  echo "Cleaning generated files..."
  rm -rf "${GEN_DIR}"
  echo "Done."
}

# Function to install required tools
install_tools() {
  echo "Installing required tools..."
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

  # Install golangci-lint if not already installed
  if ! command -v ${GOLANGCI_LINT} &> /dev/null; then
    echo "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  fi

  echo "Done."
}

# Function to run the linter
run_lint() {
  echo "Running linter..."
  ${GOLANGCI_LINT} run ./...
  echo "Done."
}

# Function to run the linter and fix issues
run_lint_fix() {
  echo "Running linter and fixing issues..."
  ${GOLANGCI_LINT} run --fix ./...
  echo "Done."
}

# Main execution
case "$1" in
  "proto")
    generate_proto
    ;;
  "clean")
    clean_generated
    ;;
  "lint")
    run_lint
    ;;
  "lint-fix")
    run_lint_fix
    ;;
  "install-tools")
    install_tools
    ;;
  "help")
    show_help
    ;;
  "all" | "")
    generate_proto
    ;;
  *)
    echo "Unknown command: $1"
    show_help
    exit 1
    ;;
esac

exit 0
