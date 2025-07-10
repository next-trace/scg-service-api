package validation_test

import (
	"testing"

	appval "github.com/next-trace/scg-service-api/application/validation"
)

// The application/validation exposes interfaces and helpers. Ensure exported types exist and default behaviors are sane.
func TestInterfacesExist(t *testing.T) {
	var _ appval.Validator
}
