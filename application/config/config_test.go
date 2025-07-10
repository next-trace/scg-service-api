package config_test

import (
	"testing"

	appcfg "github.com/next-trace/scg-service-api/application/config"
)

func TestDefaultOptions(t *testing.T) {
	opt := appcfg.DefaultOptions()
	if opt.ConfigType != "yaml" {
		t.Fatalf("unexpected ConfigType: %s", opt.ConfigType)
	}
	if opt.EnvPrefix != "" {
		t.Fatalf("unexpected EnvPrefix: %s", opt.EnvPrefix)
	}
	if !opt.AllowEnvOverride {
		t.Fatalf("expected AllowEnvOverride=true by default")
	}
	if opt.RequireConfigFile {
		t.Fatalf("expected RequireConfigFile=false by default")
	}
}
