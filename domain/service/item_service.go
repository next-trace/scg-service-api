// Package service contains domain services that implement business logic.
package service

import (
	"context"
	"fmt"

	"github.com/next-trace/scg-service-api/domain/entity"
	"github.com/next-trace/scg-service-api/domain/repository"
)

// ItemService provides business operations for items.
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService creates a new item service.
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

// GetItem retrieves an item by ID.
func (s *ItemService) GetItem(ctx context.Context, id string) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

// ListItems retrieves items based on filter criteria.
func (s *ItemService) ListItems(ctx context.Context, filter repository.ItemFilter) ([]*entity.Item, int64, error) {
	items, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list items: %w", err)
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count items: %w", err)
	}

	return items, count, nil
}

// CreateItem creates a new item.
func (s *ItemService) CreateItem(ctx context.Context, name, description string, tags []string) (*entity.Item, error) {
	item, err := entity.NewItem(name, description, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save item: %w", err)
	}

	return item, nil
}

// UpdateItem updates an existing item.
func (s *ItemService) UpdateItem(ctx context.Context, id, name, description string, tags []string, status entity.ItemStatus) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item for update: %w", err)
	}

	item.Update(name, description, tags, status)

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save updated item: %w", err)
	}

	return item, nil
}

// DeleteItem deletes an item by ID.
func (s *ItemService) DeleteItem(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("item ID cannot be empty")
	}

	// Option 1: Hard delete - remove from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	// Option 2: Soft delete - mark as deleted
	// item, err := s.repo.GetByID(ctx, id)
	// if err != nil {
	//     return fmt.Errorf("failed to get item for deletion: %w", err)
	// }
	//
	// item.Delete()
	//
	// if err := s.repo.Save(ctx, item); err != nil {
	//     return fmt.Errorf("failed to save deleted item: %w", err)
	// }

	return nil
}

// ActivateItem activates an inactive item.
func (s *ItemService) ActivateItem(ctx context.Context, id string) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item for activation: %w", err)
	}

	if item.IsActive() {
		return item, nil // Already active
	}

	item.Activate()

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save activated item: %w", err)
	}

	return item, nil
}

// DeactivateItem deactivates an active item.
func (s *ItemService) DeactivateItem(ctx context.Context, id string) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item for deactivation: %w", err)
	}

	if !item.IsActive() {
		return item, nil // Already inactive
	}

	item.Deactivate()

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save deactivated item: %w", err)
	}

	return item, nil
}

// AddTagToItem adds a tag to an item.
func (s *ItemService) AddTagToItem(ctx context.Context, id, tag string) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item for adding tag: %w", err)
	}

	item.AddTag(tag)

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save item with new tag: %w", err)
	}

	return item, nil
}

// RemoveTagFromItem removes a tag from an item.
func (s *ItemService) RemoveTagFromItem(ctx context.Context, id, tag string) (*entity.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}

	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item for removing tag: %w", err)
	}

	item.RemoveTag(tag)

	if err := s.repo.Save(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save item after removing tag: %w", err)
	}

	return item, nil
}
