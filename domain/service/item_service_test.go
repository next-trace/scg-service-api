package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/next-trace/scg-service-api/domain/entity"
	"github.com/next-trace/scg-service-api/domain/repository"
	servicepkg "github.com/next-trace/scg-service-api/domain/service"
)

type fakeRepo struct {
	items    map[string]*entity.Item
	saveN    int
	delN     int
	findN    int
	countN   int
	errGet   error
	errSave  error
	errDel   error
	findErr  error
	countErr error
}

func newFakeRepo() *fakeRepo { return &fakeRepo{items: map[string]*entity.Item{}} }

func (f *fakeRepo) GetByID(_ context.Context, id string) (*entity.Item, error) {
	if f.errGet != nil {
		return nil, f.errGet
	}
	it, ok := f.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return it, nil
}

func (f *fakeRepo) FindAll(_ context.Context, _ repository.ItemFilter) ([]*entity.Item, error) {
	f.findN++
	if f.findErr != nil {
		return nil, f.findErr
	}
	out := make([]*entity.Item, 0, len(f.items))
	for _, it := range f.items {
		out = append(out, it)
	}
	return out, nil
}

func (f *fakeRepo) Count(_ context.Context, _ repository.ItemFilter) (int64, error) {
	f.countN++
	if f.countErr != nil {
		return 0, f.countErr
	}
	return int64(len(f.items)), nil
}

func (f *fakeRepo) Save(_ context.Context, item *entity.Item) error {
	f.saveN++
	if f.errSave != nil {
		return f.errSave
	}
	f.items[item.ID] = item
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id string) error {
	f.delN++
	if f.errDel != nil {
		return f.errDel
	}
	delete(f.items, id)
	return nil
}

func TestItemService_BasicFlows(t *testing.T) {
	repo := newFakeRepo()
	s := servicepkg.NewItemService(repo)
	ctx := context.Background()

	// Create
	it, err := s.CreateItem(ctx, "n", "d", []string{"t"})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if it.ID == "" {
		t.Fatalf("expected id to be set")
	}
	if repo.saveN == 0 {
		t.Fatalf("expected Save to be called")
	}

	// Get
	got, err := s.GetItem(ctx, it.ID)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if got.ID != it.ID {
		t.Fatalf("unexpected item returned")
	}

	// List + Count
	items, total, err := s.ListItems(ctx, repository.NewItemFilter())
	if err != nil || len(items) == 0 || total == 0 {
		t.Fatalf("list/count failed: items=%d total=%d err=%v", len(items), total, err)
	}

	// Update
	updated, err := s.UpdateItem(ctx, it.ID, "n2", "d2", []string{"x"}, entity.ItemStatusInactive)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}
	if updated.Name != "n2" || updated.Status != entity.ItemStatusInactive {
		t.Fatalf("update did not apply: %#v", updated)
	}

	// Activate/Deactivate idempotency
	if _, err := s.ActivateItem(ctx, it.ID); err != nil {
		t.Fatalf("activate error: %v", err)
	}
	if _, err := s.ActivateItem(ctx, it.ID); err != nil { // already active
		t.Fatalf("activate idempotent error: %v", err)
	}
	if _, err := s.DeactivateItem(ctx, it.ID); err != nil {
		t.Fatalf("deactivate error: %v", err)
	}
	if _, err := s.DeactivateItem(ctx, it.ID); err != nil { // already inactive
		t.Fatalf("deactivate idempotent error: %v", err)
	}

	// Tag operations
	if _, err := s.AddTagToItem(ctx, it.ID, "tag1"); err != nil {
		t.Fatalf("add tag error: %v", err)
	}
	if _, err := s.RemoveTagFromItem(ctx, it.ID, "tag1"); err != nil {
		t.Fatalf("remove tag error: %v", err)
	}

	// Delete
	if err := s.DeleteItem(ctx, it.ID); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if repo.delN == 0 {
		t.Fatalf("expected Delete to be called")
	}
}

func TestItemService_ErrorsOnEmptyID(t *testing.T) {
	repo := newFakeRepo()
	s := servicepkg.NewItemService(repo)
	ctx := context.Background()

	if _, err := s.GetItem(ctx, ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.UpdateItem(ctx, "", "", "", nil, ""); err == nil {
		t.Fatalf("expected error")
	}
	if err := s.DeleteItem(ctx, ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.ActivateItem(ctx, ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.DeactivateItem(ctx, ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.AddTagToItem(ctx, "", "t"); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.AddTagToItem(ctx, "id", ""); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.RemoveTagFromItem(ctx, "", "t"); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := s.RemoveTagFromItem(ctx, "id", ""); err == nil {
		t.Fatalf("expected error")
	}
}
