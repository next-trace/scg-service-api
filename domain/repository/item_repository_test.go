package repository_test

import (
	"testing"

	"github.com/next-trace/scg-service-api/domain/entity"
	repo "github.com/next-trace/scg-service-api/domain/repository"
)

func TestItemFilterBuilders(t *testing.T) {
	f := repo.NewItemFilter()
	if f.Limit != 50 {
		t.Fatalf("expected default limit 50, got %d", f.Limit)
	}

	f = f.WithStatus(entity.ItemStatusInactive).WithTags([]string{"a", "b"}).WithSearch("q").WithPagination(10, 5)
	if f.Status != entity.ItemStatusInactive {
		t.Fatalf("status not set")
	}
	if len(f.Tags) != 2 || f.SearchTerm != "q" || f.Offset != 10 || f.Limit != 5 {
		t.Fatalf("unexpected filter: %#v", f)
	}

	// limit should not change when non-positive
	f2 := f.WithPagination(0, 0)
	if f2.Limit != f.Limit {
		t.Fatalf("expected limit unchanged when non-positive, got %d", f2.Limit)
	}
}
