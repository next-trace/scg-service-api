package entity_test

import (
	"testing"
	"time"

	entity "github.com/next-trace/scg-service-api/domain/entity"
)

func TestNewItemAndValidate(t *testing.T) {
	item, err := entity.NewItem("Widget", "A test item", []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID == "" {
		t.Fatalf("expected ID to be set")
	}
	if item.Name != "Widget" || item.Description != "A test item" {
		t.Fatalf("unexpected name/description: %#v", item)
	}
	if item.Status != entity.ItemStatusActive {
		t.Fatalf("expected status active, got %v", item.Status)
	}
	if got, want := len(item.Tags), 2; got != want {
		t.Fatalf("expected %d tags, got %d", want, got)
	}
	if item.CreatedAt.IsZero() || item.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps to be set")
	}

	// Validate should pass
	if err := item.Validate(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
}

func TestNewItem_EmptyName(t *testing.T) {
	if _, err := entity.NewItem("", "desc", nil); err == nil {
		t.Fatalf("expected error for empty name")
	}
}

func TestUpdateAndStatusTransitions(t *testing.T) {
	item, err := entity.NewItem("A", "B", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	before := item.UpdatedAt
	time.Sleep(10 * time.Millisecond) // ensure UpdatedAt will change
	item.Update("A2", "B2", []string{"x", "y"}, entity.ItemStatusInactive)

	if item.Name != "A2" || item.Description != "B2" {
		t.Fatalf("update failed for fields: %#v", item)
	}
	if item.Status != entity.ItemStatusInactive {
		t.Fatalf("expected status inactive, got %v", item.Status)
	}
	if got, want := len(item.Tags), 2; got != want {
		t.Fatalf("expected %d tags, got %d", want, got)
	}
	if !item.UpdatedAt.After(before) {
		t.Fatalf("expected UpdatedAt to move forward")
	}

	item.Activate()
	if !item.IsActive() {
		t.Fatalf("expected item active")
	}

	item.Deactivate()
	if item.IsActive() {
		t.Fatalf("expected item inactive")
	}

	item.Delete()
	if item.Status != entity.ItemStatusDeleted {
		t.Fatalf("expected deleted status")
	}
}

func TestTagsHelpers(t *testing.T) {
	item, err := entity.NewItem("T", "D", []string{"a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !item.HasTag("a") {
		t.Fatalf("expected to have tag 'a'")
	}
	if item.HasTag("b") {
		t.Fatalf("did not expect to have tag 'b'")
	}

	before := item.UpdatedAt
	item.AddTag("")       // no-op
	item.AddTag("a")      // duplicate, no-op
	item.AddTag("b")      // new
	item.RemoveTag("")    // no-op
	item.RemoveTag("zzz") // not present, no-op
	item.RemoveTag("a")   // remove existing

	if item.HasTag("a") {
		t.Fatalf("expected 'a' to be removed")
	}
	if !item.HasTag("b") {
		t.Fatalf("expected 'b' to exist")
	}
	if !item.UpdatedAt.After(before) {
		t.Fatalf("expected UpdatedAt to move forward after tag changes")
	}
}
