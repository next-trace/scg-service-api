package pagination_test

import (
	"testing"

	app "github.com/next-trace/scg-service-api/application/pagination"
)

func TestDefaultPaginationOptions(t *testing.T) {
	opt := app.DefaultPaginationOptions()
	if opt.Page != 1 {
		t.Fatalf("unexpected Page: %d", opt.Page)
	}
	if opt.PageSize != app.DefaultPageSize {
		t.Fatalf("unexpected PageSize: %d", opt.PageSize)
	}
	if opt.BaseURL != "" {
		t.Fatalf("unexpected BaseURL: %s", opt.BaseURL)
	}
}
