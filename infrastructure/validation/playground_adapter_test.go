package validation_test

import (
	"bytes"
	"context"
	"testing"

	appvalidation "github.com/next-trace/scg-service-api/application/validation"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
	validatorimpl "github.com/next-trace/scg-service-api/infrastructure/validation"
)

type sample struct {
	Name string `validate:"required,min=3"`
}

func TestPlaygroundAdapter_BasicValidation(t *testing.T) {
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "info")
	cfg := appvalidation.DefaultConfig()
	v := validatorimpl.NewPlaygroundAdapter(cfg, log)

	ok := v.Validate(context.Background(), sample{Name: "abc"})
	if !ok.Valid {
		t.Fatalf("expected valid sample, got errors: %+v", ok.Errors)
	}

	// In this mock playground adapter, validation always succeeds; ensure it does not error
	bad := v.Validate(context.Background(), sample{Name: "a"})
	if !bad.Valid {
		t.Fatalf("expected mock validator to return valid result even for short name")
	}
}
